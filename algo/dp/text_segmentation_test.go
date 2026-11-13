package dp

import (
	"math/rand/v2"
	"slices"
	"strings"
	"testing"
)

// splittableSlow is Erickson's backtracking Splittable(i) without memoisation (short texts only).
func splittableSlow(text string, words []string) bool {
	if text == "" {
		return true
	}
	for _, w := range words {
		if w != "" && strings.HasPrefix(text, w) && splittableSlow(text[len(w):], words) {
			return true
		}
	}
	return false
}

func checkPieces(t *testing.T, text string, words []string, pieces []string) {
	t.Helper()
	if strings.Join(pieces, "") != text {
		t.Fatalf("pieces %q do not concatenate to %q", pieces, text)
	}
	for _, p := range pieces {
		if p == "" || !slices.Contains(words, p) {
			t.Fatalf("piece %q of %q is not a dictionary word", p, pieces)
		}
	}
}

// Erickson, Algorithms, §2.5: the run-together text BOTHEARTHANDSATURNSPIN splits into words.
func TestTextSegmentation_ericksonExample(t *testing.T) {
	text := "BOTHEARTHANDSATURNSPIN"
	words := []string{"BOTH", "EARTH", "AND", "SATURN", "SPIN", "BOT", "HEART", "HANDS", "AT", "URNS", "PIN", "HE", "ART", "TURN"}
	pieces, ok := TextSegmentation(text, words)
	if !ok {
		t.Fatalf("TextSegmentation(%q) found no split", text)
	}
	checkPieces(t, text, words, pieces)
}

func TestTextSegmentation_cases(t *testing.T) {
	tests := []struct {
		name  string
		text  string
		words []string
		want  []string
		ok    bool
	}{
		{"empty text", "", []string{"a"}, []string{}, true},
		{"empty text, no words", "", nil, []string{}, true},
		{"no words", "abc", nil, nil, false},
		{"one word", "sun", []string{"sun"}, []string{"sun"}, true},
		{"no split", "abc", []string{"ab"}, nil, false},
		{"reuse a word", "aaa", []string{"a"}, []string{"a", "a", "a"}, true},
		{"empty dictionary word ignored", "aa", []string{"", "a"}, []string{"a", "a"}, true},
		{"longest last piece wins", "sunflowerseed", []string{"sun", "flower", "sunflower", "seed", "s", "eed"}, []string{"sunflower", "seed"}, true},
		{"needs backtracking", "abcd", []string{"a", "abc", "b", "cd"}, []string{"a", "b", "cd"}, true},
		{"duplicate dictionary words", "xyxy", []string{"xy", "xy"}, []string{"xy", "xy"}, true},
		{"multibyte bytes", "héhé", []string{"hé"}, []string{"hé", "hé"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := TextSegmentation(tt.text, tt.words)
			if ok != tt.ok || !slices.Equal(got, tt.want) || (tt.ok && got == nil) {
				t.Errorf("TextSegmentation(%q, %q) = (%#v, %v), want (%q, %v)", tt.text, tt.words, got, ok, tt.want, tt.ok)
			}
		})
	}
}

func TestTextSegmentation_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 400 {
		words := make([]string, rng.IntN(5))
		for i := range words {
			words[i] = string(randomBytes(rng, 1+rng.IntN(3), "ab"))
		}
		text := string(randomBytes(rng, rng.IntN(14), "ab"))
		pieces, ok := TextSegmentation(text, words)
		if want := splittableSlow(text, words); ok != want {
			t.Fatalf("trial %d: TextSegmentation(%q, %q) ok = %v, want %v", trial, text, words, ok, want)
		}
		if ok {
			checkPieces(t, text, words, pieces)
		}
	}
}

func BenchmarkTextSegmentation(b *testing.B) {
	words := []string{"a", "b", "ab", "ba", "aab", "abb", "baba", "bbbbbbbbbb"}
	for _, size := range linearSizes {
		rng := rand.New(rand.NewPCG(1, 2))
		text := string(randomBytes(rng, size.n, "ab"))
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _ = TextSegmentation(text, words)
			}
		})
	}
}
