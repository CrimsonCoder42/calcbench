package window

import (
	"errors"
	"math"
	"math/rand/v2"
	"slices"
	"testing"
)

// naiveMovingMedian is the oracle: copy and sort each window, take the middle (or the mean of the two
// middles by the same midpoint rule). O(n·w log w).
func naiveMovingMedian(values []float64, w int) []float64 {
	var out []float64
	buf := make([]float64, w)
	for end := w - 1; end < len(values); end++ {
		copy(buf, values[end-w+1:end+1])
		slices.Sort(buf)
		if w%2 == 1 {
			out = append(out, buf[w/2])
		} else {
			out = append(out, midpoint(buf[w/2-1], buf[w/2]))
		}
	}
	return out
}

// Hand-worked running median of span 3 (the running median smoother of Härdle & Steiger 1995 with span 3;
// the paper's own data are not reproduced, so the windows are derived here by hand):
// 2 80 6 3 1 9 4 → [2 80 6]=6, [80 6 3]=6, [6 3 1]=3, [3 1 9]=3, [1 9 4]=4.
// The spike 80 is removed, which is why a median smoother is used instead of a moving average.
func TestMovingMedian_handWorkedSpan3(t *testing.T) {
	got, err := MovingMedian(nil, []float64{2, 80, 6, 3, 1, 9, 4}, 3)
	if err != nil {
		t.Fatal(err)
	}
	if want := []float64{6, 6, 3, 3, 4}; !slices.Equal(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestMovingMedian_cases(t *testing.T) {
	tests := []struct {
		name   string
		values []float64
		w      int
		want   []float64
	}{
		{"empty", nil, 1, nil},
		{"window longer than data", []float64{1, 2}, 3, nil},
		{"one value", []float64{7}, 1, []float64{7}},
		{"window 1 is identity", []float64{3, -1, 4}, 1, []float64{3, -1, 4}},
		{"two values even window", []float64{1, 4}, 2, []float64{2.5}},
		{"even window mean of middles", []float64{5, 1, 3, 9, 7}, 4, []float64{4, 5}},
		{"duplicates", []float64{2, 2, 2, 2, 2}, 3, []float64{2, 2, 2}},
		{"sorted", []float64{1, 2, 3, 4, 5, 6}, 3, []float64{2, 3, 4, 5}},
		{"reversed", []float64{6, 5, 4, 3, 2, 1}, 3, []float64{5, 4, 3, 2}},
		{"negative and zero", []float64{-3, 0, -1, -8, 0}, 3, []float64{-1, -1, -1}},
		{"whole data one window", []float64{9, 1, 5, 3}, 4, []float64{4}},
		{"extremes even window", []float64{math.MaxFloat64, math.MaxFloat64}, 2, []float64{math.MaxFloat64}},
		{"opposite extremes", []float64{-math.MaxFloat64, math.MaxFloat64}, 2, []float64{0}},
		{"smallest subnormal", []float64{5e-324, 5e-324, 5e-324}, 3, []float64{5e-324}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MovingMedian(nil, tt.values, tt.w)
			if err != nil {
				t.Fatal(err)
			}
			if !slices.Equal(got, tt.want) {
				t.Errorf("MovingMedian(%v, %d) = %v, want %v", tt.values, tt.w, got, tt.want)
			}
		})
	}
}

// Problem-site shape, in our own words: report the median of every run of k neighbours; with an even k
// the answer is the average of the two central readings.
func TestMovingMedian_problemSiteShape(t *testing.T) {
	got, _ := MovingMedian(nil, []float64{4, 8, 1, 7, 2, 6}, 2)
	if want := []float64{6, 4.5, 4, 4.5, 4}; !slices.Equal(got, want) {
		t.Errorf("k=2: %v, want %v", got, want)
	}
	got, _ = MovingMedian(nil, []float64{4, 8, 1, 7, 2, 6}, 5)
	if want := []float64{4, 6}; !slices.Equal(got, want) {
		t.Errorf("k=5: %v, want %v", got, want)
	}
}

func TestMovingMedian_appendsToDst(t *testing.T) {
	dst := []float64{-99}
	got, err := MovingMedian(dst, []float64{1, 3, 2}, 3)
	if err != nil || !slices.Equal(got, []float64{-99, 2}) {
		t.Errorf("got %v, %v; want [-99 2]", got, err)
	}
}

func TestMovingMedian_doesNotModifyInput(t *testing.T) {
	values := []float64{5, 1, 4, 2, 3}
	before := slices.Clone(values)
	if _, err := MovingMedian(nil, values, 2); err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(values, before) {
		t.Errorf("values modified: %v, want %v", values, before)
	}
}

func TestMovingMedian_errors(t *testing.T) {
	dst := []float64{1}
	for _, w := range []int{0, -1, math.MinInt} {
		got, err := MovingMedian(dst, []float64{1, 2}, w)
		if !errors.Is(err, ErrWindowLen) || len(got) != 1 {
			t.Errorf("windowLen %d: (%v, %v), want ErrWindowLen and dst unchanged", w, got, err)
		}
	}
	for _, bad := range []float64{math.NaN(), math.Inf(1), math.Inf(-1)} {
		got, err := MovingMedian(dst, []float64{1, bad, 2}, 2)
		if !errors.Is(err, ErrNotFinite) || len(got) != 1 {
			t.Errorf("value %v: (%v, %v), want ErrNotFinite and dst unchanged", bad, got, err)
		}
	}
	if got, err := MovingMedian(nil, []float64{1}, math.MaxInt); err != nil || len(got) != 0 {
		t.Errorf("huge window: (%v, %v), want no windows", got, err)
	}
}

// The oracle covers heavy ties (values from a small set), odd and even windows, and long inputs whose
// stale entries force heap rebuilds (n up to 40·w).
func TestMovingMedian_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 400 {
		w := 1 + rng.IntN(12)
		n := rng.IntN(40 * w)
		values := make([]float64, n)
		spread := 1 + rng.IntN(20)
		for i := range values {
			values[i] = float64(rng.IntN(spread) - spread/2)
		}
		if trial%3 == 0 { // monotone runs push stale entries deep into one heap
			slices.Sort(values)
		}
		got, err := MovingMedian(nil, values, w)
		if err != nil {
			t.Fatal(err)
		}
		if want := naiveMovingMedian(values, w); !slices.Equal(got, want) {
			t.Fatalf("trial %d (w=%d, n=%d): got %v, want %v", trial, w, n, got, want)
		}
	}
}

