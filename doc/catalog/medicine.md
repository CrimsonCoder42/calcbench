# Catalog: medicine (`med/*`)

Status: BUILT 2026-11-12 (44 functions, 7 packages); every other row is HELD, MOVED, DROPPED or LATER below.
Work folder: [../medicine_2026-11-12/](../medicine_2026-11-12/README.md) · conventions:
[CONVENTIONS.md](../medicine_2026-11-12/CONVENTIONS.md) · decisions: [DECISIONS_DRAFT.md](../medicine_2026-11-12/DECISIONS_DRAFT.md)
· shared rows: [PENDING_SHARED.md](../medicine_2026-11-12/PENDING_SHARED.md).

The input list is vCalc's Medicine names only; no vCalc page was fetched. The citation bar (D-M2) is that a formula ships
only when it was READ in the primary source (abstract, full text, scan or official document) on 2026-11-12. A citation that
resolved but whose formula was seen only in a secondary source is HELD-SOURCE, and is listed under LEFT FOR DEVIN with what
must be read to release it. Results are numbers and the source's own labels; nothing diagnoses, interprets or advises.

Decisions: KEEP (built) · MERGE (into another row) · MOVE (another home) · DROP (not a formula) · HELD (policy: dose or
amount to administer) · HELD-SOURCE (citation resolves, formula not read in the primary) · SOURCE NOT FOUND · LATER.

## 1. The input list, cleaned

| # | Name as given | Decision | Home | Reason |
|---|---|---|---|---|
| 1 | Body Surface Area Calculator (22 fns) (featured) | DROP | — | A calculator. Bundles the BSA formulas: MED-ANTHRO BSA rows (Du Bois, Mosteller, Haycock) |
| 2 | Epworth Sleepiness Scale (featured) | MERGE | `med/scales` | Same as row 'Epworth Sleepiness Scale' below |
| 3 | Metabolic Health Calculator (7 fns) (featured) | DROP | — | A calculator. Bundles BMI, waist and metabolic markers; its formulas are rows of med/anthropometry and med/chemistry |
| 4 | Blood Alcohol Concentration Calculator | DROP | — | A calculator. Its formula is Widmark: MED-PHARM-001 |
| 5 | Body Surface Area (collection) | DROP | — | A collection of the BSA formulas (as row 1) |
| 6 | Dr. Kurz Calculations (collection) | DROP | — | A personal collection; no formula behind the name |
| 7 | Herbal Remedy Converter | DROP | — | Not a formula with a citable source (herbal-preparation conversions); SOURCE NOT FOUND |
| 8 | Alveolar gas equation | HELD-SOURCE | `med/respiratory` | Riley & Cournand 1949 (J Appl Physiol 1:825, PMID 18145478) and West 11th ed: equation not seen (publisher 403, book not opened). LEFT FOR DEVIN |
| 9 | Alveolar gas equation (second listing) | MERGE | `med/respiratory` | Duplicate of the row above |
| 10 | Banana Equivalent Dose | MOVE | `sci/physics` | Radiation dose comparison (⁴⁰K activity), health physics, not a clinical calculator |
| 11 | Bayes' Theorem for Disease Testing | MOVE | `num/prob` | Already catalogued there: NUM-PROB-010 PositivePredictiveValue, -011 NegativePredictiveValue (doc/catalog/prob.md row 9); med calls it |
| 12 | Blood Pressure | HELD-SOURCE | `med/cardio` | 2025 AHA/ACC guideline (Jones DW, Hypertension 82:e212, PMID 40811516) confirmed; its category table not read (403). LEFT FOR DEVIN |
| 13 | Blood Volume in a Human | HELD-SOURCE | `med/fluids` | Nadler 1962 (PMID 21936146): coefficients not seen in the paper (no abstract, full text not reached). LEFT FOR DEVIN |
| 14 | Carbon Dioxide Elimination | HELD-SOURCE | `med/respiratory` | West 11th ed not opened. LEFT FOR DEVIN |
| 15 | Cardiac Efficiency (first listing) | HELD-SOURCE | `med/cardio` | Listed three times; no primary source found for a cardiac-efficiency formula (Guyton ch. 9 not opened). LEFT FOR DEVIN |
| 16 | Cardiac Efficiency (second listing) | MERGE | `med/cardio` | Duplicate of the first Cardiac Efficiency row |
| 17 | Cardiac Efficiency (third listing) | MERGE | `med/cardio` | Duplicate of the first Cardiac Efficiency row |
| 18 | Cardiac Flow | MERGE | `med/cardio` | Flow = pressure difference / resistance: the same relation as SVR (MED-CARDIO SVR row), solved for flow |
| 19 | Cardiac Output | HELD-SOURCE | `med/cardio` | CO = HR × SV and Fick: Guyton & Hall 14th ed ch. 20 not opened. LEFT FOR DEVIN |
| 20 | Conduction Velocity | SOURCE NOT FOUND | — | Nerve conduction velocity: no primary searched to confirmation in this run. LEFT FOR DEVIN |
| 21 | Dewormer Dosage Calculator | HELD | `med/pharmacology` | Outputs an amount to administer (policy). Label: pyrantel pamoate suspension (FDA SPL) — LEFT FOR DEVIN |
| 22 | Diffusion Capacity | HELD-SOURCE | `med/respiratory` | ERS/ATS 2017 (Graham BL, Eur Respir J 49:1600016, PMID 28049168) equation not seen (403). LEFT FOR DEVIN |
| 23 | Diffusion Capacity Components | KEEP | `med/respiratory` | MED-RESP-001 DiffusingCapacityRoughtonForster (Zavorsky 2017 ERS standard; Roughton & Forster 1957) |
| 24 | Diffusion Capacity using CO | MERGE | `med/respiratory` | Into row 'Diffusion Capacity' (HELD-SOURCE) |
| 25 | Diluting Stock Solutions | MOVE | `sci/chemistry` | C₁V₁ = C₂V₂ is solution chemistry; no clinical source needed or found here |
| 26 | Direct Care FTEs | MOVE | `fin/business` | Staffing arithmetic, not a clinical formula |
| 27 | Dosage Amount(Size of Patient) | HELD | `med/pharmacology` | Dose arithmetic (D/H × Q family): held by policy; source to be supplied by Devin. LEFT FOR DEVIN |
| 28 | Dosage Calculation | HELD | `med/pharmacology` | Dose arithmetic (D/H × Q family): held by policy; source to be supplied by Devin. LEFT FOR DEVIN |
| 29 | Dosage Calculation (Amount of Drug on Hand) | HELD | `med/pharmacology` | Dose arithmetic (D/H × Q family): held by policy; source to be supplied by Devin. LEFT FOR DEVIN |
| 30 | Dosage Calculation (Amount to be Administered) | HELD | `med/pharmacology` | Dose arithmetic (D/H × Q family): held by policy; source to be supplied by Devin. LEFT FOR DEVIN |
| 31 | Dosage Calculation (Desired Dose) | HELD | `med/pharmacology` | Dose arithmetic (D/H × Q family): held by policy; source to be supplied by Devin. LEFT FOR DEVIN |
| 32 | Dosage Calculation (Quantity on Hand) | HELD | `med/pharmacology` | Dose arithmetic (D/H × Q family): held by policy; source to be supplied by Devin. LEFT FOR DEVIN |
| 33 | Dosage rate | HELD | `med/pharmacology` | Dose arithmetic (D/H × Q family): held by policy; source to be supplied by Devin. LEFT FOR DEVIN |
| 34 | Epworth Sleepiness Scale | KEEP | `med/scales` | MED-SCALE-001 score, -002 category (Johns 1991; official site bands) |
| 35 | Fibrosis-4 Score | KEEP | `med/hepatic` | MED-HEPAT-001 FIB4 (Sterling 2006) |
| 36 | Flow Rate (Gravity System) | HELD | `med/fluids` | Drip rate arithmetic; textbook source not confirmed (opentextbc.ca 403). LEFT FOR DEVIN |
| 37 | Fractional Excretion of Sodium | KEEP | `med/renal` | MED-RENAL-009 (Espinel 1976) |
| 38 | Hagen-Poiseuille | MOVE | `sci/physics` | Laminar pipe flow, physics; med/cardio does not duplicate it |
| 39 | Heart Stroke Work | HELD-SOURCE | `med/cardio` | Guyton & Hall ch. 9 not opened. LEFT FOR DEVIN |
| 40 | Ideal Body Weight for Children of ages 1-10 | HELD-SOURCE | `med/anthropometry` | APLS manual not seen (formulas only via Luscombe 2007, Seddon 2012); Tinning–Acworth 1–4 y band edge unconfirmed. Built instead for ages 1–17 by height: MED-ANTHRO-006 Traub–Kichen. LEFT FOR DEVIN |
| 41 | Ideal Body Weight for Children under 1 | KEEP | `med/anthropometry` | MED-ANTHRO-007 EstimatedWeightInfantBestGuess (Tinning & Acworth 2007) |
| 42 | IV Dosage Calculation | HELD | `med/pharmacology` | Outputs a dose: held by policy. LEFT FOR DEVIN |
| 43 | Mean Arterial Pressure | HELD-SOURCE | `med/cardio` | No primary establishing the one-third rule seen; the only open study (Mahieu 2007, Artery Res 1:71) reports the rule underestimates MAP. LEFT FOR DEVIN |
| 44 | Minute ventilation | HELD-SOURCE | `med/respiratory` | West 11th ed not opened. LEFT FOR DEVIN |
| 45 | Moens-Korteweg Equation | KEEP | `med/cardio` | MED-CARDIO-003 (Korteweg 1878 eq. 10, read in the scan) |
| 46 | Parkland Replacement Fluids | HELD | `med/fluids` | An amount of fluid to administer: held by policy. Baxter & Shires 1968 (PMID 4973463, formula not seen); current ABA 2024 guideline (Cartotto, J Burn Care Res 45:565, PMID 38051821) gives 2 mL/kg/%TBSA. LEFT FOR DEVIN |
| 47 | Pediatric Acetaminophen Dose | HELD | `med/pharmacology` | A drug dose: held by policy. Label: Children's Acetaminophen Oral Suspension 160 mg/5 mL, DailyMed setid 8231dd91-e9e6-4350-a485-fa95d65a6d25. LEFT FOR DEVIN |
| 48 | Physiologic Dead Space | HELD-SOURCE | `med/respiratory` | Enghoff 1938: no DOI, not seen. LEFT FOR DEVIN |
| 49 | Physiologic Dead Space(Bohr Method) | HELD-SOURCE | `med/respiratory` | Bohr 1891 (doi:10.1111/j.1748-1716.1891.tb00581.x) not read. LEFT FOR DEVIN |
| 50 | Physiologic Dead Space(Tidal Volume) | MERGE | `med/respiratory` | VD = VT × (VD/VT): the caller's multiply of the fraction rows; one function returns the volume |
| 51 | Pill Dosage Calculation | HELD | `med/pharmacology` | Dose arithmetic: held by policy. LEFT FOR DEVIN |
| 52 | Pyrantel Pamoate Dewormer for Puppies and Dogs | HELD | `med/pharmacology` | Veterinary drug dose: held by policy (and veterinary, not med). Label: FDA SPL pyrantel pamoate suspension (5 mL per 10 lb). LEFT FOR DEVIN |
| 53 | Reabsorption Rate [GFR] | HELD-SOURCE | `med/renal` | Guyton & Hall 14th ed ch. 28 not opened. LEFT FOR DEVIN |
| 54 | Respiratory Exchange Ratio | HELD-SOURCE | `med/respiratory` | West 11th ed not opened. LEFT FOR DEVIN |
| 55 | Respiratory Quotient [moles] | MERGE | `med/respiratory` | Same ratio as RER (moles of CO₂ produced / O₂ consumed); one function |
| 56 | Respiratory Quotient [STP Volume] | MERGE | `med/respiratory` | Same ratio in STPD volumes (equal to the mole ratio for ideal gases); one function |
| 57 | Rule of Fives (Burn Percent) | SOURCE NOT FOUND | — | No primary for an infant rule of fives; Livingston & Lee 2000 (J Surg Res 91:106, PMID 10839957) has a rule of fives for adults over 80 kg whose figures as summarised sum to 102%; not built |
| 58 | Shunt Equation | HELD-SOURCE | `med/respiratory` | Berggren 1942: no DOI, not seen. LEFT FOR DEVIN |
| 59 | Starling Equation for volume flow | HELD-SOURCE | `med/cardio` | Home decided: med/cardio. Levick & Michel 2010 (PMID 20200043) equation not seen. LEFT FOR DEVIN |
| 60 | Starling Equation for volume flow (second listing) | MERGE | `med/cardio` | Duplicate |
| 61 | Stroke Work [P, KE, SV] | HELD-SOURCE | `med/cardio` | Guyton & Hall ch. 9 not opened. LEFT FOR DEVIN |
| 62 | Time of Death | LATER | `med/forensic` (not planned) | Henssge nomogram needs its full model; noted for a later forensic slice |
| 63 | Total peripheral Resistance | HELD-SOURCE | `med/cardio` | SVR: Guyton & Hall ch. 14 not opened. LEFT FOR DEVIN |
| 64 | Wallace Rule of Nines | HELD-SOURCE | `med/fluids` | Wallace 1951 (PMID 14805109): table not seen in the paper; Lund–Browder 1944 table seen only on a hospital form. LEFT FOR DEVIN |
| 65 | Weber Fraction | LATER | psychophysics (not planned) | ΔI/I, sensory psychophysics; no home in the taxonomy |
| 66 | Data item: blood | DROP | — | A single reference number with no cited source on the page; blood volume is computed by Nadler (MED-FLUID) |
| 67 | Data item: blood per kilogram of human mass | DROP | — | A single reference number with no cited source on the page; blood volume is computed by Nadler (MED-FLUID) |
| 68 | Data item: BSA mean - child of 10 years | DROP | — | Mean BSA values: no primary table found (Du Bois 1916 gives individual subjects, not population means); SOURCE NOT FOUND |
| 69 | Data item: BSA mean - child of 12-13 years | DROP | — | Mean BSA values: no primary table found (Du Bois 1916 gives individual subjects, not population means); SOURCE NOT FOUND |
| 70 | Data item: BSA mean - child of 2 years | DROP | — | Mean BSA values: no primary table found (Du Bois 1916 gives individual subjects, not population means); SOURCE NOT FOUND |
| 71 | Data item: BSA mean - child of 9 years | DROP | — | Mean BSA values: no primary table found (Du Bois 1916 gives individual subjects, not population means); SOURCE NOT FOUND |
| 72 | Data item: BSA mean - men | DROP | — | Mean BSA values: no primary table found (Du Bois 1916 gives individual subjects, not population means); SOURCE NOT FOUND |
| 73 | Data item: BSA mean - neonate | DROP | — | Mean BSA values: no primary table found (Du Bois 1916 gives individual subjects, not population means); SOURCE NOT FOUND |
| 74 | Data item: BSA mean - women | DROP | — | Mean BSA values: no primary table found (Du Bois 1916 gives individual subjects, not population means); SOURCE NOT FOUND |
| 75 | WikiClip: Blood Velocity Magnitude Information | DROP | — | Explanatory text, not a formula |
| 76 | WikiClip: Heart and Cardiac Calculators | DROP | — | Explanatory text, not a formula |
| 77 | WikiClip: Respiratory Calculator Info | DROP | — | Explanatory text, not a formula |
Counts (77 rows: the "x2"/"x3" listings and every data item and WikiClip are expanded): KEEP 6 · MERGE 10 · MOVE 5 ·
DROP 18 · HELD (policy) 14 · HELD-SOURCE 20 · SOURCE NOT FOUND 2 · LATER 2. The KEEP rows are 6 input rows; with the
standard set and the source-driven additions they come to 44 functions (§3).

