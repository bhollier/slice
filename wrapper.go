package slice

// Wrapper is a Slice type, implemented as a very thin wrapper around a Go slice
type Wrapper[T any] []T

// Wrap creates a Slice from a Go slice
func Wrap[T any](slice []T) Slice[T] {
	return Wrapper[T](slice)
}

// EmptySlice creates an empty Slice, implemented as []T
func EmptySlice[T any](len, cap int) Slice[T] {
	return Wrap(make([]T, len, cap))
}

func (s Wrapper[T]) Append(elems ...T) Slice[T] {
	return Wrap(append(s, elems...))
}

func (s Wrapper[T]) AppendSlice(elems Slice[T]) Slice[T] {
	return s.Append(elems.ToGoSlice()...)
}

func (s Wrapper[T]) Prepend(elems ...T) Slice[T] {
	return Wrap(append(elems, s...))
}

func (s Wrapper[T]) PrependSlice(elems Slice[T]) Slice[T] {
	return s.Prepend(elems.ToGoSlice()...)
}

func (s Wrapper[T]) Slice(i, j int) Slice[T] {
	return Wrap(s[i:j])
}

func (s Wrapper[T]) Get(i int) T {
	return s[i]
}

func (s Wrapper[T]) Set(i int, elem T) {
	s[i] = elem
}

type wrapperIterator[T any] struct {
	slice []T
	index int
}

func (i *wrapperIterator[T]) HasNext() bool {
	return i.index+1 < len(i.slice)
}

func (i *wrapperIterator[T]) Next() bool {
	if i.HasNext() {
		i.index++
		return true
	}
	return false
}

func (i *wrapperIterator[T]) HasPrev() bool {
	return i.index > 0
}

func (i *wrapperIterator[T]) Prev() bool {
	if i.HasPrev() {
		i.index--
		return true
	}
	return false
}

func (i *wrapperIterator[T]) Get() T {
	return i.slice[i.index]
}

func (i *wrapperIterator[T]) Set(elem T) {
	i.slice[i.index] = elem
}

func IterStart[T any](slice []T) Iterator[T] {
	i := wrapperIterator[T]{
		slice: slice,
		index: -1,
	}
	return &i
}

func IterEnd[T any](slice []T) Iterator[T] {
	i := wrapperIterator[T]{
		slice: slice,
		index: len(slice),
	}
	return &i
}

func ReverseIterStart[T any](slice []T) Iterator[T] {
	return Reverse(IterEnd(slice))
}

func ReverseIterEnd[T any](slice []T) Iterator[T] {
	return Reverse(IterStart(slice))
}

func (s Wrapper[T]) IterStart() Iterator[T] {
	return IterStart(s)
}

func (s Wrapper[T]) IterEnd() Iterator[T] {
	return IterEnd(s)
}

func (s Wrapper[T]) ReverseIterStart() Iterator[T] {
	return ReverseIterStart(s)
}

func (s Wrapper[T]) ReverseIterEnd() Iterator[T] {
	return ReverseIterEnd(s)
}

func (s Wrapper[T]) DeepCopy() Slice[T] {
	// Create a new slice
	slice := make([]T, s.Len())
	// Copy the elements onto it
	copy(slice, s)
	// Return it (wrapped)
	return Wrap(slice)
}

func (s Wrapper[T]) Len() int {
	return len(s)
}

func (s Wrapper[T]) Cap() int {
	return cap(s)
}

func (s Wrapper[T]) ToGoSlice() []T {
	return s
}
