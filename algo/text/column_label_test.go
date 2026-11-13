package text

import (
	"errors"
	"math"
	"testing"
)

// Foster, Mathematics Magazine 21(1), 1947: in base k without a zero digit, digits run 1…k. With k = 26
// and A…Z for 1…26, hand-derived: 27 = 1·26 + 1 → AA, 702 = 26·26 + 26 → ZZ, 703 = 1·26² + 1·26 + 1 → AAA,
// 16384 = 24·26² + 6·26 + 4 → XFD.
func TestColumnLabel_bijectiveExample(t *testing.T) {
	tests := []struct {
		number int
		want   string
	}{{1, "A"}, {2, "B"}, {26, "Z"}, {27, "AA"}, {28, "AB"}, {52, "AZ"}, {53, "BA"}, {702, "ZZ"}, {703, "AAA"}, {16384, "XFD"}}
	for _, tt := range tests {
		if got, err := ColumnLabel(tt.number); got != tt.want || err != nil {
			t.Errorf("ColumnLabel(%d) = %q, %v; want %q, nil", tt.number, got, err, tt.want)
		}
	}
}

func TestColumnLabel_maximum(t *testing.T) {
	label, err := ColumnLabel(math.MaxInt)
	if err != nil || len(label) == 0 || len(label) > 14 {
		t.Fatalf("ColumnLabel(MaxInt) = %q, %v; want 1–14 letters", label, err)
	}
	if n, err := ColumnNumber(label); n != math.MaxInt || err != nil {
		t.Errorf("ColumnNumber(ColumnLabel(MaxInt)) = %d, %v; want MaxInt", n, err)
	}
}

func TestColumnLabel_errors(t *testing.T) {
	for _, n := range []int{0, -1, math.MinInt} {
		if got, err := ColumnLabel(n); got != "" || !errors.Is(err, ErrColumnNumber) {
			t.Errorf("ColumnLabel(%d) = %q, %v; want \"\", ErrColumnNumber", n, got, err)
		}
	}
}

func TestColumnLabel_oneAllocation(t *testing.T) {
	if allocs := testing.AllocsPerRun(100, func() { s, _ := ColumnLabel(16384); benchSink = len(s) }); allocs > 1 {
		t.Errorf("ColumnLabel allocates %v times, want at most 1", allocs)
	}
}

// Round trip against ColumnNumber, and labels in shortlex order: each label is the successor of the last.
func TestColumnLabel_oracle(t *testing.T) {
	prev := ""
	for n := 1; n <= 20000; n++ {
		label, err := ColumnLabel(n)
		if err != nil {
			t.Fatal(err)
		}
		if back, err := ColumnNumber(label); back != n || err != nil {
			t.Fatalf("ColumnNumber(ColumnLabel(%d) = %q) = %d, %v", n, label, back, err)
		}
		if len(label) < len(prev) || (len(label) == len(prev) && label <= prev) {
			t.Fatalf("ColumnLabel(%d) = %q does not follow %q in shortlex order", n, label, prev)
		}
		prev = label
	}
}

func BenchmarkColumnLabel(b *testing.B) {
	b.Run("scalar", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			s, _ := ColumnLabel(math.MaxInt)
			benchSink = len(s)
		}
	})
}
