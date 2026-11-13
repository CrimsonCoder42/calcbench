# CONSUMER DEFECTS (port wave 2)

Entries where the consumer shows a number its cited source does not support: what the entry computes, what the paper says, the locator, the severity.

### W1_B1
- aap_peds_htn: computes a home-made height adjustment ((height − median)/2.5 × 2 or × 1 mmHg) on 12-row age tables; Flynn 2017 (PMID 28827377) uses sex/age/height-percentile Tables 4–5. Severity high.
- diastolic_strain: shows a combined "diastolic strain index"; Nagueh 2016 (PMID 27037982) defines only a decision algorithm, no index. Severity medium.
- rescue_ihca: shows a mortality % from logit −1.5 + 0.08 × score ("approximate", self-described); Tonna 2022 Table 2 (PMID 35033471) prints no such formula. Severity high.
- cart_score: shows bands ≥ 8 moderate, ≥ 12 high, ≥ 20 very high; Churpek 2012 Table 3 / Results (PMID 22584764) report only the cut-off > 17. Severity medium.
- modified_shock_index: shows category labels (hyperdynamic / borderline / hypodynamic) and derives MAP by the one-third rule; Liu 2012 (PMID 25215048) gives only thresholds 0.7 and 1.3. Severity low.
- heart_pathway: scores age 65 as 2 points (≥ 65) and checks a single troponin; Mahler 2015 Table 2 / Methods (PMID 25737484): "> 65" for 2 points, negative troponin at 0 and 3 hours. Severity medium.
- garfield_af: computes risks from baseline survivals and coefficients not printed in Fox 2017 (PMID 29273652; HRs only, online supplementary table 2). Severity medium.
- isth_bat: an out-of-range item index scores 0 silently (isth_bat.go:16-34). Severity medium.
- rasi_calc: cites a Mayo Clinic stress-test page for RASI (wrong page); formula equals MED-CARDIO-033. Severity low.
- cpes_score: cites PMID 37100484 (an unrelated JACC TAVR editorial). Severity low.
- pcp_hf: offers "White/Other" race; Khan 2019 (PMID 31097157) derived White and Black equations only. Severity low.
- pulse_pressure: outputs 0.25 × systolic with no source found. Severity low.
- kawasaki: result text includes IVIG and aspirin doses (dose content in a criteria tool). Severity low.

### W1_B2
| entry | what it computes | what the paper says | locator | severity |
|---|---|---|---|---|
| circ_covid | additive points "approximating" a Cox model (age, RR, BMI, Hb, pulse, albumin …) | CIRC is a multistate model; no point values | Garibaldi 2021 Ann Intern Med 174:33, PMID 32960645 | high |
| vaco_index | age base rate × 0.5 (female) × CCI multiplier | logistic model with CCI × age interaction; weights only in S1 File | King 2020 PLoS One 15:e0241825, PMC7657526 | high |
| utah_covid_risk | omits Male (1) and Communities of color (2); immunocompromised 1; ages 0–10 = 1.0 | Table 6: Male 1, Communities of color 2, Severely Immunocompromised 2, 0–10 = 0.5 | Webb 2022 PLoS One, PMC8893609, Table 6 | medium |
| pef_predicted | linear −130.15 + 5.67·cm − 5.66·age / 213.41 + 2.80·cm − 4.91·age, cited to Hankinson 1999 | Hankinson: quadratic in age, height², L/s, by sex and race | PMID 9872837, Tables 4–6 | medium |
| pompe_c | RespDist 0.799, SaO2 −0.0636 | equation prints 0.800 and −0.063 (8.42 % vs 7.98 % in the test case) | Kline 2012 Thromb Res, PMC3350789, equation | low |
| ottawa_copd | risk 85.0 % for score 9 | Table 8: score 9 "NA" | Stiell 2014 CMAJ, PMC3971051, Table 8 | low |
| dissociation | Hill curve P50 26.6, n 2.7 attributed to Severinghaus | cited page has no equation; Severinghaus 1979 is not a Hill form | PMID 35496 | low |
| fio2_estimate | 0.21 + 0.04·flow | no source found | cited Cleveland Clinic page | low |
| pef_flowmetric | 5.3·cm − 433 / (cm − 100)·5 + 100 | no source found | cited Mayo Clinic page | low |
| roth_score | "High/Moderate/Lower Risk" labels | paper gives sensitivity cut-offs only | Chorin 2016 Clin Cardiol, PMC6490799 | low |
| drip_pneumonia | "Very High Risk" at ≥ 6 | paper has the single ≥ 4 threshold | Webb 2016, PMC4862530 | low |

