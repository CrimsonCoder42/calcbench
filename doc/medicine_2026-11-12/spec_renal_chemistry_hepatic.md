# Build spec: med/renal, med/chemistry, med/hepatic (medicine build 2026-11-12)

Read `CONVENTIONS.md` in this folder first; it is binding. Citations below were confirmed 2026-11-12 on PubMed
(efetch) and, where named, the PMC full text or publisher page. Use them verbatim in `Source:` lines (one line
each). Never invent a coefficient, threshold or label not written here.

Shared conversion (put it in the file of the first SI function of each package, unexported, with this comment):
- `creatinineUmolLPerMgDL = 88.4` — KDIGO 2024 CKD guideline, Kidney Int. 2024;105(4S):S117–S314,
  doi:10.1016/j.kint.2023.10.018, p. S126 "Conversion factors of conventional units to SI units" (read 2026-11-12).
- `calciumMmolLPerMgDL = 0.2495` — same table, same page.

## med/renal (SHORT RENAL)

### MED-RENAL-001 EGFRCKDEPI2021(creatinineMgDL, ageYears float64, female bool) (float64, error)
eGFR (mL/min/1.73 m²) = 142 × min(Scr/κ,1)^α × max(Scr/κ,1)^−1.200 × 0.9938^age × 1.012 [if female];
κ = 0.7 female, 0.9 male; α = −0.241 female, −0.302 male.
Source: Inker LA, Eneanya ND, Coresh J, et al. New creatinine- and cystatin C-based equations to estimate GFR without race. N Engl J Med. 2021;385(19):1737-1749. doi:10.1056/NEJMoa2102953. PMID 34554658. Table 2, eGFRcr(AS) column (full text PMC8822996). NOT the ASR-NB columns.
Range: ageYears ≥ 18 ("All the participants were 18 years of age or older"; no upper age stated); creatinineMgDL > 0 (no limits stated). State "no upper limits stated by the source" in Range.
Example (HAND-DERIVED, compute and write the arithmetic in the test comment): female, 50 y, Scr 0.7 → 142 × 1 × 1 × 0.9938^50 × 1.012; male 60 y, Scr 1.2 → 142 × (1.2/0.9)^−1.200 × 0.9938^60. Give values to 1e-9 relative from an independent computation in the test (write the formula out in the test, not by calling the function).

### MED-RENAL-002 EGFRCKDEPI2021SI(creatinineUmolL, ageYears float64, female bool)
Converts creatinine µmol/L ÷ 88.4 → mg/dL, then MED-RENAL-001. Second Source line: KDIGO 2024 p. S126.

### MED-RENAL-003 EGFRCKDEPI2021CystatinC(creatinineMgDL, cystatinCMgL, ageYears float64, female bool)
eGFR = 135 × min(Scr/κ,1)^a1 × max(Scr/κ,1)^−0.544 × min(Scys/0.8,1)^−0.323 × max(Scys/0.8,1)^−0.778 × 0.9961^age × 0.963 [if female];
κ 0.7 F / 0.9 M; a1 = −0.219 female, −0.144 male. Same source, Table 2, eGFRcr-cys(AS) column. Range as 001 plus cystatinCMgL > 0.

### MED-RENAL-004 EGFRCKDEPI2021CystatinCSI(creatinineUmolL, cystatinCMgL, ageYears float64, female bool)
µmol/L ÷ 88.4, then 003.

### MED-RENAL-005 CreatinineClearanceCockcroftGault(ageYears, weightKg, creatinineMgDL float64, female bool)
Ccr (mL/min) = (140 − age) × weightKg / (72 × Scr mg/dL); × 0.85 if female ("15% less in females").
Source: Cockcroft DW, Gault MH. Prediction of creatinine clearance from serum creatinine. Nephron. 1976;16(1):31-41. doi:10.1159/000180580. PMID 1244564. Formula as printed on the publisher abstract page: "Ccr = (140 – age) (wt kg)/72 × Scr(mg/100ml)" "(15% less in females)".
Range: ageYears 18–92 (the 249 persons aged 18–92 the formula was derived from); weightKg > 0; creatinineMgDL > 0.
Note: the derivation population was adult males; the female factor is the paper's statement.
Example (HAND-DERIVED): male 40 y, 72 kg, Scr 1.0 → 100 mL/min exactly; female same → 85.

