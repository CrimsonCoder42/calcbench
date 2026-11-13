package text

import (
	"errors"
	"math/rand/v2"
	"slices"
	"testing"
)

// Hand-derived with the pattern of Cormen et al., 3rd ed., §32.4, Figure 32.10: ababaca occurs in
// abababacaba once, at offset 2.
func TestKMPSearch_example(t *testing.T) {
	got, err := KMPSearch(nil, "abababacaba", "ababaca")
	if err != nil || !slices.Equal(got, []int{2}) {
		t.Errorf("KMPSearch = %v, %v; want [2], nil", got, err)
	}
}

// Our own wording of the usual exercise cases.
func TestKMPSearch_cases(t *testing.T) {
	tests := []struct {
		name, text, pattern string
		want                []int
	}{
		{"empty text", "", "a", nil},
		{"pattern longer than text", "ab", "abc", nil},
		{"whole text", "abc", "abc", []int{0}},
		{"absent", "abcdef", "xy", nil},
		{"at the end", "hello world", "world", []int{6}},
		{"overlapping", "aaaa", "aa", []int{0, 1, 2}},
		{"first of several", "needle in a needle stack", "needle", []int{0, 12}},
		{"one byte", "banana", "a", []int{1, 3, 5}},
		{"UTF-8 byte offsets", "naïve naïve", "ïv", []int{2, 9}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := KMPSearch(nil, tt.text, tt.pattern)
			if err != nil || !slices.Equal(got, tt.want) {
				t.Errorf("KMPSearch(%q, %q) = %v, %v; want %v, nil", tt.text, tt.pattern, got, err, tt.want)
			}
		})
	}
}

func TestKMPSearch_emptyPattern(t *testing.T) {
	dst := []int{9}
	got, err := KMPSearch(dst, "abc", "")
	if !errors.Is(err, ErrEmptyPattern) || !slices.Equal(got, []int{9}) {
		t.Errorf("KMPSearch(empty pattern) = %v, %v; want [9], ErrEmptyPattern", got, err)
	}
}

func TestKMPSearch_bytes(t *testing.T) {
	got, err := KMPSearch(nil, []byte{0, 1, 0, 1, 0}, []byte{0, 1, 0})
	if err != nil || !slices.Equal(got, []int{0, 2}) {
		t.Errorf("KMPSearch(bytes) = %v, %v; want [0 2], nil", got, err)
	}
}

func TestKMPSearch_oneAllocation(t *testing.T) {
	buf := make([]int, 0, 16)
	allocs := testing.AllocsPerRun(100, func() { buf, _ = KMPSearch(buf[:0], "abababacabababaca", "ababaca") })
	if allocs > 1 {
		t.Errorf("KMPSearch allocates %v times, want at most 1 (the prefix function)", allocs)
	}
}

func TestKMPSearch_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 500 {
		text := randomText(rng, rng.IntN(60), "abc"[:1+trial%3])
		pattern := randomText(rng, 1+rng.IntN(5), "abc"[:1+trial%3])
		got, err := KMPSearch(nil, text, pattern)
		if want := naiveSearch(text, pattern); err != nil || !slices.Equal(got, want) {
			t.Fatalf("trial %d: KMPSearch(%q, %q) = %v, %v; want %v", trial, text, pattern, got, err, want)
		}
	}
}

func BenchmarkKMPSearch(b *testing.B) {
	rng := rand.New(rand.NewPCG(1, 2))
	for _, size := range benchSizes {
		text := randomText(rng, size.n, "ab")
		pattern := text[size.n/2 : size.n/2+16]
		buf := make([]int, 0, size.n)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				buf, _ = KMPSearch(buf[:0], text, pattern)
			}
			benchSink = len(buf)
		})
	}
}
