package slice

import (
	"math/rand"
	"testing"
	"time"
)

func TestSingly_Append(t *testing.T) {
	commonSliceAppendTest(t, EmptySingly(0))
	commonSliceAppendTest(t, SinglyFrom([]int{1}))
	commonSliceAppendTest(t, SinglyFrom([]int{1, 2}))
}

func TestSingly_Prepend(t *testing.T) {
	commonSlicePrependTest(t, EmptySingly(0))
	commonSlicePrependTest(t, SinglyFrom([]int{1}))
	commonSlicePrependTest(t, SinglyFrom([]int{1, 2}))
}

func TestSingly_Slice(t *testing.T) {
	commonSliceSliceTest(t, EmptySingly(0))
	commonSliceSliceTest(t, SinglyFrom([]int{1}))
	commonSliceSliceTest(t, SinglyFrom([]int{1, 2}))
}

func TestSingly_Iter(t *testing.T) {
	commonSliceIterTest(t, EmptySingly(0))
	commonSliceIterTest(t, SinglyFrom([]int{1}))
	commonSliceIterTest(t, SinglyFrom([]int{1, 2}))
}

// BENCHMARKING

func BenchmarkSingly_Append(b *testing.B) {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	commonSliceAppendBenchmark(b, r, EmptySingly(0))
}

func BenchmarkSingly_Prepend(b *testing.B) {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	commonSlicePrependBenchmark(b, r, EmptySingly(0))
}

func BenchmarkSingly_Erase(b *testing.B) {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	commonSliceEraseBenchmark(b, r, EmptySingly(0))
}

func BenchmarkSingly_Index(b *testing.B) {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	commonSliceIndexBenchmark(b, r, EmptySingly(0))
}

func BenchmarkSingly_Iter(b *testing.B) {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	commonSliceIterBenchmark(b, r, EmptySingly(0))
}