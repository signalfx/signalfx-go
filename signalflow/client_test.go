package signalflow

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/signalfx/signalfx-go/idtool"
	"github.com/signalfx/signalfx-go/signalflow/messages"
	"github.com/stretchr/testify/require"
)

func TestAuthenticationFlow(t *testing.T) {
	fakeBackend := NewRunningFakeBackend()
	defer fakeBackend.Stop()

	c, err := NewClient(StreamURL(fakeBackend.URL()), AccessToken(fakeBackend.AccessToken))
	require.Nil(t, err)

	comp, err := c.Execute(&ExecuteRequest{
		Program: "data('cpu.utilization').publish()",
	})
	require.Nil(t, err)

	require.Equal(t, 1*time.Second, comp.Resolution())

	require.Equal(t, []map[string]interface{}{
		{"type": "authenticate",
			"token": fakeBackend.AccessToken},
		{"type": "execute",
			"channel":    "ch-1",
			"immediate":  false,
			"maxDelay":   0.,
			"program":    "data('cpu.utilization').publish()",
			"resolution": 0.,
			"start":      0.,
			"stop":       0.,
			"timezone":   ""},
	}, fakeBackend.received)
}

func TestBasicComputation(t *testing.T) {
	fakeBackend := NewRunningFakeBackend()
	defer fakeBackend.Stop()

	c, err := NewClient(StreamURL(fakeBackend.URL()), AccessToken(fakeBackend.AccessToken))
	require.Nil(t, err)

	tsids := []idtool.ID{idtool.ID(rand.Int63()), idtool.ID(rand.Int63())}
	for i, host := range []string{"host1", "host2"} {
		fakeBackend.AddTSIDMetadata(tsids[i], &messages.MetadataProperties{
			Metric: "jobs_queued",
			CustomProperties: map[string]string{
				"host": host,
			},
		})
	}

	for i, val := range []float64{5, 10} {
		fakeBackend.SetTSIDFloatData(tsids[i], val)
	}

	program := "data('cpu.utilization').publish()"
	fakeBackend.AddProgramTSIDs(program, tsids)

	comp, err := c.Execute(&ExecuteRequest{
		Program:    program,
		Resolution: time.Duration(1001) * time.Second,
	})
	require.Nil(t, err)

	require.Equal(t, time.Duration(1001)*time.Second, comp.Resolution())

	dataMsg := <-comp.Data()
	require.Len(t, dataMsg.Payloads, 2)
	require.Equal(t, dataMsg.Payloads[0].Float64(), float64(5))
	require.Equal(t, dataMsg.Payloads[1].Float64(), float64(10))
}

func TestMultipleComputations(t *testing.T) {
	fakeBackend := NewRunningFakeBackend()
	defer fakeBackend.Stop()

	c, err := NewClient(StreamURL(fakeBackend.URL()), AccessToken(fakeBackend.AccessToken))
	require.Nil(t, err)

	for i := 1; i < 50; i++ {
		comp, err := c.Execute(&ExecuteRequest{
			Program:    "data('cpu.utilization').publish()",
			Resolution: time.Duration(i) * time.Second,
		})
		require.Nil(t, err)

		require.Equal(t, time.Duration(i)*time.Second, comp.Resolution())
		require.Equal(t, fmt.Sprintf("ch-%d", i), comp.Channel().name)
	}
}

func TestShutdown(t *testing.T) {
	fakeBackend := NewRunningFakeBackend()
	defer fakeBackend.Stop()

	c, err := NewClient(StreamURL(fakeBackend.URL()), AccessToken(fakeBackend.AccessToken))
	require.Nil(t, err)

	var comps []*Computation
	for i := 1; i < 3; i++ {
		comp, err := c.Execute(&ExecuteRequest{
			Program:    "data('cpu.utilization').publish()",
			Resolution: time.Duration(i) * time.Second,
		})
		require.Nil(t, err)
		comps = append(comps, comp)

		require.Equal(t, time.Duration(i)*time.Second, comp.Resolution())
		require.Equal(t, fmt.Sprintf("ch-%d", i), comp.Channel().name)
	}

	c.Close()

	for _, comp := range comps {
		require.True(t, comp.IsFinished())
	}
}
