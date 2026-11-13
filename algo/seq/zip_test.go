package seq

import (
	"math/rand/v2"
	"slices"
	"testing"
)

type pair struct {
	a int
	b string
}

func collectPairs(s func(yield func(int, string) bool)) []pair {
	var out []pair
	for a, b := range s {
		out = append(out, pair{a, b})
	}
	return out
}

// Hand-derived: 1 2 3 zipped with "x" "y" stops at the shorter: (1, x) (2, y).
func TestZip_example(t *testing.T) {
	got := collectPairs(Zip(slices.Values([]int{1, 2, 3}), slices.Values([]string{"x", "y"})))
	want := []pair{{1, "x"}, {2, "y"}}
	if !slices.Equal(got, want) {
		t.Errorf("Zip = %v, want %v", got, want)
	}
}

func TestZip_cases(t *testing.T) {
	tests := []struct {
		name  string
		first []int
		secnd []string
		want  []pair
	}{
		{"both empty", nil, nil, nil},
		{"first empty", nil, []string{"a"}, nil},
		{"second empty", []int{1}, nil, nil},
		{"one each", []int{1}, []string{"a"}, []pair{{1, "a"}}},
		{"equal lengths", []int{1, 2}, []string{"a", "b"}, []pair{{1, "a"}, {2, "b"}}},
		{"second longer", []int{1}, []string{"a", "b"}, []pair{{1, "a"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := collectPairs(Zip(slices.Values(tt.first), slices.Values(tt.secnd)))
			if !slices.Equal(got, tt.want) {
				t.Errorf("Zip(%v, %q) = %v, want %v", tt.first, tt.secnd, got, tt.want)
			}
		})
	}
}

func TestZip_nil(t *testing.T) {
	for range Zip[int, int](nil, slices.Values([]int{1})) {
		t.Error("Zip(nil, …) yielded a pair")
	}
	for range Zip[int, int](slices.Values([]int{1}), nil) {
		t.Error("Zip(…, nil) yielded a pair")
	}
}

func TestZip_firstShorterStopsSecond(t *testing.T) {
	var p1, p2 probe
	n := 0
	for range Zip(p1.ints(2), p2.ints(-1)) {
		n++
	}
	if n != 2 || !p1.finished || p2.produced != 2 || !p2.stopped {
		t.Errorf("pairs %d, first finished %v, second produced %d stopped %v; want 2, true, 2, true", n, p1.finished, p2.produced, p2.stopped)
	}
}

func TestZip_secondShorterStopsFirst(t *testing.T) {
	var p1, p2 probe
	n := 0
	for range Zip(p1.ints(-1), p2.ints(2)) {
		n++
	}
	// first produced a 3rd element before second reported its end; first is then stopped.
	if n != 2 || !p2.finished || p1.produced != 3 || !p1.stopped {
		t.Errorf("pairs %d, second finished %v, first produced %d stopped %v; want 2, true, 3, true", n, p2.finished, p1.produced, p1.stopped)
	}
}

func TestZip_stopsBothOnBreak(t *testing.T) {
	var p1, p2 probe
	for a := range Zip(p1.ints(-1), p2.ints(-1)) {
		if a == 1 {
			break
		}
	}
	if !p1.stopped || !p2.stopped || p1.produced != 2 || p2.produced != 2 {
		t.Errorf("first produced %d stopped %v, second produced %d stopped %v; want 2 true 2 true", p1.produced, p1.stopped, p2.produced, p2.stopped)
	}
}

func TestZip_stopsSecondOnPanic(t *testing.T) {
	var p1, p2 probe
	func() {
		defer func() {
			if recover() == nil {
				t.Error("the loop body's panic was swallowed")
			}
		}()
		for range Zip(p1.ints(-1), p2.ints(-1)) {
			panic("loop body")
		}
	}()
	if !p2.stopped {
		t.Error("Zip did not call the second iterator's stop after a panic in the loop body")
	}
}

func TestZip_secondNotStartedWhenFirstEmpty(t *testing.T) {
	var p2 probe
	for range Zip(slices.Values([]int(nil)), p2.ints(-1)) {
	}
	if p2.produced != 0 {
		t.Errorf("second produced %d with an empty first, want 0", p2.produced)
	}
}

func TestZip_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for range 300 {
		a, b := randomInts(rng, rng.IntN(20)), randomInts(rng, rng.IntN(20))
		var got [][2]int
		for x, y := range Zip(slices.Values(a), slices.Values(b)) {
			got = append(got, [2]int{x, y})
		}
		n := min(len(a), len(b))
		if len(got) != n {
			t.Fatalf("Zip(%v, %v) gave %d pairs, want %d", a, b, len(got), n)
		}
		for i := range n {
			if got[i] != [2]int{a[i], b[i]} {
				t.Fatalf("Zip(%v, %v)[%d] = %v", a, b, i, got[i])
			}
		}
	}
}

// iter.Pull allocates once per range loop; no pair allocates.
func TestZip_noAllocationPerPair(t *testing.T) {
	run := func(n int) float64 {
		s := Zip(count(n), count(n))
		return testing.AllocsPerRun(20, func() {
			for range s {
			}
		})
	}
	if small, large := run(10), run(10_000); small != large {
		t.Errorf("Zip allocates %v times over 10 pairs and %v over 10000; want equal", small, large)
	}
}

func BenchmarkZip(b *testing.B) {
	for _, size := range benchSizes {
		b.Run(size.name, func(b *testing.B) {
			s := Zip(count(size.n), count(size.n))
			b.ReportAllocs()
			for b.Loop() {
				sum := 0
				for x, y := range s {
					sum += x + y
				}
				_ = sum
			}
		})
	}
}
