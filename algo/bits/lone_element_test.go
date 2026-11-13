package bits

import (
	"math"
	"math/rand/v2"
	"slices"
	"testing"
)

// shuffledMultiset returns lone once plus each of others repeated times times, shuffled.
func shuffledMultiset(rng *rand.Rand, lone int64, others []int64, times int) []int64 {
	values := []int64{lone}
	for _, o := range others {
		for range times {
			values = append(values, o)
		}
	}
	rng.Shuffle(len(values), func(i, j int) { values[i], values[j] = values[j], values[i] })
	return values
}

// Knuth Vol. 4A §7.1.3: x ⊕ x = 0, so pairs cancel wherever they stand.
func TestLoneElementAmongPairs_cases(t *testing.T) {
	tests := []struct {
		name   string
		values []int64
		want   int64
	}{
		{"empty", nil, 0},
		{"one", []int64{42}, 42},
		{"pairs around the lone value", []int64{3, 8, 5, 3, 5}, 8},
		{"lone zero", []int64{7, 0, 7}, 0},
		{"negative", []int64{-9, 4, 4}, -9},
		{"extremes", []int64{math.MinInt64, math.MaxInt64, math.MinInt64}, math.MaxInt64},
		{"value four times", []int64{2, 2, 2, 2, 6}, 6},
	}
	for _, tt := range tests {
		if got := LoneElementAmongPairs(tt.values); got != tt.want {
			t.Errorf("%s: LoneElementAmongPairs(%v) = %d, want %d", tt.name, tt.values, got, tt.want)
		}
	}
	if got := LoneElementAmongPairs([]uint8{200, 13, 200}); got != 13 {
		t.Errorf("uint8: got %d, want 13", got)
	}
}

func TestLoneElementAmongPairs_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for range 300 {
		lone := rng.Int64() - rng.Int64()
		others := make([]int64, rng.IntN(50))
		for i := range others {
			others[i] = rng.Int64()
		}
		values := shuffledMultiset(rng, lone, others, 2)
		before := slices.Clone(values)
		if got := LoneElementAmongPairs(values); got != lone {
			t.Fatalf("LoneElementAmongPairs = %d, want %d", got, lone)
		}
		if !slices.Equal(values, before) {
			t.Fatal("input modified")
		}
	}
}

// Every bit position holds 3m + b set bits; b is the lone value's bit.
func TestLoneElementAmongTriples_cases(t *testing.T) {
	tests := []struct {
		name   string
		values []int64
		want   int64
	}{
		{"empty", nil, 0},
		{"one", []int64{42}, 42},
		{"triple then lone", []int64{5, 5, 5, 11}, 11},
		{"interleaved", []int64{4, 9, 4, 1, 9, 4, 9}, 1},
		{"lone zero", []int64{3, 0, 3, 3}, 0},
		{"negative lone", []int64{-6, 2, 2, 2}, -6},
		{"negative triples", []int64{-1, -1, 7, -1}, 7},
		{"extremes", []int64{math.MinInt64, math.MaxInt64, math.MinInt64, math.MinInt64}, math.MaxInt64},
	}
	for _, tt := range tests {
		if got := LoneElementAmongTriples(tt.values); got != tt.want {
			t.Errorf("%s: LoneElementAmongTriples(%v) = %d, want %d", tt.name, tt.values, got, tt.want)
		}
	}
	if got := LoneElementAmongTriples([]uint16{65535, 3, 65535, 65535}); got != 3 {
		t.Errorf("uint16: got %d, want 3", got)
	}
}

func TestLoneElementAmongTriples_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for range 300 {
		lone := rng.Int64() - rng.Int64()
		others := make([]int64, rng.IntN(50))
		for i := range others {
			others[i] = rng.Int64() - rng.Int64()
		}
		values := shuffledMultiset(rng, lone, others, 3)
		if got := LoneElementAmongTriples(values); got != lone {
			t.Fatalf("LoneElementAmongTriples = %d, want %d", got, lone)
		}
	}
}

func TestLoneElement_noAllocation(t *testing.T) {
	values := []int64{1, 1, 1, 2}
	if a := testing.AllocsPerRun(20, func() { _ = LoneElementAmongPairs(values) + LoneElementAmongTriples(values) }); a != 0 {
		t.Errorf("allocates %v times, want 0", a)
	}
}

func benchTriples(n int) []int64 {
	rng := rand.New(rand.NewPCG(1, 2))
	others := make([]int64, n/3)
	for i := range others {
		others[i] = rng.Int64()
	}
	return shuffledMultiset(rng, 12345, others, 3)
}

func BenchmarkLoneElementAmongPairs(b *testing.B) {
	for _, size := range benchSizes {
		values := benchTriples(size.n)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			var s int64
			for b.Loop() {
				s += LoneElementAmongPairs(values)
			}
			benchSink = uint64(s)
		})
	}
}

func BenchmarkLoneElementAmongTriples(b *testing.B) {
	for _, size := range benchSizes {
		values := benchTriples(size.n)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			var s int64
			for b.Loop() {
				s += LoneElementAmongTriples(values)
			}
			benchSink = uint64(s)
		})
	}
}
