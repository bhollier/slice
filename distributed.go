package slice

import (
	"fmt"
	"unsafe"
)

// DefaultDistributedBucketCapacity is the default capacity of a distributed slice bucket,
// currently 256 bytes
const DefaultDistributedBucketCapacity = 256

//	int(unsafe.Sizeof(cpu.CacheLinePad{}) /
//		unsafe.Sizeof(interface{}(nil)))

// The capacity of a bucket, currently the number of interface{}
// that can fit on a cache line
// const BucketCapacity = int(unsafe.Sizeof(cpu.CacheLinePad{}) /
//	unsafe.Sizeof(interface{}(nil)))

type bucket[T any] []T

// Distributed is a Slice type, implemented as an array of buckets, similar to an unrolled linked
// list (see https://en.wikipedia.org/wiki/Unrolled_linked_list)
type Distributed[T any] struct {
	buckets   []bucket[T]
	bucketCap int
	start     int
	end       int
}

// EmptyDistributed creates an empty Distributed Slice. cap is the capacity of each bucket, not
// the capacity of the entire slice. If cap is 0, DefaultDistributedBucketCapacity is used
// instead
func EmptyDistributed[T any](len, cap int) Slice[T] {
	// If no capacity was given
	if cap == 0 {
		// Set it to the default
		var t T
		cap = int(DefaultDistributedBucketCapacity / unsafe.Sizeof(t))
	}

	// Calculate the initial number of buckets
	numInitialBuckets := (len + cap - 1) / cap
	s := Distributed[T]{
		buckets:   make([]bucket[T], numInitialBuckets),
		bucketCap: cap,
		start:     0,
		end:       0,
	}
	// Add the nodes
	for i := 0; i < numInitialBuckets; i++ {
		s.buckets[i] = make(bucket[T], s.bucketCap)
	}
	// If there is at least one element
	if len > 0 {
		// Set the end point
		s.end = (len - 1) % cap
	}

	return s
}

// DistributedFrom creates a Distributed Slice from a Go slice
func DistributedFrom[T any](slice []T) Slice[T] {
	return EmptyDistributed[T](0, 0).Append(slice...)
}

func (s Distributed[T]) Append(elems ...T) Slice[T] {
	return s.AppendSlice(Wrap(elems))
}

func (s Distributed[T]) AppendSlice(elems Slice[T]) Slice[T] {
	// If the slice is empty
	if len(s.buckets) == 0 {
		// Add a node
		s.buckets = append(s.buckets, make(bucket[T], s.bucketCap))
		// Reset the start and end points
		s.start = 0
		s.end = 0
	}

	// If the bucket is full
	if s.end == s.bucketCap {
		// Add a new node
		s.buckets = append(s.buckets, make(bucket[T], s.bucketCap))
		// Reset the end point
		s.end = 0
	}

	// Iterate over the elements
	iter := elems.IterStart()
	for iter.Next() {
		// Copy the element
		s.buckets[len(s.buckets)-1][s.end] = iter.Get()
		// Move the end point forward
		s.end++
		// If the bucket is full
		if s.end == s.bucketCap {
			// Add a new node
			s.buckets = append(s.buckets, make(bucket[T], s.bucketCap))
			// Reset the end point
			s.end = 0
		}
	}

	return s
}

func (s Distributed[T]) Prepend(elems ...T) Slice[T] {
	return s.PrependSlice(Wrap(elems))
}

func (s Distributed[T]) PrependSlice(elems Slice[T]) Slice[T] {
	// If the slice is empty
	if len(s.buckets) == 0 {
		// Add a node
		s.buckets = append(s.buckets, make(bucket[T], s.bucketCap))
		// Set the start and end points
		s.start = s.bucketCap
		s.end = s.bucketCap
	}

	// Iterate over the elements
	iter := elems.IterEnd()
	for iter.Prev() {
		// Move the start point backwards
		s.start--
		// If the start point is out of bounds
		if s.start < 0 {
			// Add a new bucket
			s.buckets = append([]bucket[T]{make(bucket[T], s.bucketCap)}, s.buckets...)
			// Reset the start point
			s.start = s.bucketCap - 1
		}

		// Copy the element
		s.buckets[0][s.start] = iter.Get()
	}

	return s
}

func (s Distributed[T]) Slice(i, j int) Slice[T] {
	// If the slice needs to be grown
	if j > s.Len() {
		s = s.Append(make([]T, j-s.Len())...).(Distributed[T])
	}

	if i < 0 {
		panic(fmt.Sprintf("index [%d] out of range", i))
	}
	if j < 0 {
		panic(fmt.Sprintf("index [%d] out of range", j))
	}

	// Calculate the real i and j index
	iIndex, jIndex := i+s.start, j+s.start

	// Calculate the start and end buckets
	bucketsStart := iIndex / s.bucketCap
	bucketsEnd := (jIndex + s.bucketCap - 1) / s.bucketCap

	// Reslice the buckets slice
	s.buckets = s.buckets[bucketsStart:bucketsEnd]

	// Set the start and end point
	s.start = iIndex % s.bucketCap
	s.end = ((jIndex - 1) % s.bucketCap) + 1

	// Return the slice
	return s
}

