package cardio

import "testing"

// Antzelevitch 2016, Table 2: ECG item + highest of II + highest of III + IV. Halves are exact in float64.
func TestShanghaiBrugadaScore_items(t *testing.T) {
	tests := []struct {
		name string
		ecg  ShanghaiECG
		h    ShanghaiHistory
		want float64
	}{
		{"spontaneous only", ShanghaiECGSpontaneousType1, ShanghaiHistory{}, 3.5},
		{"fever only", ShanghaiECGFeverInducedType1, ShanghaiHistory{}, 3},
		{"drug only", ShanghaiECGDrugConvertedType2Or3, ShanghaiHistory{}, 2},
		{"example", ShanghaiECGSpontaneousType1, ShanghaiHistory{SuspectedArrhythmicSyncope: true, RelativeDefiniteBrS: true}, 7.5},
		{"highest of II only", ShanghaiECGDrugConvertedType2Or3, ShanghaiHistory{CardiacArrestOrVFPolyVT: true, SyncopeUnclearMechanism: true, AFFlutterUnder30: true}, 5},
		{"highest of III only", ShanghaiECGDrugConvertedType2Or3, ShanghaiHistory{RelativeSuspiciousSCD: true, RelativeUnexplainedSCDU45: true}, 3},
		{"small items", ShanghaiECGDrugConvertedType2Or3, ShanghaiHistory{AFFlutterUnder30: true, RelativeUnexplainedSCDU45: true, ProbablePathogenicMutation: true}, 3.5},
		{"maximum", ShanghaiECGSpontaneousType1, ShanghaiHistory{CardiacArrestOrVFPolyVT: true, RelativeDefiniteBrS: true, ProbablePathogenicMutation: true}, 9},
	}
	for _, tc := range tests {
		got, err := ShanghaiBrugadaScore(tc.ecg, tc.h)
		if err != nil || got != tc.want {
			t.Errorf("%s: got (%v, %v), want %v", tc.name, got, err, tc.want)
		}
	}
	for _, e := range []ShanghaiECG{0, 4} {
		_, err := ShanghaiBrugadaScore(e, ShanghaiHistory{})
		checkErr(t, "ecg", err, ErrOutOfRange)
	}
}

func TestShanghaiBrugadaScore_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkFloat, _ = ShanghaiBrugadaScore(ShanghaiECGSpontaneousType1, ShanghaiHistory{}) }); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkShanghaiBrugadaScore(b *testing.B) {
	b.ReportAllocs()
	h := ShanghaiHistory{SuspectedArrhythmicSyncope: true}
	for b.Loop() {
		sinkFloat, _ = ShanghaiBrugadaScore(ShanghaiECGSpontaneousType1, h)
	}
}
