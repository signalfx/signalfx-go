package writer

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math"
	"sort"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/signalfx/golib/v3/trace"
	"github.com/stretchr/testify/require"

	"github.com/signalfx/golib/v3/sfxclient"
)

type traceTester struct {
	Ctx    context.Context
	Cancel context.CancelFunc

	Received []*trace.Span
	Input    chan []*trace.Span

	SendLock    sync.Mutex
	ReceiveLock sync.Mutex

	Writer *SpanWriter

	// Use atomic.Value to avoid race detection
	SendShouldFail atomic.Value
}

func (tester *traceTester) assertAllReceived(t *testing.T, expectedCount int) {
	tester.ReceiveLock.Lock()
	defer tester.ReceiveLock.Unlock()

	require.Len(t, tester.Received, int(expectedCount))

	metas := make([]int, expectedCount)
	expected := make([]int, expectedCount)
	for i := 0; i < expectedCount; i++ {
		metas[i] = tester.Received[i].Meta["i"].(int)
		expected[i] = i
	}
	require.ElementsMatch(t, metas, expected)
}

func setupSpanTesting(chanSize int) *traceTester {
	ts := &traceTester{
		Input: make(chan []*trace.Span, chanSize),
	}

	ts.Ctx, ts.Cancel = context.WithCancel(context.Background())

	filter := func(inst *trace.Span) bool {
		shouldSend, ok := inst.Meta["shouldSend"].(bool)
		return !ok || shouldSend
	}

	sender := func(ctx context.Context, insts []*trace.Span) error {
		if ts.SendShouldFail.Load().(bool) {
			return errors.New("failed")
		}
		ts.ReceiveLock.Lock()
		ts.Received = append(ts.Received, insts...)
		ts.ReceiveLock.Unlock()

		// The lock can be used to block the sender by test code
		ts.SendLock.Lock()
		//nolint:staticcheck
		ts.SendLock.Unlock()
		return nil
	}

	ts.Writer = &SpanWriter{
		PreprocessFunc: filter,
		SendFunc:       sender,
		InputChan:      ts.Input,
	}

	ts.SendShouldFail.Store(false)

	return ts
}

