# Entry list W3_B9b (wave 2): peak flow and blood-sugar ranges

Source root (read-only): `/Users/devinanderson/Desktop/advo_production/advomedix-v2/backendAdvo/`. Paths are as wave 1 recorded them: `knowledgebase/calculators/<category>/...`, `<category>/...` (= knowledgebase/calculators/<category>/), `knowledgebase/ranges/...`, `calculators/<vital>/...`, or relative to `knowledgebase/series/` for series functions. The wave 1 reason is the starting point, not the answer.

| # | entry id | source path (wave 1) | wave 1 section | wave 1 outcome | wave 1 reason |
|---|---|---|---|---|---|
| 5 | blood_sugar:all_readings.go | calculators/blood_sugar/all_readings.go:1 | W3_B10 | NOT PORTED | deferred: blood sugar slice; not researched |
| 6 | blood_sugar:bs_classifier.go | calculators/blood_sugar/bs_classifier.go:1 | W3_B10 | NOT PORTED | deferred: blood sugar slice; not researched |
| 7 | blood_sugar:bs_trend.go | calculators/blood_sugar/bs_trend.go:1 | W3_B10 | NOT PORTED | deferred: blood sugar slice; not researched |
| 8 | blood_sugar:coordinator.go | calculators/blood_sugar/coordinator.go:1 | W3_B10 | NOT PORTED | deferred: blood sugar slice; not researched |
| 41 | bs_a1c_adult | knowledgebase/ranges/vitals/blood_sugar/bs_a1c_adult.json:1 | W3_B10 | NOT PORTED | deferred: blood sugar slice (home med/endocrine is being extended by another builder); not researched |
| 42 | bs_cgm_metrics_category | knowledgebase/ranges/vitals/blood_sugar/bs_cgm_metrics_category.json:1 | W3_B10 | NOT PORTED | deferred: blood sugar slice (home med/endocrine is being extended by another builder); not researched |
| 44 | bs_fasting_child | knowledgebase/ranges/vitals/blood_sugar/bs_fasting_child.json:1 | W3_B10 | NOT PORTED | deferred: blood sugar slice (home med/endocrine is being extended by another builder); not researched |
| 45 | bs_gestational_ogtt_one_step | knowledgebase/ranges/vitals/blood_sugar/bs_gestational_ogtt_one_step.json:1 | W3_B10 | NOT PORTED | deferred: blood sugar slice (home med/endocrine is being extended by another builder); not researched |
| 48 | bs_ogtt_2h_adult | knowledgebase/ranges/vitals/blood_sugar/bs_ogtt_2h_adult.json:1 | W3_B10 | NOT PORTED | deferred: blood sugar slice (home med/endocrine is being extended by another builder); not researched |
| 88 | peak_flow:all_readings.go | calculators/peak_flow/all_readings.go:99 | W3_B10 | NOT PORTED | plumbing; diurnalFromDays (:180) as pef_diurnal.go |
| 89 | peak_flow:coordinator.go | calculators/peak_flow/coordinator.go:54 | W3_B10 | NOT PORTED | plumbing |
| 90 | pef_diurnal.go | calculators/peak_flow/pef_diurnal.go:74 | W3_B10 | NOT PORTED | daily (max − min)/mean(max, min) × 100, GINA thresholds 10%/13%: GINA report not reached (PDF 404). Local-day grouping exists as NUM-SERIES LocalDayBuckets (reuse when ported) |
| 91 | pef_exacerbation_percent_personal_best | knowledgebase/ranges/vitals/peak_expiratory_flow/pef_exacerbation_percent_personal_best.json:1 | W3_B10 | NOT PORTED | as W1_B4 pef_exacerbation: BTS/SIGN 2019 asthma guideline severity table (33/50/75%) not opened |
| 92 | pef_exacerbation_percent_predicted | knowledgebase/ranges/vitals/peak_expiratory_flow/pef_exacerbation_percent_predicted.json:1 | W3_B10 | NOT PORTED | as pef_exacerbation_percent_personal_best |
| 93 | pef_percent_predicted_category | knowledgebase/ranges/vitals/peak_expiratory_flow/pef_percent_predicted_category.json:1 | W3_B10 | NOT PORTED | cites MDCalc (calculator site) + Hankinson 1999/Nunn & Gregg 1989 (equations, not bands); no primary for 50/80/100 bands (W1_B4 pef_percent_predicted: NOT PORTED) |
| 94 | pef_population_adult_female | knowledgebase/ranges/vitals/peak_expiratory_flow/pef_population_adult_female.json:1 | W3_B10 | NOT PORTED | cites Cleveland Clinic/ALA/AAFP pages (secondary); population bands have no primary; predicted-PEF equations NOT PORTED by W1_B4 |
| 95 | pef_population_adult_male | knowledgebase/ranges/vitals/peak_expiratory_flow/pef_population_adult_male.json:1 | W3_B10 | NOT PORTED | as pef_population_adult_female |
| 96 | pef_population_adult_sex_unspecified | knowledgebase/ranges/vitals/peak_expiratory_flow/pef_population_adult_sex_unspecified.json:1 | W3_B10 | NOT PORTED | as pef_population_adult_female; a sex-unspecified band has no primary at all |
| 97 | pef_population_child_5_to_18 | knowledgebase/ranges/vitals/peak_expiratory_flow/pef_population_child_5_to_18.json:1 | W3_B10 | NOT PORTED | as pef_population_adult_female |
| 98 | pef_reversibility_category | knowledgebase/ranges/vitals/peak_expiratory_flow/pef_reversibility_category.json:1 | W3_B10 | NOT PORTED | as W1_B4 pef_reversibility: NICE NG80 / GINA thresholds not opened |
| 99 | pef_variability_adult | knowledgebase/ranges/vitals/peak_expiratory_flow/pef_variability_adult.json:1 | W3_B10 | NOT PORTED | GINA report (ginasthma.org) not reached (PDF 404); > 10% adults / > 13% children not seen in a primary; as W1_B4 pef_variability |
| 100 | pef_variability_child | knowledgebase/ranges/vitals/peak_expiratory_flow/pef_variability_child.json:1 | W3_B10 | NOT PORTED | as pef_variability_adult |

## Also yours
| # | entry | what | research already done |
|---|---|---|---|
| X1 | TABLE Hankinson 1999 | med/respiratory: NHANES III spirometry reference equations (PMID 9872837, Tables 4–6); official coefficients only | ledger/W1_B2.md (pef_predicted is SOURCE DOES NOT SUPPORT) |
| X2 | candidate diurnal.go:114 PEFDiurnalVariabilityPercent | med/respiratory: GINA Box 1-2 / Quanjer 1997 ERS | PENDING_TRIAGE.md |
| X3 | candidate mage.go:10 MAGE | med/endocrine: Service 1970 Diabetes 19:644 (PMID 5469118) | PENDING_TRIAGE.md |
| X4 | candidate risk.go:248 FINDRISCScore | med/endocrine: Lindström & Tuomilehto 2003 Diabetes Care (PMID 12610029) | PENDING_TRIAGE.md |
| X5 | candidate time_in_range.go:17 CGM time in ranges | med/endocrine: Battelino 2019 Diabetes Care (PMID 31177185, PMC6973648) | PENDING_TRIAGE.md |
| X6 | candidates compute.go:664/761/843 HOMA-IR, HOMA-β, QUICKI | REUSED: MED-ENDO-024/025/026 (W1-B3) — ledger row only | ledger/W1_B3.md |
