package generic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnqueue(t *testing.T) {
	r := NewRingQueue[int](4)

	r.Enqueue(1)
	assert.Equal(t, 1, r.head)
	assert.Equal(t, []int{1, 0, 0, 0}, r.buff)

	r.Enqueue(2)
	assert.Equal(t, 2, r.head)
	assert.Equal(t, []int{1, 2, 0, 0}, r.buff)

	r.Enqueue(3)
	assert.Equal(t, 3, r.head)
	assert.Equal(t, []int{1, 2, 3, 0}, r.buff)

	r.Enqueue(4)
	assert.Equal(t, 0, r.head)
	assert.Equal(t, []int{1, 2, 3, 4}, r.buff)

	r.Enqueue(5)
	assert.Equal(t, 1, r.head)
	assert.Equal(t, []int{5, 2, 3, 4}, r.buff)

	r.Enqueue(6)
	assert.Equal(t, 2, r.head)
	assert.Equal(t, []int{5, 6, 3, 4}, r.buff)
}

func TestDequeue(t *testing.T) {
	r := NewRingQueue[int](4)

	r.Enqueue(1)
	assert.Equal(t, 1, r.Dequeue())
	assert.Equal(t, 0, r.Dequeue())

	r.Enqueue(2)
	r.Enqueue(3)
	r.Enqueue(4)
	r.Enqueue(5)
	assert.Equal(t, 2, r.Dequeue())
	assert.Equal(t, 3, r.Dequeue())
	assert.Equal(t, 4, r.Dequeue())
	assert.Equal(t, 5, r.Dequeue())

	r.Enqueue(6)
	r.Enqueue(7)
	r.Enqueue(8)
	assert.Equal(t, 6, r.Dequeue())
	assert.Equal(t, 7, r.Dequeue())
	assert.Equal(t, 8, r.Dequeue())
}

func TestResize(t *testing.T) {
	r := NewRingQueue[int](2)

	r.Enqueue(1)
	r.Enqueue(2)
	assert.Equal(t, 1, r.Dequeue())
	assert.Equal(t, 2, r.Dequeue())
	assert.Equal(t, 1, r.Dequeue())

	r.Resize(4)
	assert.Equal(t, 2, r.Dequeue())
	assert.Equal(t, 0, r.Dequeue())

	r.Enqueue(3)
	r.Enqueue(4)
	r.Enqueue(5)
	r.Enqueue(6)
	assert.Equal(t, 3, r.Dequeue())
	assert.Equal(t, 4, r.Dequeue())
	assert.Equal(t, 5, r.Dequeue())
	assert.Equal(t, 6, r.Dequeue())
	r.visualize()

	r.Resize(2)
	r.visualize()
	assert.Equal(t, 3, r.Dequeue())
	assert.Equal(t, 6, r.Dequeue())
	assert.Equal(t, 3, r.Dequeue())
	assert.Equal(t, 6, r.Dequeue())
}

func TestEdgeCases(t *testing.T) {
	{
		r := NewRingQueue[int](1)

		r.Enqueue(1)
		assert.Equal(t, 1, r.Dequeue())
		r.Enqueue(2)
		assert.Equal(t, 2, r.Dequeue())

		r.Enqueue(3)
		r.Enqueue(4)
		r.Enqueue(5)
		assert.Equal(t, 5, r.Dequeue())
		assert.Equal(t, 5, r.Dequeue())
	}

	{
		r := NewRingQueue[int](0)

		r.Enqueue(1)
		assert.Equal(t, 1, r.Dequeue())
		r.Enqueue(2)
		assert.Equal(t, 2, r.Dequeue())

		r.Enqueue(3)
		r.Enqueue(4)
		r.Enqueue(5)
		assert.Equal(t, 5, r.Dequeue())
		assert.Equal(t, 5, r.Dequeue())
	}
}
