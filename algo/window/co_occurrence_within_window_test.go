package window

import (
	"errors"
	"math"
	"math/rand/v2"
	"slices"
	"testing"
)

// naiveCoOccurrence is the O(n²) definition: for each ordered label pair (A, B), A ≠ B, count events of A
// with some event of B at |Δt| ≤ window. |Δt| uses uint64 differences, as the function does.
func naiveCoOccurrence(events []LabeledEvent[string], window int64) []CoOccurrence[string] {
	counts := map[[2]string]int{}
	for _, a := range events {
		seen := map[string]bool{}
		for _, b := range events {
			if b.Label == a.Label || seen[b.Label] {
				continue
			}
			lo, hi := a.At, b.At
			if lo > hi {
				lo, hi = hi, lo
			}
			if uint64(hi)-uint64(lo) <= uint64(window) {
				seen[b.Label] = true
				counts[[2]string{a.Label, b.Label}]++
			}
		}
	}
	out := []CoOccurrence[string]{}
	for key, c := range counts {
		out = append(out, CoOccurrence[string]{A: key[0], B: key[1], Count: c})
	}
	slices.SortFunc(out, func(x, y CoOccurrence[string]) int {
		if x.A != y.A {
			if x.A < y.A {
				return -1
			}
			return 1
		}
		if x.B < y.B {
			return -1
		} else if x.B > y.B {
			return 1
		}
		return 0
	})
	return out
}

// Pattern example (two pointers over time-sorted events, Laaksonen §8.1), by hand, times in hours, window
// 24: a at 0 and 30, b at 10, c at 100. Both a events have b within 24 h (10 and 20 apart) → (a, b) = 2;
// the one b event has an a → (b, a) = 1; c is alone.
func TestCoOccurrenceWithinWindow_example(t *testing.T) {
	ev := []LabeledEvent[string]{{0, "a"}, {30, "a"}, {10, "b"}, {100, "c"}}
	got, err := CoOccurrenceWithinWindow(ev, 24)
	want := []CoOccurrence[string]{{"a", "b", 2}, {"b", "a", 1}}
	if err != nil || !slices.Equal(got, want) {
		t.Errorf("= (%v, %v), want %v", got, err, want)
	}
}

// parity: ledger C4 (correlation.go:172). The source's JointCount for a name pair (A, B) with A < B is the
// number of A occurrences with some B within windowHours = Count(A, B) here. Source-shaped case: headache at
// 0 h and 30 h, nausea at 10 h, fatigue at 100 h, window 24 h → JointCount(headache, nausea) = 2. The source's
// CoOccurrenceRate = 2/min(2, 1) = 2 exceeds its documented range 0–1 (a consumer defect; the rate is not
// ported). Times as hours.
func TestCoOccurrenceWithinWindow_parity(t *testing.T) {
	ev := []LabeledEvent[string]{{0, "headache"}, {30, "headache"}, {10, "nausea"}, {100, "fatigue"}}
	got, err := CoOccurrenceWithinWindow(ev, 24)
	if err != nil {
		t.Fatal(err)
	}
	var joint int
	for _, c := range got {
		if c.A == "headache" && c.B == "nausea" {
			joint = c.Count
		}
	}
	if joint != 2 {
		t.Errorf("Count(headache, nausea) = %d, want 2 (source JointCount)", joint)
	}
}

