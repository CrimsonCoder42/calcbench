package interval

import (
	"errors"
	"math"
	"math/rand/v2"
	"slices"
	"testing"
)

// gridRuns reads the maximal runs of marked points on a doubled grid (point p stands for p/2) as intervals.
func gridRuns(marked []bool) []Interval[int] {
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

// mark sets the doubled-grid points of every interval of list.
func mark(list []Interval[int], maxValue int) []bool {
	marked := make([]bool, 2*maxValue+1)
	for _, v := range list {
		for p := 2 * v.Start; p <= 2*v.End; p++ {
			marked[p] = true
		}
	}
	return marked
}

// Worked example written for this repo: [1,3] [5,9] against [2,6] [8,10]. By hand: [1,3]∩[2,6] = [2,3],
// [5,9]∩[2,6] = [5,6], [5,9]∩[8,10] = [8,9].
func TestIntersectLists_workedExample(t *testing.T) {
	got, err := IntersectLists([]Interval[int]{iv(1, 3), iv(5, 9)}, []Interval[int]{iv(2, 6), iv(8, 10)})
	if want := []Interval[int]{iv(2, 3), iv(5, 6), iv(8, 9)}; err != nil || !slices.Equal(got, want) {
		t.Errorf("IntersectLists = (%v, %v), want (%v, nil)", got, err, want)
	}
}

func TestIntersectLists_cases(t *testing.T) {
	tests := []struct {
		name string
		a, b []Interval[int]
		want []Interval[int]
	}{
		{"both empty", nil, nil, []Interval[int]{}},
		{"one empty", []Interval[int]{iv(1, 2)}, nil, []Interval[int]{}},
		{"disjoint", []Interval[int]{iv(1, 2)}, []Interval[int]{iv(4, 5)}, []Interval[int]{}},
		{"shared endpoint is a point", []Interval[int]{iv(1, 3)}, []Interval[int]{iv(3, 5)}, []Interval[int]{iv(3, 3)}},
		{"identical", []Interval[int]{iv(1, 4)}, []Interval[int]{iv(1, 4)}, []Interval[int]{iv(1, 4)}},
		{"one spans many", []Interval[int]{iv(0, 20)}, []Interval[int]{iv(1, 2), iv(5, 6), iv(9, 9)}, []Interval[int]{iv(1, 2), iv(5, 6), iv(9, 9)}},
		{"equal ends", []Interval[int]{iv(1, 5), iv(7, 8)}, []Interval[int]{iv(3, 5), iv(6, 8)}, []Interval[int]{iv(3, 5), iv(7, 8)}},
		{"negative", []Interval[int]{iv(-5, -1)}, []Interval[int]{iv(-3, 2)}, []Interval[int]{iv(-3, -1)}},
		{"extremes", []Interval[int]{iv(math.MinInt, math.MaxInt)}, []Interval[int]{iv(math.MinInt, 0)}, []Interval[int]{iv(math.MinInt, 0)}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := IntersectLists(tt.a, tt.b)
			if err != nil || got == nil || !slices.Equal(got, tt.want) {
				t.Errorf("IntersectLists(%v, %v) = (%#v, %v), want (%v, nil)", tt.a, tt.b, got, err, tt.want)
			}
		})
	}
}

// Problem-site shape, in our own words: two people's free-time lists; the common free time.
func TestIntersectLists_problemSiteShape(t *testing.T) {
	alice := []Interval[int]{iv(9, 11), iv(13, 17)}
	bob := []Interval[int]{iv(10, 14), iv(16, 18)}
	got, _ := IntersectLists(alice, bob)
	if want := []Interval[int]{iv(10, 11), iv(13, 14), iv(16, 17)}; !slices.Equal(got, want) {
		t.Errorf("common free time = %v, want %v", got, want)
	}
}

func TestIntersectLists_errors(t *testing.T) {
	good := []Interval[int]{iv(1, 2)}
	for _, bad := range [][]Interval[int]{{iv(5, 6), iv(1, 2)}, {iv(1, 3), iv(3, 4)}} {
		if _, err := IntersectLists(bad, good); !errors.Is(err, ErrNotSortedDisjoint) {
			t.Errorf("a = %v: err = %v, want ErrNotSortedDisjoint", bad, err)
		}
		if _, err := IntersectLists(good, bad); !errors.Is(err, ErrNotSortedDisjoint) {
			t.Errorf("b = %v: err = %v, want ErrNotSortedDisjoint", bad, err)
		}
	}
	if _, err := IntersectLists(good, []Interval[int]{iv(4, 3)}); !errors.Is(err, ErrInverted) {
		t.Errorf("inverted b: err = %v, want ErrInverted", err)
	}
}

func TestIntersectLists_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	const maxValue = 60
	for trial := range 400 {
		a := randomSortedDisjoint(rng, rng.IntN(12), maxValue)
		b := randomSortedDisjoint(rng, rng.IntN(12), maxValue+5)
		ma, mb := mark(a, maxValue+12), mark(b, maxValue+12)
		for p := range ma {
			ma[p] = ma[p] && mb[p]
		}
		want := gridRuns(ma)
		got, err := IntersectLists(a, b)
		if err != nil || !slices.Equal(got, want) {
			t.Fatalf("trial %d: IntersectLists(%v, %v) = (%v, %v), want %v", trial, a, b, got, err, want)
		}
	}
}

func BenchmarkIntersectLists(b *testing.B) {
	for _, size := range benchSizes {
		x := make([]Interval[int], size.n)
		y := make([]Interval[int], size.n)
		for i := range x {
			x[i] = iv(10*i, 10*i+6)
			y[i] = iv(10*i+4, 10*i+8)
		}
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _ = IntersectLists(x, y)
			}
		})
	}
}
