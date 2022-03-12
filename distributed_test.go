package slice

import (
	"math/rand"
	"testing"
	"time"
)

func TestDistributed_Append(t *testing.T) {
	commonSliceAppendTest(t, EmptyDistributed[int](0, 2))
	commonSliceAppendTest(t, DistributedFrom([]int{1}))
	commonSliceAppendTest(t, DistributedFrom([]int{1, 2}))
}

func TestDistributed_Prepend(t *testing.T) {
	commonSlicePrependTest(t, EmptyDistributed[int](0, 2))
	commonSlicePrependTest(t, DistributedFrom([]int{1}))
	commonSlicePrependTest(t, DistributedFrom([]int{1, 2}))
}

func TestDistributed_Insert(t *testing.T) {
	commonSliceInsertTest(t, EmptyDistributed[int](0, 2))
	commonSliceInsertTest(t, DistributedFrom([]int{1}))
	commonSliceInsertTest(t, DistributedFrom([]int{1, 2}))
}

func TestDistributed_Slice(t *testing.T) {
	commonSliceSliceTest(t, EmptyDistributed[int](0, 2))
	commonSliceSliceTest(t, DistributedFrom([]int{1}))
	commonSliceSliceTest(t, DistributedFrom([]int{1, 2}))
}

func TestDistributed_Iter(t *testing.T) {
	commonSliceIterTest(t, EmptyDistributed[int](0, 2))
	commonSliceIterTest(t, DistributedFrom([]int{1}))
	commonSliceIterTest(t, DistributedFrom([]int{1, 2}))
}

func TestDistributed_ReverseIter(t *testing.T) {
	commonSliceReverseIterTest(t, EmptyDistributed[int](0, 2))
	commonSliceReverseIterTest(t, DistributedFrom([]int{1}))
	commonSliceReverseIterTest(t, DistributedFrom([]int{1, 2}))
}

// BENCHMARKING

func BenchmarkDistributed_Append(b *testing.B) {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	commonSliceAppendBenchmark(b, r, EmptyDistributed[int](0, 0))
}

func BenchmarkDistributed_Prepend(b *testing.B) {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	commonSlicePrependBenchmark(b, r, EmptyDistributed[int](0, 0))
}

func BenchmarkDistributed_Erase(b *testing.B) {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	commonSliceEraseBenchmark(b, r, EmptyDistributed[int](0, 0))
}

func BenchmarkDistributed_Index(b *testing.B) {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	commonSliceIndexBenchmark(b, r, EmptyDistributed[int](0, 0))
}

func BenchmarkDistributed_Iter(b *testing.B) {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	commonSliceIterBenchmark(b, r, EmptyDistributed[int](0, 0))
}
