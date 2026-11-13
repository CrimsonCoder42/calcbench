# Entry list W1_B1 (wave 2)

Source root (read-only): `/Users/devinanderson/Desktop/advo_production/advomedix-v2/backendAdvo/`. Paths are as wave 1 recorded them: `knowledgebase/calculators/<category>/...`, `<category>/...` (= knowledgebase/calculators/<category>/), `knowledgebase/ranges/...`, `calculators/<vital>/...`, or relative to `knowledgebase/series/` for series functions. The wave 1 reason is the starting point, not the answer.

| # | entry id | source path (wave 1) | wave 1 section | wave 1 outcome | wave 1 reason |
|---|---|---|---|---|---|
| 1 | aap_peds_htn | cardiac/aap_peds_htn.go:25-87 | W1_B3 | NOT PORTED | AAP 2017 (Flynn JT, Pediatrics 2017;140(3):e20171904, doi:10.1542/peds.2017-1904, PMID 28827377) confirmed; the source uses a home-made height adjustment ((height − median)/2.5 × 2 or × 1 mmHg) and 12-row age tables that are not the guideline's height-percentile Tables 4/5; the AAP tables were not obtained verbatim |
| 2 | abi_calc | cardiac/abi_calc.json:50 | W1_B3 | NOT PORTED | entry cites Johns Hopkins (secondary); Aboyans 2012 AHA statement (PMID 23159553) not readable (403, no abstract) |
| 3 | adapt_protocol | cardiac/adapt_protocol.json:107 | W1_B3 | NOT PORTED | decision pathway (TIMI 0 + 0/2 h troponin protocol, Than 2012, PMID 22578923), not a calculation |
| 4 | altitude_adjusted_perc | cardiac/altitude_adjusted_perc.json:131 | W1_B3 | NOT PORTED | Wolf 2008 (PMID 18272098): altitude-adjusted SaO₂ cut-off not in abstract; PERC is a rule (all criteria) |
| 5 | arc_hbr | cardiac/arc_hbr.json:250 | W1_B3 | NOT PORTED | decision rule (≥1 major or ≥2 minor criteria), Urban 2019 (PMC6636810); criteria table not quoted cell by cell |
| 6 | ascvd_risk | cardiac/ascvd_risk.json:126 | W1_B3 | NOT PORTED | Goff 2014 (PMID 24222018): not in PMC; coefficients not read |
| 7 | boston_syncope | cardiac/boston_syncope.json:117 | W1_B3 | NOT PORTED | decision rule (any criterion present), no score (Grossman 2007, PMID 17976548) |
| 8 | bova_score | cardiac/bova_score.json:74 | W1_B3 | NOT PORTED | Bova 2014 ERJ (PMID 24696111): points and stage cut-offs not in abstract |
| 9 | bp_classify | cardiac/bp_classify.go (outputs sys/dia unchanged; bands in the app) | W1_B3 | NOT PORTED | held: doc/catalog/medicine.md §5 (BP categories, 2025 AHA/ACC table not read); entry cites Mayo Clinic (secondary) |
| 10 | brugada_criteria | cardiac/brugada_criteria.json:75 | W1_B3 | NOT PORTED | decision pathway: Brugada 1991 (PMID 2022022) is a 4-step algorithm, no score; steps 3–4 not seen |
| 11 | cahp_score | cardiac/cahp_score.json:103 | W1_B3 | NOT PORTED | Maupain 2016 EHJ (PMID 26497161): formula not in abstract; printed groups overlap at 150 and 200 |
| 12 | canadian_syncope | cardiac/canadian_syncope.json:118 | W1_B3 | NOT PORTED | Thiruganasambandamoorthy 2016 CMAJ (PMC5008955) read through a summariser: the transcribed Table 4 points sum to −3…9, but the paper states "−3 to 11", so at least two item weights (likely SBP and QTc) were mis-transcribed; category bands are in an image. Needs the PDF read by a person |
| 13 | cardiac_index | cardiac/cardiac_index.json:59 | W1_B3 | NOT PORTED | held: doc/catalog/medicine.md §5 (CI); also needs BSA (Du Bois/Mosteller HELD-SOURCE) |
| 14 | cardiac_output | cardiac/cardiac_output.json:81 | W1_B3 | NOT PORTED | held: doc/catalog/medicine.md §5 (CO; Guyton & Hall not opened); Fick with assumed VO₂ cites Cleveland Clinic + LaFarge 1970 (VO₂ estimate, not read) |
| 15 | cart_score | cardiac/cart_score.json:70 | W1_B3 | NOT PORTED | CORRECTED: Crit Care Med 2012;40:2102 (PMC3378796), not Chest; points table not reached |
| 16 | ccs_angina | cardiac/ccs_angina.go (grade = index + 1) | W1_B3 | NOT PORTED | classification chosen by the clinician (Campeau 1976, Circulation 54:522, PMID 947585); nothing to compute |
| 17 | chads65 | cardiac/chads65.json:97 | W1_B3 | NOT PORTED | decision pathway (CCS 2020 algorithm, PMID 33191198); text not seen |
| 18 | chadsvasc | cardiac/chadsvasc.go | W1_B3 | NOT PORTED | Lip 2010 Chest (PMID 19762550): abstract names the schema only; points table not seen (Chest 403); Friberg 2012 closed |
| 19 | cpes_score | cardiac/cpes_score.json:96 | W1_B3 | NOT PORTED | CORRECTED: PMID 37100484 is an unrelated JACC editorial; nearest is Zhang 2024 Circ Cardiovasc Interv (PMID 38994599, closed, no components in abstract) |
| 20 | cri | cardiac/cri.go (220 − age reserve) | W1_B3 | NOT PORTED | Wilkoff & Miller 1992 (PMID 1423382, citation corrected) not seen; formula seen only in a review (Brubaker 2011, PMC3065291); 220 − age not seen in a primary |
| 21 | crusade_score | cardiac/crusade_score.json:114 | W1_B3 | NOT PORTED | Subherwal 2009 (PMC3767035): Table 4 band points not reached; risk bands seen |
| 22 | dash_prediction_vte | cardiac/dash_prediction_vte.json:74 | W1_B3 | NOT PORTED | Tosetto 2012 (PMID 22489957): items seen, points not |
| 23 | diastolic_strain | cardiac/diastolic_strain.json:72 | W1_B3 | NOT PORTED | no published "diastolic strain index" formula: ASE/EACVI 2016 (Nagueh, PMID 27037982) is a decision algorithm over E/e′, LAVI, TR velocity; the entry's combined index has no primary |
| 24 | do2 | cardiac/do2.json:69 | W1_B3 | NOT PORTED | CaO₂ term (1.34 vs 1.39 mL/g) is HELD in doc/catalog/medicine.md §5 (respiratory row); Shoemaker 1988 (PMID 3191758) is a trial, not the formula's source |
| 25 | do2_simple | cardiac/do2_simple.json:49 | W1_B3 | NOT PORTED | CO × CaO₂ × 10: depends on CO (HELD, medicine.md §5); cites StatPearls (secondary) |
| 26 | doac_score | cardiac/doac_score.json:140 | W1_B3 | NOT PORTED | Aggarwal 2023 Circulation (PMC10529708) read only through a summariser; category ranges are in a supplement and the 30–60 boundary is unconfirmed |
| 27 | duke_iscvid | cardiac/duke_iscvid.json:153 | W1_B3 | NOT PORTED | decision pathway (diagnostic criteria, Fowler 2023, PMID 37138445), not a calculation |
| 28 | edacs | cardiac/edacs.json:128 | W1_B3 | NOT PORTED | Than 2014 EMA (PMID 24428678): points and the < 16 cut-off not in abstract; no open text |
| 29 | egsys_syncope | cardiac/egsys_syncope.json:96 | W1_B3 | NOT PORTED | Del Rosso 2008 Heart (PMID 18519550): points not in abstract |
| 30 | ehmrg | cardiac/ehmrg.json:135 | W1_B3 | NOT PORTED | Lee 2012 Ann Intern Med (PMID 22665814): equation not in abstract |
| 31 | euromacs_rhf | cardiac/euromacs_rhf.json:86 | W1_B3 | NOT PORTED | Soliman 2018 (PMID 28847897): points not in abstract |
| 32 | four_ts_hit | cardiac/four_ts_hit.json:75 | W1_B3 | NOT PORTED | Lo 2006 (PMID 16634744): items and points not in abstract |
| 33 | framingham_hard_chd | cardiac/framingham_hard_chd.json:104 | W1_B3 | NOT PORTED | Wilson 1998 (PMID 9603539): not in PMC; point tables not read |
| 34 | framingham_hf | cardiac/framingham_hf.json:184 | W1_B3 | NOT PORTED | decision pathway (2 major or 1 major + 2 minor criteria, McKee 1971, PMID 5122894), no score |
| 35 | garfield_af | cardiac/garfield_af.json:213 | W1_B3 | NOT PORTED | GARFIELD-AF (Fox 2017 BMJ Open, PMID 29273652) is distributed as a web/software tool; no printed equation was read in this run |
| 36 | go_far_score | cardiac/go_far_score.json:172 | W1_B3 | NOT PORTED | Ebell 2013 JAMA Intern Med (PMID 24018585): 13 variables, points not in abstract |
| 37 | grace_acs | cardiac/grace_acs.json:115 | W1_B3 | NOT PORTED | Fox 2006 BMJ (PMC1661748) nomogram points only as an image; Granger 2003 (PMID 14581255) points paywalled. Note: the app's variable list (history of MI, in-hospital PCI) matches Eagle 2004 JAMA (PMID 15187054), not Fox 2006 |
| 38 | grogan_attr_cm | cardiac/grogan_attr_cm.json:50 | W1_B3 | NOT PORTED | Grogan 2016 JACC (PMID 27585505): abstract does not say whether a value equal to 0.05 ng/mL or 3000 pg/mL is above or below |
| 39 | has_bled | cardiac/has_bled.json:130 | W1_B3 | NOT PORTED | Pisters 2010 Chest (PMID 20299623): items in abstract, points/definitions not seen |
| 40 | hcm_risk_scd | cardiac/hcm_risk_scd.json:103 | W1_B3 | NOT PORTED | O'Mahony 2014 EHJ (PMID 24126876): coefficients and 0.998 baseline not in abstract; likely in the 2014 ESC HCM guideline (not read) |
| 41 | he_macs | cardiac/he_macs.json:117 | W1_B3 | NOT PORTED | Alghamdi 2019 (PMC6728057): coefficients in an image table, cut-offs in the supplement |
| 42 | hear_score | cardiac/hear_score.json:74 | W1_B3 | NOT PORTED | Smith 2020 EMJ (PMC7952041) seen: item definitions deferred to the HEART score papers (Six 2008, Backus 2010), not read |
| 43 | heart_pathway | cardiac/heart_pathway.json:84 | W1_B3 | NOT PORTED | decision pathway (HEART score + serial troponin, Mahler 2015, PMID 25737484), not a calculation |
| 44 | hemorr2hages | cardiac/hemorr2hages.json:151 | W1_B3 | NOT PORTED | Gage 2006 AHJ (PMID 16504638): item list not in abstract |
| 45 | hep_score_hit | cardiac/hep_score_hit.json:217 | W1_B3 | NOT PORTED | Cuker 2010 (PMID 20854372): items not in abstract |
| 46 | hestia_criteria | cardiac/hestia_criteria.json:151 | W1_B3 | NOT PORTED | decision rule (any criterion); Zondag 2011 criteria not seen |
| 47 | hope_score | cardiac/hope_score.json:91 | W1_B3 | NOT PORTED | Pasquier 2018 (PMID 29481910): coefficients not in abstract |
| 48 | hr_zones | cardiac/hr_zones.go (Karvonen reserve × % + rest) | W1_B3 | NOT PORTED | held: doc/catalog/medicine.md §2/§5 (Karvonen 1957 not read); zone percentages cite Mayo Clinic (secondary); max HR 220 − age (Fox) not seen in a primary |
| 49 | improve_vte | cardiac/improve_vte.json:107 | W1_B3 | NOT PORTED | Spyropoulos 2011 Chest (PMID 21436241): weights not seen |
| 50 | improvedd_vte | cardiac/improvedd_vte.json:119 | W1_B3 | NOT PORTED | CORRECTED: Gibson TH Open 2017;1:e56 (not 2019); depends on IMPROVE weights (not seen) |
| 51 | isth_bat | cardiac/isth_bat.go:16-34 | W1_B3 | NOT PORTED | Rodeghiero 2010 (PMID 20626619, no abstract): per-item 0–4 scale and CNS item scale not seen; Elbatarny 2014 (PMC4251588) cut-offs seen. Finding kept: source scores an out-of-range index as 0 silently |
| 52 | jones_criteria | cardiac/jones_criteria.json:173 | W1_B3 | NOT PORTED | decision pathway (diagnostic criteria, Gewitz 2015 AHA, PMID 25908771), not a calculation |
| 53 | kawasaki | cardiac/kawasaki.json:118 | W1_B3 | NOT PORTED | decision pathway (diagnostic criteria, Newburger 2004 AHA, PMID 15505111), not a calculation |
| 54 | killip_classification | cardiac/killip_classification.json:41 | W1_B3 | NOT PORTED | classification chosen by the clinician from exam findings (Killip & Kimball 1967, PMID 6059183); nothing to compute |
| 55 | map_calc | cardiac/map_calc.go ((sys+2·dia)/3) | W1_B3 | NOT PORTED | held: doc/catalog/medicine.md §5 (MAP one-third rule, no primary); cites Magder 2014 (PMID 24736333), a review |
| 56 | map_inverse | cardiac/map_inverse.json:60 | W1_B3 | NOT PORTED | held: doc/catalog/medicine.md §5 (inverse of the MAP one-third rule); cites StatPearls/GlobalRPH (secondary) |
| 57 | medped_fh | cardiac/medped_fh.json:62 | W1_B3 | NOT PORTED | Williams 1993 (PMID 8328379): full age × relative-degree table not seen |
| 58 | michigan_picc_dvt | cardiac/michigan_picc_dvt.json:85 | W1_B3 | NOT PORTED | Chopra 2017 (PMID 28796444): points and class ranges not in abstract |
| 59 | modified_sgarbossa | cardiac/modified_sgarbossa.json:64 | W1_B3 | NOT PORTED | Smith 2012 (PMID 22939607) abstract gives the ST/S ≤ −0.25 criterion but not whether the concordant criteria are kept; Sgarbossa 1996 points 5/3/2 not seen |
| 60 | modified_shock_index | cardiac/modified_shock_index.json:58 | W1_B3 | NOT PORTED | cited PMC6499869 and PMID 25400383 are unrelated papers; Liu 2012 (PMID 25215048) abstract does not define MSI |
| 61 | ohfrs | cardiac/ohfrs.json:139 | W1_B3 | NOT PORTED | Stiell 2013 (PMID 23570474): 10 elements, points not in abstract |
| 62 | orbi_risk_score | cardiac/orbi_risk_score.json:151 | W1_B3 | NOT PORTED | CORRECTED: Eur Heart J 2018;39:2090 (not Am J Cardiol); points not in abstract |
| 63 | orthostatic | cardiac/orthostatic.go (two subtractions) | W1_B3 | NOT PORTED | Freeman 2011 consensus (PMID 21431947) not seen (no abstract, paywall); entry cites Cleveland Clinic |
| 64 | oxygen_pulse | cardiac/oxygen_pulse.json:48 | W1_B3 | NOT PORTED | VO₂/HR: cites Cleveland Clinic + AHA 2010 CPET statement (publisher 403, not read); formula not seen in a primary |
| 65 | pcp_hf | cardiac/pcp_hf.json:155 | W1_B3 | NOT PORTED | Khan 2019 (PMID 31097157): full text not reached |
| 66 | prevent_cvd | cardiac/prevent_cvd.json:145 | W1_B3 | NOT PORTED | Khan 2024 Circulation (PMID 37947085): coefficients not read (supplement) |
| 67 | primacy | cardiac/primacy.json:133 | W1_B3 | NOT PORTED | Miron 2020 (PMC7365676): equation in the supplement only |
| 68 | pulse_classify | cardiac/pulse_classify.go (returns hr; bands in app) | W1_B3 | NOT PORTED | resting-HR bands cite Cleveland Clinic (secondary); no primary table seen |
| 69 | pulse_pressure | cardiac/pulse_pressure.go (sys − dia; PP/sys; 0.25·sys) | W1_B3 | NOT PORTED | held: doc/catalog/medicine.md §5 (pulse pressure family); cites Cleveland Clinic (secondary) |
| 70 | pvr_calc | cardiac/pvr_calc.json:58 | W1_B3 | NOT PORTED | cites Cleveland Clinic/Mayo/StatPearls (secondary) and the 2022 ESC/ERS PH guideline (publisher page, not read); (mPAP − PAWP)/CO not seen in a primary in this run; same family as the HELD SVR row (doc/catalog/medicine.md §5) |
| 71 | qtc | cardiac/qtc.json:63 | W1_B3 | NOT PORTED | output is QTc Bazett: held: doc/catalog/medicine.md §2/§5 (Bazett 1920 not read). A Fridericia QTc is available as MED-CARDIO-002 |
| 72 | rasi_calc | cardiac/rasi_calc.json:58 | W1_B3 | NOT PORTED | duplicate of `rasi` with a Mayo Clinic stress-test page as citation (wrong page, secondary); see the `rasi` row |
| 73 | readmits_ami | cardiac/readmits_ami.json:107 | W1_B3 | NOT PORTED | Nguyen 2018 JAHA (PMC6015397) seen, but how "Age (per decade >18 y)" accrues points is not stated verbatim |
| 74 | rescue_ihca | cardiac/rescue_ihca.json:95 | W1_B3 | NOT PORTED | CORRECTED: JACC Cardiovasc Interv 2022;15:237 (PMC8837656); Table 2 seen via summariser, but rounding of age and duration (per 10) not stated and reference categories inferred |
| 75 | rome4_chest_pain | cardiac/rome4_chest_pain.go:15 | W1_B2 | NOT PORTED | Rome IV primary (Gastroenterology 2016;150(6), Rome Foundation ©) not open (Europe PMC: no full text); checklist of caller-judged criteria with no published number; criteria text licensed (D-M8) |
| 76 | rose_syncope | cardiac/rose_syncope.json:107 | W1_B3 | NOT PORTED | decision rule (BRACES: any item present), no score (Reed 2010, PMID 20170806) |
| 77 | rpp | cardiac/rpp.go (hr × sbp) | W1_B3 | NOT PORTED | held: doc/catalog/medicine.md §5 (Gobel 1978 abstract says "HR x blood pressure", not systolic) |
| 78 | save_score_ecmo | cardiac/save_score_ecmo.json:173 | W1_B3 | NOT PORTED | Schmidt 2015 EHJ (PMID 26033984): points not in abstract |
| 79 | score2_cvd | cardiac/score2_cvd.json:103 | W1_B3 | NOT PORTED | SCORE2 2021 (PMID 34120177): coefficients in supplement, not read |
| 80 | score2_diabetes | cardiac/score2_diabetes.json:133 | W1_B3 | NOT PORTED | SCORE2-Diabetes 2023 (PMC10361012): main text gives SHRs and baseline survival only; algorithm in Supplementary Methods Table 1 |
| 81 | score2_op | cardiac/score2_op.json:113 | W1_B3 | NOT PORTED | SCORE2-OP 2021 (PMID 34120185): coefficients in supplement, not read |
| 82 | seattle_heart_failure | cardiac/seattle_heart_failure.json:280 | W1_B3 | NOT PORTED | Seattle HF Model (Levy 2006 Circulation, PMID 16534009): multivariable model with medication/device adjustments partly from other trials; full coefficient set not read in this run (Circulation full text not open in PMC; not attempted) |
| 83 | sex_shock_score | cardiac/sex_shock_score.json:149 | W1_B3 | NOT PORTED | Wang 2024 EHJ (PMC11560280): points in Supplementary Table S13 only |
| 84 | sf_syncope_rule | cardiac/sf_syncope_rule.json:85 | W1_B3 | NOT PORTED | decision rule (CHESS: any item present), no score (Quinn 2004, PMID 14747812) |
| 85 | shock_index | cardiac/shock_index.json:50 | W1_B3 | NOT PORTED | Allgöwer & Burri 1967 (PMID 5299769, German, no abstract): HR/SBP not seen in the primary |
| 86 | spesi | cardiac/spesi.json:93 | W1_B3 | NOT PORTED | Jiménez 2010 (PMID 20696966): thresholds not in abstract |
| 87 | steinhart_ahf | cardiac/steinhart_ahf.json:61 | W1_B3 | NOT PORTED | CORRECTED: J Am Coll Cardiol 2009;54:1515 (not CMAJ); model not in abstract |
| 88 | svr_calc | cardiac/svr_calc.json:68 | W1_B3 | NOT PORTED | held: doc/catalog/medicine.md §5 (SVR/TPR; Guyton & Hall not opened; also uses MAP one-third rule) |
| 89 | t_macs | cardiac/t_macs.json:106 | W1_B3 | NOT PORTED | Body 2017 EMJ (PMC5502241): equation printed only as an image; the text formula under Table 1 is the original MACS model (variables a–h), not T-MACS |
| 90 | timi_stemi | cardiac/timi_stemi.json:118 | W1_B3 | NOT PORTED | Morrow 2000 Circulation 102:2031, PMID 11044416 confirmed; item points not in abstract, full text 403 |
| 91 | vancouver_chest_pain | cardiac/vancouver_chest_pain.json:95 | W1_B3 | NOT PORTED | decision rule (Scheuermeyer 2014, PMID 24626115), no score |
| 92 | villalta_pts | cardiac/villalta_pts.json:162 | W1_B3 | NOT PORTED | Kahn 2009 is a measurement-properties communication, not the original Villalta 1994 scale; items not seen |
| 93 | virsta_score | cardiac/virsta_score.json:140 | W1_B3 | NOT PORTED | Tubiana 2016 (PMID 26916042): 10 predictors, points not in abstract |
| 94 | vo2_oer | cardiac/vo2_oer.json:61 | W1_B3 | NOT PORTED | Fick VO₂ = CO × (CaO₂ − CvO₂) × 10 and O₂ER: cites StatPearls/SSC 2021 (no formula in SSC); CO family HELD (medicine.md §5) |
| 95 | vte_bleed | cardiac/vte_bleed.json:95 | W1_B3 | NOT PORTED | Klok 2016 ERJ (PMID 27471209): items and points not seen (403) |
| 96 | wells_dvt | cardiac/wells_dvt.json:140 | W1_B3 | NOT PORTED | Wells 1995/1997/2003 abstracts give no item points; full texts closed |
| 97 | wifi_classification | cardiac/wifi_classification.json:62 | W1_B3 | NOT PORTED | Mills 2014 J Vasc Surg (PMID 24126108): stage lookup table not seen |
| 98 | years_algorithm | cardiac/years_algorithm.json:85 | W1_B3 | NOT PORTED | decision pathway (van der Hulle 2017 Lancet, PMID 28549662: items + D-dimer 1000/500 ng/mL thresholds seen in abstract) |
