package enumerate

import (
	"errors"
	"fmt"
	"math/bits"
	"slices"
	"testing"
)

func gray(t *testing.T, b int) []uint64 {
	t.Helper()
	seq, err := GrayCodeSequence(b)
	if err != nil {
		t.Fatalf("GrayCodeSequence(%d): %v", b, err)
	}
	return slices.Collect(seq)
}

// Knuth, TAOCP Vol. 4A, §7.2.1.1: the binary reflected Gray code on 3 bits is 000 001 011 010 110 111
// 101 100.
func TestGrayCodeSequence_knuthExample(t *testing.T) {
	if got := gray(t, 3); !slices.Equal(got, []uint64{0, 1, 3, 2, 6, 7, 5, 4}) {
		t.Errorf("got %v", got)
	}
}

func TestGrayCodeSequence_cases(t *testing.T) {
	if got := gray(t, 0); !slices.Equal(got, []uint64{0}) {
		t.Errorf("0 bits: got %v, want [0]", got)
	}
	if got := gray(t, 1); !slices.Equal(got, []uint64{0, 1}) {
		t.Errorf("1 bit: got %v, want [0 1]", got)
	}
	if got := gray(t, 2); !slices.Equal(got, []uint64{0, 1, 3, 2}) {
		t.Errorf("2 bits: got %v", got)
	}
}

func TestGrayCodeSequence_errors(t *testing.T) {
	for _, b := range []int{-1, 65, 1 << 40} {
		if _, err := GrayCodeSequence(b); !errors.Is(err, ErrBitsRange) {
			t.Errorf("bits=%d: err = %v, want ErrBitsRange", b, err)
		}
	}
}

// 64 bits: the first words and early break (2^64 words never finish).
func TestGrayCodeSequence_64bitsEarlyBreak(t *testing.T) {
	seq, err := GrayCodeSequence(64)
	if err != nil {
		t.Fatal(err)
	}
	var got []uint64
	for g := range seq {
		got = append(got, g)
		if len(got) == 4 {
			break
		}
	}
	if !slices.Equal(got, []uint64{0, 1, 3, 2}) {
		t.Errorf("got %v", got)
	}
}

// TestGrayCodeSequence_oracle: for b ≤ 16, every word is below 2^b, all are distinct, and consecutive
// words (cyclically) differ in one bit.
func TestGrayCodeSequence_oracle(t *testing.T) {
	for b := 1; b <= 16; b++ {
		got := gray(t, b)
		if len(got) != 1<<b {
			t.Fatalf("b=%d: %d words", b, len(got))
		}
		seen := make([]bool, 1<<b)
		for i, g := range got {
			if g >= 1<<b || seen[g] {
				t.Fatalf("b=%d: word %d out of range or repeated", b, g)
			}
			seen[g] = true
			if d := bits.OnesCount64(g ^ got[(i+1)%len(got)]); d != 1 {
				t.Fatalf("b=%d: words %d and next differ in %d bits", b, i, d)
			}
		}
	}
}

// No allocation per word: a full pass allocates a fixed number of times (the iterator and the range
// loop's closures; 3 measured by the reviewer), the same for 16 words as for 65 536.
func TestGrayCodeSequence_noAllocationPerWord(t *testing.T) {
	pass := func(bits int) float64 {
		return testing.AllocsPerRun(5, func() {
			seq, _ := GrayCodeSequence(bits)
			for g := range seq {
				benchSink += int(g)
			}
		})
	}
	if small, large := pass(4), pass(16); small != large {
		t.Errorf("allocations per pass depend on bits: %v for 4 bits, %v for 16", small, large)
	}
}

func BenchmarkGrayCodeSequence(b *testing.B) {
	for _, n := range []int{10, 14, 17} {
		b.Run(fmt.Sprint(n), func(b *testing.B) {
			b.ReportAllocs()
			sink := uint64(0)
			for b.Loop() {
				seq, _ := GrayCodeSequence(n)
				for g := range seq {
					sink += g
				}
			}
			benchSink = int(sink)
		})
	}
}
