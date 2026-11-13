package stack

import (
	"math"
	"math/rand/v2"
	"slices"
	"testing"
)

// Hand-derived (catalog), our own readings: 15 18 16 14 17 20 19 19 → 1 4 2 1 1 0 0 0
// (18 waits for 20 four places on; the equal 19s never see a warmer reading).
func TestDailyTemperatures_example(t *testing.T) {
	got := DailyTemperatures(nil, []int{15, 18, 16, 14, 17, 20, 19, 19})
	if want := []int{1, 4, 2, 1, 1, 0, 0, 0}; !slices.Equal(got, want) {
		t.Errorf("DailyTemperatures = %v, want %v", got, want)
	}
}

func TestDailyTemperatures_cases(t *testing.T) {
	tests := []struct {
		name     string
		readings []float64
		want     []int
	}{
		{"empty", nil, []int{}},
		{"one", []float64{20}, []int{0}},
		{"two rising", []float64{20, 21}, []int{1, 0}},
		{"equal never warmer", []float64{20, 20, 20}, []int{0, 0, 0}},
		{"falling", []float64{30, 20, 10}, []int{0, 0, 0}},
		{"rising", []float64{10, 20, 30}, []int{1, 1, 0}},
		{"below zero", []float64{-5.5, -7, 0, -1}, []int{2, 1, 0, 0}},
		{"infinities", []float64{math.Inf(-1), math.Inf(1)}, []int{1, 0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DailyTemperatures(nil, tt.readings); !slices.Equal(got, tt.want) {
				t.Errorf("DailyTemperatures(%v) = %v, want %v", tt.readings, got, tt.want)
			}
		})
	}
}

func TestDailyTemperatures_appendsToDst(t *testing.T) {
	got := DailyTemperatures([]int{7}, []int{1, 2})
	if want := []int{7, 1, 0}; !slices.Equal(got, want) {
		t.Errorf("DailyTemperatures(dst) = %v, want %v", got, want)
	}
}

func TestDailyTemperatures_noAllocWithBuffer(t *testing.T) {
	readings := []int{5, 1, 4, 2, 3}
	buf := make([]int, 0, len(readings))
	if allocs := testing.AllocsPerRun(100, func() { _ = DailyTemperatures(buf[:0], readings) }); allocs != 0 {
		t.Errorf("DailyTemperatures with capacity allocated %v times, want 0", allocs)
	}
}

func TestDailyTemperatures_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 400 {
		readings := randomInts(rng, rng.IntN(120), 1+rng.IntN(40))
		want := make([]int, len(readings))
		for i := range readings {
			for j := i + 1; j < len(readings); j++ {
				if readings[j] > readings[i] {
					want[i] = j - i
					break
				}
			}
		}
		if got := DailyTemperatures(nil, readings); !slices.Equal(got, want) {
			t.Fatalf("trial %d: DailyTemperatures(%v) = %v, want %v", trial, readings, got, want)
		}
	}
}

func BenchmarkDailyTemperatures(b *testing.B) { benchmarkNearest(b, DailyTemperatures[int]) }
