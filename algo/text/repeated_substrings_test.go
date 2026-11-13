package text

import (
	"errors"
	"math/rand/v2"
	"slices"
	"testing"
)

// naiveRepeated is the oracle: count every window in a map, then list the first starts of those seen twice.
func naiveRepeated(text string, length int) []int {
	first := map[string]int{}
	count := map[string]int{}
	for i := 0; i+length <= len(text); i++ {
		w := text[i : i+length]
		if _, ok := first[w]; !ok {
			first[w] = i
		}
		count[w]++
	}
	var out []int
	for w, c := range count {
		if c > 1 {
			out = append(out, first[w])
		}
	}
	slices.Sort(out)
	return out
}

// Hand-derived: in xyzxyzxy every length-3 window that starts at 0, 1 or 2 occurs again 3 bytes later.
func TestRepeatedSubstrings_example(t *testing.T) {
	got, err := RepeatedSubstrings(nil, "xyzxyzxy", 3)
	if err != nil || !slices.Equal(got, []int{0, 1, 2}) {
		t.Errorf("RepeatedSubstrings(xyzxyzxy, 3) = %v, %v; want [0 1 2], nil", got, err)
	}
}

// Our own wording of the usual exercise cases.
func TestRepeatedSubstrings_cases(t *testing.T) {
	tests := []struct {
		name   string
		text   string
		length int
		want   []int
	}{
		{"empty", "", 1, nil},
		{"longer than text", "abc", 4, nil},
		{"whole text once", "abc", 3, nil},
		{"no repeats", "abcdef", 2, nil},
		{"one byte repeated", "abca", 1, []int{0}},
		{"overlapping occurrences count", "aaaa", 3, []int{0}},
		{"reported once however many repeats", "ababababab", 2, []int{0, 1}},
		{"sorted by first occurrence", "CCGTTCCAGTT", 2, []int{0, 2, 3}},
		{"repeat seen late", "qwertyuiopqw", 2, []int{0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RepeatedSubstrings(nil, tt.text, tt.length)
			if err != nil || !slices.Equal(got, tt.want) {
				t.Errorf("RepeatedSubstrings(%q, %d) = %v, %v; want %v", tt.text, tt.length, got, err, tt.want)
			}
		})
	}
}

func TestRepeatedSubstrings_errors(t *testing.T) {
	for _, length := range []int{0, -1} {
		dst := []int{4}
		if got, err := RepeatedSubstrings(dst, "abab", length); !errors.Is(err, ErrLength) || !slices.Equal(got, []int{4}) {
			t.Errorf("length %d: %v, %v; want [4], ErrLength", length, got, err)
		}
	}
}

func TestRepeatedSubstrings_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 400 {
		text := randomText(rng, rng.IntN(80), "ACGT"[:1+trial%4])
		length := 1 + rng.IntN(6)
		got, err := RepeatedSubstrings(nil, []byte(text), length)
		if want := naiveRepeated(text, length); err != nil || !slices.Equal(got, want) {
			t.Fatalf("trial %d: RepeatedSubstrings(%q, %d) = %v, %v; want %v", trial, text, length, got, err, want)
		}
	}
}

func BenchmarkRepeatedSubstrings(b *testing.B) {
	rng := rand.New(rand.NewPCG(1, 2))
	for _, size := range benchSizes {
		text := randomText(rng, size.n, "ACGT")
		buf := make([]int, 0, size.n)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				buf, _ = RepeatedSubstrings(buf[:0], text, 10)
			}
			benchSink = len(buf)
		})
	}
}
