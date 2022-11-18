// Package syncslice provides a simple goroutine-safe wrapper interface over a generic slice using [sync.RWMutex].
package syncslice

import "sync"

// SliceOf is a struct that abstracts access over a slice of R in a goroutine-safe way.
type SliceOf[R any] struct {
	mutex sync.RWMutex

	items []R
}

// SliceItem is a helper struct emitted in the channel returned from [SliceOf.Iter].
type sliceItem[R any] struct {
	Index int
	Value R
}

// Make[R](size, cap) is equivalent to make([]R, size, cap).
func Make[R any](size int, cap int) SliceOf[R] {
	return SliceOf[R]{items: make([]R, size, cap)}
}

// Append adds a new item to the end of the slice, in-place. Equivalent to slice = append(slice, item)
func (s *SliceOf[R]) Append(item ...R) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.items = append(s.items, item...)
}

// Set puts the value at the given index of the slice. Equivalent to slice[index] = value.
func (s *SliceOf[R]) Set(index int, value R) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.items[index] = value
}

// Get reads the value at the given index of the slice. Equivalent to slice[index]
func (s *SliceOf[R]) Get(index int) R {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return s.items[index]
}

// Len reads the size of the slice. Equivalent to len(slice).
func (s *SliceOf[R]) Len() int {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return len(s.items)
}

// Iter returns a channel that can be iterated on, in case you want to do anything with the full slice.
// Similar to [SliceOf.Range].
//
//     for item := range s.Iter() {
//         fmt.Println(item.Index, item.Value)
//     }
func (s *SliceOf[R]) Iter() <-chan sliceItem[R] {
	c := make(chan sliceItem[R])

	f := func() {
		s.mutex.RLock()
		defer s.mutex.RUnlock()
		for index, value := range s.items {
			c <- sliceItem[R]{index, value}
		}
		close(c)
	}
	go f()

	return c
}

// Range calls op for each element in the slice, if op returns false it stops the iteration.
func (s *SliceOf[R]) Range(op func(index int, value R) bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	for index, value := range s.items {
		proceed := op(index, value)
		if !proceed {
			break
		}
	}
}

// Slice takes a part of the existing slice and returns a new SliceOf[R]. Equivalent to slice[from:to].
func (s *SliceOf[R]) Slice(from, to int) *SliceOf[R] {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return &SliceOf[R]{items: s.items[from:to]}
}

// Do takes a function op and calls it with the full contents of the slice in its raw form, for any lower-level
// operation than the ones provided by [SliceOf].
func (s *SliceOf[R]) Do(op func(slice []R)) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	op(s.items)
}
