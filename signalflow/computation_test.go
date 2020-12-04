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

func waitForDataMsg(t *testing.T, comp *Computation) (messages.Message, error) {
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

	msg, _ := waitForDataMsg(t, comp)
	require.Equal(t, idtool.ID(4000), msg.(*messages.DataMessage).Payloads[0].TSID)

	ch.AcceptMessage(&messages.DataMessage{
		Payloads: []messages.DataPayload{
			{
				TSID: idtool.ID(4001),
			},
		},
	})
	msg, _ = waitForDataMsg(t, comp)
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

	msg, _ := waitForDataMsg(t, comp)
	require.Equal(t, idtool.ID(4000), msg.(*messages.DataMessage).Payloads[0].TSID)

	ch.AcceptMessage(&messages.DataMessage{
		Payloads: []messages.DataPayload{
			{
				TSID: idtool.ID(4001),
			},
		},
	})
	msg, _ = waitForDataMsg(t, comp)
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

func TestFindLimitedResultSetMetadata(t *testing.T) {
	ch := newChannel(context.Background(), "ch1")
	comp := newComputation(context.Background(), ch, &Client{
		defaultMetadataTimeout: 1 * time.Second,
	})
	defer comp.cancel()
	ch.AcceptMessage(mustParse(messages.ParseMessage([]byte(`{
		"type": "message",
		"message": {
			"messageCode": "FIND_LIMITED_RESULT_SET",
			"contents": {
				"matchedSize": 123456789,
				"limitSize": 50000
			}
		}
	}`), true)))

	require.Equal(t, 123456789, comp.MatchedSize())
	require.Equal(t, 50000, comp.LimitSize())
}

func TestMatchedNoTimeseriesQueryMetaData(t *testing.T) {
	ch := newChannel(context.Background(), "ch1")
	comp := newComputation(context.Background(), ch, &Client{
		defaultMetadataTimeout: 1 * time.Second,
	})
	defer comp.cancel()
	ch.AcceptMessage(mustParse(messages.ParseMessage([]byte(`{
		"type": "message",
		"message": {
			"messageCode": "FIND_MATCHED_NO_TIMESERIES",
			"contents": {
				"query": "abc"
			}
		}
	}`), true)))

	require.Equal(t, "abc", comp.MatchedNoTimeseriesQuery())
}

func TestGroupByMissingPropertyMetaData(t *testing.T) {
	ch := newChannel(context.Background(), "ch1")
	comp := newComputation(context.Background(), ch, &Client{
		defaultMetadataTimeout: 1 * time.Second,
	})
	defer comp.cancel()
	ch.AcceptMessage(mustParse(messages.ParseMessage([]byte(`{
		"type": "message",
		"message": {
			"messageCode": "GROUPBY_MISSING_PROPERTY",
			"contents": {
				"propertyNames": ["x", "y", "z"]
			}
		}
	}`), true)))

	require.Equal(t, []string{"x", "y", "z"}, comp.GroupByMissingProperties())
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

func TestComputationError(t *testing.T) {
	ch := newChannel(context.Background(), "ch1")
	comp := newComputation(context.Background(), ch, &Client{
		defaultMetadataTimeout: 1 * time.Second,
	})
	defer comp.cancel()
	ch.AcceptMessage(mustParse(messages.ParseMessage([]byte(`{
		"type": "error",
		"error": 400,
		"errorType": "ANALYTICS_PROGRAM_NAME_ERROR",
		"message": "We hit some error"
	}`), true)))

	_, err := waitForDataMsg(t, comp)
	if err == nil {
		t.Fatal("Expected computation error")
	}
	require.Equal(t, 400, err.(*ComputationError).Code)
	require.Equal(t, "ANALYTICS_PROGRAM_NAME_ERROR", err.(*ComputationError).ErrorType)
	require.Equal(t, "We hit some error", err.(*ComputationError).Message)
}

func TestComputationFinish(t *testing.T) {
	ch := newChannel(context.Background(), "ch1")
	comp := newComputation(context.Background(), ch, &Client{
		defaultMetadataTimeout: 1 * time.Second,
	})
	defer comp.cancel()
	go func() {
		ch.AcceptMessage(mustParse(messages.ParseMessage([]byte(`{
			"type": "control-message",
			"event": "JOB_START",
			"handle": "AAAABBBB"
		}`), true)))

		ch.AcceptMessage(&messages.MetadataMessage{
			TSID: idtool.ID(4000),
		})

		ch.AcceptMessage(&messages.DataMessage{
			Payloads: []messages.DataPayload{
				{
					TSID: idtool.ID(4000),
				},
			},
		})

		ch.AcceptMessage(&messages.EndOfChannelControlMessage{})
	}()

	for msg := range comp.Data() {
		require.Equal(t, idtool.ID(4000), msg.Payloads[0].TSID)
	}

	// The for loop should exit when the end of channel message comes through
}
