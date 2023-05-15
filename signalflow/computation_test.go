package signalflow

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/signalfx/signalfx-go/idtool"
	"github.com/signalfx/signalfx-go/signalflow/v2/messages"
	"github.com/stretchr/testify/require"
)

func waitForDataMsg(t *testing.T, comp *Computation) (messages.Message, error) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	for {
		select {
		case m := <-comp.Data():
			if m == nil {
				continue
			}
			return m, nil
		case <-ctx.Done():
			err := comp.Err()
			if err != nil {
				return nil, err
			}

			t.Fatal("data message didn't get buffered")
		}
	}
}

func TestBuffersDataMessages(t *testing.T) {
	t.Parallel()
	ch := make(chan messages.Message)
	comp := newComputation(ch, "ch1", &Client{
		defaultMetadataTimeout: 1 * time.Second,
	})
	defer close(ch)
	ch <- &messages.DataMessage{
		Payloads: []messages.DataPayload{
			{
				TSID: idtool.ID(4000),
			},
		},
	}
	ch <- &messages.MetadataMessage{
		TSID: idtool.ID(4000),
	}

	md, _ := comp.TSIDMetadata(context.Background(), 4000)
	require.NotNil(t, md)

	ch <- &messages.InfoMessage{}

	msg, _ := waitForDataMsg(t, comp)
	require.Equal(t, idtool.ID(4000), msg.(*messages.DataMessage).Payloads[0].TSID)

	ch <- &messages.DataMessage{
		Payloads: []messages.DataPayload{
			{
				TSID: idtool.ID(4001),
			},
		},
	}
	msg, _ = waitForDataMsg(t, comp)
	require.Equal(t, idtool.ID(4001), msg.(*messages.DataMessage).Payloads[0].TSID)
}

func waitForExpiryMsg(t *testing.T, comp *Computation) (messages.Message, error) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	for {
		select {
		case m := <-comp.Expirations():
			if m == nil {
				continue
			}
			return m, nil
		case <-ctx.Done():
			err := comp.Err()
			if err != nil {
				return nil, err
			}

			t.Fatal("data message didn't get buffered")
		}
	}
}

func TestBuffersExpiryMessages(t *testing.T) {
	t.Parallel()
	ch := make(chan messages.Message)
	comp := newComputation(ch, "ch1", &Client{
		defaultMetadataTimeout: 1 * time.Second,
	})
	defer close(ch)
	ch <- &messages.ExpiredTSIDMessage{
		TSID: idtool.ID(4000).String(),
	}
	ch <- &messages.MetadataMessage{
		TSID: idtool.ID(4000),
	}

	md, _ := comp.TSIDMetadata(context.Background(), 4000)
	require.NotNil(t, md)

	ch <- &messages.InfoMessage{}

	msg, _ := waitForExpiryMsg(t, comp)
	require.Equal(t, idtool.ID(4000).String(), msg.(*messages.ExpiredTSIDMessage).TSID)

	ch <- &messages.ExpiredTSIDMessage{
		TSID: idtool.ID(4001).String(),
	}
	msg, _ = waitForExpiryMsg(t, comp)
	require.Equal(t, idtool.ID(4001).String(), msg.(*messages.ExpiredTSIDMessage).TSID)
}

func mustParse(m messages.Message, err error) messages.Message {
	if err != nil {
		panic(err)
	}
	return m
}

func TestResolutionMetadata(t *testing.T) {
	t.Parallel()
	ch := make(chan messages.Message)
	comp := newComputation(ch, "ch1", &Client{
		defaultMetadataTimeout: 1 * time.Second,
	})
	defer close(ch)

	wg := sync.WaitGroup{}

	// Ensure multiple calls get the same result and also wait for the message
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			resolution, _ := comp.Resolution(context.Background())
			require.Equal(t, 5*time.Second, resolution)
		}()
	}

	ch <- mustParse(messages.ParseMessage([]byte(`{
		"type": "message",
		"message": {
			"messageCode": "JOB_RUNNING_RESOLUTION",
			"contents": {
				"resolutionMs": 5000
			}
		}
	}`), true))

	wg.Wait()
}

func TestMaxDelayMetadata(t *testing.T) {
	t.Parallel()
	ch := make(chan messages.Message)
	comp := newComputation(ch, "ch1", &Client{
		defaultMetadataTimeout: 1 * time.Second,
	})
	defer close(ch)
	ch <- mustParse(messages.ParseMessage([]byte(`{
		"type": "message",
		"message": {
			"messageCode": "JOB_INITIAL_MAX_DELAY",
			"contents": {
				"maxDelayMs": 1000
			}
		}
	}`), true))

	maxDelay, _ := comp.MaxDelay(context.Background())
	require.Equal(t, 1*time.Second, maxDelay)
}

func TestLagMetadata(t *testing.T) {
	t.Parallel()
	ch := make(chan messages.Message)
	comp := newComputation(ch, "ch1", &Client{
		defaultMetadataTimeout: 1 * time.Second,
	})
	defer close(ch)
	ch <- mustParse(messages.ParseMessage([]byte(`{
		"type": "message",
		"message": {
			"messageCode": "JOB_DETECTED_LAG",
			"contents": {
				"lagMs": 3500
			}
		}
	}`), true))

	lag, _ := comp.Lag(context.Background())
	require.Equal(t, 3500*time.Millisecond, lag)
}

