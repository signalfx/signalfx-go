package signalflow

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/signalfx/signalfx-go/idtool"
	"github.com/signalfx/signalfx-go/signalflow/messages"
	"github.com/stretchr/testify/require"
)

func waitForDataMsg(t *testing.T, comp *Computation) messages.Message {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	select {
	case m := <-comp.Data():
		return m
	case <-ctx.Done():
		t.Fatal("data message didn't get buffered")
		return nil
	}
}

func TestBuffersDataMessages(t *testing.T) {
	ch := newChannel(context.Background(), "ch1")
	comp := newComputation(context.Background(), ch, &Client{
		defaultMetadataTimeout: 1 * time.Second,
	})
	defer comp.cancel()
	ch.AcceptMessage(&messages.DataMessage{
		Payloads: []messages.DataPayload{
			{
				TSID: idtool.ID(4000),
			},
		},
	})
	ch.AcceptMessage(&messages.MetadataMessage{
		TSID: idtool.ID(4000),
	})

	require.NotNil(t, comp.TSIDMetadata(4000))

	ch.AcceptMessage(&messages.InfoMessage{})

	msg := waitForDataMsg(t, comp)
	require.Equal(t, idtool.ID(4000), msg.(*messages.DataMessage).Payloads[0].TSID)

	ch.AcceptMessage(&messages.DataMessage{
		Payloads: []messages.DataPayload{
			{
				TSID: idtool.ID(4001),
			},
		},
	})
	msg = waitForDataMsg(t, comp)
	require.Equal(t, idtool.ID(4001), msg.(*messages.DataMessage).Payloads[0].TSID)
}

func TestBuffersExpiryMessages(t *testing.T) {
	ch := newChannel(context.Background(), "ch1")
	comp := newComputation(context.Background(), ch, &Client{
		defaultMetadataTimeout: 1 * time.Second,
	})
	defer comp.cancel()
	ch.AcceptMessage(&messages.DataMessage{
		Payloads: []messages.DataPayload{
			{
				TSID: idtool.ID(4000),
			},
		},
	})
	ch.AcceptMessage(&messages.MetadataMessage{
		TSID: idtool.ID(4000),
	})

	require.NotNil(t, comp.TSIDMetadata(4000))

	ch.AcceptMessage(&messages.InfoMessage{})

	msg := waitForDataMsg(t, comp)
	require.Equal(t, idtool.ID(4000), msg.(*messages.DataMessage).Payloads[0].TSID)

	ch.AcceptMessage(&messages.DataMessage{
		Payloads: []messages.DataPayload{
			{
				TSID: idtool.ID(4001),
			},
		},
	})
	msg = waitForDataMsg(t, comp)
	require.Equal(t, idtool.ID(4001), msg.(*messages.DataMessage).Payloads[0].TSID)
}

func mustParse(m messages.Message, err error) messages.Message {
	if err != nil {
		panic(err)
	}
	return m
}

func TestResolutionMetadata(t *testing.T) {
	ch := newChannel(context.Background(), "ch1")
	comp := newComputation(context.Background(), ch, &Client{
		defaultMetadataTimeout: 1 * time.Second,
	})
	defer comp.cancel()

	wg := sync.WaitGroup{}

	// Ensure multiple calls get the same result and also wait for the message
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			require.Equal(t, 5*time.Second, comp.Resolution())
			wg.Done()
		}()
	}

	ch.AcceptMessage(mustParse(messages.ParseMessage([]byte(`{
		"type": "message",
		"message": {
			"messageCode": "JOB_RUNNING_RESOLUTION",
			"contents": {
				"resolutionMs": 5000
			}
		}
	}`), true)))

	wg.Wait()
}

func TestMaxDelayMetadata(t *testing.T) {
	ch := newChannel(context.Background(), "ch1")
	comp := newComputation(context.Background(), ch, &Client{
		defaultMetadataTimeout: 1 * time.Second,
	})
	defer comp.cancel()
	ch.AcceptMessage(mustParse(messages.ParseMessage([]byte(`{
		"type": "message",
		"message": {
			"messageCode": "JOB_INITIAL_MAX_DELAY",
			"contents": {
				"maxDelayMs": 1000
			}
		}
	}`), true)))

	require.Equal(t, 1*time.Second, comp.MaxDelay())
}

func TestLagMetadata(t *testing.T) {
	ch := newChannel(context.Background(), "ch1")
	comp := newComputation(context.Background(), ch, &Client{
		defaultMetadataTimeout: 1 * time.Second,
	})
	defer comp.cancel()
	ch.AcceptMessage(mustParse(messages.ParseMessage([]byte(`{
		"type": "message",
		"message": {
			"messageCode": "JOB_DETECTED_LAG",
			"contents": {
				"lagMs": 3500
			}
		}
	}`), true)))

	require.Equal(t, 3500*time.Millisecond, comp.Lag())
}

func TestHandle(t *testing.T) {
	ch := newChannel(context.Background(), "ch1")
	comp := newComputation(context.Background(), ch, &Client{
		defaultMetadataTimeout: 1 * time.Second,
	})
	defer comp.cancel()
	ch.AcceptMessage(mustParse(messages.ParseMessage([]byte(`{
		"type": "control-message",
		"event": "JOB_START",
		"handle": "AAAABBBB"
	}`), true)))

	require.Equal(t, "AAAABBBB", comp.Handle())
}
