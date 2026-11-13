package intmath

import (
	"errors"
	"math"
	"math/rand/v2"
	"testing"
)

// Go's specification, "Arithmetic operators": 5/3 = 1 r 2, −5/3 = −1 r −2, 5/−3 = −1 r 2, −5/−3 = 1 r −2.
func TestLongDivision_sourceExample(t *testing.T) {
	tests := []struct{ x, y, q, r int64 }{{5, 3, 1, 2}, {-5, 3, -1, -2}, {5, -3, -1, 2}, {-5, -3, 1, -2}}
	for _, tt := range tests {
		q, r, err := LongDivision(tt.x, tt.y)
		if err != nil || q != tt.q || r != tt.r {
			t.Errorf("LongDivision(%d, %d) = (%d, %d, %v), want (%d, %d)", tt.x, tt.y, q, r, err, tt.q, tt.r)
		}
	}
}

func TestLongDivision_cases(t *testing.T) {
	tests := []struct{ x, y, q, r int64 }{
		{0, 7, 0, 0}, {7, 1, 7, 0}, {7, -1, -7, 0}, {6, 7, 0, 6}, {100, 10, 10, 0},
		{math.MaxInt64, 1, math.MaxInt64, 0}, {math.MinInt64, 1, math.MinInt64, 0},
		{math.MinInt64, 2, math.MinInt64 / 2, 0}, {math.MinInt64, math.MinInt64, 1, 0},
		{math.MaxInt64, math.MinInt64, 0, math.MaxInt64}, {math.MinInt64, math.MaxInt64, -1, -1},
		{math.MaxInt64, math.MaxInt64, 1, 0}, {math.MaxInt64, 2, math.MaxInt64 / 2, 1},
	}
	for _, tt := range tests {
		q, r, err := LongDivision(tt.x, tt.y)
		if err != nil || q != tt.q || r != tt.r {
			t.Errorf("LongDivision(%d, %d) = (%d, %d, %v), want (%d, %d)", tt.x, tt.y, q, r, err, tt.q, tt.r)
		}
	}
	// int32 and int8 semantics, including the minimum.
	if q, r, err := LongDivision[int32](math.MinInt32, 3); err != nil || q != math.MinInt32/3 || r != math.MinInt32%3 {
		t.Errorf("int32: (%d, %d, %v)", q, r, err)
	}
	if q, r, err := LongDivision[int8](-128, 1); err != nil || q != -128 || r != 0 {
		t.Errorf("int8 −128/1: (%d, %d, %v)", q, r, err)
	}
}

func TestLongDivision_errors(t *testing.T) {
	if _, _, err := LongDivision[int64](5, 0); !errors.Is(err, ErrDivideByZero) {
		t.Errorf("5/0 error = %v, want ErrDivideByZero", err)
	}
	if _, _, err := LongDivision[int64](math.MinInt64, -1); !errors.Is(err, ErrOverflow) {
		t.Errorf("MinInt64/−1 error = %v, want ErrOverflow", err)
	}
	if _, _, err := LongDivision[int32](math.MinInt32, -1); !errors.Is(err, ErrOverflow) {
		t.Errorf("MinInt32/−1 error = %v, want ErrOverflow", err)
	}
	if _, _, err := LongDivision[int8](-128, -1); !errors.Is(err, ErrOverflow) {
		t.Errorf("MinInt8/−1 error = %v, want ErrOverflow", err)
	}
}

// Oracle: Go's / and % operators.
func TestLongDivision_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for range 1000 {
		x := int64(rng.Uint64()) >> rng.IntN(64)
		y := int64(rng.Uint64()) >> rng.IntN(64)
		if y == 0 || (x == math.MinInt64 && y == -1) {
			continue
		}
		q, r, err := LongDivision(x, y)
		if err != nil || q != x/y || r != x%y {
			t.Fatalf("LongDivision(%d, %d) = (%d, %d, %v), want (%d, %d)", x, y, q, r, err, x/y, x%y)
		}
	}
	for x := -128; x <= 127; x++ {
		for y := -128; y <= 127; y++ {
			if y == 0 || (x == -128 && y == -1) {
				continue
			}
			q, r, err := LongDivision(int8(x), int8(y))
			if err != nil || q != int8(x)/int8(y) || r != int8(x)%int8(y) {
				t.Fatalf("int8 LongDivision(%d, %d) = (%d, %d, %v)", x, y, q, r, err)
			}
		}
	}
}

func BenchmarkLongDivision(b *testing.B) {
	b.ReportAllocs()
	var s int64
	for i := int64(1); b.Loop(); i++ {
		q, _, _ := LongDivision(math.MaxInt64-i, i&0xFFFF+1)
		s += q
	}
	benchSink = uint64(s)
}
