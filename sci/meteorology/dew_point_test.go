package meteorology

import (
	"errors"
	"math"
	"testing"
)

// Hand computation with the Alduchov-Eskridge Magnus coefficients (a = 17.625, b = 243.04 °C):
// t = 21.1111 °C, RH 65 %: γ = ln 0.65 + 17.625·21.1111/264.1511 = −0.430783 + 1.408606 = 0.977823,
// Td = 243.04·0.977823/16.647177 = 14.2756 °C. RH = 100 % gives Td = t exactly (γ = a·t/(b + t)).
// Tolerance 1e-3 °C (hand arithmetic).
func TestDewPointC_cases(t *testing.T) {
	got, err := DewPointC(70.0*5/9-32.0*5/9, 65)
	if err != nil || math.Abs(got-14.2756) > 1e-3 {
		t.Errorf("70 °F 65 %%: %v %v", got, err)
	}
	if got, err := DewPointC(25, 100); err != nil || math.Abs(got-25) > 1e-12 {
		t.Errorf("saturated: %v %v", got, err)
	}
	cases := []struct {
		t, rh float64
		want  error
	}{
		{20, 0, ErrOutOfRange},
		{20, 100.1, ErrOutOfRange},
		{-41, 50, ErrOutOfRange},
		{51, 50, ErrOutOfRange},
		{math.NaN(), 50, ErrNotFinite},
	}
	for _, c := range cases {
		if _, err := DewPointC(c.t, c.rh); !errors.Is(err, c.want) {
			t.Errorf("(%v, %v): %v, want %v", c.t, c.rh, err, c.want)
		}
	}
}

// Round trip: RH → Td → RH returns the RH within 1e-9 %, over a grid of temperatures and humidities.
func TestRelativeHumidityFromDewPoint_roundTrip(t *testing.T) {
	for tc := -40.0; tc <= 50; tc += 5 {
		for rh := 5.0; rh <= 100; rh += 5 {
			td, err := DewPointC(tc, rh)
			if err != nil {
				t.Fatalf("DewPointC(%v, %v): %v", tc, rh, err)
			}
			if td < -40 {
				continue
			}
			back, err := RelativeHumidityFromDewPoint(tc, td)
			if err != nil || math.Abs(back-rh) > 1e-9 {
				t.Fatalf("t %v rh %v: td %v back %v %v", tc, rh, td, back, err)
			}
		}
	}
}

func TestRelativeHumidityFromDewPoint_cases(t *testing.T) {
	// 85 °F = 29.4444 °C, 70 °F = 21.1111 °C → 60.8997 % (hand computation, same coefficients).
	got, err := RelativeHumidityFromDewPoint((85.0-32)*5/9, (70.0-32)*5/9)
	if err != nil || math.Abs(got-60.8997) > 1e-3 {
		t.Errorf("got %v %v", got, err)
	}
	if _, err := RelativeHumidityFromDewPoint(20, 21); !errors.Is(err, ErrOutOfRange) {
		t.Errorf("dew point above t: %v", err)
	}
	if _, err := RelativeHumidityFromDewPoint(20, math.Inf(-1)); !errors.Is(err, ErrNotFinite) {
		t.Errorf("Inf: %v", err)
	}
}

func BenchmarkDewPointC(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = DewPointC(21.1111, 65)
	}
}

func BenchmarkRelativeHumidityFromDewPoint(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = RelativeHumidityFromDewPoint(29.4444, 21.1111)
	}
}