## 2. The standard set (anthropometry, vitals-derived, renal), beyond the list

| Item | Decision | Where | Note |
|---|---|---|---|
| Anion gap Na − (Cl + HCO₃) | HELD-SOURCE | med/chemistry | Emmett & Narins 1977 (Medicine 56:38, PMID 401925) abstract does not state the formula |
| Albumin-adjusted anion gap (Figge 1998) | KEEP | MED-CHEM-007, -008 | normal albumin is a caller input (the abstract gives no number) |
| Corrected calcium | KEEP (Payne form) | MED-CHEM-001, -002 | the 0.8 form: SOURCE NOT FOUND (D-M6) |
| Corrected sodium for hyperglycaemia (Katz 1973, Hillier 1999) | HELD-SOURCE | med/chemistry | factors 1.6 and 2.4 seen; the baseline glucose (100 mg/dL) only in a pointer |
| Serum osmolality | KEEP (Dorwart's recommended equation) | MED-CHEM-003, -004 | the 2Na + glucose/18 + BUN/2.8 form is not Dorwart's recommendation; the SI 2Na + glucose + urea form has no primary found |
| eAG from HbA1c (ADAG, Nathan 2008) | KEEP | MED-CHEM-005, -006 | |
| LDL Friedewald 1972 | HELD-SOURCE | med/chemistry | paper's wording and the TG < 400 limit not read (paywalled) |
| LDL Martin–Hopkins 2013 | KEEP | MED-CHEM-009, -010 | 180 cells transcribed twice independently, all agree |
| GFR estimation CKD-EPI 2021 (creatinine; creatinine–cystatin C) | KEEP | MED-RENAL-001…004 | Table 2 AS columns, not ASR-NB |
| MDRD 2006 (IDMS) | HELD-SOURCE | med/renal | coefficients not seen (abstract silent, full text paywalled) |
| Cockcroft–Gault, bedside Schwartz | KEEP | MED-RENAL-005…008 | |
| KDIGO GFR categories | KEEP | MED-RENAL-011 | read in the KDIGO 2024 PDF p. S125 |
| Target heart rate: Tanaka | KEEP | MED-CARDIO-001 | |
| Target heart rate: Karvonen 1957 | HELD-SOURCE | med/cardio | no abstract, no scan found |
| QTc Bazett | HELD-SOURCE | med/cardio | 1920 original and 1997 reprint not read |
| QTc Fridericia | KEEP | MED-CARDIO-002 | open English translation (PMC7328879) |
| BP categories (2025 AHA/ACC) | HELD-SOURCE | med/cardio | citation confirmed (PMID 40811516); table blocked (403) |
| FIB-4; MELD (UNOS form), MELDNa, MELD 3.0 | KEEP | MED-HEPAT-001…004 | |
| APRI (Wai 2003) | HELD-SOURCE | med/hepatic | formula and cut-offs only from pointers |
| Epworth, GCS (+ Rimel bands), Apgar (+ 1953 groups) | KEEP | MED-SCALE-001…006 | Epworth: no item text (licence, D-M8); AAP/ACOG Apgar bands not seen |
| Body fat (Deurenberg 1991), Haycock BSA, adjusted weight (Schwartz 1978), Traub–Kichen, infant weight (Tinning–Acworth), CDC BMI and categories | KEEP | MED-ANTHRO-001…009 | |
| Du Bois BSA, Mosteller BSA, Devine IBW, APLS weight | HELD-SOURCE | med/anthropometry | formulas only from pointers |
| VO₂max estimates | LATER | med/cardio | not researched in this run |
| BMI-for-age z-score (CDC/WHO LMS tables) | LATER | med/anthropometry | official LMS files not embedded in this run (the request's condition) |
| Moens–Korteweg | KEEP | MED-CARDIO-003 | Korteweg 1878 eq. (10), read in the scan |
| Roughton–Forster | KEEP | MED-RESP-001 | Zavorsky 2017 ERS standard (read) and the 1957 paper |

## 3. Built functions (44): one entry each, generated from the code's doc comments (the code is the truth)

### BodyMassIndex      (ID: MED-ANTHRO-001, built: med/anthropometry/body_mass_index.go:58)
Signature: `func BodyMassIndex(weightKg, heightM float64) (float64, error)`
Formula and output: BodyMassIndex returns the body mass index, in kg/m², from weight in kilograms and height in metres: weightKg / heightM².
Inputs and units: weightKg kg; heightM m; result kg/m²
Range: weightKg > 0; heightM > 0 (the source states no limits)
Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range; ErrDomain when the result is beyond float64 (a height so small that its square underflows, or a weight so large that the quotient overflows).
Primary source: Centers for Disease Control and Prevention. BMI Frequently Asked Questions. https://www.cdc.gov/bmi/faq/index.html (last reviewed 2024-06-28; read 2026-11-12): "weight (kg)/[height (m)]2"
Worked example: 70 kg, 1.75 m → 70 / 3.0625 = 22.857142857142858 (HAND-DERIVED)
Best approach: Time O(1); Space O(1), no allocation
Tags: anthropometry, bmi, body-mass-index, metric

### BodyMassIndexImperial      (ID: MED-ANTHRO-002, built: med/anthropometry/body_mass_index_imperial.go:23)
Signature: `func BodyMassIndexImperial(weightLb, heightIn float64) (float64, error)`
Formula and output: BodyMassIndexImperial returns the body mass index, in kg/m², from weight in pounds and height in inches by the CDC's imperial form: weightLb / heightIn² × 703. The factor 703 is the CDC's; the result agrees with BodyMassIndex (MED-ANTHRO-001) on the same person to within 0.01%, not exactly.
Inputs and units: weightLb lb; heightIn in; result kg/m²
Range: weightLb > 0; heightIn > 0 (the source states no limits)
Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range; ErrDomain when the result is beyond float64.
Primary source: Centers for Disease Control and Prevention. BMI Frequently Asked Questions. https://www.cdc.gov/bmi/faq/index.html (last reviewed 2024-06-28; read 2026-11-12): "weight (pounds)/[height (in)]2 x 703"
Worked example: 150 lb, 65 in → 150 / 4225 × 703 = 24.958579881656806 (HAND-DERIVED)
Best approach: Time O(1); Space O(1), no allocation
Tags: anthropometry, bmi, body-mass-index, imperial

### BMICategoryCDCAdult      (ID: MED-ANTHRO-003, built: med/anthropometry/bmi_category_cdc_adult.go:35)
Signature: `func BMICategoryCDCAdult(bmiKgM2 float64) (BMICategory, error)`
Formula and output: BMICategoryCDCAdult returns the CDC adult BMI category for a BMI in kg/m², with the CDC's own label: < 18.5 "Underweight"; 18.5 to < 25 "Healthy Weight"; 25 to < 30 "Overweight"; 30 to < 35 "Class 1 Obesity"; 35 to < 40 "Class 2 Obesity"; ≥ 40 "Class 3 Obesity (Severe Obesity)". These are the CDC's ADULT categories; the function takes no age and does not apply them to children.
Inputs and units: bmiKgM2 kg/m²; result the CDC's category label
Range: bmiKgM2 > 0; the CDC page's adult categories
Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange for bmiKgM2 ≤ 0.
Primary source: Centers for Disease Control and Prevention. Adult BMI Categories. https://www.cdc.gov/bmi/adult-calculator/bmi-categories.html (last reviewed 2024-03-19; read 2026-11-12)
Worked example: 18.4999 → "Underweight"; 18.5 → "Healthy Weight"; 40 → "Class 3 Obesity (Severe Obesity)" (CDC table)
Best approach: Time O(1); Space O(1), no allocation
Tags: anthropometry, bmi, category, cdc, adult

### BodySurfaceAreaHaycock      (ID: MED-ANTHRO-004, built: med/anthropometry/body_surface_area_haycock.go:18)
Signature: `func BodySurfaceAreaHaycock(weightKg, heightCm float64) (float64, error)`
Formula and output: BodySurfaceAreaHaycock returns the body surface area, in m², by the Haycock height-weight formula: weightKg^0.5378 × heightCm^0.3964 × 0.024265. Every finite positive input gives a finite positive result (the exponents sum to less than 1, so the result never overflows; the smallest inputs give about 2e-304).
Inputs and units: weightKg kg; heightCm cm; result m²
Range: weightKg > 0; heightCm > 0 (no hard limits stated; the 81 subjects ran from premature infants to adults, measured SA from less than 0.2 m² to greater than 2.0 m²)
Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range.
Primary source: Haycock GB, Schwartz GJ, Wisotsky DH. Geometric method for measuring body surface area: a height-weight formula validated in infants, children, and adults. J Pediatr. 1978;93(1):62-66. doi:10.1016/s0022-3476(78)80601-5. PMID 650346. Abstract: "SA (m2) = weight (kg)0.5378 X height (cm)0.3964 X 0.024265"
Worked example: 70 kg, 175 cm → 70^0.5378 × 175^0.3964 × 0.024265 (HAND-DERIVED; the test writes the formula out)
Best approach: Time O(1); Space O(1), no allocation
Tags: anthropometry, bsa, body-surface-area, haycock, pediatric

### AdjustedBodyWeight      (ID: MED-ANTHRO-005, built: med/anthropometry/adjusted_body_weight.go:19)
Signature: `func AdjustedBodyWeight(idealWeightKg, actualWeightKg float64) (float64, error)`
Formula and output: AdjustedBodyWeight returns the adjusted body weight, in kg: idealWeightKg + 0.4 × (actualWeightKg − idealWeightKg), the source's "normalized body mass plus 40% of the adipose mass". The 0.4 was derived for the distribution volume of gentamicin and tobramycin (the source's context); the function returns a weight, nothing else. The result lies between the two inputs, so it never overflows.
Inputs and units: idealWeightKg kg; actualWeightKg kg; result kg
Range: idealWeightKg > 0; actualWeightKg > idealWeightKg (the adipose mass is actual − ideal; the source's population: "13 obese and 13 non-obese subjects")
Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range.
Primary source: Schwartz SN, Pazin GJ, Lyon JA, Ho M, Pasculle AW. A controlled investigation of the pharmacokinetics of gentamicin and tobramycin in obese subjects. J Infect Dis. 1978;138(4):499-505. doi:10.1093/infdis/138.4.499. PMID 712111. Abstract
Worked example: ideal 70 kg, actual 120 kg → 70 + 0.4 × 50 = 90 kg (HAND-DERIVED)
Best approach: Time O(1); Space O(1), no allocation
Tags: anthropometry, body-weight, adjusted-body-weight, dosing-weight

### IdealBodyMassTraubKichen      (ID: MED-ANTHRO-006, built: med/anthropometry/ideal_body_mass_traub_kichen.go:22)
Signature: `func IdealBodyMassTraubKichen(heightCm, ageYears float64) (float64, error)`
Formula and output: IdealBodyMassTraubKichen returns the ideal body mass of a child, in kg, by the Traub–Kichen formula: 2.396 × e^(0.01863 × heightCm), fitted to the 50th-percentile weight for height of NCHS data, ages 1–17. ageYears is used only to enforce the source's age range; the formula itself takes height alone.
Inputs and units: heightCm cm; ageYears years (decimal); result kg
Range: 1 ≤ ageYears ≤ 17 (the source's ages); heightCm > 0 (no height limits stated)
Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range; ErrDomain when the result is beyond float64 (heightCm above about 38,100).
Primary source: Traub SL, Kichen L. Estimating ideal body mass in children. Am J Hosp Pharm. 1983;40(1):107-110. PMID 6823980. Abstract: "IBM(kg) = 2.396e0.01863(ht)", ages 1–17, the 50th-percentile weight for height from NCHS data
Worked example: 100 cm → 2.396 × e^1.863 (HAND-DERIVED; the test writes the formula out)
Best approach: Time O(1); Space O(1), no allocation
Tags: anthropometry, ideal-body-weight, pediatric, traub-kichen

### EstimatedWeightInfantBestGuess      (ID: MED-ANTHRO-007, built: med/anthropometry/estimated_weight_infant_best_guess.go:18)
Signature: `func EstimatedWeightInfantBestGuess(ageMonths float64) (float64, error)`
Formula and output: EstimatedWeightInfantBestGuess returns the estimated weight of an infant under 12 months, in kg, by the Best Guess formula for that band: (ageMonths + 9) / 2. The paper's other bands (1–4 years, 5–14 years) are not built here.
Inputs and units: ageMonths months (decimal); result kg
Range: 0 ≤ ageMonths < 12 (the "<12 months" band; the study had 70,181 emergency department presentations aged 0–14 years)
Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range.
Primary source: Tinning K, Acworth J. Make your Best Guess: an updated method for paediatric weight estimation in emergencies. Emerg Med Australas. 2007;19(6):528-534. doi:10.1111/j.1742-6723.2007.01026.x. PMID 18021105. Abstract, the "<12 months" formula
Worked example: 6 months → (6 + 9) / 2 = 7.5 kg; 0 months → 4.5 kg (HAND-DERIVED)
Best approach: Time O(1); Space O(1), no allocation
Tags: anthropometry, weight-estimation, pediatric, infant, best-guess

### BodyFatPercentDeurenbergAdult      (ID: MED-ANTHRO-008, built: med/anthropometry/body_fat_percent_deurenberg.go:35)
Signature: `func BodyFatPercentDeurenbergAdult(bmiKgM2, ageYears float64, male bool) (float64, error)`
Formula and output: BodyFatPercentDeurenbergAdult returns the body-fat percentage predicted from BMI, age and sex by the Deurenberg adult formula: 1.20 × BMI + 0.23 × age − 10.8 × sex − 5.4, with sex = 1 for male and 0 for female (the source's coding). The source reports R² 0.79 and SEE 4.1% for adults.
Inputs and units: bmiKgM2 kg/m²; ageYears years (decimal); male the source's sex term; result percent body fat
Range: 15 < ageYears ≤ 83; 13.9 ≤ bmiKgM2 ≤ 40.9 (the study's 1229 subjects, ages 7–83, BMI 13.9–40.9; the adult formula above age 15)
Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range.
Primary source: Deurenberg P, Weststrate JA, Seidell JC. Body mass index as a measure of body fatness: age- and sex-specific prediction formulas. Br J Nutr. 1991;65(2):105-114. doi:10.1079/bjn19910073. PMID 2043597. Abstract (adults, R² 0.79, SEE 4.1%)
Worked example: male, 40 y, BMI 25 → 30 + 9.2 − 10.8 − 5.4 = 23.0% (HAND-DERIVED)
Best approach: Time O(1); Space O(1), no allocation
Tags: anthropometry, body-fat, bmi, deurenberg, adult

### BodyFatPercentDeurenbergChild      (ID: MED-ANTHRO-009, built: med/anthropometry/body_fat_percent_deurenberg.go:64)
Signature: `func BodyFatPercentDeurenbergChild(bmiKgM2, ageYears float64, male bool) (float64, error)`
Formula and output: BodyFatPercentDeurenbergChild returns the body-fat percentage predicted from BMI, age and sex by the Deurenberg child formula: 1.51 × BMI − 0.70 × age − 3.6 × sex + 1.4, with sex = 1 for male and 0 for female (the source's coding). The source reports R² 0.38 and SEE 4.4% for children aged 15 and under.
Inputs and units: bmiKgM2 kg/m²; ageYears years (decimal); male the source's sex term; result percent body fat
Range: 7 ≤ ageYears ≤ 15; 13.9 ≤ bmiKgM2 ≤ 40.9 (the study's subjects, ages 7–83, BMI 13.9–40.9; the child formula to age 15)
Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range.
Primary source: Deurenberg P, Weststrate JA, Seidell JC. Body mass index as a measure of body fatness: age- and sex-specific prediction formulas. Br J Nutr. 1991;65(2):105-114. doi:10.1079/bjn19910073. PMID 2043597. Abstract (children ≤ 15 y, R² 0.38, SEE 4.4%)
Worked example: female, 10 y, BMI 18 → 27.18 − 7 + 1.4 = 21.58% (HAND-DERIVED)
Best approach: Time O(1); Space O(1), no allocation
Tags: anthropometry, body-fat, bmi, deurenberg, pediatric

### MaxHeartRateTanaka      (ID: MED-CARDIO-001, built: med/cardio/max_heart_rate_tanaka.go:56)
Signature: `func MaxHeartRateTanaka(ageYears float64) (float64, error)`
Formula and output: MaxHeartRateTanaka returns the age-predicted maximal heart rate, in beats/min, by Tanaka's equation HRmax = 208 − 0.7 × age.
Inputs and units: ageYears years (decimal allowed); result beats/min
Range: ageYears 18–81 (the laboratory study's 514 subjects, ages 18 to 81; the meta-analysis pooled 18,712 adults in 351 studies)
Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range.
Primary source: Tanaka H, Monahan KD, Seals DR. Age-predicted maximal heart rate revisited. J Am Coll Cardiol. 2001;37(1):153-156. doi:10.1016/S0735-1097(00)01054-8. PMID 11153730. Abstract: "HRmax was strongly related to age (r = −0.90), using the equation of 208 − 0.7 × age"
Worked example: Tanaka 2001 (full text): at 70 years, 208 − 0.7 × 70 = 159 beats/min
Best approach: Time O(1); Space O(1), no allocation
Tags: cardio, heart-rate, maximal-heart-rate, exercise, tanaka

### QTcFridericia      (ID: MED-CARDIO-002, built: med/cardio/qtc_fridericia.go:20)
Signature: `func QTcFridericia(qtMs, rrSeconds, ageYears float64) (float64, error)`
Formula and output: QTcFridericia returns the heart-rate-corrected QT interval, in milliseconds, by Fridericia's cube-root relation QTc = QT / RR^(1/3), RR in seconds: the QT interval scaled to a pulse period of 1 s. Fridericia found the duration of electrical systole proportional to the cube root of the pulse period; dividing by RR^(1/3) with RR in seconds is that proportionality normalised to RR = 1 s.
Inputs and units: qtMs milliseconds; rrSeconds seconds; ageYears years; result milliseconds
Range: rrSeconds 60/135 to 60/51 s (the paper's pulse frequencies, 51 to 135 per minute); ageYears 2–87 (the paper's age range; 50 persons); qtMs > 0. ageYears only enforces the Range
Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range.
Primary source: Fridericia LS. The duration of systole in an electrocardiogram in normal humans and in patients with heart disease (English translation of Acta Med Scand. 1920;53:469-486). Ann Noninvasive Electrocardiol. 2003;8(4):343-351. doi:10.1046/j.1542-474X.2003.08413.x. PMID 14516292. PMC7328879, "Relationship of the Pulse Frequency and the Duration of the Ventricular Electrocardiogram": "the duration of electrical systole in an electrocardiogram is proportional to the cube root of the duration of the pulse period"
Primary source: Fridericia LS. Die Systolendauer im Elektrokardiogramm bei normalen Menschen und bei Herzkranken. Acta Med Scand. 1920;53(1):469-486. doi:10.1111/j.0954-6820.1920.tb18266.x
Worked example: QT 400 ms at RR 0.8 s → 400 / 0.8^(1/3) = 430.887… ms (hand-derived)
Best approach: Time O(1); Space O(1), no allocation
Tags: cardio, ecg, qt, qtc, fridericia

### PulseWaveVelocityMoensKorteweg      (ID: MED-CARDIO-003, built: med/cardio/pulse_wave_velocity_moens_korteweg.go:25)
Signature: `func PulseWaveVelocityMoensKorteweg(wallThicknessM, elasticModulusPa, densityKgM3, radiusM float64) (float64, error)`
Formula and output: PulseWaveVelocityMoensKorteweg returns the propagation speed of a pressure wave, in m/s, in a thin-walled elastic tube filled with an incompressible liquid, by Korteweg's eq. (10), β = √(a₁E₁ / (2ρR₁)), with a₁ the wall thickness, E₁ the elastic modulus of the wall, ρ the liquid's density and R₁ the tube radius. Korteweg states the conditions (pp. 530–531): the wavelength is very long compared with the tube diameter and the kinetic energy of the transverse motion is negligible against the longitudinal one.
Inputs and units: SI throughout: wallThicknessM m; elasticModulusPa Pa; densityKgM3 kg/m³; radiusM m; result m/s
Range: every input > 0; the physics has no population (Korteweg's conditions above apply)
Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange for a non-positive input; ErrDomain when the result overflows or underflows float64.
Primary source: Korteweg DJ. Ueber die Fortpflanzungsgeschwindigkeit des Schalles in elastischen Röhren. Ann Phys Chem. 1878;241(12):525-542. doi:10.1002/andp.18782411206. Eq. (10), p. 530; symbols p. 526 (scan read 2026-11-12, zenodo.org/record/1423762)
Primary source: Moens AI. Die Pulscurve. Leiden: E. J. Brill; 1878 (cited by Korteweg p. 530, footnote 5, as the experimental counterpart)
Worked example: a₁ 0.001 m, E₁ 4×10⁵ Pa, ρ 1050 kg/m³, R₁ 0.01 m → √(400/21) = 4.36435… m/s (hand-derived)
Best approach: Time O(1); Space O(1), no allocation
Tags: cardio, pulse-wave-velocity, arterial-stiffness, moens-korteweg, hemodynamics

### CorrectedCalciumPayne      (ID: MED-CHEM-001, built: med/chemistry/corrected_calcium_payne.go:65)
Signature: `func CorrectedCalciumPayne(calciumMgDL, albuminGDL float64) (float64, error)`
Formula and output: CorrectedCalciumPayne returns the albumin-adjusted serum calcium, in mg/dL, by Payne's formula: calcium − albumin + 4.0 (calcium in mg/100 mL, albumin in g/100 mL). The source states no limits; inputs far from those it studied can give a result ≤ 0, which is returned as computed.
Inputs and units: calciumMgDL mg/dL (= mg/100 mL); albuminGDL g/dL (= g/100 mL); result mg/dL
Range: calciumMgDL > 0; albuminGDL > 0 (population: 200 consecutive liver-function-test specimens; limits not stated by the source)
Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range.
Primary source: Payne RB, Little AJ, Williams RB, Milner JR. Interpretation of serum calcium in patients with abnormal serum proteins. Br Med J. 1973;4(5893):643-646. doi:10.1136/bmj.4.5893.643. PMID 4758544. Abstract: "Adjusted calcium = calcium − albumin + 4.0, where calcium is in mg/100 ml and albumin in g/100 ml."
Worked example: calcium 8.0, albumin 2.0 → 8.0 − 2.0 + 4.0 = 10.0 mg/dL (hand-derived)
Best approach: Time O(1); Space O(1), no allocation
Tags: chemistry, calcium, albumin, corrected-calcium, payne

### CorrectedCalciumPayneSI      (ID: MED-CHEM-002, built: med/chemistry/corrected_calcium_payne.go:86)
Signature: `func CorrectedCalciumPayneSI(calciumMmolL, albuminGL float64) (float64, error)`
Formula and output: CorrectedCalciumPayneSI returns the albumin-adjusted serum calcium, in mmol/L, by Payne's formula (MED-CHEM-001): calcium mmol/L ÷ 0.2495 gives mg/dL and albumin g/L ÷ 10 gives g/dL; the mg/dL result is multiplied by 0.2495 back to mmol/L.
Inputs and units: calciumMmolL mmol/L; albuminGL g/L; result mmol/L
Range: calciumMmolL > 0; albuminGL > 0 (limits not stated by the source)
Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range; ErrDomain when the result is beyond float64.
Primary source: Payne RB, Little AJ, Williams RB, Milner JR. Interpretation of serum calcium in patients with abnormal serum proteins. Br Med J. 1973;4(5893):643-646. doi:10.1136/bmj.4.5893.643. PMID 4758544. Abstract: "Adjusted calcium = calcium − albumin + 4.0, where calcium is in mg/100 ml and albumin in g/100 ml."
Primary source: Kidney Disease: Improving Global Outcomes (KDIGO) CKD Work Group. KDIGO 2024 Clinical Practice Guideline for the Evaluation and Management of Chronic Kidney Disease. Kidney Int. 2024;105(4S):S117–S314. doi:10.1016/j.kint.2023.10.018. p. S126, "Conversion factors of conventional units to SI units" (calcium 0.2495)
Worked example: calcium 1.996 mmol/L (= 8.0 mg/dL), albumin 20 g/L (= 2.0 g/dL) → 10.0 × 0.2495 = 2.495 mmol/L (hand-derived)
Best approach: Time O(1); Space O(1), no allocation
Tags: chemistry, calcium, albumin, corrected-calcium, payne, si

### SerumOsmolalityDorwart      (ID: MED-CHEM-003, built: med/chemistry/serum_osmolality_dorwart.go:15)
Signature: `func SerumOsmolalityDorwart(sodiumMmolL, glucoseMgDL, ureaNitrogenMgDL float64) (float64, error)`
Formula and output: SerumOsmolalityDorwart returns the calculated serum osmolality, in mOsm/kg, by the Dorwart–Chalmers equation: 1.86 × Na + glucose/18 + BUN/2.8 + 9.
Inputs and units: sodiumMmolL mmol/L; glucoseMgDL mg/dL; ureaNitrogenMgDL mg/dL (blood urea nitrogen); result mOsm/kg
Range: every input > 0 (715 hospital sera; limits not stated by the source)
Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range; ErrDomain when the result is beyond float64.
Primary source: Dorwart WV, Chalmers L. Comparison of methods for calculating serum osmolality from chemical concentrations, and the prognostic value of such calculations. Clin Chem. 1975;21(2):190-194. doi:10.1093/clinchem/21.2.190. PMID 1112025. Abstract: "Osmolality = 1.86 Na + (Glu/18) + (BUN/2.8) + 9."
Worked example: Na 140, glucose 90, BUN 14 → 260.4 + 5 + 5 + 9 = 279.4 mOsm/kg (hand-derived)
Best approach: Time O(1); Space O(1), no allocation
Tags: chemistry, osmolality, sodium, glucose, bun

### SerumOsmolalityDorwartSI      (ID: MED-CHEM-004, built: med/chemistry/serum_osmolality_dorwart.go:40)
Signature: `func SerumOsmolalityDorwartSI(sodiumMmolL, glucoseMmolL, ureaMmolL float64) (float64, error)`
Formula and output: SerumOsmolalityDorwartSI returns the calculated serum osmolality, in mOsm/kg, from glucose and urea in mmol/L: 1.86 × Na + glucose + urea + 9. This is the Dorwart–Chalmers equation (MED-CHEM-003) with the paper's own divisors 18 and 2.8 read as its mg/dL → mmol/L conversions for glucose and urea nitrogen; it is exact algebra on the paper's equation, with no other conversion factor.
Inputs and units: sodiumMmolL mmol/L; glucoseMmolL mmol/L; ureaMmolL mmol/L; result mOsm/kg
Range: every input > 0 (limits not stated by the source)
Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range; ErrDomain when the result is beyond float64.
Primary source: Dorwart WV, Chalmers L. Comparison of methods for calculating serum osmolality from chemical concentrations, and the prognostic value of such calculations. Clin Chem. 1975;21(2):190-194. doi:10.1093/clinchem/21.2.190. PMID 1112025. Abstract: "Osmolality = 1.86 Na + (Glu/18) + (BUN/2.8) + 9."
Worked example: Na 140, glucose 5 mmol/L (= 90 mg/dL), urea 5 mmol/L (= BUN 14 mg/dL) → 260.4 + 5 + 5 + 9 = 279.4 mOsm/kg (hand-derived)
Best approach: Time O(1); Space O(1), no allocation
Tags: chemistry, osmolality, sodium, glucose, urea, si

### EstimatedAverageGlucose      (ID: MED-CHEM-005, built: med/chemistry/estimated_average_glucose.go:19)
Signature: `func EstimatedAverageGlucose(a1cPercent float64) (float64, error)`
Formula and output: EstimatedAverageGlucose returns the estimated average glucose (eAG), in mg/dL, from HbA1c in percent by the ADAG regression: 28.7 × A1C − 46.7. It is evaluated as (287 × A1C − 467) / 10, the same equation with the decimal coefficients scaled to integers, so a whole-percent A1C gives the correctly rounded value (6 % → 125.5 exactly, where 28.7 × 6 − 46.7 in float64 gives 125.49999999999999).
Inputs and units: a1cPercent percent (NGSP); result mg/dL
Range: a1cPercent 4–12 (the study's lowest reported A1C band starts at 4 %; its table tabulates 5–12; ASSUMPTION A-M4)
Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range.
Primary source: Nathan DM, Kuenen J, Borg R, Zheng H, Schoenfeld D, Heine RJ; A1c-Derived Average Glucose Study Group. Translating the A1C assay into estimated average glucose values. Diabetes Care. 2008;31(8):1473-1478. doi:10.2337/dc08-0545. PMID 18540046. Results: "AG(mg/dl) = 28.7 × A1C − 46.7" (full text PMC2742903)
Worked example: A1C 7 → 28.7 × 7 − 46.7 = 154.2 mg/dL; the paper's table prints 154
Best approach: Time O(1); Space O(1), no allocation
Tags: chemistry, glucose, a1c, hba1c, eag, diabetes

### EstimatedAverageGlucoseMmolL      (ID: MED-CHEM-006, built: med/chemistry/estimated_average_glucose.go:39)
Signature: `func EstimatedAverageGlucoseMmolL(a1cPercent float64) (float64, error)`
Formula and output: EstimatedAverageGlucoseMmolL returns the estimated average glucose (eAG), in mmol/L, from HbA1c in percent by the ADAG regression in SI units: 1.5944 × A1C − 2.594. It is evaluated as (15944 × A1C − 25940) / 10000, the same equation with the decimal coefficients scaled to integers.
Inputs and units: a1cPercent percent (NGSP); result mmol/L
Range: a1cPercent 4–12 (ASSUMPTION A-M4)
Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range.
Primary source: Nathan DM, Kuenen J, Borg R, Zheng H, Schoenfeld D, Heine RJ; A1c-Derived Average Glucose Study Group. Translating the A1C assay into estimated average glucose values. Diabetes Care. 2008;31(8):1473-1478. doi:10.2337/dc08-0545. PMID 18540046. Results, the mmol/L form 1.5944 × A1C − 2.594 (full text PMC2742903)
Worked example: A1C 7 → 1.5944 × 7 − 2.594 = 8.5668 mmol/L; the paper's table prints 8.6
Best approach: Time O(1); Space O(1), no allocation
Tags: chemistry, glucose, a1c, hba1c, eag, diabetes, si

### AdjustedAnionGapFigge      (ID: MED-CHEM-007, built: med/chemistry/adjusted_anion_gap_figge.go:16)
Signature: `func AdjustedAnionGapFigge(observedAnionGapMmolL, albuminGL, normalAlbuminGL float64) (float64, error)`
Formula and output: AdjustedAnionGapFigge returns the albumin-adjusted anion gap, in mmol/L, by Figge's correction: observed anion gap + 0.25 × (normal albumin − observed albumin), albumin in g/L. "Normal albumin" is the source's own term for the reference albumin; its abstract gives no number for it, so the caller passes it.
Inputs and units: observedAnionGapMmolL mmol/L; albuminGL g/L; normalAlbuminGL g/L; result mmol/L
Range: albuminGL > 0; normalAlbuminGL > 0; observedAnionGapMmolL any finite value (limits not stated by the source)
Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range; ErrDomain when the result is beyond float64.
Primary source: Figge J, Jabor A, Kazda A, Fencl V. Anion gap and hypoalbuminemia. Crit Care Med. 1998;26(11):1807-1810. doi:10.1097/00003246-199811000-00019. PMID 9824071. Abstract
Worked example: anion gap 12, albumin 20 g/L, normal albumin 44 g/L → 12 + 0.25 × 24 = 18 mmol/L (hand-derived)
Best approach: Time O(1); Space O(1), no allocation
Tags: chemistry, anion-gap, albumin, acid-base, figge

### AdjustedAnionGapFiggeGDL      (ID: MED-CHEM-008, built: med/chemistry/adjusted_anion_gap_figge.go:40)
Signature: `func AdjustedAnionGapFiggeGDL(observedAnionGapMmolL, albuminGDL, normalAlbuminGDL float64) (float64, error)`
Formula and output: AdjustedAnionGapFiggeGDL returns the albumin-adjusted anion gap, in mmol/L, from albumin in g/dL: observed anion gap + 2.5 × (normal albumin − observed albumin), the source's factor for g/dL ("if given in g/dL, the factor is 2.5"). "Normal albumin" is the source's own term; the caller passes it.
Inputs and units: observedAnionGapMmolL mmol/L; albuminGDL g/dL; normalAlbuminGDL g/dL; result mmol/L
Range: albuminGDL > 0; normalAlbuminGDL > 0; observedAnionGapMmolL any finite value (limits not stated by the source)
Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range; ErrDomain when the result is beyond float64.
Primary source: Figge J, Jabor A, Kazda A, Fencl V. Anion gap and hypoalbuminemia. Crit Care Med. 1998;26(11):1807-1810. doi:10.1097/00003246-199811000-00019. PMID 9824071. Abstract
Worked example: anion gap 12, albumin 2.0 g/dL, normal albumin 4.4 g/dL → 12 + 2.5 × 2.4 = 18 mmol/L (hand-derived)
Best approach: Time O(1); Space O(1), no allocation
Tags: chemistry, anion-gap, albumin, acid-base, figge

### LDLCholesterolMartinHopkins      (ID: MED-CHEM-009, built: med/chemistry/ldl_cholesterol_martin_hopkins.go:76)
Signature: `func LDLCholesterolMartinHopkins(totalCholesterolMgDL, hdlCholesterolMgDL, triglyceridesMgDL float64) (float64, error)`
Formula and output: LDLCholesterolMartinHopkins returns the estimated LDL cholesterol, in mg/dL, by the Martin–Hopkins method: non-HDL-C − TG / factor, with non-HDL-C = total cholesterol − HDL-C and the factor read from the 180-cell table of Figure 2 by triglyceride stratum (row) and non-HDL-C stratum (column). A value belongs to the first stratum whose printed upper bound it does not exceed: the integer strata are read as half-open at the top, so TG 49.5 is in "50–56" and non-HDL-C 99.5 is in "100–129" (ASSUMPTION A-M5). High TG with low non-HDL-C can give a result ≤ 0, which is returned as computed.
Inputs and units: totalCholesterolMgDL, hdlCholesterolMgDL, triglyceridesMgDL mg/dL; result mg/dL
Range: triglyceridesMgDL 7–13975 (the table's span); totalCholesterolMgDL > 0; hdlCholesterolMgDL > 0; non-HDL-C = total − HDL > 0
Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range.
Primary source: Martin SS, Blaha MJ, Elshazly MB, et al. Comparison of a novel method vs the Friedewald equation for estimating low-density lipoprotein cholesterol levels from the standard lipid profile. JAMA. 2013;310(19):2061-2068. doi:10.1001/jama.2013.280532. PMID 24240933. Figure 2 (180-cell table; PMC4226221, transcribed twice independently 2026-11-12, all cells agree)
Worked example: TC 200, HDL 50, TG 150 → non-HDL 150, factor 5.7 (TG 147–154, non-HDL 130–159) → 150 − 150/5.7 = 123.68421052631578 mg/dL (hand-derived)
Best approach: Time O(1) (a scan of at most 30 + 6 fixed strata); Space O(1), no allocation
Tags: chemistry, lipids, ldl, cholesterol, triglycerides, martin-hopkins

### LDLCholesterolMartinHopkinsSI      (ID: MED-CHEM-010, built: med/chemistry/ldl_cholesterol_martin_hopkins.go:116)
Signature: `func LDLCholesterolMartinHopkinsSI(totalCholesterolMmolL, hdlCholesterolMmolL, triglyceridesMmolL float64) (float64, error)`
Formula and output: LDLCholesterolMartinHopkinsSI returns the estimated LDL cholesterol, in mmol/L, by the Martin–Hopkins method (MED-CHEM-009) with the paper's own Figure 2 caption factors: cholesterol mg/dL = mmol/L ÷ 0.0259, triglycerides mg/dL = mmol/L ÷ 0.0113; the mg/dL result is multiplied by 0.0259 back to mmol/L. The strata are applied to the converted mg/dL values.
Inputs and units: totalCholesterolMmolL, hdlCholesterolMmolL, triglyceridesMmolL mmol/L; result mmol/L
Range: triglyceridesMmolL ÷ 0.0113 within 7–13975 mg/dL; totalCholesterolMmolL > 0; hdlCholesterolMmolL > 0; total − HDL > 0
Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range; ErrDomain when the result is beyond float64.
Primary source: Martin SS, Blaha MJ, Elshazly MB, et al. Comparison of a novel method vs the Friedewald equation for estimating low-density lipoprotein cholesterol levels from the standard lipid profile. JAMA. 2013;310(19):2061-2068. doi:10.1001/jama.2013.280532. PMID 24240933. Figure 2 (180-cell table; PMC4226221, transcribed twice independently 2026-11-12, all cells agree); Figure 2 caption conversion factors 0.0259 and 0.0113
Worked example: TC 5.18, HDL 1.295, TG 1.695 mmol/L (= 200, 50, 150 mg/dL) → (150 − 150/5.7) × 0.0259 = 3.2034210526315787 mmol/L (hand-derived)
Best approach: Time O(1) (a scan of at most 30 + 6 fixed strata); Space O(1), no allocation
Tags: chemistry, lipids, ldl, cholesterol, triglycerides, martin-hopkins, si

### FIB4      (ID: MED-HEPAT-001, built: med/hepatic/fib4.go:50)
Signature: `func FIB4(ageYears, astUL, altUL, plateletsPer10e9L float64) (float64, error)`
Formula and output: FIB4 returns the FIB-4 index, dimensionless: (age × AST) / (platelets × √ALT). No category is returned; the source's abstract states its two cutoffs for its own population (832 HIV/HCV-coinfected adults, APRICOT), which are not applied here.
Inputs and units: ageYears years (decimal); astUL U/L; altUL U/L; plateletsPer10e9L 10^9/L; result dimensionless
Range: every input > 0 (population: 832 HIV/HCV-coinfected adults (APRICOT); input limits not stated in the source)
Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range; ErrDomain when the result is beyond float64.
Primary source: Sterling RK, Lissen E, Clumeck N, et al. Development of a simple noninvasive index to predict significant fibrosis in patients with HIV/HCV coinfection. Hepatology. 2006;43(6):1317-1325. doi:10.1002/hep.21178. PMID 16729309. Abstract formula "age ([yr] x AST [U/L]) / ((PLT [10^9/L]) x (ALT [U/L])^(1/2))"
Worked example: age 50, AST 40, ALT 25, platelets 200 → 50 × 40 / (200 × 5) = 2.0 (hand-derived)
Best approach: Time O(1); Space O(1), no allocation
Tags: hepatic, liver, fibrosis, fib-4, ast, alt, platelets

### MELD      (ID: MED-HEPAT-002, built: med/hepatic/meld.go:19)
Signature: `func MELD(bilirubinMgDL, inrRatio, creatinineMgDL float64, dialysisTwiceInPastWeek bool) (float64, error)`
Formula and output: MELD returns the MELD score, dimensionless and not rounded: 9.57 × ln(creatinine) + 3.78 × ln(bilirubin) + 11.20 × ln(INR) + 6.43. Values below 1 are set to 1; creatinine is capped at 4.0 mg/dL, and set to 4.0 when dialysisTwiceInPastWeek. The quoted form states no rounding, so none is applied.
Inputs and units: bilirubinMgDL mg/dL; inrRatio dimensionless; creatinineMgDL mg/dL; dialysisTwiceInPastWeek true when dialysed at least twice in the past week; result dimensionless score
Range: every numeric input > 0 (limits not stated by the source; the bounds above are its rules, applied as clamps)
Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range.
Primary source: Kim WR, Mannalithara A, Heimbach JK, et al. MELD 3.0: the Model for End-Stage Liver Disease updated for the modern era. Gastroenterology. 2021;161(6):1887-1895.e4. doi:10.1053/j.gastro.2021.08.050. PMID 34481845. Methods (full text PMC8608337), the MELD formula as quoted: "MELD = 9.57 × loge(creatinine) + 3.78 × loge(bilirubin) + 11.20 × loge(INR) + 6.43"
Primary source: Kamath PS, Wiesner RH, Malinchoc M, et al. A model to predict survival in patients with end-stage liver disease. Hepatology. 2001;33(2):464-470. doi:10.1053/jhep.2001.22172. PMID 11172350 (the original model; its own coefficient form not read)
Worked example: bilirubin 1, INR 1, creatinine 1 → 6.43 (hand-derived: every log is 0)
Best approach: Time O(1); Space O(1), no allocation
Tags: hepatic, liver, meld, bilirubin, inr, creatinine

### MELDNa      (ID: MED-HEPAT-003, built: med/hepatic/meld_na.go:20)
Signature: `func MELDNa(bilirubinMgDL, inrRatio, creatinineMgDL, sodiumMmolL float64, dialysisTwiceInPastWeek bool) (float64, error)`
Formula and output: MELDNa returns the MELD-Na score, dimensionless and not rounded: MELD + 1.32 × (137 − Na) − 0.033 × MELD × (137 − Na), with MELD from MED-HEPAT-002 (its bounds and dialysis rule) and sodium bounded to 125–137 mmol/L (values outside are clamped, the source's rule). The sodium term is applied at every MELD value: the formula as quoted in the source states no MELD threshold, so none is applied here (the OPTN allocation policy's condition on MELD is not part of this function; ASSUMPTION A-M9).
Inputs and units: bilirubinMgDL mg/dL; inrRatio dimensionless; creatinineMgDL mg/dL; sodiumMmolL mmol/L; dialysisTwiceInPastWeek true when dialysed at least twice in the past week; result dimensionless score
Range: every numeric input > 0 (limits not stated by the source; the bounds above are its rules, applied as clamps)
Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range.
Primary source: Kim WR, Mannalithara A, Heimbach JK, et al. MELD 3.0: the Model for End-Stage Liver Disease updated for the modern era. Gastroenterology. 2021;161(6):1887-1895.e4. doi:10.1053/j.gastro.2021.08.050. PMID 34481845. Methods (full text PMC8608337), the MELDNa formula as quoted there
Worked example: bilirubin 1, INR 1, creatinine 1, Na 125 → 6.43 + 1.32 × 12 − 0.033 × 6.43 × 12 = 19.72372 (hand-derived)
Best approach: Time O(1); Space O(1), no allocation
Tags: hepatic, liver, meld, meld-na, sodium, bilirubin, inr, creatinine

### MELD3      (ID: MED-HEPAT-004, built: med/hepatic/meld3.go:22)
Signature: `func MELD3(female bool, bilirubinMgDL, sodiumMmolL, inrRatio, creatinineMgDL, albuminGDL float64) (int, error)`
Formula and output: MELD3 returns the MELD 3.0 score, an integer: 1.33 (if female) + 4.56 × ln(bilirubin) + 0.82 × (137 − Na) − 0.24 × (137 − Na) × ln(bilirubin) + 9.09 × ln(INR) + 11.14 × ln(creatinine) + 1.85 × (3.5 − albumin) − 1.83 × (3.5 − albumin) × ln(creatinine) + 6, rounded to the nearest integer (the source's rule; it states no tie rule, so math.Round, half away from zero, is used: ASSUMPTION A-M6). The source's bounds are applied as clamps: bilirubin, INR and creatinine below 1.0 set to 1.0; creatinine at most 3.0; sodium 125–137; albumin 1.5–3.5. There is no dialysis input: the equation and bounds as quoted in the source state no dialysis rule (ASSUMPTION A-M10).
Inputs and units: female the source's sex term; bilirubinMgDL mg/dL; sodiumMmolL mmol/L; inrRatio dimensionless; creatinineMgDL mg/dL; albuminGDL g/dL; result dimensionless integer score
Range: every numeric input > 0 (limits not stated by the source; the bounds above are its rules, applied as clamps)
Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range.
Primary source: Kim WR, Mannalithara A, Heimbach JK, et al. MELD 3.0: the Model for End-Stage Liver Disease updated for the modern era. Gastroenterology. 2021;161(6):1887-1895.e4. doi:10.1053/j.gastro.2021.08.050. PMID 34481845. Methods/Results, the MELD 3.0 equation (full text PMC8608337)
Worked example: male, bilirubin 1, Na 137, INR 1, creatinine 1, albumin 3.5 → 6 (hand-derived: every other term is 0)
Best approach: Time O(1); Space O(1), no allocation
Tags: hepatic, liver, meld, meld-3, sodium, albumin, bilirubin, inr, creatinine

### EGFRCKDEPI2021      (ID: MED-RENAL-001, built: med/renal/egfr_ckd_epi_2021.go:71)
Signature: `func EGFRCKDEPI2021(creatinineMgDL, ageYears float64, female bool) (float64, error)`
Formula and output: EGFRCKDEPI2021 returns the estimated glomerular filtration rate, in mL/min/1.73 m², from serum creatinine in mg/dL by the 2021 CKD-EPI creatinine equation (race-free): 142 × min(Scr/κ, 1)^α × max(Scr/κ, 1)^−1.200 × 0.9938^age × 1.012 [if female], with κ = 0.7 (female) or 0.9 (male) and α = −0.241 (female) or −0.302 (male).
Inputs and units: creatinineMgDL mg/dL; ageYears years (decimal, completed age as the caller has it); female the source's sex coefficient; result mL/min/1.73 m²
Range: ageYears ≥ 18 ("All the participants were 18 years of age or older"); creatinineMgDL > 0; no upper limits stated by the source
Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range (every in-range input gives a finite result).
Primary source: Inker LA, Eneanya ND, Coresh J, et al. New creatinine- and cystatin C-based equations to estimate GFR without race. N Engl J Med. 2021;385(19):1737-1749. doi:10.1056/NEJMoa2102953. PMID 34554658. Table 2, eGFRcr(AS) column (full text PMC8822996)
Worked example: female, 50 y, Scr 0.7 → 142 × 0.9938^50 × 1.012 = 105.29760114914401 (hand-derived from Table 2)
Best approach: Time O(1); Space O(1), no allocation
Tags: renal, egfr, ckd-epi, creatinine

### EGFRCKDEPI2021SI      (ID: MED-RENAL-002, built: med/renal/egfr_ckd_epi_2021.go:100)
Signature: `func EGFRCKDEPI2021SI(creatinineUmolL, ageYears float64, female bool) (float64, error)`
Formula and output: EGFRCKDEPI2021SI returns the estimated glomerular filtration rate, in mL/min/1.73 m², from serum creatinine in µmol/L: creatinine ÷ 88.4 gives mg/dL, then the 2021 CKD-EPI creatinine equation (MED-RENAL-001).
Inputs and units: creatinineUmolL µmol/L; ageYears years (decimal); female the source's sex coefficient; result mL/min/1.73 m²
Range: ageYears ≥ 18; creatinineUmolL > 0; no upper limits stated by the source
Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range; ErrDomain when the result is beyond float64.
Primary source: Inker LA, Eneanya ND, Coresh J, et al. New creatinine- and cystatin C-based equations to estimate GFR without race. N Engl J Med. 2021;385(19):1737-1749. doi:10.1056/NEJMoa2102953. PMID 34554658. Table 2, eGFRcr(AS) column (full text PMC8822996)
Primary source: Kidney Disease: Improving Global Outcomes (KDIGO) CKD Work Group. KDIGO 2024 Clinical Practice Guideline for the Evaluation and Management of Chronic Kidney Disease. Kidney Int. 2024;105(4S):S117–S314. doi:10.1016/j.kint.2023.10.018. p. S126, "Conversion factors of conventional units to SI units" (creatinine 88.4)
Worked example: female, 50 y, Scr 61.88 µmol/L (= 0.7 mg/dL) → 105.29760114914401 (MED-RENAL-001), to 1e-12 relative
Best approach: Time O(1); Space O(1), no allocation
Tags: renal, egfr, ckd-epi, creatinine, si

### EGFRCKDEPI2021CystatinC      (ID: MED-RENAL-003, built: med/renal/egfr_ckd_epi_2021_cystatin_c.go:20)
Signature: `func EGFRCKDEPI2021CystatinC(creatinineMgDL, cystatinCMgL, ageYears float64, female bool) (float64, error)`
Formula and output: EGFRCKDEPI2021CystatinC returns the estimated glomerular filtration rate, in mL/min/1.73 m², from serum creatinine in mg/dL and serum cystatin C in mg/L by the 2021 CKD-EPI creatinine–cystatin C equation (race-free): 135 × min(Scr/κ, 1)^a1 × max(Scr/κ, 1)^−0.544 × min(Scys/0.8, 1)^−0.323 × max(Scys/0.8, 1)^−0.778 × 0.9961^age × 0.963 [if female], with κ = 0.7 (female) or 0.9 (male) and a1 = −0.219 (female) or −0.144 (male).
Inputs and units: creatinineMgDL mg/dL; cystatinCMgL mg/L; ageYears years (decimal); female the source's sex coefficient; result mL/min/1.73 m²
Range: ageYears ≥ 18 ("All the participants were 18 years of age or older"); creatinineMgDL > 0; cystatinCMgL > 0; no upper limits stated by the source
Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range (every in-range input gives a finite result).
Primary source: Inker LA, Eneanya ND, Coresh J, et al. New creatinine- and cystatin C-based equations to estimate GFR without race. N Engl J Med. 2021;385(19):1737-1749. doi:10.1056/NEJMoa2102953. PMID 34554658. Table 2, eGFRcr-cys(AS) column (full text PMC8822996)
Worked example: female, 50 y, Scr 0.7, Scys 0.8 → 135 × 0.9961^50 × 0.963 = 106.9318450861178 (hand-derived from Table 2)
Best approach: Time O(1); Space O(1), no allocation
Tags: renal, egfr, ckd-epi, creatinine, cystatin-c

### EGFRCKDEPI2021CystatinCSI      (ID: MED-RENAL-004, built: med/renal/egfr_ckd_epi_2021_cystatin_c.go:56)
Signature: `func EGFRCKDEPI2021CystatinCSI(creatinineUmolL, cystatinCMgL, ageYears float64, female bool) (float64, error)`
Formula and output: EGFRCKDEPI2021CystatinCSI returns the estimated glomerular filtration rate, in mL/min/1.73 m², from serum creatinine in µmol/L and serum cystatin C in mg/L: creatinine ÷ 88.4 gives mg/dL, then the 2021 CKD-EPI creatinine–cystatin C equation (MED-RENAL-003).
Inputs and units: creatinineUmolL µmol/L; cystatinCMgL mg/L; ageYears years (decimal); female the source's sex coefficient; result mL/min/1.73 m²
Range: ageYears ≥ 18; creatinineUmolL > 0; cystatinCMgL > 0; no upper limits stated by the source
Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range; ErrDomain when the result is beyond float64.
Primary source: Inker LA, Eneanya ND, Coresh J, et al. New creatinine- and cystatin C-based equations to estimate GFR without race. N Engl J Med. 2021;385(19):1737-1749. doi:10.1056/NEJMoa2102953. PMID 34554658. Table 2, eGFRcr-cys(AS) column (full text PMC8822996)
Primary source: Kidney Disease: Improving Global Outcomes (KDIGO) CKD Work Group. KDIGO 2024 Clinical Practice Guideline for the Evaluation and Management of Chronic Kidney Disease. Kidney Int. 2024;105(4S):S117–S314. doi:10.1016/j.kint.2023.10.018. p. S126, "Conversion factors of conventional units to SI units" (creatinine 88.4)
Worked example: female, 50 y, Scr 61.88 µmol/L (= 0.7 mg/dL), Scys 0.8 → 106.9318450861178 (MED-RENAL-003), to 1e-12 relative
Best approach: Time O(1); Space O(1), no allocation
Tags: renal, egfr, ckd-epi, creatinine, cystatin-c, si

### CreatinineClearanceCockcroftGault      (ID: MED-RENAL-005, built: med/renal/creatinine_clearance_cockcroft_gault.go:18)
Signature: `func CreatinineClearanceCockcroftGault(ageYears, weightKg, creatinineMgDL float64, female bool) (float64, error)`
Formula and output: CreatinineClearanceCockcroftGault returns the predicted creatinine clearance, in mL/min, by the Cockcroft–Gault formula: (140 − age) × weightKg / (72 × Scr mg/dL), × 0.85 if female (the paper's "15% less in females"; the derivation population was adult males).
Inputs and units: ageYears years (decimal); weightKg kg; creatinineMgDL mg/dL (= mg/100 mL); female the source's sex factor; result mL/min
Range: ageYears 18–92 (the 249 persons aged 18–92 the formula was derived from); weightKg > 0; creatinineMgDL > 0
Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range; ErrDomain when the result is beyond float64.
Primary source: Cockcroft DW, Gault MH. Prediction of creatinine clearance from serum creatinine. Nephron. 1976;16(1):31-41. doi:10.1159/000180580. PMID 1244564. Formula as printed on the publisher abstract page: "Ccr = (140 – age) (wt kg)/72 × Scr(mg/100ml)" "(15% less in females)"
Worked example: male, 40 y, 72 kg, Scr 1.0 → (140 − 40) × 72 / (72 × 1.0) = 100 mL/min; female, same → 85 (hand-derived)
Best approach: Time O(1); Space O(1), no allocation
Tags: renal, creatinine-clearance, cockcroft-gault, creatinine

### CreatinineClearanceCockcroftGaultSI      (ID: MED-RENAL-006, built: med/renal/creatinine_clearance_cockcroft_gault.go:42)
Signature: `func CreatinineClearanceCockcroftGaultSI(ageYears, weightKg, creatinineUmolL float64, female bool) (float64, error)`
Formula and output: CreatinineClearanceCockcroftGaultSI returns the predicted creatinine clearance, in mL/min, from serum creatinine in µmol/L: creatinine ÷ 88.4 gives mg/dL, then the Cockcroft–Gault formula (MED-RENAL-005).
Inputs and units: ageYears years (decimal); weightKg kg; creatinineUmolL µmol/L; female the source's sex factor; result mL/min
Range: ageYears 18–92; weightKg > 0; creatinineUmolL > 0
Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range; ErrDomain when the result is beyond float64.
Primary source: Cockcroft DW, Gault MH. Prediction of creatinine clearance from serum creatinine. Nephron. 1976;16(1):31-41. doi:10.1159/000180580. PMID 1244564. Formula as printed on the publisher abstract page: "Ccr = (140 – age) (wt kg)/72 × Scr(mg/100ml)" "(15% less in females)"
Primary source: Kidney Disease: Improving Global Outcomes (KDIGO) CKD Work Group. KDIGO 2024 Clinical Practice Guideline for the Evaluation and Management of Chronic Kidney Disease. Kidney Int. 2024;105(4S):S117–S314. doi:10.1016/j.kint.2023.10.018. p. S126, "Conversion factors of conventional units to SI units" (creatinine 88.4)
Worked example: male, 40 y, 72 kg, Scr 88.4 µmol/L (= 1.0 mg/dL) → 100 mL/min (hand-derived)
Best approach: Time O(1); Space O(1), no allocation
Tags: renal, creatinine-clearance, cockcroft-gault, creatinine, si

### EGFRSchwartzBedside      (ID: MED-RENAL-007, built: med/renal/egfr_schwartz_bedside.go:18)
Signature: `func EGFRSchwartzBedside(heightCm, creatinineMgDL, ageYears float64) (float64, error)`
Formula and output: EGFRSchwartzBedside returns the estimated glomerular filtration rate, in mL/min/1.73 m², of a child by the bedside CKiD Schwartz equation: 0.413 × heightCm / creatinineMgDL. ageYears is taken only to enforce the Range.
Inputs and units: heightCm cm; creatinineMgDL mg/dL; ageYears years (decimal); result mL/min/1.73 m²
Range: ageYears 1–16 (349 children aged 1 to 16, CKiD); heightCm > 0; creatinineMgDL > 0
Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range; ErrDomain when the result is beyond float64.
Primary source: Schwartz GJ, Muñoz A, Schneider MF, et al. New equations to estimate GFR in children with CKD. J Am Soc Nephrol. 2009;20(3):629-637. doi:10.1681/ASN.2008030287. PMID 19158356. Abstract: "bedside calculation of 0.413*(height/serum creatinine)"
Worked example: height 100 cm, Scr 0.413 → 0.413 × 100 / 0.413 = 100 (hand-derived)
Best approach: Time O(1); Space O(1), no allocation
Tags: renal, egfr, schwartz, pediatric, creatinine

### EGFRSchwartzBedsideSI      (ID: MED-RENAL-008, built: med/renal/egfr_schwartz_bedside.go:42)
Signature: `func EGFRSchwartzBedsideSI(heightCm, creatinineUmolL, ageYears float64) (float64, error)`
Formula and output: EGFRSchwartzBedsideSI returns the estimated glomerular filtration rate, in mL/min/1.73 m², of a child from serum creatinine in µmol/L: creatinine ÷ 88.4 gives mg/dL, then the bedside Schwartz equation (MED-RENAL-007).
Inputs and units: heightCm cm; creatinineUmolL µmol/L; ageYears years (decimal); result mL/min/1.73 m²
Range: ageYears 1–16; heightCm > 0; creatinineUmolL > 0
Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range; ErrDomain when the result is beyond float64.
Primary source: Schwartz GJ, Muñoz A, Schneider MF, et al. New equations to estimate GFR in children with CKD. J Am Soc Nephrol. 2009;20(3):629-637. doi:10.1681/ASN.2008030287. PMID 19158356. Abstract: "bedside calculation of 0.413*(height/serum creatinine)"
Primary source: Kidney Disease: Improving Global Outcomes (KDIGO) CKD Work Group. KDIGO 2024 Clinical Practice Guideline for the Evaluation and Management of Chronic Kidney Disease. Kidney Int. 2024;105(4S):S117–S314. doi:10.1016/j.kint.2023.10.018. p. S126, "Conversion factors of conventional units to SI units" (creatinine 88.4)
Worked example: height 100 cm, Scr 88.4 µmol/L (= 1.0 mg/dL) → 0.413 × 100 / 1.0 = 41.3 (hand-derived)
Best approach: Time O(1); Space O(1), no allocation
Tags: renal, egfr, schwartz, pediatric, creatinine, si

### FractionalExcretionSodium      (ID: MED-RENAL-009, built: med/renal/fractional_excretion_sodium.go:17)
Signature: `func FractionalExcretionSodium(urineSodiumMmolL, plasmaSodiumMmolL, urineCreatinine, plasmaCreatinine float64) (float64, error)`
Formula and output: FractionalExcretionSodium returns the fractional excretion of sodium (FENa), in percent: (UNa × PCr) / (PNa × UCr) × 100. The two creatinine values may be in any one unit (mg/dL or µmol/L), the same for both: the ratio cancels it. The source reports group findings in oliguric acute renal failure, not a classification table, so no category is returned.
Inputs and units: urineSodiumMmolL, plasmaSodiumMmolL mmol/L; urineCreatinine, plasmaCreatinine any one unit, the same for both; result percent
Range: every input > 0 (population: oliguric acute renal failure; no input limits stated by the source)
Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range; ErrDomain when the result is beyond float64.
Primary source: Espinel CH. The FENa test. Use in the differential diagnosis of acute renal failure. JAMA. 1976;236(6):579-581. doi:10.1001/jama.236.6.579. PMID 947239
Worked example: UNa 20, PNa 140, UCr 100, PCr 2 → (20 × 2) / (140 × 100) × 100 = 0.2857142857142857 % (hand-derived)
Best approach: Time O(1); Space O(1), no allocation
Tags: renal, fena, sodium, creatinine, urine

### FractionalExcretionUrea      (ID: MED-RENAL-010, built: med/renal/fractional_excretion_urea.go:17)
Signature: `func FractionalExcretionUrea(urineUrea, plasmaUrea, urineCreatinine, plasmaCreatinine float64) (float64, error)`
Formula and output: FractionalExcretionUrea returns the fractional excretion of urea (FEUrea), in percent: (UUrea × PCr) / (PUrea × UCr) × 100. The urea pair may be in any one unit (mg/dL of urea nitrogen or mmol/L), the same for both, and the creatinine pair likewise: each ratio cancels its unit. No category is returned.
Inputs and units: urineUrea, plasmaUrea any one unit, the same for both; urineCreatinine, plasmaCreatinine any one unit, the same for both; result percent
Range: every input > 0 (no input limits stated by the source)
Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range; ErrDomain when the result is beyond float64.
Primary source: Carvounis CP, Nisar S, Guro-Razuman S. Significance of the fractional excretion of urea in the differential diagnosis of acute renal failure. Kidney Int. 2002;62(6):2223-2229. doi:10.1046/j.1523-1755.2002.00683.x. PMID 12427149
Worked example: UUrea 300, PUrea 30, UCr 100, PCr 2 → (300 × 2) / (30 × 100) × 100 = 20 % (hand-derived)
Best approach: Time O(1); Space O(1), no allocation
Tags: renal, feurea, urea, creatinine, urine

### GFRCategoryKDIGO      (ID: MED-RENAL-011, built: med/renal/gfr_category_kdigo.go:52)
Signature: `func GFRCategoryKDIGO(egfrMLMin173 float64) (GFRCategory, error)`
Formula and output: GFRCategoryKDIGO returns the KDIGO GFR category for a GFR in mL/min/1.73 m²: ≥ 90 G1, ≥ 60 G2, ≥ 45 G3a, ≥ 30 G3b, ≥ 15 G4, else G5. The table is printed in integers (90+, 60–89, 45–59, 30–44, 15–29, < 15); a value between two printed integers (89.5) takes the lower category (ASSUMPTION A-M3).
Inputs and units: egfrMLMin173 mL/min/1.73 m²; result a category code (Term gives the source's label)
Range: egfrMLMin173 ≥ 0
Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange for egfrMLMin173 < 0.
Primary source: Kidney Disease: Improving Global Outcomes (KDIGO) CKD Work Group. KDIGO 2024 Clinical Practice Guideline for the Evaluation and Management of Chronic Kidney Disease. Kidney Int. 2024;105(4S):S117–S314. doi:10.1016/j.kint.2023.10.018. p. S125, "Current chronic kidney disease (CKD) nomenclature used by KDIGO", GFR categories
Worked example: 89.999 → G2 "Mildly decreased"; 90 → G1 "Normal or high" (KDIGO 2024 p. S125)
Best approach: Time O(1); Space O(1), no allocation
Tags: renal, egfr, kdigo, ckd, category

### DiffusingCapacityRoughtonForster      (ID: MED-RESP-001, built: med/respiratory/diffusing_capacity_roughton_forster.go:47)
Signature: `func DiffusingCapacityRoughtonForster(membraneDLMLMinMmHg, thetaMLMinMmHgPerMLBlood, capillaryBloodVolumeML float64) (float64, error)`
Formula and output: DiffusingCapacityRoughtonForster returns the lung diffusing capacity, in mL/min/mmHg, from its two components by the Roughton–Forster equation 1/DL = 1/DM + 1/(θ·VC): DM the membrane diffusing capacity, θ the gas uptake of 1 mL of blood per minute per mmHg, VC the pulmonary capillary blood volume. It is computed as DM·θVC / (DM + θVC), the form without three reciprocals.
Inputs and units: membraneDLMLMinMmHg mL/min/mmHg; thetaMLMinMmHgPerMLBlood mL/min/mmHg per mL of blood; capillaryBloodVolumeML mL; result mL/min/mmHg
Range: every input > 0 (the relation is physiological, with no population limits stated in the source)
Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange for a non-positive input; ErrDomain when θ·VC or DM + θ·VC overflows float64.
Primary source: Zavorsky GS, Hsia CCW, Hughes JMB, et al. Standardisation and application of the single-breath determination of nitric oxide uptake in the lung. Eur Respir J. 2017;49(2):1600962. doi:10.1183/13993003.00962-2016. p. 2 "Origins of DLNO" and Figure 1 legend: "Roughton–Forster equation, 1/DL=1/DM+1/(θ·VC)"
Primary source: Roughton FJW, Forster RE. Relative importance of diffusion and chemical reaction rates in determining rate of exchange of gases in the human lung, with special reference to true diffusing capacity of pulmonary membrane and volume of blood in the lung capillaries. J Appl Physiol. 1957;11(2):290-302. doi:10.1152/jappl.1957.11.2.290. PMID 13475180
Worked example: DM 50, θ 1.0, VC 75 → 50·75/125 = 30 mL/min/mmHg (hand-derived)
Best approach: Time O(1); Space O(1), no allocation
Tags: respiratory, diffusing-capacity, dlco, roughton-forster, gas-exchange

### EpworthSleepinessScore      (ID: MED-SCALE-001, built: med/scales/epworth_sleepiness_score.go:30)
Signature: `func EpworthSleepinessScore(items [8]int) (int, error)`
Formula and output: EpworthSleepinessScore returns the Epworth Sleepiness Scale total, 0–24: the sum of the 8 item ratings, each rated 0 to 3 on the response key of Johns 1991, Table 1. The item wording is not embedded here (the scale is © Murray W Johns; the caller holds the licensed text); items[i] is the rating of item i+1 in the source's order.
Inputs and units: items ratings, points 0–3 each; result points 0–24
Range: each item 0–3 (population: 180 adults aged 18–78 with sleep disorders, and controls)
Errors: ErrOutOfRange when any item is outside 0–3, naming the item index.
Primary source: Johns MW. A new method for measuring daytime sleepiness: the Epworth sleepiness scale. Sleep. 1991;14(6):540-545. doi:10.1093/sleep/14.6.540. PMID 1798888. Table 1
Worked example: all 8 items rated 1 → 8; all rated 3 → 24; all 0 → 0 (HAND-DERIVED; Johns 1991 controls scored 2–10)
Best approach: Time O(1); Space O(1), no allocation
Tags: scales, epworth, ess, sleep, sleepiness, score

### EpworthCategory      (ID: MED-SCALE-002, built: med/scales/epworth_category.go:33)
Signature: `func EpworthCategory(score int) (EpworthSleepinessCategory, error)`
Formula and output: EpworthCategory returns the Epworth Sleepiness Scale band for a total, with the official site's own label: 0–5 "Lower Normal Daytime Sleepiness"; 6–10 "Higher Normal Daytime Sleepiness"; 11–12 "Mild Excessive Daytime Sleepiness"; 13–15 "Moderate Excessive Daytime Sleepiness"; 16–24 "Severe Excessive Daytime Sleepiness".
Inputs and units: score points 0–24; result the official site's band label
Range: 0 ≤ score ≤ 24
Errors: ErrOutOfRange when score is outside 0–24.
Primary source: Johns MW. About the ESS. The Epworth Sleepiness Scale official website. https://epworthsleepinessscale.com/about-the-ess/ (read 2026-11-12)
Worked example: 5 → "Lower Normal Daytime Sleepiness"; 6 → "Higher Normal Daytime Sleepiness"; 16 → "Severe Excessive Daytime Sleepiness"
Best approach: Time O(1); Space O(1), no allocation
Tags: scales, epworth, ess, sleep, sleepiness, category

### GlasgowComaScore      (ID: MED-SCALE-003, built: med/scales/glasgow_coma_score.go:61)
Signature: `func GlasgowComaScore(eye, verbal, motor int) (int, error)`
Formula and output: GlasgowComaScore returns the Glasgow Coma Scale total, 3–15: eye (1–4) + verbal (1–5) + motor (1–6). A component rated "Non testable" is passed as GCSNotTestable (0); then no total exists and the function returns ErrNotTestable. Teasdale et al. 2014 (abstract): "Individual patients are best described by the three components of the coma scale; whereas the derived total coma score should be used to characterise groups."
Inputs and units: eye points 1–4; verbal points 1–5; motor points 1–6 (0 = GCSNotTestable); result points 3–15
Range: eye 1–4, verbal 1–5, motor 1–6, each or GCSNotTestable (0)
Errors: ErrOutOfRange when a component is outside its range and not GCSNotTestable (checked for all three first); ErrNotTestable when any component is GCSNotTestable.
Primary source: Teasdale G. Glasgow Coma Scale: Assessment Aid (English). 2015. https://www.glasgowcomascale.org/downloads/GCS-Assessment-Aid-English.pdf (read 2026-11-12): eye 4 Spontaneous, 3 To sound, 2 To pressure, 1 None; verbal 5 Orientated, 4 Confused, 3 Words, 2 Sounds, 1 None; motor 6 Obeys commands, 5 Localising, 4 Normal flexion, 3 Abnormal flexion, 2 Extension, 1 None; NT Non testable.
Primary source: Teasdale G, Jennett B. Assessment of coma and impaired consciousness. A practical scale. Lancet. 1974;2(7872):81-84. doi:10.1016/s0140-6736(74)91639-0. PMID 4136544.
Primary source: Teasdale G, Maas A, Lecky F, Manley G, Stocchetti N, Murray G. The Glasgow Coma Scale at 40 years: standing the test of time. Lancet Neurol. 2014;13(8):844-854. doi:10.1016/S1474-4422(14)70120-6. PMID 25030516.
Worked example: eye 3 (To sound), verbal 4 (Confused), motor 5 (Localising) → 12; 4 + 5 + 6 → 15; 1 + 1 + 1 → 3 (HAND-DERIVED from the Assessment Aid)
Best approach: Time O(1); Space O(1), no allocation
Tags: scales, gcs, glasgow-coma-scale, consciousness, neurology, score

### GlasgowComaSeverityRimel      (ID: MED-SCALE-004, built: med/scales/glasgow_coma_severity_rimel.go:29)
Signature: `func GlasgowComaSeverityRimel(total int) (GCSSeverity, error)`
Formula and output: GlasgowComaSeverityRimel returns Rimel et al.'s 1982 head-injury group for a Glasgow Coma Scale total, with the source's own word: 3–8 "severe"; 9–12 "moderate"; 13–15 "minor". These are Rimel's head-injury groups; the official Glasgow Coma Scale site states no bands.
Inputs and units: total points 3–15; result the source's group word
Range: 3 ≤ total ≤ 15
Errors: ErrOutOfRange when total is outside 3–15.
Primary source: Rimel RW, Giordani B, Barth JT, Jane JA. Moderate head injury: completing the clinical spectrum of brain trauma. Neurosurgery. 1982;11(3):344-351. doi:10.1227/00006123-198209000-00002. PMID 7133350. Abstract: "severe, 3-8; moderate, 9-12; and minor, 13-15"
Worked example: 8 → "severe"; 9 → "moderate"; 13 → "minor" (Rimel 1982 abstract)
Best approach: Time O(1); Space O(1), no allocation
Tags: scales, gcs, glasgow-coma-scale, head-injury, category

### ApgarScore      (ID: MED-SCALE-005, built: med/scales/apgar_score.go:27)
Signature: `func ApgarScore(heartRate, respiratoryEffort, reflexIrritability, muscleTone, color int) (int, error)`
Formula and output: ApgarScore returns the Apgar score, 0–10: the sum of the five signs (heart rate, respiratory effort, reflex irritability, muscle tone, color), each rated 0, 1 or 2 (Apgar 1953: "A rating of zero, one or two, was given to each sign"; "A score of ten indicated a baby in the best possible condition"). The paper rated the signs at "sixty seconds after the complete birth of the baby"; the function sums whatever ratings it is given.
Inputs and units: each sign points 0–2; result points 0–10
Range: heartRate, respiratoryEffort, reflexIrritability, muscleTone, color each 0–2
Errors: ErrOutOfRange when a sign is outside 0–2, naming the sign.
Primary source: Apgar V. A proposal for a new method of evaluation of the newborn infant. Curr Res Anesth Analg. 1953;32(4):260-267. PMID 13083014. pp. 261–262 (read in the NLM scan)
Worked example: all five rated 2 → 10; all 0 → 0; 2, 1, 1, 1, 0 → 5 (HAND-DERIVED from pp. 261–262)
Best approach: Time O(1); Space O(1), no allocation
Tags: scales, apgar, newborn, neonatal, score

### ApgarCondition1953      (ID: MED-SCALE-006, built: med/scales/apgar_condition_1953.go:28)
Signature: `func ApgarCondition1953(score int) (ApgarCondition, error)`
Formula and output: ApgarCondition1953 returns Apgar's 1953 group for an Apgar score, with the source's own words: 8–10 "good condition"; 3–7 "fair condition"; 0–2 "poor condition". Later bands (AAP/ACOG) are not built here.
Inputs and units: score points 0–10; result the source's group words
Range: 0 ≤ score ≤ 10
Errors: ErrOutOfRange when score is outside 0–10.
Primary source: Apgar V. A proposal for a new method of evaluation of the newborn infant. Curr Res Anesth Analg. 1953;32(4):260-267. PMID 13083014. p. 267: "groups 8, 9, and 10 indicate infants in good condition, 0, 1, and 2, poor condition, and the remaining scores, fair condition."
Worked example: 2 → "poor condition"; 3 → "fair condition"; 8 → "good condition" (Apgar 1953 p. 267)
Best approach: Time O(1); Space O(1), no allocation
Tags: scales, apgar, newborn, neonatal, category

## 4. HELD by policy (drug dose or amount to administer) — catalogued, not built (D-M5)

| Row | What it would compute | Source it would be built from | Status |
|---|---|---|---|
| Dosage Calculation (and the six named variants: Amount of Drug on Hand, Amount to be Administered, Desired Dose, Quantity on Hand, Dosage Amount (Size of Patient), Dosage rate) | D/H × Q family; dose by weight; rate | a pharmaceutical-calculations text chosen by Devin (opentextbc "Math for Medication Administration" was not reachable: 403) | LEFT FOR DEVIN |
| Pill Dosage Calculation | tablets = desired dose / tablet strength | as above | LEFT FOR DEVIN |
| IV Dosage Calculation | dose from rate and concentration | as above | LEFT FOR DEVIN |
| Pediatric Acetaminophen Dose | volume by weight band | DailyMed, Children's Acetaminophen Oral Suspension 160 mg/5 mL, setid 8231dd91-e9e6-4350-a485-fa95d65a6d25 (bands 24–35 lb 5 mL … 72–95 lb 15 mL; under 24 lb "ask a doctor"; every 4 h, at most 6 doses in 24 h) | LEFT FOR DEVIN |
| Dewormer Dosage Calculator; Pyrantel Pamoate Dewormer for Puppies and Dogs | 5 mL per 10 lb (4.54 mg/mL pyrantel base) | FDA SPL pyrantel pamoate suspension labels (e.g. accessdata.fda.gov/spl/data/c746ab3c-62e6-4dfe-9a4f-a9d3090a3b67); veterinary, so a home outside `med/` is also Devin's call | LEFT FOR DEVIN |
| Parkland Replacement Fluids | 4 mL × kg × %TBSA / 24 h | Baxter & Shires 1968 (Ann N Y Acad Sci 150:874, PMID 4973463; formula not seen) or ABA 2024 (Cartotto R, J Burn Care Res 45:565, PMID 38051821: "2 mL/kg/% TBSA") | LEFT FOR DEVIN |
| Maintenance fluids (Holliday–Segar 100/50/20; hourly 4-2-1) | mL/day, mL/h | Holliday & Segar 1957 (Pediatrics 19:823, PMID 13431307); Oh 1980 (Anesthesiology 53:351, PMID 7425366); neither read | LEFT FOR DEVIN |
| mg → mL at a given concentration (pure conversion, no drug named) | volume = mass / concentration | permitted by the request, but no citable text was reachable for it in this run | LEFT FOR DEVIN |

## 5. HELD-SOURCE: what must be read to release each (citation already resolved)

| Function (planned home) | Citation confirmed | What was not seen |
|---|---|---|
| Du Bois BSA (anthropometry) | Du Bois D, Du Bois EF. Arch Intern Med. 1916;17(6_2):863-871. doi:10.1001/archinte.1916.00080130010002; reprint Nutrition 1989;5:303, PMID 2520314 | the formula in the paper |
| Mosteller BSA (anthropometry) | Mosteller RD. N Engl J Med. 1987;317(17):1098. doi:10.1056/NEJM198710223171717. PMID 3657876 | the letter's text (403) |
| Devine IBW (anthropometry) | McCarron MM, Devine BJ. Drug Intell Clin Pharm. 1974;8(11):650-655. doi:10.1177/106002807400801104 | the formula (not in PubMed; publisher 403) |
| APLS weight 1–10 y (anthropometry) | APLS 5th ed. (Samuels M, ed., Wiley-Blackwell/ALSG 2011) via Seddon 2012 PMID 23050162 | the manual |
| Nadler blood volume (fluids) | Nadler SB, Hidalgo JH, Bloch T. Surgery. 1962;51(2):224-232. PMID 21936146 | coefficients, population |
| Rule of nines; Lund–Browder (fluids) | Wallace AB. Lancet. 1951;1(6653):501-504. PMID 14805109; Lund CC, Browder NC. Surg Gynecol Obstet. 1944;79:352-358 | both tables in the papers |
| MAP, pulse pressure, CO, CI, SVR/TPR, stroke work, cardiac efficiency (cardio) | Hall JE, Hall ME. Guyton and Hall Textbook of Medical Physiology. 14th ed. Elsevier; 2021. ISBN 978-0-323-59712-8, ch. 9, 14, 15, 20 | the book (not opened); MAP one-third rule has no primary found |
| Rate–pressure product (cardio) | Gobel FL et al. Circulation. 1978;57(3):549-556. doi:10.1161/01.CIR.57.3.549. PMID 624164 | "systolic" (abstract says "HR x blood pressure") |
| QTc Bazett (cardio) | Bazett HC. Heart. 1920;7:353-370; reprint Ann Noninvasive Electrocardiol. 1997;2(2):177-194. doi:10.1111/j.1542-474X.1997.tb00325.x | either text |
| Karvonen target HR (cardio) | Karvonen MJ, Kentala E, Mustala O. Ann Med Exp Biol Fenn. 1957;35(3):307-315. PMID 13470504 | any text |
| BP categories (cardio) | Jones DW et al. Hypertension. 2025;82(10):e212-e316. doi:10.1161/HYP.0000000000000249. PMID 40811516 | the category table (403) |
| Starling filtration (cardio) | Starling EH. J Physiol. 1896;19(4):312, PMID 16992325; Levick JR, Michel CC. Cardiovasc Res. 2010;87(2):198-210. doi:10.1093/cvr/cvq062 | the modern equation |
| Alveolar gas equation, A–a gradient, CO₂ elimination, minute ventilation, RER/RQ, CaO₂ (respiratory) | Riley RL, Cournand A. J Appl Physiol. 1949;1(12):825-847. doi:10.1152/jappl.1949.1.12.825; West JB, Luks AM. West's Respiratory Physiology. 11th ed. Wolters Kluwer; 2021. ISBN 978-1-975139-18-6 | the equations (403; book not opened); 1.34 vs 1.39 to be chosen |
| Dead space Bohr / Enghoff; shunt (respiratory) | Bohr C. Skand Arch Physiol. 1891;2:236, doi:10.1111/j.1748-1716.1891.tb00581.x; Enghoff 1938; Berggren 1942 | the papers |
| DLCO single-breath, Hb correction (respiratory) | Graham BL et al. Eur Respir J. 2017;49:1600016. doi:10.1183/13993003.00016-2016. PMID 28049168 | the equations (403) |
| Widmark BAC (pharmacology) | Widmark 1932 (book); Watson PE, Watson ID, Batt RD. J Stud Alcohol. 1981;42(7):547. doi:10.15288/jsa.1981.42.547. PMID 7289599; Watson 1980 AJCN 33:27, PMID 6986753 | r and β values; TBW coefficients |
| MDRD; Friedewald; APRI; anion gap; corrected sodium (renal/chemistry/hepatic) | see §2 | see §2 |
| Tubular reabsorption rate (renal) | Guyton and Hall 14th ed. ch. 28 | the book |
| MELD-Na threshold (hepatic; MED-HEPAT-003 built without it) | Kim WR et al. N Engl J Med. 2008;359:1018-1026 (original MELD-Na) | whether a MELD > 11 condition belongs in the function (A-M9) |
| MELD 3.0 dialysis rule (MED-HEPAT-004 built without it) | Kim 2021 Methods | whether a dialysis rule exists in the source (A-M10); if it does, a new ID with the extra parameter |

## 6. SOURCE NOT FOUND · LATER · MOVE

- SOURCE NOT FOUND: Rule of Fives (no infant primary; Livingston & Lee 2000 J Surg Res 91:106, PMID 10839957, is for adults
  over 80 kg, and its figures as summarised sum to 102%); Conduction Velocity (not researched to a primary); Herbal Remedy
  Converter; BSA mean data items; corrected calcium 0.8 form; osmolality SI 2Na + glucose + urea form.
- LATER: Time of Death (Henssge C. Forensic Sci Int. 1988;38(3-4):209-236. doi:10.1016/0379-0738(88)90168-5; needs the
  whole nomogram model); Weber fraction (psychophysics, no home); VO₂max; BMI-for-age LMS.
- MOVE: Hagen–Poiseuille and Banana Equivalent Dose → `sci/physics` (BED has no peer-reviewed origin: a 1995 mailing-list
  post; EPA FGR 11 / ICRP 72 hold the K-40 dose figures); Diluting Stock Solutions → `sci/chemistry`; Direct Care FTEs →
  `fin/business`; Bayes for disease testing → `num/prob` NUM-PROB-010/-011.
