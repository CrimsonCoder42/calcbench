package respiratory

import (
	"errors"
	"testing"
)

// HAND-DERIVED from Gajic 2011 Table 3 (PMID 20802164): sepsis 1 + pneumonia 1.5 + tachypnea 1.5 + SpO₂ 1 +
// diabetes with sepsis −1 = 4. Exact comparison (multiples of 0.5).
func TestLIPSScore_tableExample(t *testing.T) {
	f := LIPSFindings{Sepsis: true, Pneumonia: true, TachypneaRROver30: true, SpO2Below95: true, DiabetesMellitus: true}
	if got, err := LIPSScore(f); err != nil || got != 4 {
		t.Errorf("got (%v, %v), want 4", got, err)
	}
}

// Each row's points alone, and the diabetes rule ("only if sepsis").
func TestLIPSScore_rows(t *testing.T) {
	cases := []struct {
		name string
		f    LIPSFindings
		want float64
	}{
		{"none", LIPSFindings{}, 0},
		{"shock", LIPSFindings{Shock: true}, 2}, {"aspiration", LIPSFindings{Aspiration: true}, 2},
		{"sepsis", LIPSFindings{Sepsis: true}, 1}, {"pneumonia", LIPSFindings{Pneumonia: true}, 1.5},
		{"spine", LIPSFindings{HighRiskSurgery: LIPSSurgeryOrthoSpine}, 1},
		{"acute abdomen", LIPSFindings{HighRiskSurgery: LIPSSurgeryAcuteAbdomen}, 2},
		{"cardiac", LIPSFindings{HighRiskSurgery: LIPSSurgeryCardiac}, 2.5},
		{"aortic", LIPSFindings{HighRiskSurgery: LIPSSurgeryAorticVascular}, 3.5},
		{"aortic emergency", LIPSFindings{HighRiskSurgery: LIPSSurgeryAorticVascular, EmergencySurgery: true}, 5},
		{"TBI", LIPSFindings{TraumaticBrainInj: true}, 2}, {"smoke", LIPSFindings{SmokeInhalation: true}, 2},
		{"drowning", LIPSFindings{NearDrowning: true}, 2}, {"contusion", LIPSFindings{LungContusion: true}, 1.5},
		{"fractures", LIPSFindings{MultipleFractures: true}, 1.5}, {"alcohol", LIPSFindings{AlcoholAbuse: true}, 1},
		{"obesity", LIPSFindings{ObesityBMIOver30: true}, 1}, {"albumin", LIPSFindings{Hypoalbuminemia: true}, 1},
		{"chemo", LIPSFindings{Chemotherapy: true}, 1}, {"FiO2", LIPSFindings{FiO2Over035: true}, 2},
		{"tachypnea", LIPSFindings{TachypneaRROver30: true}, 1.5}, {"SpO2", LIPSFindings{SpO2Below95: true}, 1},
		{"acidosis", LIPSFindings{AcidosisPHBelow735: true}, 1.5},
		{"diabetes without sepsis", LIPSFindings{DiabetesMellitus: true}, 0},
		{"diabetes with sepsis", LIPSFindings{DiabetesMellitus: true, Sepsis: true}, 0},
	}
	for _, c := range cases {
		if got, err := LIPSScore(c.f); err != nil || got != c.want {
			t.Errorf("%s: got (%v, %v), want %v", c.name, got, err, c.want)
		}
	}
}

// parity: ledger lips_score, source case 1: shock + aspiration + FiO₂ > 0.35 + acidosis → 2 + 2 + 2 + 1.5 = 7.5,
// as the source computes. The source takes high-risk trauma as one choice and diabetes as one "diabetes with
// sepsis" input; here the trauma rows add independently and diabetes counts −1 only with Sepsis, as Table 3 reads.
func TestLIPSScore_parity(t *testing.T) {
	f := LIPSFindings{Shock: true, Aspiration: true, FiO2Over035: true, AcidosisPHBelow735: true}
	if got, err := LIPSScore(f); err != nil || got != 7.5 {
		t.Errorf("got (%v, %v), want 7.5", got, err)
	}
}

func TestLIPSScore_errors(t *testing.T) {
	for _, bad := range []LIPSHighRiskSurgery{-1, 5} {
		if _, err := LIPSScore(LIPSFindings{HighRiskSurgery: bad}); !errors.Is(err, ErrOutOfRange) {
			t.Errorf("surgery %d: err = %v, want ErrOutOfRange", bad, err)
		}
	}
}

func TestLIPSScore_noAlloc(t *testing.T) {
	f := LIPSFindings{Sepsis: true, Pneumonia: true}
	if n := testing.AllocsPerRun(100, func() { sinkFloat, _ = LIPSScore(f) }); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkLIPSScore(b *testing.B) {
	f := LIPSFindings{Sepsis: true, Pneumonia: true}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sinkFloat, _ = LIPSScore(f)
	}
}
