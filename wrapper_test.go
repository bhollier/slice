package slice

import (
	"math/rand"
	"testing"
	"time"
)

func TestWrapper_Append(t *testing.T) {
	commonSliceAppendTest(t, EmptySlice(0, 0))
	commonSliceAppendTest(t, Wrap([]int{1}))
	commonSliceAppendTest(t, Wrap([]int{1, 2}))
}

func TestWrapper_Prepend(t *testing.T) {
	commonSlicePrependTest(t, EmptySlice(0, 0))
	commonSlicePrependTest(t, Wrap([]int{1}))
	commonSlicePrependTest(t, Wrap([]int{1, 2}))
}

func TestWrapper_Slice(t *testing.T) {
	commonSliceSliceTest(t, EmptySlice(0, 0))
	commonSliceSliceTest(t, Wrap([]int{1}))
	commonSliceSliceTest(t, Wrap([]int{1, 2}))
}

func TestWrapper_Iter(t *testing.T) {
	commonSliceIterTest(t, EmptySlice(0, 0))
	commonSliceIterTest(t, Wrap([]int{1}))
	commonSliceIterTest(t, Wrap([]int{1, 2}))
}

func TestWrapper_ReverseIter(t *testing.T) {
	commonSliceReverseIterTest(t, EmptySlice(0, 0))
	commonSliceReverseIterTest(t, Wrap([]int{1}))
	commonSliceReverseIterTest(t, Wrap([]int{1, 2}))
}

// BENCHMARKING

func BenchmarkWrapper_Append(b *testing.B) {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	commonSliceAppendBenchmark(b, r, EmptySlice(0, 0))
}

func BenchmarkWrapper_Prepend(b *testing.B) {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	commonSlicePrependBenchmark(b, r, EmptySlice(0, 0))
}

func BenchmarkWrapper_Erase(b *testing.B) {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	commonSliceEraseBenchmark(b, r, EmptySlice(0, 0))
}

func BenchmarkWrapper_Index(b *testing.B) {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	commonSliceIndexBenchmark(b, r, EmptySlice(0, 0))
}

func BenchmarkWrapper_Iter(b *testing.B) {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	commonSliceIterBenchmark(b, r, EmptySlice(0, 0))
}