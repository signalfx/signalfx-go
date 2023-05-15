package signalflow

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"runtime"
	"runtime/pprof"
	"sync"
	"testing"
	"time"

	"github.com/signalfx/signalfx-go/idtool"
	"github.com/signalfx/signalfx-go/signalflow/v2/messages"
	"github.com/stretchr/testify/require"
)

func TestAuthenticationFlow(t *testing.T) {
	t.Parallel()
	fakeBackend := NewRunningFakeBackend()
	defer fakeBackend.Stop()

	c, err := NewClient(StreamURL(fakeBackend.URL()), AccessToken(fakeBackend.AccessToken))
	require.Nil(t, err)
	defer c.Close()

	comp, err := c.Execute(context.Background(), &ExecuteRequest{
		Program: "data('cpu.utilization').publish()",
	})
	require.Nil(t, err)

	resolution, _ := comp.Resolution(context.Background())
	require.Equal(t, 1*time.Second, resolution)

	require.Equal(t, []map[string]interface{}{
		{
			"type":  "authenticate",
			"token": fakeBackend.AccessToken,
		},
		{
			"type":       "execute",
			"channel":    "ch-1",
			"immediate":  false,
			"maxDelay":   0.,
			"program":    "data('cpu.utilization').publish()",
			"resolution": 0.,
			"start":      0.,
			"stop":       0.,
			"timezone":   "",
		},
	}, fakeBackend.received)
}

func TestBasicComputation(t *testing.T) {
	t.Parallel()
	fakeBackend := NewRunningFakeBackend()
	defer fakeBackend.Stop()

	c, err := NewClient(StreamURL(fakeBackend.URL()), AccessToken(fakeBackend.AccessToken))
	require.Nil(t, err)
	defer c.Close()

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

	comp, err := c.Execute(context.Background(), &ExecuteRequest{
		Program:    program,
		Resolution: 1 * time.Second,
	})
	require.Nil(t, err)

	resolution, _ := comp.Resolution(context.Background())
	require.Equal(t, 1*time.Second, resolution)

	dataMsg := <-comp.Data()
	require.Len(t, dataMsg.Payloads, 2)
	require.Equal(t, dataMsg.Payloads[0].Float64(), float64(5))
	require.Equal(t, dataMsg.Payloads[1].Float64(), float64(10))
}

func TestMultipleComputations(t *testing.T) {
	t.Parallel()
	fakeBackend := NewRunningFakeBackend()
	defer fakeBackend.Stop()

	c, err := NewClient(StreamURL(fakeBackend.URL()), AccessToken(fakeBackend.AccessToken))
	require.Nil(t, err)
	defer c.Close()

	for i := 1; i < 50; i++ {
		comp, err := c.Execute(context.Background(), &ExecuteRequest{
			Program:    "data('cpu.utilization').publish()",
			Resolution: time.Duration(i) * time.Second,
		})
		require.Nil(t, err)

		resolution, _ := comp.Resolution(context.Background())
		require.Equal(t, time.Duration(i)*time.Second, resolution)
		require.Equal(t, fmt.Sprintf("ch-%d", i), comp.name)
	}
}

func TestShutdown(t *testing.T) {
	fakeBackend := NewRunningFakeBackend()
	defer fakeBackend.Stop()

	c, err := NewClient(StreamURL(fakeBackend.URL()), AccessToken(fakeBackend.AccessToken))
	require.Nil(t, err)

	var comps []*Computation
	for i := 1; i < 3; i++ {
		comp, err := c.Execute(context.Background(), &ExecuteRequest{
			Program:    "data('cpu.utilization').publish()",
			Resolution: time.Duration(i) * time.Second,
		})
		require.Nil(t, err)
		comps = append(comps, comp)

		resolution, _ := comp.Resolution(context.Background())
		require.Equal(t, time.Duration(i)*time.Second, resolution)
		require.Equal(t, fmt.Sprintf("ch-%d", i), comp.name)
	}

	c.Close()

	for _, comp := range comps {
		_, ok := <-comp.Data()
		require.False(t, ok)
	}
}

