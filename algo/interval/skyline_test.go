package interval

import (
	"errors"
	"math"
	"math/rand/v2"
	"slices"
	"testing"
)

func bd(left, right, height int) Building[int] { return Building[int]{left, right, height} }
func kp(x, height int) KeyPoint[int]           { return KeyPoint[int]{x, height} }

// Worked example written for this repo (Manber §5.6 states the problem; its figure is not reproduced):
// buildings (1,5,3) (2,4,6) (4,7,2) (9,10,4). By hand: height 3 from 1, 6 from 2, back to 3 at 4 (the
// tall one ends), 2 at 5, ground at 7, 4 from 9, ground at 10.
func TestSkyline_workedExample(t *testing.T) {
	got, err := Skyline([]Building[int]{bd(1, 5, 3), bd(2, 4, 6), bd(4, 7, 2), bd(9, 10, 4)})
	want := []KeyPoint[int]{kp(1, 3), kp(2, 6), kp(4, 3), kp(5, 2), kp(7, 0), kp(9, 4), kp(10, 0)}
	if err != nil || !slices.Equal(got, want) {
		t.Errorf("Skyline = (%v, %v), want (%v, nil)", got, err, want)
	}
}

func TestSkyline_cases(t *testing.T) {
	tests := []struct {
		name string
		in   []Building[int]
		want []KeyPoint[int]
	}{
		{"empty", nil, []KeyPoint[int]{}},
		{"one", []Building[int]{bd(2, 5, 7)}, []KeyPoint[int]{kp(2, 7), kp(5, 0)}},
		{"zero width and zero height add nothing", []Building[int]{bd(3, 3, 9), bd(1, 4, 0)}, []KeyPoint[int]{}},
		{"touching, same height, join", []Building[int]{bd(1, 3, 5), bd(3, 6, 5)}, []KeyPoint[int]{kp(1, 5), kp(6, 0)}},
		{"touching, step up", []Building[int]{bd(1, 3, 2), bd(3, 6, 5)}, []KeyPoint[int]{kp(1, 2), kp(3, 5), kp(6, 0)}},
		{"same left, taller wins", []Building[int]{bd(0, 4, 2), bd(0, 2, 8)}, []KeyPoint[int]{kp(0, 8), kp(2, 2), kp(4, 0)}},
		{"same right", []Building[int]{bd(0, 5, 2), bd(3, 5, 8)}, []KeyPoint[int]{kp(0, 2), kp(3, 8), kp(5, 0)}},
		{"hidden behind", []Building[int]{bd(0, 10, 9), bd(2, 4, 3)}, []KeyPoint[int]{kp(0, 9), kp(10, 0)}},
		{"duplicates", []Building[int]{bd(1, 2, 3), bd(1, 2, 3)}, []KeyPoint[int]{kp(1, 3), kp(2, 0)}},
		{"negative x", []Building[int]{bd(-4, -1, 2)}, []KeyPoint[int]{kp(-4, 2), kp(-1, 0)}},
		{"extremes", []Building[int]{bd(math.MinInt, math.MaxInt, math.MaxInt)}, []KeyPoint[int]{kp(math.MinInt, math.MaxInt), kp(math.MaxInt, 0)}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Skyline(tt.in)
			if err != nil || got == nil || !slices.Equal(got, tt.want) {
				t.Errorf("Skyline(%v) = (%#v, %v), want (%v, nil)", tt.in, got, err, tt.want)
			}
		})
	}
}

// Problem-site shape, in our own words: a row of three overlapping towers where the middle one is the
// tallest, and a lone tower further along.
func TestSkyline_problemSiteShape(t *testing.T) {
	got, _ := Skyline([]Building[int]{bd(0, 6, 4), bd(3, 8, 9), bd(5, 12, 6), bd(15, 18, 3)})
	want := []KeyPoint[int]{kp(0, 4), kp(3, 9), kp(8, 6), kp(12, 0), kp(15, 3), kp(18, 0)}
	if !slices.Equal(got, want) {
		t.Errorf("Skyline = %v, want %v", got, want)
	}
}

func TestSkyline_floats(t *testing.T) {
	got, err := Skyline([]Building[float64]{{0.5, 1.5, 2.25}, {1, 3, 1}})
	want := []KeyPoint[float64]{{0.5, 2.25}, {1.5, 1}, {3, 0}}
	if err != nil || !slices.Equal(got, want) {
		t.Errorf("floats = (%v, %v), want %v", got, err, want)
	}
}

func TestSkyline_errors(t *testing.T) {
	if got, err := Skyline([]Building[int]{bd(1, 2, 3), bd(5, 4, 1)}); !errors.Is(err, ErrInverted) || got != nil {
		t.Errorf("inverted: (%v, %v), want ErrInverted", got, err)
	}
	if _, err := Skyline([]Building[int]{bd(1, 2, -1)}); !errors.Is(err, ErrNegativeHeight) {
		t.Errorf("negative height: err = %v, want ErrNegativeHeight", err)
	}
	if _, err := Skyline([]Building[float64]{{1, 2, math.NaN()}}); !errors.Is(err, ErrNegativeHeight) {
		t.Errorf("NaN height: err = %v, want ErrNegativeHeight", err)
	}
}

// Oracle: on integer coordinates the outline only changes at integers, so take the height of each unit
// cell [x, x+1) as the tallest building over it and emit a key point wherever it changes.
func TestSkyline_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	const width = 40
	for trial := range 400 {
		in := make([]Building[int], rng.IntN(15))
		for i := range in {
			left := rng.IntN(width)
			in[i] = bd(left, left+rng.IntN(min(10, width-left+1)), rng.IntN(6))
		}
		cell := make([]int, width+1)
		for _, b := range in {
			for x := b.Left; x < b.Right; x++ {
				cell[x] = max(cell[x], b.Height)
			}
		}
		want := []KeyPoint[int]{}
		prev := 0
		for x, h := range cell {
			if h != prev {
				want = append(want, kp(x, h))
				prev = h
			}
		}
		before := slices.Clone(in)
		got, err := Skyline(in)
		if err != nil || !slices.Equal(got, want) {
			t.Fatalf("trial %d: Skyline(%v) = (%v, %v), want %v", trial, in, got, err, want)
		}
		if !slices.Equal(in, before) {
			t.Fatalf("trial %d: input modified", trial)
		}
	}
}

func BenchmarkSkyline(b *testing.B) {
	for _, size := range benchSizes {
		rng := rand.New(rand.NewPCG(1, 2))
		in := make([]Building[int], size.n)
		for i := range in {
			left := rng.IntN(10 * size.n)
			in[i] = bd(left, left+1+rng.IntN(200), 1+rng.IntN(1000))
		}
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _ = Skyline(in)
			}
		})
	}
}
