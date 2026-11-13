# Build spec: med/anthropometry, med/scales (medicine build 2026-11-12)

Read `CONVENTIONS.md` first; it is binding. Every citation below was confirmed 2026-11-12 (PubMed record via
NCBI E-utilities, Crossref, or the document read as noted). Use the `Source:` text verbatim, one line each.
Never add a coefficient, threshold or label that is not written here.

## med/anthropometry (SHORT ANTHRO)

Conversion constants (unexported, in the file of the first function that needs them, with this comment):
none needed — every function here takes the unit its source uses.

### MED-ANTHRO-001 BodyMassIndex(weightKg, heightM float64) (float64 kg/m², error)
BMI = weight (kg) / [height (m)]².
Source: Centers for Disease Control and Prevention. BMI Frequently Asked Questions. https://www.cdc.gov/bmi/faq/index.html (last reviewed 2024-06-28; read 2026-11-12): "weight (kg)/[height (m)]2".
Range: weightKg > 0, heightM > 0 (the source states no limits). Example: 70 kg, 1.75 m → 22.857142857… (HAND-DERIVED).

### MED-ANTHRO-002 BodyMassIndexImperial(weightLb, heightIn float64) (float64 kg/m², error)
BMI = weight (lb) / [height (in)]² × 703 (the CDC's own factor; not exact — the test documents that it agrees with
001 on converted inputs to within 0.1 %, since 703 vs 703.0696 exact).
Source: same CDC page: "weight (pounds)/[height (in)]2 x 703".
Range: as 001. Example: 150 lb, 65 in → 150/4225 × 703 = 24.958579… (HAND-DERIVED).

