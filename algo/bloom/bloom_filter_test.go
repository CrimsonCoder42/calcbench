package bloom

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"math"
	"math/big"
	"math/rand/v2"
	"slices"
	"testing"
)

// benchSizes are the collection sizes of .claude/rules/go-complexity.md. Shared by the benchmarks.
var benchSizes = []struct {
	name string
	n    int
}{{"1e3", 1_000}, {"1e4", 10_000}, {"1e5", 100_000}}

func mustFilter(t testing.TB, bitCount uint64, hashCount int) *BloomFilter {
	t.Helper()
	f, err := NewBloomFilter(bitCount, hashCount)
	if err != nil {
		t.Fatal(err)
	}
	return f
}

// slowIndices is the oracle for the k indices: (h₁ + i·h₂ + (i³ − i)/6) mod m in arbitrary precision,
// the closed form of Dillinger–Manolios §5.2 eq. (6), with h₁, h₂ the first two big-endian words of
// SHA-256(item).
func slowIndices(item []byte, m uint64, k int) []uint64 {
	d := sha256.Sum256(item)
	h1 := new(big.Int).SetBytes(d[0:8])
	h2 := new(big.Int).SetBytes(d[8:16])
	mod := new(big.Int).SetUint64(m)
	out := make([]uint64, k)
	for i := range k {
		g := new(big.Int).Mul(big.NewInt(int64(i)), h2)
		g.Add(g, h1).Add(g, big.NewInt(int64((i*i*i-i)/6))).Mod(g, mod)
		out[i] = g.Uint64()
	}
	return out
}

// setBitList lists the indices of the 1 bits.
func setBitList(f *BloomFilter) []uint64 {
	var out []uint64
	for i := range f.bitCount {
		if f.words[i>>6]&(1<<(i&63)) != 0 {
			out = append(out, i)
		}
	}
	return out
}

// Worked example (catalog/bloom.md): m = 64, k = 3. The indices were computed independently with Python's
// hashlib.sha256 and (h₁ + i·h₂ + (i³ − i)/6) mod 64: "apple" → 41, 48, 56; "banana" → 13, 35, 58.
func TestBloomFilter_example(t *testing.T) {
	f := mustFilter(t, 64, 3)
	for word, want := range map[string][]uint64{"apple": {41, 48, 56}, "banana": {13, 35, 58}} {
		if got := slowIndices([]byte(word), 64, 3); !slices.Equal(got, want) {
			t.Fatalf("indices of %q = %v, want %v", word, got, want)
		}
	}
	f.Add([]byte("apple"))
	f.Add([]byte("banana"))
	if got, want := setBitList(f), []uint64{13, 35, 41, 48, 56, 58}; !slices.Equal(got, want) {
		t.Errorf("set bits = %v, want %v", got, want)
	}
	if f.SetBits() != 6 || f.BitCount() != 64 || f.HashCount() != 3 {
		t.Errorf("SetBits, BitCount, HashCount = %d, %d, %d, want 6, 64, 3", f.SetBits(), f.BitCount(), f.HashCount())
	}
	if !f.MayContain([]byte("apple")) || !f.MayContain([]byte("banana")) {
		t.Error("an added item tests absent")
	}
	if got, want := f.EstimatedFPR(), math.Pow(6.0/64, 3); got != want {
		t.Errorf("EstimatedFPR = %v, want %v", got, want)
	}
}

