package slice

import "fmt"

type singlyNode[T any] struct {
	elem T
	next *singlyNode[T]
}

func (n *singlyNode[T]) Next() LinkedListNode[T] {
	return n.next
}

func (n *singlyNode[T]) Get() T {
	return n.elem
}

func (n *singlyNode[T]) Set(elem T) {
	n.elem = elem
}

// Singly is a Slice type, implemented as a singly linked list
type Singly[T any] struct {
	len   int
	start *singlyNode[T]
	end   *singlyNode[T]
}

// EmptySingly creates an empty Singly Slice
func EmptySingly[T any]() Slice[T] {
	return Singly[T]{}
}

// SinglyFrom creates a Singly Slice from any type of slice
func SinglyFrom[T any](elems []T) Slice[T] {
	s := Singly[T]{}

	// Iterate over the elements
	for i := 0; i < len(elems); i++ {
		// If the list is empty
		if i == 0 {
			s.start = new(singlyNode[T])
			s.start.elem = elems[i]
			s.end = s.start
		} else {
			s.end.next = new(singlyNode[T])
			s.end.next.elem = elems[i]
			s.end = s.end.next
		}
	}
	// Set the length
	s.len = len(elems)

	return s
}

func joinSingly[T any](lhs, rhs Singly[T]) Singly[T] {
	// If the left linked list is empty
	if lhs.len == 0 {
		// Just return the right list
		return rhs

		// If the right linked list is empty
	} else if rhs.len == 0 {
		// Just return the left list
		return lhs

		// Otherwise
	} else {
		// Connect the pointers from the right to the left
		lhs.end.next = rhs.start
		lhs.end = rhs.end

		// Update the length
		lhs.len = lhs.len + rhs.len
		return lhs
	}
}

func (s Singly[T]) Append(elems ...T) Slice[T] {
	return s.AppendSlice(SinglyFrom(elems))
}

func (s Singly[T]) AppendSlice(elems Slice[T]) Slice[T] {
	// Try to convert the elements slice to a linked list
	rhs, ok := elems.(Singly[T])
	// If it isn't a linked list
	if !ok {
		// Convert it to a linked list
		rhs = SinglyFrom(elems.ToGoSlice()).(Singly[T])
	}

	// Connect the linked lists
	return joinSingly(s, rhs)
}

func (s Singly[T]) Prepend(elems ...T) Slice[T] {
	return s.PrependSlice(SinglyFrom(elems))
}

func (s Singly[T]) PrependSlice(elems Slice[T]) Slice[T] {
	// Try to convert the elements slice to a linked list
	lhs, ok := elems.(Singly[T])
	// If it isn't a linked list
	if !ok {
		// Convert it to a linked list
		lhs = SinglyFrom(elems.ToGoSlice()).(Singly[T])
	}

	// Connect the linked lists
	return joinSingly(lhs, s)
}

func (s Singly[T]) node(i int) *singlyNode[T] {
	if i < 0 || i >= s.len {
		panic(fmt.Sprintf("index [%d] out of range", i))
	}

	// Create a counter
	ctr := 0
	// Iterate over the list
	node := s.start
	for node != nil {
		// If the node is a match
		if ctr == i {
			return node
		}
		// Increment the counter and go to the next node
		ctr++
		node = node.next
	}

	return nil
}

func (s Singly[T]) Node(i int) LinkedListNode[T] {
	return s.node(i)
}

func (s Singly[T]) Slice(i, j int) Slice[T] {
	if j < i {
		panic(fmt.Sprintf("invalid slice index: %d > %d", i, j))
	}

	if j-i == 0 {
		return EmptySingly[T]()
	}

	// If the slice needs to be grown
	if j > s.Len() {
		s = s.Append(make([]T, j-s.Len())...).(Singly[T])
	}

	original := s

	// If i is the start
	if i == 0 {
		// The start is the start of the slice, nothing needs to be done

		// If i is the end, the start is the end of the slice
	} else if i == s.len {
		s.start = s.end
	} else {
		s.start = s.node(i)
	}

	// If j is at the start, the end is the start of the slice
	if j == 0 {
		s.end = s.start
		// If j is at the end
	} else if j == s.len {
		// The end is the end of the slice, nothing needs to be done

	} else {
		s.end = original.node(j - 1)
	}

	// Set the length
	s.len = j - i

	return s
}

func (s Singly[T]) Get(i int) T {
	return s.Node(i).Get()
}

func (s Singly[T]) Set(i int, elem T) {
	s.Node(i).Set(elem)
}

type singlyIterator[T any] struct {
	node *singlyNode[T]
	end  *singlyNode[T]
}

func (i *singlyIterator[T]) HasNext() bool {
	return i.node != i.end
}

func (i *singlyIterator[T]) Next() bool {
	if i.HasNext() {
		i.node = i.node.next
		return true
	}
	return false
}

// HasPrev always returns false
func (i *singlyIterator[T]) HasPrev() bool {
	return false
}

// Prev always returns false
func (i *singlyIterator[T]) Prev() bool {
	return false
}

// Node gets the node the iterator is currently pointed to
func (i *singlyIterator[T]) Node() LinkedListNode[T] {
	return i.node
}

func (i *singlyIterator[T]) Get() T {
	return i.Node().Get()
}

func (i *singlyIterator[T]) Set(elem T) {
	i.Node().Set(elem)
}

func (s Singly[T]) IterStart() Iterator[T] {
	if s.Len() == 0 {
		return &singlyIterator[T]{}
	}
	return &singlyIterator[T]{
		// Set the node as a new one that is one element out of bounds
		node: &singlyNode[T]{next: s.start},
		end:  s.end,
	}
}

func (s Singly[T]) IterEnd() Iterator[T] {
	return &singlyIterator[T]{}
}

func (s Singly[T]) ReverseIterStart() Iterator[T] {
	return Reverse(s.IterEnd())
}

func (s Singly[T]) ReverseIterEnd() Iterator[T] {
	return Reverse(s.IterStart())
}

func (s Singly[T]) DeepCopy() Slice[T] {
	// Make sure to append the slice as a Go slice (to make sure it's a deep
	// copy)
	return EmptySingly[T]().Append(s.ToGoSlice()...)
}

func (s Singly[T]) Len() int {
	return s.len
}

func (s Singly[T]) Cap() int {
	return s.Len()
}

func (s Singly[T]) ToGoSlice() []T {
	return ToGoSlice[T](s)
}
