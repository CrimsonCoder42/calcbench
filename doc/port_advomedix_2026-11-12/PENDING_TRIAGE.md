# Pending triage (W3-B11, 2026-11-12): TRIAGE ONLY, nothing built

Scope: every non-test `.go` file under the reference's `knowledgebase/pending/` (49 files: analytics 10, blood_pressure 1,
blood_sugar 2, calculators 21, knowledge 10, oxygen_saturation 1, peak_flow 1, pulse 1, respiration_rate 1,
temperature 1). The only `_test.go` file is `calculators/context_test.go` (it covers ComputeContext, which is plumbing).
Context: `pending/PENDING.md` says these bodies were moved unchanged out of the old `analytics/` tree and do not reach the
knowledgebase `Run`. `pending_frozen.txt` lists them as frozen (D-441). Several `calculators/*.go` (bayesian, medication,
risk, severity, symptom_scoring, vital_trends) are design notes around a `TODO: Implement` stub that returns nil.

Duplicate column: a calcbench ID (`grep "^// ID:"`), a knowledgebase calculator id (with its outcome in this run's
ledgers), or a function of the knowledgebase `series/` package (the W3-B9 scope, D-P8). "W3-B10" means the age-banded
reference ranges and classify pipelines that W3-B10 is porting from `ranges/`.

Decisions: **PORT LATER** = PORT IN A LATER WAVE (home package · primary to look for; a home in num/stats, num/prob or
algo is a PENDING_SHARED request to its owner, D-P8) · **DROP** (reason).

## Rows