func TestFindLimitedResultSetMetadata(t *testing.T) {
	t.Parallel()
	ch := make(chan messages.Message)
	comp := newComputation(ch, "ch1", &Client{
		defaultMetadataTimeout: 1 * time.Second,
	})
	defer close(ch)
	ch <- mustParse(messages.ParseMessage([]byte(`{
		"type": "message",
		"message": {
			"messageCode": "FIND_LIMITED_RESULT_SET",
			"contents": {
				"matchedSize": 123456789,
				"limitSize": 50000
			}
		}
	}`), true))

	matchedSize, _ := comp.MatchedSize(context.Background())
	require.Equal(t, 123456789, matchedSize)

	limitSize, _ := comp.LimitSize(context.Background())
	require.Equal(t, 50000, limitSize)
}

func TestMatchedNoTimeseriesQueryMetaData(t *testing.T) {
	t.Parallel()
	ch := make(chan messages.Message)
	comp := newComputation(ch, "ch1", &Client{
		defaultMetadataTimeout: 1 * time.Second,
	})
	defer close(ch)
	ch <- mustParse(messages.ParseMessage([]byte(`{
		"type": "message",
		"message": {
			"messageCode": "FIND_MATCHED_NO_TIMESERIES",
			"contents": {
				"query": "abc"
			}
		}
	}`), true))

	noMatched, _ := comp.MatchedNoTimeseriesQuery(context.Background())
	require.Equal(t, "abc", noMatched)
}

func TestGroupByMissingPropertyMetaData(t *testing.T) {
	t.Parallel()
	ch := make(chan messages.Message)
	comp := newComputation(ch, "ch1", &Client{
		defaultMetadataTimeout: 1 * time.Second,
	})
	defer close(ch)
	ch <- mustParse(messages.ParseMessage([]byte(`{
		"type": "message",
		"message": {
			"messageCode": "GROUPBY_MISSING_PROPERTY",
			"contents": {
				"propertyNames": ["x", "y", "z"]
			}
		}
	}`), true))

	missingProps, _ := comp.GroupByMissingProperties(context.Background())
	require.Equal(t, []string{"x", "y", "z"}, missingProps)
}

func TestHandle(t *testing.T) {
	t.Parallel()
	ch := make(chan messages.Message)
	comp := newComputation(ch, "ch1", &Client{
		defaultMetadataTimeout: 1 * time.Second,
	})
	defer close(ch)
	ch <- mustParse(messages.ParseMessage([]byte(`{
		"type": "control-message",
		"event": "JOB_START",
		"handle": "AAAABBBB"
	}`), true))

	handle, _ := comp.Handle(context.Background())
	require.Equal(t, "AAAABBBB", handle)
}

func TestComputationError(t *testing.T) {
	t.Parallel()
	ch := make(chan messages.Message)
	comp := newComputation(ch, "ch1", &Client{
		defaultMetadataTimeout: 1 * time.Second,
	})
	defer close(ch)
	ch <- mustParse(messages.ParseMessage([]byte(`{
		"type": "error",
		"error": 400,
		"errorType": "ANALYTICS_PROGRAM_NAME_ERROR",
		"message": "We hit some error"
	}`), true))

	_, err := waitForDataMsg(t, comp)
	if err == nil {
		t.Fatal("Expected computation error")
	}

	var ce *ComputationError
	if !errors.As(err, &ce) {
		t.FailNow()
	}
	require.Equal(t, 400, ce.Code)
	require.Equal(t, "ANALYTICS_PROGRAM_NAME_ERROR", ce.ErrorType)
	require.Equal(t, "We hit some error", ce.Message)
}

func TestComputationErrorWithNullMessage(t *testing.T) {
	t.Parallel()
	ch := make(chan messages.Message)
	comp := newComputation(ch, "ch1", &Client{
		defaultMetadataTimeout: 1 * time.Second,
	})
	defer close(ch)
	ch <- mustParse(messages.ParseMessage([]byte(`{
		"type": "error",
		"error": 400,
		"errorType": "ANALYTICS_INTERNAL_ERROR",
		"message": null
	}`), true))

	_, err := waitForDataMsg(t, comp)
	if err == nil {
		t.Fatal("Expected computation error")
	}

	var ce *ComputationError
	if !errors.As(err, &ce) {
		t.FailNow()
	}
	require.Equal(t, 400, ce.Code)
	require.Equal(t, "ANALYTICS_INTERNAL_ERROR", ce.ErrorType)
	require.Equal(t, "", ce.Message)
}

func TestComputationFinish(t *testing.T) {
	t.Parallel()
	ch := make(chan messages.Message)
	comp := newComputation(ch, "ch1", &Client{
		defaultMetadataTimeout: 1 * time.Second,
	})
	defer close(ch)
	go func() {
		ch <- mustParse(messages.ParseMessage([]byte(`{
			"type": "control-message",
			"event": "JOB_START",
			"handle": "AAAABBBB"
		}`), true))

		ch <- &messages.MetadataMessage{
			TSID: idtool.ID(4000),
		}

		ch <- &messages.DataMessage{
			Payloads: []messages.DataPayload{
				{
					TSID: idtool.ID(4000),
				},
			},
		}

		ch <- &messages.EndOfChannelControlMessage{}
	}()

	for msg := range comp.Data() {
		require.Equal(t, idtool.ID(4000), msg.Payloads[0].TSID)
	}

	// The for loop should exit when the end of channel message comes through
}
