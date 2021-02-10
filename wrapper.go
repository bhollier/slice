package slice

import "reflect"

// Slice type, implemented as a very thin wrapper around a Go slice
type Wrapper []interface{}

// Create a Slice from any type of slice, and store it as an []interface{} type
func Wrap(slice interface{}) Slice {
	s := Wrapper{}
	switch v := slice.(type) {
	case []interface{}:
		s = v
	default:
		val := reflect.ValueOf(slice)
		// Make sure the slice is a slice
		if val.Kind() != reflect.Slice {
			panic("given slice is not slice type")
		}
		s = make([]interface{}, val.Len(), val.Cap())
		// Iterate over the elements of the slice
		for i := 0; i < val.Len(); i++ {
			// Copy it
			s[i] = val.Index(i).Interface()
		}
	}
	return &s
}

func wrap(slice []interface{}) *Wrapper {
	s := Wrapper(slice)
	return &s
}

// Create an empty Slice, implemented as an interface{} slice
func EmptySlice(len, cap int) Slice {
	return wrap(make([]interface{}, len, cap))
}

func (s *Wrapper) Append(elems ...interface{}) Slice {
	return wrap(append(*s, elems...))
}

func (s *Wrapper) AppendSlice(elems Slice) Slice {
	return s.Append(elems.ToGoSlice()...)
}

func (s *Wrapper) Prepend(elems ...interface{}) Slice {
	return wrap(append(elems, *s...))
}

func (s *Wrapper) PrependSlice(elems Slice) Slice {
	return s.Prepend(elems.ToGoSlice()...)
}

func (s *Wrapper) Slice(i, j int) Slice {
	return wrap((*s)[i:j])
}

func (s *Wrapper) Index(i int) interface{} {
	return (*s)[i]
}

type wrapperIterator struct {
	slice []interface{}
	index int
}

func (i *wrapperIterator) HasNext() bool {
	return i.index + 1 < len(i.slice)
}

func (i *wrapperIterator) Next() bool {
	if i.HasNext() {
		i.index++
		return true
	}
	return false
}

func (i *wrapperIterator) HasPrev() bool {
	return i.index > 0
}

func (i *wrapperIterator) Prev() bool {
	if i.HasPrev() {
		i.index--
		return true
	}
	return false
}

func (i *wrapperIterator) Elem() interface{} {
	return i.slice[i.index]
}

func IterStart(slice []interface{}) Iterator {
	i := wrapperIterator{
		slice: slice,
		index: -1,
	}
	return &i
}

func IterEnd(slice []interface{}) Iterator {
	i := wrapperIterator{
		slice: slice,
		index: len(slice),
	}
	return &i
}

func ReverseIterStart(slice []interface{}) Iterator {
	return Reverse(IterEnd(slice))
}

func ReverseIterEnd(slice []interface{}) Iterator {
	return Reverse(IterStart(slice))
}

func (s *Wrapper) IterStart() Iterator {
	return IterStart(*s)
}

func (s *Wrapper) IterEnd() Iterator {
	return IterEnd(*s)
}

func (s *Wrapper) ReverseIterStart() Iterator {
	return ReverseIterStart(*s)
}

func (s *Wrapper) ReverseIterEnd() Iterator {
	return ReverseIterEnd(*s)
}

func (s *Wrapper) DeepCopy() Slice {
	// Create a new slice
	slice := make([]interface{}, s.Len())
	// Copy the elements onto it
	copy(slice, *s)
	// Return it (wrapped)
	return wrap(slice)
}

func (s *Wrapper) Len() int {
	return len(*s)
}

func (s *Wrapper) Cap() int {
	return cap(*s)
}

func (s *Wrapper) ToGoSlice() []interface{} {
	return *s
}