func TestCoOccurrenceWithinWindow_cases(t *testing.T) {
	tests := []struct {
		name   string
		events []LabeledEvent[string]
		window int64
		want   []CoOccurrence[string]
	}{
		{"empty", nil, 5, []CoOccurrence[string]{}},
		{"one event", []LabeledEvent[string]{{1, "a"}}, 5, []CoOccurrence[string]{}},
		{"one label", []LabeledEvent[string]{{1, "a"}, {2, "a"}}, 5, []CoOccurrence[string]{}},
		{"window 0, same instant", []LabeledEvent[string]{{7, "x"}, {7, "y"}}, 0, []CoOccurrence[string]{{"x", "y", 1}, {"y", "x", 1}}},
		{"window 0, one apart", []LabeledEvent[string]{{7, "x"}, {8, "y"}}, 0, []CoOccurrence[string]{}},
		{"edge inclusive", []LabeledEvent[string]{{0, "x"}, {5, "y"}}, 5, []CoOccurrence[string]{{"x", "y", 1}, {"y", "x", 1}}},
		{"edge exclusive", []LabeledEvent[string]{{0, "x"}, {6, "y"}}, 5, []CoOccurrence[string]{}},
		{"extreme instants", []LabeledEvent[string]{{math.MinInt64, "x"}, {math.MaxInt64, "y"}}, math.MaxInt64, []CoOccurrence[string]{}},
		{"extreme instants, full window", []LabeledEvent[string]{{math.MinInt64, "x"}, {-1, "y"}}, math.MaxInt64, []CoOccurrence[string]{{"x", "y", 1}, {"y", "x", 1}}},
	}
	for _, tc := range tests {
		got, err := CoOccurrenceWithinWindow(tc.events, tc.window)
		if err != nil || !slices.Equal(got, tc.want) {
			t.Errorf("%s: = (%v, %v), want %v", tc.name, got, err, tc.want)
		}
	}
}

// Oracle: the O(n²) definition on 300 random event sets (fixed seed, n ≤ 150), with few labels (dense
// counters) and with many labels (k² > 4n + 1024, hash map); the input is unchanged.
func TestCoOccurrenceWithinWindow_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	names := make([]string, 60)
	for i := range names {
		names[i] = string(rune('A'+i/26)) + string(rune('a'+i%26))
	}
	for iter := range 300 {
		n := rng.IntN(151)
		k := 2 + rng.IntN(5)
		if iter%2 == 1 {
			k = 60
		}
		ev := make([]LabeledEvent[string], n)
		for i := range ev {
			ev[i] = LabeledEvent[string]{At: rng.Int64N(200) - 100, Label: names[rng.IntN(k)]}
		}
		before := slices.Clone(ev)
		w := rng.Int64N(30)
		got, err := CoOccurrenceWithinWindow(ev, w)
		if err != nil {
			t.Fatal(err)
		}
		if want := naiveCoOccurrence(ev, w); !slices.Equal(got, want) {
			t.Fatalf("n %d k %d w %d: got %v, want %v", n, k, w, got, want)
		}
		if !slices.Equal(ev, before) {
			t.Fatal("input modified")
		}
	}
}

// Generic over other ordered labels.
func TestCoOccurrenceWithinWindow_intLabels(t *testing.T) {
	got, err := CoOccurrenceWithinWindow([]LabeledEvent[int]{{5, 2}, {6, 1}, {50, 3}}, 3)
	want := []CoOccurrence[int]{{1, 2, 1}, {2, 1, 1}}
	if err != nil || !slices.Equal(got, want) {
		t.Errorf("= (%v, %v), want %v", got, err, want)
	}
}

func TestCoOccurrenceWithinWindow_errors(t *testing.T) {
	if _, err := CoOccurrenceWithinWindow([]LabeledEvent[string]{{1, "a"}}, -1); !errors.Is(err, ErrNegativeWindow) {
		t.Errorf("negative window: %v", err)
	}
}

var sinkCoOccurrence []CoOccurrence[int]

func BenchmarkCoOccurrenceWithinWindow(b *testing.B) {
	for _, n := range []int{1000, 10000, 100000} {
		rng := rand.New(rand.NewPCG(1, 2))
		ev := make([]LabeledEvent[int], n)
		for i := range ev {
			// About 10 events per window on average, 20 labels.
			ev[i] = LabeledEvent[int]{At: rng.Int64N(int64(n) * 100), Label: rng.IntN(20)}
		}
		name := map[int]string{1000: "1e3", 10000: "1e4", 100000: "1e5"}[n]
		b.Run(name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				sinkCoOccurrence, _ = CoOccurrenceWithinWindow(ev, 500)
			}
		})
	}
}
