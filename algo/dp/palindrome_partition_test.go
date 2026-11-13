package dp

import (
	"math/rand/v2"
	"slices"
	"testing"
)

func isPalindrome[T comparable](s []T) bool {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		if s[i] != s[j] {
			return false
		}
	}
	return true
}

// minCutsSlow is the prefix recurrence of Fici et al. §2 with a direct palindrome test: O(n³).
func minCutsSlow(s []byte) int {
	if len(s) == 0 {
		return 0
	}
	f := make([]int, len(s)+1) // f[j] = fewest palindromic factors of s[:j]
	for j := 1; j <= len(s); j++ {
		f[j] = j
		for i := range j {
			if isPalindrome(s[i:j]) {
				f[j] = min(f[j], f[i]+1)
			}
		}
	}
	return f[len(s)] - 1
}

// Cases worked by hand: racecar|x|y|z needs 3 cuts; noon|abba|d needs 2.
func TestPalindromePartitionMinCuts_cases(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want int
	}{
		{"empty", "", 0},
		{"one", "a", 0},
		{"two equal", "aa", 0},
		{"two different", "ab", 1},
		{"palindrome", "racecar", 0},
		{"palindrome then singles", "racecarxyz", 3},
		{"two palindromes then one", "noonabbad", 2},
		{"no repeats", "abcde", 4},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PalindromePartitionMinCuts([]byte(tt.s)); got != tt.want {
				t.Errorf("PalindromePartitionMinCuts(%q) = %d, want %d", tt.s, got, tt.want)
			}
		})
	}
}

func TestPalindromePartitionMinCuts_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 400 {
		s := randomBytes(rng, rng.IntN(30), "ab")
		orig := slices.Clone(s)
		if got, want := PalindromePartitionMinCuts(s), minCutsSlow(s); got != want {
			t.Fatalf("trial %d: %q → %d, want %d", trial, s, got, want)
		}
		if !slices.Equal(s, orig) {
			t.Fatalf("input modified")
		}
	}
}

// All-equal input is the worst case (every centre expands to the boundary); random input is typical.
func BenchmarkPalindromePartitionMinCuts(b *testing.B) {
	for _, size := range quadraticSizes {
		rng := rand.New(rand.NewPCG(1, 2))
		for _, in := range []struct {
			name string
			s    []byte
		}{{"random-", randomBytes(rng, size.n, "ab")}, {"same-", randomBytes(rng, size.n, "a")}} {
			b.Run(in.name+size.name, func(b *testing.B) {
				b.ReportAllocs()
				for b.Loop() {
					_ = PalindromePartitionMinCuts(in.s)
				}
			})
		}
	}
}
