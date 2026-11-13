package interest

import "testing"

func TestRuleOf72_SourceExample(t *testing.T) {
	// Principles of Finance (OpenStax 2022) §7.4, Eq. 7.37: at 9 % money doubles in 72/9 = 8 years;
	// at 6 %, 12 years; an economy growing 4 % doubles in 18 years.
	got, err := RuleOf72(0.09)
	wantFloat(t, "§7.4 9 %", got, err, 8, 1e-12)
	got, err = RuleOf72(0.06)
	wantFloat(t, "§7.4 6 %", got, err, 12, 1e-12)
	got, err = RuleOf72(0.04)
	wantFloat(t, "§7.4 4 %", got, err, 18, 1e-12)
	// DeFusco et al., Quantitative Investment Analysis 3rd ed. ch. 1 footnote 9: 72/7 = 10.3.
	got, err = RuleOf72(0.07)
	wantFloat(t, "QIA fn. 9", got, err, 10.3, 0.05)
}

func TestRuleOf72_Table(t *testing.T) {
	got, err := RuleOf72(0.72)
	wantFloat(t, "72 %", got, err, 1, 1e-15)
	got, err = RuleOf72(1e-6)
	wantRel(t, "tiny rate", got, err, 720000, 1e-15)
}

func TestRuleOf72_Errors(t *testing.T) {
	for _, r := range badRates {
		_, err := RuleOf72(r)
		wantErr(t, "bad rate", err, ErrRate)
	}
	_, err := RuleOf72(0)
	wantErr(t, "zero", err, ErrDomain)
	_, err = RuleOf72(-0.05)
	wantErr(t, "negative", err, ErrDomain)
	_, err = RuleOf72(1e-310)
	wantErr(t, "beyond float64", err, ErrOverflow)
}

func TestRuleOf72_Oracle(t *testing.T) {
	// Against the exact DoublingTime the rule is within 2 % for rates from 4 % to 11 % (it is exact
	// near 7.85 %; QIA ch. 1 fn. 9 calls it an approximation).
	for r := 0.04; r < 0.115; r += 0.01 {
		rule, err := RuleOf72(r)
		exact, err2 := DoublingTime(r)
		if err2 != nil {
			t.Fatal(err2)
		}
		wantRel(t, "rule vs exact", rule, err, exact, 0.02)
	}
}

func TestRuleOf72_NoAllocs(t *testing.T) {
	wantNoAllocs(t, "RuleOf72", func() { sinkFloat, _ = RuleOf72(0.09) })
}

func BenchmarkRuleOf72(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = RuleOf72(0.09)
	}
}
