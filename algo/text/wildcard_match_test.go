package text

import (
	"math/rand/v2"
	"strings"
	"testing"
)

// naiveWildcard is the oracle: direct recursion on runes (exponential in the worst case, so the oracle
// inputs stay short).
func naiveWildcard(text, pattern []rune) bool {
	if len(pattern) == 0 {
		return len(text) == 0
	}
	switch pattern[0] {
	case '*':
		return naiveWildcard(text, pattern[1:]) || (len(text) > 0 && naiveWildcard(text[1:], pattern))
	case '?':
		return len(text) > 0 && naiveWildcard(text[1:], pattern[1:])
	default:
		return len(text) > 0 && text[0] == pattern[0] && naiveWildcard(text[1:], pattern[1:])
	}
}

// POSIX.1-2017 XCU §2.13.1–2.13.2, applied by hand: ? is one character and * any string, so
// "report-2026.txt" matches "report-????.*" and "*.txt", but not "*.csv".
func TestWildcardMatch_posixExample(t *testing.T) {
	text := "report-2026.txt"
	for pattern, want := range map[string]bool{"report-????.*": true, "*.txt": true, "*.csv": false} {
		if got := WildcardMatch(text, pattern); got != want {
			t.Errorf("WildcardMatch(%q, %q) = %v, want %v", text, pattern, got, want)
		}
	}
}

// Our own wording of the usual exercise cases.
func TestWildcardMatch_cases(t *testing.T) {
	tests := []struct {
		name, text, pattern string
		want                bool
	}{
		{"both empty", "", "", true},
		{"empty pattern, text left", "a", "", false},
		{"star matches empty", "", "*", true},
		{"stars only", "", "***", true},
		{"question needs a rune", "", "?", false},
		{"literal", "abc", "abc", true},
		{"literal mismatch", "abc", "abd", false},
		{"one question", "abc", "a?c", true},
		{"too few runes", "ab", "a?c", false},
		{"star in the middle", "abxyzc", "ab*c", true},
		{"star must leave a suffix", "abxyz", "ab*c", false},
		{"must match all of text", "xaby", "ab", false},
		{"backtracking case", "aaab", "*a*b", true},
		{"star crosses slash", "dir/sub/file", "dir*file", true},
		{"no brackets", "[a]", "[a]", true},
		{"no escapes", `\*`, `\*`, true},
		{"question is one rune", "é", "?", true},
		{"question is not two runes", "éé", "?", false},
		{"invalid UTF-8 reads as U+FFFD", "\xff", "�", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WildcardMatch(tt.text, tt.pattern); got != tt.want {
				t.Errorf("WildcardMatch(%q, %q) = %v, want %v", tt.text, tt.pattern, got, tt.want)
			}
		})
	}
}

func TestWildcardMatch_oneAllocation(t *testing.T) {
	if allocs := testing.AllocsPerRun(100, func() { _ = WildcardMatch("report-2026.txt", "report-????.*") }); allocs > 1 {
		t.Errorf("WildcardMatch allocates %v times, want at most 1", allocs)
	}
}

func TestWildcardMatch_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 800 {
		text := randomText(rng, rng.IntN(10), "ab")
		pattern := randomText(rng, rng.IntN(7), "ab?*")
		if trial%5 == 0 {
			text = strings.ReplaceAll(text, "a", "é")
			pattern = strings.ReplaceAll(pattern, "a", "é")
		}
		want := naiveWildcard([]rune(text), []rune(pattern))
		if got := WildcardMatch(text, pattern); got != want {
			t.Fatalf("trial %d: WildcardMatch(%q, %q) = %v, want %v", trial, text, pattern, got, want)
		}
	}
}

// BenchmarkWildcardMatch: a 10-byte pattern with stars against n runes, so time grows as n·m = 10n.
func BenchmarkWildcardMatch(b *testing.B) {
	rng := rand.New(rand.NewPCG(1, 2))
	for _, size := range benchSizes {
		text := randomText(rng, size.n, "abc")
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				if WildcardMatch(text, "*a?b*c*a?*") {
					benchSink++
				}
			}
		})
	}
}
