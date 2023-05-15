package signalflow

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/signalfx/signalfx-go/idtool"
	"github.com/signalfx/signalfx-go/signalflow/v2/messages"
)

// Computation is a single running SignalFlow job
type Computation struct {
	sync.Mutex
	channel <-chan messages.Message
	name    string
	client  *Client
	dataCh  chan *messages.DataMessage
	// An intermediate channel for data messages where they can be buffered if
	// nothing is currently pulling data messages.
	dataChBuffer       chan *messages.DataMessage
	expirationCh       chan *messages.ExpiredTSIDMessage
	expirationChBuffer chan *messages.ExpiredTSIDMessage

	errMutex  sync.RWMutex
	lastError error

	handle                   asyncMetadata[string]
	resolutionMS             asyncMetadata[int]
	lagMS                    asyncMetadata[int]
	maxDelayMS               asyncMetadata[int]
	matchedSize              asyncMetadata[int]
	limitSize                asyncMetadata[int]
	matchedNoTimeseriesQuery asyncMetadata[string]
	groupByMissingProperties asyncMetadata[[]string]

	tsidMetadata map[idtool.ID]*asyncMetadata[*messages.MetadataProperties]
}

// ComputationError exposes the underlying metadata of a computation error
type ComputationError struct {
	Code      int
	Message   string
	ErrorType string
}

func (e *ComputationError) Error() string {
	err := fmt.Sprintf("%v", e.Code)
	if e.ErrorType != "" {
		err = fmt.Sprintf("%v (%v)", e.Code, e.ErrorType)
	}
	if e.Message != "" {
		err = fmt.Sprintf("%v: %v", err, e.Message)
	}
	return err
}

func newComputation(channel <-chan messages.Message, name string, client *Client) *Computation {
	comp := &Computation{
		channel:            channel,
		name:               name,
		client:             client,
		dataCh:             make(chan *messages.DataMessage),
		dataChBuffer:       make(chan *messages.DataMessage),
		expirationCh:       make(chan *messages.ExpiredTSIDMessage),
		expirationChBuffer: make(chan *messages.ExpiredTSIDMessage),
		tsidMetadata:       make(map[idtool.ID]*asyncMetadata[*messages.MetadataProperties]),
	}

	go comp.bufferDataMessages()
	go comp.bufferExpirationMessages()
	go func() {
		err := comp.watchMessages()

		if !errors.Is(err, errChannelClosed) {
			comp.errMutex.Lock()
			comp.lastError = err
			comp.errMutex.Unlock()
		}

		comp.shutdown()
	}()

	return comp
}

// Handle of the computation. Will wait as long as the given ctx is not closed. If ctx is closed an
// error will be returned.
func (c *Computation) Handle(ctx context.Context) (string, error) {
	return c.handle.Get(ctx)
}

// Resolution of the job. Will wait as long as the given ctx is not closed. If ctx is closed an
// error will be returned.
func (c *Computation) Resolution(ctx context.Context) (time.Duration, error) {
	resMS, err := c.resolutionMS.Get(ctx)
	return time.Duration(resMS) * time.Millisecond, err
}

// Lag detected for the job. Will wait as long as the given ctx is not closed. If ctx is closed an
// error will be returned.
func (c *Computation) Lag(ctx context.Context) (time.Duration, error) {
	lagMS, err := c.lagMS.Get(ctx)
	return time.Duration(lagMS) * time.Millisecond, err
}

// MaxDelay detected of the job. Will wait as long as the given ctx is not closed. If ctx is closed an
// error will be returned.
func (c *Computation) MaxDelay(ctx context.Context) (time.Duration, error) {
	maxDelayMS, err := c.maxDelayMS.Get(ctx)
	return time.Duration(maxDelayMS) * time.Millisecond, err
}

// MatchedSize detected of the job. Will wait as long as the given ctx is not closed. If ctx is closed an
// error will be returned.
func (c *Computation) MatchedSize(ctx context.Context) (int, error) {
	return c.matchedSize.Get(ctx)
}

// LimitSize detected of the job. Will wait as long as the given ctx is not closed. If ctx is closed an
// error will be returned.
func (c *Computation) LimitSize(ctx context.Context) (int, error) {
	return c.limitSize.Get(ctx)
}

// MatchedNoTimeseriesQuery if it matched no active timeseries. Will wait as long as the given ctx
// is not closed. If ctx is closed an error will be returned.
func (c *Computation) MatchedNoTimeseriesQuery(ctx context.Context) (string, error) {
	return c.matchedNoTimeseriesQuery.Get(ctx)
}

