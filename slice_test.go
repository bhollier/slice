package slice

import (
	"github.com/stretchr/testify/assert"
	"math/rand"
	"reflect"
	"testing"
)

var intType = reflect.TypeOf((*int)(nil)).Elem()

func commonSliceIndexTest(t *testing.T, s Slice, index int, expected interface{}) {
	assert.Equal(t, expected, s.Index(index))
}

func commonSliceLenTest(t *testing.T, s Slice, expected int) {
	assert.Equal(t, expected, s.Len())
}

func commonSliceAppendTest(t *testing.T, s Slice) {
	s1 := s.Append(1)
	commonSliceLenTest(t, s1, s.Len() + 1)
	commonSliceIndexTest(t, s1, s1.Len() - 1, 1)

	s1 = s.Append(2, 3)
	commonSliceLenTest(t, s1, s.Len() + 2)
	commonSliceIndexTest(t, s1, s1.Len() - 2, 2)
	commonSliceIndexTest(t, s1, s1.Len() - 1, 3)
}

func commonSlicePrependTest(t *testing.T, s Slice) {
	s1 := s.Prepend(1)
	commonSliceLenTest(t, s1, s.Len() + 1)
	commonSliceIndexTest(t, s1, 0, 1)

	s1 = s.Prepend(2, 3)
	commonSliceLenTest(t, s1, s.Len() + 2)
	commonSliceIndexTest(t, s1, 0, 2)
	commonSliceIndexTest(t, s1, 1, 3)
}

func commonSliceSliceTest(t *testing.T, s Slice) {
	s1 := s.Append(2, 3, 4)
	commonSliceLenTest(t, s1, s.Len() + 3)
	slice := s1.Slice(s1.Len() - 3, s1.Len())
	commonSliceLenTest(t, slice, 3)
	commonSliceIndexTest(t, slice, 0, 2)
	commonSliceIndexTest(t, slice, 1, 3)
	commonSliceIndexTest(t, slice, 2, 4)

	emptySlice := s1.Slice(0, 0)
	commonSliceLenTest(t, emptySlice, 0)

	emptySlice = s1.Slice(s1.Len(), s1.Len())
	commonSliceLenTest(t, emptySlice, 0)

	cpy := slice.DeepCopy()
	erased := Erase(slice, 1)
	commonSliceLenTest(t, erased, 2)
	commonSliceIndexTest(t, erased, 0, 2)
	commonSliceIndexTest(t, erased, 1, 4)

	erased = EraseRange(cpy, 0, 1)
	commonSliceLenTest(t, erased, 1)
	commonSliceIndexTest(t, erased, 0, 4)
}

func commonSliceIterTest(t *testing.T, s Slice) {
	elems := []int{2, 3, 4}
	s1 := s.AppendSlice(Wrap(elems))
	iter := s1.Slice(s1.Len() - 3, s1.Len()).IterStart()
	i := 0
	for iter.Next() {
		assert.Equal(t, elems[i], iter.Elem())
		i++
	}
	assert.Equal(t, len(elems), i)
}

func commonSliceReverseIterTest(t *testing.T, s Slice) {
	elems := []int{2, 3, 4}
	s1 := s.AppendSlice(Wrap(elems))
	iter := s1.Slice(s1.Len() - 3, s1.Len()).ReverseIterStart()
	i := len(elems) - 1
	for iter.Next() {
		assert.Equal(t, elems[i], iter.Elem())
		i--
	}
	assert.Equal(t, -1, i)
}

// BENCHMARKING

const benchmarkMaxSliceInserts = 100
const benchmarkMinSliceLen = 10000
const benchmarkMinReadAmount = 1000

func commonSliceAppendBenchmark(b *testing.B, r *rand.Rand, s Slice) {
	for i := 0; i < b.N; i++ {
		end := make([]interface{}, atLeast(1, r.Intn(benchmarkMaxSliceInserts)))
		for k := 0; k < len(end); k++ {
			end[k] = r.Int()
		}
		s = s.Append(end...)
	}
}

func commonSlicePrependBenchmark(b *testing.B, r *rand.Rand, s Slice) {
	for i := 0; i < b.N; i++ {
		end := make([]interface{}, atLeast(1, r.Intn(benchmarkMaxSliceInserts)))
		for k := 0; k < len(end); k++ {
			end[k] = r.Int()
		}
		s = s.Prepend(end...)
	}
}

func addElems(b *testing.B, r *rand.Rand, s Slice) Slice {
	b.StopTimer()
	for k := 0; k < benchmarkMinSliceLen; k++ {
		s = s.Append(r.Int())
	}
	b.StartTimer()
	return s
}

func commonSliceEraseBenchmark(b *testing.B, r *rand.Rand, s Slice) {
	s = addElems(b, r, s)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if s.Len() == 0 {
			s = addElems(b, r, s)
		}
		s = Erase(s, r.Intn(s.Len()))
	}
}

func commonSliceIndexBenchmark(b *testing.B, r *rand.Rand, s Slice) {
	s = addElems(b, r, s)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		readAmount := atLeast(1, r.Intn(benchmarkMinReadAmount))
		start := r.Intn(s.Len() - readAmount)
		sum := 0
		for j := 0; j < readAmount; j++ {
			sum += s.Index(j + start).(int)
		}
	}
}

func commonSliceIterBenchmark(b *testing.B, r *rand.Rand, s Slice) {
	s = addElems(b, r, s)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		readAmount := atLeast(1, r.Intn(benchmarkMinReadAmount))
		start := r.Intn(s.Len() - readAmount)
		iter := s.Slice(start, start + readAmount).IterStart()
		sum := 0
		for iter.Next() {
			sum += iter.Elem().(int)
		}
	}
}