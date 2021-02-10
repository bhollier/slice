package slice

import "fmt"

type singlyNode struct {
	elem interface{}
	next *singlyNode
}

func (n *singlyNode) Next() LinkedListNode {
	return n.next
}

func (n *singlyNode) Elem() interface{} {
	return n.elem
}

// Slice type, implemented as a singly linked list
type Singly struct {
	len int
	start *singlyNode
	end *singlyNode
}

// Create an empty Singly Slice
func EmptySingly(len int) Slice {
	var s Slice
	s = &Singly{}
	for i := 0; i < len; i++ {
		s = s.Append(nil)
	}
	return s
}

// Create a Singly Slice from any type of slice
func SinglyFrom(slice interface{}) Slice {
	return appendNativeSliceToSlice(EmptySingly(0), slice)
}

func singlyFrom(elems []interface{}) *Singly {
	s := Singly{}

	// Iterate over the elements
	for i := 0; i < len(elems); i++ {
		// If the list is empty
		if i == 0 {
			s.start = new(singlyNode)
			s.start.elem = elems[i]
			s.end = s.start
		} else {
			s.end.next = new(singlyNode)
			s.end.next.elem = elems[i]
			s.end = s.end.next
		}
	}
	// Set the length
	s.len = len(elems)

	return &s
}

func (s* Singly) join(rhs *Singly) {
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
		s.end = rhs.end

		// Update the length
		s.len = s.len + rhs.len
	}
}

func (s *Singly) Append(elems ...interface{}) Slice {
	return s.AppendSlice(singlyFrom(elems))
}

func (s *Singly) AppendSlice(elems Slice) Slice {
	// Copy the linked list
	lhs := *s

	// Try to convert the elements slice to a linked list
	rhs, ok := elems.(*Singly)
	// If it isn't a linked list
	if !ok {
		// Convert it to a linked list
		rhs = singlyFrom(elems.ToGoSlice())
	}

	// Join the linked lists
	lhs.join(rhs)

	// Return the slice
	return &lhs
}

func (s *Singly) Prepend(elems ...interface{}) Slice {
	return s.PrependSlice(singlyFrom(elems))
}

func (s *Singly) PrependSlice(elems Slice) Slice {
	// Try to convert the elements slice to a linked list
	lhs, ok := elems.(*Singly)
	// If it isn't a linked list
	if !ok {
		// Convert it to a linked list
		lhs = singlyFrom(elems.ToGoSlice())
	}

	// Connect the linked lists
	lhs.join(s)

	// Return the linked list
	return lhs
}

func (s *Singly) node(i int) *singlyNode {
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

func (s *Singly) Node(i int) LinkedListNode {
	return s.node(i)
}

func (s *Singly) Slice(i, j int) Slice {
	// Create the slice
	slice := *s

	if j < i {panic(fmt.Sprintf("invalid slice index: %d > %d", i, j))}

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

func (s *Singly) Index(i int) interface{} {
	return s.Node(i).Elem()
}

type singlyIterator struct {
	node *singlyNode
	end *singlyNode
}

func (i *singlyIterator) HasNext() bool {
	return i.node != i.end
}

func (i *singlyIterator) Next() bool {
	if i.HasNext() {
		i.node = i.node.next
		return true
	}
	return false
}

// Always returns false
func (i *singlyIterator) HasPrev() bool {
	return false
}

// Always returns false
func (i *singlyIterator) Prev() bool {
	return false
}

// Gets the node the iterator is currently pointed to
func (i *singlyIterator) Node() LinkedListNode {
	return i.node
}

func (i *singlyIterator) Elem() interface{} {
	return i.Node().Elem()
}

func (s *Singly) IterStart() Iterator {
	return &singlyIterator{
		// Set the node as a new one that is one element out of bounds
		node: &singlyNode{next: s.start},
		end: s.end,
	}
}

func (s *Singly) IterEnd() Iterator {
	return &singlyIterator{}
}

func (s *Singly) ReverseIterStart() Iterator {
	return Reverse(s.IterEnd())
}

func (s *Singly) ReverseIterEnd() Iterator {
	return Reverse(s.IterStart())
}

func (s *Singly) DeepCopy() Slice {
	// Make sure to append the slice as a Go slice (to make sure it's a deep
	// copy)
	return EmptySingly(0).Append(s.ToGoSlice()...)
}

func (s *Singly) Len() int {
	return s.len
}

func (s *Singly) Cap() int {
	return s.Len()
}

func (s *Singly) ToGoSlice() []interface{} {
	return ToGoSlice(s)
}