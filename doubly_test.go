package slice

import (
	"math/rand"
	"testing"
	"time"
)

func TestDoubly_Append(t *testing.T) {
	commonSliceAppendTest(t, EmptyDoubly[int]())
	commonSliceAppendTest(t, DoublyFrom([]int{1}))
	commonSliceAppendTest(t, DoublyFrom([]int{1, 2}))
}

func TestDoubly_Prepend(t *testing.T) {
	commonSlicePrependTest(t, EmptyDoubly[int]())
	commonSlicePrependTest(t, DoublyFrom([]int{1}))
	commonSlicePrependTest(t, DoublyFrom([]int{1, 2}))
}

func TestDoubly_Insert(t *testing.T) {
	commonSliceInsertTest(t, EmptyDoubly[int]())
	commonSliceInsertTest(t, DoublyFrom([]int{1}))
	commonSliceInsertTest(t, DoublyFrom([]int{1, 2}))
}

func TestDoubly_Slice(t *testing.T) {
	commonSliceSliceTest(t, EmptyDoubly[int]())
	commonSliceSliceTest(t, DoublyFrom([]int{1}))
	commonSliceSliceTest(t, DoublyFrom([]int{1, 2}))
}

func TestDoubly_Iter(t *testing.T) {
	commonSliceIterTest(t, EmptyDoubly[int]())
	commonSliceIterTest(t, DoublyFrom([]int{1}))
	commonSliceIterTest(t, DoublyFrom([]int{1, 2}))
}

func TestDoubly_ReverseIter(t *testing.T) {
	commonSliceReverseIterTest(t, EmptyDoubly[int]())
	commonSliceReverseIterTest(t, DoublyFrom([]int{1}))
	commonSliceReverseIterTest(t, DoublyFrom([]int{1, 2}))
}

// BENCHMARKING

func BenchmarkDoubly_Append(b *testing.B) {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	commonSliceAppendBenchmark(b, r, EmptyDoubly[int]())
}

func BenchmarkDoubly_Prepend(b *testing.B) {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	commonSlicePrependBenchmark(b, r, EmptyDoubly[int]())
}

func BenchmarkDoubly_Erase(b *testing.B) {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	commonSliceEraseBenchmark(b, r, EmptyDoubly[int]())
}

func BenchmarkDoubly_Index(b *testing.B) {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	commonSliceIndexBenchmark(b, r, EmptyDoubly[int]())
}

func BenchmarkDoubly_Iter(b *testing.B) {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	commonSliceIterBenchmark(b, r, EmptyDoubly[int]())
}
