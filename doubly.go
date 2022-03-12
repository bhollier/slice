package slice

import "fmt"

type doublyNode[T any] struct {
	elem T
	next *doublyNode[T]
	prev *doublyNode[T]
}

func (n *doublyNode[T]) Next() LinkedListNode[T] {
	return n.next
}

func (n *doublyNode[T]) Prev() LinkedListNode[T] {
	return n.prev
}

func (n *doublyNode[T]) Get() T {
	return n.elem
}

func (n *doublyNode[T]) Set(elem T) {
	n.elem = elem
}

// Doubly is a Slice type, implemented as a doubly linked list
type Doubly[T any] struct {
	len   int
	start *doublyNode[T]
	end   *doublyNode[T]
}

// EmptyDoubly creates an empty Doubly Slice
func EmptyDoubly[T any]() Slice[T] {
	return Doubly[T]{}
}

// DoublyFrom creates a Doubly Slice from any type of slice
func DoublyFrom[T any](elems []T) Slice[T] {
	s := Doubly[T]{}

	// Iterate over the elements
	for i := 0; i < len(elems); i++ {
		// If the list is empty
		if i == 0 {
			s.start = new(doublyNode[T])
			s.start.elem = elems[i]
			s.end = s.start
		} else {
			s.end.next = new(doublyNode[T])
			s.end.next.prev = s.end
			s.end.next.elem = elems[i]
			s.end = s.end.next
		}
	}
	// Set the length
	s.len = len(elems)

	return s
}

func joinDoubly[T any](lhs, rhs Doubly[T]) Doubly[T] {
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
		rhs.start.prev = lhs.end
		lhs.end = rhs.end

		// Update the length
		lhs.len += rhs.len
		return lhs
	}
}

func (s Doubly[T]) Append(elems ...T) Slice[T] {
	return s.AppendSlice(DoublyFrom(elems))
}

func (s Doubly[T]) AppendSlice(elems Slice[T]) Slice[T] {
	// Try to convert the elements slice to a linked list
	rhs, ok := elems.(Doubly[T])
	// If it isn't a linked list
	if !ok {
		// Convert it to a linked list
		rhs = DoublyFrom(elems.ToGoSlice()).(Doubly[T])
	}

	// Join the linked lists
	return joinDoubly(s, rhs)
}

func (s Doubly[T]) Prepend(elems ...T) Slice[T] {
	return s.PrependSlice(DoublyFrom(elems))
}

func (s Doubly[T]) PrependSlice(elems Slice[T]) Slice[T] {
	// Try to convert the elements slice to a linked list
	lhs, ok := elems.(Doubly[T])
	// If it isn't a linked list
	if !ok {
		// Convert it to a linked list
		lhs = DoublyFrom(elems.ToGoSlice()).(Doubly[T])
	}

	// Connect the linked lists
	return joinDoubly(lhs, s)
}

func (s Doubly[T]) node(i int) *doublyNode[T] {
	if i < 0 || i >= s.len {
		panic(fmt.Sprintf("index [%d] out of range", i))
	}

	// If the index is in the first half, iterate forwards
	if i <= s.len/2 {
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

		// Otherwise, iterate backwards
	} else {
		// Create a counter
		ctr := s.len - 1
		// Iterate over the list
		node := s.end
		for node != nil {
			// If the node is a match
			if ctr == i {
				return node
			}
			// Decrement the counter and go to the next node
			ctr--
			node = node.prev
		}
	}
	return nil
}

func (s Doubly[T]) Node(i int) LinkedListNode[T] {
	return s.node(i)
}

func (s Doubly[T]) Slice(i, j int) Slice[T] {
	if j < i {
		panic(fmt.Sprintf("invalid slice index: %d > %d", i, j))
	}

	if j-i == 0 {
		return EmptyDoubly[T]()
	}

	// If the slice needs to be grown
	if j > s.Len() {
		s = s.Append(make([]T, j-s.Len())...).(Doubly[T])
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

func (s Doubly[T]) Get(i int) T {
	return s.Node(i).Get()
}

func (s Doubly[T]) Set(i int, elem T) {
	s.Node(i).Set(elem)
}

type doublyIterator[T any] struct {
	node  *doublyNode[T]
	start *doublyNode[T]
	end   *doublyNode[T]
}

func (i *doublyIterator[T]) HasNext() bool {
	return i.node != i.end
}

func (i *doublyIterator[T]) Next() bool {
	if i.HasNext() {
		i.node = i.node.next
		return true
	}
	return false
}

func (i *doublyIterator[T]) HasPrev() bool {
	return i.node != i.start
}

func (i *doublyIterator[T]) Prev() bool {
	if i.HasPrev() {
		i.node = i.node.prev
		return true
	}
	return false
}

func (i *doublyIterator[T]) Node() LinkedListNode[T] {
	return i.node
}

func (i *doublyIterator[T]) Get() T {
	return i.Node().Get()
}

func (i *doublyIterator[T]) Set(elem T) {
	i.Node().Set(elem)
}

func (s Doubly[T]) IterStart() Iterator[T] {
	if s.Len() == 0 {
		return &doublyIterator[T]{}
	}
	return &doublyIterator[T]{
		// Set the node as a new one that is one element out of bounds
		node:  &doublyNode[T]{next: s.start},
		start: s.start,
		end:   s.end,
	}
}

func (s Doubly[T]) IterEnd() Iterator[T] {
	if s.Len() == 0 {
		return &doublyIterator[T]{}
	}
	return &doublyIterator[T]{
		// Set the node as a new one that is one element out of bounds
		node:  &doublyNode[T]{prev: s.end},
		start: s.start,
		end:   s.end,
	}
}

func (s Doubly[T]) ReverseIterStart() Iterator[T] {
	return Reverse(s.IterEnd())
}

func (s Doubly[T]) ReverseIterEnd() Iterator[T] {
	return Reverse(s.IterStart())
}

func (s Doubly[T]) DeepCopy() Slice[T] {
	// Make sure to append the slice as a Go slice (to make sure it's a deep
	// copy)
	return EmptyDoubly[T]().Append(s.ToGoSlice()...)
}

func (s Doubly[T]) Len() int {
	return s.len
}

func (s Doubly[T]) Cap() int {
	return s.Len()
}

func (s Doubly[T]) ToGoSlice() []T {
	return ToGoSlice[T](s)
}
