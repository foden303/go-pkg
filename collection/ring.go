package collection

import "sync"

// Ring is a fixed-size circular buffer that overwrites the oldest elements when full.
type Ring struct {
	lock     sync.RWMutex
	elements []any
	index    int
}

// NewRing creates and returns a new Ring buffer with the specified size.
func NewRing(size int) *Ring {
	if size < 1 {
		panic("n should be greater than 0")
	}
	return &Ring{
		elements: make([]any, size),
	}
}

// Add adds an element to the Ring buffer, overwriting the oldest element if necessary.
func (r *Ring) Add(element any) {
	r.lock.Lock()
	defer r.lock.Unlock()

	rLength := len(r.elements)
	r.elements[r.index%rLength] = element
	r.index++

	// prevent ring index overflow
	if r.index >= rLength<<1 { // rLength << 1 == rLength * 2
		r.index -= rLength
	}
}

// Take returns a slice of the current elements in the Ring buffer in the order they were added.
func (r *Ring) Take() []any {
	r.lock.RLock()
	defer r.lock.RUnlock()

	var size int
	var start int
	rLength := len(r.elements)

	if r.index > rLength {
		size = rLength
		start = r.index % rLength
	} else {
		size = r.index
	}
	// safety copy to avoid data race
	elements := make([]any, size)
	for i := 0; i < size; i++ {
		elements[i] = r.elements[(start+i)%rLength]
	}

	return elements
}
