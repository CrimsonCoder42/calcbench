package series

import (
	"errors"
	"fmt"
	"math/rand/v2"
	"testing"
	"time"
)

// parity: ledger LocalDayBuckets, source case 1 (readings per New York day across the DST end:
// 10-31: 1, 11-01: 7, 11-02: 1).
func TestLocalDayBuckets_parity(t *testing.T) {
	loc := newYork(t)
	utc := func(d, h, m int) time.Time { return time.Date(2026, 11, d, h, m, 0, 0, time.UTC) }
	ts := []time.Time{utc(1, 2, 0), utc(1, 6, 0), utc(1, 9, 0), utc(1, 14, 0), utc(1, 18, 0), utc(1, 23, 0),
		utc(2, 1, 0), utc(2, 4, 30), utc(2, 15, 0)}
	pts := make([]Point, len(ts))
	for i, x := range ts {
		pts[i] = Point{At: x, Value: 1}
	}
	got, err := LocalDayBuckets(pts, loc)
	if err != nil {
		t.Fatal(err)
	}
	want := []DayBucket{{2026, 10, 31, 0, 1}, {2026, 11, 1, 1, 8}, {2026, 11, 2, 8, 9}}
	if len(got) != len(want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("bucket %d: got %v, want %v", i, got[i], want[i])
		}
	}
}

func TestLocalDayBuckets_cases(t *testing.T) {
	if got, err := LocalDayBuckets(nil, time.UTC); err != nil || len(got) != 0 {
		t.Errorf("empty: %v, %v", got, err)
	}
	t0 := time.Date(2026, 5, 1, 12, 0, 0, 0, time.UTC)
	// Days with no point are absent.
	got, _ := LocalDayBuckets([]Point{{t0, 1}, {t0.Add(72 * time.Hour), 2}}, time.UTC)
	if len(got) != 2 || got[1].Day != 4 || got[1].First != 1 {
		t.Errorf("gap: %v", got)
	}
}

func TestLocalDayBuckets_errors(t *testing.T) {
	t0 := time.Date(2026, 5, 1, 12, 0, 0, 0, time.UTC)
	if _, err := LocalDayBuckets(nil, nil); !errors.Is(err, ErrNilLocation) {
		t.Errorf("nil loc: %v", err)
	}
	if _, err := LocalDayBuckets([]Point{{t0, 1}, {t0.Add(-time.Nanosecond), 1}}, time.UTC); !errors.Is(err, ErrUnsorted) {
		t.Errorf("unsorted: %v", err)
	}
}

// Oracle: every point lies in exactly one bucket whose date is its local date, buckets are contiguous and
// their dates strictly increase.
func TestLocalDayBuckets_oracle(t *testing.T) {
	loc := newYork(t)
	r := rand.New(rand.NewPCG(1, 2))
	t0 := time.Date(2026, 3, 1, 0, 0, 0, 0, loc)
	for trial := 0; trial < 200; trial++ {
		pts := make([]Point, r.IntN(80))
		at := t0
		for i := range pts {
			at = at.Add(time.Duration(r.IntN(20*60)) * time.Minute)
			pts[i] = Point{At: at}
		}
		got, err := LocalDayBuckets(pts, loc)
		if err != nil {
			t.Fatal(err)
		}
		next := 0
		for k, b := range got {
			if b.First != next || b.End <= b.First {
				t.Fatalf("trial %d: bucket %d not contiguous: %v", trial, k, b)
			}
			for i := b.First; i < b.End; i++ {
				if y, m, d := pts[i].At.In(loc).Date(); y != b.Year || m != b.Month || d != b.Day {
					t.Fatalf("trial %d: point %d in wrong bucket", trial, i)
				}
			}
			next = b.End
		}
		if next != len(pts) {
			t.Fatalf("trial %d: covered %d of %d", trial, next, len(pts))
		}
	}
}

func BenchmarkLocalDayBuckets(b *testing.B) {
	loc := newYork(b)
	for _, n := range []int{1_000, 10_000, 100_000} {
		pts := benchPoints(n)
		b.Run(fmt.Sprint(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				if _, err := LocalDayBuckets(pts, loc); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