| file:line | what it computes | duplicate of | decision |
|---|---|---|---|
| pending/calculators/baseline.go:31 | ComputePersonalBaseline: median of the trailing windowDays, refit after dropping values more than 2 SD from the median | series.ComputeBaseline (W3-B9); pieces NUM-STATS-012 Median, NUM-STATS-022 | DROP: app policy (2-SD exclusion rule and window are the app's; the cited WHO pulse-oximetry manual quote gives no formula) |
| pending/calculators/baseline.go:96 | ComputeBaselineDelta: min and mean deviation from a baseline, count/percent of readings ≥3 % below it | — | DROP: app policy (3 % threshold from a training-manual sentence, no primary) |
| pending/calculators/bayesian.go:196 | ComputeBayesian (stub): odds-form Bayes update, posterior odds = prior odds × LR₁ × … × LRₙ; LR+ = Se/(1−Sp), LR− = (1−Se)/Sp | NUM-PROB-008/009 (Bayes over a partition), NUM-PROB-012/013 (odds ⇄ probability); LR+/LR− listed LATER in doc/catalog/prob.md | PORT LATER: num/prob (request to prob owner): LikelihoodRatioPositive/Negative and PosteriorProbabilityFromLikelihoodRatios · primary Deeks & Altman, "Diagnostic tests 4: likelihood ratios", BMJ 2004;329:168. The prevalence and LR tables it needs are "TO BE BUILT" and are not ported |
| pending/calculators/classify.go:26 | ClassifySingleVital: age/sex/context band lookup against the vitals reference tables | knowledgebase ranges/ (W3-B10) | DROP: duplicate (W3-B10 ports the cited range tables and the classify pipeline) |
| pending/calculators/classify.go:72 | ClassifyBloodPressure: SBP/DBP band lookup by age (adult and pediatric tables) | knowledgebase ranges/blood_pressure (W3-B10) | DROP: duplicate (W3-B10) |
| pending/calculators/classify.go:137 | ReferenceSex, findAgeGroup, valueInBand, valueInRange: lookup helpers | — | DROP: plumbing |
| pending/calculators/classify.go:295 | ClassifyBMI: adult WHO cut-offs; under 20 y fixed cut-offs 14/22/27 | MED-ANTHRO-003 BMICategoryCDCAdult (adult part) | DROP: duplicate for adults; the pediatric cut-offs are app-made (CDC uses BMI-for-age percentiles), no primary |
| pending/calculators/classify.go:328 | ClassifyFeverDays: fever-day count → acute (≤2) / persistent (≤7) / prolonged | — | DROP: app policy, no primary |
| pending/calculators/classify.go:342 | ClassifyWeightChange: ±lb change within ≤2 days → rapid/moderate gain/loss | — | DROP: app policy (thresholds are not the heart-failure guideline's; decision rule) |
| pending/calculators/clustering.go:37 | TimeOfDayClustering / FromHourCounts: 4 day-part counts, Pearson chi-square vs uniform, flag if > 7.815 (df 3, α 0.05) | — (chi-square is LATER in doc/catalog/stats.md "inference slice") | PORT LATER: num/stats (request to stats owner): ChiSquareGoodnessOfFit(observed, expected) and, in num/prob, ChiSquareQuantile for the critical value · primary NIST/SEMATECH e-Handbook §1.3.5.15 (chi-square goodness-of-fit) and §1.3.6.6.6 (chi-square distribution). The day-part buckets are app policy |
| pending/calculators/clustering.go:121 | DayOfWeekPattern / FromCounts: 7 weekday counts, chi-square vs uniform, flag if > 12.592 (df 6) | same as row above | DROP: duplicate of the chi-square row above (same statistic, k = 7) |
| pending/calculators/clustering.go:180 | EpisodeCadence: mean and population SD of gaps between start dates, periodic if SD/mean < 0.3 | NUM-STATS-001, NUM-STATS-022, NUM-STATS-003 | DROP: duplicate (mean/SD/CV exist); the 0.3 "periodic" cut-off is app policy |
| pending/calculators/context.go:139 | ComputeContext: 30-day mean, P10–P90 "typical range" (Postgres percentile_cont), split-half direction | NUM-STATS-015 Percentile (Hyndman–Fan type 7), NUM-STATS-001; series.SplitHalfTrend (W3-B9) | DROP: product presentation plumbing; the maths is reused |
| pending/calculators/context.go:257 | ContextWindowEnd, LocalDay, InZone, ContextKeyPart, civilDayIndex: day boundaries and cache keys | num/series local_day_window.go, local_day_buckets.go | DROP: plumbing |
| pending/calculators/context.go:338 | ClassifyDisplayed / ClassifyDisplayedBP: classify a rounded displayed value, attach a short citation | W3-B10 classify pipeline | DROP: plumbing |
| pending/calculators/correlation.go:35 | DetectLiebermeister: day-over-day temperature rise ≥1 °F with any pulse rise | knowledgebase temp_liebermeister (ledger W1_B1: NOT PORTED, primary not reached; app constants match neither source) | DROP: duplicate of a NOT PORTED entry; app rule |
| pending/calculators/correlation.go:66 | DetectRespiratoryCompensation: SpO2 drop ≥2 % with pulse rise ≥10 bpm day over day | — | DROP: app policy, no primary |
| pending/calculators/correlation.go:114 | DetectRapidWeightChange: ≥2 lb gain / ≥3 lb loss within ≤2 days ("ACC/AHA weight monitoring") | same thresholds as classify.go:342 | DROP: app policy (decision rule; the cited guideline is not specific) |
| pending/calculators/correlation.go:172 | DetectSymptomCoOccurrence: counts unordered pairs of distinct names whose timestamps fall within windowHours | — | PORT LATER: algo/window (request to S1 owner): CoOccurringPairsWithinWindow, sort + two-pointer sweep · primary: algorithm textbook (sliding window over sorted events); no clinical claim |
| pending/calculators/distribution.go:27 | ComputeDistribution: count and percent per class label | NUM-STATS-002 RelativeFrequency, NUM-STATS-008 FrequencyDistribution | DROP: duplicate |
| pending/calculators/distribution.go:60 | DateBucketKey, IsHourlyPeriod, BucketByDate, CountByDate: day/hour bucket means and counts | num/series fixed_stride_buckets.go, local_day_buckets.go; ALG-WINDOW-002 | DROP: plumbing (duplicate of num/series bucketing) |
| pending/calculators/diurnal.go:34 | ComputeDiurnalBP / DiurnalBPFromWindowSums: mean SBP/DBP in 05–12 h and 17–24 h windows, their difference; flags surge > 20 mmHg and non-dipping | — | PORT LATER: num/series: generic hour-of-day window means and their difference (no clinical threshold) · primary: none needed beyond the arithmetic mean (cite NIST e-Handbook §1.3.5.1). The 20 mmHg surge and "non-dipping" flags are app policy (non-dipping is a nocturnal ABPM concept, not AM vs PM) and are not ported |
| pending/calculators/diurnal.go:114 | ComputeDiurnalPEF: per day (max PM − min AM)/mean × 100, averaged; GINA >10 % adult / >13 % child; control bands 10/20 % | knowledgebase pef_variability (ledger: NOT PORTED, GINA not opened) | PORT LATER: med/respiratory: PEFDiurnalVariabilityPercent ((day's highest − day's lowest)/mean of the two × 100, averaged over the days) · primary GINA Global Strategy report (Box 1-2) or the ERS working party on PEF (Quanjer et al., Eur Respir J 1997 suppl). Note: the app's max-PM/min-AM form is not GINA's and must not be the parity target |
| pending/calculators/diurnal.go:243 | ComputeMorningDip: per day (PM mean − AM mean)/PM mean × 100; detected if ≥3 days > 20 % | — | DROP: app policy (Hetzel & Clark 1980 describe the phenomenon, not this rule) |
| pending/calculators/diurnal.go:356 | ComputeDiurnalO2 / DiurnalO2FromWindowSums: day vs night SpO2 means, flag delta > 2 % | generic part = diurnal.go:34 row | DROP: duplicate of the window-means row; the 2 % flag cites a consumer-watch study (secondary), app policy |
| pending/calculators/estimated_a1c.go:9 | EstimateA1C: A1C = (mean glucose + 46.7)/28.7, then ADA bands | knowledgebase bs_a1c_from_eag (ledger: NOT PORTED, an inverted regression is not a published equation); published alternative MED-ENDO-001 GMI | DROP: duplicate of a NOT PORTED entry |
| pending/calculators/mage.go:10 | ComputeMAGE: mean of excursions ≥1 SD (turning point vs its immediate neighbour; fallback consecutive differences) | series.MAGE (W3-B9 scope) | PORT LATER: med/endocrine (only if W3-B9 does not land it): MAGE per Service et al., Diabetes 1970;19:644–655 (PMID 5469118); look for the peak-to-nadir definition there. The source's algorithm is not Service's (it compares a turning point with its neighbour, not the previous turning point) and must not be the parity target |
| pending/calculators/medication.go:62 | ComputeMedicationFlags (stub): allergy match, class duplication, polypharmacy ≥5/≥10, implied conditions | — | DROP: stub; app policy and vocabulary lookups |
| pending/calculators/risk.go:283 | ComputeRisk (stub) documenting the Framingham general CVD 10-year risk 1 − S₀^exp(Σβx − mean) | none in calcbench; knowledgebase framingham_hard_chd (Wilson 1998, NOT PORTED) and ascvd_risk (NOT PORTED) are different models | PORT LATER: med/cardio: FraminghamGeneralCVD10y · primary D'Agostino et al., Circulation 2008;117:743–753, doi:10.1161/CIRCULATIONAHA.107.699579 (PMID 18212285; free full text, coefficients in Table 2) |
| pending/calculators/risk.go:248 | (reference [7] in the same stub) FINDRISC diabetes risk questionnaire | none | PORT LATER: med/endocrine: FINDRISCScore · primary Lindström & Tuomilehto, Diabetes Care 2003;26:725–731 (PMID 12610029) |
| pending/calculators/risk.go:208 | (same stub) "Metabolic Risk Score" composite and pediatric growth trajectory | — | DROP: app policy ("no published regression"); growth percentiles belong to the CDC tables (W3-B10 / med/pediatrics) |
| pending/calculators/severity.go:592 | ComputeSeverity (stub): routes age <18 → PEWS, ≥18 → NEWS2/qSOFA | — | DROP: stub; routing is app policy |
| pending/calculators/severity.go:6 | (notes) qSOFA | MED-CRIT-008 QSOFAScore | DROP: duplicate |
| pending/calculators/severity.go:600 | (notes) NEWS2 | MED-CRIT-001 NEWS2, MED-CRIT-002 | DROP: duplicate |
| pending/calculators/severity.go:281 | (notes) Brighton PEWS | knowledgebase pews (ledger: NOT PORTED, Monaghan 2005 chart not seen) | DROP: duplicate of a NOT PORTED entry |
| pending/calculators/severity.go:440 | (notes) SOFA component bands and mortality bands | knowledgebase sepsis/sofa_score (W2-B6 scope, med/critical) | DROP: duplicate |
| pending/calculators/symptom_scoring.go:313 | ComputeSymptomScore (stub): severity × duration × recency, system diversity bonus, 0–100 bands | — | DROP: stub; app composite, no primary |
| pending/calculators/time_in_range.go:17 | ComputeTimeInRange: % of readings <54, 54–69, 70–180, 181–250, >250 mg/dL; target TIR > 70 % and TBR < 4 % | series.TimeInRange (W3-B9); counting = NUM-STATS-008 FrequencyDistribution; knowledgebase bs_cgm_metrics (NOT PORTED: Battelino not read) | PORT LATER: med/endocrine: CGMTimeInRangesBattelino (fixed consensus edges over NUM-STATS-008) · primary Battelino et al., Diabetes Care 2019;42:1593–1603 (PMID 31177185, PMC6973648). The "meets target" flag is a decision rule; port the bands only |
| pending/calculators/trend.go:28 | ComputeTrend / TrendFromAggregates: split-half means, % change, "stable" if abs(Δ) < 2 % | series.SplitHalfTrend (W3-B9) | DROP: app policy (2 % band); the means reuse NUM-STATS-001 |
| pending/calculators/trend.go:116 | LinearSlope / LinearSlopeFromSums: OLS slope on index and R² | NUM-STATS-009 LinearRegression (slope), NUM-STATS-014 PearsonCorrelation (R² = r²), NUM-STATS-025 SumOfSquares | DROP: duplicate |
| pending/calculators/types.go:20 | snapshot / result / entry types | — | DROP: plumbing |
| pending/calculators/unit_convert.go:13 | KgToLbs, LbsToKg, CmToInches, InchesToCm | none (no unit/mass, unit/length yet) | PORT LATER: unit/mass and unit/length · primary NIST SP 811 (2008) Appendix B (1 lb = 0.453 592 37 kg, 1 in = 2.54 cm exactly) |
| pending/calculators/unit_convert.go:37 | FToC, CToF | UNT-TEMP-001 FahrenheitToCelsius, UNT-TEMP-002 CelsiusToFahrenheit | DROP: duplicate |
| pending/calculators/unit_convert.go:50 | ComputeBMI | MED-ANTHRO-001 BodyMassIndex | DROP: duplicate |
| pending/calculators/unit_convert.go:61 | ComputeMAP: DBP + (SBP − DBP)/3 | knowledgebase map_calc (held: doc/catalog/medicine.md §5, one-third rule, no primary) | DROP: duplicate of a held row |
| pending/calculators/unit_convert.go:67 | ComputeAgeYears, ComputeAgeFromTime, RoundTo | — | DROP: plumbing (reads the clock; rounding helper) |
| pending/calculators/variability.go:12 | ComputeSD: population SD rounded to 2 dp | NUM-STATS-022 PopulationStandardDeviation | DROP: duplicate |
| pending/calculators/variability.go:42 | ComputeARV: Σ abs(xᵢ₊₁ − xᵢ)/(n − 1) | series.ARV (W3-B9 scope) | PORT LATER: num/series (only if W3-B9 does not land it): AverageRealVariability · CORRECTED primary: Mena et al., J Hypertens 2005;23:505–511 (ARV defined there; the source cites Rothwell 2010, which uses SD/CV/VIM) |
| pending/calculators/variability.go:60 | ComputeCV: rounded SD / mean × 100 | NUM-STATS-003 CoefficientOfVariation | DROP: duplicate |
| pending/calculators/velocity.go:32 | ComputeVelocity: (last − first)/(Δdays), and per 30.44-day month; nil if span < 7 d | — | PORT LATER: num/series: AverageRateOfChangePerDay (secant slope between two timestamped values) · primary: any calculus text (average rate of change); the 7-day floor and 30.44-day month are app choices, not ported |
| pending/calculators/velocity.go:58 | GrowthVelocity: velocity row above on heights in inches and weights in lb | velocity.go:32 row + unit/length, unit/mass row | DROP: duplicate (composition) |
| pending/calculators/velocity.go:103 | ComputeSeverityProgression: split-half trend + OLS slope per day + direction labels | trend.go rows | DROP: duplicate (composition) and app labels |
| pending/calculators/vital_trends.go:339 | ComputeVitalAnalysis (stub): good/bad direction tally, Liebermeister and compensation pairs, alerts | correlation.go rows | DROP: stub; product report plumbing |
| pending/calculators/weighted.go:44 | WeightedSD / WeightedCV: population SD and CV over (value, count) pairs | NUM-STATS-005 GroupedMean (mean only); no frequency-weighted SD in num/stats | PORT LATER: num/stats (request to stats owner): FrequencyWeightedStandardDeviation (population) · primary NIST Dataplot Reference Manual, "WEIGHTED STANDARD DEVIATION" (or NIST e-Handbook) |
| pending/calculators/weighted.go:74 | WeightedDistribution: class percentages over counts | NUM-STATS-002 RelativeFrequency | DROP: duplicate |
| pending/calculators/weighted.go:107 | weightedMedian: median over (value, count) pairs | NUM-STATS-006 GroupedMedian interpolates bins; no exact weighted median | PORT LATER: num/stats (request to stats owner): WeightedMedian (lower/upper convention stated) · primary NIST Dataplot "WEIGHTED MEDIAN" or Hyndman & Fan 1996 for the convention |
| pending/calculators/weighted.go:137 | PersonalBaselineFromCounts / BaselineDeltaFromCounts | baseline.go rows | DROP: app policy (same as baseline.go) |
| pending/blood_sugar/compute.go:50 | ClassifyBloodSugar (FPG, A1C, OGTT, random, hypoglycemia levels) | knowledgebase bs_classifier (NOT PORTED: website thresholds); W3-B10 blood-sugar slice | DROP: duplicate (W3-B10 blood-sugar slice owns the cited bands) |
| pending/blood_sugar/compute.go:191 | DetectCrisisRisk, DetectWarnings, GetDoctorQuestions | — | DROP: app policy and product text |
| pending/blood_sugar/compute.go:288 | ScreenGDM | MED-ENDO-005/006 GDMThresholdIADPSG | DROP: duplicate |
| pending/blood_sugar/compute.go:405 | ComputeEAGFromA1C, ComputeEAGMmolFromA1C | MED-CHEM-005, MED-CHEM-006 EstimatedAverageGlucose | DROP: duplicate |
| pending/blood_sugar/compute.go:424 | ComputeA1CFromEAG | knowledgebase bs_a1c_from_eag (NOT PORTED, inverted regression) | DROP: duplicate of a NOT PORTED entry |
| pending/blood_sugar/compute.go:487 | ComputeGMI, AnalyzeGMI, AnalyzeEAG | MED-ENDO-001/002 GlucoseManagementIndicator | DROP: duplicate (Analyze* = report plumbing) |
| pending/blood_sugar/compute.go:566 | ComputeCV(sd, mean), AnalyzeCGMMetrics (CV ≤ 36 % target, TIR targets) | NUM-STATS-003; time_in_range.go row; knowledgebase bs_cgm_metrics (NOT PORTED) | DROP: duplicate (targets are decision rules) |
| pending/blood_sugar/compute.go:664 | ComputeHOMAIR / ComputeHOMAIRMmol: glucose × insulin / 405 (mmol: / 22.5) | knowledgebase bs_homa_ir, homa_ir (ledger: NOT PORTED, Matthews formula not seen) | PORT LATER: med/endocrine: HOMAIR (retry the primary) · primary Matthews et al., Diabetologia 1985;28:412–419 (PMID 3899825), full text; or Wallace et al., Diabetes Care 2004;27:1487 (PMID 15161807) |
| pending/blood_sugar/compute.go:761 | ComputeHOMAB: 20 × insulin / (glucose − 3.5) (mmol) | knowledgebase bs_homa_b (NOT PORTED) | PORT LATER: med/endocrine: HOMABeta · primary Matthews 1985 as above |
| pending/blood_sugar/compute.go:843 | ComputeQUICKI: 1/(log₁₀ insulin + log₁₀ glucose) | knowledgebase bs_quicki (NOT PORTED: units and log base not in the abstract) | PORT LATER: med/endocrine: QUICKI · primary Katz et al., J Clin Endocrinol Metab 2000;85:2402 (PMID 10902785), full text for units/base |
| pending/blood_sugar/compute.go:680 | ClassifyHOMAIR/B, ClassifyQUICKI, Detect*Warnings, Analyze* | — | DROP: app bands (no cited cut-off table) and product text |
| pending/blood_sugar/compute.go:924 | ComputeCorrectedSodium (Katz/Hillier) | knowledgebase bs_corrected_na (held: doc/catalog/medicine.md §5) | DROP: duplicate of a held row |
| pending/blood_sugar/compute.go:968 | ComputeOsmolality: 2Na + glucose/18 + BUN/2.8 | MED-CHEM-003/004 SerumOsmolalityDorwart (ledger bs_osmolality NOT PORTED: app form is not Dorwart's) | DROP: duplicate |
| pending/blood_sugar/compute.go:1011 | ComputeTDD, ComputeISF, ComputeICR, ComputeBolus | knowledgebase bs_insulin_dosing, bs_bolus (held D-M5) | DROP: dose / amount to administer (D-M5) |
| pending/blood_sugar/compute.go:1159 | ComputeGIR | knowledgebase bs_gir, glucose_infusion_rate (held D-M5) | DROP: dose / amount to administer (D-M5) |
| pending/blood_sugar/knowledge.go:245 | citations, adult classification table, warning signs, doctor questions for each metric | — | DROP: vocabulary/reference data, not a calculation |
| pending/blood_pressure/knowledge.go:120 | citations, AHA classifications, warning signs, secondary-HTN clues, doctor questions (BP, PP, MAP, CPP, SVR, PVR, ABI) | — | DROP: vocabulary/reference data, not a calculation |
| pending/blood_pressure/knowledge.go:236 | ResolveAgeGroup, PediatricRangeFor, older-adult and sex-context constants | knowledgebase ranges/blood_pressure (W3-B10) | DROP: duplicate (reference ranges, W3-B10) |
| pending/oxygen_saturation/knowledge.go:121 | citations, SpO2 classifications, COPD and neonatal notes, warning signs (P/F, S/F, OI, OSI, A–a), doctor questions | MED-RESP-013/015/018/019 for the ratios themselves | DROP: vocabulary/reference data, not a calculation |
| pending/oxygen_saturation/knowledge.go:1207 | NEWS2SpO2ScoreScale1 / Scale2 | MED-RESP-009, MED-RESP-010 | DROP: duplicate |
| pending/peak_flow/knowledge.go:143 | citations, zone/predicted classifications, zone action plans, variability/reversibility/exacerbation notes, doctor questions | knowledgebase pef_zone, pef_reversibility, pef_exacerbation (all NOT PORTED) | DROP: vocabulary/reference data and duplicates of NOT PORTED entries |
| pending/pulse/knowledge.go:165 | ResolveAgeGroup, PulseRangeFor | knowledgebase ranges/pulse (W3-B10) | DROP: duplicate (W3-B10) |
| pending/pulse/knowledge.go:462 | ExerciseTargetRange: (220 − age) × low/high % | knowledgebase hr_zones (held: medicine.md §2/§5, Fox not seen) | DROP: duplicate of a held row |
| pending/pulse/knowledge.go:474 | KarvonenTargetHR: (max − rest) × intensity + rest | knowledgebase hr_zones (held: Karvonen 1957 not read) | DROP: duplicate of a held row |
| pending/pulse/knowledge.go:83 | citations and warning signs (CRI, HRR1, QTc, CO, CI, DO2, SI, RPP), doctor questions | MED-CARDIO-028 (HRR1), MED-CARDIO-002 (QTc) for the formulas | DROP: vocabulary/reference data, not a calculation |
| pending/respiration_rate/knowledge.go:492 | ResolveAgeGroup, RRRangeFor | knowledgebase ranges/respiration_rate (W3-B10) | DROP: duplicate (W3-B10) |
| pending/respiration_rate/knowledge.go:376 | citations and warning signs (VE, VA, VD/VT, RSBI, ROX, qSOFA, SIRS, NEWS2, MEWS, CURB-65, CRB-65, PSI, SMART-COP, PESI, RASI, APGAR, Silverman) | the scores themselves: MED-RESP-011/012/017, MED-CRIT-001/008, MED-CARDIO-033, MED-SCALE-005, knowledgebase entries | DROP: vocabulary/reference data, not a calculation |
| pending/temperature/knowledge.go:254 | ResolveAgeGroup, NormalRangeFor | knowledgebase ranges/temperature (W3-B10) | DROP: duplicate (W3-B10) |
| pending/temperature/knowledge.go:518 | PediatricFeverUrgencyFor: age band → urgency | knowledgebase acep_peds_fever / pecarn (MED-CRIT-007) | DROP: decision pathway, not a calculation |
| pending/temperature/knowledge.go:154 | citations and warning signs (fever, Liebermeister, ABG correction, metabolic cost) | knowledgebase temp_liebermeister, temp_metabolic_cost (both NOT PORTED) | DROP: vocabulary/reference data |
| pending/temperature/knowledge.go:1098 | SIRSTempCriterionMet (> 38 °C or < 36 °C) | knowledgebase sepsis/sirs (W2-B6 scope) | DROP: duplicate |
| pending/temperature/knowledge.go:1136 | NEWS2TempScore | MED-CRIT-001 NEWS2 (temperature band) | DROP: duplicate |
| pending/temperature/knowledge.go:1179 | MEWSTempScore | knowledgebase sepsis/mews (W2-B6 scope) | DROP: duplicate |
| pending/temperature/knowledge.go:1238 | APACHEIITempScore | knowledgebase apache_ii (ledger: NOT PORTED, Knaus 1985 point tables not open) | DROP: duplicate of a NOT PORTED entry |
| pending/knowledge/allergy_reference.go:46 | allergy list load/search (embeds allergy_reference/*.json) | — | DROP: vocabulary, not a calculation |
| pending/knowledge/conditions_reference.go:85 | conditions list, by domain/age/symptom (embeds conditions_reference/*.json) | — | DROP: vocabulary, not a calculation |
| pending/knowledge/icd10_codes.go:33 | ICD-10 lookup/validate (embeds icd10_codes.json) | — | DROP: vocabulary, not a calculation |
| pending/knowledge/infant_care_reference.go:158 | feeding, sleep, diaper, growth benchmarks by age (embeds infant_care_reference/*.json) | MED-PEDS-002 RecommendedSleepHoursAASM (sleep hours); knowledgebase infant_feeding_calculator (held D-M5) | DROP: reference data, not a calculation |
| pending/knowledge/medications_reference.go:39 | medication list (embeds medications_reference.json) | — | DROP: vocabulary, not a calculation |
| pending/knowledge/prn_reference.go:77 | PRN medication list (embeds prn_reference/*.json) | — | DROP: vocabulary, not a calculation |
| pending/knowledge/scheduled_medications_reference.go:85 | scheduled medication list (embeds scheduled_reference/*.json) | — | DROP: vocabulary, not a calculation |
| pending/knowledge/symptom_key.go:26 | symptom key and categories | — | DROP: vocabulary, not a calculation |
| pending/knowledge/test_results_reference.go:149 | test-result reference and articles (embeds test_results_reference/*.json) | — | DROP: vocabulary, not a calculation |
| pending/knowledge/vitals_reference.go:97 | vital reference tables (embeds vitals_reference/*.json: BP, glucose, SpO2, PEF, pulse, RR) | knowledgebase ranges/ (W3-B10) | DROP: duplicate (reference ranges, W3-B10) |
| pending/analytics/cardiovascular_reports.go:616 | ComputeCardiovascularReport: BP/pulse/SpO2 deep analytics, ICD-10 mapping, alerts, guidance, summary card | rows above (trend, SD/ARV, diurnal, distribution, velocity) | DROP: product report plumbing |
| pending/analytics/fever_infection_reports.go:480 | ComputeFeverInfectionReport: temperature analytics, consecutive fever days, antipyretic context, co-occurrence, Liebermeister | rows above | DROP: product report plumbing |
| pending/analytics/fitness.go:143 | FitnessVolume: Σ weight × reps | LIF-FIT-008 VolumeLoad | DROP: duplicate |
| pending/analytics/fitness.go:437 | FitnessNextLoad: W × 1.05 rounded to a plate step, accepted if within +2…+10 % | LIF-FIT-001 LoadProgressionACSM | DROP: duplicate (plate rounding is app policy) |
| pending/analytics/fitness.go:74 | plan week, session days, previous/best set, suggestions, page size | — | DROP: product plumbing |
| pending/analytics/fitness_summary.go:922 | FitnessCardioLoadOf: Edwards TRIMP (minutes in 50/60/70/80/90 % HRmax zones × 1…5) | series.EdwardsTRIMP / EdwardsZone (W3-B9 scope, life/fitness) | DROP: duplicate (W3-B9) |
| pending/analytics/fitness_summary.go:676 | fitnessPaceOrSpeed | series.Pace (W3-B9 scope) | DROP: duplicate (W3-B9) |
| pending/analytics/fitness_summary.go:338 | FitnessSummaryCompute, schedules, muscle groups/load, activities | — | DROP: product report plumbing |
| pending/analytics/infant_care_reports.go:485 | ComputeInfantCareReport: feeding/sleep/diaper trends, wake windows, benchmarks | infant_care_reference row | DROP: product report plumbing |
| pending/analytics/infant_care_reports.go:1759 | ComputeWetDiaperCheck: wet diapers per completed day vs age benchmark | — | DROP: app policy against reference data, no primary |
| pending/analytics/infant_care_reports.go:1675 | ComputeCompletedDayAverage: per-day mean over completed days | NUM-STATS-001; num/series local_day_buckets.go | DROP: plumbing |
| pending/analytics/metabolic_reports.go:550 | ComputeMetabolicReport: growth, BMI trend, growth velocity, blood-sugar analytics, ICD-10, alerts | MED-ANTHRO-001; velocity.go and time_in_range.go rows | DROP: product report plumbing |
| pending/analytics/metabolic_reports.go:1167 | indexSums: Sxx, Syy, Sxy on index | NUM-STATS-025 SumOfSquares, NUM-STATS-009 | DROP: duplicate |
| pending/analytics/respiratory_reports.go:750 | ComputeRespiratoryReport: SpO2/RR/PEF analytics, desaturation days, personal best, occupational pattern | diurnal.go rows | DROP: product report plumbing |
| pending/analytics/sleep_score.go:352 | ComputeSleepScore: home-made 0–100 composite (duration, efficiency, latency, WASO, wakings, regularity points) | MED-PEDS-002 (recommended hours) | DROP: app composite, no primary |
| pending/analytics/sleep.go:128 | SleepKindOf, SleepStagesFromLaps, SleepStagesFromFellAsleep | — | DROP: plumbing |
| pending/analytics/vital_series.go:96 | seriesFromRows, *SeriesFromAggregates, statsSeries, VitalSeries methods | num/series bucketing | DROP: plumbing |

## Summary

| decision | rows |
|---|---|
| PORT IN A LATER WAVE | 17 |
| DROP: duplicate (calcbench ID, NOT PORTED / held knowledgebase entry, or W2-B6 / W3-B9 / W3-B10 scope) | 47 |
| DROP: vocabulary / reference data | 16 |
| DROP: plumbing / product report plumbing | 18 |
| DROP: app policy, no primary (incl. stubs and decision pathways) | 18 |
| DROP: dose / amount to administer (D-M5) | 2 |
| **total** | **118** |

Later wave, by home package:

| home | function to write | primary to find | owner |
|---|---|---|---|
| med/endocrine | HOMAIR, HOMABeta | Matthews 1985 Diabetologia (PMID 3899825), full text | this run (W1 NOT PORTED → retry) |
| med/endocrine | QUICKI | Katz 2000 JCEM (PMID 10902785), full text | this run (retry) |
| med/endocrine | CGMTimeInRangesBattelino | Battelino 2019 Diabetes Care (PMID 31177185, PMC6973648) | this run (if W3-B9 does not) |
| med/endocrine | MAGE | Service 1970 Diabetes 19:644 (PMID 5469118) | this run (if W3-B9 does not) |
| med/endocrine | FINDRISCScore | Lindström & Tuomilehto 2003 Diabetes Care (PMID 12610029) | this run |
| med/cardio | FraminghamGeneralCVD10y | D'Agostino 2008 Circulation (PMID 18212285, free full text) | this run |
| med/respiratory | PEFDiurnalVariabilityPercent | GINA report Box 1-2 / ERS PEF working party (Quanjer 1997) | this run (W1 NOT PORTED → retry) |
| num/series | AverageRealVariability | Mena 2005 J Hypertens 23:505 (corrects the source's Rothwell citation) | this run (if W3-B9 does not) |
| num/series | AverageRateOfChangePerDay | calculus text | this run |
| num/series | hour-of-day window means and their difference | NIST e-Handbook §1.3.5.1 | this run |
| num/stats | ChiSquareGoodnessOfFit | NIST e-Handbook §1.3.5.15 | stats subject (PENDING_SHARED request) |
| num/stats | FrequencyWeightedStandardDeviation | NIST Dataplot "weighted standard deviation" | stats subject (request) |
| num/stats | WeightedMedian | NIST Dataplot / Hyndman & Fan 1996 | stats subject (request) |
| num/prob | ChiSquareQuantile (critical values 7.815, 12.592) | NIST e-Handbook §1.3.6.6.6 | prob subject (request) |
| num/prob | LikelihoodRatioPositive/Negative, PosteriorProbabilityFromLikelihoodRatios | Deeks & Altman 2004 BMJ 329:168 | prob subject (request; prob.md lists it LATER) |
| unit/mass, unit/length | kg ⇄ lb, cm ⇄ in | NIST SP 811 (2008) App. B | this run (new packages from doc/taxonomy.md) |
| algo/window | CoOccurringPairsWithinWindow | algorithm textbook (sorted sweep) | S1 (request) |

(17 later-wave rows map to the 17 table lines above; ChiSquareQuantile rides on the clustering.go row, and HOMAIR and
HOMABeta are two rows.)
