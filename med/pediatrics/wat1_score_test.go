package pediatrics

import (
	"errors"
	"math"
	"testing"
)

// Franck 2008 form (PMC2775493): ten 0/1 items and the 0/1/2 recovery item; recovery edges both sides.
func TestWAT1Score_form(t *testing.T) {
	cases := []struct {
		name string
		f    WAT1Findings
		want int
	}{
		{"none", WAT1Findings{}, 0},
		{"calm 1.9 min", WAT1Findings{MinutesToCalmState: 1.9}, 0},
		{"calm 2 min", WAT1Findings{MinutesToCalmState: 2}, 1},
		{"calm 5 min", WAT1Findings{MinutesToCalmState: 5}, 1},
		{"calm 5.1 min", WAT1Findings{MinutesToCalmState: 5.1}, 2},
		{"stools", WAT1Findings{LooseWateryStools: true}, 1},
		{"temperature", WAT1Findings{TemperatureAbove378C: true}, 1},
		{"doc example", WAT1Findings{LooseWateryStools: true, TremorModerateSevere: true, Sweating: true, StartleToTouch: true, MinutesToCalmState: 3}, 5},
		{"maximum", WAT1Findings{true, true, true, true, true, true, true, true, true, true, 10}, 12},
	}
	for _, c := range cases {
		got, err := WAT1Score(c.f)
		if err != nil || got != c.want {
			t.Errorf("%s: got %d, %v; want %d", c.name, got, err, c.want)
		}
	}
}

// parity: ledger wat1, source case 1: the consumer's worked example "watery stools (1), tremor (1), sweating
// (1), startle to touch (1), and recovery time 2-5" min (1) → 5.
func TestWAT1Score_parity(t *testing.T) {
	f := WAT1Findings{LooseWateryStools: true, TremorModerateSevere: true, Sweating: true, StartleToTouch: true, MinutesToCalmState: 4}
	if got, err := WAT1Score(f); err != nil || got != 5 {
		t.Errorf("got %d, %v; want 5", got, err)
	}
}

func TestWAT1Score_errors(t *testing.T) {
	for _, c := range []struct {
		v    float64
		want error
	}{{math.NaN(), ErrNotFinite}, {math.Inf(1), ErrNotFinite}, {math.Inf(-1), ErrNotFinite}, {-0.1, ErrOutOfRange}} {
		if _, err := WAT1Score(WAT1Findings{MinutesToCalmState: c.v}); !errors.Is(err, c.want) {
			t.Errorf("%v: err = %v, want %v", c.v, err, c.want)
		}
	}
}

func TestWAT1Score_noAlloc(t *testing.T) {
	f := WAT1Findings{Sweating: true, MinutesToCalmState: 3}
	if n := testing.AllocsPerRun(100, func() { sinkInt, _ = WAT1Score(f) }); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkWAT1Score(b *testing.B) {
	f := WAT1Findings{Sweating: true, MinutesToCalmState: 3}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sinkInt, _ = WAT1Score(f)
	}
}
