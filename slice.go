package slice

// Iterator is an interface type for an iterator
type Iterator[T any] interface {
	// HasNext returns whether there is a next element
	HasNext() bool

	// Next goes to the next element, then returns false if the end of the slice
	// was reached
	Next() bool

	// HasPrev returns whether there is a previous element
	HasPrev() bool

	// Prev goes to the previous element, then returns false if the start of the
	// slice was reached
	Prev() bool

	// Get gets the element the iterator is currently pointed to
	Get() T

	// Set sets the value of the element the iterator is currently pointed to
	Set(T)
}

// Slice is an interface type for a generic data structure that behaves like a []T
// type
type Slice[T any] interface {
	// Append copies the given element(s) onto the end of the slice. This function is
	// roughly equivalent to `append(slice, elems...)`
	Append(...T) Slice[T]

	// AppendSlice copies the elements in the given slice onto the end of this slice
	AppendSlice(Slice[T]) Slice[T]

	// Prepend copies the given elements onto the start of the slice. This function is
	// roughly equivalent to `append(elems, slice...)`
	Prepend(...T) Slice[T]

	// PrependSlice copies the elements in the given slice onto the start of this slice
	PrependSlice(Slice[T]) Slice[T]

	// Slice gets a subset of the slice. This function is roughly equivalent to
	// `slice[i:j]`
	Slice(int, int) Slice[T]

	// Get gets the element at the given index. This function is roughly equivalent
	// to `slice[i]`
	Get(int) T

	// Set sets the value of the element at the given index. This function is
	// roughly equivalent to `slice[i] = elem`
	Set(int, T)

	// IterStart creates an iterator, pointed to the first element
	IterStart() Iterator[T]

	// ReverseIterStart creates a reverse iterator, pointed to the first element
	ReverseIterStart() Iterator[T]

	// IterEnd creates an iterator, pointed to the last elements
	IterEnd() Iterator[T]

	// ReverseIterEnd creates a reverse iterator, pointed to the last element
	ReverseIterEnd() Iterator[T]

	// DeepCopy creates a deep copy of the slice
	DeepCopy() Slice[T]

	// Len gets the slice's length. This function is roughly equivalent to
	// `len(slice)`
	Len() int

	// Cap gets the slice's capacity. This function is roughly equivalent to
	// `cap(slice)`
	Cap() int

	// ToGoSlice converts the slice to an []T type
	ToGoSlice() []T
}

// Essentially math.Max for ints
func atLeast(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Essentially math.Min for ints
func atMost(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// ToGoSlice converts a slice to a []T type
func ToGoSlice[T any](s Slice[T]) []T {
	// Create a slice
	slice := make([]T, 0, atLeast(1, s.Len()))
	// Iterate over the elements
	iter := s.IterStart()
	for iter.Next() {
		// Add the element to the slice
		slice = append(slice, iter.Get())
	}
	return slice
}

// Erase returns a slice where the element at the given index is erased. Equivalent to
// `append(s[:index], s[index + 1:]...)`. Warning: this function does not copy
// s, so the contents of s can be (and probably will be) modified. The
// Slice.Erase() function should be used instead of this function where
// possible, as it can be faster
func Erase[T any](s Slice[T], index int) Slice[T] {
	return s.Slice(0, index).DeepCopy().AppendSlice(s.Slice(index+1, s.Len()))
}

// EraseRange returns a slice where the range of elements are erased. Equivalent to
// `append(s[:i], s[j + 1:]...)`. Warning: this function does not copy s, so the
// contents of s can be modified. The Slice.EraseRange function should be used
// instead of this function where possible, as it can be faster
func EraseRange[T any](s Slice[T], i, j int) Slice[T] {
	return s.Slice(0, i).DeepCopy().AppendSlice(s.Slice(j+1, s.Len()))
}

// Insert returns a slice where the given element was inserted at the given index.
// Equivalent to:
//
//  s = append(s[:index + 1], s[index:])
//  s[index] = elem
//
// The Slice.Insert() function should be used instead of this function where
// possible, as it can be faster
func Insert[T any](s Slice[T], index int, elem T) Slice[T] {
	s = s.Slice(0, index+1).DeepCopy().AppendSlice(s.Slice(index, s.Len()))
	s.Set(index, elem)
	return s
}

// InsertSlice returns a slice where the given slice of elements were inserted at the given
// index. Equivalent to:
//
//  s = append(s[:index + len(elems)], s[index:])
//  copy(s[index:], elems)
//
// The Slice.InsertSlice() function should be used instead of this function where
// possible, as it can be faster
func InsertSlice[T any](s Slice[T], index int, elems Slice[T]) Slice[T] {
	s = s.Slice(0, index+elems.Len()).DeepCopy().AppendSlice(s.Slice(index, s.Len()))
	iter := elems.IterStart()
	i := 0
	for iter.Next() {
		s.Set(index+i, iter.Get())
		i++
	}
	return s
}

// ReverseIterator is a reverse iterator, essentially an inverted iterator
type ReverseIterator[T any] struct {
	Iterator[T]
}

// Reverse creates a reverse iterator
func Reverse[T any](i Iterator[T]) Iterator[T] {
	// If the given iterator is a reverse iterator
	reverse, ok := i.(ReverseIterator[T])
	if ok {
		// The reverse iterator is just the unwrapped iterator
		return reverse.Iterator
	}
	// Otherwise return a new iterator
	return ReverseIterator[T]{i}
}

func (i ReverseIterator[T]) HasNext() bool {
	return i.Iterator.HasPrev()
}

func (i ReverseIterator[T]) Next() bool {
	return i.Iterator.Prev()
}

func (i ReverseIterator[T]) HasPrev() bool {
	return i.Iterator.HasNext()
}

func (i ReverseIterator[T]) Prev() bool {
	return i.Iterator.Next()
}
