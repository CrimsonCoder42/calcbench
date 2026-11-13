package heap

import (
	"errors"
	"math"
	"math/rand/v2"
	"slices"
	"testing"
)

// Knuth's 16 keys (TAOCP Vol. 3, 2nd ed., ch. 5) fed in Knuth's order; the running lower and upper
// medians after each of the first five keys, by hand from the sorted prefix:
// 503 → 503,503 · 87 503 → 87,503 · 87 503 512 → 503,503 · 61 87 503 512 → 87,503 ·
// 61 87 503 512 908 → 503,503. After all 16 the medians are 509 and 512 (as ALG-SELECT-002).
func TestMedianStream_knuthExample(t *testing.T) {
	want := [][2]int{{503, 503}, {87, 503}, {503, 503}, {87, 503}, {503, 503}}
	var m MedianStream[int]
	for i, v := range knuth16 {
		if !m.Add(v) {
			t.Fatalf("Add(%d) = false", v)
		}
		lo, hi, err := m.Median()
		if err != nil {
			t.Fatal(err)
		}
		if i < len(want) && (lo != want[i][0] || hi != want[i][1]) {
			t.Errorf("after %d keys: (%d, %d), want %v", i+1, lo, hi, want[i])
		}
	}
	if lo, hi, _ := m.Median(); lo != 509 || hi != 512 || m.Len() != 16 {
		t.Errorf("after 16 keys: (%d, %d), Len %d; want (509, 512), 16", lo, hi, m.Len())
	}
}

func TestMedianStream_cases(t *testing.T) {
	tests := []struct {
		name   string
		in     []int
		lo, hi int
	}{
		{"one", []int{4}, 4, 4},
		{"two", []int{9, 2}, 2, 9},
		{"duplicates", []int{3, 3, 3, 1}, 3, 3},
		{"sorted", []int{1, 2, 3, 4, 5, 6}, 3, 4},
		{"reversed", []int{6, 5, 4, 3, 2, 1}, 3, 4},
		{"negative and zero", []int{-2, 0, -7}, -2, -2},
		{"extremes", []int{math.MaxInt, math.MinInt}, math.MinInt, math.MaxInt},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, _ := NewMedianStream[int](len(tt.in))
			for _, v := range tt.in {
				m.Add(v)
			}
			if lo, hi, err := m.Median(); err != nil || lo != tt.lo || hi != tt.hi {
				t.Errorf("Median(%v) = (%d, %d, %v), want (%d, %d)", tt.in, lo, hi, err, tt.lo, tt.hi)
			}
		})
	}
}

// Problem-site shape, in our own words: add values one at a time and ask for the median after each,
// averaging the two middles when the count is even (the caller's arithmetic, not ours).
func TestMedianStream_problemSiteShape(t *testing.T) {
	var m MedianStream[float64]
	wantMean := []float64{5, 7.5, 5, 6.5}
	for i, v := range []float64{5, 10, 2, 8} {
		m.Add(v)
		lo, hi, _ := m.Median()
		if got := (lo + hi) / 2; got != wantMean[i] {
			t.Errorf("after %d values mean of middles = %v, want %v", i+1, got, wantMean[i])
		}
	}
}

func TestMedianStream_floatsNaN(t *testing.T) {
	var m MedianStream[float64]
	for _, v := range []float64{math.NaN(), 1, math.Inf(1)} {
		m.Add(v)
	}
	// D-002: NaN is the least value, so the sorted order is NaN 1 +Inf and the median is 1.
	if lo, hi, _ := m.Median(); lo != 1 || hi != 1 {
		t.Errorf("Median = (%v, %v), want (1, 1)", lo, hi)
	}
}

func TestMedianStream_errors(t *testing.T) {
	var m MedianStream[int]
	if _, _, err := m.Median(); !errors.Is(err, ErrEmpty) {
		t.Errorf("empty Median err = %v, want ErrEmpty", err)
	}
	var nilStream *MedianStream[int]
	if nilStream.Add(1) || nilStream.Len() != 0 {
		t.Error("nil stream: Add stored a value")
	}
	if _, _, err := nilStream.Median(); !errors.Is(err, ErrEmpty) {
		t.Errorf("nil Median err = %v, want ErrEmpty", err)
	}
	if s, err := NewMedianStream[int](-3); err != nil || s.Len() != 0 {
		t.Errorf("NewMedianStream(-3) = (Len %d, %v), want an empty stream", s.Len(), err)
	}
}

// Rule L41: a capacity whose buffers exceed one allocation's limit is an error, never a makeslice panic.
func TestNewMedianStream_tooLarge(t *testing.T) {
	for _, c := range []int{1 << 60, math.MaxInt} {
		if s, err := NewMedianStream[int](c); !errors.Is(err, ErrTooLarge) || s != nil {
			t.Errorf("NewMedianStream(%d) = (%v, %v), want (nil, ErrTooLarge)", c, s, err)
		}
	}
}

func TestMedianStream_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 300 {
		var m MedianStream[int]
		var seen []int
		for step := range 1 + rng.IntN(60) {
			v := rng.IntN(40) - 20
			m.Add(v)
			i, _ := slices.BinarySearch(seen, v)
			seen = slices.Insert(seen, i, v)
			n := len(seen)
			lo, hi, err := m.Median()
			if err != nil || lo != seen[(n-1)/2] || hi != seen[n/2] || m.Len() != n {
				t.Fatalf("trial %d step %d: (%d, %d, %v), want (%d, %d)", trial, step, lo, hi, err, seen[(n-1)/2], seen[n/2])
			}
		}
	}
}

func TestMedianStream_noAllocWithinCapacity(t *testing.T) {
	m, _ := NewMedianStream[int](400) // room for the warm-up run and the measured run
	allocs := testing.AllocsPerRun(1, func() {
		for i := range 100 {
			m.Add(i * 7 % 13)
			_, _, _ = m.Median()
		}
	})
	if allocs != 0 {
		t.Errorf("Add/Median within capacity allocated %v times, want 0", allocs)
	}
}

func BenchmarkNewMedianStream(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_, _ = NewMedianStream[int](16)
	}
}

// n Adds into a presized stream: O(n log n).
func BenchmarkMedianStreamAdd(b *testing.B) {
	for _, size := range benchSizes {
		values := randomInts(size.n)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				m, _ := NewMedianStream[int](size.n)
				for _, v := range values {
					m.Add(v)
				}
			}
		})
	}
}

func BenchmarkMedianStreamMedian(b *testing.B) {
	var m MedianStream[int]
	for _, v := range randomInts(1_000) {
		m.Add(v)
	}
	b.ReportAllocs()
	for b.Loop() {
		_, _, _ = m.Median()
	}
}

func BenchmarkMedianStreamLen(b *testing.B) {
	var m MedianStream[int]
	m.Add(1)
	b.ReportAllocs()
	for b.Loop() {
		_ = m.Len()
	}
}
