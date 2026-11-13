package invest

import "testing"

// OpenStax Principles of Finance §15.1: Target $128.74 → $176.53 (capital gain $47.79); 3M $175.00 →
// $169.72 → −3.02 % (displayed to 2 decimals of a percent, ±5e-5).
func TestCapitalGainsYield_openStax(t *testing.T) {
	got, err := CapitalGainsYield(12874, 17653)
	if err != nil || got != 4779.0/12874.0 {
		t.Errorf("Target: %v, %v; want 4779/12874", got, err)
	}
	got, err = CapitalGainsYield(17500, 16972)
	if err != nil || !near(got, -0.0302, 5e-5) {
		t.Errorf("3M: %v, %v; want −0.0302", got, err)
	}
}

func TestCapitalGainsYield_cases(t *testing.T) {
	for _, c := range []struct {
		name       string
		begin, end int64
		want       float64
	}{{"flat", 500, 500, 0}, {"to zero", 500, 0, -1}, {"double", 500, 1000, 1}, {"max", 1, maxCents, float64(maxCents - 1)}} {
		got, err := CapitalGainsYield(c.begin, c.end)
		if err != nil || got != c.want {
			t.Errorf("%s: %v, %v; want %v", c.name, got, err, c.want)
		}
	}
}

func TestCapitalGainsYield_errors(t *testing.T) {
	_, err := CapitalGainsYield(0, 10)
	wantErr(t, "begin 0", err, ErrZeroDenominator)
	_, err = CapitalGainsYield(-1, 10)
	wantErr(t, "begin < 0", err, ErrNegative)
	_, err = CapitalGainsYield(10, -1)
	wantErr(t, "end < 0", err, ErrNegative)
	_, err = CapitalGainsYield(10, maxCents+1)
	wantErr(t, "end > 2^53", err, ErrOutOfRange)
}

func TestCapitalGainsYield_noAlloc(t *testing.T) {
	if a := testing.AllocsPerRun(100, func() { sinkFloat, sinkErr = CapitalGainsYield(12874, 17653) }); a != 0 {
		t.Errorf("allocs = %v, want 0", a)
	}
}

func BenchmarkCapitalGainsYield(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, sinkErr = CapitalGainsYield(12874, 17653)
	}
}
