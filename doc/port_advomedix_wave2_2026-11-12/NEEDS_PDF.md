# NEEDS PDF (port wave 2)

Drop each paper into `papers/` named `<PMID>.pdf` (the folder is gitignored); the next run reads it and ports the entry. Exact citations per builder section below.

### Lead: withdrawn after review (D-075, D-076): these are NEEDS PDF; drop the primary to release them
- bp_classify, bp_adult_18_to_60, bp_adult_60_plus, bp_classifier.go: Whelton PK, Carey RM, Aronow WS, et al. 2017 ACC/AHA … Guideline for the Prevention, Detection, Evaluation, and Management of High Blood Pressure in Adults. Hypertension. 2018;71(6):e13-e115. doi:10.1161/HYP.0000000000000065. PMID 29133356
- FraminghamGeneralCVD10y: D'Agostino RB Sr, Vasan RS, Pencina MJ, et al. General cardiovascular risk profile for use in primary care. Circulation. 2008;117(6):743-753. doi:10.1161/CIRCULATIONAHA.107.699579. PMID 18212285 (Table 2)
- braden_score: Bergstrom 1987 Nurs Res 36(4):205-210, PMID 3299278 · mentzer_index: Mentzer 1973 Lancet 1(7808):882, PMID 4123424
- hospital_score, dipss, scorma_index, temp_mean_body, oconner_1rm, air_score, pucai (bands), decaf_copd, caspar_psa, aldrete_score, rcvs2_score, mnihss: primary citations in each builder section's row and NEEDS PDF list; the failing secondaries are named in out/verify_A.md and out/verify_B.md.

