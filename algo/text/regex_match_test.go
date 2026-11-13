package text

import (
	"errors"
	"math/rand/v2"
	"regexp"
	"strings"
	"testing"
)

// Aho, Lam, Sethi, Ullman, Compilers, 2nd ed., §3.3: concatenation and Kleene closure. Hand-derived on
// that semantics: a*b denotes {b, ab, aab, …}, so it matches "aaab" and "b" but not "aba".
func TestRegexMatch_closureExample(t *testing.T) {
	for text, want := range map[string]bool{"aaab": true, "b": true, "aba": false} {
		if got, err := RegexMatch(text, "a*b"); got != want || err != nil {
			t.Errorf("RegexMatch(%q, a*b) = %v, %v; want %v, nil", text, got, err, want)
		}
	}
}

// Our own wording of the usual exercise cases.
func TestRegexMatch_cases(t *testing.T) {
	tests := []struct {
		name, text, pattern string
		want                bool
	}{
		{"both empty", "", "", true},
		{"empty pattern, text left", "a", "", false},
		{"starred elements match empty", "", "x*y*", true},
		{"unstarred element needs a rune", "", "x", false},
		{"literal", "xyz", "xyz", true},
		{"one element is one rune", "xx", "x", false},
		{"star repeats", "xxy", "x*y", true},
		{"star allows zero", "y", "x*y", true},
		{"dot is any rune", "xyz", "x.z", true},
		{"dot star is anything", "anything at all", ".*", true},
		{"dot star then a suffix", "xyz", ".*z", true},
		{"whole text must match", "xyz", "x.", false},
		{"star binds to the rune before it", "abab", "ab*", false},
		{"many stars, polynomial time", strings.Repeat("a", 40), strings.Repeat("a*", 20) + "b", false},
		{"dot is one rune, not one byte", "é", ".", true},
		{"starred multibyte rune", "ééé", "é*", true},
		{"dot does not escape", "x", `\.`, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := RegexMatch(tt.text, tt.pattern); got != tt.want || err != nil {
				t.Errorf("RegexMatch(%q, %q) = %v, %v; want %v, nil", tt.text, tt.pattern, got, err, tt.want)
			}
		})
	}
}

func TestRegexMatch_errors(t *testing.T) {
	for _, pattern := range []string{"*", "*a", "a**", "a*b**", "é**"} {
		if got, err := RegexMatch("a", pattern); got || !errors.Is(err, ErrPattern) {
			t.Errorf("RegexMatch(a, %q) = %v, %v; want false, ErrPattern", pattern, got, err)
		}
	}
}

func TestRegexMatch_oneAllocation(t *testing.T) {
	if allocs := testing.AllocsPerRun(100, func() { _, _ = RegexMatch("xxxyz", "x*y.") }); allocs > 1 {
		t.Errorf("RegexMatch allocates %v times, want at most 1", allocs)
	}
}

// randomRegex returns a valid pattern of up to maxElems elements over a, b and '.', each starred or not.
func randomRegex(rng *rand.Rand, maxElems int) string {
	var b strings.Builder
	for range rng.IntN(maxElems + 1) {
		b.WriteByte("ab."[rng.IntN(3)])
		if rng.IntN(2) == 0 {
			b.WriteByte('*')
		}
	}
	return b.String()
}

// The oracle is the standard regexp package, anchored to the whole text; on this alphabet its '.' and
// '*' mean the same as ours.
func TestRegexMatch_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 800 {
		text := randomText(rng, rng.IntN(10), "ab")
		pattern := randomRegex(rng, 6)
		want := regexp.MustCompile("^(?:" + pattern + ")$").MatchString(text)
		if got, err := RegexMatch(text, pattern); got != want || err != nil {
			t.Fatalf("trial %d: RegexMatch(%q, %q) = %v, %v; want %v", trial, text, pattern, got, err, want)
		}
	}
}

// BenchmarkRegexMatch: a fixed 10-byte pattern against n runes, so time grows as n·m = 10n.
func BenchmarkRegexMatch(b *testing.B) {
	rng := rand.New(rand.NewPCG(1, 2))
	for _, size := range benchSizes {
		text := randomText(rng, size.n, "abc")
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				if ok, _ := RegexMatch(text, "a*.*b.c*.*"); ok {
					benchSink++
				}
			}
		})
	}
}
