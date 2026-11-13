# fitness-and-health: discrepancies (A-FIT)

| slug | inputs | ours | theirs | verdict | citation / note |
|---|---|---|---|---|---|
| body-fat-calculator | m waist 34, neck 15, h 70 in | 17.5131 | 17.4 | THEIRS DIFFERS (code withdrawn, NEEDS SOURCE) | ours: DoD direct %fat eq., Foulis et al. Front Physiol 2023;14:1183836 Methods. Theirs reproduces the Hodgdon–Beckett 1984 body-density form + Siri (495/D − 450). |
| body-fat-calculator | f waist 30, hip 38, neck 13, h 65 in | 28.5563 | 28.3 | THEIRS DIFFERS (code withdrawn, NEEDS SOURCE) | same as above |
| bmi-calculator | 100 lb, 72 in | "Underweight" | "Severe thinness" | THEIRS DIFFERS (labels) | ours: CDC adult categories (MED-ANTHRO-003); theirs: WHO-style labels, same cut points |
| ideal-weight-calculator | man 70 in / woman 64 in, Hamwi | 166.0 / 120.0 lb | 165.3 / 119.7 lb | THEIRS DIFFERS (code withdrawn, NEEDS SOURCE) | ours: Hamwi lb form, Peterson et al. AJCN 2016;103:1197 Table 3. Theirs: metric restatement 48 + 2.7 kg/in, 45.5 + 2.2 kg/in |
| macro-calculator | m 25 y, 1740 kcal | AMDR protein 43.5–152.25, carb 195.75–282.75, fat 38.67–67.67 g | 62–149, 186–303, 40–69 g | THEIRS DIFFERS | ours: NASEM 2024 Letter Report Table 3-1 (IOM 2002/2005 AMDR), LIF-NUTR-003. Theirs: unstated scheme, not the AMDR % |
| carbohydrate-calculator | 2060.4 kcal, 65 % | 334.81 g | 357 g | THEIRS DIFFERS | ours: 4 kcal/g, 21 CFR 101.9(c)(1)(i)(B); theirs: 3.75 kcal/g (monosaccharide-equivalent factor); their 40–75 % band is not the AMDR 45–65 % |
| protein-calculator | m 25 y 160 lb 5ft10 ×1.375 | 59.02–206.58 g (10–35 % of 2360.9 kcal) | 58–202 g | THEIRS DIFFERS | same % rule (NASEM AMDR, LIF-NUTR-003); their implied energy ≈ 2310 kcal, basis not shown |
| fat-intake-calculator | 2060.4 kcal, 20–35 % | 45.79–80.13 g | 47–82 g | THEIRS DIFFERS | ours: 9 kcal/g, 21 CFR 101.9(c)(1)(i)(B); theirs reproduced by ≈ 8.8 kcal/g |
| gfr-calculator | m 50 y SCr 0.9 | 104.049 (CKD-EPI 2021) | 99.2 | THEIRS DIFFERS | ours: Inker NEJM 2021 (MED-RENAL-001); theirs: CKD-EPI 2009 with race term (Levey Ann Intern Med 2009) |
| due-date-calculator | LMP 2026-07-04, cycle 32 | EDD 2027-04-10, GA 18w6d | Apr 14, 2027, 18w2d | THEIRS DIFFERS | ours: ACOG CO 700 LMP + 280 d, no cycle adjustment stated; theirs adds (cycle − 28) days |
