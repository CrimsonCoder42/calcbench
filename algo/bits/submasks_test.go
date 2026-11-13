package bits

import (
	mathbits "math/bits"
	"math/rand/v2"
	"slices"
	"testing"
)

// Knuth Vol. 4A §7.1.3: the subsets of {0, 2, 3} as masks of 0b1101, from the full set down to 0.
func TestSubmasks_sourceExample(t *testing.T) {
	got := slices.Collect(Submasks(0b1101))
	want := []uint64{0b1101, 0b1100, 0b1001, 0b1000, 0b0101, 0b0100, 0b0001, 0}
	if !slices.Equal(got, want) {
		t.Errorf("Submasks(0b1101) = %b, want %b", got, want)
	}
}

func TestSubmasks_cases(t *testing.T) {
	if got := slices.Collect(Submasks(0)); !slices.Equal(got, []uint64{0}) {
		t.Errorf("Submasks(0) = %v, want [0]", got)
	}
	if got := slices.Collect(Submasks(1 << 63)); !slices.Equal(got, []uint64{1 << 63, 0}) {
		t.Errorf("Submasks(2^63) = %v", got)
	}
	n := 0
	for range Submasks(0xFFFF) {
		n++
	}
	if n != 1<<16 {
		t.Errorf("Submasks(0xFFFF) yields %d values, want 65536", n)
	}
}

func TestSubmasks_earlyStop(t *testing.T) {
	var got []uint64
	for s := range Submasks(0b111) {
		got = append(got, s)
		if len(got) == 3 {
			break
		}
	}
	if !slices.Equal(got, []uint64{0b111, 0b110, 0b101}) {
		t.Errorf("early stop got %b", got)
	}
}

// Oracle: every value 0..mask that is a submask, in descending order, on masks with at most 12 bits set
// placed anywhere in the word.
func TestSubmasks_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for range 200 {
		var mask uint64
		for range rng.IntN(13) {
			mask |= 1 << rng.IntN(64)
		}
		got := slices.Collect(Submasks(mask))
		if len(got) != 1<<mathbits.OnesCount64(mask) {
			t.Fatalf("Submasks(%#x) has %d values, want %d", mask, len(got), 1<<mathbits.OnesCount64(mask))
		}
		for i, s := range got {
			if s&^mask != 0 || (i > 0 && s >= got[i-1]) {
				t.Fatalf("Submasks(%#x)[%d] = %#x: not a submask or not descending", mask, i, s)
			}
		}
	}
}

func BenchmarkSubmasks(b *testing.B) {
	for _, c := range []struct {
		name string
		mask uint64
	}{{"k=10", 0x3FF}, {"k=14", 0x3FFF}, {"k=17", 0x1FFFF}} {
		b.Run(c.name, func(b *testing.B) {
			b.ReportAllocs()
			var s uint64
			for b.Loop() {
				for m := range Submasks(c.mask) {
					s += m
				}
			}
			benchSink = s
		})
	}
}
