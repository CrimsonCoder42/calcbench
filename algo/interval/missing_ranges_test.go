package interval

import (
	"errors"
	"math"
	"math/rand/v2"
	"slices"
	"testing"
)

// Worked example written for this repo: values 2 3 7 9 within [0, 10]. By hand the absent integers are
// 0 1, 4 5 6, 8 and 10.
func TestMissingRanges_workedExample(t *testing.T) {
	got, err := MissingRanges([]int{2, 3, 7, 9}, 0, 10)
	if want := []Interval[int]{iv(0, 1), iv(4, 6), iv(8, 8), iv(10, 10)}; err != nil || !slices.Equal(got, want) {
		t.Errorf("MissingRanges = (%v, %v), want (%v, nil)", got, err, want)
	}
}

func TestMissingRanges_cases(t *testing.T) {
	tests := []struct {
		name         string
		values       []int
		lower, upper int
		want         []Interval[int]
	}{
		{"no values", nil, 3, 7, []Interval[int]{iv(3, 7)}},
		{"single point range, present", []int{4}, 4, 4, []Interval[int]{}},
		{"single point range, absent", nil, 4, 4, []Interval[int]{iv(4, 4)}},
		{"all present", []int{1, 2, 3}, 1, 3, []Interval[int]{}},
		{"duplicates", []int{2, 2, 2, 5, 5}, 1, 6, []Interval[int]{iv(1, 1), iv(3, 4), iv(6, 6)}},
		{"values outside the range", []int{-10, 0, 3, 50}, 1, 5, []Interval[int]{iv(1, 2), iv(4, 5)}},
		{"first and last present", []int{1, 5}, 1, 5, []Interval[int]{iv(2, 4)}},
		{"negative", []int{-3, -1}, -4, 0, []Interval[int]{iv(-4, -4), iv(-2, -2), iv(0, 0)}},
		{"maximum present", []int{math.MaxInt - 1, math.MaxInt}, math.MaxInt - 3, math.MaxInt, []Interval[int]{iv(math.MaxInt-3, math.MaxInt-2)}},
		{"minimum absent", []int{math.MinInt + 1}, math.MinInt, math.MinInt + 1, []Interval[int]{iv(math.MinInt, math.MinInt)}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MissingRanges(tt.values, tt.lower, tt.upper)
			if err != nil || got == nil || !slices.Equal(got, tt.want) {
				t.Errorf("MissingRanges(%v, %d, %d) = (%#v, %v), want (%v, nil)", tt.values, tt.lower, tt.upper, got, err, tt.want)
			}
		})
	}
}

func TestMissingRanges_smallTypes(t *testing.T) {
	all := make([]uint8, 256)
	for i := range all {
		all[i] = uint8(i)
	}
	if got, err := MissingRanges(all, 0, 255); err != nil || len(got) != 0 {
		t.Errorf("all of uint8 = (%v, %v), want empty", got, err)
	}
	if got, _ := MissingRanges([]int8{}, math.MinInt8, math.MaxInt8); !slices.Equal(got, []Interval[int8]{{math.MinInt8, math.MaxInt8}}) {
		t.Errorf("empty int8 = %v, want the whole type", got)
	}
	if got, _ := MissingRanges([]int8{127}, 120, 127); !slices.Equal(got, []Interval[int8]{{120, 126}}) {
		t.Errorf("int8 top present = %v, want [120,126]", got)
	}
}

// Problem-site shape, in our own words: a sorted list of ticket numbers sold out of 1–20; which are left.
func TestMissingRanges_problemSiteShape(t *testing.T) {
	got, _ := MissingRanges([]int{1, 4, 5, 6, 12, 20}, 1, 20)
	if want := []Interval[int]{iv(2, 3), iv(7, 11), iv(13, 19)}; !slices.Equal(got, want) {
		t.Errorf("tickets = %v, want %v", got, want)
	}
}

func TestMissingRanges_errors(t *testing.T) {
	if _, err := MissingRanges([]int{1}, 5, 4); !errors.Is(err, ErrInverted) {
		t.Errorf("upper < lower: err = %v, want ErrInverted", err)
	}
	for _, values := range [][]int{{3, 1}, {1, 2, 9, 4}, {0, 12, 2}, {0, 9, 3}} { // the last two: unsorted after a value above upper, and after upper itself
		if got, err := MissingRanges(values, 0, 9); !errors.Is(err, ErrNotSorted) || got != nil {
			t.Errorf("MissingRanges(%v) = (%v, %v), want ErrNotSorted", values, got, err)
		}
	}
}

func TestMissingRanges_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 400 {
		values := make([]int, rng.IntN(30))
		for i := range values {
			values[i] = rng.IntN(50) - 5
		}
		slices.Sort(values)
		lower := rng.IntN(20)
		upper := lower + rng.IntN(30)
		present := map[int]bool{}
		for _, v := range values {
			present[v] = true
		}
		marked := make([]bool, 2*(upper+1)+1) // doubled grid; a gap of one integer shows as an unmarked odd point
		for x := lower; x <= upper; x++ {
			if !present[x] {
				marked[2*x] = true
				if x+1 <= upper && !present[x+1] {
					marked[2*x+1] = true
				}
			}
		}
		want := gridRuns(marked)
		got, err := MissingRanges(values, lower, upper)
		if err != nil || !slices.Equal(got, want) {
			t.Fatalf("trial %d: MissingRanges(%v, %d, %d) = (%v, %v), want %v", trial, values, lower, upper, got, err, want)
		}
	}
}

func BenchmarkMissingRanges(b *testing.B) {
	for _, size := range benchSizes {
		values := make([]int, size.n)
		for i := range values {
			values[i] = 2*i + i%3 // gaps of varying width
		}
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _ = MissingRanges(values, 0, 3*size.n)
			}
		})
	}
}
