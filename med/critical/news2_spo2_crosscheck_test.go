package critical_test

import (
	"testing"

	"github.com/CrimsonCoder42/calcbench/med/critical"
	"github.com/CrimsonCoder42/calcbench/med/respiratory"
)

// The NEWS2 SpO₂ bands live twice: inside MED-CRIT-001 NEWS2 and as MED-RESP-009/010 (D-P14: med packages do not
// import each other in production code). This test keeps the two copies identical over every whole SpO₂ value
// 70–100 % on both scales, on air and on oxygen, so neither can drift (review finding).
func TestNEWS2_spo2MatchesRespiratoryPoints(t *testing.T) {
	for p := 70; p <= 100; p++ {
		for _, scale2 := range []bool{false, true} {
			for _, oxygen := range []bool{false, true} {
				// Every other parameter at its 0-point value (RR 16, SBP 120, pulse 70, alert, 37 °C).
				got, err := critical.NEWS2(16, p, scale2, oxygen, 120, 70, true, 37)
				if err != nil {
					t.Fatalf("NEWS2(spo2 %d, scale2 %v, oxygen %v): %v", p, scale2, oxygen, err)
				}
				var want int
				if scale2 {
					want, err = respiratory.NEWS2SpO2Scale2Points(p, oxygen)
				} else {
					want, err = respiratory.NEWS2SpO2Scale1Points(p)
				}
				if err != nil {
					t.Fatalf("respiratory points(spo2 %d, scale2 %v, oxygen %v): %v", p, scale2, oxygen, err)
				}
				if got.SpO2 != want {
					t.Errorf("spo2 %d scale2 %v oxygen %v: critical %d, respiratory %d", p, scale2, oxygen, got.SpO2, want)
				}
			}
		}
	}
}
