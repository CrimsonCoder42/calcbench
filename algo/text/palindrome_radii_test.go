package text

import (
	"math/rand/v2"
	"slices"
	"testing"
)

// naiveRadii is the oracle: expand around each of the 2n−1 centres.
func naiveRadii(s string) []int {
	n := len(s)
	if n == 0 {
		return []int{}
	}
	out := make([]int, 2*n-1)
	for c := range out {
		lo, hi := c/2, (c+1)/2 // equal for a byte centre, adjacent for a gap
		for lo >= 0 && hi < n && s[lo] == s[hi] {
			lo--
			hi++
		}
		out[c] = hi - lo - 1
	}
	return out
}

// Hand-derived: abaaba has the even palindrome abaaba at centre 5 and odd ones aba at centres 2 and 8.
func TestPalindromeRadii_example(t *testing.T) {
	got := PalindromeRadii(nil, "abaaba")
	if want := []int{1, 0, 3, 0, 1, 6, 1, 0, 3, 0, 1}; !slices.Equal(got, want) {
		t.Errorf("PalindromeRadii(abaaba) = %v, want %v", got, want)
	}
}

func TestPalindromeRadii_cases(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want []int
	}{
		{"empty", "", []int{}},
		{"one", "z", []int{1}},
		{"two equal", "zz", []int{1, 2, 1}},
		{"two different", "zy", []int{1, 0, 1}},
		{"all equal", "aaaa", []int{1, 2, 3, 4, 3, 2, 1}},
		{"no palindrome longer than 1", "abc", []int{1, 0, 1, 0, 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := PalindromeRadii(nil, tt.s)
			if got == nil {
				got = []int{}
			}
			if !slices.Equal(got, tt.want) {
				t.Errorf("PalindromeRadii(%q) = %v, want %v", tt.s, got, tt.want)
			}
		})
	}
}

func TestPalindromeRadii_appendKeepsPrefix(t *testing.T) {
	if got, want := PalindromeRadii([]int{-1, -2}, []byte("aa")), []int{-1, -2, 1, 2, 1}; !slices.Equal(got, want) {
		t.Errorf("PalindromeRadii([-1 -2], aa) = %v, want %v", got, want)
	}
}

func TestPalindromeRadii_noAllocation(t *testing.T) {
	buf := make([]int, 0, 2*6+1) // the working size is 2n+1
	if allocs := testing.AllocsPerRun(100, func() { buf = PalindromeRadii(buf[:0], "abaaba") }); allocs != 0 {
		t.Errorf("PalindromeRadii with capacity 2n+1 allocates %v times, want 0", allocs)
	}
}

func TestPalindromeRadii_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 500 {
		s := randomText(rng, rng.IntN(40), "abc"[:1+trial%3])
		if got, want := PalindromeRadii([]int{}, s), naiveRadii(s); !slices.Equal(got, want) {
			t.Fatalf("trial %d: PalindromeRadii(%q) = %v, want %v", trial, s, got, want)
		}
	}
}

func BenchmarkPalindromeRadii(b *testing.B) {
	rng := rand.New(rand.NewPCG(1, 2))
	for _, size := range benchSizes {
		s := randomText(rng, size.n, "ab")
		buf := make([]int, 0, 2*size.n+1)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				buf = PalindromeRadii(buf[:0], s)
			}
			benchSink = len(buf)
		})
	}
}
