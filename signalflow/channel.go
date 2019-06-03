package signalflow

import "github.com/signalfx/signalfx-go/signalflow/messages"

// Channel is a queue of messages that all pertain to the same computation.
type Channel struct {
	name     string
	messages chan messages.Message
}

func newChannel(name string) *Channel {
	c := &Channel{
		name:     name,
		messages: make(chan messages.Message),
	}
	return c
}

// AcceptMessage from a websocket.  This might block if nothing is reading from
// the channel but generally a compuatation should always be doing so.
func (c *Channel) AcceptMessage(msg messages.Message) {
	c.messages <- msg
}

// Messages returns a Go chan that will be pushed all of the deserialized
// SignalFlow messages from the websocket.
func (c *Channel) Messages() <-chan messages.Message {
	return c.messages
}

// Close the channel.  This does not actually stop a job in SignalFlow, for
// that use Computation.Stop().
func (c *Channel) Close() {
	close(c.messages)
}
