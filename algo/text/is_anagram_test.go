package text

import (
	"math/rand/v2"
	"slices"
	"strings"
	"testing"
)

// sortedRunes is the oracle's signature (Bentley, Programming Pearls, 2nd ed., Column 2): the runes, sorted.
func sortedRunes(s string) string {
	r := []rune(s)
	slices.Sort(r)
	return string(r)
}

// Bentley, Programming Pearls, 2nd ed., Column 2: words with the same signature are anagrams, e.g. pots,
// stop and tops (signature opst).
func TestIsAnagram_bentleyExample(t *testing.T) {
	if !IsAnagram("pots", "stop") || !IsAnagram("stop", "tops") {
		t.Error("pots, stop and tops should be anagrams")
	}
}

// Our own wording of the usual exercise cases.
func TestIsAnagram_cases(t *testing.T) {
	tests := []struct {
		name string
		a, b string
		want bool
	}{
		{"both empty", "", "", true},
		{"one empty", "a", "", false},
		{"same string", "listen", "listen", true},
		{"rearranged", "listen", "silent", true},
		{"different counts", "aab", "abb", false},
		{"different lengths", "abc", "abcc", false},
		{"case matters", "Ab", "ba", false},
		{"spaces count", "a b", "ab", false},
		{"accented runes", "résumé", "émusér", true},
		{"accent differs", "résumé", "resume", false},
		{"mixed ASCII and not", "añob", "boña", true},
		{"invalid UTF-8 reads as U+FFFD", "\xffa", "a�", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsAnagram(tt.a, tt.b); got != tt.want {
				t.Errorf("IsAnagram(%q, %q) = %v, want %v", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestIsAnagram_noAllocationASCII(t *testing.T) {
	allocs := testing.AllocsPerRun(100, func() {
		if IsAnagram("dormitory", "dirtyroom") {
			benchSink++
		}
	})
	if allocs != 0 {
		t.Errorf("IsAnagram on ASCII allocates %v times, want 0", allocs)
	}
}

func TestIsAnagram_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	alphabets := []string{"ab", "abc", "aé"}
	for trial := range 600 {
		alpha := alphabets[trial%3]
		a := randomText(rng, rng.IntN(10), alpha)
		var b string
		if trial%2 == 0 { // a shuffle of a, so true cases are common
			bs := []byte(a)
			rng.Shuffle(len(bs), func(i, j int) { bs[i], bs[j] = bs[j], bs[i] })
			b = string(bs)
		} else {
			b = randomText(rng, rng.IntN(10), alpha)
		}
		want := sortedRunes(a) == sortedRunes(b)
		if got := IsAnagram(a, b); got != want {
			t.Fatalf("trial %d: IsAnagram(%q, %q) = %v, want %v", trial, a, b, got, want)
		}
	}
}

func BenchmarkIsAnagram(b *testing.B) {
	rng := rand.New(rand.NewPCG(1, 2))
	for _, size := range benchSizes {
		a := randomText(rng, size.n, "abcdefghijklmnopqrstuvwxyz")
		bs := []byte(a)
		slices.Reverse(bs)
		rev := string(bs)
		b.Run(size.name+"/ascii", func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				if IsAnagram(a, rev) {
					benchSink++
				}
			}
		})
		ua := strings.Repeat("é", size.n/2)
		b.Run(size.name+"/utf8", func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				if IsAnagram(ua, ua) {
					benchSink++
				}
			}
		})
	}
}
