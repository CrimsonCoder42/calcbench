package dp

import (
	"math"
	"testing"
)

// Sizes shared by the package's benchmarks: linear and n log n functions grow over linearSizes,
// quadratic ones (n² cells) over quadraticSizes.
var (
	linearSizes = []struct {
		name string
		n    int
	}{{"1e3", 1_000}, {"1e4", 10_000}, {"1e5", 100_000}}
	quadraticSizes = []struct {
		name string
		n    int
	}{{"1e2", 100}, {"1e3", 1_000}, {"1e4", 10_000}}
)

func TestAddInt64_cases(t *testing.T) {
	tests := []struct {
		a, b int64
		want int64
		ok   bool
	}{
		{1, 2, 3, true},
		{-5, 3, -2, true},
		{math.MaxInt64, 0, math.MaxInt64, true},
		{math.MaxInt64, 1, 0, false},
		{math.MinInt64, -1, 0, false},
		{math.MinInt64, math.MaxInt64, -1, true},
	}
	for _, tt := range tests {
		got, ok := addInt64(tt.a, tt.b)
		if ok != tt.ok || (ok && got != tt.want) {
			t.Errorf("addInt64(%d, %d) = (%d, %v), want (%d, %v)", tt.a, tt.b, got, ok, tt.want, tt.ok)
		}
	}
}

func TestAddUint64_cases(t *testing.T) {
	if s, ok := addUint64(math.MaxUint64-1, 1); !ok || s != math.MaxUint64 {
		t.Errorf("addUint64(max−1, 1) = (%d, %v)", s, ok)
	}
	if _, ok := addUint64(math.MaxUint64, 1); ok {
		t.Error("addUint64(max, 1) reported no overflow")
	}
}

func TestMulPositiveInt64_cases(t *testing.T) {
	if p, ok := mulPositiveInt64(1<<31, 1<<31); !ok || p != 1<<62 {
		t.Errorf("mulPositiveInt64(2³¹, 2³¹) = (%d, %v)", p, ok)
	}
	if _, ok := mulPositiveInt64(1<<32, 1<<31); ok {
		t.Error("mulPositiveInt64(2³², 2³¹) reported no overflow")
	}
}

func TestTableFits_cases(t *testing.T) {
	tests := []struct {
		rows, cols, cell int
		want             bool
	}{
		{0, 5, 8, true},
		{1 << 24, 1, 8, true},
		{1<<24 + 1, 1, 8, false},
		{math.MaxInt, math.MaxInt, 8, false},
		{1 << 40, 1 << 40, 1, false},
	}
	for _, tt := range tests {
		if got := tableFits(tt.rows, tt.cols, tt.cell); got != tt.want {
			t.Errorf("tableFits(%d, %d, %d) = %v, want %v", tt.rows, tt.cols, tt.cell, got, tt.want)
		}
	}
}