func TestSpanWriter(t *testing.T) {
	t.Run("Should send all traces received", func(t *testing.T) {
		t.Parallel()
		ts := setupSpanTesting(1000)

		ts.Writer.MaxBuffered = 20000
		ts.Writer.Start(ts.Ctx)

		count := 0
		for i := 0; i < 10000; i++ {
			ts.Input <- []*trace.Span{{Meta: map[interface{}]interface{}{"i": i}}}
			count++
		}

		ts.Cancel()
		ts.Writer.WaitForShutdown()

		ts.assertAllReceived(t, count)
	})

	t.Run("Should panic if waiting without starting", func(t *testing.T) {
		t.Parallel()
		ts := setupSpanTesting(1000)
		require.Panics(t, ts.Writer.WaitForShutdown)
	})

	for _, inBatchSize := range []int{1, 2, 3, 5, 8, 13} {
		t.Run(fmt.Sprintf("Should cycle buffer without losing anything (inBatchSize: %d)", inBatchSize), func(t *testing.T) {
			t.Parallel()
			ts := setupSpanTesting(0)
			ts.Writer.MaxBuffered = 3000
			ts.Writer.MaxRequests = 100
			ts.Writer.MaxBatchSize = 1000
			ts.Writer.Start(ts.Ctx)

			inputCount := 100000

			i := 0
			var batch []*trace.Span
			for i <= inputCount {
				if i == inputCount || len(batch) == inBatchSize {
					// Slow it down a bit so it doesn't wraparound the buffer
					time.Sleep(100 * time.Nanosecond)
					ts.Input <- batch
					batch = nil
					time.Sleep(100 * time.Nanosecond)
				}
				batch = append(batch, &trace.Span{Meta: map[interface{}]interface{}{"i": i}})
				i++
			}

			ts.Cancel()
			ts.Writer.WaitForShutdown()

			spew.Dump(ts.Writer.InternalMetrics(""))

			ts.assertAllReceived(t, inputCount)

			require.Equal(t, findInternalMetricWithName(ts.Writer, "trace_spans_received"), inputCount)
			require.Equal(t, findInternalMetricWithName(ts.Writer, "trace_spans_sent"), inputCount)

		})
	}

	for _, maxBuff := range []int{100, 2500, 4999, 9997, 9999} {
		for _, inputSize := range []int{1, 11} {
			t.Run(fmt.Sprintf("Should overflow cleanly with %d max, inputSize: %d", maxBuff, inputSize), func(t *testing.T) {
				t.Parallel()
				ts := setupSpanTesting(1000)
				ts.Writer.MaxBuffered = maxBuff
				ts.Writer.MaxRequests = 2
				ts.Writer.MaxBatchSize = 5
				ts.Writer.Start(ts.Ctx)

				// Prevent things from being sent
				ts.SendLock.Lock()

				count := 0
				for i := 0; i < 100000; i += inputSize {
					var toSend []*trace.Span
					for j := 0; j < inputSize; j++ {
						toSend = append(toSend, &trace.Span{Meta: map[interface{}]interface{}{"i": i + j}})
					}
					ts.Input <- toSend
					count += len(toSend)
				}

				initialReceivedCount := 0

				go func() {
					// Wait to let input get processed a bit before letting things
					// through so that the buffer gets backed up
					for {
						ts.ReceiveLock.Lock()
						initialReceivedCount = len(ts.Received)
						ts.ReceiveLock.Unlock()
						if initialReceivedCount > 0 {
							time.Sleep(1 * time.Second)
							break
						}
						time.Sleep(100 * time.Millisecond)
					}
					ts.SendLock.Unlock()
					ts.Cancel()
				}()
				ts.Writer.WaitForShutdown()

				require.Greater(t, atomic.LoadInt64(&ts.Writer.TotalOverwritten), int64(0))

				sortSpansByMeta(ts.Received)

				expectedReceived := int(math.Min(float64(ts.Writer.MaxBuffered+initialReceivedCount), float64(count)))
				log.Printf("len(ts.Received)=%d; expectedReceive=%d; ts.Writer.MaxBuffered=%d; initialReceivedCount=%d; count=%d", len(ts.Received), expectedReceived, ts.Writer.MaxBuffered, initialReceivedCount, count)
				ts.ReceiveLock.Lock()
				require.Len(t, ts.Received, expectedReceived)

				for i := 0; i < initialReceivedCount; i++ {
					require.Equal(t, ts.Received[i].Meta["i"].(int), i)
				}

				for i := 0; i < expectedReceived-initialReceivedCount; i++ {
					require.Equal(t, ts.Received[i+initialReceivedCount].Meta["i"].(int), i+initialReceivedCount+count-expectedReceived)
				}

				require.Equal(t, atomic.LoadInt64(&ts.Writer.TotalOverwritten), int64(count-expectedReceived))
			})
		}
	}

	t.Run("Should filter out traces", func(t *testing.T) {
		t.Parallel()
		ts := setupSpanTesting(1000)
		ts.Writer.Start(ts.Ctx)

		count := 0
		for i := 0; i < 10000; i += 4 {
			ts.Input <- []*trace.Span{
				{Meta: map[interface{}]interface{}{"i": i, "shouldSend": true}},
				{Meta: map[interface{}]interface{}{"i": i + 1, "shouldSend": false}},
				{Meta: map[interface{}]interface{}{"i": i + 2, "shouldSend": true}},
				{Meta: map[interface{}]interface{}{"i": i + 3, "shouldSend": false}},
			}
			count++
		}

		ts.Cancel()
		ts.Writer.WaitForShutdown()

		require.Len(t, ts.Received, 5000)
		require.Equal(t, findInternalMetricWithName(ts.Writer, "trace_spans_filtered"), 5000)
		require.Equal(t, findInternalMetricWithName(ts.Writer, "trace_spans_sent"), 5000)
	})

	t.Run("Should report internal metrics", func(t *testing.T) {
		t.Parallel()
		ts := setupSpanTesting(1000)
		ts.Writer.Start(ts.Ctx)

		inputCount := 10000
		failPoint := 9990
		for i := 0; i < inputCount; i += 2 {
			if i >= failPoint {
				if i == failPoint {
					time.Sleep(1 * time.Second)
				}
				ts.SendShouldFail.Store(true)
			}
			ts.Input <- []*trace.Span{
				{Meta: map[interface{}]interface{}{"i": i}},
				{Meta: map[interface{}]interface{}{"i": i + 1}},
			}
		}

		ts.Cancel()
		ts.Writer.WaitForShutdown()

		require.Equal(t, findInternalMetricWithName(ts.Writer, "trace_spans_received"), inputCount)
		require.Equal(t, findInternalMetricWithName(ts.Writer, "trace_spans_sent"), len(ts.Received))
		require.Equal(t, findInternalMetricWithName(ts.Writer, "trace_spans_failed"), 10000-len(ts.Received))
		require.Equal(t, findInternalMetricWithName(ts.Writer, "trace_spans_failed"), 10000-failPoint)
		require.Equal(t, findInternalMetricWithName(ts.Writer, "trace_spans_filtered"), 0)
		require.Equal(t, findInternalMetricWithName(ts.Writer, "trace_spans_overwritten"), 0)
		require.Equal(t, findInternalMetricWithName(ts.Writer, "trace_span_requests_active"), 0)
		require.Equal(t, findInternalMetricWithName(ts.Writer, "trace_spans_in_flight"), 0)
		require.Equal(t, findInternalMetricWithName(ts.Writer, "trace_spans_waiting"), 0)
		require.Equal(t, findInternalMetricWithName(ts.Writer, "trace_spans_buffered"), 0)
		require.Equal(t, findInternalMetricWithName(ts.Writer, "trace_spans_max_buffered"), ts.Writer.MaxBuffered)
	})

	t.Run("Should call OverwriteFunc on overwrites", func(t *testing.T) {
		t.Parallel()
		ts := setupSpanTesting(1000)
		ts.Writer.MaxBatchSize = 1
		ts.Writer.MaxRequests = 1
		ts.Writer.MaxBuffered = 100

		overwrittenCount := int64(0)
		ts.Writer.OverwriteFunc = func() {
			atomic.AddInt64(&overwrittenCount, 1)
		}
		// Prevent things from being sent
		ts.SendLock.Lock()

		ts.Writer.Start(ts.Ctx)

		for i := 0; i < ts.Writer.MaxBuffered+11; i++ {
			ts.Input <- []*trace.Span{{}}
		}

		require.Eventually(t, func() bool { return atomic.LoadInt64(&overwrittenCount) == 10 }, 3*time.Second, 50*time.Millisecond)
	})

	t.Run("Can handle big inputs", func(t *testing.T) {
		t.Parallel()
		ts := setupSpanTesting(1000)
		ts.Writer.MaxBatchSize = 100
		ts.Writer.MaxRequests = 10
		ts.Writer.MaxBuffered = 100

		overwrittenCount := int64(0)
		ts.Writer.OverwriteFunc = func() {
			atomic.AddInt64(&overwrittenCount, 1)
		}
		ts.Writer.Start(ts.Ctx)

		var in []*trace.Span
		// The first 10 100-count batches of MaxBuffered instances should be sent in 10
		// requests interpersed with the processing of the input.
		for i := 0; i < ts.Writer.MaxBuffered*10; i++ {
			in = append(in, &trace.Span{Meta: map[interface{}]interface{}{"i": i}})
		}
		ts.Input <- in

		require.Eventually(t, func() bool {
			ts.ReceiveLock.Lock()
			defer ts.ReceiveLock.Unlock()
			return len(ts.Received) == ts.Writer.MaxBuffered*10 && atomic.LoadInt64(&ts.Writer.TotalOverwritten) == 0
		}, 5*time.Second, 500*time.Millisecond)
	})
}

