package endocrine

import "testing"

// parity: ledger bs_gestational_ogtt_one_step, source cases 1–3 (each value against its own threshold,
// fasting ≥ 92, 1 h ≥ 180, 2 h ≥ 153 mg/dL; any one met → criteria met). Reused by MED-ENDO-005; IADPSG
// 2010 Table 1 agrees with the source on every threshold.
func TestGDMThresholdIADPSGMgDL_parityOneStep(t *testing.T) {
	cases := []struct {
		f, h1, h2 float64
		want      bool
	}{
		{91.9, 179.9, 152.9, false}, // every value "Below threshold"
		{92, 100, 100, true},        // case 1: fasting "Threshold met"
		{80, 180, 100, true},        // case 2: 1 h "Threshold met"
		{80, 100, 153, true},        // case 3: 2 h "Threshold met"
		{92, 180, 153, true},
	}
	for _, c := range cases {
		got, err := GDMThresholdIADPSGMgDL(c.f, c.h1, c.h2)
		if err != nil || got != c.want {
			t.Errorf("(%v, %v, %v) = %v, %v; want %v", c.f, c.h1, c.h2, got, err, c.want)
		}
	}
}
