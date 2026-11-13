package infectious

import (
	"errors"
	"math"
	"testing"
)

var sinkFloat float64

// Tolerance 1e-15 relative: one product of exact small integers and one division.
func closeRel(got, want float64) bool { return math.Abs(got-want) <= 1e-15*math.Abs(want) }

// Vertesi 2003, Examples 1 and 2 and Fig. 1.
func TestHIVRASPTotalRisk_examples(t *testing.T) {
	tests := []struct {
		name string
		a    RASPSource
		b    RASPInoculum
		c    RASPMethod
		d    RASPVolume
		want float64
	}{
		// Example 1: A = 10, B = 1, C = 100, small-bore needle multiplier 3 → 3/1000. The text's C value (100)
		// is the "Deep transcutaneous with visible bleeding" row.
		{"example 1", RASPSourceAsymptomatic, RASPInoculumFreshBlood, RASPMethodDeepWithBleeding, RASPVolumeSmall, 3.0 / 1000},
		// Example 2: A = 1000, B = 100, C = 200, Z = 20,000,000; large-bore modifier 5 → 1 in 4 million.
		{"example 2", RASPSourceUnknownLowRisk, RASPInoculumDriedOldBlood, RASPMethodSuperficialNoBleeding, RASPVolumeModerate, 1.0 / 4e6},
		{"maximum: all 1s with massive volume", RASPSourceAcuteAIDS, RASPInoculumFreshBlood, RASPMethodIntravenous, RASPVolumeMassive, 100},
		{"minimum", RASPSourceUnknownLowRisk, RASPInoculumLowRiskSecretions, RASPMethodIntactSkin, RASPVolumeTrace, 1e-9},
		{"mucosal, measurable", RASPSourceUnknownHighRisk, RASPInoculumBodyFluidsAtRisk, RASPMethodMucosal, RASPVolumeMeasurable, 10.0 / 500000},
		{"deep intramuscular, trace", RASPSourceAsymptomatic, RASPInoculumFreshBlood, RASPMethodDeepIntramuscular, RASPVolumeTrace, 1.0 / 100},
	}
	for _, tc := range tests {
		got, err := HIVRASPTotalRisk(tc.a, tc.b, tc.c, tc.d)
		if err != nil || !closeRel(got, tc.want) {
			t.Errorf("%s: got (%v, %v), want %v", tc.name, got, err, tc.want)
		}
	}
}

// parity: ledger hiv_rasp. The source entry has no _test.go and its JSON no example; the consumer computes
// volume / (source × inoculum × method) from the same four tables, so case 1 repeats Example 1.
func TestHIVRASPTotalRisk_parity(t *testing.T) {
	got, err := HIVRASPTotalRisk(RASPSourceAsymptomatic, RASPInoculumFreshBlood, RASPMethodDeepWithBleeding, RASPVolumeSmall)
	if err != nil || !closeRel(got, 0.003) {
		t.Errorf("source case 1: got (%v, %v), want 0.003", got, err)
	}
}

func TestHIVRASPTotalRisk_errors(t *testing.T) {
	ok := func() (RASPSource, RASPInoculum, RASPMethod, RASPVolume) {
		return RASPSourceAsymptomatic, RASPInoculumFreshBlood, RASPMethodMucosal, RASPVolumeSmall
	}
	a, b, c, d := ok()
	for _, bad := range []RASPSource{0, 2, -1, 10000} {
		if _, err := HIVRASPTotalRisk(bad, b, c, d); !errors.Is(err, ErrOutOfRange) {
			t.Errorf("source %d: error = %v, want ErrOutOfRange", bad, err)
		}
	}
	for _, bad := range []RASPInoculum{0, 5} {
		if _, err := HIVRASPTotalRisk(a, bad, c, d); !errors.Is(err, ErrOutOfRange) {
			t.Errorf("inoculum %d: error = %v, want ErrOutOfRange", bad, err)
		}
	}
	for _, bad := range []RASPMethod{0, 300} {
		if _, err := HIVRASPTotalRisk(a, b, bad, d); !errors.Is(err, ErrOutOfRange) {
			t.Errorf("method %d: error = %v, want ErrOutOfRange", bad, err)
		}
	}
	for _, bad := range []RASPVolume{0, 2, 1000} {
		if _, err := HIVRASPTotalRisk(a, b, c, bad); !errors.Is(err, ErrOutOfRange) {
			t.Errorf("volume %d: error = %v, want ErrOutOfRange", bad, err)
		}
	}
}

func TestHIVRASPTotalRisk_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() {
		sinkFloat, _ = HIVRASPTotalRisk(RASPSourceAsymptomatic, RASPInoculumFreshBlood, RASPMethodDeepWithBleeding, RASPVolumeSmall)
	}); n != 0 {
		t.Errorf("allocates %v times, want 0", n)
	}
}

func BenchmarkHIVRASPTotalRisk(b *testing.B) {
	b.ReportAllocs()
	a, i, m, v := RASPSourceAsymptomatic, RASPInoculumFreshBlood, RASPMethodDeepWithBleeding, RASPVolumeSmall
	for b.Loop() {
		sinkFloat, _ = HIVRASPTotalRisk(a, i, m, v)
	}
}