func (s Distributed[T]) Get(i int) T {
	index := i + s.start
	return s.buckets[index/s.bucketCap][index%s.bucketCap]
}

func (s Distributed[T]) Set(i int, elem T) {
	index := i + s.start
	s.buckets[index/s.bucketCap][index%s.bucketCap] = elem
}

type distributedIterator[T any] struct {
	slice       Distributed[T]
	bucketIndex int
	index       int
}

func (i *distributedIterator[T]) HasNext() bool {
	// There is no next if the slice is empty
	if len(i.slice.buckets) == 0 {
		return false
	}

	// Set the capacity of the current bucket
	bucketCap := i.slice.bucketCap
	// If this is the last bucket
	if i.bucketIndex == len(i.slice.buckets)-1 {
		// Set the capacity as the end of the slice
		bucketCap = i.slice.end
	}

	// There is another element if the next index is still inside the bounds of
	// the bucket
	return i.index+1 < bucketCap ||
		// Or if there is a next bucket
		i.bucketIndex+1 < len(i.slice.buckets)
}

func (i *distributedIterator[T]) Next() bool {
	if i.HasNext() {
		// Go to the next element
		i.index++
		// If the index is out of bounds. This if statement doesn't care if this
		// is the last bucket and therefore if index is equal to i.slice.end,
		// because that is checked by i.HasNext anyway
		if i.index == i.slice.bucketCap {
			// Reset the index
			i.index = 0
			// Go to the next bucket
			i.bucketIndex++
		}
		return true
	}
	return false
}

func (i *distributedIterator[T]) HasPrev() bool {
	// There is no next if the slice is empty
	if len(i.slice.buckets) == 0 {
		return false
	}

	// Set the start point of the current bucket
	start := 0
	// If this is the first bucket
	if i.bucketIndex == 0 {
		// Set the start as the start of the slice
		start = i.slice.start
	}

	// There is another element if the next index is still inside the bounds of
	// the bucket
	return i.index > start ||
		// Or if there is a previous bucket
		i.bucketIndex > 0
}

func (i *distributedIterator[T]) Prev() bool {
	if i.HasPrev() {
		// Go to the previous element
		i.index--
		// If the index is out of bounds. This if statement doesn't care if this
		// is the first bucket and therefore if index is equal to i.slice.start,
		// because that is checked by i.HasPrev anyway
		if i.index < 0 {
			// Reset the index
			i.index = i.slice.bucketCap - 1
			// Go to the previous bucket
			i.bucketIndex--
		}
		return true
	}
	return false
}

func (i *distributedIterator[T]) Get() T {
	return i.slice.buckets[i.bucketIndex][i.index]
}

func (i *distributedIterator[T]) Set(elem T) {
	i.slice.buckets[i.bucketIndex][i.index] = elem
}

func (s Distributed[T]) IterStart() Iterator[T] {
	return &distributedIterator[T]{
		slice:       s,
		bucketIndex: 0,
		index:       s.start - 1,
	}
}

func (s Distributed[T]) IterEnd() Iterator[T] {
	return &distributedIterator[T]{
		slice:       s,
		bucketIndex: len(s.buckets) - 1,
		index:       s.end,
	}
}

func (s Distributed[T]) ReverseIterStart() Iterator[T] {
	return Reverse(s.IterEnd())
}

func (s Distributed[T]) ReverseIterEnd() Iterator[T] {
	return Reverse(s.IterStart())
}

func (s Distributed[T]) DeepCopy() Slice[T] {
	return EmptyDistributed[T](0, s.bucketCap).AppendSlice(s)
}

func (s Distributed[T]) Len() int {
	// If there are no buckets
	if len(s.buckets) == 0 {
		return 0

		// If there is only one block
	} else if len(s.buckets) == 1 {
		// The length is the number of nodes between the start and end
		return s.end - s.start
	} else {
		// The length is the number of elements in the first node
		return (s.bucketCap - s.start) +
			// Plus the number of elements in the middle nodes
			((len(s.buckets) - 2) * s.bucketCap) +
			// Plus the number in the final block
			s.end
	}
}

// Cap gets the total capacity of the slice, not the bucket capacity
func (s Distributed[T]) Cap() int {
	return len(s.buckets) * s.bucketCap
}

func (s Distributed[T]) ToGoSlice() []T {
	return ToGoSlice[T](s)
}
