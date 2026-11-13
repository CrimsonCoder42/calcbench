package interval

import (
	"errors"
	"math"
	"math/rand/v2"
	"slices"
	"testing"
)

// Worked example written for this repo: [3,8] [1,4] [10,12] against [2,7] and [2,11]. By hand [1,4]
// and [3,8] together cover [1,8], so [2,7] is covered; (8,10) is open, so [2,11] is not.
func TestCovers_workedExample(t *testing.T) {
	in := []Interval[int]{iv(3, 8), iv(1, 4), iv(10, 12)}
	if got, err := Covers(in, iv(2, 7)); err != nil || !got {
		t.Errorf("Covers([2,7]) = (%v, %v), want (true, nil)", got, err)
	}
	if got, err := Covers(in, iv(2, 11)); err != nil || got {
		t.Errorf("Covers([2,11]) = (%v, %v), want (false, nil)", got, err)
	}
}

func TestCovers_cases(t *testing.T) {
	tests := []struct {
		name   string
		in     []Interval[int]
		target Interval[int]
		want   bool
	}{
		{"no intervals", nil, iv(1, 2), false},
		{"exact", []Interval[int]{iv(1, 5)}, iv(1, 5), true},
		{"point target covered", []Interval[int]{iv(5, 5)}, iv(5, 5), true},
		{"point target missed", []Interval[int]{iv(1, 4), iv(6, 9)}, iv(5, 5), false},
		{"shared endpoint joins", []Interval[int]{iv(2, 4), iv(1, 2)}, iv(1, 4), true},
		{"integer neighbours leave a gap", []Interval[int]{iv(1, 2), iv(3, 4)}, iv(1, 4), false},
		{"starts too late", []Interval[int]{iv(2, 9)}, iv(1, 5), false},
		{"ends too early", []Interval[int]{iv(0, 4)}, iv(1, 5), false},
		{"irrelevant intervals before", []Interval[int]{iv(-9, -5), iv(0, 9)}, iv(1, 5), true},
		{"duplicates", []Interval[int]{iv(1, 3), iv(1, 3), iv(3, 6)}, iv(1, 6), true},
		{"contained ones do not extend", []Interval[int]{iv(1, 6), iv(2, 3)}, iv(1, 7), false},
		{"extremes", []Interval[int]{iv(math.MinInt, 0), iv(0, math.MaxInt)}, iv(math.MinInt, math.MaxInt), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := Covers(tt.in, tt.target); err != nil || got != tt.want {
				t.Errorf("Covers(%v, %v) = (%v, %v), want (%v, nil)", tt.in, tt.target, got, err, tt.want)
			}
		})
	}
}

func TestCovers_floats(t *testing.T) {
	in := []Interval[float64]{{0, 0.5}, {0.5, math.Inf(1)}}
	if got, _ := Covers(in, Interval[float64]{0.25, 1e300}); !got {
		t.Error("float union [0,+Inf] should cover [0.25, 1e300]")
	}
	if got, _ := Covers([]Interval[float64]{{0, 0.5}, {0.5000001, 1}}, Interval[float64]{0, 1}); got {
		t.Error("a tiny float gap must not count as covered")
	}
}

// Problem-site shape, in our own words: video clips of a recording; do they together span the whole
// recording?
func TestCovers_problemSiteShape(t *testing.T) {
	clips := []Interval[int]{iv(0, 2), iv(4, 6), iv(8, 10), iv(1, 9), iv(1, 5), iv(5, 9)}
	if got, _ := Covers(clips, iv(0, 10)); !got {
		t.Error("clips should span [0,10]")
	}
	if got, _ := Covers(clips[:3], iv(0, 10)); got {
		t.Error("three clips leave gaps")
	}
}

func TestCovers_errors(t *testing.T) {
	if _, err := Covers([]Interval[int]{iv(1, 2)}, iv(3, 1)); !errors.Is(err, ErrInverted) {
		t.Errorf("inverted target: err = %v", err)
	}
	if _, err := Covers([]Interval[int]{iv(1, 2), iv(5, 4)}, iv(1, 2)); !errors.Is(err, ErrInverted) {
		t.Errorf("inverted interval: err = %v", err)
	}
}

func TestCovers_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	const maxValue = 40
	for trial := range 400 {
		in := make([]Interval[int], rng.IntN(12))
		for i := range in {
			start := rng.IntN(maxValue)
			in[i] = iv(start, start+rng.IntN(min(10, maxValue-start+1)))
		}
		start := rng.IntN(maxValue)
		target := iv(start, start+rng.IntN(min(10, maxValue-start+1)))
		marked := mark(in, maxValue)
		want := true
		for p := 2 * target.Start; p <= 2*target.End; p++ {
			want = want && marked[p]
		}
		before := slices.Clone(in)
		if got, err := Covers(in, target); err != nil || got != want {
			t.Fatalf("trial %d: Covers(%v, %v) = (%v, %v), want %v", trial, in, target, got, err, want)
		}
		if !slices.Equal(in, before) {
			t.Fatalf("trial %d: input modified", trial)
		}
	}
}

func BenchmarkCovers(b *testing.B) {
	for _, size := range benchSizes {
		rng := rand.New(rand.NewPCG(1, 2))
		in := make([]Interval[int], size.n)
		for i := range in {
			in[i] = iv(10*i, 10*i+10) // a chain: the whole list is swept
		}
		rng.Shuffle(len(in), func(i, j int) { in[i], in[j] = in[j], in[i] })
		target := iv(0, 10*size.n)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _ = Covers(in, target)
			}
		})
	}
}
