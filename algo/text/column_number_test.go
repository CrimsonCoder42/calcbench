package text

import (
	"errors"
	"math"
	"testing"
)

// Foster 1947 (bijective base k), k = 26, A…Z = 1…26; hand-derived values as in column_label_test.go.
func TestColumnNumber_bijectiveExample(t *testing.T) {
	tests := []struct {
		label string
		want  int
	}{{"A", 1}, {"Z", 26}, {"AA", 27}, {"AZ", 52}, {"BA", 53}, {"ZZ", 702}, {"AAA", 703}, {"XFD", 16384}}
	for _, tt := range tests {
		if got, err := ColumnNumber(tt.label); got != tt.want || err != nil {
			t.Errorf("ColumnNumber(%q) = %d, %v; want %d, nil", tt.label, got, err, tt.want)
		}
	}
}

func TestColumnNumber_errors(t *testing.T) {
	maxLabel, _ := ColumnLabel(math.MaxInt)
	tests := []struct {
		name  string
		label string
		want  error
	}{
		{"empty", "", ErrColumnLabel},
		{"lower case", "ab", ErrColumnLabel},
		{"digit", "A1", ErrColumnLabel},
		{"space", " A", ErrColumnLabel},
		{"at sign (one below A)", "@", ErrColumnLabel},
		{"bracket (one above Z)", "[", ErrColumnLabel},
		{"non-ASCII", "É", ErrColumnLabel},
		{"one past MaxInt", maxLabel[:len(maxLabel)-1] + string(maxLabel[len(maxLabel)-1]+1), ErrOverflow},
		{"far too long", "ZZZZZZZZZZZZZZZZ", ErrOverflow},
	}
	if maxLabel[len(maxLabel)-1] == 'Z' {
		t.Fatalf("test assumption broken: ColumnLabel(MaxInt) = %q ends in Z", maxLabel)
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := ColumnNumber(tt.label); got != 0 || !errors.Is(err, tt.want) {
				t.Errorf("ColumnNumber(%q) = %d, %v; want 0, %v", tt.label, got, err, tt.want)
			}
		})
	}
}

func TestColumnNumber_noAllocation(t *testing.T) {
	if allocs := testing.AllocsPerRun(100, func() { benchSink, _ = ColumnNumber("XFD") }); allocs != 0 {
		t.Errorf("ColumnNumber allocates %v times, want 0", allocs)
	}
}

// The round-trip oracle against ColumnLabel is TestColumnLabel_oracle (20 000 numbers).

func BenchmarkColumnNumber(b *testing.B) {
	maxLabel, _ := ColumnLabel(math.MaxInt)
	b.Run("scalar", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			benchSink, _ = ColumnNumber(maxLabel)
		}
	})
}
