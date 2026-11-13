package text

import (
	"math/rand/v2"
	"strings"
	"testing"
)

// naiveLCP is the oracle: the longest prefix of the first string that every string starts with.
func naiveLCP(strs []string) int {
	if len(strs) == 0 {
		return 0
	}
	for p := len(strs[0]); p >= 0; p-- {
		all := true
		for _, s := range strs {
			if !strings.HasPrefix(s, strs[0][:p]) {
				all = false
				break
			}
		}
		if all {
			return p
		}
	}
	return 0
}

// Hand-derived: carpet, carpool and cartoon share "car".
func TestLongestCommonPrefix_example(t *testing.T) {
	if got := LongestCommonPrefix([]string{"carpet", "carpool", "cartoon"}); got != 3 {
		t.Errorf("LongestCommonPrefix = %d, want 3", got)
	}
}

// Our own wording of the usual exercise cases.
func TestLongestCommonPrefix_cases(t *testing.T) {
	tests := []struct {
		name string
		strs []string
		want int
	}{
		{"no strings", nil, 0},
		{"one string", []string{"solo"}, 4},
		{"two equal", []string{"same", "same"}, 4},
		{"nothing shared", []string{"dog", "cat"}, 0},
		{"an empty string", []string{"abc", ""}, 0},
		{"one is a prefix of the other", []string{"inter", "internet"}, 5},
		{"shorter string last", []string{"internet", "inter"}, 5},
		{"differs at the last string", []string{"abc", "abc", "abd"}, 2},
		{"UTF-8 stops inside a rune", []string{"é", "è"}, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := LongestCommonPrefix(tt.strs); got != tt.want {
				t.Errorf("LongestCommonPrefix(%q) = %d, want %d", tt.strs, got, tt.want)
			}
		})
	}
}

func TestLongestCommonPrefix_bytes(t *testing.T) {
	if got := LongestCommonPrefix([][]byte{{1, 2, 3}, {1, 2, 4}}); got != 2 {
		t.Errorf("LongestCommonPrefix(bytes) = %d, want 2", got)
	}
}

func TestLongestCommonPrefix_noAllocation(t *testing.T) {
	strs := []string{"carpet", "carpool", "cartoon"}
	if allocs := testing.AllocsPerRun(100, func() { benchSink = LongestCommonPrefix(strs) }); allocs != 0 {
		t.Errorf("LongestCommonPrefix allocates %v times, want 0", allocs)
	}
}

func TestLongestCommonPrefix_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 500 {
		strs := make([]string, rng.IntN(6))
		for i := range strs {
			strs[i] = randomText(rng, rng.IntN(8), "ab"[:1+trial%2])
		}
		if got, want := LongestCommonPrefix(strs), naiveLCP(strs); got != want {
			t.Fatalf("trial %d: LongestCommonPrefix(%q) = %d, want %d", trial, strs, got, want)
		}
	}
}

// BenchmarkLongestCommonPrefix uses n strings sharing a 16-byte prefix.
func BenchmarkLongestCommonPrefix(b *testing.B) {
	rng := rand.New(rand.NewPCG(1, 2))
	for _, size := range benchSizes {
		strs := make([]string, size.n)
		for i := range strs {
			strs[i] = "common-prefix-16" + randomText(rng, 8, "xyz")
		}
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				benchSink = LongestCommonPrefix(strs)
			}
		})
	}
}
