package slice

import "reflect"

// Interface type for an iterator
type Iterator interface {
	// Returns whether there is a next element
	HasNext() bool

	// Goes to the next element, then returns false if the end of the slice
	// was reached
	Next() bool

	// Returns whether there is a previous element
	HasPrev() bool

	// Goes to the previous element, then returns false if the start of the
	// slice was reached
	Prev() bool

	// Gets the element the iterator is currently pointed to
	Elem() interface{}
}

// Interface type for a generic data structure that behaves like a []interface
// type
type Slice interface {
	// Copies the given element(s) onto the end of the slice. This function is
	// roughly equivalent to `append(slice, elems...)`
	Append(...interface{}) Slice

	// Copies the elements in the given slice onto the end of this slice
	AppendSlice(Slice) Slice

	// Copies the given elements onto the start of the slice. This function is
	// roughly equivalent to `append(elems, slice...)`
	Prepend(...interface{}) Slice

	// Copies the elements in the given slice onto the start of this slice
	PrependSlice(Slice) Slice

	// Gets a subset of the slice. This function is roughly equivalent to
	// `slice[i:j]`
	Slice(int, int) Slice

	// Gets the element at the given index. This function is roughly equivalent
	// to `slice[i]`
	Index(int) interface{}

	// Creates an iterator, pointed to the first element
	IterStart() Iterator

	// Creates a reverse iterator, pointed to the first element
	ReverseIterStart() Iterator

	// Creates an iterator, pointed to the last elements
	IterEnd() Iterator

	// Creates a reverse iterator, pointed to the last element
	ReverseIterEnd() Iterator

	// Creates a deep copy of the slice
	DeepCopy() Slice

	// Gets the slice's length. This function is roughly equivalent to
	// `len(slice)`
	Len() int

	// Gets the slice's capacity. This function is roughly equivalent to
	// `cap(slice)`
	Cap() int

	// Converts the slice to an []interface{} type
	ToGoSlice() []interface{}
}

type bucket []interface{}

// Essentially math.Max for ints
func atLeast(a, b int) int {
	if a > b {return a}
	return b
}

// Appends an interface{} to the given slice. intf must be some kind of slice
func appendNativeSliceToSlice(s Slice, slice interface{}) Slice {
	switch v := slice.(type) {
	case []interface{}:
		return s.Append(v...)
	default:
		val := reflect.ValueOf(slice)
		// Make sure the slice is a slice
		if val.Kind() != reflect.Slice {
			panic("given slice is not slice type")
		}
		// Iterate over the elements of the slice
		for i := 0; i < val.Len(); i++ {
			// Add it
			s = s.Append(val.Index(i).Interface())
		}
		return s
	}
}

// Converts a slice to an []interface{} type
func ToGoSlice(s Slice) []interface{} {
	// Create a slice
	slice := make([]interface{}, 0, s.Len())
	// Iterate over the elements
	iter := s.IterStart()
	for iter.Next() {
		// Add the element to the slice
		slice = append(slice, iter.Elem())
	}
	return slice
}

// Returns a slice where the element at the given index is erased. Equivalent to
// `append(s[:index], s[index + 1:]...)`. Warning: this function does not copy
// s, so the contents of s can be (and probably will be) modified. The
// Slice.Erase() function should be used instead of this function where
// possible, as it can be faster
func Erase(s Slice, index int) Slice {
	return s.Slice(0, index).AppendSlice(s.Slice(index + 1, s.Len()))
}

// Returns a slice where the range of elements are erased. Equivalent to
// `append(s[:i], s[j + 1:]...)`. Warning: this function does not copy s, so the
// contents of s can be modified. The Slice.EraseRange function should be used
// instead of this function where possible, as it can be faster
func EraseRange(s Slice, i, j int) Slice {
	return s.Slice(0, i).AppendSlice(s.Slice(j + 1, s.Len()))
}

// A reverse iterator, essentially an inverted iterator
type ReverseIterator struct {
	Iterator
}

// Create a reverse iterator
func Reverse(i Iterator) Iterator {
	// If the given iterator is a reverse iterator
	reverse, ok := i.(*ReverseIterator)
	if ok {
		// The reverse iterator is just the unwrapped iterator
		return reverse.Iterator
	}
	// Otherwise return a new iterator
	return &ReverseIterator{i}
}

func (i *ReverseIterator) HasNext() bool {
	return i.Iterator.HasPrev()
}

func (i *ReverseIterator) Next() bool {
	return i.Iterator.Prev()
}

func (i *ReverseIterator) HasPrev() bool {
	return i.Iterator.HasNext()
}

func (i *ReverseIterator) Prev() bool {
	return i.Iterator.Next()
}