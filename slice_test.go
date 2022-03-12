package slice

import (
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
)

func commonSliceGetTest(t *testing.T, s Slice[int], index int, expected int) {
	assert.Equal(t, expected, s.Get(index))
}

func commonSliceLenTest(t *testing.T, s Slice[int], expected int) {
	assert.Equal(t, expected, s.Len())
}

func commonSliceAppendTest(t *testing.T, s Slice[int]) {
	s1 := s.Append(1)
	commonSliceLenTest(t, s1, s.Len()+1)
	commonSliceGetTest(t, s1, s1.Len()-1, 1)

	s1 = s.Append(2, 3)
	commonSliceLenTest(t, s1, s.Len()+2)
	commonSliceGetTest(t, s1, s1.Len()-2, 2)
	commonSliceGetTest(t, s1, s1.Len()-1, 3)

	s1 = s.Append(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
	assert.Equal(t,
		[]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		s1.Slice(s1.Len()-10, s1.Len()).ToGoSlice())
}

func commonSlicePrependTest(t *testing.T, s Slice[int]) {
	s1 := s.Prepend(1)
	commonSliceLenTest(t, s1, s.Len()+1)
	commonSliceGetTest(t, s1, 0, 1)

	s1 = s.Prepend(3, 2)
	commonSliceLenTest(t, s1, s.Len()+2)
	commonSliceGetTest(t, s1, 0, 3)
	commonSliceGetTest(t, s1, 1, 2)

	s1 = s.Prepend(10, 9, 8, 7, 6, 5, 4, 3, 2, 1)
	assert.Equal(t,
		[]int{10, 9, 8, 7, 6, 5, 4, 3, 2, 1},
		s1.Slice(0, 10).ToGoSlice())
}

func commonSliceInsertTest(t *testing.T, s Slice[int]) {
	s1 := s.Append(1, 5)
	s1 = Insert(s1, s.Len()+1, 2)
	commonSliceLenTest(t, s1, s.Len()+3)
	commonSliceGetTest(t, s1, s.Len()+0, 1)
	commonSliceGetTest(t, s1, s.Len()+1, 2)
	commonSliceGetTest(t, s1, s.Len()+2, 5)

	s1 = InsertSlice(s1, s.Len()+2, Wrap([]int{3, 4}))
	commonSliceLenTest(t, s1, s.Len()+5)
	commonSliceGetTest(t, s1, s.Len()+0, 1)
	commonSliceGetTest(t, s1, s.Len()+1, 2)
	commonSliceGetTest(t, s1, s.Len()+2, 3)
	commonSliceGetTest(t, s1, s.Len()+3, 4)
	commonSliceGetTest(t, s1, s.Len()+4, 5)
}

func commonSliceSliceTest(t *testing.T, s Slice[int]) {
	s1 := s.Append(2, 3, 4)
	commonSliceLenTest(t, s1, s.Len()+3)
	slice := s1.Slice(s1.Len()-3, s1.Len())
	commonSliceLenTest(t, slice, 3)
	commonSliceGetTest(t, slice, 0, 2)
	commonSliceGetTest(t, slice, 1, 3)
	commonSliceGetTest(t, slice, 2, 4)

	emptySlice := s1.Slice(0, 0)
	commonSliceLenTest(t, emptySlice, 0)

	emptySlice = s1.Slice(s1.Len(), s1.Len())
	commonSliceLenTest(t, emptySlice, 0)

	cpy := slice.DeepCopy()
	erased := Erase(slice, 1)
	commonSliceLenTest(t, erased, 2)
	commonSliceGetTest(t, erased, 0, 2)
	commonSliceGetTest(t, erased, 1, 4)

	erased = EraseRange(cpy, 0, 1)
	commonSliceLenTest(t, erased, 1)
	commonSliceGetTest(t, erased, 0, 4)
}

func commonSliceEraseTest(t *testing.T, s Slice[int]) {
	// todo
}

func commonSliceIterTest(t *testing.T, s Slice[int]) {
	elems := []int{2, 3, 4}
	s1 := s.AppendSlice(Wrap(elems))
	iter := s1.Slice(s1.Len()-3, s1.Len()).IterStart()
	i := 0
	for iter.Next() {
		assert.Equal(t, elems[i], iter.Get())
		i++
	}
	assert.Equal(t, len(elems), i)
}

func commonSliceReverseIterTest(t *testing.T, s Slice[int]) {
	elems := []int{2, 3, 4}
	s1 := s.AppendSlice(Wrap(elems))
	iter := s1.Slice(s1.Len()-3, s1.Len()).ReverseIterStart()
	i := len(elems) - 1
	for iter.Next() {
		assert.Equal(t, elems[i], iter.Get())
		i--
	}
	assert.Equal(t, -1, i)
}

// BENCHMARKING

const benchmarkMaxSliceInserts = 100
const benchmarkMinSliceLen = 10000
const benchmarkMinReadAmount = 1000

func commonSliceAppendBenchmark(b *testing.B, r *rand.Rand, s Slice[int]) {
	for i := 0; i < b.N; i++ {
		end := make([]int, atLeast(1, r.Intn(benchmarkMaxSliceInserts)))
		for k := 0; k < len(end); k++ {
			end[k] = r.Int()
		}
		s = s.Append(end...)
	}
}

func commonSlicePrependBenchmark(b *testing.B, r *rand.Rand, s Slice[int]) {
	for i := 0; i < b.N; i++ {
		end := make([]int, atLeast(1, r.Intn(benchmarkMaxSliceInserts)))
		for k := 0; k < len(end); k++ {
			end[k] = r.Int()
		}
		s = s.Prepend(end...)
	}
}

func commonSliceInsertBenchmark(b *testing.B, r *rand.Rand, s Slice[int]) {
	for i := 0; i < b.N; i++ {
		end := make([]int, atLeast(1, r.Intn(benchmarkMaxSliceInserts)))
		for k := 0; k < len(end); k++ {
			end[k] = r.Int()
		}
		s = InsertSlice(s, r.Intn(s.Len()), Wrap(end))
	}
}

func addElems(b *testing.B, r *rand.Rand, s Slice[int]) Slice[int] {
	b.StopTimer()
	for k := 0; k < benchmarkMinSliceLen; k++ {
		s = s.Append(r.Int())
	}
	b.StartTimer()
	return s
}

func commonSliceEraseBenchmark(b *testing.B, r *rand.Rand, s Slice[int]) {
	s = addElems(b, r, s)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if s.Len() == 0 {
			s = addElems(b, r, s)
		}
		s = Erase(s, r.Intn(s.Len()))
	}
}

func commonSliceIndexBenchmark(b *testing.B, r *rand.Rand, s Slice[int]) {
	s = addElems(b, r, s)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		readAmount := atLeast(1, r.Intn(benchmarkMinReadAmount))
		start := r.Intn(s.Len() - readAmount)
		sum := 0
		for j := 0; j < readAmount; j++ {
			sum += s.Get(j + start)
		}
	}
}

func commonSliceIterBenchmark(b *testing.B, r *rand.Rand, s Slice[int]) {
	s = addElems(b, r, s)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		readAmount := atLeast(1, r.Intn(benchmarkMinReadAmount))
		start := r.Intn(s.Len() - readAmount)
		iter := s.Slice(start, start+readAmount).IterStart()
		sum := 0
		for iter.Next() {
			sum += iter.Get()
		}
	}
}
