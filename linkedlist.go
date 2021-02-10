package slice

// Interface type for a singly linked list node
type LinkedListNode interface {
	// Get the next node
	Next() LinkedListNode

	// Get the element
	Elem() interface{}
}

// Interface type for a doubly linked list node
type DoublyLinkedListNode interface{
	LinkedListNode

	// Get the previous node
	Prev() LinkedListNode
}

// Interface type for a linked list iterator
type LinkedListIterator interface {
	Iterator

	// Gets the node the iterator is currently pointed to
	Node() LinkedListNode
}

// Interface Slice type for a linked list
type LinkedList interface {
	Slice

	// Gets the node at the given index. Returns nil if the node couldn't be
	// found
	Node(i int) LinkedListNode
}