### MED-RENAL-006 CreatinineClearanceCockcroftGaultSI(ageYears, weightKg, creatinineUmolL float64, female bool)
÷ 88.4 then 005. Test: 88.4 µmol/L equals 1.0 mg/dL result.

### MED-RENAL-007 EGFRSchwartzBedside(heightCm, creatinineMgDL, ageYears float64) (float64, error)
eGFR (mL/min/1.73 m²) = 0.413 × heightCm / creatinineMgDL.
Source: Schwartz GJ, Muñoz A, Schneider MF, et al. New equations to estimate GFR in children with CKD. J Am Soc Nephrol. 2009;20(3):629-637. doi:10.1681/ASN.2008030287. PMID 19158356. Abstract: "bedside calculation of 0.413*(height/serum creatinine)".
Range: ageYears 1–16 (349 children aged 1 to 16, CKiD); heightCm > 0; creatinineMgDL > 0. ageYears is taken only to enforce the Range.
Example: height 100 cm, Scr 0.413 → 100 (HAND-DERIVED).

### MED-RENAL-008 EGFRSchwartzBedsideSI(heightCm, creatinineUmolL, ageYears float64)  ÷ 88.4 then 007.

### MED-RENAL-009 FractionalExcretionSodium(urineSodiumMmolL, plasmaSodiumMmolL, urineCreatinine, plasmaCreatinine float64)
FENa (%) = (UNa × PCr)/(PNa × UCr) × 100. Creatinine inputs in ONE unit, any (the ratio cancels it): document that.
Source: Espinel CH. The FENa test. Use in the differential diagnosis of acute renal failure. JAMA. 1976;236(6):579-581. doi:10.1001/jama.236.6.579. PMID 947239.
Range: all > 0 (population: oliguric acute renal failure; no input limits stated). No category function (the paper reports group findings, not a classification table); the doc comment may quote the abstract sentence with its citation, nothing more.
Example: UNa 20, PNa 140, UCr 100, PCr 2 → (20×2)/(140×100)×100 = 0.285714… (HAND-DERIVED).

### MED-RENAL-010 FractionalExcretionUrea(urineUrea, plasmaUrea, urineCreatinine, plasmaCreatinine float64)
FEUrea (%) = (UUrea × PCr)/(PUrea × UCr) × 100. Urea pair in one unit, creatinine pair in one unit.
Source: Carvounis CP, Nisar S, Guro-Razuman S. Significance of the fractional excretion of urea in the differential diagnosis of acute renal failure. Kidney Int. 2002;62(6):2223-2229. doi:10.1046/j.1523-1755.2002.00683.x. PMID 12427149.
Range: all > 0. No category function.

### MED-RENAL-011 GFRCategoryKDIGO(egfrMLMin173 float64) (GFRCategory, error)
`type GFRCategory string` with constants G1 "G1", G2 "G2", G3a "G3a", G3b "G3b", G4 "G4", G5 "G5", and
`func (c GFRCategory) Term() string` returning the source's term verbatim:
G1 "Normal or high" (≥90); G2 "Mildly decreased" (60–89); G3a "Mildly to moderately decreased" (45–59);
G3b "Moderately to severely decreased" (30–44); G4 "Severely decreased" (15–29); G5 "Kidney failure" (<15).
Edges (ASSUMPTION A-M3, the table is printed in integers): ≥90 G1, ≥60 G2, ≥45 G3a, ≥30 G3b, ≥15 G4, else G5.
Source: Kidney Disease: Improving Global Outcomes (KDIGO) CKD Work Group. KDIGO 2024 Clinical Practice Guideline for the Evaluation and Management of Chronic Kidney Disease. Kidney Int. 2024;105(4S):S117–S314. doi:10.1016/j.kint.2023.10.018. p. S125, "Current chronic kidney disease (CKD) nomenclature used by KDIGO", GFR categories.
Range: egfr ≥ 0. Tests: both sides of 90, 60, 45, 30, 15 (e.g. 89.999 and 90).

## med/chemistry (SHORT CHEM)

