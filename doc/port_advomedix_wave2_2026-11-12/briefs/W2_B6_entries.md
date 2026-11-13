# Entry list W2_B6 (wave 2)

Source root (read-only): `/Users/devinanderson/Desktop/advo_production/advomedix-v2/backendAdvo/`. Paths are as wave 1 recorded them: `knowledgebase/calculators/<category>/...`, `<category>/...` (= knowledgebase/calculators/<category>/), `knowledgebase/ranges/...`, `calculators/<vital>/...`, or relative to `knowledgebase/series/` for series functions. The wave 1 reason is the starting point, not the answer.

| # | entry id | source path (wave 1) | wave 1 section | wave 1 outcome | wave 1 reason |
|---|---|---|---|---|---|
| 1 | amt10 | knowledgebase/calculators/neuro/amt10.json:141 | W2_B5 | NOT PORTED | primary not seen: Hodkinson 1972 Age Ageing PMID 4669880 has no abstract scoring; item points and cut-off not seen |
| 2 | amt4 | knowledgebase/calculators/neuro/amt4.json:74 | W2_B5 | NOT PORTED | Swain 1997 PMID 9360037 abstract: items (age, DOB, place, year) and "score of less than four"; points per item not stated |
| 3 | aspects_stroke | knowledgebase/calculators/neuro/aspects_stroke.json:140 | W2_B5 | NOT PORTED | Barber 2000 Lancet PMID 10905241 abstract: "ten regions of interest"; the subtraction rule not stated in the abstract (full text not open). The 0–10 bound is used by MED-NEURO-011 |
| 4 | asrs_adhd | knowledgebase/calculators/neuro/asrs_adhd.json:121, neuro/asrs_adhd.go:13 | W2_B5 | NOT PORTED | Kessler 2005 Psychol Med PMID 15841682 abstract: 6-question screener confirmed but the per-item thresholds (Sometimes+/Often+) and the 4-of-6 cut-off not seen (full text not open); abstract scores all 18 items unweighted |
| 5 | bacterial_meningitis_score | knowledgebase/calculators/neuro/bacterial_meningitis_score.json:85 | W2_B5 | NOT PORTED | Nigrovic 2007 JAMA PMID 17200475 abstract: five criteria seen as a rule ("lack all of the following"); the score points (Gram stain 2) not seen |
| 6 | bars | knowledgebase/calculators/neuro/bars.json:41 | W2_B5 | NOT PORTED | primary not seen: Swift 2002 PMID 11777497 abstract has no levels |
| 7 | bjh_stroke_dysphagia | knowledgebase/calculators/neuro/bjh_stroke_dysphagia.json:84 | W2_B5 | NOT PORTED | REQUESTED by W1-B2. Edmiaston 2010 PMID 19875722 abstract: no screening items or failure rule stated |
| 8 | cam_icu | knowledgebase/calculators/neuro/cam_icu.json:86 | W2_B5 | NOT PORTED | decision rule (feature algorithm), no score; Inouye 1990 PMID 2240918 abstract gives the CAM algorithm, not CAM-ICU (Ely 2001 PMID 11445689) |
| 9 | canadian_ct_head | knowledgebase/calculators/neuro/canadian_ct_head.json:108 | W2_B5 | NOT PORTED | decision rule (any criterion present), no score; Stiell 2001 Lancet PMID 11356436 (criteria seen in abstract) |
| 10 | canadian_tia_score | knowledgebase/calculators/neuro/canadian_tia_score.json:173 | W2_B5 | NOT PORTED | primary not seen: Perry 2014 Stroke PMID 24262323 abstract names 13 variables, no points |
| 11 | capd | knowledgebase/calculators/neuro/capd.json:129 | W2_B5 | NOT PORTED | Traube 2014 Crit Care Med PMID 24145848 / PMC5527829: 8 items and cut-off 9 seen; response-option points and reverse items not in text (figure) |
| 12 | catch_head_injury | knowledgebase/calculators/neuro/catch_head_injury.json:108 | W2_B5 | NOT PORTED | decision rule (any criterion present), no score; Osmond 2010 CMAJ PMID 20142371 (criteria seen in abstract) |
| 13 | chalice_rule | knowledgebase/calculators/neuro/chalice_rule.json:192 | W2_B5 | NOT PORTED | decision rule (any criterion present), no score; Dunning 2006 Arch Dis Child; CITATION WRONG: the JSON PMID 16861481 resolves to Acerini 2006 (endocrine disruptors) |
| 14 | chip_prediction | knowledgebase/calculators/neuro/chip_prediction.json:194 | W2_B5 | NOT PORTED | primary not seen: Smits 2007 Ann Intern Med PMID 17371884 abstract lists no criteria/points |
| 15 | ciwa_ar | knowledgebase/calculators/neuro/ciwa_ar.json:140 | W2_B5 | NOT PORTED | primary not seen: Sullivan 1989 PMID 2597811 abstract: "10-item scale", item ranges and total not stated |
| 16 | cns_ipi | knowledgebase/calculators/neuro/cns_ipi.json:96 | W2_B5 | NOT PORTED | Schmitz 2016 JCO PMID 27382100 abstract: IPI factors + kidney/adrenal and three risk groups named; group cut-offs (factor counts) not stated |
| 17 | cpp_calc | knowledgebase/calculators/neuro/cpp_calc.json:58 | W2_B5 | NOT PORTED | primary not seen: Steiner 2006 BJA PMID 16698860 is a review (not open); CPP = MAP − ICP not seen in a primary; other citations StatPearls/Cleveland Clinic (secondary) |
| 18 | cssrs | knowledgebase/calculators/neuro/cssrs.json:96 | W2_B5 | NOT PORTED | Posner 2011 AJP PMID 22193671 / PMC3893686: no computed total in the abstract; licensed instrument (D-M8); severity levels not quoted |
| 19 | cts_6 | knowledgebase/calculators/neuro/cts_6.json:96 | W2_B5 | NOT PORTED | primary not seen: Graham 2006 J Hand Surg PMID 16886290 abstract has no item weights |
| 20 | dsm5_mdd | knowledgebase/calculators/neuro/dsm5_mdd.json:173 | W2_B5 | NOT PORTED | primary not seen: DSM-5 (APA 2013, doi:10.1176/appi.books.9780890425596) not open; criteria counts not read |
| 21 | dsm5_ptsd | knowledgebase/calculators/neuro/dsm5_ptsd.json:152 | W2_B5 | NOT PORTED | primary not seen: DSM-5 (APA 2013) not open; criteria counts not read |
| 22 | ed_safe_pss3 | knowledgebase/calculators/neuro/ed_safe_pss3.json:74 | W2_B5 | NOT PORTED | primary not seen: Boudreaux 2013 PMID 23707435 is the trial design paper; PSS-3 items/positive rule not in abstract |
| 23 | finnegan_nas | knowledgebase/calculators/neuro/finnegan_nas.json:262 | W2_B5 | NOT PORTED | primary not seen: cites Kocherlakota 2014 review (PMID 25070299); Finnegan 1975 item weights not seen |
| 24 | fisher_sah | knowledgebase/calculators/neuro/fisher_sah.json:41 | W2_B5 | NOT PORTED | primary not seen: Fisher 1980 Neurosurgery PMID 7354892 no abstract grades |
| 25 | flacc_scale | knowledgebase/calculators/neuro/flacc_scale.json:85 | W2_B5 | NOT PORTED | Merkel 1997 PMID 9220806 abstract: five categories named; 0–2 per item and 0–10 total not stated |
| 26 | four_score_coma | knowledgebase/calculators/neuro/four_score_coma.json:74 | W2_B5 | NOT PORTED | Wijdicks 2005 Ann Neurol PMID 16178024 abstract: four components "maximal score of 4"; component minimum (0) and grade definitions not stated |
| 27 | func_score | knowledgebase/calculators/neuro/func_score.json:84 | W2_B5 | NOT PORTED | Rost 2008 Stroke PMID 18556582 abstract: "sum of individual points (0-11)"; points not stated |
| 28 | gds15 | knowledgebase/calculators/neuro/gds15.json:222, neuro/gds15.go:13 | W2_B5 | NOT PORTED | CITATION WRONG: Yesavage 1982 PMID 7183759 is the 30-item GDS (abstract: "30 questions"); GDS-15 (Sheikh & Yesavage 1986, Clin Gerontol) not in PubMed, not seen |
| 29 | gmaws | knowledgebase/calculators/neuro/gmaws.json:85 | W2_B5 | NOT PORTED | CITATION WRONG: JSON PMID 22345174 resolves to Dell'Amore 2012 (aortic valve replacement); McPherson GMAWS primary not located this pass |
| 30 | ham_a | knowledgebase/calculators/neuro/ham_a.json:184 | W2_B5 | NOT PORTED | primary not seen: Hamilton 1959 PMID 13638508 no abstract; item 0–4 scoring and bands not seen (bands are not Hamilton's) |
| 31 | ham_d | knowledgebase/calculators/neuro/ham_d.json:217 | W2_B5 | NOT PORTED | Hamilton 1960 JNNP PMID 14399272 (PMC495331 page scan, not readable as text); bands not in the primary; licensed form (D-M8) |
| 32 | hat_score | knowledgebase/calculators/neuro/hat_score.json:65 | W2_B5 | NOT PORTED | Lou 2008 Neurology PMID 18955684, PMC2676961: Table 1 points not readable (PMC page behind a browser check; Europe PMC full text not OA); only sICH rates by score seen |
| 33 | heads_ed | knowledgebase/calculators/neuro/heads_ed.json:107 | W2_B5 | NOT PORTED | Cappelli 2012 Pediatrics PMID 22826567 abstract: seven domains named; 0–2 scoring not stated |
| 34 | hii_impairment | knowledgebase/calculators/neuro/hii_impairment.json:85 | W2_B5 | NOT PORTED | CITATION WRONG: JSON PMID 24405226 resolves to Ahamad 2014 (injecting drug use); Hack HII primary not located |
| 35 | hints_stroke | knowledgebase/calculators/neuro/hints_stroke.json:64 | W2_B5 | NOT PORTED | bedside examination algorithm, no score (Kattah 2009 PMID 19762709) |
| 36 | hunt_hess | knowledgebase/calculators/neuro/hunt_hess.json:42 | W2_B5 | NOT PORTED | primary not seen: Hunt & Hess 1968 PMID 5635959 no abstract grades |
| 37 | impact_score | knowledgebase/calculators/neuro/impact_score.json:148 | W2_B5 | NOT PORTED | Steyerberg 2008 PLoS Med PMID 18684008 (OA): LP formulas in Fig 1 caption (−2.55 + 0.275 × sum score core; −1.62 + 0.299 ×) seen, but the core sum-score point chart is an image only |
| 38 | infant_scalp_score | knowledgebase/calculators/neuro/infant_scalp_score.json:63 | W2_B5 | NOT PORTED | primary not seen: Greenes 2001 PMID 11334100 abstract has no points |
| 39 | lams | knowledgebase/calculators/neuro/lams.json:63 | W2_B5 | NOT PORTED | Llanes 2004 PMID 14691787 abstract: items (facial weakness, arm strength, grip) and 0–5 total; per-item points not stated |
| 40 | major_depression_index | knowledgebase/calculators/neuro/major_depression_index.json:193, neuro/major_depression_index.go:17 | W2_B5 | NOT PORTED | Bech 2001 J Affect Disord PMID 11578668 abstract: "simple total sum of the items", cut-off 26; the paired-item max (8a/8b, 10a/10b), 0–5 ranges and DSM symptom thresholds not seen (full text not open) |
| 41 | mbig | knowledgebase/calculators/neuro/mbig.json:86 | W2_B5 | NOT PORTED | decision pathway, not a calculation; brain injury guideline triage categories (Khan 2020 PMID 32537518) |
| 42 | mdq_bipolar | knowledgebase/calculators/neuro/mdq_bipolar.json:200, neuro/mdq_bipolar.go:17 | W2_B5 | NOT PORTED | Hirschfeld 2000 AJP PMID 11058490 abstract: "score of 7 or more items"; item count (13), co-occurrence and impairment conditions not stated |
| 43 | midas_migraine | knowledgebase/calculators/neuro/midas_migraine.json:105, neuro/midas_migraine.go:17 | W2_B5 | NOT PORTED | CITATION WRONG: Stewart 1999 Pain PMID 10068175 is the HImQ (8-item), not MIDAS. MIDAS primary Stewart 2001 Neurology PMID 11294956 abstract: five questions counting days in 3 months; summation and grades I–IV not stated |
| 44 | mnihss | knowledgebase/calculators/neuro/mnihss.json:151 | W2_B5 | NOT PORTED | Lyden 2001 Stroke PMID 11387492 abstract: deleted items named; retained items' points not stated |
| 45 | moca_score | knowledgebase/calculators/neuro/moca_score.json:111 | W2_B5 | NOT PORTED | Nasreddine 2005 JAGS PMID 15817019 abstract: cut-off 26 only; total and education point not stated; licensed (D-M8) |
| 46 | modified_fisher_sah | knowledgebase/calculators/neuro/modified_fisher_sah.json:53 | W2_B5 | NOT PORTED | Frontera 2006 PMID 16823296 abstract: grades based on cisternal blood and IVH; grade definitions not stated |
| 47 | modified_rankin_scale | knowledgebase/calculators/neuro/modified_rankin_scale.json:41 | W2_B5 | NOT PORTED | van Swieten 1988 Stroke PMID 3363593 abstract: "six grades (0-5)"; grade descriptions not stated |
| 48 | mrs_9q | knowledgebase/calculators/neuro/mrs_9q.json:168, neuro/mrs_9q.go:13 | W2_B5 | NOT PORTED | Patel 2012 Neurosurgery PMID 22843133 abstract: 9 yes/no questions; the question→grade mapping not stated (full text not open) |
| 49 | mtoq4 | knowledgebase/calculators/neuro/mtoq4.json:100, neuro/mtoq4.go:17 | W2_B5 | NOT PORTED | CITATION WRONG: Lipton 2009 Cephalalgia PMID 19239676 derives M-TOQ-15 and M-TOQ-5, not M-TOQ-4; mTOQ-4 scoring and bands not seen |
| 50 | new_orleans_head_ct | knowledgebase/calculators/neuro/new_orleans_head_ct.json:107 | W2_B5 | NOT PORTED | decision rule (any criterion present), no score; Haydel 2000 NEJM PMID 10891517 (seven findings seen in abstract) |
| 51 | nexus_head_ct | knowledgebase/calculators/neuro/nexus_head_ct.json:117 | W2_B5 | NOT PORTED | decision rule (any criterion present), no score; Mower 2005 PMID 16374287 (criteria seen in abstract) |
| 52 | nps_pain | knowledgebase/calculators/neuro/nps_pain.json:172, neuro/nps_pain.go:17 | W2_B5 | NOT PORTED | Galer & Jensen 1997 Neurology PMID 9040716 abstract: no items, scoring or bands; the source's average bands have no primary |
| 53 | ottawa_sah_rule | knowledgebase/calculators/neuro/ottawa_sah_rule.json:97 | W2_B5 | NOT PORTED | decision rule (any criterion present), no score; Perry 2013 JAMA PMID 24065011 (criteria seen in abstract) |
| 54 | painad | knowledgebase/calculators/neuro/painad.json:85 | W2_B5 | NOT PORTED | Warden 2003 PMID 12807591 abstract: five items, 0–10 total; per-item 0–2 not stated |
| 55 | palchak_pediatric_head_trauma | knowledgebase/calculators/neuro/palchak_pediatric_head_trauma.json:95 | W2_B5 | NOT PORTED | decision rule (any criterion present), no score; Palchak 2003 PMID 14520320 (criteria seen in abstract) |
| 56 | pecarn | knowledgebase/calculators/neuro/pecarn.json:151 | W2_B5 | NOT PORTED | decision rule (any criterion present), no score; Kuppermann 2009 Lancet PMID 19758692 (criteria seen in abstract) |
| 57 | pecarn_cspine | knowledgebase/calculators/neuro/pecarn_cspine.json:119 | W2_B5 | NOT PORTED | decision rule (any criterion present), no score; Leonard 2024 PMID 38843852 |
| 58 | ped_nihss | knowledgebase/calculators/neuro/ped_nihss.json:195 | W2_B5 | NOT PORTED | Ichord 2011 Stroke PMID 21317270 (PMC3065389): item adaptations and points not read this pass |
| 59 | pediatric_gcs | knowledgebase/calculators/neuro/pediatric_gcs.json:64 | W2_B5 | NOT PORTED | primary not seen: James 1986 PMID 3951884 no abstract scoring |
| 60 | pediatric_nexus_ii | knowledgebase/calculators/neuro/pediatric_nexus_ii.json:108 | W2_B5 | NOT PORTED | decision rule (any criterion present), no score; Oman 2006 Pediatrics PMID 16418311 |
| 61 | qids_depression | knowledgebase/calculators/neuro/qids_depression.json:253, neuro/qids_depression.go:19 | W2_B5 | NOT PORTED | Rush 2003 Biol Psychiatry PMID 12946886 abstract: "16-item"; domain-max scoring and severity bands not stated (full text not open) |
| 62 | rass | knowledgebase/calculators/neuro/rass.json:41 | W2_B5 | NOT PORTED | Sessler 2002 AJRCCM PMID 12421743 abstract: 10 levels, only "+4 combative" and "-5 unarousable" labels stated |
| 63 | rcvs2_score | knowledgebase/calculators/neuro/rcvs2_score.json:85 | W2_B5 | NOT PORTED | REQUESTED by W1-B2. Rocha 2019 Neurology PMID 30635475 abstract: components named, range −2 to +10, cut-offs ≥5 / ≤2 / 3–4; item points not stated |
| 64 | rome_iv_child_abdominal_migraine | neuro/rome_iv_child_abdominal_migraine.go:18 | W1_B2 | NOT PORTED | CORRECTED: entry is titled Rome IV but cites Rome III child/adolescent (Rasquin 2006, PMC7104693); not read (PMC captcha / Europe PMC full text 500); Rome IV child (Hyams 2016, PMID 27144632) not open; checklist of caller-judged criteria, criteria text © Rome Foundation (D-M8) |
| 65 | rosier_scale | knowledgebase/calculators/neuro/rosier_scale.json:106 | W2_B5 | NOT PORTED | Nor 2005 Lancet Neurol PMID 16239179 abstract: items named, range −2 to +5, cut-off >0; per-item points not stated |
| 66 | rule_of_7s_lyme | knowledgebase/calculators/neuro/rule_of_7s_lyme.json:62 | W2_B5 | NOT PORTED | Garro 2009 Pediatrics PMID 19403476 abstract: risk categories only; the three criteria not stated |
| 67 | setscore | knowledgebase/calculators/neuro/setscore.json:196 | W2_B5 | NOT PORTED | Schönenberger 2016 PMID 26842719 abstract: cut-off 8 only; items not stated |
| 68 | sudbury_vertigo | knowledgebase/calculators/neuro/sudbury_vertigo.json:107 | W2_B5 | NOT PORTED | Ohle 2025 Ann Emerg Med PMID 39093245 abstract: bands <5, 5–8, >8 with risks; item points not stated |
| 69 | tdn_grade | knowledgebase/calculators/neuro/tdn_grade.json:64 | W2_B5 | NOT PORTED | primary not seen: Terrapon 2021 Neurosurgery PMID 33887774 abstract has no grade definitions |

## Routed from W1-B4 (BUILD IN med/psychiatry; W1-B4's research is in ledger/W1_B4.md — read its row first; your row supersedes it)

| # | entry id | source path (wave 1) | wave 1 section | wave 1 outcome | wave 1 reason |
|---|---|---|---|---|---|
| R1 | audit_c | general/audit_c | W1_B1 | NOT PORTED | see ledger/W1_B4.md |
| R2 | brief_addiction_monitor | general/brief_addiction_monitor | W1_B1 | NOT PORTED | open at PMC3602977, not read; see ledger/W1_B4.md |
| R3 | cage_alcohol | general/cage_alcohol | W1_B1 | NOT PORTED | see ledger/W1_B4.md |
| R4 | cas_covid_anxiety | general/cas_covid_anxiety | W1_B1 | NOT PORTED | see ledger/W1_B4.md |
| R5 | ciwa_b | general/ciwa_b | W1_B1 | NOT PORTED | see ledger/W1_B4.md |
| R6 | cows_score | general/cows_score | W1_B1 | NOT PORTED | see ledger/W1_B4.md |
| R7 | dast_10 | general/dast_10 | W1_B1 | NOT PORTED | cited source is the 28-item DAST; see ledger/W1_B4.md |
