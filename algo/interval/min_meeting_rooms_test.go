package interval

import (
	"errors"
	"math"
	"math/rand/v2"
	"testing"
)

// Worked example written for this repo: [1,4] [2,5] [6,8] [7,9] [3,3]. By hand, the point 3 lies in
// [1,4], [2,5] and [3,3], and no point lies in four, so the depth (Kleinberg–Tardos §4.1) is 3.
func TestMinMeetingRooms_workedExample(t *testing.T) {
	in := []Interval[int]{iv(1, 4), iv(2, 5), iv(6, 8), iv(7, 9), iv(3, 3)}
	if got, err := MinMeetingRooms(in); err != nil || got != 3 {
		t.Errorf("MinMeetingRooms = (%d, %v), want (3, nil)", got, err)
	}
}

func TestMinMeetingRooms_cases(t *testing.T) {
	tests := []struct {
		name string
		in   []Interval[int]
		want int
	}{
		{"empty", nil, 0},
		{"one", []Interval[int]{iv(1, 2)}, 1},
		{"two apart", []Interval[int]{iv(1, 2), iv(3, 4)}, 1},
		{"shared endpoint needs two (closed)", []Interval[int]{iv(1, 3), iv(3, 5)}, 2},
		{"nested", []Interval[int]{iv(1, 10), iv(2, 3), iv(4, 5)}, 2},
		{"duplicates", []Interval[int]{iv(2, 4), iv(2, 4), iv(2, 4)}, 3},
		{"points", []Interval[int]{iv(5, 5), iv(5, 5)}, 2},
		{"sorted chain", []Interval[int]{iv(1, 2), iv(2, 3), iv(3, 4)}, 2},
		{"reversed chain", []Interval[int]{iv(3, 4), iv(2, 3), iv(1, 2)}, 2},
		{"negative", []Interval[int]{iv(-5, -1), iv(-2, 0)}, 2},
		{"extremes", []Interval[int]{iv(math.MinInt, math.MaxInt), iv(0, 0)}, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := MinMeetingRooms(tt.in); err != nil || got != tt.want {
				t.Errorf("MinMeetingRooms(%v) = (%d, %v), want (%d, nil)", tt.in, got, err, tt.want)
			}
		})
	}
}

// Problem-site shape, in our own words: one long meeting with two short ones inside it at different
// times needs two rooms; meetings that never meet need one. Half-open bookings are passed as End-1.
func TestMinMeetingRooms_problemSiteShape(t *testing.T) {
	if got, _ := MinMeetingRooms([]Interval[int]{iv(0, 40), iv(10, 15), iv(25, 30)}); got != 2 {
		t.Errorf("long with two inside = %d, want 2", got)
	}
	// Bookings 9:00–10:00 and 10:00–11:00 in minutes, half-open, passed as [540, 599] and [600, 659].
	if got, _ := MinMeetingRooms([]Interval[int]{iv(540, 599), iv(600, 659)}); got != 1 {
		t.Errorf("back-to-back half-open = %d, want 1", got)
	}
}

func TestMinMeetingRooms_errors(t *testing.T) {
	if got, err := MinMeetingRooms([]Interval[int]{iv(1, 2), iv(4, 3)}); !errors.Is(err, ErrInverted) || got != 0 {
		t.Errorf("inverted: (%d, %v), want (0, ErrInverted)", got, err)
	}
}

// Oracle: the maximum, over the doubled grid, of the number of intervals containing each point.
func TestMinMeetingRooms_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	const maxValue = 40
	for trial := range 400 {
		in := make([]Interval[int], rng.IntN(25))
		for i := range in {
			start := rng.IntN(maxValue)
			in[i] = iv(start, start+rng.IntN(min(8, maxValue-start+1)))
		}
		depth := make([]int, 2*maxValue+1)
		want := 0
		for _, v := range in {
			for p := 2 * v.Start; p <= 2*v.End; p++ {
				depth[p]++
				want = max(want, depth[p])
			}
		}
		if got, err := MinMeetingRooms(in); err != nil || got != want {
			t.Fatalf("trial %d: MinMeetingRooms(%v) = (%d, %v), want %d", trial, in, got, err, want)
		}
	}
}

func BenchmarkMinMeetingRooms(b *testing.B) {
	for _, size := range benchSizes {
		rng := rand.New(rand.NewPCG(1, 2))
		in := make([]Interval[int], size.n)
		for i := range in {
			start := rng.IntN(10 * size.n)
			in[i] = iv(start, start+rng.IntN(100))
		}
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _ = MinMeetingRooms(in)
			}
		})
	}
}