### MED-CHEM-001 CorrectedCalciumPayne(calciumMgDL, albuminGDL float64)
Adjusted calcium (mg/dL) = calcium − albumin + 4.0 (calcium mg/100 mL, albumin g/100 mL).
Source: Payne RB, Little AJ, Williams RB, Milner JR. Interpretation of serum calcium in patients with abnormal serum proteins. Br Med J. 1973;4(5893):643-646. doi:10.1136/bmj.4.5893.643. PMID 4758544. Abstract: "Adjusted calcium = calcium − albumin + 4.0, where calcium is in mg/100 ml and albumin in g/100 ml."
Range: both > 0 (population: 200 consecutive liver-function-test specimens; limits not stated).
NOT the 0.8 coefficient form (no primary source found; decision D-M6).
Example: Ca 8.0, albumin 2.0 → 10.0 (HAND-DERIVED).

### MED-CHEM-002 CorrectedCalciumPayneSI(calciumMmolL, albuminGL float64) (float64 mmol/L, error)
Convert: calcium mg/dL = mmol/L ÷ 0.2495 (KDIGO 2024 p. S126); albumin g/dL = g/L ÷ 10 (definition); apply 001;
convert back × 0.2495. Second Source line KDIGO 2024 p. S126. Test agreement with 001.

### MED-CHEM-003 SerumOsmolalityDorwart(sodiumMmolL, glucoseMgDL, ureaNitrogenMgDL float64)
Osmolality (mOsm/kg) = 1.86 × Na + glucose/18 + BUN/2.8 + 9.
Source: Dorwart WV, Chalmers L. Comparison of methods for calculating serum osmolality from chemical concentrations, and the prognostic value of such calculations. Clin Chem. 1975;21(2):190-194. PMID 1112025. Abstract: "Osmolality = 1.86 Na + (Glu/18) + (BUN/2.8) + 9."
Range: all > 0 (715 hospital sera; limits not stated). Example: Na 140, Glu 90, BUN 14 → 260.4 + 5 + 5 + 9 = 279.4.

