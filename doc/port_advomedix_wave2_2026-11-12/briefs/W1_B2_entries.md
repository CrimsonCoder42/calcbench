# Entry list W1_B2 (wave 2)

Source root (read-only): `/Users/devinanderson/Desktop/advo_production/advomedix-v2/backendAdvo/`. Paths are as wave 1 recorded them: `knowledgebase/calculators/<category>/...`, `<category>/...` (= knowledgebase/calculators/<category>/), `knowledgebase/ranges/...`, `calculators/<vital>/...`, or relative to `knowledgebase/series/` for series functions. The wave 1 reason is the starting point, not the answer.

| # | entry id | source path (wave 1) | wave 1 section | wave 1 outcome | wave 1 reason |
|---|---|---|---|---|---|
| 1 | aa_gradient | respiratory/aa_gradient.json:1 | W1_B4 | NOT PORTED | held: doc/catalog/medicine.md §1 row 8 (alveolar gas equation, HELD-SOURCE); cites StatPearls NBK482268 (secondary); equation not seen in a primary this pass |
| 2 | abg_analyzer | respiratory/abg_analyzer.json:1 | W1_B4 | NOT PORTED | Baillie 2008 Thorax 63:289 (PMID 18308967) rules of thumb: not opened this pass; interpretive algorithm, compensation rules not seen |
| 3 | acute_asthma_severity | respiratory/acute_asthma_severity.json:1 | W1_B4 | NOT PORTED | NAEPP EPR-3 severity table not opened (NHLBI site); classification table only via secondaries |
| 4 | ais_inhalation | respiratory/ais_inhalation.json:1 | W1_B4 | NOT PORTED | PMID 17211205 abstract read: grades 0-4 named, criteria not stated (bronchoscopic grading; Endorf 2007 is not the grading's origin) |
| 5 | asthma_predictive_index | respiratory/asthma_predictive_index.json:1 | W1_B4 | NOT PORTED | citation is Guilbert 2006 NEJM trial (PMID 16687711), not the mAPI derivation; abstract has no index criteria; mAPI primary not found this pass |
| 6 | bap65 | respiratory/bap65.json:1 | W1_B4 | NOT PORTED | PMID 19786679 abstract read: Tabak 2009 abstract states BUN > 25 mg/dL, acute mental status change, pulse > 109/min, age 65 splits lowest group, but no class I–V assignment |
| 7 | bcrss_covid | respiratory/bcrss_covid.json:1 | W1_B4 | NOT PORTED | PMID 32297727 abstract read: items/points/thresholds not stated in the abstract; full text not open (Europe PMC: not OA); algorithm (decision pathway) |
| 8 | blast_lung_injury | respiratory/blast_lung_injury.json:1 | W1_B4 | NOT PORTED | Pizov 1999 PMID 9925079 abstract: categories combine P/F (<60, 60–200, >200) with radiographic infiltrates and bronchopleural fistula/pneumothorax; not a single computed number and undefined when findings disagree |
| 9 | bode_index | respiratory/bode_index.json:1 | W1_B4 | NOT PORTED | PMID 14999112 abstract read: '10-point scale', four factors; point table not in abstract; NEJM full text not open |
| 10 | cao2 | respiratory/cao2.json:1 | W1_B4 | NOT PORTED | held: doc/catalog/medicine.md §1 §5 note: CaO2 1.34 vs 1.39 held; cites StatPearls (secondary); no primary seen |
| 11 | chosen_risk_score | respiratory/chosen_risk_score.json:1 | W1_B4 | NOT PORTED | PMID 33274414 abstract read: age, SpO2, albumin; range 0–55, cut-off 30; points per variable not stated |
| 12 | circ_covid | respiratory/circ_covid.json:1 | W1_B4 | NOT PORTED | PMID 32960645 abstract read: multistate model; no point values |
| 13 | copd_cat | respiratory/copd_cat.go:15 | W1_B4 | NOT PORTED | Jones 2009 PMID 19720809 abstract: 'eight items', '40-point scale'; per-item 0–5 scale and bands not seen (catestonline.org page silent; instrument licensed, GSK). Sum is trivial but item scale unconfirmed |
| 14 | covid_gram | respiratory/covid_gram.json:1 | W1_B4 | NOT PORTED | PMID 32396163 abstract read: 10 predictors with ORs; nomogram, no points |
| 15 | crb65_calc | respiratory/crb65_calc.json:1 | W1_B4 | NOT PORTED | Lim 2003 (PMC1746657) read: CRB-65 not defined in the text reached; Dwyer 2014 (PMC4212804) is a secondary validation; PMID 25478185 abstract adds points but no base definition |
| 16 | decaf_copd | respiratory/decaf_copd.json:1 | W1_B4 | NOT PORTED | PMID 22895999 abstract read: five predictors named, no points/thresholds (not OA in Europe PMC) |
| 17 | dissociation | respiratory/dissociation.json:1 | W1_B4 | NOT PORTED | SpO2↔PaO2 conversion cites Cleveland Clinic (secondary); Severinghaus 1979 (PMID 35496) equation not seen this pass |
| 18 | drip_pneumonia | respiratory/drip_pneumonia.json:1 | W1_B4 | NOT PORTED | PMID 26856838 abstract read: threshold ≥ 4 points stated; item points not stated |
| 19 | egri_difficult_airway | respiratory/egri_difficult_airway.json:1 | W1_B4 | NOT PORTED | PMID 8638791 abstract read: 0/1/2 weighting mentioned; per-criterion thresholds not stated |
| 20 | epiphany_index | respiratory/epiphany_index.json:1 | W1_B4 | NOT PORTED | PMID 28267709 abstract read: decision tree; six variables, no rules |
| 21 | feverpain | respiratory/feverpain.json:1 | W1_B4 | NOT PORTED | JSON cites Little 2012 (PMID 23211183, a cohort with ORs); FeverPAIN items in Little 2013 BMJ (PMID 24114306) abstract but point values not stated explicitly → not ported |
| 22 | fio2_estimate | respiratory/fio2_estimate.json:1 | W1_B4 | NOT PORTED | cites Cleveland Clinic/StatPearls (secondary); '21% + 4%/L' rule has no primary seen |
| 23 | gold_copd | respiratory/gold_copd.json:1 | W1_B4 | NOT PORTED | GOLD 2023 executive summary (Agustí, Eur Respir J 61:2300239, PMC10066569) read: FEV1/FVC < 0.7 seen; GOLD 1–4 FEV1 % bounds only in a figure not transcribed → not ported |
| 24 | heaven_criteria | respiratory/heaven_criteria.json:1 | W1_B4 | NOT PORTED | PMID 28739243 abstract read: six criteria named, no thresholds |
| 25 | lights_criteria | respiratory/lights_criteria.json:1 | W1_B4 | NOT PORTED | Light 1972 Ann Intern Med 77:507 (PMID 4642731): PubMed has no abstract; not OA → ratios 0.5/0.6/LDH not seen in the primary |
| 26 | lips_score | respiratory/lips_score.json:1 | W1_B4 | NOT PORTED | PMID 20562130 abstract read: items/points/thresholds not stated in the abstract; full text not open (Europe PMC: not OA) |
| 27 | mallampati_classification | respiratory/mallampati_classification.json:1 | W1_B4 | NOT PORTED | observed class (I–IV) entered directly; no computation |
| 28 | mmrc_dyspnea | respiratory/mmrc_dyspnea.go:15 | W1_B4 | NOT PORTED | source returns the picked grade unchanged (identity, no computation); also Mahler & Wells 1988 (PMID 3342669) is not the mMRC 0–4 wording's origin |
| 29 | newsom_score | respiratory/newsom_score.json:1 | W1_B4 | NOT PORTED | citation PMID 29356220 resolves to an unrelated paper (B. anthracis sampling); Newsom 2018 Acad Emerg Med primary not found this pass |
| 30 | ottawa_copd | respiratory/ottawa_copd.json:1 | W1_B4 | NOT PORTED | PMID 24549125 abstract read: five variables, no points |
| 31 | paradise_criteria | respiratory/paradise_criteria.json:1 | W1_B4 | NOT PORTED | Paradise 1984 NEJM (PMID 6700642) is a trial; abstract has no criteria list; eligibility criteria (decision rule, not a calculation) |
| 32 | pars | respiratory/pars.json:1 | W1_B4 | NOT PORTED | citation PMID 30389456 resolves to an unrelated paper (HCC metabolites); Biagini Myers 2019 JACI not reached |
| 33 | pass_asthma | respiratory/pass_asthma.json:1 | W1_B4 | NOT PORTED | PMID 14709423 abstract read: three items named, no points |
| 34 | pediatric_asthma_score | respiratory/pediatric_asthma_score.json:1 | W1_B4 | NOT PORTED | Kelly 2000 (PMID 10831004) is a clinical-pathway study; abstract has no score items |
| 35 | pediatric_ett_size | respiratory/pediatric_ett_size.json:1 | W1_B4 | NOT PORTED | citation PMID 10148666 resolves to an unrelated paper; age/4 + 4 rule (Cole) primary not seen; sizing is a device setting |
| 36 | pef_exacerbation | respiratory/pef_exacerbation.json:1 | W1_B4 | NOT PORTED | BTS/SIGN severity bands (PEF % best) not opened (site guideline PDF not reached) |
| 37 | pef_flowmetric | respiratory/pef_flowmetric.go:16 | W1_B4 | NOT PORTED | cites Mayo Clinic (secondary); 5.3·cm − 433 and (cm − 100)·5 + 100 have no primary seen (PMIDs 2399052/7734867 not checked for these exact forms) |
| 38 | pef_percent_predicted | respiratory/pef_percent_predicted.go:15 | W1_B4 | NOT PORTED | measured/predicted × 100 cites Hankinson 1999 (PMID 9872837, not OA), which does not define % predicted as such; no primary definition seen |
| 39 | pef_predicted | respiratory/pef_predicted.go:16 | W1_B4 | NOT PORTED | cites Hankinson 1999 (PMID 9872837) but its linear coefficients (−130.15 + 5.67 cm − 5.66 age; 213.41 + 2.80 cm − 4.91 age) are not Hankinson's (quadratic age, height², L/s, Table 4/5); paper not OA → not seen. Pediatric branch (cm − 100)·5 + 100 unsourced |
| 40 | pef_predicted_eu | respiratory/pef_predicted_eu.go:16 | W1_B4 | NOT PORTED | cites Nunn & Gregg 1989 BMJ 298:1068 (PMID 2497892, PMC1836460: scanned images only; PDF blocked); the linear form is Gregg & Nunn 1973; neither seen |
| 41 | pef_reversibility | respiratory/pef_reversibility.json:1 | W1_B4 | NOT PORTED | NICE NG80 / GINA reversibility thresholds not opened |
| 42 | pef_variability | respiratory/pef_variability.go:19 | W1_B4 | NOT PORTED | GINA report not opened (ginasthma.org); definition (max − min)/mean × 100 not seen in a primary |
| 43 | pef_zone | respiratory/pef_zone.go:22 | W1_B4 | NOT PORTED | cites Cleveland Clinic (secondary); NAEPP EPR-3 zone percentages not opened |
| 44 | pesi | respiratory/pesi.json:1 | W1_B4 | NOT PORTED | PMID 16020800 abstract read: '11 simple patient characteristics', no variables or points |
| 45 | piama_risk_score | respiratory/piama_risk_score.go:16 | W1_B4 | NOT PORTED | cited Brunekreef 2002 (PMID 12688626) is the cohort design paper, no score; the PIAMA risk score (Caudri 2009 JACI) not reached → points (2/1/4/1/4,7/2/6) not seen |
| 46 | pneumonia_risk_score | respiratory/pneumonia_risk_score.json:1 | W1_B4 | NOT PORTED | PMID 34694254 abstract read: variables named, no points |
| 47 | pompe_c | respiratory/pompe_c.json:1 | W1_B4 | NOT PORTED | PMID 22475313 abstract read: eight predictors; logistic model, no points |
| 48 | pram_asthma | respiratory/pram_asthma.json:1 | W1_B4 | NOT PORTED | PMID 11113831 abstract read: five variables named, no scoring |
| 49 | preset_ecmo | respiratory/preset_ecmo.json:1 | W1_B4 | NOT PORTED | PMID 29233160 abstract read: items/points/thresholds not stated in the abstract; full text not open (Europe PMC: not OA) |
| 50 | priest_covid | respiratory/priest_covid.json:1 | W1_B4 | NOT PORTED | PMID 33481930 abstract read: based on NEWS2, age, sex, performance status; allocation not stated |
| 51 | psi_calc | respiratory/psi_calc.json:1 | W1_B4 | NOT PORTED | PMID 8995086 abstract read: Fine 1997 NEJM point table not in abstract; not OA |
| 52 | psi_port_score | respiratory/psi_port_score.json:1 | W1_B4 | NOT PORTED | duplicate of psi_calc; same finding (Fine 1997 PMID 8995086) |
| 53 | qcsi_covid | respiratory/qcsi_covid.json:1 | W1_B4 | NOT PORTED | PMID 33012378 abstract read: items/points/thresholds not stated in the abstract; full text not open (Europe PMC: not OA) |
| 54 | rads_inhalation | respiratory/rads_inhalation.json:1 | W1_B4 | NOT PORTED | cited PMID 22210063 abstract: radiologist's CT score, methodology not stated |
| 55 | rdos_score | respiratory/rdos_score.json:1 | W1_B4 | NOT PORTED | PMID 20078243 abstract read: '8-item scale', items not stated |
| 56 | resect_90 | respiratory/resect_90.json:1 | W1_B4 | NOT PORTED | PMID 34324664 abstract read: items/points/thresholds not stated in the abstract; full text not open (Europe PMC: not OA) |
| 57 | resp_ecmo | respiratory/resp_ecmo.json:1 | W1_B4 | NOT PORTED | PMID 24693864 abstract read: items/points/thresholds not stated in the abstract; full text not open (Europe PMC: not OA) |
| 58 | respiratory_score_asthma | respiratory/respiratory_score_asthma.json:1 | W1_B4 | NOT PORTED | PMID 14966818 abstract read: four parameters named, no values |
| 59 | roth_score | respiratory/roth_score.json:1 | W1_B4 | NOT PORTED | PMID 27701750 abstract read: counting method, no bands |
| 60 | rr_classify | respiratory/rr_classify.go:15 | W1_B4 | NOT PORTED | identity (returns the rate); age bands from Cleveland Clinic (secondary). NEWS2 Chart 1 respiration-rate row was seen (RCP 2017) and could be ported as a separate NEWS2 function if wanted |
| 61 | rse_asthma | respiratory/rse_asthma.json:1 | W1_B4 | NOT PORTED | PMID 25258144 abstract read: '5 selected predictors', no points |
| 62 | scap_score | respiratory/scap_score.json:1 | W1_B4 | NOT PORTED | PMID 16973986 abstract read: eight factors, no points |
| 63 | silverman_calc | respiratory/silverman_calc.json:1 | W1_B4 | NOT PORTED | Silverman & Andersen 1956 (PMID 13353856): no abstract; item grades not seen in a primary |
| 64 | smart_cop | respiratory/smart_cop.json:1 | W1_B4 | NOT PORTED | Charles 2008 CID (PMID 18558884) abstract: points seen (SBP 2, multilobar 1, albumin 1, RR 1, HR 1, confusion 1, oxygenation 2, pH 2) but the thresholds ('low', 'high', age-adjusted) are not → not ported |
| 65 | stumbl_battle_score | respiratory/stumbl_battle_score.json:1 | W1_B4 | NOT PORTED | PMID 24887537 abstract read: predictors named, no points |
| 66 | utah_covid_risk | respiratory/utah_covid_risk.json:1 | W1_B4 | NOT PORTED | Webb 2022 PLoS One (PMID 35239664) abstract: points seen, but 'age (0.5 points per decade)' is ambiguous (floor? from 0?) → not computable without the full text (PLoS is OA; not opened this pass) |
| 67 | va_calc | respiratory/va_calc.json:1 | W1_B4 | NOT PORTED | held: doc/catalog/medicine.md §1 rows 44/48–50 (minute ventilation / dead space HELD-SOURCE); cites StatPearls (secondary) |
| 68 | vaco_index | respiratory/vaco_index.json:1 | W1_B4 | NOT PORTED | PMID 33175863 abstract read: items/points/thresholds not stated in the abstract; full text not open (Europe PMC: not OA) |
| 69 | vdvt_calc | respiratory/vdvt_calc.json:1 | W1_B4 | NOT PORTED | held: doc/catalog/medicine.md §1 rows 48–49 (Bohr/Enghoff dead space HELD-SOURCE) |
| 70 | ve_calc | respiratory/ve_calc.json:1 | W1_B4 | NOT PORTED | held: doc/catalog/medicine.md §1 row 44 (minute ventilation HELD-SOURCE) |
| 71 | westley_croup | respiratory/westley_croup.json:1 | W1_B4 | NOT PORTED | Westley 1978 Am J Dis Child (PMID 347921): no abstract; item scores not seen |
| 72 | winters_formula | respiratory/winters_formula.json:1 | W1_B4 | NOT PORTED | Albert, Dell & Winters 1967 Ann Intern Med 66:312, doi:10.7326/0003-4819-66-2-312 (PMID 6016545): no abstract; 1.5·HCO3 + 8 ± 2 not seen |
