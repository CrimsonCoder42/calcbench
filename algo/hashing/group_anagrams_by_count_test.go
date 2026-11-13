package hashing

import (
	"errors"
	"math/rand/v2"
	"testing"
)

const lowercase = "abcdefghijklmnopqrstuvwxyz"

// Bentley, Programming Pearls, 2nd ed., Column 2: the "deiopst" anagram class, counted over a–z.
func TestGroupAnagramsByCount_bentleySignature(t *testing.T) {
	words := []string{"deposit", "pots", "dopiest", "stop", "posited", "topside", "tops"}
	want := [][]string{{"deposit", "dopiest", "posited", "topside"}, {"pots", "stop", "tops"}}
	if got, err := GroupAnagramsByCount(words, lowercase); err != nil || !nestedEqual(got, want) {
		t.Errorf("GroupAnagramsByCount = (%v, %v), want (%v, nil)", got, err, want)
	}
}

func TestGroupAnagramsByCount_cases(t *testing.T) {
	tests := []struct {
		name     string
		words    []string
		alphabet string
		want     [][]string
	}{
		{"empty", nil, lowercase, [][]string{}},
		{"empty alphabet, empty words", []string{"", ""}, "", [][]string{{"", ""}}},
		{"counts matter", []string{"aab", "abb", "bab"}, "ab", [][]string{{"aab"}, {"abb", "bab"}}},
		{"counts above one byte", []string{repeat('a', 300), repeat('a', 300)}, "a",
			[][]string{{repeat('a', 300), repeat('a', 300)}}},
		{"unicode alphabet", []string{"éte", "été", "tée"}, "éte", [][]string{{"éte", "tée"}, {"été"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GroupAnagramsByCount(tt.words, tt.alphabet)
			if err != nil || got == nil || !nestedEqual(got, tt.want) {
				t.Errorf("GroupAnagramsByCount(%q, %q) = (%#v, %v), want %q", tt.words, tt.alphabet, got, err, tt.want)
			}
		})
	}
}

func repeat(r rune, n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = r
	}
	return string(b)
}

func TestGroupAnagramsByCount_errors(t *testing.T) {
	tests := []struct {
		name     string
		words    []string
		alphabet string
		want     error
	}{
		{"repeated ASCII rune", nil, "aba", ErrAlphabet},
		{"repeated other rune", nil, "éé", ErrAlphabet},
		{"invalid UTF-8 alphabet", nil, "a\xff", ErrAlphabet},
		{"U+FFFD in alphabet", nil, "a�", ErrAlphabet},
		{"rune outside alphabet", []string{"ab", "ac"}, "ab", ErrNotInAlphabet},
		{"non-ASCII outside alphabet", []string{"é"}, "e", ErrNotInAlphabet},
		{"invalid UTF-8 word", []string{"a\xff"}, "a", ErrNotInAlphabet},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := GroupAnagramsByCount(tt.words, tt.alphabet); !errors.Is(err, tt.want) || got != nil {
				t.Errorf("GroupAnagramsByCount(%q, %q) = (%v, %v), want (nil, %v)", tt.words, tt.alphabet, got, err, tt.want)
			}
		})
	}
}

// Oracle: the count key gives the same groups as the sorted key and as pairwise comparison.
func TestGroupAnagramsByCount_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 300 {
		words := randomWords(rng, rng.IntN(60), 5, "abcé")
		got, err := GroupAnagramsByCount(words, "abcé")
		if want := anagramOracle(words); err != nil || !nestedEqual(got, want) || !nestedEqual(got, GroupAnagrams(words)) {
			t.Fatalf("trial %d: GroupAnagramsByCount(%q) = (%q, %v), want %q", trial, words, got, err, want)
		}
	}
}

func BenchmarkGroupAnagramsByCount(b *testing.B) {
	for _, size := range benchSizes {
		words := randomWords(rand.New(rand.NewPCG(1, 2)), size.n, 8, lowercase)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _ = GroupAnagramsByCount(words, lowercase)
			}
		})
	}
}