### W1_B3
- calcium_correction: computes Ca + 0.8 × (4.0 − albumin) citing Payne 1973; Payne says Ca − albumin + 4.0 (coefficient 1.0); calcium_correction.tsx:50. Severity medium.
- bs_corrected_na: above glucose 400 mg/dL switches to 2.4 per 100 mg/dL (calculators.tsx:1361) and calls it Hillier's; Hillier 1999 abstract says the mean was 2.4 overall and "a correction factor of 4.0 was better" above 400. Severity medium.
- bs_osmolality: "total" = 2Na + glucose/18 + BUN/2.8 (calculators.tsx:1483) cited to Kitabchi 2009, whose total form is 2(Na + K) + urea + glucose in mmol/L (Table 2 footnote). Severity low.
- bs_homa_ir / homa_ir: bands 1.0 / 2.0 / 3.0 with labels ("Optimal", "Early Insulin Resistance", ...) (calculators.tsx:830-849) appear in no cited source. Severity low.
- bs_cgm_metrics: TIR target tested as ≥ 70 (calculators.tsx:694); Battelino Table 3 says ">70%"; "High Risk"/"Below Goal" labels not in the source. Severity low.
- bs_a1c_from_eag: inverts the ADAG regression (calculators.tsx:486); Nathan 2008 publishes eAG on A1C only. Severity low.
- serum_anion_gap: delta gap / delta ratio with normal AG 12 and HCO3 24 (serum_anion_gap.tsx:35-37) have no cited source. Severity low.
- rebound_hyperbilirubinemia: citation URL points to PMID 28219964 (an unrelated article); correct PMID 28196932. Severity low.

### W1_B4
- cdiva: 4 items scored 0–2; van Loon 2016 A-DIVA has 5 predictors 1 point each (PMC4845841 Table 5). medium.
- ethanol_serum_concentration: ethanol mass = %×10×mL with no density (0.789 g/mL), Vd 0.6 L/kg for all alcohols; the cited fomepizole review (PMID 19458366) gives no such formula. medium.
- ace_score: sums 10 categories citing Felitti 1998, which defines 7 (abstract "range: 0-7"); bands "Intermediate Risk/High Risk/Severely Elevated" in no source. medium.
- esas_r: bands at 18/36/63 "Symptom Burden" in no source; Watanabe 2011 abstract defines no total. medium.
- katz_adl: count of independent functions mapped to Katz grades A–G; Katz grades C–F name specific functions (to confirm with the 1963 PDF). medium.
- brue_criteria: age > 2 months instead of > 60 days; concerning history/examination features omitted (Nama 2025 PMC11791710 restates the AAP higher-risk list). medium.
- braden_score: bands ≤ 12 / ≤ 15 / ≤ 17 match neither secondary (≤ 9, 10–12, 13–14, 15–18). low.
- abuse_assessment_screen: bands "Two Concerns"/"Multiple Concerns" in no source. low.
- duval_cibmtr: 3-year OS 28% (score 1) and 15% (score 2) not in Table 5 or the abstract (only 42% for 0, 6% for ≥ 3). low.
- hiri_msm: insertive anal sex with HIV-positive partner scored on any yes; Wilton 2016 Table 2 scores ≥ 5 times. low.
- hiv_rasp: citation URL PMID 17688678 is an unrelated Breast Cancer Res editorial; Vertesi 2003 is PMID 17659153. low.
- icc_systemic_mastocytosis: no SM-AMN exception for the tryptase minor criterion (Arber 2022 Table 11). low.
- age_adjusted_esr_crp: CRP formula applied at ages 18–100; Wener 2000 proposes it for 25–70. low.
- crs_grading: no path for grading after antipyretic or anti-cytokine therapy without fever (Lee 2019 text). low.
- breast_milk_storage: citation link PMID 29958002 is the 2018 correction, not the protocol (PMID 29624432). low.

