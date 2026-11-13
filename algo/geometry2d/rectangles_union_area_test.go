package geometry2d

import (
	"errors"
	"math"
	"math/rand/v2"
	"testing"
)

// Worked by hand with inclusion–exclusion (Stanley, Enumerative Combinatorics Vol. 1, 2nd ed., §2.1):
// [−3,3]×[0,4] has area 24, [0,9]×[−1,2] has 27, their overlap [0,3]×[0,2] has 6: 24 + 27 − 6 = 45.
func TestRectanglesUnionArea_example(t *testing.T) {
	got, err := RectanglesUnionArea(Rect{-3, 0, 3, 4}, Rect{0, -1, 9, 2})
	if err != nil || got != 45 {
		t.Errorf("RectanglesUnionArea(example) = (%d, %v), want (45, nil)", got, err)
	}
}

func TestRectanglesUnionArea_cases(t *testing.T) {
	const maxI, minI = math.MaxInt64, math.MinInt64
	tests := []struct {
		name   string
		r1, r2 Rect
		want   uint64
	}{
		{"disjoint", Rect{0, 0, 2, 2}, Rect{5, 5, 6, 7}, 6},
		{"touching edges", Rect{0, 0, 2, 2}, Rect{2, 0, 4, 2}, 8},
		{"touching corners", Rect{0, 0, 2, 2}, Rect{2, 2, 3, 3}, 5},
		{"one inside the other", Rect{0, 0, 10, 10}, Rect{2, 2, 3, 3}, 100},
		{"identical", Rect{1, 1, 4, 5}, Rect{1, 1, 4, 5}, 12},
		{"cross shape", Rect{0, 2, 6, 4}, Rect{2, 0, 4, 6}, 20},
		{"both degenerate", Rect{0, 0, 0, 5}, Rect{3, 3, 3, 3}, 0},
		{"one degenerate", Rect{1, 1, 1, 5}, Rect{0, 0, 2, 2}, 4},
		{"negative", Rect{-4, -4, -2, -1}, Rect{-3, -3, 0, 0}, 13}, // 6 + 9 − overlap [−3,−2]×[−3,−1] = 2
		{"largest square that fits", Rect{0, 0, math.MaxUint32, math.MaxUint32}, Rect{0, 0, 1, 1}, math.MaxUint32 * math.MaxUint32},
		{"extremes, mostly overlapping", Rect{minI, minI, maxI, minI + 1}, Rect{minI, minI, maxI, minI + 1}, math.MaxUint64},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RectanglesUnionArea(tt.r1, tt.r2)
			if err != nil || got != tt.want {
				t.Errorf("RectanglesUnionArea(%+v, %+v) = (%d, %v), want (%d, nil)", tt.r1, tt.r2, got, err, tt.want)
			}
			if got, err := RectanglesUnionArea(tt.r2, tt.r1); err != nil || got != tt.want {
				t.Errorf("RectanglesUnionArea swapped = (%d, %v), want (%d, nil)", got, err, tt.want)
			}
		})
	}
}

func TestRectanglesUnionArea_errors(t *testing.T) {
	const maxI, minI = math.MaxInt64, math.MinInt64
	ok := Rect{0, 0, 1, 1}
	inverted := []Rect{{2, 0, 1, 1}, {0, 2, 1, 1}}
	for _, r := range inverted {
		if _, err := RectanglesUnionArea(r, ok); !errors.Is(err, ErrInvertedRect) {
			t.Errorf("RectanglesUnionArea(%+v, ok) error = %v, want ErrInvertedRect", r, err)
		}
		if _, err := RectanglesUnionArea(ok, r); !errors.Is(err, ErrInvertedRect) {
			t.Errorf("RectanglesUnionArea(ok, %+v) error = %v, want ErrInvertedRect", r, err)
		}
	}
	overflows := [][2]Rect{
		{{minI, minI, maxI, maxI}, ok},                   // (2⁶⁴ − 1)²
		{{0, 0, 1 << 32, 1 << 32}, {0, 0, 0, 0}},         // exactly 2⁶⁴
		{{minI, 0, maxI, 1}, {minI, 2, maxI, 3}},         // 2·(2⁶⁴ − 1), disjoint
		{{0, 0, math.MaxUint32, math.MaxUint32 + 3}, ok}, // (2³² − 1)(2³² + 2) = 2⁶⁴ + 2³² − 2
	}
	for _, rr := range overflows {
		if got, err := RectanglesUnionArea(rr[0], rr[1]); !errors.Is(err, ErrOverflow) {
			t.Errorf("RectanglesUnionArea(%+v, %+v) = (%d, %v), want ErrOverflow", rr[0], rr[1], got, err)
		}
	}
}

// slowUnionArea counts the unit cells covered by either rectangle on a small grid.
func slowUnionArea(r1, r2 Rect, lo, hi int64) uint64 {
	in := func(r Rect, x, y int64) bool { return r.MinX <= x && x+1 <= r.MaxX && r.MinY <= y && y+1 <= r.MaxY }
	var cells uint64
	for x := lo; x < hi; x++ {
		for y := lo; y < hi; y++ {
			if in(r1, x, y) || in(r2, x, y) {
				cells++
			}
		}
	}
	return cells
}

// Oracle: unit cells counted on a 21×21 grid.
func TestRectanglesUnionArea_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	rect := func() Rect {
		x1, x2 := rng.Int64N(21)-10, rng.Int64N(21)-10
		y1, y2 := rng.Int64N(21)-10, rng.Int64N(21)-10
		return Rect{min(x1, x2), min(y1, y2), max(x1, x2), max(y1, y2)}
	}
	for trial := range 2000 {
		r1, r2 := rect(), rect()
		got, err := RectanglesUnionArea(r1, r2)
		if want := slowUnionArea(r1, r2, -10, 10); err != nil || got != want {
			t.Fatalf("trial %d: RectanglesUnionArea(%+v, %+v) = (%d, %v), want %d", trial, r1, r2, got, err, want)
		}
	}
}

func TestRectanglesUnionArea_noAllocation(t *testing.T) {
	r1, r2 := Rect{-3, 0, 3, 4}, Rect{0, -1, 9, 2}
	if allocs := testing.AllocsPerRun(100, func() { _, _ = RectanglesUnionArea(r1, r2) }); allocs != 0 {
		t.Errorf("RectanglesUnionArea allocated %v times per call, want 0", allocs)
	}
}

func BenchmarkRectanglesUnionArea(b *testing.B) {
	r1, r2 := Rect{-3, 0, 3, 4}, Rect{0, -1, 9, 2}
	b.Run("1", func(b *testing.B) {
		b.ReportAllocs()
		var s uint64
		for b.Loop() {
			a, _ := RectanglesUnionArea(r1, r2)
			s += a
		}
		_ = s
	})
}
