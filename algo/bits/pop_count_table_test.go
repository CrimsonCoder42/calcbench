package bits

import (
	"errors"
	"math"
	mathbits "math/bits"
	"slices"
	"testing"
)

// Knuth Vol. 4A §7.1.3, sideways addition: ν(0..8) = 0 1 1 2 1 2 2 3 1.
func TestPopCountTable_sourceExample(t *testing.T) {
	got, err := PopCountTable(nil, 8)
	if want := []uint8{0, 1, 1, 2, 1, 2, 2, 3, 1}; err != nil || !slices.Equal(got, want) {
		t.Errorf("PopCountTable(nil, 8) = (%v, %v), want %v", got, err, want)
	}
}

func TestPopCountTable_cases(t *testing.T) {
	if got, err := PopCountTable(nil, 0); err != nil || !slices.Equal(got, []uint8{0}) {
		t.Errorf("n = 0: (%v, %v), want [0]", got, err)
	}
	if got, err := PopCountTable(nil, 5); err != nil || !slices.Equal(got, []uint8{0, 1, 1, 2, 1, 2}) {
		t.Errorf("n = 5: (%v, %v)", got, err)
	}
	// Appends after existing entries; the recurrence must index from the start of this table.
	if got, err := PopCountTable([]uint8{9, 9}, 3); err != nil || !slices.Equal(got, []uint8{9, 9, 0, 1, 1, 2}) {
		t.Errorf("append: (%v, %v), want [9 9 0 1 1 2]", got, err)
	}
}

// Rule L41: a table beyond one allocation's limit is ErrTooLarge, never a makeslice panic.
func TestPopCountTable_errors(t *testing.T) {
	dst := []uint8{7}
	for _, tt := range []struct {
		n       int
		wantErr error
	}{{-1, ErrNegative}, {math.MinInt, ErrNegative}, {math.MaxInt, ErrTooLarge}, {int(min(maxAllocBytes, math.MaxInt)), ErrTooLarge}, {math.MaxInt / 2, ErrTooLarge}} {
		got, err := PopCountTable(dst, tt.n)
		if !errors.Is(err, tt.wantErr) || !slices.Equal(got, []uint8{7}) {
			t.Errorf("PopCountTable(n=%d) = (%v, %v), want ([7], %v)", tt.n, got, err, tt.wantErr)
		}
	}
}

func TestPopCountTable_oracle(t *testing.T) {
	got, err := PopCountTable(nil, 70_000)
	if err != nil {
		t.Fatal(err)
	}
	for i, c := range got {
		if int(c) != mathbits.OnesCount(uint(i)) {
			t.Fatalf("entry %d = %d, want %d", i, c, mathbits.OnesCount(uint(i)))
		}
	}
}

func TestPopCountTable_noAllocationWithBuffer(t *testing.T) {
	buf := make([]uint8, 0, 1001)
	if allocs := testing.AllocsPerRun(20, func() { buf, _ = PopCountTable(buf[:0], 1000) }); allocs != 0 {
		t.Errorf("PopCountTable with a reused buffer allocates %v times, want 0", allocs)
	}
}

func BenchmarkPopCountTable(b *testing.B) {
	for _, size := range benchSizes {
		buf := make([]uint8, 0, size.n+1)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				buf, _ = PopCountTable(buf[:0], size.n)
			}
			benchSink = uint64(buf[len(buf)-1])
		})
	}
}
