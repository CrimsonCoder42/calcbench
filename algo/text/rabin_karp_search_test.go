package text

import (
	"errors"
	"math/rand/v2"
	"slices"
	"testing"
)

// Cormen et al., 3rd ed., §32.2, Figure 32.5: radix 10, q = 13, P = 31415. The valid match is at shift 6;
// the window 67399 at shift 12 is a spurious hit, which the byte check rejects.
func TestRabinKarpSearch_clrsFigure32_5(t *testing.T) {
	got, err := RabinKarpSearch(nil, clrsDigits, "31415", 10, 13)
	if err != nil || !slices.Equal(got, []int{6}) {
		t.Errorf("RabinKarpSearch(figure 32.5) = %v, %v; want [6], nil", got, err)
	}
}

func TestRabinKarpSearch_cases(t *testing.T) {
	tests := []struct {
		name, text, pattern string
		want                []int
	}{
		{"empty text", "", "a", nil},
		{"pattern longer than text", "ab", "abc", nil},
		{"whole text", "abc", "abc", []int{0}},
		{"absent", "abcdef", "xy", nil},
		{"overlapping", "aaaa", "aa", []int{0, 1, 2}},
		{"one byte", "banana", "a", []int{1, 3, 5}},
		{"at both ends", "abxab", "ab", []int{0, 3}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RabinKarpSearch(nil, tt.text, tt.pattern, DefaultRollingBase, DefaultRollingModulus)
			if err != nil || !slices.Equal(got, tt.want) {
				t.Errorf("RabinKarpSearch(%q, %q) = %v, %v; want %v", tt.text, tt.pattern, got, err, tt.want)
			}
		})
	}
}

func TestRabinKarpSearch_errors(t *testing.T) {
	dst := []int{1}
	if got, err := RabinKarpSearch(dst, "abc", "", 10, 13); !errors.Is(err, ErrEmptyPattern) || len(got) != 1 {
		t.Errorf("empty pattern: %v, %v; want [1], ErrEmptyPattern", got, err)
	}
	if got, err := RabinKarpSearch(dst, "abc", "a", 13, 13); !errors.Is(err, ErrHashParams) || len(got) != 1 {
		t.Errorf("base = modulus: %v, %v; want [1], ErrHashParams", got, err)
	}
	if _, err := RabinKarpSearch(dst, "abc", "a", 1, 1); !errors.Is(err, ErrHashParams) {
		t.Errorf("modulus 1: %v, want ErrHashParams", err)
	}
}

func TestRabinKarpSearch_noAllocation(t *testing.T) {
	buf := make([]int, 0, 8)
	allocs := testing.AllocsPerRun(100, func() {
		buf, _ = RabinKarpSearch(buf[:0], clrsDigits, "31415", DefaultRollingBase, DefaultRollingModulus)
	})
	if allocs != 0 {
		t.Errorf("RabinKarpSearch allocates %v times, want 0", allocs)
	}
}

// The oracle runs with a tiny modulus too, so that spurious hits are frequent and verification is tested.
func TestRabinKarpSearch_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	params := []struct{ base, modulus uint64 }{{2, 3}, {10, 13}, {DefaultRollingBase, DefaultRollingModulus}}
	for trial := range 600 {
		p := params[trial%len(params)]
		text := randomText(rng, rng.IntN(60), "abc"[:1+(trial/3)%3])
		pattern := randomText(rng, 1+rng.IntN(5), "abc"[:1+(trial/3)%3])
		got, err := RabinKarpSearch(nil, []byte(text), []byte(pattern), p.base, p.modulus)
		if want := naiveSearch(text, pattern); err != nil || !slices.Equal(got, want) {
			t.Fatalf("trial %d: RabinKarpSearch(%q, %q, %d, %d) = %v, %v; want %v", trial, text, pattern, p.base, p.modulus, got, err, want)
		}
	}
}

func BenchmarkRabinKarpSearch(b *testing.B) {
	rng := rand.New(rand.NewPCG(1, 2))
	for _, size := range benchSizes {
		text := randomText(rng, size.n, "ab")
		pattern := text[size.n/2 : size.n/2+16]
		buf := make([]int, 0, size.n)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				buf, _ = RabinKarpSearch(buf[:0], text, pattern, DefaultRollingBase, DefaultRollingModulus)
			}
			benchSink = len(buf)
		})
	}
}
