package dp

import (
	"errors"
	"math/rand/v2"
	"testing"
)

// occurrencesSlow counts index sets recursively (short inputs only).
func occurrencesSlow(text, pattern []byte) uint64 {
	if len(pattern) == 0 {
		return 1
	}
	if len(text) < len(pattern) {
		return 0
	}
	n := occurrencesSlow(text[1:], pattern)
	if text[0] == pattern[0] {
		n += occurrencesSlow(text[1:], pattern[1:])
	}
	return n
}

// Lothaire, Combinatorics on Words, ch. 6, §6.3 (binomial coefficients of words); counted by hand:
// "see" occurs 6 times in "subsequence" (two s's, each before C(3, 2) = 3 pairs of e's).
func TestSubsequenceOccurrences_wordBinomial(t *testing.T) {
	if got, err := SubsequenceOccurrences([]byte("subsequence"), []byte("see")); err != nil || got != 6 {
		t.Errorf("SubsequenceOccurrences(subsequence, see) = (%d, %v), want (6, nil)", got, err)
	}
}

func TestSubsequenceOccurrences_cases(t *testing.T) {
	tests := []struct {
		name          string
		text, pattern string
		want          uint64
	}{
		{"both empty", "", "", 1},
		{"empty pattern", "abc", "", 1},
		{"empty text", "", "a", 0},
		{"pattern longer", "ab", "abc", 0},
		{"one match", "a", "a", 1},
		{"no match", "a", "b", 0},
		{"equal", "abc", "abc", 1},
		{"repeats", "abab", "ab", 3},
		{"all same", "aaaa", "aa", 6}, // C(4, 2)
		{"reversed order does not count", "ba", "ab", 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := SubsequenceOccurrences([]byte(tt.text), []byte(tt.pattern)); err != nil || got != tt.want {
				t.Errorf("SubsequenceOccurrences(%q, %q) = (%d, %v), want %d", tt.text, tt.pattern, got, err, tt.want)
			}
		})
	}
}

func TestSubsequenceOccurrences_errors(t *testing.T) {
	// "a"×200 choose "a"×100 = C(200, 100) ≈ 9·10⁵⁸.
	text, pattern := make([]byte, 200), make([]byte, 100)
	for i := range text {
		text[i] = 'a'
	}
	for i := range pattern {
		pattern[i] = 'a'
	}
	if _, err := SubsequenceOccurrences(text, pattern); !errors.Is(err, ErrOverflow) {
		t.Errorf("C(200, 100): err = %v, want ErrOverflow", err)
	}
}

func TestSubsequenceOccurrences_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 400 {
		text := randomBytes(rng, rng.IntN(16), "ab")
		pattern := randomBytes(rng, rng.IntN(5), "ab")
		if got, err := SubsequenceOccurrences(text, pattern); err != nil || got != occurrencesSlow(text, pattern) {
			t.Fatalf("trial %d: (%q, %q) = (%d, %v), want %d", trial, text, pattern, got, err, occurrencesSlow(text, pattern))
		}
	}
}

// Pattern of length 100 over a 4-letter text: n·m = 10⁵, 10⁶, 10⁷ cell visits. Prefix counts would
// overflow on long texts, so the pattern starts with a letter the text lacks: every count past the
// first stays 0 while every cell is still visited.
func BenchmarkSubsequenceOccurrences(b *testing.B) {
	rng := rand.New(rand.NewPCG(1, 2))
	pattern := append([]byte{'Z'}, randomBytes(rng, 99, "ACGT")...)
	for _, size := range linearSizes {
		text := randomBytes(rng, size.n, "ACGT")
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _ = SubsequenceOccurrences(text, pattern)
			}
		})
	}
}
