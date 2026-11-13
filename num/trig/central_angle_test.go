package trig

import (
	"math"
	"testing"
)

// Worked example, hand-derived from DLMF 4.42.8: two equator points a quarter turn apart are π/2 apart;
// a pole and any equator point are π/2 apart; two equator points 1° apart are π/180 apart.
func TestCentralAngle_example(t *testing.T) {
	tests := []struct{ lat1, lon1, lat2, lon2, want float64 }{
		{0, 0, 0, math.Pi / 2, math.Pi / 2},
		{math.Pi / 2, 0, 0, 1, math.Pi / 2},
		{0, 0, 0, radOf(1), math.Pi / 180},
	}
	for _, tc := range tests {
		if got, err := CentralAngle(tc.lat1, tc.lon1, tc.lat2, tc.lon2); err != nil || !relNear(got, tc.want, 1e-15) {
			t.Errorf("CentralAngle(%v, %v, %v, %v) = (%v, %v), want %v", tc.lat1, tc.lon1, tc.lat2, tc.lon2, got, err, tc.want)
		}
	}
}

func TestCentralAngle_cases(t *testing.T) {
	tests := []struct {
		name                         string
		lat1, lon1, lat2, lon2, want float64
		tol                          float64
	}{
		{"same point", 0.3, 1.2, 0.3, 1.2, 0, 0},
		{"antipodes on the equator", 0, 0, 0, math.Pi, math.Pi, 1e-15},
		{"pole to pole", math.Pi / 2, 0, -math.Pi / 2, 0, math.Pi, 1e-15},
		// Small: two points 1e-9 rad apart along a meridian. The law of cosines would give 0 or 2e-8
		// here; the haversine form keeps the digits.
		{"tiny separation", 0.5, 0.1, 0.5 + 1e-9, 0.1, 1e-9, 1e-6},
		{"longitudes a full turn apart", 0.2, 0.1, 0.2, 0.1 + 2*math.Pi, 0, 1e-15},
		// Δλ = 2e308 is beyond float64, but Δλ/2 = 1e308 is exact, and on the equator σ = 2·arcsin|sin(Δλ/2)|.
		{"huge longitudes do not overflow", 0, -1e308, 0, 1e308, 2 * math.Asin(math.Abs(math.Sin(1e308))), 1e-15},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got, err := CentralAngle(tc.lat1, tc.lon1, tc.lat2, tc.lon2); err != nil || !near(got, tc.want, tc.tol) {
				t.Errorf("got (%v, %v), want %v", got, err, tc.want)
			}
		})
	}
}

func TestCentralAngle_errors(t *testing.T) {
	call := func(lat1, lon1, lat2, lon2 float64) func() error {
		return func() error { _, err := CentralAngle(lat1, lon1, lat2, lon2); return err }
	}
	runErrCases(t, []errCase{
		{"lat1 NaN", call(math.NaN(), 0, 0, 0), ErrNotFinite},
		{"lon1 +Inf", call(0, math.Inf(1), 0, 0), ErrNotFinite},
		{"lat2 -Inf", call(0, 0, math.Inf(-1), 0), ErrNotFinite},
		{"lon2 NaN", call(0, 0, 0, math.NaN()), ErrNotFinite},
		{"lat1 past the pole", call(2, 0, 0, 0), ErrDomain},
		{"lat2 past the south pole", call(0, 0, -1.6, 0), ErrDomain},
	})
}

// Oracle: the spherical law of cosines, cos σ = sin φ1 sin φ2 + cos φ1 cos φ2 cos Δλ (DLMF 4.42.8 with
// colatitudes), on pairs with σ in [0.1, 3.0], where arccos is well conditioned (|d acos| ≤ 1/sin 0.1 ≈
// 10): tolerance 1e-13. Also symmetry, exactly.
func TestCentralAngle_oracle(t *testing.T) {
	rng := newRNG()
	checked := 0
	for range 2000 {
		lat1, lat2 := (rng.Float64()-0.5)*math.Pi, (rng.Float64()-0.5)*math.Pi
		lon1, lon2 := (rng.Float64()-0.5)*4*math.Pi, (rng.Float64()-0.5)*4*math.Pi
		want := math.Acos(math.Sin(lat1)*math.Sin(lat2) + math.Cos(lat1)*math.Cos(lat2)*math.Cos(lon2-lon1))
		got, err := CentralAngle(lat1, lon1, lat2, lon2)
		back, _ := CentralAngle(lat2, lon2, lat1, lon1)
		if err != nil || got < 0 || got > math.Pi || got != back {
			t.Fatalf("CentralAngle(%v, %v, %v, %v) = (%v, %v), reversed %v", lat1, lon1, lat2, lon2, got, err, back)
		}
		if want < 0.1 || want > 3.0 {
			continue
		}
		checked++
		if !near(got, want, 1e-13) {
			t.Fatalf("CentralAngle(%v, %v, %v, %v) = %v, law of cosines %v", lat1, lon1, lat2, lon2, got, want)
		}
	}
	if checked < 1000 {
		t.Fatalf("only %d pairs checked", checked)
	}
}

func TestCentralAngle_noAllocation(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sink, _ = CentralAngle(0.7, -1.3, -0.2, 2.1) }); n != 0 {
		t.Errorf("CentralAngle allocates %v times, want 0", n)
	}
}

func BenchmarkCentralAngle(b *testing.B) {
	b.ReportAllocs()
	lat1, lon1, lat2, lon2 := 0.7, -1.3, -0.2, 2.1
	for b.Loop() {
		sink, _ = CentralAngle(lat1, lon1, lat2, lon2)
	}
}
