package slice

// LinkedListNode is an interface type for a singly linked list node
type LinkedListNode[T any] interface {
	// Next gets the next node
	Next() LinkedListNode[T]

	// Get gets the element
	Get() T

	// Set sets the element
	Set(T)
}

// DoublyLinkedListNode is an interface type for a doubly linked list node
type DoublyLinkedListNode[T any] interface {
	LinkedListNode[T]

	// Prev gets the previous node
	Prev() DoublyLinkedListNode[T]
}

// LinkedListIterator is an interface type for a linked list iterator
type LinkedListIterator[T any] interface {
	Iterator[T]

	// Node gets the node the iterator is currently pointed to
	Node() LinkedListNode[T]
}

// LinkedList is an interface Slice type for a linked list
type LinkedList[T any] interface {
	Slice[T]

	// Node gets the node at the given index. Returns nil if the node couldn't be
	// found
	Node(i int) LinkedListNode[T]
}
