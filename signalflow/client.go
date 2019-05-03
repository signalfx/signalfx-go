package signalflow

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/url"
	"path"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
)

// Client for SignalFlow via websockets (SSE is not currently supported).
type Client struct {
	// Access token for the org
	token            string
	nextChannelNum   int64
	conn             *websocket.Conn
	outgoingMessages chan interface{}
	// How long to wait for writes to the websocket to finish
	writeTimeout   time.Duration
	streamURL      *url.URL
	channelsByName map[string]*Channel

	isShutdown bool
	ctx        context.Context
	cancel     context.CancelFunc
	lock       sync.Mutex
}

type ClientParam func(*Client) error

// StreamURL lets you set the full URL to the stream endpoint, including the
// path.
func StreamURL(streamEndpoint string) ClientParam {
	return func(c *Client) error {
		var err error
		c.streamURL, err = url.Parse(streamEndpoint)
		return err
	}
}

func StreamURLFromRealm(realm string) ClientParam {
	return func(c *Client) error {
		var err error
		c.streamURL, err = url.Parse(fmt.Sprintf("wss://stream.%s.signalfx.com/v2/signalflow", realm))
		return err
	}
}

func AccessToken(token string) ClientParam {
	return func(c *Client) error {
		c.token = token
		return nil
	}
}

func NewClient(options ...ClientParam) (*Client, error) {
	c := &Client{
		streamURL: &url.URL{
			Scheme: "wss",
			Host:   "stream.us0.signalfx.com",
			Path:   "/v2/signalflow",
		},
		writeTimeout:     5 * time.Second,
		outgoingMessages: make(chan interface{}),
		channelsByName:   make(map[string]*Channel),
	}

	for i := range options {
		if err := options[i](c); err != nil {
			return nil, err
		}
	}

	return c, nil
}

func (c *Client) initialize() error {
	c.lock.Lock()
	defer c.lock.Unlock()

	authenticatedCond := sync.NewCond(&c.lock)

	if c.isShutdown {
		return errors.New("cannot initialize client after shutdown")
	}

	c.ctx, c.cancel = context.WithCancel(context.Background())

	if c.conn == nil {
		if err := c.connect(); err != nil {
			return err
		}

		go c.keepWritingMessages()
		go c.keepReadingMessages(authenticatedCond)
		c.authenticate()
		authenticatedCond.Wait()
	}
	return nil
}

func (c *Client) newUniqueChannelName() string {
	name := fmt.Sprintf("ch-%d", c.nextChannelNum)
	atomic.AddInt64(&c.nextChannelNum, 1)
	return name
}

// Writes all messages from a single goroutine since that is required by
// websocket library.
func (c *Client) keepWritingMessages() {
	for {
		select {
		case <-c.ctx.Done():
			return
		case message := <-c.outgoingMessages:
			err := c.conn.SetWriteDeadline(time.Now().Add(c.writeTimeout))
			if err != nil {
				log.Printf("Error setting write timeout for SignalFlow request: %v", err)
				continue
			}

			msgBytes, err := json.Marshal(message)
			if err != nil {
				log.Printf("Error marshaling SignalFlow request: %v", err)
				continue
			}

			err = c.conn.WriteMessage(websocket.TextMessage, msgBytes)
			if err != nil {
				log.Printf("Error writing SignalFlow request: %v", err)
				continue
			}
		}
	}
}

// Reads all messages from a single goroutine and distributes them where
// needed.
func (c *Client) keepReadingMessages(authenticatedCond *sync.Cond) {
	for {
		msgTyp, msgBytes, err := c.conn.ReadMessage()
		if err != nil {
			// this means we are shutdown
			if c.ctx.Err() != nil {
				return
			}
			c.broadcastError(err)
			continue
		}

		message, err := parseMessage(msgBytes, msgTyp == websocket.TextMessage)
		if err != nil {
			log.Printf("Error parsing SignalFlow message: %v", err)
			continue
		}

		if cm, ok := message.(ChannelMessage); ok {
			channelName := cm.Channel()
			channel, ok := c.channelsByName[channelName]
			if !ok || channelName == "" {
				log.Printf("SignalFlow message received for unknown channel: %s", channelName)
				continue
			}
			channel.AcceptMessage(cm)
		} else {
			c.acceptMessage(message, authenticatedCond)
		}
	}
}

func (c *Client) acceptMessage(message Message, authenticatedCond *sync.Cond) {
	if _, ok := message.(*AuthenticatedMessage); ok {
		authenticatedCond.Signal()
	}
}

func (c *Client) connect() error {
	connectURL := *c.streamURL
	connectURL.Path = path.Join(c.streamURL.Path, "connect")
	conn, _, err := websocket.DefaultDialer.DialContext(c.ctx, connectURL.String(), nil)
	if err != nil {
		return fmt.Errorf("could not connect Signalflow websocket: %v", err)
	}
	c.conn = conn
	return nil
}

func (c *Client) authenticate() {
	c.outgoingMessages <- &AuthRequest{
		Token: c.token,
	}
}

func (c *Client) broadcastError(err error) {
	for _, ch := range c.channelsByName {
		ch.AcceptMessage(&WebsocketErrorMessage{
			Error: err,
		})
	}
}

// Execute a SignalFlow job and return a channel upon which informational
// messages and data will flow.
func (c *Client) Execute(req *ExecuteRequest) (*Channel, error) {
	if err := c.initialize(); err != nil {
		return nil, err
	}

	if req.Channel == "" {
		req.Channel = c.newUniqueChannelName()
	}

	c.outgoingMessages <- req

	return c.registerChannel(req.Channel), nil
}

func (c *Client) registerChannel(name string) *Channel {
	ch := newChannel(name)

	c.channelsByName[name] = ch
	return ch
}

// Close the client and shutdown any ongoing connections and goroutines.
func (c *Client) Close() {
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.cancel != nil {
		c.cancel()
	}
	if c.conn != nil {
		c.conn.Close()
	}
	for _, ch := range c.channelsByName {
		ch.Close()
	}
	c.isShutdown = true
}

// The way to distinguish between JSON and binary messages is the websocket
// message type.
func parseMessage(msg []byte, isText bool) (Message, error) {
	if isText {
		var baseMessage BaseMessage
		if err := json.Unmarshal(msg, &baseMessage); err != nil {
			return nil, fmt.Errorf("couldn't unmarshal JSON websocket message: %v", err)
		}
		return parseJSONMessage(&baseMessage, msg)
	}
	return parseBinaryMessage(msg)
}
