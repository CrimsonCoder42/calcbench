package text

import (
	"math/rand/v2"
	"testing"
)

// naiveLongestPalindrome is the oracle: try every substring, longest first, leftmost first.
func naiveLongestPalindrome(s string) (int, int) {
	for l := len(s); l > 0; l-- {
		for i := 0; i+l <= len(s); i++ {
			ok := true
			for k := range l / 2 {
				if s[i+k] != s[i+l-1-k] {
					ok = false
					break
				}
			}
			if ok {
				return i, l
			}
		}
	}
	return 0, 0
}

// Hand-derived: racecar is the longest palindrome inside xyracecarz.
func TestLongestPalindromicSubstring_example(t *testing.T) {
	if start, length := LongestPalindromicSubstring("xyracecarz"); start != 2 || length != 7 {
		t.Errorf("LongestPalindromicSubstring(xyracecarz) = (%d, %d), want (2, 7)", start, length)
	}
}

// Our own wording of the usual exercise cases.
func TestLongestPalindromicSubstring_cases(t *testing.T) {
	tests := []struct {
		name          string
		s             string
		start, length int
	}{
		{"empty", "", 0, 0},
		{"one", "q", 0, 1},
		{"two equal", "qq", 0, 2},
		{"two different: leftmost single", "qr", 0, 1},
		{"even palindrome", "xabbay", 1, 4},
		{"whole string", "abaaba", 0, 6},
		{"tie goes left", "abacdc", 0, 3},
		{"all equal", "zzzzz", 0, 5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if start, length := LongestPalindromicSubstring(tt.s); start != tt.start || length != tt.length {
				t.Errorf("LongestPalindromicSubstring(%q) = (%d, %d), want (%d, %d)", tt.s, start, length, tt.start, tt.length)
			}
		})
	}
}

func TestLongestPalindromicSubstring_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 500 {
		s := randomText(rng, rng.IntN(30), "abc"[:1+trial%3])
		gs, gl := LongestPalindromicSubstring([]byte(s))
		if ws, wl := naiveLongestPalindrome(s); gs != ws || gl != wl {
			t.Fatalf("trial %d: LongestPalindromicSubstring(%q) = (%d, %d), want (%d, %d)", trial, s, gs, gl, ws, wl)
		}
	}
}

func BenchmarkLongestPalindromicSubstring(b *testing.B) {
	rng := rand.New(rand.NewPCG(1, 2))
	for _, size := range benchSizes {
		s := randomText(rng, size.n, "ab")
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, benchSink = LongestPalindromicSubstring(s)
			}
		})
	}
}