// GroupByMissingProperties are timeseries that don't contain the required dimensions. Will wait as
// long as the given ctx is not closed. If ctx is closed an error will be returned.
func (c *Computation) GroupByMissingProperties(ctx context.Context) ([]string, error) {
	return c.groupByMissingProperties.Get(ctx)
}

// TSIDMetadata for a particular tsid. Will wait as long as the given ctx is not closed. If ctx is closed an
// error will be returned.
func (c *Computation) TSIDMetadata(ctx context.Context, tsid idtool.ID) (*messages.MetadataProperties, error) {
	c.Lock()
	if _, ok := c.tsidMetadata[tsid]; !ok {
		c.tsidMetadata[tsid] = &asyncMetadata[*messages.MetadataProperties]{}
	}
	md := c.tsidMetadata[tsid]
	c.Unlock()
	return md.Get(ctx)
}

// Err returns the last fatal error that caused the computation to stop, if
// any.  Will be nil if the computation stopped in an expected manner.
func (c *Computation) Err() error {
	c.errMutex.RLock()
	defer c.errMutex.RUnlock()

	return c.lastError
}

func (c *Computation) watchMessages() error {
	for {
		m, ok := <-c.channel
		if !ok {
			return nil
		}
		if err := c.processMessage(m); err != nil {
			return err
		}
	}
}

var errChannelClosed = errors.New("computation channel is closed")

func (c *Computation) processMessage(m messages.Message) error {
	switch v := m.(type) {
	case *messages.JobStartControlMessage:
		c.handle.Set(v.Handle)
	case *messages.EndOfChannelControlMessage, *messages.ChannelAbortControlMessage:
		return errChannelClosed
	case *messages.DataMessage:
		c.dataChBuffer <- v
	case *messages.ExpiredTSIDMessage:
		c.Lock()
		delete(c.tsidMetadata, idtool.IDFromString(v.TSID))
		c.Unlock()
		c.expirationChBuffer <- v
	case *messages.InfoMessage:
		switch v.MessageBlock.Code {
		case messages.JobRunningResolution:
			c.resolutionMS.Set(v.MessageBlock.Contents.(messages.JobRunningResolutionContents).ResolutionMS())
		case messages.JobDetectedLag:
			c.lagMS.Set(v.MessageBlock.Contents.(messages.JobDetectedLagContents).LagMS())
		case messages.JobInitialMaxDelay:
			c.maxDelayMS.Set(v.MessageBlock.Contents.(messages.JobInitialMaxDelayContents).MaxDelayMS())
		case messages.FindLimitedResultSet:
			c.matchedSize.Set(v.MessageBlock.Contents.(messages.FindLimitedResultSetContents).MatchedSize())
			c.limitSize.Set(v.MessageBlock.Contents.(messages.FindLimitedResultSetContents).LimitSize())
		case messages.FindMatchedNoTimeseries:
			c.matchedNoTimeseriesQuery.Set(v.MessageBlock.Contents.(messages.FindMatchedNoTimeseriesContents).MatchedNoTimeseriesQuery())
		case messages.GroupByMissingProperty:
			c.groupByMissingProperties.Set(v.MessageBlock.Contents.(messages.GroupByMissingPropertyContents).GroupByMissingProperties())
		}
	case *messages.ErrorMessage:
		rawData := v.RawData()
		computationError := ComputationError{}
		if code, ok := rawData["error"]; ok {
			computationError.Code = int(code.(float64))
		}
		if msg, ok := rawData["message"]; ok && msg != nil {
			computationError.Message = msg.(string)
		}
		if errType, ok := rawData["errorType"]; ok {
			computationError.ErrorType = errType.(string)
		}
		return &computationError
	case *messages.MetadataMessage:
		c.Lock()
		if _, ok := c.tsidMetadata[v.TSID]; !ok {
			c.tsidMetadata[v.TSID] = &asyncMetadata[*messages.MetadataProperties]{}
		}
		c.tsidMetadata[v.TSID].Set(&v.Properties)
		c.Unlock()
	}
	return nil
}