func ExampleSpanWriter() {
	client := sfxclient.NewHTTPSink()
	filterFunc := func(dp *trace.Span) bool {
		return dp.Meta["shouldSend"].(bool)
	}

	in := make(chan []*trace.Span, 1)

	// filterFunc can also be nil if no filtering/modification is needed.
	writer := &SpanWriter{
		PreprocessFunc: filterFunc,
		SendFunc:       client.AddSpans,
		InputChan:      in,
	}

	ctx, cancel := context.WithCancel(context.Background())
	writer.Start(ctx)

	// Send traces with the writer
	in <- []*trace.Span{}

	// Close the context passed to Run()
	cancel()
	// Will wait for all pending traces to be written.
	writer.WaitForShutdown()
}

func BenchmarkSpanWriter(b *testing.B) {
	var doneSignal atomic.Value
	received := int64(0)

	total := 1000002

	in := make(chan []*trace.Span)
	writer := &SpanWriter{
		SendFunc: func(ctx context.Context, batch []*trace.Span) error {
			newReceived := atomic.AddInt64(&received, int64(len(batch)))
			if newReceived == int64(total) {
				close(doneSignal.Load().(chan struct{}))
				atomic.StoreInt64(&received, int64(0))
			}
			return nil
		},
		InputChan: in,
	}

	ctx, cancel := context.WithCancel(context.Background())
	writer.Start(ctx)
	defer cancel()

	input1 := []*trace.Span{{}}
	input2 := []*trace.Span{{}, {}}
	input3 := []*trace.Span{{}, {}, {}}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		doneSignal.Store(make(chan struct{}))
		for j := 0; j < total/6; j++ {
			in <- input1
			in <- input3
			in <- input2
		}
		<-doneSignal.Load().(chan struct{})
	}
}

func sortSpansByMeta(received []*trace.Span) {
	sort.SliceStable(received, func(i, j int) bool {
		return received[i].Meta["i"].(int) < received[j].Meta["i"].(int)
	})
}
