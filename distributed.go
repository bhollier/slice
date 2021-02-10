package slice

import (
	"fmt"
	"unsafe"
)

// The default capacity of a distributed slice bucket, currently 256 bytes
const DefaultDistributedBucketCapacity =
	int(256 / unsafe.Sizeof(interface{}(nil)))
//	int(unsafe.Sizeof(cpu.CacheLinePad{}) /
//		unsafe.Sizeof(interface{}(nil)))

// The capacity of a bucket, currently the number of interface{}
// that can fit on a cache line
//const BucketCapacity = int(unsafe.Sizeof(cpu.CacheLinePad{}) /
//	unsafe.Sizeof(interface{}(nil)))

// Slice type, implemented as an array of buckets, similar to an unrolled linked
// list (see https://en.wikipedia.org/wiki/Unrolled_linked_list)
type Distributed struct {
	buckets []bucket
	bucketCap int
	start int
	end int
}

// Create an empty Distributed Slice. cap is the capacity of each bucket, not
// the capacity of the entire slice. If cap is 0, DefaultBucketCapacity is used
// instead
func EmptyDistributed(len, cap int) Slice {
	// If no capacity was given
	if cap == 0 {
		// Set it to the default
		cap = DefaultDistributedBucketCapacity
	}

	// Calculate the initial number of buckets
	numInitialBuckets := (len / cap) + 1
	s := &Distributed{
		buckets:   make([]bucket, numInitialBuckets),
		bucketCap: cap,
		start:     0,
		end:       0,
	}
	// Add the nodes
	for i := 0; i < numInitialBuckets; i++ {
		s.buckets[i] = make(bucket, s.bucketCap)
	}
	// If there is at least one element
	if len > 0 {
		// Set the end point
		s.end = (len - 1) % cap
	}

	return s
}

// Create a Distributed Slice from any type of slice
func DistributedFrom(slice interface{}) Slice {
	return appendNativeSliceToSlice(EmptyDistributed(0, 0), slice)
}

func (s *Distributed) Append(elems ...interface{}) Slice {
	return s.AppendSlice(wrap(elems))
}

func (s *Distributed) AppendSlice(elems Slice) Slice {
	// Copy the slice
	slice := *s

	// If the bucket is full
	if slice.end == s.bucketCap {
		// Add a new node
		slice.buckets = append(slice.buckets, make(bucket, s.bucketCap))
		// Reset the end point
		slice.end = 0
	}

	// Iterate over the elements
	iter := elems.IterStart()
	for iter.Next() {
		// Copy the element
		slice.buckets[len(slice.buckets) - 1][slice.end] = iter.Elem()
		// Move the end point forward
		slice.end++
		// If the bucket is full
		if slice.end == s.bucketCap {
			// Add a new node
			slice.buckets = append(slice.buckets, make(bucket, s.bucketCap))
			// Reset the end point
			slice.end = 0
		}
	}

	return &slice
}

func (s *Distributed) Prepend(elems ...interface{}) Slice {
	return s.PrependSlice(wrap(elems))
}

func (s *Distributed) PrependSlice(elems Slice) Slice {
	// Copy the slice
	slice := *s

	// Iterate over the elements
	iter := elems.IterEnd()
	for iter.Prev() {
		// Move the start point backwards
		slice.start--
		// If the start point is out of bounds
		if slice.start < 0 {
			// Add a new bucket
			slice.buckets = append([]bucket{make(bucket, s.bucketCap)}, slice.buckets...)
			// Reset the start point
			slice.start = slice.bucketCap - 1
		}

		// Copy the element
		slice.buckets[0][slice.start] = iter.Elem()
	}

	return &slice
}

func (s *Distributed) Slice(i, j int) Slice {
	// Copy the slice
	slice := *s

	slen := s.Len()
	if i < 0 || i > slen {
		panic(fmt.Sprintf("index [%d] out of range", i))
	}
	if j < 0 || j > slen {
		panic(fmt.Sprintf("index [%d] out of range", j))
	}

	// Calculate the real i and j index
	iIndex, jIndex := i + slice.start, j + slice.start

	// Calculate the start and end buckets
	bucketsStart, bucketsEnd := iIndex / slice.bucketCap, jIndex / slice.bucketCap

	// Make sure there is at least one bucket
	if bucketsStart == bucketsEnd {bucketsEnd++}

	// Reslice the buckets slice
	slice.buckets = slice.buckets[bucketsStart:bucketsEnd]

	// Set the start and end point
	slice.start = iIndex % slice.bucketCap
	slice.end = jIndex % slice.bucketCap

	// Return the slice
	return &slice
}

func (s *Distributed) Index(i int) interface{} {
	index := i + s.start
	return s.buckets[index / s.bucketCap][index % s.bucketCap]
}

type distributedIterator struct {
	slice Distributed
	bucketIndex int
	index int
}

func (i *distributedIterator) HasNext() bool {
	// Set the capacity of the current bucket
	bucketCap := i.slice.bucketCap
	// If this is the last bucket
	if i.bucketIndex == len(i.slice.buckets) - 1 {
		// Set the capacity as the end of the slice
		bucketCap = i.slice.end
	}

	// There is another element if the next index is still inside the bounds of
	// the bucket
	return i.index + 1 < bucketCap ||
		// Or if there is a next bucket
		i.bucketIndex + 1 < len(i.slice.buckets)
}

func (i *distributedIterator) Next() bool {
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

func (i *distributedIterator) HasPrev() bool {
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

func (i *distributedIterator) Prev() bool {
	if i.HasPrev() {
		// Go to the previous element
		i.index--
		// If the index is out of bounds. This if statement doesn't care if this
		// is the first bucket and therefore if index is equal to i.slice.start,
		// because that is checked by i.HasPrev anyway
		if i.index < 0 {
			// Reset the index
			i.index = i.slice.bucketCap
			// Go to the previous bucket
			i.bucketIndex--
		}
		return true
	}
	return false
}

func (i *distributedIterator) Elem() interface{} {
	return i.slice.buckets[i.bucketIndex][i.index]
}

func (s *Distributed) IterStart() Iterator {
	return &distributedIterator{
		slice: *s,
		bucketIndex: 0,
		index: s.start - 1,
	}
}

func (s *Distributed) IterEnd() Iterator {
	return &distributedIterator{
		slice: *s,
		bucketIndex: len(s.buckets) - 1,
		index: s.end,
	}
}

func (s *Distributed) ReverseIterStart() Iterator {
	return Reverse(s.IterEnd())
}

func (s *Distributed) ReverseIterEnd() Iterator {
	return Reverse(s.IterStart())
}

func (s *Distributed) DeepCopy() Slice {
	return EmptyDistributed(0, s.bucketCap).AppendSlice(s)
}

func (s *Distributed) Len() int {
	// If there is only one block
	if len(s.buckets) == 1 {
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

// Gets the total capacity of the slice, not the bucket capacity
func (s *Distributed) Cap() int {
	return len(s.buckets) * s.bucketCap
}

func (s *Distributed) ToGoSlice() []interface{} {
	return ToGoSlice(s)
}