- braden_score: Bergstrom N, Braden BJ, Laguzza A, Holman V. The Braden Scale for predicting pressure sore risk. Nurs Res. 1987;36(4):205-210. PMID 3299278 (added by the lead, D-075)
- mentzer_index: Mentzer WC Jr. Differentiation of iron deficiency from thalassaemia trait. Lancet. 1973;1(7808):882. doi:10.1016/s0140-6736(73)91446-3. PMID 4123424 (added by the lead, D-075)

### W2_B5
- pews: adds 1 for quarter-hourly nebulisers and 1 for post-operative vomiting; sources give 2 each (Shafi 2020 PMC7719480 "Score 2 extra"); citation URL PMID 15960531 is a J Org Chem paper (true PMID 15751446). medium.
- parc_appendicitis: age coefficients (girl 8–11 +0.6513, boy 8–13 −0.0654) and a community recalibration (−0.615 + 1.1x) not in Kharbanda 2018 Table 3 as rendered (age 8–13.9 −.72). medium (to confirm with the supplement).
- myxedema_coma: at-risk band from 25; primary abstract (PMID 24518183) says 45–59. medium (provisional).
- glasgow_blatchford: BUN band "> 70 mg/dL" for 6 points; urea ≥ 25 mmol/L (= 70 mg/dL) in Adrych 2025 / Orpen-Palmer 2023 reading. low.
- temp_abg_correction: pH slope term 0.065 × (pH − 7.40); usually printed 0.0065 (to confirm with Ashwood 1983). medium (provisional).
- temp_classifier / temp_clinical_scores: "Normal" age bands from a patient-information page, no primary. medium.
- temp_q10: fixed Q10 = 2.5 with no primary. low. temp_metabolic_cost: 2000 kcal/day baseline and 10 %/°C with no primary. low.
- temp_liebermeister: cites Cunha's 10 bpm per °F rule but applies 10 bpm per °C. low.
- temp_mean_body: bands "Low/Normal/Elevated Mean Body Temperature" (35/38 °C) in no source. low.
- smast_g: bands 2–5 "Suggestive", ≥ 5 "Strongly Suggestive" in no source (sources: ≥ 2 positive). low. smast: cites the 25-item MAST (1971) for the 13-item SMAST (1975). low.
- wast: 8-item cut-off ≥ 13 and short-form ≥ 5 not in the CDC compendium (short form: 0–2, cut-off 1). low.
- ost_osteoporosis: lb ÷ 2.205 instead of × 0.45359237 (no case changes); women's bands applied to men (Adler 2003 male bands differ). low.
- revised_trauma_score: cites the 1981 Trauma Score for the 1989 Revised Trauma Score. low. ovat_ipv: cites the 2003 OAS paper for the 2004 OVAT. low. nyha: cites Goldman 1981 (Specific Activity Scale). low.
- nutritional_risk_index: cites PMID 3124597 (companion paper) for the protocol (PMID 3124598). low.

