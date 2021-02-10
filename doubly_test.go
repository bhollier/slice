package slice

import (
	"math/rand"
	"testing"
	"time"
)

func TestDoubly_Append(t *testing.T) {
	commonSliceAppendTest(t, EmptyDoubly(0))
	commonSliceAppendTest(t, DoublyFrom([]int{1}))
	commonSliceAppendTest(t, DoublyFrom([]int{1, 2}))
}

func TestDoubly_Prepend(t *testing.T) {
	commonSlicePrependTest(t, EmptyDoubly(0))
	commonSlicePrependTest(t, DoublyFrom([]int{1}))
	commonSlicePrependTest(t, DoublyFrom([]int{1, 2}))
}

func TestDoubly_Slice(t *testing.T) {
	commonSliceSliceTest(t, EmptyDoubly(0))
	commonSliceSliceTest(t, DoublyFrom([]int{1}))
	commonSliceSliceTest(t, DoublyFrom([]int{1, 2}))
}

func TestDoubly_Iter(t *testing.T) {
	commonSliceIterTest(t, EmptyDoubly(0))
	commonSliceIterTest(t, DoublyFrom([]int{1}))
	commonSliceIterTest(t, DoublyFrom([]int{1, 2}))
}

func TestDoubly_ReverseIter(t *testing.T) {
	commonSliceReverseIterTest(t, EmptyDoubly(0))
	commonSliceReverseIterTest(t, DoublyFrom([]int{1}))
	commonSliceReverseIterTest(t, DoublyFrom([]int{1, 2}))
}

// BENCHMARKING

func BenchmarkDoubly_Append(b *testing.B) {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	commonSliceAppendBenchmark(b, r, EmptyDoubly(0))
}

func BenchmarkDoubly_Prepend(b *testing.B) {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	commonSlicePrependBenchmark(b, r, EmptyDoubly(0))
}

func BenchmarkDoubly_Erase(b *testing.B) {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	commonSliceEraseBenchmark(b, r, EmptyDoubly(0))
}

func BenchmarkDoubly_Index(b *testing.B) {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	commonSliceIndexBenchmark(b, r, EmptyDoubly(0))
}

func BenchmarkDoubly_Iter(b *testing.B) {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	commonSliceIterBenchmark(b, r, EmptyDoubly(0))
}