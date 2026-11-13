# PORT candidates from the AdvoMedix knowledgebase (read-only inventory, 2026-11-12)

Reference: `/Users/devinanderson/Desktop/advo_production/advomedix-v2/backendAdvo/knowledgebase` — root, `series/`,
`calculators/`, `formula/`, `ranges/`; `pending/` and `analysis/` not opened; no env, vendor or node_modules file opened.
Nothing there was edited. No porting was done in this run (a later slice, with its own tests).

Shape of the reference: Go. 692 JSON entries in `calculators/<category>/<id>.json`; 128 have a Go formula registered
through `formula/formula.go:31`; range checks for every runnable formula live in one place, `run.go` checkInputs
(~lines 72–110, JSON [min,max], no unit conversion). The other 564 entries are cite-only JSON (inputs, citation,
pointer to the mobile `.tsx`), with no backend formula. `series/` is plain statistics, no clinical cut-offs.

Column "calcbench": the ID that replaces it, NEW (no calcbench function yet), or HELD (catalogued here, held back).
"Finding" is a difference between the reference and its cited source, for Devin; nothing was changed there.

| Reference file:line | What it computes | calcbench | Finding |
|---|---|---|---|
| calculators/metabolic/bmi_bsa.go:26 | BMI, lb×703/in² | MED-ANTHRO-002 (and -001 metric) | inputs 50–500 lb, 36–96 in in JSON; calcbench takes any positive value (CDC states no limits) |
| calculators/metabolic/bmi_bsa.go:27 | BSA Mosteller (converts lb/in inline) | HELD (Mosteller 1987 letter not seen) | inline unit conversion 0.453592 / 2.54 |
| calculators/metabolic/bmi_bsa.go:28 | BSA Du Bois | HELD (Du Bois 1916 formula not seen in the paper) | — |
| calculators/metabolic/ideal_body_weight.go:31-43 | Devine IBW + tidal volume 6 and 8 mL/kg | HELD (Devine IBW: McCarron & Devine 1974 not seen); tidal volume NEW | extrapolates below 60 in by subtracting 2.3 kg/in, outside Devine's stated use |
| calculators/metabolic/basal_energy_expenditure.go:29-39 | Harris–Benedict BEE/TDEE | NEW | male age coefficient 6.775 vs the published 6.7550; constants rounded (66.5, 655.1) |
| calculators/nutrition/mifflin_st_jeor_bmr.go:20 | Mifflin–St Jeor REE | NEW | matches the abstract |
| calculators/metabolic/body_roundness_index.go:21-28 | Body roundness index (Thomas 2013) | NEW | — |
| calculators/general/relative_fat_mass.go:20 | Relative fat mass (Woolcott 2018) | NEW | — |
| calculators/metabolic/osiris_osteoporosis.go:20-31 | OSIRIS | NEW | — |
| calculators/general/ost_osteoporosis.go:23-25 | OST | NEW | — |
| calculators/metabolic/bs_eag.go:18-19 | A1C → eAG (28.7/46.7; 1.59/2.59) | MED-CHEM-005, -006 | the mmol/L form uses the rounded 1.59/2.59; ADAG full text gives 1.5944/2.594 |
| calculators/cardiac/map_calc.go:17-21 | MAP (SBP + 2·DBP)/3 | HELD (one-third rule: no primary seen) | — |
| calculators/cardiac/pulse_pressure.go:17-25 | pulse pressure, PPI, narrow threshold SBP×0.25 | NEW (PP); threshold has no cited primary | — |
| calculators/cardiac/max_hr.go:18-30 | Fox, Tanaka, Gulati, HUNT max HR | MED-CARDIO-001 (Tanaka); others NEW | the "recommended" chooser (Gulati if female, Tanaka ≥40 else Fox) is the reference's own rule, not a source's |
| calculators/cardiac/hr_zones.go:21-38 | Karvonen zones | HELD (Karvonen 1957 not seen) | — |
| calculators/cardiac/cri.go:21-31 | chronotropic response index | NEW | source marked unverified in the reference |
| calculators/cardiac/hrr1.go:20 | heart-rate recovery at 1 min | NEW | — |
| calculators/cardiac/rpp.go:23 | rate–pressure product HR×SBP | HELD (Gobel 1978 abstract says "HR x blood pressure", not systolic) | — |
| calculators/cardiac/orthostatic.go:17-18 | orthostatic drop | NEW | — |
| calculators/cardiac/aap_peds_htn.go:25-87 | AAP 2017 pediatric BP level | NEW | height adjustment is a home-made approximation, not the AAP height-percentile tables |
| calculators/cardiac/chads2_score.go:16-24 | CHADS2 + stroke rates | NEW (med/scales) | rates match Gage 2001 |
| calculators/cardiac/chadsvasc.go:23-30 | CHA2DS2-VASc | NEW (med/scales) | — |
| calculators/cardiac/isth_bat.go:16-34 | ISTH-BAT | NEW (med/scales) | out-of-range item index scores 0 silently |
| calculators/respiratory/pef_predicted.go:17-31 | predicted PEF (US) | NEW | coefficients are not Hankinson 1999's (which are quadratic in age, height², L/s) |
| calculators/respiratory/pef_predicted_eu.go:17-29 | predicted PEF (EU) | NEW | linear form is Gregg & Nunn 1973, cited as Nunn & Gregg 1989 |
| calculators/respiratory/six_minute_walk.go:18-29 | 6MWD Enright & Sherrill 1998 | NEW | matches |
| calculators/respiratory/pack_years.go:27-31 | pack-years | NEW | — |
| calculators/respiratory/spo2_classify.go:17-48 | NEWS2 SpO2 points | NEW | scale 2 ignores air vs oxygen (SpO₂ ≥ 93 on air scores 0 in NEWS2) |
| calculators/general/news2.go:20-131 | NEWS2 total | NEW (med/scales) | same scale-2 finding (lines 36-50) |
| calculators/general/temp_liebermeister.go:16-25 | Liebermeister expected HR | NEW | 10 bpm per °C matches neither Liebermeister (~8/°C) nor Cunha (10/°F) |
| calculators/general/temp_mean_body.go:16-17 | mean body temperature (Burton 1935) | NEW | — |
| calculators/metabolic/ausdrisk.go:27-60 | AUSDRISK | NEW | the fruit/vegetable item is missing (max 37 instead of 38) |
| calculators/obstetric/pregnancy_dates.go:29-38 | Naegele dates | NEW (needs time.Time policy) | — |
| calculators/renal/cockcroft_gault.json (cite-only) | Cockcroft–Gault | MED-RENAL-005, -006 | no backend formula today |
| calculators/renal/* cite-only (MDRD, revised Schwartz, CKiD U25) | eGFR | MED-RENAL-007/-008 (bedside Schwartz); MDRD HELD | — |
| cite-only: QTc, anion gap, corrected Ca/Na, osmolality, Martin LDL, FIB-4, MELD, GCS, Apgar, Holliday–Segar, dosing | various | MED-CARDIO-002, MED-CHEM-001…010, MED-HEPAT-001…004, MED-SCALE-003…006; HELD: anion gap, corrected Na, Friedewald, Holliday–Segar, dosing | — |

Counts: 35 rows; 15 map to a calcbench ID now (fully or partly), 7 are HELD, the rest NEW. Not computed anywhere in the
reference: CKD-EPI (any year), Epworth, APRI, Parkland, Friedewald, VO2max, LMS / BMI-for-age tables.

What a consumer (AdvoMedix backend) must change to use calcbench — a list, never an edit:
1. Import `github.com/CrimsonCoder42/calcbench/med/<pkg>` and call the function in place of the formula body.
2. Pass inputs in the unit the calcbench function names (e.g. `weightKg`, `heightCm`), converting once at the
   boundary; the reference's inline conversions (bmi_bsa.go:27) go away.
3. Map calcbench errors (`ErrOutOfRange`, `ErrNotFinite`, `ErrNotTestable`) to the reference's run.go error path;
   calcbench Ranges are the sources' and can be narrower than the JSON [min,max] (e.g. Cockcroft–Gault ages 18–92).
4. Category labels are the sources' verbatim strings (CDC, KDIGO, Epworth site, Rimel, Apgar 1953); the
   `ranges/` JSON bands are separate and are not replaced by calcbench.