### W2_B6
- canadian_tia_score: 13 items (unilateral weakness 3, prior TIA 1, multiple TIA 2, carotid bruit 1, imaging ischemia 3, symptomatic carotid 3, non-cerebrovascular −3, …) and bands; the Canadian TIA Score (Perry 2021 BMJ Table 1, PMC7859838) has different items (antiplatelet 3, DBP ≥110 3, glucose ≥15 3, platelets ≥400 2, vertigo −3, …). high.
- mnihss: sums 1a LOC and dysarthria, rates sensory 0–2, omits 1c LOC commands; Meyer & Lyden 2009 Table 1 (PMC2729912): those items removed, sensory 0–1, 1c retained, total 31. high.
- ed_safe_pss3: positive on any lifetime attempt; adds "Indeterminate"; Boudreaux 2013 (PMC3979300 Methods): positive = active ideation or attempt within the past 6 months. medium.
- cssrs: High/Moderate/Low/No Risk tiers from q1–q6; Posner 2011 (PMC3893686) has a 1–5 ideation severity rating and no tiers. medium.
- brief_addiction_monitor: Use/Risk/Protective subscales and total; Cacciola 2013 (PMC3602977) defines none. medium.
- modified_fisher_sah: no SAH with IVH → grade 0; neither source grades that case (Patel 2021 Table 6, Zhou 2022 §2.3). low.
- chalice_rule (PMID 16861481 → correct 17056862), gmaws (22345174 → 22328545), hii_impairment (24405226 is Ahamad 2014), midas_migraine (cites the HImQ paper 10068175; MIDAS is 11294956): wrong citations. low.
- gds15, dast_10, mtoq4, cam_icu, finnegan_nas, infant_scalp_score: cited paper is a different instrument or a review; corrected citations in the rows above. low.
- cage_alcohol, cas_covid_anxiety, audit_c, aspects_stroke, four_score_coma, lams, new_orleans_head_ct, nexus_head_ct, ottawa_sah_rule, nps_pain: bands or labels with no source (rows above). low.
- cows_score: shows "Buprenorphine Induction: Appropriate" at ≥ 12 (treatment advice; not in the scale, Wesson & Ling 2003 Appendix 1). low.

### W2_B7
- hiv_ckd_prediction: age 1 point per decade (max 4, adds ≥ 70) and 1 point for glucose/hypertension/proteinuria; Scherzer 2014 Table 2 gives 0/2/4/6 and 2 each (PMC4188545). HIGH.
- pediatric_sirs_sepsis: SBP column 59/69/75/74/83/90; Goldstein 2005 Table 3 as reproduced (Wynn 2014 PMC7122308; Schmoch 2021 PMC8238315) 65/75/100/94/105/117. HIGH.
- sirs: omits the PaCO2 < 32 mmHg and > 10% bands arms; Singer 2016 Box 1 includes both. MEDIUM.
- mcmahon_rhabdomyolysis: bands < 6 / ≥ 6 (86%/68%) / ≥ 10 (52%); McMahon 2013 Results report < 5 (2.3%) and > 10 (61.2%) (PMC5152583). MEDIUM.
- mess_score: age exactly 50 → 2; source 30–50 → 1 (Topal 2010; Okereke 2022). "No ischemia (0)" and 5–6 "borderline" not in the sources. MEDIUM.
- ruis_incontinence: score 8 = "Moderate"; all three secondaries 4–8 mild, 9–12 moderate (PMC6988664, PMC6174897, PMC6287547). MEDIUM.
- mews / mews_calc: HR 40, temperature 35.0 and a 38.5–38.9 band differ from Kim 2015 Table 1 (PMC4476665); bands > 8 / ≥ 7 seen nowhere. MEDIUM.
- neonatal_eos: antibiotic categories and exam LRs (0.41/5/21.2) differ from Puopolo 2011 (≥ 4 h) and Escobar 2014 (0.36/3.75/14.5); fixed incidence. MEDIUM.
- ain_risk: ln(prevalence/0.23) adjustment and 15/30/50% bands not in Moledina 2022. MEDIUM.
- bun_creatinine_ratio: bands 10/20/40 have no source (KDIGO 2012 CKD and AKI state none). MEDIUM.
- urine_output_fluid_balance: < 0.1 anuria, 1.0 under 40 kg, > 3.0 polyuria, "normal" 0.5–1.0 have no source; KDIGO AKI 2012 Rec 2.1.1 / Table 2 tie 0.5 and 0.3 mL/kg/h to durations (PMC4089654). MEDIUM.
- kt_v_dialysis: V = 0.6 × weight unsourced (Churchill 2021 PMC8240937 says 55% of post-dialysis weight). MEDIUM.
- triss: cites Boyd 1987 for coefficients that are Champion 1995's (PMID 7745669). LOW.
- step_by_step_febrile: leukocyturia and PCT ≥ 0.5 placed in intermediate; Rybak 2022 (PMC9359674) places them in high. LOW (one secondary).
- rhmp30: worked example (tsx:211) x = 0.494 / 62.1%; its own coefficients give 0.4434 / 60.9%. LOW.
- tlics: "None" morphology (0) not in Lee 2005 Table 1. LOW.
- pediatric_trauma_score: JSON applies_when 18–120 on a paediatric score. LOW.
- slics: cites Patel 2008 only; the instrument is Vaccaro 2007. LOW.
- mdrd_gfr: cites Levey 1999 (186) while computing Levey 2006 (175). LOW.