### W1_B1
- abi_calc: Aboyans V, Criqui MH, Abraham P, et al. Measurement and interpretation of the ankle-brachial index. Circulation. 2012;126(24):2890-2909. doi:10.1161/CIR.0b013e318276fbcb. PMID 23159553 — PMID 23159553
- adapt_protocol: Than M, Cullen L, Aldous S, et al. J Am Coll Cardiol. 2012;59(23):2091-2098. doi:10.1016/j.jacc.2012.02.035. PMID 22578923. Explicit rule (TIMI 0 + non-ischaemic ECG + 0/2 h troponin) expected in Methods; would reuse MED-CARDIO-004 — PMID 22578923
- altitude_adjusted_perc: Wolf SJ, McCubbin TR, Nordenholz KE, Naviaux NW, Haukoos JS. Am J Emerg Med. 2008;26(2):181-185. doi:10.1016/j.ajem.2007.04.026. PMID 18272098; PERC: Kline JA, et al. J Thromb Haemost. 2004;2(8):1247-1255. doi:10.1111/j.1538-7836.2004.00790.x. PMID 15304025 — PMID 18272098
- ascvd_risk: Goff DC Jr, Lloyd-Jones DM, Bennett G, et al. 2013 ACC/AHA guideline on the assessment of cardiovascular risk. Circulation. 2014;129(25 Suppl 2):S49-S73. doi:10.1161/01.cir.0000437741.48606.98. PMID 24222018 (Pooled Cohort coefficients, Table A) — PMID 24222018
- bova_score: Bova C, Sanchez O, Prandoni P, et al. Eur Respir J. 2014;44(3):694-703. doi:10.1183/09031936.00006114. PMID 24696111 — PMID 24696111
- brugada_criteria: Brugada P, Brugada J, Mont L, Smeets J, Andries EW. A new approach to the differential diagnosis of a regular tachycardia with a wide QRS complex. Circulation. 1991;83(5):1649-1659. doi:10.1161/01.cir.83.5.1649. PMID 2022022 (four-step rule; steps 3–4 morphology criteria not seen) — PMID 2022022
- cahp_score: Maupain C, Bougouin W, Lamhaut L, et al. Eur Heart J. 2016;37(42):3222-3228. doi:10.1093/eurheartj/ehv556. PMID 26497161 — PMID 26497161
- canadian_syncope: Thiruganasambandamoorthy V, Kwong K, Wells GA, et al. CMAJ. 2016;188(12):E289-E298. doi:10.1503/cmaj.151469. PMID 27378464. Score range −3 to 11 and bands (≤ −2 very low; −1 to 3 low-medium; ≥ 4 high/very high) seen in text; per-item points not confirmed (a summariser's transcription summed to −3…14, contradicting the paper) — PMID 27378464
- cardiac_index: Hall JE, Hall ME. Guyton and Hall Textbook of Medical Physiology. 14th ed. Elsevier; 2021. ISBN 978-0-323-59712-8 (book not opened), ch. 9 (CI = CO/BSA); BSA Du Bois/Mosteller also unconfirmed (see report) — PMID none (book or unindexed)
- cardiac_output: Hall JE, Hall ME. Guyton and Hall Textbook of Medical Physiology. 14th ed. Elsevier; 2021. ISBN 978-0-323-59712-8 (book not opened); assumed VO₂: LaFarge CG, Miettinen OS. Cardiovasc Res. 1970;4(1):23-30. doi:10.1093/cvr/4.1.23. PMID 5416840 — PMID 5416840
- chads65: Andrade JG, Aguilar M, Atzema C, et al. The 2020 Canadian Cardiovascular Society/Canadian Heart Rhythm Society comprehensive guidelines for the management of atrial fibrillation. Can J Cardiol. 2020;36(12):1847-1948. doi:10.1016/j.cjca.2020.09.001. PMID 33191198 (CHADS-65 algorithm figure) — PMID 33191198
- chadsvasc: Lip GY, Nieuwlaat R, Pisters R, Lane DA, Crijns HJ. Chest. 2010;137(2):263-272. doi:10.1378/chest.09-1584. PMID 19762550 (point table) — PMID 19762550
- cpes_score: CORRECTED: cited PMID 37100484 is Pibarot P, JACC 2023;81:1675 (unrelated). Nearest instrument: Zhang RS, et al. Pulmonary Embolism Shock Score. Circ Cardiovasc Interv. 2024;17:e014088. doi:10.1161/CIRCINTERVENTIONS.124.014088. PMID 38994599 — PMID 37100484
- cri: Wilkoff BL, Miller RE. Exercise testing for chronotropic assessment. Cardiol Clin. 1992;10(4):705-717. doi:10.1016/S0733-8651(18)30211-X. PMID 1423382 — PMID 1423382
- dash_prediction_vte: Tosetto A, Iorio A, Marcucci M, et al. J Thromb Haemost. 2012;10(6):1019-1025. doi:10.1111/j.1538-7836.2012.04735.x. PMID 22489957 — PMID 22489957
- do2: West JB, Luks AM. West's Respiratory Physiology. 11th ed. 2021. ISBN 978-1-975139-18-6 (CaO₂ 1.34 vs 1.39); Hall JE, Hall ME. Guyton and Hall Textbook of Medical Physiology. 14th ed. Elsevier; 2021. ISBN 978-0-323-59712-8 (book not opened); Shoemaker 1988 (PMID 3191758) is a trial, not the formula — PMID 3191758
- do2_simple: as do2 (CO × CaO₂ × 10); cites StatPearls (secondary) — PMID none (book or unindexed)
- doac_score: Aggarwal R, Ruff CT, Virdone S, et al. Circulation. 2023;148(12):936-946. doi:10.1161/CIRCULATIONAHA.123.064556. PMID 37621213. PMC10529708 — PMID 37621213
- edacs: Than M, Flaws D, Sanders S, et al. Emerg Med Australas. 2014;26(1):34-44. doi:10.1111/1742-6723.12164. PMID 24428678 — PMID 24428678
- egsys_syncope: Del Rosso A, Ungar A, Maggi R, et al. Heart. 2008;94(12):1620-1626. doi:10.1136/hrt.2008.143123. PMID 18519550 — PMID 18519550
- ehmrg: Lee DS, Stitt A, Austin PC, et al. Ann Intern Med. 2012;156(11):767-775. doi:10.7326/0003-4819-156-11-201206050-00003. PMID 22665814 — PMID 22665814
- euromacs_rhf: Soliman OII, Akin S, Muslem R, et al. Circulation. 2018;137(9):891-906. doi:10.1161/CIRCULATIONAHA.117.030543. PMID 28847897 — PMID 28847897
- four_ts_hit: Lo GK, Juhl D, Warkentin TE, Sigouin CS, Eichler P, Greinacher A. J Thromb Haemost. 2006;4(4):759-765. doi:10.1111/j.1538-7836.2006.01787.x. PMID 16634744 — PMID 16634744
- framingham_hard_chd: Wilson PW, D'Agostino RB, Levy D, Belanger AM, Silbershatz H, Kannel WB. Circulation. 1998;97(18):1837-1847. doi:10.1161/01.cir.97.18.1837. PMID 9603539 (also NCEP ATP III 2002 point tables, NHLBI PDF — not attempted) — PMID 9603539
- framingham_hf: McKee PA, Castelli WP, McNamara PM, Kannel WB. N Engl J Med. 1971;285(26):1441-1446. doi:10.1056/NEJM197112232852601. PMID 5122894 (explicit rule: 2 major or 1 major + 2 minor; criteria list not seen) — PMID 5122894
- garfield_af: Fox KAA, Lucas JE, Pieper KS, et al. BMJ Open. 2017;7(12):e017157. doi:10.1136/bmjopen-2017-017157. PMID 29273652 — PMID 29273652
- go_far_score: Ebell MH, Jang W, Shen Y, Geocadin RG. JAMA Intern Med. 2013;173(20):1872-1878. doi:10.1001/jamainternmed.2013.10037. PMID 24018585 — PMID 24018585
- grace_acs: Granger CB, Goldberg RJ, Dabbous O, et al. Arch Intern Med. 2003;163(19):2345-2353. doi:10.1001/archinte.163.19.2345. PMID 14581255; app variables match Eagle KA, et al. JAMA. 2004;291(22):2727-2733. doi:10.1001/jama.291.22.2727. PMID 15187054 — PMID 14581255
- grogan_attr_cm: Grogan M, Scott CG, Kyle RA, et al. J Am Coll Cardiol. 2016;68(10):1014-1020. doi:10.1016/j.jacc.2016.06.033. PMID 27585505 (whether 0.05 ng/mL / 3000 pg/mL are inclusive) — PMID 27585505
- has_bled: Pisters R, Lane DA, Nieuwlaat R, de Vos CB, Crijns HJ, Lip GY. Chest. 2010;138(5):1093-1100. doi:10.1378/chest.10-0134. PMID 20299623 — PMID 20299623
- hcm_risk_scd: O'Mahony C, Jichi F, Pavlou M, et al. Eur Heart J. 2014;35(30):2010-2020. doi:10.1093/eurheartj/eht439. PMID 24126876; formula also in Elliott PM, et al. 2014 ESC HCM guideline. Eur Heart J. 2014;35(39):2733-2779. doi:10.1093/eurheartj/ehu284. PMID 25173338 — PMID 24126876
- he_macs: Alghamdi A, Howard L, Reynard C, et al. Eur J Emerg Med. 2019;26(5):356-361. doi:10.1097/MEJ.0000000000000575. PMID 30289775 — PMID 30289775
- hemorr2hages: Gage BF, Yan Y, Milligan PE, et al. Am Heart J. 2006;151(3):713-719. doi:10.1016/j.ahj.2005.04.017. PMID 16504638 — PMID 16504638
- hep_score_hit: Cuker A, Arepally G, Crowther MA, et al. J Thromb Haemost. 2010;8(12):2642-2650. doi:10.1111/j.1538-7836.2010.04059.x. PMID 20854372 — PMID 20854372
- hestia_criteria: Zondag W, Mos IC, Creemers-Schild D, et al. J Thromb Haemost. 2011;9(8):1500-1507. doi:10.1111/j.1538-7836.2011.04388.x. PMID 21645235 (11 criteria; any present) — PMID 21645235
- hope_score: Pasquier M, Hugli O, Paal P, et al. Resuscitation. 2018;126:58-64. doi:10.1016/j.resuscitation.2018.02.026. PMID 29481910 — PMID 29481910
- hr_zones: Karvonen MJ, Kentala E, Mustala O. Ann Med Exp Biol Fenn. 1957;35(3):307-315. PMID 13470504; zone percentages cite Mayo Clinic (secondary); 220 − age not seen in a primary (MED-CARDIO-001/029/030 offer sourced HRmax) — PMID 13470504
- improve_vte: Spyropoulos AC, Anderson FA Jr, FitzGerald G, et al. Chest. 2011;140(3):706-714. doi:10.1378/chest.10-1944. PMID 21436241 — PMID 21436241
- improvedd_vte: Gibson CM, Spyropoulos AC, Cohen AT, et al. TH Open. 2017;1(1):e56-e65. doi:10.1055/s-0037-1603929. PMID 31249911. Needs Spyropoulos 2011 (PMID 21436241) weights — PMID 31249911
- isth_bat: Rodeghiero F, Tosetto A, Abshire T, et al. J Thromb Haemost. 2010;8(9):2063-2065. doi:10.1111/j.1538-7836.2010.03975.x. PMID 20626619; cut-offs Elbatarny M, et al. Haemophilia. 2014;20(6):831-835, PMC4251588 — PMID 20626619
- jones_criteria: Gewitz MH, Baltimore RS, Tani LY, et al. Circulation. 2015;131(20):1806-1818. doi:10.1161/CIR.0000000000000205. PMID 25908771 (explicit rule, population-risk-specific criteria; table not seen) — PMID 25908771
- kawasaki: McCrindle BW, Rowley AH, Newburger JW, et al. Circulation. 2017;135(17):e927-e999. doi:10.1161/CIR.0000000000000484. PMID 28356445 (current); Newburger JW, et al. Circulation. 2004;110(17):2747-2771. doi:10.1161/01.CIR.0000145143.19711.78. PMID 15505111 (cited). Explicit rule (fever ≥5 d + ≥4 of 5 principal features) not seen in a primary — PMID 28356445
- map_calc: Hall JE, Hall ME. Guyton and Hall Textbook of Medical Physiology. 14th ed. Elsevier; 2021. ISBN 978-0-323-59712-8 (book not opened), ch. 15; Magder SA. Crit Care Med. 2014;42(5):1241-1251. doi:10.1097/CCM.0000000000000324. PMID 24736333 is a review — PMID 24736333
- map_inverse: as map_calc (inverse of the one-third rule); cites StatPearls/GlobalRPH (secondary) — PMID none (book or unindexed)
- medped_fh: Williams RR, Hunt SC, Schumacher MC, et al. Am J Cardiol. 1993;72(2):171-176. doi:10.1016/0002-9149(93)90155-6. PMID 8328379 (age × relative-degree table) — PMID 8328379
- michigan_picc_dvt: Chopra V, Kaatz S, Conlon A, et al. J Thromb Haemost. 2017;15(10):1951-1962. doi:10.1111/jth.13794. PMID 28796444 — PMID 28796444
- modified_sgarbossa: Smith SW, Dodd KW, Henry TD, Dvorak DM, Pearce LA. Ann Emerg Med. 2012;60(6):766-776. doi:10.1016/j.annemergmed.2012.07.119. PMID 22939607 — PMID 22939607
- ohfrs: Stiell IG, Clement CM, Brison RJ, et al. Acad Emerg Med. 2013;20(1):17-26. doi:10.1111/acem.12056. PMID 23570474 — PMID 23570474
- orbi_risk_score: Auffret V, Cottin Y, Leurent G, et al. Eur Heart J. 2018;39(22):2090-2102. doi:10.1093/eurheartj/ehy127. PMID 29554243 — PMID 29554243
- orthostatic: Freeman R, Wieling W, Axelrod FB, et al. Clin Auton Res. 2011;21(2):69-72. doi:10.1007/s10286-011-0119-5. PMID 21431947 — PMID 21431947
- oxygen_pulse: Balady GJ, Arena R, Sietsema K, et al. Clinician's guide to cardiopulmonary exercise testing in adults. Circulation. 2010;122(2):191-225. doi:10.1161/CIR.0b013e3181e52e69. PMID 20585013 — PMID 20585013
- prevent_cvd: Khan SS, Matsushita K, Sang Y, et al. Circulation. 2024;149(6):430-449. doi:10.1161/CIRCULATIONAHA.123.067626. PMID 37947085 (supplemental coefficient tables) — PMID 37947085
- primacy: Miron A, Lafreniere-Roula M, Steve Fan CP, et al. Circulation. 2020;142(3):217-229. doi:10.1161/CIRCULATIONAHA.120.047235. PMID 32418493 — PMID 32418493
- pulse_classify: Entry cites Cleveland Clinic only; candidate primary for brady/tachy cut-offs: Kusumoto FM, et al. 2018 ACC/AHA/HRS bradycardia guideline. Circulation. 2019;140(8):e382-e482. doi:10.1161/CIR.0000000000000628. PMID 30586772 (uses < 50 bpm, which may contradict the app's < 60) — PMID 30586772
- pulse_pressure: Hall JE, Hall ME. Guyton and Hall Textbook of Medical Physiology. 14th ed. Elsevier; 2021. ISBN 978-0-323-59712-8 (book not opened); the 0.25 × systolic output has no source found — PMID none (book or unindexed)
- pvr_calc: Humbert M, Kovacs G, Hoeper MM, et al. 2022 ESC/ERS PH guidelines. Eur Heart J. 2022;43(38):3618-3731. doi:10.1093/eurheartj/ehac237. PMID 36017548; Eur Respir J. 2023;61(1):2200879. doi:10.1183/13993003.00879-2022. PMID 36028254 — PMID 36017548
- qtc: Bazett HC. Heart. 1920;7:353-370; reprint Ann Noninvasive Electrocardiol. 1997;2(2):177-194. doi:10.1111/j.1542-474X.1997.tb00325.x. (Fridericia available: MED-CARDIO-002) — PMID none (book or unindexed)
- rose_syncope: Reed MJ, Newby DE, Coull AJ, Prescott RJ, Jacques KG, Gray AJ. J Am Coll Cardiol. 2010;55(8):713-721. doi:10.1016/j.jacc.2009.09.049. PMID 20170806 — PMID 20170806
- rpp: Gobel FL, Norstrom LA, Nelson RR, Jorgensen CR, Wang Y. Circulation. 1978;57(3):549-556. doi:10.1161/01.CIR.57.3.549. PMID 624164 (abstract says "HR x blood pressure"; systolic not seen) — PMID 624164
- save_score_ecmo: Schmidt M, Burrell A, Roberts L, et al. Eur Heart J. 2015;36(33):2246-2256. doi:10.1093/eurheartj/ehv194. PMID 26033984 — PMID 26033984
- score2_cvd: SCORE2 working group and ESC Cardiovascular risk collaboration. Eur Heart J. 2021;42(25):2439-2454. doi:10.1093/eurheartj/ehab309. PMID 34120177 — PMID 34120177
- score2_diabetes: SCORE2-Diabetes Working Group and the ESC Cardiovascular Risk Collaboration. Eur Heart J. 2023;44(28):2544-2556. doi:10.1093/eurheartj/ehad260. PMID 37247330 — PMID 37247330
- score2_op: SCORE2-OP working group and ESC Cardiovascular risk collaboration. Eur Heart J. 2021;42(25):2455-2467. doi:10.1093/eurheartj/ehab312. PMID 34120185 — PMID 34120185
- seattle_heart_failure: Levy WC, Mozaffarian D, Linker DT, et al. Circulation. 2006;113(11):1424-1433. doi:10.1161/CIRCULATIONAHA.105.584102. PMID 16534009 — PMID 16534009
- sex_shock_score: Wang Y, Zeller M, Auffret V, et al. Eur Heart J. 2024;45(43):4564-4578. doi:10.1093/eurheartj/ehae593. PMID 39217456 — PMID 39217456
- shock_index: Allgöwer M, Burri C. ["Shock index"]. Dtsch Med Wochenschr. 1967;92(43):1947-1950. doi:10.1055/s-0028-1106070. PMID 5299769 — PMID 5299769
- spesi: Jiménez D, Aujesky D, Moores L, et al. Arch Intern Med. 2010;170(15):1383-1389. doi:10.1001/archinternmed.2010.199. PMID 20696966 — PMID 20696966
- steinhart_ahf: Steinhart B, Thorpe KE, Bayoumi AM, Moe G, Januzzi JL Jr, Mazer CD. J Am Coll Cardiol. 2009;54(16):1515-1521. doi:10.1016/j.jacc.2009.05.065. PMID 19815122 — PMID 19815122
- svr_calc: Hall JE, Hall ME. Guyton and Hall Textbook of Medical Physiology. 14th ed. Elsevier; 2021. ISBN 978-0-323-59712-8 (book not opened), ch. 14 (SVR = 80 × (MAP − CVP)/CO); also uses the unsourced MAP one-third rule — PMID none (book or unindexed)
- t_macs: Body R, Carlton E, Sperrin M, et al. Emerg Med J. 2017;34(6):349-356. doi:10.1136/emermed-2016-205983. PMID 27565197 — PMID 27565197
- timi_stemi: Morrow DA, Antman EM, Charlesworth A, et al. Circulation. 2000;102(17):2031-2037. doi:10.1161/01.cir.102.17.2031. PMID 11044416 — PMID 11044416
- vancouver_chest_pain: Scheuermeyer FX, Wong H, Yu E, et al. CJEM. 2014;16(2):106-119. doi:10.2310/8000.2013.130938. PMID 24626115 — PMID 24626115
- villalta_pts: Kahn SR, Partsch H, Vedantham S, Prandoni P, Kearon C. J Thromb Haemost. 2009;7(5):879-883. doi:10.1111/j.1538-7836.2009.03294.x. PMID 19175497 (ISTH definition with Villalta scoring); original Villalta S, et al. Haemostasis. 1994;24(suppl 1):158a (abstract, not indexed) — PMID 19175497
- virsta_score: Tubiana S, Duval X, Alla F, et al. J Infect. 2016;72(5):544-553. doi:10.1016/j.jinf.2016.02.003. PMID 26916042 — PMID 26916042
- vo2_oer: Fick VO₂ = CO × (CaO₂ − CvO₂) × 10 and O₂ER: Hall JE, Hall ME. Guyton and Hall Textbook of Medical Physiology. 14th ed. Elsevier; 2021. ISBN 978-0-323-59712-8 (book not opened); West's Respiratory Physiology 11th ed. (CaO₂ constant); cites StatPearls / SSC 2021 (no formula) — PMID none (book or unindexed)
- vte_bleed: Klok FA, Barco S, Konstantinides SV. Eur Respir J. 2016;48(5):1369-1376. doi:10.1183/13993003.00280-2016. PMID 27471209 — PMID 27471209
- wells_dvt: Wells PS, Anderson DR, Rodger M, et al. N Engl J Med. 2003;349(13):1227-1235. doi:10.1056/NEJMoa023153. PMID 14507948; Wells PS, et al. Lancet. 1997;350(9094):1795-1798. PMID 9428249 — PMID 14507948
- wifi_classification: Mills JL Sr, Conte MS, Armstrong DG, et al. J Vasc Surg. 2014;59(1):220-234.e2. doi:10.1016/j.jvs.2013.08.003. PMID 24126108 (stage lookup table) — PMID 24126108

### W1_B2
| entry | citation | PMID |
|---|---|---|
| aa_gradient | Riley RL, Cournand A. J Appl Physiol. 1949;1:825 (alveolar gas equation) | 18145478 |
| abg_analyzer | Baillie JK. Thorax. 2008;63(3):289-290. doi:10.1136/thx.2007.091223 | 18308967 |
| acute_asthma_severity | NAEPP Expert Panel Report 3, NIH Publication No. 07-4051, 2007, Section 5 (exacerbation severity figure) | — |
| asthma_predictive_index | Castro-Rodriguez JA, et al. Am J Respir Crit Care Med. 2000;162(4):1403-1406. doi:10.1164/ajrccm.162.4.9912111; Guilbert TW, et al. J Allergy Clin Immunol. 2004;114(6):1282-1287. doi:10.1016/j.jaci.2004.09.020 | 11029352; 15577824 |
| bap65 | Tabak YP, et al. Arch Intern Med. 2009;169(17):1595-1602. doi:10.1001/archinternmed.2009.270 | 19786679 |
| bcrss_covid | Duca A, et al. Emerg Med Pract. 2020;22(5 Suppl):CD1-CD2. doi:10.1097/01.emp.0000698131.66206.c8 | 32297727 |
| blast_lung_injury | Pizov R, et al. Chest. 1999;115(1):165-172. doi:10.1378/chest.115.1.165 | 9925079 |
| cao2 | West JB. Respiratory Physiology: The Essentials. 11th ed. 2021, ch. 6 (O2 content; 1.34 vs 1.39) | — |
| covid_gram | Liang W, et al. JAMA Intern Med. 2020;180(8):1081-1089. doi:10.1001/jamainternmed.2020.2033 — the supplement (β coefficients) | 32396163 |
| egri_difficult_airway | el-Ganzouri AR, et al. Anesth Analg. 1996;82(6):1197-1204. doi:10.1097/00000539-199606000-00017 | 8638791 |
| epiphany_index | Carmona-Bayonas A, et al. Br J Cancer. 2017;116(8):994-1001. doi:10.1038/bjc.2017.48 — the tree figure | 28267709 |
| gold_copd | GOLD 2024 Report, ch. 2 figures (grades; ABE tool), goldcopd.org | — |
| heaven_criteria | Davis DP, Olvera DJ. Air Med J. 2017;36(4):195-197. doi:10.1016/j.amj.2017.04.001 | 28739243 |
| newsom_score | Newsom C, et al. Acad Emerg Med. 2018;25(6):650-656. doi:10.1111/acem.13386 | 29427301 |
| paradise_criteria | Paradise JL, et al. N Engl J Med. 1984;310(11):674-683. doi:10.1056/NEJM198403153101102 | 6700642 |
| pass_asthma | Gorelick MH, et al. Acad Emerg Med. 2004;11(1):10-18. doi:10.1197/j.aem.2003.07.015 | 14709423 |
| pediatric_asthma_score | Kelly CS, Andersen CL, Pestian JP. Ann Allergy Asthma Immunol. 2000;84(5):509-516 | 10831004 |
| pediatric_ett_size | Cole F. Pediatric formulas for the anesthesiologist. AMA J Dis Child. 1957;94:672 (unconfirmed) | — |
| pef_exacerbation | BTS/SIGN British guideline on the management of asthma (SIGN 158), 2019 | — |
| pef_predicted_eu | Gregg I, Nunn AJ. Br Med J. 1973;3(5874):282-284. doi:10.1136/bmj.3.5874.282 (PMC1586779, scan) | 4723469 |
| pef_reversibility | NICE/BTS/SIGN NG245 (2024), bronchodilator reversibility | — |
| pef_variability | GINA 2024 report, Box 1-2 | — |
| piama_risk_score | Caudri D, et al. J Allergy Clin Immunol. 2009;124(5):903-910. doi:10.1016/j.jaci.2009.06.045 | 19665765 |
| pneumonia_risk_score | Lipsett SC, et al. Pediatr Infect Dis J. 2022;41(1):24-30. doi:10.1097/INF.0000000000003361 | 34694254 |
| pram_asthma | Chalut DS, et al. J Pediatr. 2000;137(6):762-768. doi:10.1067/mpd.2000.110121 | 11113831 |
| priest_covid | Goodacre S, et al. PLoS One. 2021;16(1):e0245840 (PMC7822515) — Fig 2 needs transcription | 33481930 |
| psi_calc, psi_port_score | Fine MJ, et al. N Engl J Med. 1997;336(4):243-250. doi:10.1056/NEJM199701233360402 | 8995086 |
| rads_inhalation | Oh JS, et al. J Burn Care Res. 2012;33(4):532-538. doi:10.1097/BCR.0b013e318237455f | 22210063 |
| rdos_score | Campbell ML, et al. J Palliat Med. 2010;13(3):285-290. doi:10.1089/jpm.2009.0229 | 20078243 |
| resp_ecmo | Schmidt M, et al. Am J Respir Crit Care Med. 2014;189(11):1374-1382. doi:10.1164/rccm.201311-2023OC | 24693864 |
| respiratory_score_asthma | Liu LL, et al. Pediatr Pulmonol. 2004;37(3):243-248. doi:10.1002/ppul.10425 | 14966818 |
| rse_asthma | Bateman ED, et al. J Allergy Clin Immunol. 2015;135(6):1457-1464.e4. doi:10.1016/j.jaci.2014.08.015 | 25258144 |
| scap_score | España PP, et al. Am J Respir Crit Care Med. 2006;174(11):1249-1256. doi:10.1164/rccm.200602-177OC | 16973986 |
| smart_cop | Charles PG, et al. Clin Infect Dis. 2008;47(3):375-384. doi:10.1086/589754 | 18558884 |
| va_calc, ve_calc, vdvt_calc | West JB. Respiratory Physiology: The Essentials. 11th ed. 2021, ch. 2 (ventilation, dead space) | — |
| westley_croup | Westley CR, Cotton EK, Brooks JG. Am J Dis Child. 1978;132(5):484-487. doi:10.1001/archpedi.1978.02120300044008 | 347921 |
| winters_formula | Albert MS, Dell RB, Winters RW. Ann Intern Med. 1967;66(2):312-322. doi:10.7326/0003-4819-66-2-312 | 6016545 |

### W1_B3
- bwps_thyrotoxicosis: Burch HB, Wartofsky L. Endocrinol Metab Clin North Am. 1993;22(2):263-277. PMID 8325286.
- canrisk: Robinson CA, Agarwal G, Nerenberg K. Chronic Dis Inj Can. 2011;32(1):19-31. doi:10.24095/hpcdp.32.1.04. PMID 22153173 (Table 3).
- caroc_system: Siminoski K et al. J Clin Densitom. 2007;10(2):120-123. PMID 17485028; Leslie WD et al. Osteoporos Int 2011, PMID 20967422.
- child_pugh: Pugh RN et al. Br J Surg. 1973;60(8):646-649. doi:10.1002/bjs.1800600817. PMID 4541913 (and a source for the INR form).
- fibrotic_nash_index: Tavaglione F et al. Clin Gastroenterol Hepatol. 2023;21(6):1523-1532.e1. doi:10.1016/j.cgh.2022.03.044. PMID 35421583.
- ho_index: Ho GT et al. Aliment Pharmacol Ther. 2004;19(10):1079-1087. doi:10.1111/j.1365-2036.2004.01945.x. PMID 15142197.
- hypoglycemia_risk: Karter AJ et al. JAMA Intern Med. 2017;177(10):1461-1470. PMID 28828479 (Figure 1 tree).
- ideal_body_weight: McCarron MM, Devine BJ. Drug Intell Clin Pharm. 1974;8(11):650-655. doi:10.1177/106002807400801104; Pai MP, Paloucek FP. Ann Pharmacother. 2000;34(9):1066-1069. PMID 10981254.
- kings_college_apap: O'Grady JG et al. Gastroenterology. 1989;97(2):439-445. PMID 2490426.
- nirudak_score: Levine AC et al. PLoS Negl Trop Dis. 2021;15(3):e0009266. doi:10.1371/journal.pntd.0009266. PMID 33690646 (Table 5 image).
- orai_osteoporosis: Cadarette SM et al. CMAJ. 2000;162(9):1289-1294. PMID 10813010.
- osiris_osteoporosis: Sedrine WB et al. Gynecol Endocrinol. 2002;16(3):245-250. PMID 12192897.
- peld_cr: Hsu E et al. Am J Transplant. 2021;21(1):222-228. doi:10.1111/ajt.15925. PMID 32306489; OPTN Policy 9 PDF.
- peld_score: McDiarmid SV et al. Transplantation. 2002;74(2):173-181. doi:10.1097/00007890-200207270-00006. PMID 12151728.
- pfib_c: Lefere S et al. Hepatology. 2025;81(4):1276-1287. doi:10.1097/HEP.0000000000001016. PMID 39028885.
- rebound_hyperbilirubinemia: Chang PW et al. Pediatrics. 2017;139(3):e20162896. PMID 28196932.

### W1_B4
- abuse_assessment_screen: McFarlane J, Parker B, Soeken K, Bullock L. JAMA. 1992;267(23):3176-3178. doi:10.1001/jama.267.23.3176. PMID 1593739
- ace_score: Dong M, Anda RF, Felitti VJ, et al. Child Abuse Negl. 2004;28(7):771-784. doi:10.1016/j.chiabu.2004.01.008. PMID 15261471
- age_adjusted_esr_crp (ESR part): Miller A, Green M, Robinson D. BMJ (Clin Res Ed). 1983;286(6361):266. doi:10.1136/bmj.286.6361.266. PMID 6402065
- aldrete_score (1970 and 1998 versions): Aldrete JA, Kroulik D. Anesth Analg. 1970;49(6):924-934. PMID 5534693; Aldrete JA. J Perianesth Nurs. 1998;13(3):148-155. PMID 9801540
- apache_ii: Knaus WA, Draper EA, Wagner DP, Zimmerman JE. Crit Care Med. 1985;13(10):818-829. PMID 3928249
- bops: Hesselgard K, et al. Pediatr Crit Care Med. 2007;8(2):102-108. doi:10.1097/01.PCC.0000257098.32268.AA. PMID 17273124
- breast_cancer_brain_met: Graesslin O, et al. J Clin Oncol. 2010;28(12):2032-2037. doi:10.1200/JCO.2009.24.6314. PMID 20308667
- brue_2: Nama N, et al. Hosp Pediatr. 2022;12(9):772-785. doi:10.1542/hpeds.2022-006637. PMID 35965279
- brue_criteria: Tieder JS, et al. Pediatrics. 2016;137(5):e20160590. doi:10.1542/peds.2016-0590. PMID 27244835
- cedocs: Weiss SJ, et al. Am J Emerg Med. 2014;32(11):1357-1363. doi:10.1016/j.ajem.2014.08.035. PMID 25234796
- cheops: McGrath PJ, et al. CHEOPS. Advances in Pain Research and Therapy vol 9. Raven Press; 1985:395-402 (book chapter, no PMID)
- cirs_g: Miller MD, et al. Psychiatry Res. 1992;41(3):237-248. doi:10.1016/0165-1781(92)90005-n. PMID 1594710
- corrected_reticulocyte_rpi: Hillman RS. J Clin Invest. 1969;48(3):443-453. doi:10.1172/JCI106001. PMID 5773082 (PMC535708 is a page scan)
- dipss (DIPSS Plus part): Gangat N, et al. J Clin Oncol. 2011;29(4):392-397. doi:10.1200/JCO.2010.32.2446. PMID 21149668
- dire_score: Belgrade MJ, Schamber CD, Lindgren BR. J Pain. 2006;7(9):671-681. doi:10.1016/j.jpain.2006.03.001. PMID 16942953
- esas_r: Watanabe SM, et al. J Pain Symptom Manage. 2011;41(2):456-468. doi:10.1016/j.jpainsymman.2010.04.020. PMID 20832987
- flipi24: Maurer MJ, et al. J Clin Oncol. 2026;44(2):117-128. doi:10.1200/JCO-25-00892. PMID 41329901 (Data Supplement Table S5)
- fracture_index: Black DM, et al. Osteoporos Int. 2001;12(7):519-528. doi:10.1007/s001980170072. PMID 11527048
- hiri_msm: Smith DK, et al. J Acquir Immune Defic Syndr. 2012;60(4):421-427. doi:10.1097/QAI.0b013e318256b2f6. PMID 22487585
- hits_domestic_violence: Sherin KM, et al. Fam Med. 1998;30(7):508-512. PMID 9669164
- katz_adl: Katz S, et al. JAMA. 1963;185:914-919. doi:10.1001/jama.1963.03060120024016. PMID 14044222
- mipss70: Guglielmelli P, et al. J Clin Oncol. 2018;36(4):310-318. doi:10.1200/JCO.2017.76.4886. PMID 29226763; Tefferi A, et al. J Clin Oncol. 2018;36(17):1769-1770. doi:10.1200/JCO.2018.78.9867. PMID 29708808
- (BUILD IN, provisional) air_score PMID 18553045; alvarado_appendicitis PMID 3963537; audit_c PMID 9738608; cage_alcohol PMID 6471323; cas_covid_anxiety PMID 32299304; ciwa_b PMID 2574193; cows_score PMID 12924748; dast_10 (10-item source); glasgow_blatchford PMID 11073021; hospital_score PMID 23529115; kruis_ibs PMID 6724251; myxedema_coma PMID 24518183

### W2_B5
- kruis_ibs: Kruis W, Thieme C, Weinzierl M, Schüssler P, Holl J, Paulus W. Gastroenterology. 1984;87(1):1-7. PMID 6724251
- pcdai: Hyams JS, Ferry GD, Mandel FS, et al. J Pediatr Gastroenterol Nutr. 1991;12(4):439-447. doi:10.1097/00005176-199105000-00005. PMID 1678008 (haematocrit thresholds)
- parc_appendicitis: Kharbanda AB, et al. Pediatrics. 2018;141(4):e20172699. doi:10.1542/peds.2017-2699. PMID 29535251 (Table 3 age coding + supplement; PMC5869337)
- rapid_pleural: Rahman NM, et al. Chest. 2014;145(4):848-855. doi:10.1378/chest.13-1558. PMID 24264558
- neonatal_hyperbilirubinemia: Bhutani VK, Johnson L, Sivieri EM. Pediatrics. 1999;103(1):6-14. doi:10.1542/peds.103.1.6. PMID 9917432; preferred Kemper AR, et al. Pediatrics. 2022;150(3):e2022058859. doi:10.1542/peds.2022-058859. PMID 35927462 (supplement tables)
- nutritional_risk_index: Buzby GP, et al. Am J Clin Nutr. 1988;47(2 Suppl):366-381. doi:10.1093/ajcn/47.2.366. PMID 3124598; VA TPN Cooperative Study Group. N Engl J Med. 1991;325(8):525-532. doi:10.1056/NEJM199108223250801. PMID 1906987
- ovat_ipv: Ernst AA, Weiss SJ, Cham E, Hall L, Nick TG. Violence Vict. 2004;19(3):375-384. doi:10.1891/vivi.19.3.375.65769. PMID 15631287
- pedocs: Weiss SJ, Ernst AA, Sills MR, Quinn BJ, Johnson A, Nick TG. Pediatr Emerg Care. 2007;23(9):641-645. doi:10.1097/PEC.0b013e31814a69e2. PMID 17876254
- saps_ii: Le Gall JR, Lemeshow S, Saulnier F. JAMA. 1993;270(24):2957-2963. doi:10.1001/jama.270.24.2957. PMID 8254858
- smast: Selzer ML, Vinokur A, van Rooijen L. J Stud Alcohol. 1975;36(1):117-126. doi:10.15288/jsa.1975.36.117. PMID 238068
- smast_g: SAMHSA TIP 26 (PEP20-02-01-011), NBK571032 (item text); Blow FC, et al. SMAST-G, RSA meeting 1998 (abstract, no PMID)
- temp_abg_correction: Ashwood ER, Kost G, Kenny M. Clin Chem. 1983;29(11):1877-1885. PMID 6354511
- temp_liebermeister: Seneta E, Seif FJ, Liebermeister H, Dietz K. J Med Biogr. 2004;12(4):215-221. doi:10.1177/096777200401200411. PMID 15486620
- (partial) myxedema_coma bands: Popoveniuc G, et al. Endocr Pract. 2014;20(8):808-817. doi:10.4158/EP13460.OR. PMID 24518183; nhfs 2012 equation: Moppett IK, et al. Br J Anaesth. 2012;109(4):546-550. doi:10.1093/bja/aes187. PMID 22728204

### W2_B6
- amt10: Hodkinson HM. Age Ageing. 1972;1(4):233-238. doi:10.1093/ageing/1.4.233. PMID 4669880
- amt4: Swain DG, Nightingale PG. Clin Rehabil. 1997;11(3):243-248. doi:10.1177/026921559701100308. PMID 9360037
- capd: Traube C, et al. Crit Care Med. 2014;42(3):656-663. doi:10.1097/CCM.0b013e3182a66b76. PMID 24145848 (Figure 1)
- chalice_rule: Dunning J, et al. Arch Dis Child. 2006;91(11):885-891. doi:10.1136/adc.2005.083980. PMID 17056862
- chip_prediction: Smits M, et al. Ann Intern Med. 2007;146(6):397-405. doi:10.7326/0003-4819-146-6-200703200-00004. PMID 17371884
- ciwa_ar: Sullivan JT, et al. Br J Addict. 1989;84(11):1353-1357. doi:10.1111/j.1360-0443.1989.tb00737.x. PMID 2597811
- ciwa_b: Busto UE, Sykora K, Sellers EM. J Clin Psychopharmacol. 1989;9(6):412-416. PMID 2574193
- cts_6: Graham B, et al. J Hand Surg Am. 2006;31(6):919-924. PMID 16886290
- func_score: Rost NS, et al. Stroke. 2008;39(8):2304-2309. doi:10.1161/STROKEAHA.107.512202. PMID 18556582
- gmaws: McPherson A, Benson G, Forrest EH. QJM. 2012;105(7):649-656. doi:10.1093/qjmed/hcs020. PMID 22328545
- ham_a: Hamilton M. Br J Med Psychol. 1959;32(1):50-55. doi:10.1111/j.2044-8341.1959.tb00467.x. PMID 13638508
- ham_d: Hamilton M. J Neurol Neurosurg Psychiatry. 1960;23(1):56-62. doi:10.1136/jnnp.23.1.56. PMID 14399272
- hat_score: Lou M, et al. Neurology. 2008;71(18):1417-1423. doi:10.1212/01.wnl.0000330297.58334.dd. PMID 18955684 (Table 1)
- heads_ed: Cappelli M, et al. Pediatrics. 2012;130(2):e321-e327. doi:10.1542/peds.2011-3798. PMID 22826567
- hii_impairment: Hack JB, Goldlust EJ, Gibbs F, Zink B. Am J Drug Alcohol Abuse. 2014;40(2):111-117 (PMID to locate)
- impact_score: Steyerberg EW, et al. PLoS Med. 2008;5(8):e165. doi:10.1371/journal.pmed.0050165. PMID 18684008 (Figure 1 point chart)
- infant_scalp_score: Schutzman SA, Nigrovic LE, Mannix R. Acad Emerg Med. 2021;28(1):92-97. doi:10.1111/acem.14087. PMID 32673432
- major_depression_index: Bech P, et al. J Affect Disord. 2001;66(2-3):159-164. doi:10.1016/s0165-0327(00)00309-8. PMID 11578668
- mbig: Khan AD, et al. Trauma Surg Acute Care Open. 2020;5(1):e000483. doi:10.1136/tsaco-2020-000483. PMID 32537518 (Figure 2); Khan AD, et al. J Trauma Acute Care Surg. 2022;93(1):106-112. PMID 35358157
- midas_migraine: Stewart WF, et al. Neurology. 2001;56(6 Suppl 1):S20-S28. doi:10.1212/wnl.56.suppl_1.s20. PMID 11294956
- mrs_9q: Patel N, et al. Neurosurgery. 2012;71(5):971-975. doi:10.1227/NEU.0b013e31826a8a56. PMID 22843133
- nps_pain: Galer BS, Jensen MP. Neurology. 1997;48(2):332-338. doi:10.1212/wnl.48.2.332. PMID 9040716
- ped_nihss: Ichord RN, et al. Stroke. 2011;42(3):613-617. doi:10.1161/STROKEAHA.110.607192. PMID 21317270 (online supplement)
- pediatric_gcs: James HE. Pediatr Ann. 1986;15(1):16-22. doi:10.3928/0090-4481-19860101-05. PMID 3951884
- pediatric_nexus_ii: Oman JA, et al. Pediatrics. 2006;117(2):e238-e246. doi:10.1542/peds.2005-1063. PMID 16418311
- qids_depression: Rush AJ, et al. Biol Psychiatry. 2003;54(5):573-583. doi:10.1016/s0006-3223(02)01866-8. PMID 12946886
- rosier_scale: Nor AM, et al. Lancet Neurol. 2005;4(11):727-734. doi:10.1016/S1474-4422(05)70201-5. PMID 16239179
- setscore: Schönenberger S, et al. Neurocrit Care. 2016;25(1):94-104. doi:10.1007/s12028-015-0235-5. PMID 26842719
- sudbury_vertigo: Ohle R, et al. Ann Emerg Med. 2025;85(2):122-131. doi:10.1016/j.annemergmed.2024.06.003. PMID 39093245
- tdn_grade: Terrapon APR, et al. Neurosurgery. 2021;89(2):236-245. doi:10.1093/neuros/nyab121. PMID 33887774
- (BUILD IN, provisional) cns_ipi PMID 27382100; finnegan_nas PMID 1163358; rule_of_7s_lyme PMID 19403476 / 22184651; bacterial_meningitis_score points: a second PMC source citing PMID 17200475

### W2_B7
- ain_risk: Moledina DG, et al. Nephrol Dial Transplant. 2022;37(11):2214-2222. doi:10.1093/ndt/gfab346. PMID 34865148 (supplement: unrounded betas, SG coding).
- igan_prediction: Barbour SJ, et al. JAMA Intern Med. 2019;179(7):942-952. doi:10.1001/jamainternmed.2019.0600. PMID 30980653 (eTable 3 of the Supplement).
- mskcc_rcc_recurrence: Kattan MW, Reuter V, Motzer RJ, Katz J, Russo P. J Urol. 2001;166(1):63-67. doi:10.1016/S0022-5347(05)66077-6. PMID 11435824.
- ruis_incontinence: Sansoni J, Hawthorne G, Marosszeky N, Fleming G. Technical Manual … Revised Incontinence and Patient Satisfaction Tools. University of Wollongong; 2011 (ro.uow.edu.au/ahsri/616). No PMID.
- vurx: Kirsch AJ, et al. J Pediatr Urol. 2014;10(6):1249-1254. doi:10.1016/j.jpurol.2014.06.019. PMID 25511573.
- mews, mews_calc: Subbe CP, Kruger M, Rutherford P, Gemmel L. QJM. 2001;94(10):521-526. doi:10.1093/qjmed/94.10.521. PMID 11588210.
- neonatal_eos: Kuzniewicz MW, et al. Jt Comm J Qual Patient Saf. 2016;42(5):232-239. doi:10.1016/s1553-7250(16)42030-1. PMID 27066927 (and Puopolo 2011 PMID 22025590).
- pediatric_sirs_sepsis: Goldstein B, Giroir B, Randolph A. Pediatr Crit Care Med. 2005;6(1):2-8. doi:10.1097/01.PCC.0000149131.72248.E6. PMID 15636651.
- step_by_step_febrile: Gomez B, et al. Pediatrics. 2016;138(2):e20154381. doi:10.1542/peds.2015-4381. PMID 27382134; Mintegi S, et al. Emerg Med J. 2014;31(e1):e19-e24. doi:10.1136/emermed-2013-202449. PMID 23851127.
- bastion_classification: Jacobs N, et al. Injury. 2014;45(9):1422-1428. doi:10.1016/j.injury.2012.05.001. PMID 22613453.
- rhmp30: de Jong L, et al. J Bone Joint Surg Am. 2025;107(5):459-468. doi:10.2106/JBJS.23.01397. PMID 39836737.
- ribs_score: Buchholz CJ, et al. Am Surg. 2023;89(11):4668-4674. doi:10.1177/00031348221123087. PMID 36120831.
- rip_score_patella: Hevesi M, et al. Arthroscopy. 2019;35(2):537-543. doi:10.1016/j.arthro.2018.09.017. PMID 30612768.
- rule_of_nines: Wallace AB. Lancet. 1951;1(6653):501-504. doi:10.1016/s0140-6736(51)91975-7. PMID 14805109.
- slics: Vaccaro AR, et al. Spine. 2007;32(21):2365-2374. doi:10.1097/BRS.0b013e3181557b92. PMID 17906580.
- triss: Champion HR, Sacco WJ, Copes WS. J Trauma. 1995;38(1):94-95. doi:10.1097/00005373-199501000-00024. PMID 7745669; Champion HR, et al. J Trauma. 1989;29(5):623-629. PMID 2657085; Boyd CR, et al. J Trauma. 1987;27(4):370-378. PMID 3106646.

### W2_B8
- bishop_score: Bishop EH. Obstet Gynecol. 1964;24:266-268 — PMID 14199536
- dapre_strength: Knight KL. Am J Sports Med. 1979;7(6):336-337. doi:10.1177/036354657900700605 — PMID 507268
- dga_estimated_calorie_needs: USDA/HHS. Dietary Guidelines for Americans, 2020-2025, 9th ed., Appendix 2 Tables A2-1/A2-2 — no PMID (drop as DGA_2020-2025.pdf)
- grade_adjusted_pace: Minetti AE, et al. J Appl Physiol. 2002;93(3):1039-1046. doi:10.1152/japplphysiol.01177.2001 — PMID 12183501
- iota_simple_rules (risk % model only): Timmerman D, et al. Am J Obstet Gynecol. 2016;214(4):424-437. doi:10.1016/j.ajog.2016.01.007 — PMID 26800772
- lander_1rm: Lander J. NSCA J. 1984;6(6):60-61. doi:10.1519/0744-0049(1984)006<0060:MBOR>2.3.CO;2 — no PMID
- lombardi_1rm: Lombardi VP. Beginning Weight Training. W.C. Brown; 1989 — book, no PMID
- mayhew_1rm: Mayhew JL, Ball TE, Arnold MD, Bowen JC. J Appl Sport Sci Res. 1992;6(4):200-206. doi:10.1519/00124278-199211000-00002 — no PMID
- pregnancy_dates: ACOG CO 700. Obstet Gynecol. 2017;129(5):e150-e154. doi:10.1097/AOG.0000000000002046 — PMID 28426621
- puqe_score: Koren G, et al. Am J Obstet Gynecol. 2002;186(5 Suppl):S228-S231. doi:10.1067/mob.2002.123054 — PMID 12011891
- reid_colposcopic_index: Reid R, Scalzi P. Am J Obstet Gynecol. 1985;153(6):611-618. doi:10.1016/s0002-9378(85)80244-1 — PMID 2998190
- riegel_race_prediction: Riegel PS. Am Sci. 1981;69(3):285-290 — PMID 7235349
- two_for_two_rule: Haff GG, Triplett NT, eds. Essentials of Strength Training and Conditioning. 4th ed. Human Kinetics; 2016. ISBN 978-1-4925-0162-6 — book, no PMID
- vbac_flamm: Flamm BL, Geiger AM. Obstet Gynecol. 1997;90(6):907-910. doi:10.1016/s0029-7844(97)00531-0 — PMID 9397100

### W3_B9a
- pulse_* paediatric (7) and rr_* paediatric (7): Fleming S, Thompson M, Stevens R, et al. Lancet. 2011;377(9770):1011-1018. doi:10.1016/S0140-6736(10)62226-X. PMID 21411136 — the Web Appendix (PMC3789232, NIHMS54714-supplement-Web_Appendix.pdf: Web Tables 4 and 5). Drop as `papers/21411136_appendix.pdf`; the PMC link is behind a browser check this agent may not pass.
- computeARV: Mena L, Pintos S, Queipo NV, et al. J Hypertens. 2005;23(3):505-511. doi:10.1097/01.hjh.0000160205.81652.5a. PMID 15716690.
- bp_ankle_brachial_index: Aboyans V, et al. Circulation. 2012;126(24):2890-2909. doi:10.1161/CIR.0b013e318276fbcb. PMID 23159553.
- bp_orthostatic_change: Freeman R, et al. Clin Auton Res. 2011;21(2):69-72. doi:10.1007/s10286-011-0119-5. PMID 21431947.
- bp_map.go, bp_mean_arterial_pressure, bp_pulse_pressure, bp_pulse_pressure.go, bp_systemic_vascular_resistance, pulse_cardiac_index, pulse_cardiac_output: Hall JE, Hall ME. Guyton and Hall Textbook of Medical Physiology. 14th ed. Elsevier; 2021. ISBN 978-0-323-59712-8 (no PMID).
- rr_alveolar_ventilation, rr_minute_ventilation: West JB. Respiratory Physiology: The Essentials. 11th ed. 2021, ch. 2 (no PMID). rr_dead_space_ratio: Bohr 1891 / Enghoff 1938 (no PMID).
- infant_sleep_0_to_1/1_to_3/3_to_4_months: Hirshkowitz M, et al. Sleep Health. 2015;1(1):40-43. doi:10.1016/j.sleh.2014.12.010. PMID 29073412.
- pulse_chronotropic_response_index: Wilkoff BL, Miller RE. Cardiol Clin. 1992;10(4):705-717. PMID 1423382; Lauer MS, et al. JAMA. 1999;281(6):524-529. PMID 10022108. _beta_blocker: Khan MN, Pothier CE, Lauer MS. Am J Cardiol. 2005;96(9):1328-1333. doi:10.1016/j.amjcard.2005.06.082. PMID 16253608.
- pulse_qtc_female/male: Rautaharju PM, et al. Circulation. 2009;119(10):e241-e250. PMID 19228821; Bazett HC. Heart. 1920;7:353-370.
- pulse_shock_index: Allgöwer M, Burri C. Dtsch Med Wochenschr. 1967;92(43):1947-1950. PMID 5299769.
- spo2_osi_category: PALICC. Pediatr Crit Care Med. 2015;16(5):428-439. PMID 25647235 (PMC5253180, Figure 2).
- temperature_liebermeister_dissociation: Seneta E, et al. J Med Biogr. 2004;12(4):215-221. doi:10.1177/096777200401200411. PMID 15486620.
- (replacement table, not an entry state) paediatric BP: Flynn JT, et al. Pediatrics. 2017;140(3):e20171904. doi:10.1542/peds.2017-1904. PMID 28827377, Tables 4–5.

### W3_B9b
| entry | citation | PMID |
|---|---|---|
| X1 TABLE Hankinson 1999 | Hankinson JL, Odencrantz JR, Fedan KB. Spirometric reference values from a sample of the general U.S. population. Am J Respir Crit Care Med. 1999;159(1):179-187. doi:10.1164/ajrccm.159.1.9712108, Tables 4–6 | 9872837 |
| X3 MAGE | Service FJ, Molnar GD, Rosevear JW, Ackerman E, Gatewood LC, Taylor WF. Mean amplitude of glycemic excursions, a measure of diabetic instability. Diabetes. 1970;19(9):644-655. doi:10.2337/diab.19.9.644 | 5469118 |
| X4 FINDRISC | Lindström J, Tuomilehto J. The diabetes risk score: a practical tool to predict type 2 diabetes risk. Diabetes Care. 2003;26(3):725-731. doi:10.2337/diacare.26.3.725 | 12610029 |
| (optional, X2 second source) | Quanjer PH, Lebowitz MD, Gregg I, Miller MR, Pedersen OF. Peak expiratory flow: conclusions and recommendations of a Working Party of the European Respiratory Society. Eur Respir J Suppl. 1997;24:2S-8S | 9098701 |

### W3_B10
- EdwardsZone / EdwardsTRIMP: Edwards S. The Heart Rate Monitor Book. Sacramento: Fleet Feet Press; 1993 (book; no PMID).
  Home life/fitness.
- MAGE: Service FJ, Molnar GD, Rosevear JW, Ackerman E, Gatewood LC, Taylor WF. Diabetes. 1970;19(9):644-655.
  doi:10.2337/diab.19.9.644. PMID 5469118. Home med/endocrine.
- ARV, time-weighted form (only if wanted beside NUM-SERIES-009): Mena L, Pintos S, Queipo NV, Aizpúrua JA, Maestre G,
  Sulbarán T. A reliable index for the prognostic significance of blood pressure variability. J Hypertens.
  2005;23(3):505-511. doi:10.1097/01.hjh.0000160205.81652.5a. PMID 15716690.

### W3_B11
- nihss: Brott T, Adams HP Jr, Olinger CP, et al. Measurements of acute cerebral infarction: a clinical examination scale. Stroke. 1989;20(7):864-870. doi:10.1161/01.str.20.7.864. PMID 2749846 — and the NINDS NIH Stroke Scale form (NIH_Stroke_Scale_508C.pdf) for the UN rule.
- pawss: Maldonado JR, Sher Y, Das S, et al. Alcohol Alcohol. 2015;50(5):509-518. doi:10.1093/alcalc/agv043. PMID 25999438 (and Maldonado JR, et al. Alcohol. 2014;48(4):375-390. doi:10.1016/j.alcohol.2014.01.004. PMID 24657098).
- nhfs 2012 mortality equation: Moppett IK, Parker M, Griffiths R, Bowers T, White SM, Moran CG. Br J Anaesth. 2012;109(4):546-550. doi:10.1093/bja/aes187. PMID 22728204.
- (optional, to confirm SAPS 3 combination rule) Moreno 2005 ESM Appendix C, Intensive Care Med 2005;31(10):1345-1355, PMID 16132892.
