package float

import (
	"errors"
	"math"
	"testing"
)

var sinkF float64

func TestUlp_Values(t *testing.T) {
	// Muller et al. 2018 §2.6.1: ulp(x) = 2^(e−p+1) for x in [2^e, 2^(e+1)), p = 53; IEEE 754 binary64.
	cases := []struct {
		name    string
		x, want float64
	}{
		{"one", 1, 0x1p-52},
		{"minus one", -1, 0x1p-52},
		{"just below two", math.Nextafter(2, 0), 0x1p-52},
		{"two", 2, 0x1p-51},
		{"1e300 binade", 0x1p996, 0x1p944},
		{"zero", 0, 0x1p-1074},
		{"negative zero", math.Copysign(0, -1), 0x1p-1074},
		{"subnormal", 0x1p-1060, 0x1p-1074},
		{"smallest normal", SmallestNormal, 0x1p-1074},
		{"max float", math.MaxFloat64, 0x1p971},
	}
	for _, c := range cases {
		got, err := Ulp(c.x)
		if err != nil || got != c.want {
			t.Errorf("%s: Ulp(%v) = %v, %v; want %v", c.name, c.x, got, err, c.want)
		}
	}
	if MachineEpsilon != 0x1p-52 || SmallestNormal != 0x1p-1022 {
		t.Errorf("constants: %v, %v", MachineEpsilon, SmallestNormal)
	}
}

func TestUlp_Errors(t *testing.T) {
	for _, x := range []float64{math.NaN(), math.Inf(1), math.Inf(-1)} {
		if _, err := Ulp(x); !errors.Is(err, ErrNotFinite) {
			t.Errorf("Ulp(%v): err = %v; want ErrNotFinite", x, err)
		}
	}
}

func BenchmarkUlp(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkF, _ = Ulp(123.456)
	}
}
