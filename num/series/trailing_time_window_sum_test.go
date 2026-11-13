package series

import (
	"errors"
	"fmt"
	"math"
	"math/rand/v2"
	"testing"
	"time"
)

// parity: ledger RollingSum, source case 1 (trailing 7-day minutes over irregular sessions, New York).
func TestTrailingTimeWindowSum_parity(t *testing.T) {
	loc := newYork(t)
	at := func(m, d, h, mi int) time.Time { return time.Date(2026, time.Month(m), d, h, mi, 0, 0, loc) }
	pts := []Point{{at(10, 1, 7, 0), 30}, {at(10, 2, 18, 0), 45}, {at(10, 2, 19, 10), 20},
		{at(10, 5, 7, 0), 60}, {at(10, 8, 6, 30), 40}, {at(10, 9, 7, 0), 30}, {at(10, 15, 7, 0), 50}}
	got, err := TrailingTimeWindowSum(pts, 7*24*time.Hour)
	if err != nil {
		t.Fatal(err)
	}
	want := []WindowTotal{{30, 1}, {75, 2}, {95, 3}, {155, 4}, {195, 5}, {195, 5}, {80, 2}}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("i=%d: got %v, want %v", i, got[i], want[i])
		}
	}
}

func TestTrailingTimeWindowSum_cases(t *testing.T) {
	t0 := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	h := func(n int) time.Time { return t0.Add(time.Duration(n) * time.Hour) }
	if got, err := TrailingTimeWindowSum(nil, time.Hour); err != nil || len(got) != 0 {
		t.Errorf("empty: %v, %v", got, err)
	}
	// The left edge is open: a point exactly one window back has left.
	got, _ := TrailingTimeWindowSum([]Point{{h(0), 1}, {h(1), 2}}, time.Hour)
	if got[1] != (WindowTotal{2, 1}) {
		t.Errorf("open left edge: %v", got[1])
	}
	// Equal times: each point's window holds the points up to and including itself.
	got, _ = TrailingTimeWindowSum([]Point{{h(0), 1}, {h(0), 2}, {h(0), -4}}, time.Hour)
	if got[0] != (WindowTotal{1, 1}) || got[1] != (WindowTotal{3, 2}) || got[2] != (WindowTotal{-1, 3}) {
		t.Errorf("equal times: %v", got)
	}
	// Compensation: 1e16 enters and leaves; the small values survive.
	got, _ = TrailingTimeWindowSum([]Point{{h(0), 1e16}, {h(1), 1}, {h(2), 1}, {h(3), 1}}, 150*time.Minute)
	if got[3] != (WindowTotal{3, 3}) {
		t.Errorf("compensated: %v, want {3 3}", got[3])
	}
}

func TestTrailingTimeWindowSum_errors(t *testing.T) {
	t0 := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	cases := []struct {
		name string
		pts  []Point
		w    time.Duration
		want error
	}{
		{"zero window", []Point{{t0, 1}}, 0, ErrOutOfRange},
		{"negative window", []Point{{t0, 1}}, -time.Hour, ErrOutOfRange},
		{"unsorted", []Point{{t0, 1}, {t0.Add(-time.Second), 1}}, time.Hour, ErrUnsorted},
		{"NaN", []Point{{t0, math.NaN()}}, time.Hour, ErrNotFinite},
		{"Inf", []Point{{t0, math.Inf(1)}}, time.Hour, ErrNotFinite},
		{"overflow", []Point{{t0, math.MaxFloat64}, {t0.Add(time.Second), math.MaxFloat64}}, time.Hour, ErrOverflow},
	}
	for _, c := range cases {
		if _, err := TrailingTimeWindowSum(c.pts, c.w); !errors.Is(err, c.want) {
			t.Errorf("%s: err = %v, want %v", c.name, err, c.want)
		}
	}
}

// Oracle: brute-force window sums on random irregular series (fixed seed).
func TestTrailingTimeWindowSum_oracle(t *testing.T) {
	r := rand.New(rand.NewPCG(1, 2))
	t0 := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	for trial := 0; trial < 300; trial++ {
		n := 1 + r.IntN(60)
		pts := make([]Point, n)
		at := t0
		for i := range pts {
			at = at.Add(time.Duration(r.IntN(48)) * time.Hour)
			pts[i] = Point{At: at, Value: r.Float64() * 200}
		}
		w := time.Duration(1+r.IntN(14*24)) * time.Hour
		got, err := TrailingTimeWindowSum(pts, w)
		if err != nil {
			t.Fatal(err)
		}
		for i := range pts {
			var s float64
			c := 0
			for j := 0; j <= i; j++ {
				if pts[j].At.After(pts[i].At.Add(-w)) {
					s += pts[j].Value
					c++
				}
			}
			// Tolerance 1e-9: the brute force sums naively, the function with compensation.
			if got[i].Count != c || math.Abs(got[i].Sum-s) > 1e-9 {
				t.Fatalf("trial %d i=%d: got %v, want {%v %d}", trial, i, got[i], s, c)
			}
		}
	}
}

func BenchmarkTrailingTimeWindowSum(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		pts := benchPoints(n)
		b.Run(fmt.Sprint(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				if _, err := TrailingTimeWindowSum(pts, 24*time.Hour); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
