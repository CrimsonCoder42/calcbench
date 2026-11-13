# Park LP

Citation (as in the doc comment): Park LP, Chu VH, Peterson G, et al. Validated risk score for predicting 6-month mortality in infective endocarditis. J Am Heart Assoc. 2016;5(4):e003016. doi:10.1161/JAHA.115.003016. PMID 27091179. PMC4859286, Table 5 "Simplified Risk Score Calculation for 6-Month Mortality in IE"

Used by: MED-CARDIO-035

How checked: read in the primary as the builder's ledger row records (the S14 port ledger (LEDGER.md in the S14 work folder listed in doc/README.md), the row of each ID above); tests reproduce its values; reviewer sample 2026-11-12.

Lockdown 2026-11-12: Table 5 (Europe PMC fullTextXML PMC4859286, table-wrap jah31437-tbl-0005, caption
"Simplified Risk Score Calculation for 6‐Month Mortality in IE") read row by row: Constant 4; Age ≤45 y 0;
46–60 y 2; 61–70 y 3; >70 y 4; History of dialysis 3; Nosocomial IE 2; Prosthetic IE 1; Symptoms >1 month
before admission −1; Staphylococcus aureus 1; Viridans group streptococci −2; Aortic vegetation 1; Mitral
vegetation 1; NYHA class 3 or 4 heart failure 3; Stroke 2; Paravalvular complication 2; Persistent bacteremia 2;
Surgical treatment −2. Footnote: "Probability of 6‐month mortality=2.416×score+0.109×score2−4.849". Every cell
equals MED-CARDIO-035's code. The probability footnote is not built.