// Buffer up data messages indefinitely until another goroutine reads them off of c.messages, which
// is an unbuffered channel. They need to be buffered because metadata messages can come _after_
// data messages.
func (c *Computation) bufferDataMessages() {
	buffer := make([]*messages.DataMessage, 0)
	var nextMessage *messages.DataMessage

	defer func() {
		if nextMessage != nil {
			c.dataCh <- nextMessage
		}
		for i := range buffer {
			c.dataCh <- buffer[i]
		}

		select {
		case msg, ok := <-c.dataChBuffer:
			if ok {
				c.dataCh <- msg
			}
		default:
		}

		close(c.dataCh)
	}()

	for {
		if len(buffer) > 0 {
			if nextMessage == nil {
				nextMessage, buffer = buffer[0], buffer[1:]
			}
			select {
			case c.dataCh <- nextMessage:
				nextMessage = nil
			case msg, ok := <-c.dataChBuffer:
				if !ok {
					return
				}
				buffer = append(buffer, msg)
			}
		} else {
			msg, ok := <-c.dataChBuffer
			if !ok {
				return
			}
			buffer = append(buffer, msg)
		}
	}
}

// Buffer up expiration messages indefinitely until another goroutine reads
// them off of c.expirationCh, which is an unbuffered channel.
func (c *Computation) bufferExpirationMessages() {
	buffer := make([]*messages.ExpiredTSIDMessage, 0)
	var nextMessage *messages.ExpiredTSIDMessage

	defer func() {
		if nextMessage != nil {
			c.expirationCh <- nextMessage
		}
		for i := range buffer {
			c.expirationCh <- buffer[i]
		}

		close(c.expirationCh)
	}()
	for {
		if len(buffer) > 0 {
			if nextMessage == nil {
				nextMessage, buffer = buffer[0], buffer[1:]
			}

			select {
			case c.expirationCh <- nextMessage:
				nextMessage = nil
			case msg, ok := <-c.expirationChBuffer:
				if !ok {
					return
				}
				buffer = append(buffer, msg)
			}
		} else {
			msg, ok := <-c.expirationChBuffer
			if !ok {
				return
			}
			buffer = append(buffer, msg)
		}
	}
}

// Data returns the channel on which data messages come.  This channel will be closed when the
// computation is finished.  To prevent goroutine leaks, you should read all messages from this
// channel until it is closed.
func (c *Computation) Data() <-chan *messages.DataMessage {
	return c.dataCh
}

// Expirations returns a channel that will be sent messages about expired TSIDs, i.e. time series
// that are no longer valid for this computation. This channel will be closed when the computation
// is finished. To prevent goroutine leaks, you should read all messages from this channel until it
// is closed.
func (c *Computation) Expirations() <-chan *messages.ExpiredTSIDMessage {
	return c.expirationCh
}

// Detach the computation on the backend
func (c *Computation) Detach(ctx context.Context) error {
	return c.DetachWithReason(ctx, "")
}

// DetachWithReason detaches the computation with a given reason. This reason will
// be reflected in the control message that signals the end of the job/channel
func (c *Computation) DetachWithReason(ctx context.Context, reason string) error {
	return c.client.Detach(ctx, &DetachRequest{
		Reason:  reason,
		Channel: c.name,
	})
}

// Stop the computation on the backend.
func (c *Computation) Stop(ctx context.Context) error {
	return c.StopWithReason(ctx, "")
}

// StopWithReason stops the computation with a given reason. This reason will
// be reflected in the control message that signals the end of the job/channel.
func (c *Computation) StopWithReason(ctx context.Context, reason string) error {
	handle, err := c.handle.Get(ctx)
	if err != nil {
		return err
	}
	return c.client.Stop(ctx, &StopRequest{
		Reason: reason,
		Handle: handle,
	})
}

func (c *Computation) shutdown() {
	close(c.dataChBuffer)
	close(c.expirationChBuffer)
}

var ErrMetadataTimeout = errors.New("metadata value did not come in time")

type asyncMetadata[T any] struct {
	sync.Mutex
	sig chan struct{}
	val T
}

func (a *asyncMetadata[T]) ensureInit() {
	a.Lock()
	if a.sig == nil {
		a.sig = make(chan struct{})
	}
	a.Unlock()
}

func (a *asyncMetadata[T]) Set(val T) {
	a.ensureInit()
	a.Lock()
	a.val = val
	close(a.sig)
	a.Unlock()
}

func (a *asyncMetadata[T]) Get(ctx context.Context) (T, error) {
	a.ensureInit()
	select {
	case <-ctx.Done():
		var t T
		return t, ErrMetadataTimeout
	case <-a.sig:
		return a.val, nil
	}
}
