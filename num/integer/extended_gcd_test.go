package integer

import (
	"errors"
	"math"
	"math/big"
	"math/rand/v2"
	"testing"
)

// Knuth §4.5.2, Algorithm X: for u = 40902, v = 24140 it ends with gcd 34 = 40902·337 − 24140·571.
func TestExtendedGCD_knuthExample(t *testing.T) {
	g, x, y, err := ExtendedGCD(40902, 24140)
	if err != nil || g != 34 || x != 337 || y != -571 {
		t.Errorf("ExtendedGCD = (%d, %d, %d, %v), want (34, 337, -571)", g, x, y, err)
	}
}

func TestExtendedGCD_cases(t *testing.T) {
	tests := []struct {
		name          string
		a, b, g, x, y int64
	}{
		{"both zero", 0, 0, 0, 0, 0},
		{"b zero", 7, 0, 7, 1, 0},
		{"negative a, b zero", -7, 0, 7, -1, 0},
		{"a zero", 0, -5, 5, 0, -1},
		{"coprime", 3, 7, 1, -2, 1},
		{"equal", 6, 6, 6, 0, 1},
	}
	for _, tc := range tests {
		g, x, y, err := ExtendedGCD(tc.a, tc.b)
		if err != nil || g != tc.g || x != tc.x || y != tc.y {
			t.Errorf("%s: ExtendedGCD(%d, %d) = (%d, %d, %d, %v), want (%d, %d, %d)",
				tc.name, tc.a, tc.b, g, x, y, err, tc.g, tc.x, tc.y)
		}
	}
}

// Bézout identity checked in big.Int on extremes and 500 random signed pairs, plus the gcd against GCD.
func TestExtendedGCD_identity(t *testing.T) {
	check := func(a, b int64) {
		t.Helper()
		g, x, y, err := ExtendedGCD(a, b)
		if err != nil {
			t.Fatalf("ExtendedGCD(%d, %d): %v", a, b, err)
		}
		lhs := new(big.Int).Mul(big.NewInt(a), big.NewInt(x))
		lhs.Add(lhs, new(big.Int).Mul(big.NewInt(b), big.NewInt(y)))
		if lhs.Cmp(big.NewInt(g)) != 0 || g < 0 {
			t.Fatalf("ExtendedGCD(%d, %d) = (%d, %d, %d): a·x + b·y = %v", a, b, g, x, y, lhs)
		}
		ua, ub := a, b
		if ua < 0 {
			ua = -ua
		}
		if ub < 0 {
			ub = -ub
		}
		if uint64(g) != GCD(uint64(ua), uint64(ub)) {
			t.Fatalf("ExtendedGCD(%d, %d) g = %d, GCD = %d", a, b, g, GCD(uint64(ua), uint64(ub)))
		}
	}
	check(math.MaxInt64, math.MaxInt64-1)
	check(-math.MaxInt64, math.MaxInt64)
	check(7540113804746346429, 4660046610375530309)
	r := rand.New(rand.NewPCG(1, 2))
	for range 500 {
		check(r.Int64()/2-r.Int64()/2, r.Int64()/2-r.Int64()/2)
	}
}

func TestExtendedGCD_overflow(t *testing.T) {
	for _, p := range [][2]int64{{math.MinInt64, 3}, {3, math.MinInt64}} {
		if _, _, _, err := ExtendedGCD(p[0], p[1]); !errors.Is(err, ErrOverflow) {
			t.Errorf("ExtendedGCD(%d, %d): err = %v, want ErrOverflow", p[0], p[1], err)
		}
	}
}

func BenchmarkExtendedGCD(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_, _, _, _ = ExtendedGCD(7540113804746346429, 4660046610375530309)
	}
}