### MED-CHEM-004 SerumOsmolalityDorwartSI(sodiumMmolL, glucoseMmolL, ureaMmolL float64)
= 1.86 × Na + glucose + urea + 9, i.e. the same equation with the paper's own divisors 18 and 2.8 taken as its
mg/dL→mmol/L conversions (document this as the reading; it is exact algebra on the paper's equation).
Test: Glu 90 mg/dL ↔ 5 mmol/L, BUN 14 ↔ 5 → same 279.4.

### MED-CHEM-005 EstimatedAverageGlucose(a1cPercent float64) (mg/dL)
eAG (mg/dL) = 28.7 × A1C − 46.7.
Source: Nathan DM, Kuenen J, Borg R, Zheng H, Schoenfeld D, Heine RJ; A1c-Derived Average Glucose Study Group. Translating the A1C assay into estimated average glucose values. Diabetes Care. 2008;31(8):1473-1478. doi:10.2337/dc08-0545. PMID 18540046. Results: "AG(mg/dl) = 28.7 × A1C − 46.7" (full text PMC2742903).
Range: a1cPercent 4–12 (ASSUMPTION A-M4: the study's lowest reported A1C band starts at 4%; its table tabulates 5–12; no upper limit stated).
Test the paper's table (rounded to integer): 5→97, 6→126, 7→154, 8→183, 9→212, 10→240, 11→269, 12→298 (compare math.Round(result)).

### MED-CHEM-006 EstimatedAverageGlucoseMmolL(a1cPercent float64)
eAG (mmol/L) = 1.5944 × A1C − 2.594 (full text). Same source. Table (1 decimal): 5→5.4, 6→7.0, 7→8.6, 8→10.2, 9→11.8, 10→13.4, 11→14.9, 12→16.5.

### MED-CHEM-007 AdjustedAnionGapFigge(observedAnionGapMmolL, albuminGL, normalAlbuminGL float64)
Adjusted AG = observed AG + 0.25 × (normal albumin − observed albumin), albumin in g/L.
Source: Figge J, Jabor A, Kazda A, Fencl V. Anion gap and hypoalbuminemia. Crit Care Med. 1998;26(11):1807-1810. doi:10.1097/00003246-199811000-00019. PMID 9824071. Abstract quoted. The paper's abstract gives no number for normal albumin, so the caller passes it.
Range: albuminGL > 0, normalAlbuminGL > 0, observed AG finite (may be negative? allow any finite).
### MED-CHEM-008 AdjustedAnionGapFiggeGDL(observedAnionGapMmolL, albuminGDL, normalAlbuminGDL float64) factor 2.5 ("if given in g/dL, the factor is 2.5"). Test agreement with 007.

### MED-CHEM-009 LDLCholesterolMartinHopkins(totalCholesterolMgDL, hdlCholesterolMgDL, triglyceridesMgDL float64) (mg/dL)
LDL-C = non-HDL-C − TG/factor, non-HDL-C = TC − HDL-C; factor from the 180-cell table below (rows TG strata,
columns non-HDL-C strata <100, 100–129, 130–159, 160–189, 190–219, ≥220).
Source: Martin SS, Blaha MJ, Elshazly MB, et al. Comparison of a novel method vs the Friedewald equation for estimating low-density lipoprotein cholesterol levels from the standard lipid profile. JAMA. 2013;310(19):2061-2068. doi:10.1001/jama.2013.280532. PMID 24240933. Figure 2 (180-cell table; PMC4226221, transcribed twice independently 2026-11-12, all cells agree).
Stratum rule (ASSUMPTION A-M5): a value belongs to the first stratum whose printed upper bound it does not exceed
(TG 49.5 → "50-56"; non-HDL 99.5 → "100-129"), the table's integer strata read as half-open at the top.
Range: triglyceridesMgDL 7–13975 (the table's span); non-HDL-C > 0 (TC > HDL). Embed the table as a [30][6]float64 plus the TG upper bounds [30]float64{49,56,61,66,71,75,79,83,87,92,96,100,105,110,115,120,126,132,138,146,154,163,173,185,201,220,247,292,399,13975} and non-HDL upper bounds {99,129,159,189,219,+Inf}.
Table (rows in that TG order; columns in that non-HDL order):
3.5 3.4 3.3 3.3 3.2 3.1 / 4.0 3.9 3.7 3.6 3.6 3.4 / 4.3 4.1 4.0 3.9 3.8 3.6 / 4.5 4.3 4.1 4.0 3.9 3.9 /
4.7 4.4 4.3 4.2 4.1 3.9 / 4.8 4.6 4.4 4.2 4.2 4.1 / 4.9 4.6 4.5 4.3 4.3 4.2 / 5.0 4.8 4.6 4.4 4.3 4.2 /
5.1 4.8 4.6 4.5 4.4 4.3 / 5.2 4.9 4.7 4.6 4.4 4.3 / 5.3 5.0 4.8 4.7 4.5 4.4 / 5.4 5.1 4.8 4.7 4.5 4.3 /
5.5 5.2 5.0 4.7 4.6 4.5 / 5.6 5.3 5.0 4.8 4.6 4.5 / 5.7 5.4 5.1 4.9 4.7 4.5 / 5.8 5.5 5.2 5.0 4.8 4.6 /
6.0 5.5 5.3 5.0 4.8 4.6 / 6.1 5.7 5.3 5.1 4.9 4.7 / 6.2 5.8 5.4 5.2 5.0 4.7 / 6.3 5.9 5.6 5.3 5.0 4.8 /
6.5 6.0 5.7 5.4 5.1 4.8 / 6.7 6.2 5.8 5.4 5.2 4.9 / 6.8 6.3 5.9 5.5 5.3 5.0 / 7.0 6.5 6.0 5.7 5.4 5.1 /
7.3 6.7 6.2 5.8 5.5 5.2 / 7.6 6.9 6.4 6.0 5.6 5.3 / 8.0 7.2 6.6 6.2 5.9 5.4 / 8.5 7.6 7.0 6.5 6.1 5.6 /
9.5 8.3 7.5 7.0 6.5 5.9 / 11.9 10.0 8.8 8.1 7.5 6.7
Tests: every TG row boundary on both sides for one column; every column boundary; corner cells (3.5, 3.1, 11.9, 6.7);
example TC 200, HDL 50, TG 150 → non-HDL 150, factor 5.7 → 150 − 150/5.7 (HAND-DERIVED).

### MED-CHEM-010 LDLCholesterolMartinHopkinsSI(totalCholesterolMmolL, hdlCholesterolMmolL, triglyceridesMmolL float64) (mmol/L)
Convert with the paper's own caption factors (Figure 2 caption: "To convert HDL-C, LDL-C, and total cholesterol, multiply by 0.0259; triglycerides, multiply by 0.0113"): mg/dL = mmol/L ÷ 0.0259 (cholesterol), ÷ 0.0113 (TG); compute 009; × 0.0259 back. Range: the converted TG within 7–13975.

## med/hepatic (SHORT HEPAT)

### MED-HEPAT-001 FIB4(ageYears, astUL, altUL, plateletsPer10e9L float64)
FIB-4 = (age × AST) / (platelets × √ALT).
Source: Sterling RK, Lissen E, Clumeck N, et al. Development of a simple noninvasive index to predict significant fibrosis in patients with HIV/HCV coinfection. Hepatology. 2006;43(6):1317-1325. doi:10.1002/hep.21178. PMID 16729309. Abstract formula "age ([yr] x AST [U/L]) / ((PLT [10^9/L]) x (ALT [U/L])^(1/2))".
Range: all > 0; age ≥ 18 is NOT stated in the abstract — enforce ageYears > 0 only and say "population: 832 HIV/HCV-coinfected adults (APRICOT); input limits not stated in the source". No category function; the doc may quote the abstract's two cutoff sentences verbatim with the citation.
Example: age 50, AST 40, ALT 25, PLT 200 → 50×40/(200×5) = 2.0 (HAND-DERIVED).

### MED-HEPAT-002 MELD(bilirubinMgDL, inrRatio, creatinineMgDL float64, dialysisTwiceInPastWeek bool)
MELD = 9.57 × ln(creatinine) + 3.78 × ln(bilirubin) + 11.20 × ln(INR) + 6.43; values below 1 set to 1; creatinine
capped at 4.0, and set to 4.0 when dialysisTwiceInPastWeek. Not rounded (the quoted form states no rounding).
Source: Kim WR, Mannalithara A, Heimbach JK, et al. MELD 3.0: the Model for End-Stage Liver Disease updated for the modern era. Gastroenterology. 2021;161(6):1887-1895.e4. doi:10.1053/j.gastro.2021.08.050. PMID 34481845. Methods (full text PMC8608337), the MELD formula as quoted: "MELD = 9.57 × loge(creatinine) + 3.78 × loge(bilirubin) + 11.20 × loge(INR) + 6.43".
Source (second line): Kamath PS, Wiesner RH, Malinchoc M, et al. A model to predict survival in patients with end-stage liver disease. Hepatology. 2001;33(2):464-470. doi:10.1053/jhep.2001.22172. PMID 11172350 (the original model; its own coefficient form not read).
Range: all > 0. Example: all inputs 1 → 6.43.

### MED-HEPAT-003 MELDNa(bilirubinMgDL, inrRatio, creatinineMgDL, sodiumMmolL float64, dialysisTwiceInPastWeek bool)
MELDNa = MELD + 1.32 × (137 − Na) − 0.033 × MELD × (137 − Na), Na bounded to 125–137 (clamped, the source's rule), MELD from 002. Same Kim 2021 source (as quoted there). Range: sodium > 0.
Example: MELD 6.43 (inputs 1), Na 137 → 6.43; Na 125 → 6.43 + 15.84 − 0.033×6.43×12.

### MED-HEPAT-004 MELD3(female bool, bilirubinMgDL, sodiumMmolL, inrRatio, creatinineMgDL, albuminGDL float64) (int, error)
MELD 3.0 = 1.33 (if female) + 4.56 × ln(bilirubin) + 0.82 × (137 − Na) − 0.24 × (137 − Na) × ln(bilirubin) + 9.09 × ln(INR) + 11.14 × ln(creatinine) + 1.85 × (3.5 − albumin) − 1.83 × (3.5 − albumin) × ln(creatinine) + 6, rounded to the nearest integer.
Bounds (the source's): bilirubin, INR, creatinine below 1.0 set to 1.0; creatinine upper bound 3.0; Na bounded 125–137; albumin bounded 1.5–3.5.
Source: Kim 2021 as 002, Methods/Results equation (verbatim above).
Rounding: math.Round (half away from zero) — document it ("rounded to the nearest integer"; tie rule not stated; A-M6).
Range: all > 0. Example: male, all lab values at their lower bounds (bili 1, INR 1, Cr 1, Na 137, albumin 3.5) → 6.
