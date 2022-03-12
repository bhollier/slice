package slice

import (
	"math/rand"
	"testing"
	"time"
)

func BenchmarkBaselineSlice_Append(b *testing.B) {
	s := make([]int, 0, 0)
	r := rand.New(rand.NewSource(time.Now().Unix()))
	for i := 0; i < b.N; i++ {
		n := atLeast(1, r.Intn(benchmarkMaxSliceInserts))
		for k := 0; k < n; k++ {
			s = append(s, r.Int())
		}
	}
}

func BenchmarkBaselineInterfaceSlice_Append(b *testing.B) {
	s := make([]interface{}, 0, 0)
	r := rand.New(rand.NewSource(time.Now().Unix()))
	for i := 0; i < b.N; i++ {
		n := atLeast(1, r.Intn(benchmarkMaxSliceInserts))
		for k := 0; k < n; k++ {
			s = append(s, r.Int())
		}
	}
}

func BenchmarkBaselineSlice_Prepend(b *testing.B) {
	s := make([]int, 0, 0)
	r := rand.New(rand.NewSource(time.Now().Unix()))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		n := atLeast(1, r.Intn(benchmarkMaxSliceInserts))
		lhs := make([]int, n, n+len(s))
		for k := 0; k < n; k++ {
			lhs[k] = r.Int()
		}
		s = append(lhs, s...)
	}
}

func BenchmarkBaselineInterfaceSlice_Prepend(b *testing.B) {
	s := make([]interface{}, 0, 0)
	r := rand.New(rand.NewSource(time.Now().Unix()))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		n := atLeast(1, r.Intn(benchmarkMaxSliceInserts))
		lhs := make([]interface{}, n, n+len(s))
		for k := 0; k < n; k++ {
			lhs[k] = r.Int()
		}
		s = append(lhs, s...)
	}
}

func intSliceAddElems(b *testing.B, r *rand.Rand, s []int) []int {
	b.StopTimer()
	for k := 0; k < benchmarkMinSliceLen; k++ {
		s = append(s, r.Int())
	}
	b.StartTimer()
	return s
}

func interfaceSliceAddElems(b *testing.B, r *rand.Rand, s []interface{}) []interface{} {
	b.StopTimer()
	for k := 0; k < benchmarkMinSliceLen; k++ {
		s = append(s, r.Int())
	}
	b.StartTimer()
	return s
}

func BenchmarkBaselineSlice_Erase(b *testing.B) {
	s := make([]int, 0, benchmarkMinSliceLen)
	r := rand.New(rand.NewSource(time.Now().Unix()))
	s = intSliceAddElems(b, r, s)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if len(s) == 0 {
			s = intSliceAddElems(b, r, s)
		}
		elemToRemove := r.Intn(len(s))
		s = append(s[:elemToRemove], s[elemToRemove+1:]...)
	}
}

func BenchmarkBaselineInterfaceSlice_Erase(b *testing.B) {
	s := make([]interface{}, 0, benchmarkMinSliceLen)
	r := rand.New(rand.NewSource(time.Now().Unix()))
	s = interfaceSliceAddElems(b, r, s)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if len(s) == 0 {
			s = interfaceSliceAddElems(b, r, s)
		}
		elemToRemove := r.Intn(len(s))
		s = append(s[:elemToRemove], s[elemToRemove+1:]...)
	}
}

func BenchmarkBaselineSlice_Index(b *testing.B) {
	s := make([]int, 0, benchmarkMinSliceLen)
	r := rand.New(rand.NewSource(time.Now().Unix()))
	s = intSliceAddElems(b, r, s)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		readAmount := atLeast(1, r.Intn(benchmarkMinReadAmount))
		start := r.Intn(len(s) - readAmount)
		sum := 0
		for j := 0; j < readAmount; j++ {
			sum += s[j+start]
		}
	}
}

func BenchmarkBaselineInterfaceSlice_Index(b *testing.B) {
	s := make([]interface{}, 0, benchmarkMinSliceLen)
	r := rand.New(rand.NewSource(time.Now().Unix()))
	s = interfaceSliceAddElems(b, r, s)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		readAmount := r.Intn(benchmarkMinReadAmount)
		start := r.Intn(len(s) - readAmount)
		sum := 0
		for j := 0; j < readAmount; j++ {
			sum += s[j+start].(int)
		}
	}
}

func BenchmarkBaselineSlice_Iter(b *testing.B) {
	s := make([]int, 0, benchmarkMinSliceLen)
	r := rand.New(rand.NewSource(time.Now().Unix()))
	s = intSliceAddElems(b, r, s)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		readAmount := r.Intn(benchmarkMinReadAmount)
		start := r.Intn(len(s) - readAmount)
		sum := 0
		for _, n := range s[start : start+readAmount] {
			sum += n
		}
	}
}

func BenchmarkBaselineInterfaceSlice_Iter(b *testing.B) {
	s := make([]interface{}, 0, benchmarkMinSliceLen)
	r := rand.New(rand.NewSource(time.Now().Unix()))
	s = interfaceSliceAddElems(b, r, s)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		readAmount := r.Intn(benchmarkMinReadAmount)
		start := r.Intn(len(s) - readAmount)
		sum := 0
		for _, n := range s[start : start+readAmount] {
			sum += n.(int)
		}
	}
}