### MED-ANTHRO-003 BMICategoryCDCAdult(bmi float64) (BMICategory, error)
`type BMICategory string`; constants, values VERBATIM from the CDC page (the source's labels; "Healthy" is quoted, not ours):
BMIUnderweight "Underweight" (less than 18.5); BMIHealthyWeight "Healthy Weight" (18.5 to less than 25);
BMIOverweight "Overweight" (25 to less than 30); BMIClass1Obesity "Class 1 Obesity" (30 to less than 35);
BMIClass2Obesity "Class 2 Obesity" (35 to less than 40); BMIClass3Obesity "Class 3 Obesity (Severe Obesity)" (40 or greater).
(The page's umbrella label "Obesity: 30 or greater" covers classes 1–3; the function returns the class.)
Source: Centers for Disease Control and Prevention. Adult BMI Categories. https://www.cdc.gov/bmi/adult-calculator/bmi-categories.html (last reviewed 2024-03-19; read 2026-11-12).
Range: bmi > 0 and finite; the page is the CDC's ADULT categories (say so; the function takes no age).
Tests: 18.4999/18.5, 24.9999/25, 29.9999/30, 34.9999/35, 39.9999/40.

### MED-ANTHRO-004 BodySurfaceAreaHaycock(weightKg, heightCm float64) (float64 m², error)
SA (m²) = weight^0.5378 × height^0.3964 × 0.024265.
Source: Haycock GB, Schwartz GJ, Wisotsky DH. Geometric method for measuring body surface area: a height-weight formula validated in infants, children, and adults. J Pediatr. 1978;93(1):62-66. doi:10.1016/s0022-3476(78)80601-5. PMID 650346. Abstract: "SA (m2) = weight (kg)0.5378 X height (cm)0.3964 X 0.024265".
Range: weightKg > 0, heightCm > 0; population 81 subjects "premature infants to adults", measured SA from less than 0.2 m² to greater than 2.0 m²; no hard input limits stated.
Example: 70 kg, 175 cm → compute independently in the test (write the formula).

### MED-ANTHRO-005 AdjustedBodyWeight(idealWeightKg, actualWeightKg float64) (float64 kg, error)
Adjusted weight = ideal + 0.4 × (actual − ideal) ("normalized body mass plus 40% of the adipose mass").
Source: Schwartz SN, Pazin GJ, Lyon JA, Ho M, Pasculle AW. A controlled investigation of the pharmacokinetics of gentamicin and tobramycin in obese subjects. J Infect Dis. 1978;138(4):499-505. doi:10.1093/infdis/138.4.499. PMID 712111. Abstract.
Range: idealWeightKg > 0; actualWeightKg > idealWeightKg (the formula's adipose mass is actual − ideal; population: 13 obese and 13 non-obese subjects). Doc note: the 0.4 was derived for aminoglycoside distribution volume (the source's context); the function returns a weight, nothing else.
Example: ideal 70, actual 120 → 90 (HAND-DERIVED).

### MED-ANTHRO-006 IdealBodyMassTraubKichen(heightCm, ageYears float64) (float64 kg, error)
IBM (kg) = 2.396 × e^(0.01863 × height cm).
Source: Traub SL, Kichen L. Estimating ideal body mass in children. Am J Hosp Pharm. 1983;40(1):107-110. PMID 6823980. Abstract: "IBM(kg) = 2.396e0.01863(ht)", ages 1–17, the 50th-percentile weight for height from NCHS data.
(NOT Traub & Johnson 1980, which is a creatinine-clearance paper.)
Range: ageYears 1–17; heightCm > 0 (no height limits stated). ageYears enforces the Range only.
Example: 100 cm → 2.396 × e^1.863 (compute in the test).

### MED-ANTHRO-007 EstimatedWeightInfantBestGuess(ageMonths float64) (float64 kg, error)
Weight (kg) = (age in months + 9) / 2, for infants under 12 months.
Source: Tinning K, Acworth J. Make your Best Guess: an updated method for paediatric weight estimation in emergencies. Emerg Med Australas. 2007;19(6):528-534. doi:10.1111/j.1742-6723.2007.01026.x. PMID 18021105. Abstract, the "<12 months" formula.
Range: 0 ≤ ageMonths < 12 (70,181 ED presentations, ages 0–14 y). The paper's other bands (1–4 y, 5–14 y) are NOT built (a band edge is unconfirmed).
Example: 6 months → 7.5 kg; 0 → 4.5.

### MED-ANTHRO-008 BodyFatPercentDeurenbergAdult(bmi, ageYears float64, male bool) (float64 %, error)
BF% = 1.20 × BMI + 0.23 × age − 10.8 × sex − 5.4, sex = 1 male, 0 female.
Source: Deurenberg P, Weststrate JA, Seidell JC. Body mass index as a measure of body fatness: age- and sex-specific prediction formulas. Br J Nutr. 1991;65(2):105-114. doi:10.1079/bjn19910073. PMID 2043597. Abstract (adults, R² 0.79, SEE 4.1%).
Range: ageYears > 15 and ≤ 83; bmi 13.9–40.9 (the study's 1229 subjects, ages 7–83, BMI 13.9–40.9).
Example: male, 40 y, BMI 25 → 30 + 9.2 − 10.8 − 5.4 = 23.0.

### MED-ANTHRO-009 BodyFatPercentDeurenbergChild(bmi, ageYears float64, male bool) (float64 %, error)
BF% = 1.51 × BMI − 0.70 × age − 3.6 × sex + 1.4. Same source (children ≤ 15 y, R² 0.38, SEE 4.4%).
Range: 7 ≤ ageYears ≤ 15; bmi 13.9–40.9. Example: female, 10 y, BMI 18 → 27.18 − 7 + 1.4 = 21.58.
Edge test: age 15 → child formula valid, adult refuses; age 15.0001 → adult valid, child refuses.

## med/scales (SHORT SCALE)

No item TEXT is embedded for Epworth (the scale is © Murray W Johns; the official site says a licence is needed;
decision D-M9: ship the sum and the published bands only, item text left to the licensed caller).

### MED-SCALE-001 EpworthSleepinessScore(items [8]int) (int, error)
Total = sum of the 8 item ratings, each 0–3 (Table 1 response key "0 = would never doze" … "3 = high chance of
dozing"); total 0–24.
Source: Johns MW. A new method for measuring daytime sleepiness: the Epworth sleepiness scale. Sleep. 1991;14(6):540-545. doi:10.1093/sleep/14.6.540. PMID 1798888. Table 1.
Range: each item 0–3 (ErrOutOfRange naming the item index). Population: 180 adults (aged 18–78, with sleep disorders, and controls).
Example: Johns 1991 controls range 2–10; test {1,1,1,1,1,1,1,1} → 8, {3×8} → 24, zeros → 0.

### MED-SCALE-002 EpworthCategory(score int) (EpworthCategory, error)
`type EpworthSleepinessCategory string` (renamed in the build: Go forbids a type and a function with the same name), values VERBATIM (official site):
0–5 "Lower Normal Daytime Sleepiness"; 6–10 "Higher Normal Daytime Sleepiness"; 11–12 "Mild Excessive Daytime
Sleepiness"; 13–15 "Moderate Excessive Daytime Sleepiness"; 16–24 "Severe Excessive Daytime Sleepiness".
Source: Johns MW. About the ESS. The Epworth Sleepiness Scale official website. https://epworthsleepinessscale.com/about-the-ess/ (read 2026-11-12).
Range: 0–24. Tests at 5/6, 10/11, 12/13, 15/16, 0, 24, −1, 25.

### MED-SCALE-003 GlasgowComaScore(eye, verbal, motor int) (int, error)
Total = eye (1–4) + verbal (1–5) + motor (1–6), range 3–15. A component the assessor rated "Non testable" (NT)
is passed as the exported constant `GCSNotTestable = 0`; then no total exists: return ErrNotTestable (new sentinel,
"scales: a component is Non testable (NT), so no total exists"). Doc may quote Teasdale 2014's abstract:
"Individual patients are best described by the three components of the coma scale; whereas the derived total
coma score should be used to characterise groups."
Source: Teasdale G. Glasgow Coma Scale: Assessment Aid (English). 2015. https://www.glasgowcomascale.org/downloads/GCS-Assessment-Aid-English.pdf (read 2026-11-12): eye 4 Spontaneous, 3 To sound, 2 To pressure, 1 None; verbal 5 Orientated, 4 Confused, 3 Words, 2 Sounds, 1 None; motor 6 Obeys commands, 5 Localising, 4 Normal flexion, 3 Abnormal flexion, 2 Extension, 1 None; NT Non testable.
Source: Teasdale G, Jennett B. Assessment of coma and impaired consciousness. A practical scale. Lancet. 1974;2(7872):81-84. doi:10.1016/s0140-6736(74)91639-0. PMID 4136544.
Source: Teasdale G, Maas A, Lecky F, Manley G, Stocchetti N, Murray G. The Glasgow Coma Scale at 40 years: standing the test of time. Lancet Neurol. 2014;13(8):844-854. doi:10.1016/S1474-4422(14)70120-6. PMID 25030516.
Export the rating labels as typed constants per component if useful (e.g. `GCSEyeSpontaneous = 4`) with the
official rating word in the comment. Range: eye 1–4, verbal 1–5, motor 1–6 (0 = NT).

### MED-SCALE-004 GlasgowComaSeverityRimel(total int) (GCSSeverity, error)
`type GCSSeverity string`, values VERBATIM lower-case: "severe" (3–8), "moderate" (9–12), "minor" (13–15).
Source: Rimel RW, Giordani B, Barth JT, Jane JA. Moderate head injury: completing the clinical spectrum of brain trauma. Neurosurgery. 1982;11(3):344-351. doi:10.1227/00006123-198209000-00002. PMID 7133350. Abstract: "severe, 3-8; moderate, 9-12; and minor, 13-15".
Range: 3–15. Tests 8/9, 12/13, 3, 15, 2, 16. Doc: these are Rimel's head-injury groups; the official GCS site states no bands.

### MED-SCALE-005 ApgarScore(heartRate, respiratoryEffort, reflexIrritability, muscleTone, color int) (int, error)
Total = sum of five signs, each rated 0, 1 or 2 ("A rating of zero, one or two, was given to each sign"); 0–10;
"A score of ten indicated a baby in the best possible condition."
Source: Apgar V. A proposal for a new method of evaluation of the newborn infant. Curr Res Anesth Analg. 1953;32(4):260-267. PMID 13083014. pp. 261–262 (read in the NLM scan).
Range: each 0–2. Doc: the paper rated at "sixty seconds after the complete birth of the baby".

### MED-SCALE-006 ApgarCondition1953(score int) (ApgarCondition, error)
`type ApgarCondition string`: "good condition" (8, 9, 10), "fair condition" (3–7, "the remaining scores"),
"poor condition" (0, 1, 2) — Apgar 1953 p. 267: "groups 8, 9, and 10 indicate infants in good condition, 0, 1,
and 2, poor condition, and the remaining scores, fair condition."
Source: Apgar 1953 as above, p. 267.
Range: 0–10. Tests 2/3, 7/8, 0, 10, −1, 11. (The AAP/ACOG 2015 bands are NOT built: their wording was not seen.)
