package text

import (
	"math/rand/v2"
	"slices"
	"strings"
	"testing"
)

// naiveReverseWords is the oracle: split on ASCII white space, reverse, join with one space.
func naiveReverseWords(s string) string {
	words := strings.FieldsFunc(s, func(r rune) bool { return r < 0x80 && isASCIISpace(byte(r)) })
	slices.Reverse(words)
	return strings.Join(words, " ")
}

// Bentley, Programming Pearls, 2nd ed., Column 2.3: reversing the whole and then each piece swaps the
// pieces; on words, hand-derived: "the quick brown fox" → "fox brown quick the".
func TestReverseWords_example(t *testing.T) {
	if got := ReverseWords("the quick brown fox"); got != "fox brown quick the" {
		t.Errorf("ReverseWords = %q, want %q", got, "fox brown quick the")
	}
}

// Our own wording of the usual exercise cases.
func TestReverseWords_cases(t *testing.T) {
	tests := []struct{ name, s, want string }{
		{"empty", "", ""},
		{"only spaces", "   \t\n ", ""},
		{"one word", "word", "word"},
		{"one word padded", "  word  ", "word"},
		{"two words", "left right", "right left"},
		{"runs of white space collapse", " a \t\n b\r\v\fc ", "c b a"},
		{"punctuation stays in its word", "hello, world!", "world! hello,"},
		{"UTF-8 words intact", "naïve café", "café naïve"},
		{"non-ASCII space is not a separator", "a b c", "c a b"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ReverseWords(tt.s); got != tt.want {
				t.Errorf("ReverseWords(%q) = %q, want %q", tt.s, got, tt.want)
			}
		})
	}
}

func TestReverseWords_oneAllocation(t *testing.T) {
	var out string
	if allocs := testing.AllocsPerRun(100, func() { out = ReverseWords("the quick brown fox") }); allocs != 1 {
		t.Errorf("ReverseWords allocates %v times, want 1", allocs)
	}
	benchSink = len(out)
}

func TestReverseWords_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 500 {
		s := randomText(rng, rng.IntN(40), "ab \t\né")
		if got, want := ReverseWords(s), naiveReverseWords(s); got != want {
			t.Fatalf("trial %d: ReverseWords(%q) = %q, want %q", trial, s, got, want)
		}
	}
}

func BenchmarkReverseWords(b *testing.B) {
	rng := rand.New(rand.NewPCG(1, 2))
	for _, size := range benchSizes {
		s := randomText(rng, size.n, "abcdefg  ")
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				benchSink = len(ReverseWords(s))
			}
		})
	}
}