func TestBloomFilter_cases(t *testing.T) {
	t.Run("empty filter holds nothing", func(t *testing.T) {
		f := mustFilter(t, 1000, 5)
		if f.MayContain([]byte("x")) || f.MayContain(nil) || f.EstimatedFPR() != 0 || f.SetBits() != 0 {
			t.Error("a new filter reports content")
		}
	})
	t.Run("nil and empty items are the same item", func(t *testing.T) {
		f := mustFilter(t, 1000, 5)
		f.Add(nil)
		if !f.MayContain([]byte{}) {
			t.Error("empty item absent after Add(nil)")
		}
	})
	t.Run("duplicates change nothing", func(t *testing.T) {
		f := mustFilter(t, 1000, 5)
		f.Add([]byte("dup"))
		before := slices.Clone(f.words)
		f.Add([]byte("dup"))
		if !slices.Equal(before, f.words) || f.SetBits() > 5 {
			t.Error("adding an item twice changed the bits")
		}
	})
	t.Run("one bit", func(t *testing.T) {
		f := mustFilter(t, 1, 1)
		f.Add([]byte("a"))
		if f.SetBits() != 1 || !f.MayContain([]byte("anything")) || f.EstimatedFPR() != 1 {
			t.Errorf("one-bit filter: SetBits %d, EstimatedFPR %v", f.SetBits(), f.EstimatedFPR())
		}
	})
	t.Run("as many hashes as bits", func(t *testing.T) {
		f := mustFilter(t, 64, 64)
		f.Add([]byte("a"))
		if f.SetBits() > 64 || !f.MayContain([]byte("a")) {
			t.Errorf("SetBits %d", f.SetBits())
		}
	})
	t.Run("MaxHashCount", func(t *testing.T) {
		f := mustFilter(t, 4096, MaxHashCount)
		f.Add([]byte("a"))
		if !f.MayContain([]byte("a")) {
			t.Error("false negative at MaxHashCount")
		}
	})
	t.Run("bit count not a multiple of 64", func(t *testing.T) {
		f := mustFilter(t, 65, 3)
		for i := range 200 {
			f.Add(fmt.Appendf(nil, "item-%d", i))
		}
		if f.words[1]>>1 != 0 {
			t.Errorf("bits beyond m set: last word %b", f.words[1])
		}
	})
	t.Run("input not modified", func(t *testing.T) {
		f := mustFilter(t, 100, 3)
		item := []byte("keep")
		f.Add(item)
		f.MayContain(item)
		if string(item) != "keep" {
			t.Errorf("item modified: %q", item)
		}
	})
	t.Run("zero value", func(t *testing.T) {
		var f BloomFilter
		f.Add([]byte("a"))
		if f.MayContain([]byte("a")) || f.EstimatedFPR() != 0 || f.SetBits() != 0 || f.BitCount() != 0 {
			t.Error("zero value holds something")
		}
	})
}

func TestBloomFilter_errors(t *testing.T) {
	tests := []struct {
		name      string
		bitCount  uint64
		hashCount int
		wantErr   error
	}{
		{"no bits", 0, 3, ErrBitCount},
		{"no hashes", 64, 0, ErrHashCount},
		{"negative hashes", 64, math.MinInt, ErrHashCount},
		{"more hashes than bits", 64, 65, ErrHashCount},
		{"one bit, two hashes", 1, 2, ErrHashCount},
		{"above MaxHashCount", 1 << 20, MaxHashCount + 1, ErrHashCount},
		{"maximum of int hashes", 1 << 20, math.MaxInt, ErrHashCount},
		{"above MaxBits", MaxBits + 1, 3, ErrTooLarge},
		{"maximum of uint64", math.MaxUint64, 3, ErrTooLarge},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := NewBloomFilter(tt.bitCount, tt.hashCount)
			if f != nil || !errors.Is(err, tt.wantErr) {
				t.Errorf("NewBloomFilter(%d, %d) = (%v, %v), want (nil, %v)", tt.bitCount, tt.hashCount, f, err, tt.wantErr)
			}
		})
	}
}

// Oracle: a filter built from slowIndices. Same bits, same SetBits, no false negatives, and MayContain
// agrees with "all k slow indices set" on fresh probes.
func TestBloomFilter_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	randItem := func() []byte {
		b := make([]byte, rng.IntN(40))
		for i := range b {
			b[i] = byte(rng.Uint32())
		}
		return b
	}
	for trial := range 300 {
		m := 1 + rng.Uint64N(5000)
		k := 1 + rng.IntN(20)
		f := mustFilter(t, m, k)
		slow := make(map[uint64]bool)
		items := make([][]byte, rng.IntN(100))
		for i := range items {
			items[i] = randItem()
			f.Add(items[i])
			for _, g := range slowIndices(items[i], m, k) {
				slow[g] = true
			}
		}
		got := setBitList(f)
		if len(got) != len(slow) || f.SetBits() != uint64(len(slow)) {
			t.Fatalf("trial %d (m=%d, k=%d): %d bits set (SetBits %d), oracle %d", trial, m, k, len(got), f.SetBits(), len(slow))
		}
		for _, g := range got {
			if !slow[g] {
				t.Fatalf("trial %d: bit %d set, oracle has it clear", trial, g)
			}
		}
		for _, item := range items {
			if !f.MayContain(item) {
				t.Fatalf("trial %d: false negative for %x", trial, item)
			}
		}
		for range 50 {
			probe := randItem()
			want := true
			for _, g := range slowIndices(probe, m, k) {
				want = want && slow[g]
			}
			if f.MayContain(probe) != want {
				t.Fatalf("trial %d: MayContain(%x) = %v, oracle %v", trial, probe, !want, want)
			}
		}
	}
}

