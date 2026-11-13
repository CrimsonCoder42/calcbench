package integer

import (
	"errors"
	"math"
	"testing"
)

// Knuth §4.5.2: gcd(u, v)·lcm(u, v) = u·v; with gcd(40902, 24140) = 34, lcm = 40902·24140/34 = 29040420.
func TestLCM_knuthExample(t *testing.T) {
	if got, err := LCM(40902, 24140); err != nil || got != 29040420 {
		t.Errorf("LCM(40902, 24140) = (%d, %v), want 29040420", got, err)
	}
}

func TestLCM_cases(t *testing.T) {
	tests := []struct {
		name       string
		a, b, want uint64
	}{
		{"zero", 0, 5, 0},
		{"both zero", 0, 0, 0},
		{"one", 1, 7, 7},
		{"coprime", 4, 9, 36},
		{"multiple", 6, 3, 6},
		{"max and 1", math.MaxUint64, 1, math.MaxUint64},
		{"max and its divisor", math.MaxUint64, 5, math.MaxUint64},
		{"just fits", 1 << 32, (1 << 32) - 1, (1 << 64) - (1 << 32)},
	}
	for _, tc := range tests {
		if got, err := LCM(tc.a, tc.b); err != nil || got != tc.want {
			t.Errorf("%s: LCM(%d, %d) = (%d, %v), want %d", tc.name, tc.a, tc.b, got, err, tc.want)
		}
	}
}

func TestLCM_overflow(t *testing.T) {
	if _, err := LCM(1<<32+1, 1<<32+3); !errors.Is(err, ErrOverflow) {
		t.Errorf("err = %v, want ErrOverflow", err)
	}
	if _, err := LCM(math.MaxUint64, 2); !errors.Is(err, ErrOverflow) {
		t.Errorf("max·2: err = %v, want ErrOverflow", err)
	}
}

func BenchmarkLCM(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_, _ = LCM(40902, 24140)
	}
}
