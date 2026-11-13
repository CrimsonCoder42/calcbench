package cardio

import "math"

// PCPHFInputs are the inputs of the Pooled Cohort equations to Prevent Heart Failure (Khan 2019 Table 2).
// Female and Black select the race- and sex-specific equation (the paper's White / Black, male / female groups).
type PCPHFInputs struct {
	Female               bool
	Black                bool
	AgeYears             float64
	SystolicMmHg         float64
	TreatedBP            bool // systolic pressure on antihypertensive treatment
	CurrentSmoker        bool
	GlucoseMgDL          float64 // fasting glucose
	TreatedGlucose       bool    // glucose on treatment
	TotalCholesterolMgDL float64
	HDLMgDL              float64
	BMIKgM2              float64
	QRSDurationMs        float64
}

// pcpHFCoef holds one column of Khan 2019 Table 2; zero where the table prints N/A.
type pcpHFCoef struct {
	lnAge, lnAgeSq, lnTrSBP, lnAgeTrSBP, lnUnSBP, lnAgeUnSBP, smoker, lnAgeSmoker,
	lnTrGlu, lnUnGlu, lnTC, lnHDL, lnBMI, lnAgeBMI, lnQRS, meanCV, s0 float64
}

// pcpHFTable is Khan 2019 Table 2, columns White male, White female, Black male, Black female.
var pcpHFTable = [4]pcpHFCoef{
	{41.94, -0.88, 1.03, 0, 0.91, 0, 0.74, 0, 0.90, 0.78, 0.49, -0.44, 37.2, -8.83, 0.63, 171.5, 0.98752},
	{20.55, 0, 12.95, -2.96, 11.86, -2.73, 11.02, -2.50, 1.04, 0.91, 0, -0.07, 1.33, 0, 1.06, 99.73, 0.99348},
	{2.88, 0, 2.31, 0, 2.17, 0, 1.66, -0.25, 0.64, 0.58, 0, -0.81, 1.16, 0, 0.73, 28.73, 0.98295},
	{51.75, 0, 29.0, -6.59, 28.18, -6.42, 0.76, 0, 0.97, 0.80, 0.32, 0, 21.24, -5.0, 1.27, 233.9, 0.99260},
}

// PCPHF10YearRiskPercent returns the 10-year risk of incident heart failure, in percent, by the PCP-HF
// equations: 100 × (1 − S0^exp(IndX − MeanCV)), IndX = Σ coefficient × value over the race- and sex-specific
// terms of Table 2 (natural logs of age, systolic pressure, glucose, total and HDL cholesterol, BMI and QRS
// duration, with the age interactions the table lists).
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range.
//
// ID:      MED-CARDIO-060
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Khan SS, Ning H, Shah SJ, et al. 10-year risk equations for incident heart failure in the general population. J Am Coll Cardiol. 2019;73(19):2388-2397. doi:10.1016/j.jacc.2019.02.057. PMID 31097157. PMC6527121, Table 2 (coefficients, MeanCV, S0) and Methods ("1−S0^e(IndX−MeanCV)")
// Tags:    cardio, heart-failure, risk-equation, pcp-hf, pooled-cohort
// Units:   AgeYears years; SystolicMmHg mmHg; GlucoseMgDL, TotalCholesterolMgDL, HDLMgDL mg/dL; BMIKgM2 kg/m²; QRSDurationMs ms; result percent
// Range:   AgeYears 30–80 ("between 30 to 80 years", free of HF at baseline); other numeric inputs > 0 (no limits stated)
// Example: White man, 55 y, SBP 120 untreated, non-smoker, glucose 100, TC 200, HDL 50, BMI 27, QRS 100 → 1.439% (HAND-DERIVED from Table 2; see the test)
func PCPHF10YearRiskPercent(in PCPHFInputs) (float64, error) {
	for _, p := range [...]struct {
		name string
		v    float64
	}{
		{"AgeYears", in.AgeYears}, {"SystolicMmHg", in.SystolicMmHg}, {"GlucoseMgDL", in.GlucoseMgDL},
		{"TotalCholesterolMgDL", in.TotalCholesterolMgDL}, {"HDLMgDL", in.HDLMgDL}, {"BMIKgM2", in.BMIKgM2},
		{"QRSDurationMs", in.QRSDurationMs},
	} {
		if err := checkFinite(p.name, p.v); err != nil {
			return 0, err
		}
	}
	if err := checkRange("AgeYears", in.AgeYears, 30, 80); err != nil {
		return 0, err
	}
	for _, p := range [...]struct {
		name string
		v    float64
	}{
		{"SystolicMmHg", in.SystolicMmHg}, {"GlucoseMgDL", in.GlucoseMgDL},
		{"TotalCholesterolMgDL", in.TotalCholesterolMgDL}, {"HDLMgDL", in.HDLMgDL}, {"BMIKgM2", in.BMIKgM2},
		{"QRSDurationMs", in.QRSDurationMs},
	} {
		if err := checkPositive(p.name, p.v); err != nil {
			return 0, err
		}
	}
	g := 0
	if in.Female {
		g = 1
	}
	if in.Black {
		g += 2
	}
	c := &pcpHFTable[g]
	lnAge := math.Log(in.AgeYears)
	lnSBP := math.Log(in.SystolicMmHg)
	lnBMI := math.Log(in.BMIKgM2)
	x := c.lnAge*lnAge + c.lnAgeSq*lnAge*lnAge
	if in.TreatedBP {
		x += c.lnTrSBP*lnSBP + c.lnAgeTrSBP*lnAge*lnSBP
	} else {
		x += c.lnUnSBP*lnSBP + c.lnAgeUnSBP*lnAge*lnSBP
	}
	if in.CurrentSmoker {
		x += c.smoker + c.lnAgeSmoker*lnAge
	}
	if in.TreatedGlucose {
		x += c.lnTrGlu * math.Log(in.GlucoseMgDL)
	} else {
		x += c.lnUnGlu * math.Log(in.GlucoseMgDL)
	}
	x += c.lnTC*math.Log(in.TotalCholesterolMgDL) + c.lnHDL*math.Log(in.HDLMgDL)
	x += c.lnBMI*lnBMI + c.lnAgeBMI*lnAge*lnBMI + c.lnQRS*math.Log(in.QRSDurationMs)
	return 100 * (1 - math.Pow(c.s0, math.Exp(x-c.meanCV))), nil
}
