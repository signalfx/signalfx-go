package signalflow

type Channel struct {
	name     string
	messages chan Message
	// Channel that is always read from to pull messages out and buffer them
	// without blocking, even if nobody is listening on `messages`.
	incoming chan Message
}

func newChannel(name string) *Channel {
	c := &Channel{
		name:     name,
		messages: make(chan Message),
		incoming: make(chan Message, 10),
	}
	go c.processBufferedMessages()
	return c
}

// Buffer up messages indefinitely until another goroutine reads them off of
// c.messages, which is an unbuffered channel.
func (c *Channel) processBufferedMessages() {
	buffer := make([]Message, 0)
	for {
		if len(buffer) > 0 {
			var nextMessage Message
			nextMessage, buffer = buffer[0], buffer[1:]
			select {
			case c.messages <- nextMessage:
				continue
			case msg := <-c.incoming:
				buffer = append(buffer, msg)
			}
		} else {
			buffer = append(buffer, <-c.incoming)
		}
	}
}

func (c *Channel) AcceptMessage(msg Message) {
	// This shouldn't block since it should always be read from
	c.incoming <- msg
}

func (c *Channel) Messages() <-chan Message {
	return c.messages
}

func (c *Channel) Close() {
	close(c.messages)
}
