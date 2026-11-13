package respiratory

import (
	"errors"
	"testing"
)

// RCP NEWS2 (2017) Chart 1, row "SpO2 Scale 1 (%)": ≤91 3, 92–93 2, 94–95 1, ≥96 0 — both sides of every edge.
func TestNEWS2SpO2Scale1Points_chart(t *testing.T) {
	cases := []struct{ spo2, want int }{
		{0, 3}, {91, 3}, {92, 2}, {93, 2}, {94, 1}, {95, 1}, {96, 0}, {100, 0},
	}
	for _, c := range cases {
		if got, err := NEWS2SpO2Scale1Points(c.spo2); err != nil || got != c.want {
			t.Errorf("%d: got (%d, %v), want %d", c.spo2, got, err, c.want)
		}
	}
}

// RCP NEWS2 (2017) Chart 1, row "SpO2 Scale 2 (%)": ≤83 3, 84–85 2, 86–87 1, 88–92 0, ≥93 on air 0,
// 93–94 on oxygen 1, 95–96 on oxygen 2, ≥97 on oxygen 3 — both sides of every edge, on air and on oxygen.
func TestNEWS2SpO2Scale2Points_chart(t *testing.T) {
	cases := []struct {
		spo2     int
		onOxygen bool
		want     int
	}{
		{0, false, 3}, {83, false, 3}, {83, true, 3}, {84, true, 2}, {85, false, 2}, {86, true, 1},
		{87, false, 1}, {88, true, 0}, {92, true, 0}, {92, false, 0},
		{93, false, 0}, {100, false, 0},
		{93, true, 1}, {94, true, 1}, {95, true, 2}, {96, true, 2}, {97, true, 3}, {100, true, 3},
	}
	for _, c := range cases {
		if got, err := NEWS2SpO2Scale2Points(c.spo2, c.onOxygen); err != nil || got != c.want {
			t.Errorf("%d (oxygen %v): got (%d, %v), want %d", c.spo2, c.onOxygen, got, err, c.want)
		}
	}
}

// parity: ledger spo2_classify, source cases 1–6. The source's condition 0 (standard) and 2 (neonatal) use
// Scale 1; condition 1 (COPD) uses Scale 2 but ignores air vs oxygen. The source's case 6 (99 %, COPD, 3 points)
// agrees with the chart only on oxygen; on air the chart gives 0 — both asserted, ledger note.
func TestNEWS2SpO2Points_parity(t *testing.T) {
	scale1 := []struct{ spo2, want int }{
		{97, 0}, // source case 1
		{93, 2}, // source case 2
		{95, 1}, // source case 5 (neonatal condition; Scale 1 as the source computes it)
	}
	for _, c := range scale1 {
		if got, err := NEWS2SpO2Scale1Points(c.spo2); err != nil || got != c.want {
			t.Errorf("scale 1 %d: got (%d, %v), want %d", c.spo2, got, err, c.want)
		}
	}
	scale2 := []struct {
		spo2     int
		onOxygen bool
		want     int
	}{
		{90, false, 0}, // source case 3
		{84, false, 2}, // source case 4
		{99, true, 3},  // source case 6, on oxygen: the chart and the source agree
		{99, false, 0}, // source case 6, on air: the chart gives 0 (source gives 3)
	}
	for _, c := range scale2 {
		if got, err := NEWS2SpO2Scale2Points(c.spo2, c.onOxygen); err != nil || got != c.want {
			t.Errorf("scale 2 %d (oxygen %v): got (%d, %v), want %d", c.spo2, c.onOxygen, got, err, c.want)
		}
	}
}

func TestNEWS2SpO2Points_errors(t *testing.T) {
	for _, bad := range []int{-1, 101} {
		if _, err := NEWS2SpO2Scale1Points(bad); !errors.Is(err, ErrOutOfRange) {
			t.Errorf("scale 1 %d: err = %v, want ErrOutOfRange", bad, err)
		}
		if _, err := NEWS2SpO2Scale2Points(bad, true); !errors.Is(err, ErrOutOfRange) {
			t.Errorf("scale 2 %d: err = %v, want ErrOutOfRange", bad, err)
		}
	}
}

func TestNEWS2SpO2Points_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkInt, _ = NEWS2SpO2Scale1Points(94) }); n != 0 {
		t.Errorf("scale 1 allocs = %v, want 0", n)
	}
	if n := testing.AllocsPerRun(100, func() { sinkInt, _ = NEWS2SpO2Scale2Points(95, true) }); n != 0 {
		t.Errorf("scale 2 allocs = %v, want 0", n)
	}
}

func BenchmarkNEWS2SpO2Scale1Points(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkInt, _ = NEWS2SpO2Scale1Points(94)
	}
}

func BenchmarkNEWS2SpO2Scale2Points(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkInt, _ = NEWS2SpO2Scale2Points(95, true)
	}
}
