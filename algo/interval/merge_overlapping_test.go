package interval

import (
	"errors"
	"math"
	"math/rand/v2"
	"slices"
	"testing"
)

var benchSizes = []struct {
	name string
	n    int
}{{"1e3", 1_000}, {"1e4", 10_000}, {"1e5", 100_000}}

// iv is shorthand for an int interval in test tables.
func iv(start, end int) Interval[int] { return Interval[int]{Start: start, End: end} }

// Worked example written for this repo (no source prints one for this algorithm). Sorted by Start the
// input is [1,4] [3,5] [6,6] [9,12] [12,13] [20,25]. By hand: [1,4] and [3,5] overlap, giving [1,5];
// [6,6] is a single point apart from it; [9,12] and [12,13] share 12, giving [9,13]; [20,25] stands alone.
func TestMergeOverlapping_workedExample(t *testing.T) {
	in := []Interval[int]{iv(9, 12), iv(1, 4), iv(3, 5), iv(12, 13), iv(20, 25), iv(6, 6)}
	want := []Interval[int]{iv(1, 5), iv(6, 6), iv(9, 13), iv(20, 25)}
	got, err := MergeOverlapping(in)
	if err != nil || !slices.Equal(got, want) {
		t.Errorf("MergeOverlapping(%v) = (%v, %v), want (%v, nil)", in, got, err, want)
	}
}

func TestMergeOverlapping_cases(t *testing.T) {
	tests := []struct {
		name string
		in   []Interval[int]
		want []Interval[int]
	}{
		{"nil", nil, []Interval[int]{}},
		{"empty", []Interval[int]{}, []Interval[int]{}},
		{"one interval", []Interval[int]{iv(2, 5)}, []Interval[int]{iv(2, 5)}},
		{"one point", []Interval[int]{iv(3, 3)}, []Interval[int]{iv(3, 3)}},
		{"two disjoint", []Interval[int]{iv(5, 6), iv(1, 2)}, []Interval[int]{iv(1, 2), iv(5, 6)}},
		{"shared endpoint merges", []Interval[int]{iv(1, 3), iv(3, 5)}, []Interval[int]{iv(1, 5)}},
		{"integer neighbours stay apart", []Interval[int]{iv(1, 2), iv(3, 4)}, []Interval[int]{iv(1, 2), iv(3, 4)}},
		{"contained", []Interval[int]{iv(1, 10), iv(3, 4)}, []Interval[int]{iv(1, 10)}},
		{"duplicates", []Interval[int]{iv(2, 4), iv(2, 4), iv(2, 4)}, []Interval[int]{iv(2, 4)}},
		{"same start, longer second", []Interval[int]{iv(2, 3), iv(2, 8)}, []Interval[int]{iv(2, 8)}},
		{"chain", []Interval[int]{iv(4, 6), iv(1, 3), iv(3, 4), iv(6, 9)}, []Interval[int]{iv(1, 9)}},
		{"already sorted", []Interval[int]{iv(1, 2), iv(4, 5), iv(7, 8)}, []Interval[int]{iv(1, 2), iv(4, 5), iv(7, 8)}},
		{"reversed", []Interval[int]{iv(7, 8), iv(4, 5), iv(1, 2)}, []Interval[int]{iv(1, 2), iv(4, 5), iv(7, 8)}},
		{"negative and zero", []Interval[int]{iv(-5, -2), iv(-3, 0)}, []Interval[int]{iv(-5, 0)}},
		{"extremes of the type", []Interval[int]{iv(0, math.MaxInt), iv(math.MinInt, 0)}, []Interval[int]{iv(math.MinInt, math.MaxInt)}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MergeOverlapping(tt.in)
			if err != nil || got == nil || !slices.Equal(got, tt.want) {
				t.Errorf("MergeOverlapping(%v) = (%#v, %v), want (%v, nil)", tt.in, got, err, tt.want)
			}
		})
	}
}

func TestMergeOverlapping_floats(t *testing.T) {
	in := []Interval[float64]{{0.5, 1.5}, {1.5, math.Inf(1)}, {math.Inf(-1), -1}, {math.Copysign(0, -1), 0.25}}
	got, err := MergeOverlapping(in)
	// -0 equals +0 and sorts with it; [-Inf,-1] is apart; [-0,0.25] and [0.5,1.5] are apart; 1.5 is shared.
	want := []Interval[float64]{{math.Inf(-1), -1}, {0, 0.25}, {0.5, math.Inf(1)}}
	if err != nil || !slices.Equal(got, want) {
		t.Errorf("MergeOverlapping(floats) = (%v, %v), want (%v, nil)", got, err, want)
	}
}

func TestMergeOverlapping_errors(t *testing.T) {
	tests := []struct {
		name string
		in   []Interval[int]
	}{
		{"inverted", []Interval[int]{iv(5, 1)}},
		{"inverted after a valid one", []Interval[int]{iv(1, 2), iv(3, 2)}},
		{"inverted at the extremes", []Interval[int]{iv(math.MaxInt, math.MinInt)}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := MergeOverlapping(tt.in); !errors.Is(err, ErrInverted) || got != nil {
				t.Errorf("MergeOverlapping(%v) = (%v, %v), want (nil, ErrInverted)", tt.in, got, err)
			}
		})
	}
	// NaN sorts before every value (D-002), so [x, NaN] is inverted for any non-NaN x.
	if _, err := MergeOverlapping([]Interval[float64]{{1, math.NaN()}}); !errors.Is(err, ErrInverted) {
		t.Errorf("MergeOverlapping([1, NaN]) error = %v, want ErrInverted", err)
	}
}

// coverage is the slow obvious oracle: mark every covered point on a doubled integer grid (2x), so a
// gap between neighbours such as [1,2] and [3,4] shows as the unmarked 5, then read off the runs.
func coverage(in []Interval[int], maxValue int) []Interval[int] {
	marked := make([]bool, 2*maxValue+1)
	for _, v := range in {
		for p := 2 * v.Start; p <= 2*v.End; p++ {
			marked[p] = true
		}
	}
	out := []Interval[int]{}
	for p := 0; p < len(marked); p++ {
		if !marked[p] {
			continue
		}
		start := p
		for p+1 < len(marked) && marked[p+1] {
			p++
		}
		out = append(out, iv(start/2, p/2))
	}
	return out
}

func TestMergeOverlapping_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	const maxValue = 60
	for trial := range 400 {
		in := make([]Interval[int], rng.IntN(25))
		for i := range in {
			start := rng.IntN(maxValue)
			in[i] = iv(start, start+rng.IntN(min(6, maxValue-start+1)))
		}
		before := slices.Clone(in)
		got, err := MergeOverlapping(in)
		if want := coverage(in, maxValue); err != nil || !slices.Equal(got, want) {
			t.Fatalf("trial %d: MergeOverlapping(%v) = (%v, %v), want %v", trial, in, got, err, want)
		}
		if !slices.Equal(in, before) {
			t.Fatalf("trial %d: MergeOverlapping modified its input", trial)
		}
	}
}

func BenchmarkMergeOverlapping(b *testing.B) {
	for _, size := range benchSizes {
		rng := rand.New(rand.NewPCG(1, 2))
		in := make([]Interval[int], size.n)
		for i := range in {
			start := rng.IntN(10 * size.n)
			in[i] = iv(start, start+rng.IntN(20))
		}
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _ = MergeOverlapping(in)
			}
		})
	}
}
