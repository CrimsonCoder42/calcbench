# Entry list W2_B5 (wave 2)

Source root (read-only): `/Users/devinanderson/Desktop/advo_production/advomedix-v2/backendAdvo/`. Paths are as wave 1 recorded them: `knowledgebase/calculators/<category>/...`, `<category>/...` (= knowledgebase/calculators/<category>/), `knowledgebase/ranges/...`, `calculators/<vital>/...`, or relative to `knowledgebase/series/` for series functions. The wave 1 reason is the starting point, not the answer.

| # | entry id | source path (wave 1) | wave 1 section | wave 1 outcome | wave 1 reason |
|---|---|---|---|---|---|
| 1 | ndi | knowledgebase/calculators/general/ndi.json:18, general/ndi.go:15 | W1_B1 | NOT PORTED | Vernon & Mior 1991 PMID 1834753 confirmed: abstract says '10-item scaled questionnaire'; 0–5 item scoring and the 50-point total not seen in the primary |
| 2 | nedocs | knowledgebase/calculators/general/nedocs.json:14 | W1_B1 | NOT PORTED | primary not seen: Weiss 2004 PMID 14709427: regression equation not in abstract |
| 3 | neonatal_hyperbilirubinemia | knowledgebase/calculators/general/neonatal_hyperbilirubinemia.json:16 | W1_B1 | NOT PORTED | primary not seen: Bhutani 1999 PMID 9917432: hour-specific nomogram percentile curves not in abstract (graph only) |
| 4 | nhfs | general/nhfs.json:106 | W1_B2 | NOT PORTED | Maxwell 2008 (PMID 18723517) abstract states no points; not open |
| 5 | nihss | knowledgebase/calculators/general/nihss.json:16 | W1_B1 | NOT PORTED | primary not seen: Brott 1989 PMID 2749846: item scoring not in abstract; NIHSS item text is the NINDS form, not read |
| 6 | nutritional_risk_index | general/nutritional_risk_index.json:60 | W1_B2 | NOT PORTED | Buzby 1988 (PMID 3124597) and VA 1991 (PMID 1906987) abstracts state no formula; neither open |
| 7 | nyha | knowledgebase/calculators/general/nyha.json:15 | W1_B1 | NOT PORTED | citation mismatch: Goldman 1981 PMID 7296795 is the Specific Activity Scale; the NYHA classes are from the NYHA Criteria Committee (1964/1994 book), not seen; classes are a lookup, not a calculation |
| 8 | oakland_score | general/oakland_score.json:107 | W1_B2 | NOT PORTED | Oakland 2017 (PMID 28651935) abstract gives only "8 or less"; points table not open |
| 9 | ost_osteoporosis | general/ost_osteoporosis.go:22 | W1_B2 | NOT PORTED | OST formula not seen: Koh 2001 (PMID 11580084) abstract states none and is not open; Adler 2003 (PMID 12934782) abstract gives only the male bands (≥ 4, −1 to 3, ≤ −2) and the cut-off 3 |
| 10 | ovat_ipv | knowledgebase/calculators/general/ovat_ipv.json:16 | W1_B1 | NOT PORTED | primary not seen: Weiss 2003 PMID 12816400: item scoring not in abstract |
| 11 | parc_appendicitis | general/parc_appendicitis.json:137 | W1_B2 | NOT PORTED | Kharbanda 2018 (PMID 29535251) model coefficients not in abstract; not open |
| 12 | pas_ii_ra | knowledgebase/calculators/general/pas_ii_ra.json:15 | W1_B1 | NOT PORTED | primary not seen: Wolfe 2005 PMID 16331773: PAS-II formula not in abstract, full text not open |
| 13 | pas_ra | knowledgebase/calculators/general/pas_ra.json:16 | W1_B1 | NOT PORTED | primary not seen: Wolfe 2005 PMID 16331773: PAS formula not in abstract, full text not open |
| 14 | pawss | knowledgebase/calculators/general/pawss.json:17 | W1_B1 | NOT PORTED | Maldonado 2014 PMID 24657098 confirmed: '10 items', 'threshold score of 4'; item points not in abstract |
| 15 | pcdai | general/pcdai.json:169 | W1_B2 | NOT PORTED | Hyams 1991 (PMID 1678008) abstract states no item table; not open |
| 16 | pedocs | knowledgebase/calculators/general/pedocs.json:15 | W1_B1 | NOT PORTED | primary not seen: Weiss 2007 PMID 17876254: regression equation not in abstract |
| 17 | pews | knowledgebase/calculators/general/pews.json:15 | W1_B1 | NOT PORTED | primary not seen: Monaghan 2005 PMID 15960531 (Paediatr Nurs): score chart not seen |
| 18 | plasmic_ttp | knowledgebase/calculators/general/plasmic_ttp.json:16 | W1_B1 | NOT PORTED | primary not seen: Bendapudi 2017 PMID 28259520: 'seven-component' tool; components not in abstract, full text not open |
| 19 | pucai | general/pucai.go:28 | W1_B2 | NOT PORTED | Turner 2007 (PMID 17681163) abstract states no item points; not open |
| 20 | rapid_pleural | knowledgebase/calculators/general/rapid_pleural.json:15 | W1_B1 | NOT PORTED | Rahman 2014 PMID 24264558 doi:10.1378/chest.13-1558 confirmed: variables and groups 0-2/3-4/5-7 in the abstract; per-variable points not in the abstract |
| 21 | revised_trauma_score | knowledgebase/calculators/general/revised_trauma_score.json:10 | W1_B1 | NOT PORTED | citation mismatch: Champion 1981 PMID 7273818 is the original Trauma Score; the Revised Trauma Score is Champion 1989 (J Trauma 29:623); coded values/weights not seen |
| 22 | rome4_anorectal_pain | general/rome4_anorectal_pain.json:86 | W1_B2 | NOT PORTED | Rao 2016 Rome IV anorectal (PMC5035713, author manuscript) not read (PMC captcha); checklist of caller-judged criteria with no published number; criteria text © Rome Foundation (D-M8) |
| 23 | rome_ii_ibs | general/rome_ii_ibs.go:24 | W1_B2 | NOT PORTED | Thompson 1999 Gut Suppl (PMC1766683) not read (PMC holds a PDF only, download not fetchable); rule is 12 weeks + 2 of 3 caller-judged features: retry |
| 24 | rome_iii_ibs | general/rome_iii_ibs.go:15 | W1_B2 | NOT PORTED | Rome III (Drossman 2006, PMID 16678553) not open |
| 25 | rome_iv_biliary | general/rome_iv_biliary.json:129 | W1_B2 | NOT PORTED | Rome IV primary (Gastroenterology 2016;150(6), Rome Foundation ©) not open (Europe PMC: no full text); checklist of caller-judged criteria with no published number; criteria text licensed (D-M8) |
| 26 | rome_iv_caps | general/rome_iv_caps.go:15 | W1_B2 | NOT PORTED | Rome IV primary (Gastroenterology 2016;150(6), Rome Foundation ©) not open (Europe PMC: no full text); checklist of caller-judged criteria with no published number; criteria text licensed (D-M8) |
| 27 | rome_iv_child_aerophagia | general/rome_iv_child_aerophagia.go:15 | W1_B2 | NOT PORTED | CORRECTED: entry is titled Rome IV but cites Rome III child/adolescent (Rasquin 2006, PMC7104693); not read (PMC captcha / Europe PMC full text 500); Rome IV child (Hyams 2016, PMID 27144632) not open; checklist of caller-judged criteria, criteria text © Rome Foundation (D-M8) |
| 28 | rome_iv_child_constipation | general/rome_iv_child_constipation.go:15 | W1_B2 | NOT PORTED | CORRECTED: entry is titled Rome IV but cites Rome III child/adolescent (Rasquin 2006, PMC7104693); not read (PMC captcha / Europe PMC full text 500); Rome IV child (Hyams 2016, PMID 27144632) not open; checklist of caller-judged criteria, criteria text © Rome Foundation (D-M8) |
| 29 | rome_iv_child_fap | general/rome_iv_child_fap.go:15 | W1_B2 | NOT PORTED | Rome IV primary (Gastroenterology 2016;150(6), Rome Foundation ©) not open (Europe PMC: no full text); checklist of caller-judged criteria with no published number; criteria text licensed (D-M8) |
| 30 | rome_iv_child_ibs | general/rome_iv_child_ibs.go:21 | W1_B2 | NOT PORTED | Rome IV primary (Gastroenterology 2016;150(6), Rome Foundation ©) not open (Europe PMC: no full text); checklist of caller-judged criteria with no published number; criteria text licensed (D-M8) |
| 31 | rome_iv_child_nausea_vomiting | general/rome_iv_child_nausea_vomiting.go:15 | W1_B2 | NOT PORTED | functional nausea/vomiting are Rome IV-only, but the entry cites Rome III (Rasquin 2006); Rome IV child (Hyams 2016) not open; checklist of caller-judged criteria (D-M8) |
| 32 | rome_iv_child_rumination | general/rome_iv_child_rumination.go:15 | W1_B2 | NOT PORTED | CORRECTED: entry is titled Rome IV but cites Rome III child/adolescent (Rasquin 2006, PMC7104693); not read (PMC captcha / Europe PMC full text 500); Rome IV child (Hyams 2016, PMID 27144632) not open; checklist of caller-judged criteria, criteria text © Rome Foundation (D-M8) |
| 33 | rome_iv_chs | general/rome_iv_chs.go:15 | W1_B2 | NOT PORTED | Rome IV primary (Gastroenterology 2016;150(6), Rome Foundation ©) not open (Europe PMC: no full text); checklist of caller-judged criteria with no published number; criteria text licensed (D-M8) |
| 34 | rome_iv_diarrhea | general/rome_iv_diarrhea.go:15 | W1_B2 | NOT PORTED | Rome IV primary (Gastroenterology 2016;150(6), Rome Foundation ©) not open (Europe PMC: no full text); checklist of caller-judged criteria with no published number; criteria text licensed (D-M8) |
| 35 | rome_iv_fecal_incontinence | general/rome_iv_fecal_incontinence.go:15 | W1_B2 | NOT PORTED | Rao 2016 Rome IV anorectal (PMC5035713, author manuscript) not read (PMC captcha); checklist of caller-judged criteria with no published number; criteria text © Rome Foundation (D-M8) |
| 36 | rome_iv_functional_constipation | general/rome_iv_functional_constipation.go:15 | W1_B2 | NOT PORTED | Rome IV primary (Gastroenterology 2016;150(6), Rome Foundation ©) not open (Europe PMC: no full text); checklist of caller-judged criteria with no published number; criteria text licensed (D-M8) |
| 37 | rome_iv_heartburn | general/rome_iv_heartburn.go:15 | W1_B2 | NOT PORTED | Rome IV primary (Gastroenterology 2016;150(6), Rome Foundation ©) not open (Europe PMC: no full text); checklist of caller-judged criteria with no published number; criteria text licensed (D-M8) |
| 38 | rome_iv_infant_colic | general/rome_iv_infant_colic.go:15 | W1_B2 | NOT PORTED | Rome IV primary (Gastroenterology 2016;150(6), Rome Foundation ©) not open (Europe PMC: no full text); checklist of caller-judged criteria with no published number; criteria text licensed (D-M8) |
| 39 | saps3_score | knowledgebase/calculators/general/saps3_score.json:14 | W1_B1 | NOT PORTED | Moreno 2005 PMID 16132892 (PMC1315315) confirmed; ~20-variable point table not transcribed in this pass |
| 40 | saps_ii | knowledgebase/calculators/general/saps_ii.json:15 | W1_B1 | NOT PORTED | primary not seen: Le Gall 1993 PMID 8254858 JAMA: point table in full text, not open |
| 41 | scorma_index | knowledgebase/calculators/general/scorma_index.json:15 | W1_B1 | NOT PORTED | primary not seen: Heide 2001 PMID 11720175: formula not in abstract |
| 42 | smast | knowledgebase/calculators/general/smast.json:17, general/smast.go:18 | W1_B1 | NOT PORTED | citation mismatch: Selzer 1971 PMID 5565851 is the 25-item MAST; the 13-item SMAST is Selzer 1975 PMID 238068; item keying not seen (no abstract scoring, full text not open) |
| 43 | smast_g | knowledgebase/calculators/general/smast_g.json:19, general/smast_g.go:15 | W1_B1 | NOT PORTED | citation mismatch: Selzer 1975 PMID 238068 is the SMAST; SMAST-G is Blow 1998 (Michigan Alcoholism Screening Test—Geriatric Version, short form), not cited; keying not seen |
| 44 | temp_abg_correction | knowledgebase/calculators/general/temp_abg_correction.json:10 | W1_B1 | NOT PORTED | no primary cited (LITFL is a secondary); formula not seen in a primary |
| 45 | temp_classifier | knowledgebase/calculators/general/temp_classifier.json:15, general/temp_classifier.go:6 | W1_B1 | NOT PORTED | no primary cited (Cleveland Clinic page); age-specific 'normal' ranges not traceable to a primary |
| 46 | temp_clinical_scores | knowledgebase/calculators/general/temp_clinical_scores.json:10 | W1_B1 | NOT PORTED | no primary cited (Cleveland Clinic page) |
| 47 | temp_liebermeister | knowledgebase/calculators/general/temp_liebermeister.json:16, general/temp_liebermeister.go:15 | W1_B1 | NOT PORTED | Cunha 2000 PMID 11284920 doi:10.1046/j.1469-0691.2000.0194f.x confirmed, no abstract and no open full text; the source uses 10 bpm per °C (temp_liebermeister.go:19), which matches neither Liebermeister's rule as usually quoted (~8 bpm/°C) nor Cunha's 10 bpm per °F; Liebermeister's 1875 Handbuch not reached. Discrepancy for Devin |
| 48 | temp_mean_body | knowledgebase/calculators/general/temp_mean_body.json:15, general/temp_mean_body.go:15 | W1_B1 | NOT PORTED | Burton 1935 J Nutr 9:261-280 doi:10.1093/jn/9.3.261 confirmed via Crossref; no abstract, publisher page 404; the 0.64/0.36 weights not seen |
| 49 | temp_metabolic_cost | knowledgebase/calculators/general/temp_metabolic_cost.json:15, general/temp_metabolic_cost.go:18 | W1_B1 | NOT PORTED | no primary cited (StatPearls NBK562334 is a secondary); 10%/°C and the 2000 kcal baseline not traceable |
| 50 | temp_q10 | knowledgebase/calculators/general/temp_q10.json:17, general/temp_q10.go:19 | W1_B1 | NOT PORTED | no primary cited (PhysiologyWeb); Q10 = 2.5 for human metabolism not traceable to a primary |
| 51 | thakar_score | knowledgebase/calculators/general/thakar_score.json:15 | W1_B1 | NOT PORTED | Thakar 2005 PMID 15563569 doi:10.1681/ASN.2004040331 confirmed: 'ranges between 0 and 17 points'; per-variable points not in abstract |
| 52 | wast | knowledgebase/calculators/general/wast.json:15 | W1_B1 | NOT PORTED | primary not seen: Brown 1996 PMID 8791071: scoring not in abstract |
| 53 | wat1 | knowledgebase/calculators/general/wat1.json:16 | W1_B1 | NOT PORTED | Franck 2008 PMID 18838937 (PMC2775493, doi:10.1097/PCC.0b013e31818c8328) confirmed: 11 items, total 0–12, cut-off ≥3; per-item 0–1/0–2 assignment read only through a summary, not verbatim |

## Routed from W1-B4 (BUILD IN; W1-B4's research is in ledger/W1_B4.md — read its row first; your row supersedes it)

| # | entry id | source path (wave 1) | home | wave 1 reason |
|---|---|---|---|---|
| R1 | air_score | general/air_score | med/gastro | see ledger/W1_B4.md |
| R2 | alvarado_appendicitis | general/alvarado_appendicitis | med/gastro | see ledger/W1_B4.md |
| R3 | glasgow_blatchford | general/glasgow_blatchford | med/gastro | see ledger/W1_B4.md |
| R4 | kruis_ibs | general/kruis_ibs | med/gastro | see ledger/W1_B4.md |
| R5 | four_peps | general/four_peps | med/cardio | open at PMC7931139 |
| R6 | hospital_score | general/hospital_score | med/cardio | beside LACEIndex |
| R7 | myxedema_coma | general/myxedema_coma | med/endocrine | see ledger/W1_B4.md |