// The measured false-positive rate matches the prediction. Deterministic (SHA-256 is fixed): a Python
// model of the same construction (enhanced double hashing) measured 0.01038 and EstimatedFPR 0.01008. The band is ±20% of the
// predicted 0.0100: binomial σ for 100 000 probes is 0.0003, so the band only fails on a real defect.
func TestBloomFilter_falsePositiveRate(t *testing.T) {
	const n, probes = 1000, 100_000
	m, k, err := OptimalParams(n, 0.01)
	if err != nil {
		t.Fatal(err)
	}
	f := mustFilter(t, m, k)
	for i := range n {
		f.Add(fmt.Appendf(nil, "member-%d", i))
	}
	fp := 0
	for i := range probes {
		if f.MayContain(fmt.Appendf(nil, "other-%d", i)) {
			fp++
		}
	}
	predicted, _ := FalsePositiveRate(m, n, k)
	measured := float64(fp) / probes
	if measured < 0.8*predicted || measured > 1.2*predicted {
		t.Errorf("measured FPR %v, predicted %v (m=%d, k=%d)", measured, predicted, m, k)
	}
	// The fill-based estimate concentrates around the prediction (KM §2); 5% covers its spread at n = 1000.
	if est := f.EstimatedFPR(); math.Abs(est-predicted) > 0.05*predicted {
		t.Errorf("EstimatedFPR %v, predicted %v", est, predicted)
	}
}

// Plain double hashing h₁ + i·h₂ collapses when the step h₂ mod m is 0 (one index) or shares a factor
// with m (Dillinger–Manolios §5.1); nextProbe's added (i³ − i)/6 keeps the indices apart. Sequences and
// distinct counts computed independently in Python (catalog/bloom.md, M7).
func TestBloomFilter_stepZeroOrSharedFactor(t *testing.T) {
	tests := []struct {
		name           string
		first, step, m uint64
		k              int
		want           []uint64
		minDistinct    int
	}{
		{"step 0, prime m", 5, 0, 1009, 10, []uint64{5, 5, 6, 9, 15, 25, 40, 61, 89, 125}, 9},
		{"step 32 divides m = 64", 5, 32, 64, 8, []uint64{5, 37, 6, 41, 15, 57, 40, 29}, 8},
		{"step 0 at index 0", 0, 0, 64, 8, []uint64{0, 0, 1, 4, 10, 20, 35, 56}, 7},
		{"one bit", 0, 0, 1, 1, []uint64{0}, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := make([]uint64, 0, tt.k)
			g, step := tt.first, tt.step
			for i := range tt.k {
				got = append(got, g)
				g, step = nextProbe(g, step, uint64(i), tt.m)
				if g >= tt.m || step >= tt.m {
					t.Fatalf("probe %d left [0, m): g = %d, step = %d", i, g, step)
				}
			}
			distinct := len(slices.Compact(slices.Sorted(slices.Values(got))))
			if !slices.Equal(got, tt.want) || distinct < tt.minDistinct {
				t.Errorf("indices = %v (%d distinct), want %v (at least %d)", got, distinct, tt.want, tt.minDistinct)
			}
		})
	}
	// The step after probe 0 is never 0 when it started at 0 and m > 1: it becomes 1.
	if _, step := nextProbe(7, 0, 0, 64); step != 1 {
		t.Errorf("step after a zero step = %d, want 1", step)
	}
}

func TestBloomFilter_noAllocation(t *testing.T) {
	f := mustFilter(t, 10_000, 7)
	item := []byte("a reasonably ordinary key")
	allocs := testing.AllocsPerRun(100, func() {
		f.Add(item)
		f.MayContain(item)
		f.EstimatedFPR()
	})
	if allocs != 0 {
		t.Errorf("Add and MayContain allocate %v times per call, want 0", allocs)
	}
}

// Each iteration adds size items and tests size items in a filter sized for them at 1%: O(size·(L + k)).
func BenchmarkBloomFilter(b *testing.B) {
	for _, size := range benchSizes {
		m, k, _ := OptimalParams(uint64(size.n), 0.01)
		f := mustFilter(b, m, k)
		items := make([][]byte, size.n)
		for i := range items {
			items[i] = fmt.Appendf(nil, "key-%08d", i)
		}
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				for _, item := range items {
					f.Add(item)
				}
				for _, item := range items {
					f.MayContain(item)
				}
			}
		})
	}
}
