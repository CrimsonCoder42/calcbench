package hashing

import (
	"math/rand/v2"
	"slices"
	"strings"
	"testing"
)

func nestedEqual(a, b [][]string) bool {
	return slices.EqualFunc(a, b, func(x, y []string) bool { return slices.Equal(x, y) })
}

// isAnagram is the slow oracle: equal sorted runes.
func isAnagram(x, y string) bool {
	rx, ry := []rune(x), []rune(y)
	slices.Sort(rx)
	slices.Sort(ry)
	return slices.Equal(rx, ry)
}

// anagramOracle groups by pairwise comparison with each group's first word: O(n²).
func anagramOracle(words []string) [][]string {
	out := [][]string{}
	for _, w := range words {
		found := false
		for g := range out {
			if isAnagram(out[g][0], w) {
				out[g] = append(out[g], w)
				found = true
				break
			}
		}
		if !found {
			out = append(out, []string{w})
		}
	}
	return out
}

// randomWords returns n words of 0..maxLen letters drawn from alphabet. Shared with the ByCount tests.
func randomWords(rng *rand.Rand, n, maxLen int, alphabet string) []string {
	letters := []rune(alphabet)
	words := make([]string, n)
	for i := range words {
		var sb strings.Builder
		for range rng.IntN(maxLen + 1) {
			sb.WriteRune(letters[rng.IntN(len(letters))])
		}
		words[i] = sb.String()
	}
	return words
}

// Bentley, Programming Pearls, 2nd ed., Column 2: deposit, dopiest, posited and topside share the
// signature "deiopst". The second group is our own.
func TestGroupAnagrams_bentleySignature(t *testing.T) {
	words := []string{"deposit", "pots", "dopiest", "stop", "posited", "topside", "tops"}
	want := [][]string{{"deposit", "dopiest", "posited", "topside"}, {"pots", "stop", "tops"}}
	if got := GroupAnagrams(words); !nestedEqual(got, want) {
		t.Errorf("GroupAnagrams = %v, want %v", got, want)
	}
}

func TestGroupAnagrams_cases(t *testing.T) {
	tests := []struct {
		name  string
		words []string
		want  [][]string
	}{
		{"empty", nil, [][]string{}},
		{"one empty word", []string{""}, [][]string{{""}}},
		{"two empty words", []string{"", ""}, [][]string{{"", ""}}},
		{"no anagrams", []string{"ab", "abc"}, [][]string{{"ab"}, {"abc"}}},
		{"duplicates kept", []string{"tea", "eat", "tea"}, [][]string{{"tea", "eat", "tea"}}},
		{"letter counts matter", []string{"aab", "abb"}, [][]string{{"aab"}, {"abb"}}},
		{"unicode runes", []string{"éte", "été", "tée"}, [][]string{{"éte", "tée"}, {"été"}}},
		{"case sensitive", []string{"Ab", "bA", "ab"}, [][]string{{"Ab", "bA"}, {"ab"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GroupAnagrams(tt.words); got == nil || !nestedEqual(got, tt.want) {
				t.Errorf("GroupAnagrams(%q) = %#v, want %q", tt.words, got, tt.want)
			}
		})
	}
}

func TestGroupAnagrams_doesNotModifyInput(t *testing.T) {
	words := []string{"b", "a", "ab", "ba"}
	_ = GroupAnagrams(words)
	if !slices.Equal(words, []string{"b", "a", "ab", "ba"}) {
		t.Errorf("GroupAnagrams modified its input: %v", words)
	}
}

func TestGroupAnagrams_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 300 {
		words := randomWords(rng, rng.IntN(60), 4, "abcé")
		if got, want := GroupAnagrams(words), anagramOracle(words); !nestedEqual(got, want) {
			t.Fatalf("trial %d: GroupAnagrams(%q) = %q, want %q", trial, words, got, want)
		}
	}
}

func BenchmarkGroupAnagrams(b *testing.B) {
	for _, size := range benchSizes {
		words := randomWords(rand.New(rand.NewPCG(1, 2)), size.n, 8, "abcdefghijklmnopqrstuvwxyz")
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_ = GroupAnagrams(words)
			}
		})
	}
}
