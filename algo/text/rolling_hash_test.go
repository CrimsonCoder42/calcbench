package text

import (
	"errors"
	"math/rand/v2"
	"testing"
)

// clrsDigits is the text of Cormen et al., Introduction to Algorithms, 3rd ed., §32.2, Figure 32.5.
const clrsDigits = "2359023141526739921"

// Cormen et al., 3rd ed., §32.2, Figure 32.5: with radix 10 and q = 13, the windows 31415 (shift 6) and
// 67399 (shift 12) both hash to 7. Over bytes ('0' = 48 …) every window's hash is shifted by the same
// constant, so the windows that collide in the figure collide here too, and exactly when their decimal
// values are congruent mod 13.
func TestNewRollingHash_clrsFigure32_5(t *testing.T) {
	h, err := NewRollingHash(10, 13, 5)
	if err != nil {
		t.Fatal(err)
	}
	var sums []uint64
	for i := range 5 {
		h.Append(clrsDigits[i])
	}
	sums = append(sums, h.Sum())
	for s := 0; s+5 < len(clrsDigits); s++ {
		h.Roll(clrsDigits[s], clrsDigits[s+5])
		sums = append(sums, h.Sum())
	}
	if sums[6] != sums[12] {
		t.Errorf("windows 31415 and 67399 hash to %d and %d, want equal", sums[6], sums[12])
	}
	decimalMod13 := func(w string) int {
		v := 0
		for i := range len(w) {
			v = (v*10 + int(w[i]-'0')) % 13
		}
		return v
	}
	if got := decimalMod13("31415"); got != 7 {
		t.Fatalf("31415 mod 13 = %d, the figure says 7", got)
	}
	for i := range sums {
		for j := range sums {
			same := decimalMod13(clrsDigits[i:i+5]) == decimalMod13(clrsDigits[j:j+5])
			if (sums[i] == sums[j]) != same {
				t.Errorf("windows %d and %d: equal hashes %v, equal residues %v", i, j, sums[i] == sums[j], same)
			}
		}
	}
}

func TestNewRollingHash_errors(t *testing.T) {
	tests := []struct {
		name          string
		base, modulus uint64
		windowLen     int
		want          error
	}{
		{"modulus 0", 1, 0, 1, ErrHashParams},
		{"modulus 1", 1, 1, 1, ErrHashParams},
		{"base 0", 0, 13, 1, ErrHashParams},
		{"base equals modulus", 13, 13, 1, ErrHashParams},
		{"window 0", 10, 13, 0, ErrWindowLen},
		{"window negative", 10, 13, -4, ErrWindowLen},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := NewRollingHash(tt.base, tt.modulus, tt.windowLen); !errors.Is(err, tt.want) {
				t.Errorf("NewRollingHash(%d, %d, %d) error = %v, want %v", tt.base, tt.modulus, tt.windowLen, err, tt.want)
			}
		})
	}
}

func TestRollingHash_lenAndReset(t *testing.T) {
	h, _ := NewRollingHash(DefaultRollingBase, DefaultRollingModulus, 3)
	h.Append('a')
	h.Append('b')
	if h.Len() != 2 {
		t.Errorf("Len = %d, want 2", h.Len())
	}
	h.Reset()
	if h.Len() != 0 || h.Sum() != 0 {
		t.Errorf("after Reset: Len = %d, Sum = %d; want 0, 0", h.Len(), h.Sum())
	}
}

// Every rolled fingerprint equals the fingerprint of the same window built from scratch, for small
// moduli, the default Mersenne prime, and the largest 64-bit prime (where a+b overflows 64 bits).
func TestRollingHash_rollMatchesFresh(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	params := []struct{ base, modulus uint64 }{
		{10, 13}, {256, 257}, {DefaultRollingBase, DefaultRollingModulus},
		{1<<64 - 60, 1<<64 - 59}, {1 << 63, 1<<64 - 59},
	}
	for _, p := range params {
		for trial := range 50 {
			text := randomText(rng, 40, "\x00\x01\xfeab")
			w := 1 + trial%8
			h, err := NewRollingHash(p.base, p.modulus, w)
			if err != nil {
				t.Fatal(err)
			}
			for i := range w {
				h.Append(text[i])
			}
			for s := 0; ; s++ {
				fresh, _ := NewRollingHash(p.base, p.modulus, w)
				for i := s; i < s+w; i++ {
					fresh.Append(text[i])
				}
				if h.Sum() != fresh.Sum() || h.Sum() >= p.modulus {
					t.Fatalf("base %d mod %d window %d at %d: rolled %d, fresh %d", p.base, p.modulus, w, s, h.Sum(), fresh.Sum())
				}
				if s+w >= len(text) {
					break
				}
				h.Roll(text[s], text[s+w])
			}
		}
	}
}

func TestRollingHash_noAllocation(t *testing.T) {
	allocs := testing.AllocsPerRun(100, func() {
		h, _ := NewRollingHash(DefaultRollingBase, DefaultRollingModulus, 4)
		for i := range 4 {
			h.Append(clrsDigits[i])
		}
		h.Roll(clrsDigits[0], clrsDigits[4])
		benchSink = int(h.Sum() & 1)
	})
	if allocs != 0 {
		t.Errorf("RollingHash allocates %v times, want 0", allocs)
	}
}

// BenchmarkNewRollingHash measures a full slide over n bytes: n Roll calls, O(1) each.
func BenchmarkNewRollingHash(b *testing.B) {
	rng := rand.New(rand.NewPCG(1, 2))
	for _, size := range benchSizes {
		text := randomText(rng, size.n, "abcdefgh")
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			var sum uint64
			for b.Loop() {
				h, _ := NewRollingHash(DefaultRollingBase, DefaultRollingModulus, 16)
				for i := range 16 {
					h.Append(text[i])
				}
				for s := 0; s+16 < len(text); s++ {
					h.Roll(text[s], text[s+16])
				}
				sum += h.Sum()
			}
			benchSink = int(sum & 1)
		})
	}
}

// The zero value is inert: Append and Roll never divide by its modulus 0, and Sum stays 0.
func TestRollingHash_zeroValueIsInert(t *testing.T) {
	var h RollingHash
	h.Append('a')
	h.Roll('a', 'b')
	if h.Sum() != 0 || h.Len() != 0 {
		t.Errorf("zero value after Append and Roll: Sum = %d, Len = %d; want 0, 0", h.Sum(), h.Len())
	}
}
