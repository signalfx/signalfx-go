package writer

import (
	"math/rand"
	"testing"
	"time"

	"github.com/signalfx/golib/v3/trace"
	"github.com/stretchr/testify/require"
)

func TestSpanBuffer(t *testing.T) {
	t.Parallel()

	t.Run("Return elements added", func(t *testing.T) {
		t.Parallel()
		buffer := NewSpanRingBuffer(100)

		for i := 0; i < buffer.Size(); i++ {
			require.Equal(t, buffer.UnprocessedCount(), i)
			overwrote := buffer.Add(&trace.Span{
				Meta: map[interface{}]interface{}{"i": i},
			})

			require.False(t, overwrote)
		}

		for i := 0; i < 10; i++ {
			require.Equal(t, buffer.UnprocessedCount(), buffer.Size()-i*10)
			batch := buffer.NextBatch(10)
			for j := 0; j < len(batch); j++ {
				require.Equal(t, batch[j].Meta["i"].(int), i*10+j)
			}
		}
	})

	t.Run("Overwrites older elements", func(t *testing.T) {
		t.Parallel()
		buffer := NewSpanRingBuffer(100)

		for i := 0; i < buffer.Size(); i++ {
			overwrote := buffer.Add(&trace.Span{
				Meta: map[interface{}]interface{}{"i": i},
			})

			require.False(t, overwrote)
		}

		// Pull out the first 9 elements
		batch := buffer.NextBatch(9)

		for i := 0; i < 9; i++ {
			require.Equal(t, batch[i].Meta["i"].(int), i)
		}

		// We should be able to add 9 more elements without overwriting
		// anything.
		for i := 0; i < 9; i++ {
			overwrote := buffer.Add(&trace.Span{
				Meta: map[interface{}]interface{}{"i": i + 100},
			})

			require.False(t, overwrote)
		}

		// Overwrite elements [9..59).  This drags the read cursor to the
		// first non-overwritten element, 59.
		for i := 0; i < 50; i++ {
			overwrote := buffer.Add(&trace.Span{
				Meta: map[interface{}]interface{}{"i": i + 109},
			})

			require.True(t, overwrote)
		}

		var rest []*trace.Span
		for i := 0; i < 11; i++ {
			rest = append(rest, buffer.NextBatch(10)...)
		}

		require.Len(t, rest, 100)
		for i := 0; i < 100; i++ {
			require.Equal(t, rest[i].Meta["i"].(int), 59+i)
		}

		require.Equal(t, buffer.UnprocessedCount(), 0)

		// Fire up a bunch of random sized reads after a random number of adds
		// and make sure things stay consistent.
		totalCount := 1000000
		expectedCount := totalCount
		var overwroteIs []int
		var receivedIs []int

		rng := rand.New(rand.NewSource(time.Now().UnixNano()))

		totalAdded := 0
		done := false
		for !done {
			addCount := rng.Intn(200) // 0-199
			var maxBatchSize int

			for i := totalAdded; i < totalAdded+addCount && i < totalCount; i++ {
				overwrote := buffer.Add(&trace.Span{
					Meta: map[interface{}]interface{}{"i": i},
				})
				if overwrote {
					expectedCount--
					overwroteIs = append(overwroteIs, i-buffer.Size())
				}
			}
			totalAdded += addCount

			if totalAdded >= totalCount {
				maxBatchSize = buffer.Size()
				done = true
			} else {
				maxBatchSize = rng.Intn(100) + 1 // 1-100
			}

			received := buffer.NextBatch(maxBatchSize)
			for i := range received {
				receivedIs = append(receivedIs, received[i].Meta["i"].(int))
			}
		}
		// Do one final receive in case the last receive hit the top of the
		// buffer (it doesn't wrap around when returning batches).
		received := buffer.NextBatch(buffer.Size())
		for i := range received {
			receivedIs = append(receivedIs, received[i].Meta["i"].(int))
		}

		require.Less(t, expectedCount, totalCount)
		require.Equal(t, len(overwroteIs), totalCount-expectedCount)
		require.Equal(t, len(overwroteIs)+len(receivedIs), totalCount)

		for i := 0; i < totalCount; i++ {
			nextReceivedI := receivedIs[0]
			if i == nextReceivedI {
				receivedIs = receivedIs[1:]
				continue
			}
			nextOverwroteI := overwroteIs[0]
			if i == nextOverwroteI {
				overwroteIs = overwroteIs[1:]
				continue
			}
			require.Fail(t, "an element was neither marked as received or overwritten: %d", i)
		}
	})
}
