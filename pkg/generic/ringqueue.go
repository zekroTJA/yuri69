package generic

import (
	"fmt"
	"sync"
)

type RingQueue[T any] struct {
	m    sync.Mutex
	head int
	tail int
	buff []T
}

func NewRingQueue[T any](size uint) *RingQueue[T] {
	if size == 0 {
		size = 1
	}

	return &RingQueue[T]{
		head: 0,
		tail: 0,
		buff: make([]T, size),
	}
}

func (t *RingQueue[T]) Enqueue(v T) {
	t.m.Lock()
	defer t.m.Unlock()

	t.buff[t.head] = v
	t.advanceHead()
}

func (t *RingQueue[T]) Dequeue() T {
	t.m.Lock()
	defer t.m.Unlock()

	v := t.buff[t.tail]
	t.advanceTail()

	return v
}

func (t *RingQueue[T]) Reset() {
	t.m.Lock()
	defer t.m.Unlock()

	t.head = 0
	t.tail = 0

	var d T
	for i := 0; i < len(t.buff); i++ {
		t.buff[i] = d
	}
}

func (t *RingQueue[T]) Snapshot() []T {
	t.m.Lock()
	defer t.m.Unlock()

	out := make([]T, len(t.buff))
	copy(out, t.buff)

	return out
}

func (t *RingQueue[T]) Size() int {
	return len(t.buff)
}

func (t *RingQueue[T]) Resize(size int) {
	t.m.Lock()
	defer t.m.Unlock()

	if size < 1 {
		size = 1
	}

	delta := size - len(t.buff)

	if delta == 0 {
		return
	}

	if delta > 0 {
		t.buff = append(t.buff, make([]T, delta)...)
		return
	}

	newBuff := make([]T, size)
	copy(newBuff, t.buff[len(t.buff)+delta:])

	t.buff = newBuff

	if t.head >= size {
		t.head = size - 1
	}
	if t.tail >= size {
		t.tail = size - 1
	}

	return
}

// --- Helpers ---

func (t *RingQueue[T]) advanceHead() {
	t.head++
	if t.head == len(t.buff) {
		t.head = 0
	}
}

func (t *RingQueue[T]) advanceTail() {
	if t.tail == t.head {
		t.advanceHead()
	}

	t.tail++
	if t.tail == len(t.buff) {
		t.tail = 0
	}
}

// --- Helpers for testing

func (t *RingQueue[T]) visualize() {
	var r0, r1, r2, r3 string
	for i, v := range t.buff {
		r0 += fmt.Sprintf("%d ", i)
		r1 += fmt.Sprintf("%v ", v)
		if i == t.head {
			r2 += "H "
		} else {
			r2 += "  "
		}
		if i == t.tail {
			r3 += "T "
		} else {
			r3 += "  "
		}
	}
	fmt.Print("----------\n", r0, "\n", r1, "\n", r2, "\n", r3, "\n")
}
