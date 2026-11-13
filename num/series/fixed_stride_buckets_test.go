package series

import (
	"errors"
	"fmt"
	"math"
	"math/rand/v2"
	"testing"
	"time"
)

// parity: ledger BucketFixed, source case 1 (5-minute buckets: {2 points, 60/65/70}, {1 point}).
func TestFixedStrideBuckets_parity(t *testing.T) {
	base := time.Date(2026, 10, 1, 10, 0, 0, 0, time.UTC)
	pts := []Point{{base.Add(10 * time.Second), 60}, {base.Add(4*time.Minute + 59*time.Second), 70}, {base.Add(5 * time.Minute), 80}}
	got, err := FixedStrideBuckets(pts, 5*time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	want := []StrideBucket{{base, 60, 65, 70, 2}, {base.Add(5 * time.Minute), 80, 80, 80, 1}}
	if len(got) != 2 {
		t.Fatalf("got %v", got)
	}
	for i := range want {
		if !got[i].Start.Equal(want[i].Start) || got[i].Min != want[i].Min || got[i].Mean != want[i].Mean ||
			got[i].Max != want[i].Max || got[i].Count != want[i].Count {
			t.Errorf("bucket %d: got %+v, want %+v", i, got[i], want[i])
		}
	}
}

func TestFixedStrideBuckets_cases(t *testing.T) {
	if got, err := FixedStrideBuckets(nil, time.Minute); err != nil || len(got) != 0 {
		t.Errorf("empty: %v, %v", got, err)
	}
	// A zone with a non-hour offset still floors on UTC instants.
	kolkata := time.FixedZone("UTC+5:30", 5*3600+1800)
	at := time.Date(2026, 1, 1, 10, 7, 0, 0, kolkata) // 04:37 UTC
	got, _ := FixedStrideBuckets([]Point{{at, 1}}, time.Hour)
	if !got[0].Start.Equal(time.Date(2026, 1, 1, 4, 0, 0, 0, time.UTC)) {
		t.Errorf("UTC floor: %v", got[0].Start)
	}
	// Means of values near the float64 limit stay finite.
	got, _ = FixedStrideBuckets([]Point{{at, math.MaxFloat64}, {at, math.MaxFloat64}}, time.Hour)
	if got[0].Mean != math.MaxFloat64 {
		t.Errorf("large mean: %v", got[0].Mean)
	}
	// Negative values.
	got, _ = FixedStrideBuckets([]Point{{at, -3}, {at, -1}}, time.Hour)
	if got[0].Min != -3 || got[0].Max != -1 || got[0].Mean != -2 {
		t.Errorf("negative: %+v", got[0])
	}
}

func TestFixedStrideBuckets_errors(t *testing.T) {
	t0 := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	cases := []struct {
		name   string
		pts    []Point
		stride time.Duration
		want   error
	}{
		{"zero stride", nil, 0, ErrOutOfRange},
		{"negative stride", nil, -time.Minute, ErrOutOfRange},
		{"unsorted", []Point{{t0, 1}, {t0.Add(-time.Second), 1}}, time.Minute, ErrUnsorted},
		{"NaN", []Point{{t0, math.NaN()}}, time.Minute, ErrNotFinite},
		{"-Inf", []Point{{t0, math.Inf(-1)}}, time.Minute, ErrNotFinite},
	}
	for _, c := range cases {
		if _, err := FixedStrideBuckets(c.pts, c.stride); !errors.Is(err, c.want) {
			t.Errorf("%s: err = %v, want %v", c.name, err, c.want)
		}
	}
}

// Oracle: counts and sums per stride against a map-based brute force (fixed seed).
func TestFixedStrideBuckets_oracle(t *testing.T) {
	r := rand.New(rand.NewPCG(1, 2))
	t0 := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	for trial := 0; trial < 200; trial++ {
		pts := make([]Point, 1+r.IntN(100))
		at := t0
		for i := range pts {
			at = at.Add(time.Duration(r.IntN(300)) * time.Second)
			pts[i] = Point{At: at, Value: float64(r.IntN(200))}
		}
		got, err := FixedStrideBuckets(pts, 5*time.Minute)
		if err != nil {
			t.Fatal(err)
		}
		n := map[int64]int{}
		s := map[int64]float64{}
		for _, p := range pts {
			k := p.At.Unix() / 300
			n[k]++
			s[k] += p.Value
		}
		if len(got) != len(n) {
			t.Fatalf("trial %d: %d buckets, want %d", trial, len(got), len(n))
		}
		for _, b := range got {
			k := b.Start.Unix() / 300
			// Integer values, small sums: the mean is exact to 1e-12.
			if b.Count != n[k] || math.Abs(b.Mean-s[k]/float64(n[k])) > 1e-12 {
				t.Fatalf("trial %d: bucket %v, want n=%d mean=%v", trial, b, n[k], s[k]/float64(n[k]))
			}
		}
	}
}

func BenchmarkFixedStrideBuckets(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		pts := benchPoints(n)
		b.Run(fmt.Sprint(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				if _, err := FixedStrideBuckets(pts, time.Hour); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