// Exact float comparison is sound: a median is one of the inputs (no arithmetic) and an even window's
// mean is computed by the same midpoint function in both, from the same two middle values.
func TestMovingMedian_oracleFloats(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for range 200 {
		w := 1 + rng.IntN(30)
		values := make([]float64, rng.IntN(300))
		for i := range values {
			values[i] = 2000*rng.Float64() - 1000
		}
		got, _ := MovingMedian(nil, values, w)
		if want := naiveMovingMedian(values, w); !slices.Equal(got, want) {
			t.Fatalf("w=%d: got %v, want %v", w, got, want)
		}
	}
}

func BenchmarkMovingMedian(b *testing.B) {
	for _, size := range benchSizes {
		values := benchFloats(size.n)
		dst := make([]float64, 0, size.n)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				dst, _ = MovingMedian(dst[:0], values, benchWindow)
			}
		})
	}
}

// Growth in w at fixed n: the log w factor.
func BenchmarkMovingMedianWindow(b *testing.B) {
	values := benchFloats(100_000)
	dst := make([]float64, 0, len(values))
	for _, w := range []struct {
		name string
		w    int
	}{{"w=10", 10}, {"w=1e3", 1_000}, {"w=1e4", 10_000}} {
		b.Run(w.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				dst, _ = MovingMedian(dst[:0], values, w.w)
			}
		})
	}
}
