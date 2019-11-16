package template

// nolint: dupl

import (
	"github.com/mauricelam/genny/generic"
)

// Instance is the element type
type Instance generic.Type

// InstanceRingBuffer is a ring buffer that supports inserting and reading
// chunks of elements in an orderly fashion.  It is NOT thread-safe and the
// returned batches are not copied, they are a slice against the original
// backing array of this instance.  This means that if the buffer wraps around,
// elements in the slice returned by NextBatch will be changed, and you are
// subject to all of the rules of Go's memory model if accessing the data in a
// separate goroutine.
type InstanceRingBuffer struct {
	// The main buffer
	buffer []*Instance
	// Store length explicitly as optimization
	bufferLen int

	nextIdx int
	// How many times around the ring buffer we have gone when putting
	// datapoints onto the buffer
	writtenCircuits int64

	// The index that indicates the last read position in the buffer.  It is
	// one greater than the actual index, to match the golang slice high range.
	readHigh int

	// How many elements are in the buffer on which processing has not begun.
	// This could be calculated from readHigh and nextIdx on demand, but
	// precalculate it in Add and NextBatch since it tends to get read often.
	// Also by precalculating it, we can tell if the buffer was completely
	// overwritten since the last read.
	unprocessed int
}

// NewInstanceRingBuffer creates a new initialized buffer ready for use.
func NewInstanceRingBuffer(size int) *InstanceRingBuffer {
	return &InstanceRingBuffer{
		// Preallocate the buffer to its maximum length
		buffer:    make([]*Instance, size),
		bufferLen: size,
	}
}

// Add an Instance to the buffer. It will overwrite any existing element in the
// buffer as the buffer wraps around.  Returns whether the new element
// overwrites an uncommitted element already in the buffer.
func (b *InstanceRingBuffer) Add(inst *Instance) (isOverwrite bool) {
	if b.unprocessed >= b.bufferLen {
		isOverwrite = true
		// Drag the read cursor along with the overwritten elements
		b.readHigh++
		if b.readHigh > b.bufferLen {
			// Wrap around to cover the 0th element of the buffer
			b.readHigh = 1
		}
	} else {
		b.unprocessed++
	}

	b.buffer[b.nextIdx] = inst
	b.nextIdx++

	if b.nextIdx == b.bufferLen { // Wrap around the buffer
		b.nextIdx = 0
		b.writtenCircuits++
	}

	return isOverwrite
}

// Size returns how many elements can fit in the buffer at once.
func (b *InstanceRingBuffer) Size() int {
	return b.bufferLen
}

// UnprocessedCount returns the number of elements that have been written to
// the buffer but not read via NextBatch.
func (b *InstanceRingBuffer) UnprocessedCount() int {
	return b.unprocessed
}

// NextBatch returns the next batch of unprocessed elements.  If there are
// none, this can return nil.
func (b *InstanceRingBuffer) NextBatch(maxSize int) []*Instance {
	prevReadHigh := b.readHigh
	if prevReadHigh == b.bufferLen {
		// Wrap around
		prevReadHigh = 0
	}

	if b.unprocessed == 0 {
		return nil
	}

	targetSize := b.unprocessed
	if targetSize > maxSize {
		targetSize = maxSize
	}

	b.readHigh = prevReadHigh + targetSize
	if b.readHigh > b.bufferLen {
		// Wrap around happened, just take what we have left until wrap around
		// so that we can take a single slice of it since slice ranges can't
		// wrap around.
		b.readHigh = b.bufferLen
	}

	b.unprocessed -= b.readHigh - prevReadHigh

	out := b.buffer[prevReadHigh:b.readHigh]

	return out
}
