package collection

// Set is a generic set data structure that holds unique elements of type T.
type Set[T comparable] struct {
	data map[T]struct{}
}

// NewSet creates and returns a new instance of Set.
func NewSet[T comparable]() *Set[T] {
	return &Set[T]{
		data: make(map[T]struct{}),
	}
}

// Add adds one or more items to the Set.
func (s *Set[T]) Add(items ...T) {
	for _, item := range items {
		s.data[item] = struct{}{}
	}
}

// Remove removes one or more items from the Set.
func (s *Set[T]) Remove(items ...T) {
	for _, item := range items {
		delete(s.data, item)
	}
}

// Constrains checks if the item exists in the Set.
func (s *Set[T]) Constrains(item T) bool {
	_, exists := s.data[item]
	return exists
}

// Clear removes all items from the Set.
func (s *Set[T]) Clear() {
	s.data = make(map[T]struct{})
}

// Size returns the number of items in the Set.
func (s *Set[T]) Size() int {
	return len(s.data)
}

// Keys returns a slice of all items in the Set.
func (s *Set[T]) Keys() []T {
	keys := make([]T, 0, len(s.data))
	for k := range s.data {
		keys = append(keys, k)
	}
	return keys
}

// Range iterates over all items in the Set and applies the given function.
func (s *Set[T]) Range(f func(item T) bool) {
	for k := range s.data {
		if !f(k) {
			break
		}
	}
}