func TestReconnect(t *testing.T) {
	t.Parallel()
	fakeBackend := NewRunningFakeBackend()
	defer fakeBackend.Stop()

	c, err := NewClient(StreamURL(fakeBackend.URL()), AccessToken(fakeBackend.AccessToken))
	require.Nil(t, err)
	defer c.Close()

	comp, err := c.Execute(context.Background(), &ExecuteRequest{
		Program: "data('cpu.utilization').publish()",
	})
	require.Nil(t, err)

	resolution, _ := comp.Resolution(context.Background())
	require.Equal(t, 1*time.Second, resolution)

	require.Equal(t, []map[string]interface{}{
		{
			"type":  "authenticate",
			"token": fakeBackend.AccessToken,
		},
		{
			"type":       "execute",
			"channel":    "ch-1",
			"immediate":  false,
			"maxDelay":   0.,
			"program":    "data('cpu.utilization').publish()",
			"resolution": 0.,
			"start":      0.,
			"stop":       0.,
			"timezone":   "",
		},
	}, fakeBackend.received)

	fakeBackend.KillExistingConnections()

	for {
		_, ok := <-comp.Data()
		if !ok {
			break
		}
	}

	comp, err = c.Execute(context.Background(), &ExecuteRequest{
		Program: "data('cpu.utilization').publish()",
	})
	require.Nil(t, err)

	resolution, _ = comp.Resolution(context.Background())
	require.Equal(t, 1*time.Second, resolution)

	log.Printf("%v", fakeBackend.received)
	require.Equal(t, []map[string]interface{}{
		{
			"type":  "authenticate",
			"token": fakeBackend.AccessToken,
		},
		{
			"type":       "execute",
			"channel":    "ch-1",
			"immediate":  false,
			"maxDelay":   0.,
			"program":    "data('cpu.utilization').publish()",
			"resolution": 0.,
			"start":      0.,
			"stop":       0.,
			"timezone":   "",
		},
		{
			"type":  "authenticate",
			"token": fakeBackend.AccessToken,
		},
		{
			"type":       "execute",
			"channel":    "ch-2",
			"immediate":  false,
			"maxDelay":   0.,
			"program":    "data('cpu.utilization').publish()",
			"resolution": 0.,
			"start":      0.,
			"stop":       0.,
			"timezone":   "",
		},
	}, fakeBackend.received)
}

func TestReconnectAfterBackendDown(t *testing.T) {
	t.Parallel()
	fakeBackend := NewRunningFakeBackend()
	defer fakeBackend.Stop()

	c, err := NewClient(StreamURL(fakeBackend.URL()), AccessToken(fakeBackend.AccessToken))
	require.Nil(t, err)

	defer c.Close()

	comp, err := c.Execute(context.Background(), &ExecuteRequest{
		Program: "data('cpu.utilization').publish()",
	})
	require.Nil(t, err)

	resolution, _ := comp.Resolution(context.Background())
	require.Equal(t, 1*time.Second, resolution)

	require.Equal(t, []map[string]interface{}{
		{
			"type":  "authenticate",
			"token": fakeBackend.AccessToken,
		},
		{
			"type":       "execute",
			"channel":    "ch-1",
			"immediate":  false,
			"maxDelay":   0.,
			"program":    "data('cpu.utilization').publish()",
			"resolution": 0.,
			"start":      0.,
			"stop":       0.,
			"timezone":   "",
		},
	}, fakeBackend.received)

	fakeBackend.Stop()
	for {
		_, ok := <-comp.Data()
		if !ok {
			break
		}
	}

	time.Sleep(7 * time.Second)
	fakeBackend.Restart()

	comp, err = c.Execute(context.Background(), &ExecuteRequest{
		Program: "data('cpu.utilization').publish()",
	})
	require.Nil(t, err)

	resolution, _ = comp.Resolution(context.Background())
	require.Equal(t, 1*time.Second, resolution)

	require.Equal(t, []map[string]interface{}{
		{
			"type":  "authenticate",
			"token": fakeBackend.AccessToken,
		},
		{
			"type":       "execute",
			"channel":    "ch-1",
			"immediate":  false,
			"maxDelay":   0.,
			"program":    "data('cpu.utilization').publish()",
			"resolution": 0.,
			"start":      0.,
			"stop":       0.,
			"timezone":   "",
		},
		{
			"type":  "authenticate",
			"token": fakeBackend.AccessToken,
		},
		{
			"type":       "execute",
			"channel":    "ch-2",
			"immediate":  false,
			"maxDelay":   0.,
			"program":    "data('cpu.utilization').publish()",
			"resolution": 0.,
			"start":      0.,
			"stop":       0.,
			"timezone":   "",
		},
	}, fakeBackend.received)
}

func TestFailedConnGoroutineShutdown(t *testing.T) {
	defer func() {
		time.Sleep(2 * time.Second)
		pprof.Lookup("goroutine").WriteTo(os.Stdout, 1)
	}()

	fakeBackend := NewRunningFakeBackend()
	fakeBackend.Stop()

	startingGoroutines := runtime.NumGoroutine()
	clients := make([]*Client, 100)
	var wg sync.WaitGroup
	for i := range clients {
		var err error
		clients[i], err = NewClient(StreamURL(fakeBackend.URL()), AccessToken(fakeBackend.AccessToken))
		require.Nil(t, err)

		wg.Add(1)
		go func(c *Client) {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			_, err = c.Execute(ctx, &ExecuteRequest{
				Program: "data('cpu.utilization').publish()",
			})
			cancel()
			t.Logf("execute error: %v", err)
			require.Error(t, err)
			wg.Done()
		}(clients[i])
	}
	wg.Wait()

	for _, c := range clients {
		c.Close()
	}
	time.Sleep(1 * time.Second)

	require.InDelta(t, startingGoroutines, runtime.NumGoroutine(), 10)
}