### W2_B8
- dapre_endurance: offers an "endurance" DAPRE variant; Knight 1979/1985 (PMID 507268, 4079735) and Greig 2020 (PMC7575491) describe no endurance variant. Severity low (entry is cite-only; Run refuses).
- riegel_race_prediction: uses exponent 1.06 as "Riegel 1981"; Blödorn 2025 (PMC12181339) says the 1981 exponent is 1.07732 and 1.06 is an online calculator's constant; Vickers 2016 (PMC5000509) gives 1.08 elite / 1.05–1.06 recreational men 40–70. Severity low (unconfirmed until the 1981 paper is read).
- iota_simple_rules: entry advertises "Malignancy Risk %" from the 2016 model but carries no coefficients (cite-only). Severity low.

### W3_B9a
- bp_child_* / bp_infant / bp_newborn / bp_adolescent (18): one sex-averaged band table ("50th–90th percentile values averaged across sexes"); Flynn 2017 classifies by sex, age and height percentile (Tables 3–5). medium.
- pulse_* / rr_* paediatric bands, pulse_classifier.go, rr_classifier.go: Cleveland Clinic / Iowa bands with Low/Normal/High; Fleming 2011 centiles differ (abstract medians: HR 127 at birth, 145 at ~1 month; RR 44 at birth, 26 at 2 y). medium.
- spo2_oxygenation_index_category: OI edges 5/15/25/40; PALICC 2015 uses 4/8/16 (PMC5253180). medium.
- spo2_adult/child/infant/older/neonatal: unsourced hypoxaemia categories (BTS 2017 gives targets only). medium.
- pulse_qtc_female/male: edges 440/460/470/500; AHA/ACCF/HRS 2009 uses ≥ 450 men / ≥ 460 women (per wave 1; PDF pending). medium.
- bs_random_adult/child: "Normal" 70–140, "Elevated" 141–199, "Critically High" ≥ 300 not in ADA (only ≥ 200 with symptoms). medium. bs_fasting_adult: "Critically High" ≥ 300 unsourced. low.
- temperature_oral_* (6): "Normal"/fever age bands from a patient page; hypothermia core-temperature table applied to oral readings. medium.
- infant_growth_* (8): single oz/in-per-month ranges for both sexes; WHO velocity standards are sex-specific percentiles (MED-PEDS-008 for weight 0–12 months). medium.
- infant_diaper_* / infant_feeding_* (16): edges read from parent-guidance prose. low.
- pulse_chronotropic_response_index_beta_blocker: < 0.62 vs Khan 2005 "≤ 62%". low.
- blood_pressure:all_readings.go: CV from an already-rounded SD (precision loss). low. ComputeTrend / all_readings.go / measurement_weight_trend.go: 2 % stable band, 50 mg/dL swing, CV > 15 %, 5 lb swing unsourced. low.
- bp_adult_*, bp_classifier.go: "Low (Hypotension)" < 90/60, crisis > 180/120, older-adult diastolic floor 55, paediatric 90/50, 120/80 unsourced. low.
- remaining band-only rows (CPP, PVR 3/5 WU, PPI, RPP, shock index, ABI 0.4, adult pulse/RR, athlete, HOMA/QUICKI bands, S/F 235, COPD SpO2, fever metabolic cost): bands with no primary. low.

