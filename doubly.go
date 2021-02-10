package slice

import "fmt"

type doublyNode struct {
	elem interface{}
	next *doublyNode
	prev *doublyNode
}

func (n *doublyNode) Next() LinkedListNode {
	return n.next
}

func (n *doublyNode) Prev() LinkedListNode {
	return n.prev
}

func (n *doublyNode) Elem() interface{} {
	return n.elem
}

// Slice type, implemented as a doubly linked list
type Doubly struct {
	len int
	start *doublyNode
	end *doublyNode
}

// Create an empty Doubly Slice
func EmptyDoubly(len int) Slice {
	var s Slice
	s = &Doubly{}
	for i := 0; i < len; i++ {
		s = s.Append(nil)
	}
	return s
}

// Create a Doubly Slice from any type of slice
func DoublyFrom(slice interface{}) Slice {
	return appendNativeSliceToSlice(EmptyDoubly(0), slice)
}

func doublyFrom(elems []interface{}) *Doubly {
	s := Doubly{}

	// Iterate over the elements
	for i := 0; i < len(elems); i++ {
		// If the list is empty
		if i == 0 {
			s.start = new(doublyNode)
			s.start.elem = elems[i]
			s.end = s.start
		} else {
			s.end.next = new(doublyNode)
			s.end.next.prev = s.end
			s.end.next.elem = elems[i]
			s.end = s.end.next
		}
	}
	// Set the length
	s.len = len(elems)

	return &s
}

func (s* Doubly) join(rhs *Doubly) {
	// If the left linked list is empty
	if s.len == 0 {
		// The left linked list is just the right one
		*s = *rhs

		// If the right linked list is empty
	} else if rhs.len == 0 {
		// Nothing needs to be done

		// Otherwise
	} else {
		// Connect the pointers
		s.end.next = rhs.start
		rhs.start.prev = s.end
		s.end = rhs.end

		// Update the length
		s.len = s.len + rhs.len
	}
}

func (s *Doubly) Append(elems ...interface{}) Slice {
	return s.AppendSlice(doublyFrom(elems))
}

func (s *Doubly) AppendSlice(elems Slice) Slice {
	// Copy the linked list
	lhs := *s

	// Try to convert the elements slice to a linked list
	rhs, ok := elems.(*Doubly)
	// If it isn't a linked list
	if !ok {
		// Convert it to a linked list
		rhs = doublyFrom(elems.ToGoSlice())
	}

	// Join the linked lists
	lhs.join(rhs)

	// Return the slice
	return &lhs
}

func (s *Doubly) Prepend(elems ...interface{}) Slice {
	return s.PrependSlice(doublyFrom(elems))
}

func (s *Doubly) PrependSlice(elems Slice) Slice {
	// Try to convert the elements slice to a linked list
	lhs, ok := elems.(*Doubly)
	// If it isn't a linked list
	if !ok {
		// Convert it to a linked list
		lhs = doublyFrom(elems.ToGoSlice())
	}

	// Connect the linked lists
	lhs.join(s)

	// Return the linked list
	return lhs
}

func (s *Doubly) node(i int) *doublyNode {
	if i < 0 || i >= s.len {
		panic(fmt.Sprintf("index [%d] out of range", i))
	}

	// If the index is in the first half, iterate forwards
	if i <= s.len / 2 {
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

func (s *Doubly) Node(i int) LinkedListNode {
	return s.node(i)
}

func (s *Doubly) Slice(i, j int) Slice {
	// Create the slice
	slice := *s

	if j < i {
		panic(fmt.Sprintf("invalid slice index: %d > %d", i, j))
	}

	// If i is the start, the start is the start of the slice
	if i == 0 {slice.start = s.start
		// If i is the end, the start is the end of the slice
	} else if i == slice.len {slice.start = s.end
	} else {slice.start = s.node(i)}

	// If j is at the start, the end is the start of the slice
	if j == 0 {slice.end = s.start
		// If j is at the end, the end is the end of the slice
	} else if j == slice.len {slice.end = s.end
	} else {slice.end = s.node(j - 1)}

	// Set the length
	slice.len = j - i

	return &slice
}

func (s *Doubly) Index(i int) interface{} {
	return s.Node(i).Elem()
}

type doublyIterator struct {
	node *doublyNode
	start *doublyNode
	end *doublyNode
}

func (i *doublyIterator) HasNext() bool {
	return i.node != i.end
}

func (i *doublyIterator) Next() bool {
	if i.HasNext() {
		i.node = i.node.next
		return true
	}
	return false
}

func (i *doublyIterator) HasPrev() bool {
	return i.node != i.start
}

func (i *doublyIterator) Prev() bool {
	if i.HasPrev() {
		i.node = i.node.prev
		return true
	}
	return false
}

func (i *doublyIterator) Node() LinkedListNode {
	return i.node
}

func (i *doublyIterator) Elem() interface{} {
	return i.Node().Elem()
}

func (s *Doubly) IterStart() Iterator {
	return &doublyIterator{
		// Set the node as a new one that is one element out of bounds
		node: &doublyNode{next: s.start},
		start: s.start,
		end: s.end,
	}
}

func (s *Doubly) IterEnd() Iterator {
	return &doublyIterator{
		// Set the node as a new one that is one element out of bounds
		node: &doublyNode{prev: s.end},
		start: s.start,
		end: s.end,
	}
}

func (s *Doubly) ReverseIterStart() Iterator {
	return Reverse(s.IterEnd())
}

func (s *Doubly) ReverseIterEnd() Iterator {
	return Reverse(s.IterStart())
}

func (s *Doubly) DeepCopy() Slice {
	// Make sure to append the slice as a Go slice (to make sure it's a deep
	// copy)
	return EmptyDoubly(0).Append(s.ToGoSlice()...)
}

func (s *Doubly) Len() int {
	return s.len
}

func (s *Doubly) Cap() int {
	return s.Len()
}

func (s *Doubly) ToGoSlice() []interface{} {
	return ToGoSlice(s)
}