### W3_B9b
| entry | what it computes | what the source says | locator | severity |
|---|---|---|---|---|
| pef_exacerbation_percent_personal_best / _predicted | bands < 33 / 33–< 50 / 50–< 75 / ≥ 75 "Mild Exacerbation", PEF alone | 50 % is in "PEF 33–50%", 75 % in "PEF >50–75%"; no mild level; life-threatening needs SpO2 <92 % plus a feature (May 2025 amendment) | SIGN 158 (2024/2025) Table 15 | medium |
| pef_variability_adult / _child | > 20 % "Elevated Variability — Poorly Controlled"; child threshold 13 % for ages 5–17 | GINA gives only > 10 % adults, > 13 % children 6–11; no 20 % band | GINA 2023 Box 1-2 | medium |
| pef_reversibility_category | increase ≥ 20 % only | "an increase or decrease … a change in PEF of at least 20%" | GINA 2023 §"How much variation…" | low |
| pef_percent_predicted_category | 50/80/100 % predicted bands | no primary; cited MDCalc | — | medium |
| pef_population_adult_female / _male / _sex_unspecified / child_5_to_18 | fixed L/min population bands | no primary; PEF reference values depend on age, height, sex (Hankinson 1999) | consumer health pages only | medium |
| pef_diurnal.go | rounds each day to 0.1 before averaging | GINA averages the daily values (no rounding) | GINA 2023 Box 1-2 footnote * | low |
| X2 diurnal.go:114 | (max PM − min AM) / mean; "mild" 10–20 %, "poor" > 20 % | day's highest − day's lowest; no 10–20/20 control bands | GINA 2023 Box 1-2 footnote * | medium |
| X3 mage.go:10 | excursion = turning point vs its immediate neighbour; fallback consecutive differences | peak-to-nadir between consecutive turning points, one direction (first qualifying excursion) (secondaries PMC5501494, PMC2945154) | Service 1970 (NEEDS PDF) | medium |
| bs_trend.go | 50 mg/dL swing, 2 % split-half stable band | no primary (as wave 1 W3_B12) | — | low |

### W3_B10
- C4 correlation.go:172 DetectSymptomCoOccurrence: CoOccurrenceRate = JointCount/min(CountA, CountB) is documented
  "0.0-1.0" but exceeds 1 (headache ×2 each near the one nausea → 2/1 = 2). No source defines the rate. Severity low.
- C6 unit_convert.go:13-21 KgToLbs/LbsToKg use 2.20462 instead of the exact 1/0.453 592 37 = 2.204 622 62…: relative
  error 1.2e-6 (0.0001 lb at 100 kg). NIST SP 811 App. B.8. Severity low.
- C7 variability.go:42 and bp_trend.go:435 cite Rothwell et al., Lancet Neurol 2010 for ARV; ARV is Mena et al.
  (J Hypertens 2005;23:505; JAHA 2017;6:e006895 eq. 1). Number unchanged; citation wrong. Severity low.
- C5 diurnal.go:34 labels evening ≥ morning as "non-dipping"; dipping is a nocturnal (sleep vs awake) ABPM concept, not
  a morning-vs-evening comparison (triage note; no source supports the AM/PM use). Severity medium (a clinical label
  shown on an unsupported comparison).

### W3_B11
- saps3_score: picker points disagree with Moreno 2005 Table 1 in most rows (AIDS 14, metastatic cancer 17, ER 8, other ICU 5, ward 6, age ≥ 70 lumped at 13, GCS 6–8 13, HR "Not available" 11, …); single-choice comorbidity/reason where the paper adds all that apply. PMC1315315 Tables 1–2. high.
- pecarn_cspine: applies_when 18–120 for a rule enrolled in children < 18 (Leonard 2024 Methods); output labels give imaging advice. PMC11261431. medium.
- nhfs: mortality constant 5.012 taken from a garbled abstract (Karres 2023 prints 5.0122); mortality bands ≥ 20/10/5 % unsourced. PMID 22728204. low.
- revised_trauma_score: cites the 1981 Trauma Score for the 1989 RTS; T-RTS bands (≤ 3, ≤ 10, 11, 12 with triage labels) unsourced. PMID 2657085. low.
- flacc_scale: bands ≥ 7 / ≥ 4 / ≥ 1 sourced only to Univ. Michigan instructions. low.
- painad: bands ≥ 7 / ≥ 4 / ≥ 1 differ from Malara 2016 (overlapping at 7). low.
