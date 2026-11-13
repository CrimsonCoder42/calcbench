# gfr-calculator

## FIELDS (2026-11-13)
calform (adult): creatinine 0.9 · creatinineunit mgdL | molL · age 50 · gender male|female · race 1 (Black) | 0 (Not Black, default).
calform2 (child): ccreatinine 0.9 · ccreatinineunit mgdL | µmol/L · height 110 · heightunit inch | cm.
Results: MDRD 4-variable, "CKD-EPI", Mayo quadratic (0.1); child: one eGFR; a CKD stage table.

## MATCH
HAVE MED-RENAL-022 EGFRMDRD, MED-RENAL-001 EGFRCKDEPI2021, MED-RENAL-007 EGFRSchwartzBedside.
Not ours: CKD-EPI 2009 (race term) and the Mayo quadratic (Rice/Rule 2004) — not built; whether to add the superseded 2009 equation is a
decision for Devin (LEAD.md).

## COMPARE (scratch: checks/scratch/gfr-calculator.go)
| inputs | ours | theirs | verdict |
|---|---|---|---|
| m 50 y, SCr 0.9, not Black — MDRD | 89.321 | 89.3 | MATCH |
| f 70 y, SCr 1.4, Black — MDRD | 45.057 | 45.1 | MATCH |
| m 50 y, SCr 0.9 — CKD-EPI | 104.049 (2021, race-free) | 99.2 | THEIRS DIFFERS |
| f 70 y, SCr 1.4, Black — CKD-EPI | 40.473 (2021) | 43.9 | THEIRS DIFFERS |
| child 110 cm, SCr 0.9 — bedside Schwartz | 50.478 | 50.5 | MATCH |
Their CKD-EPI is the 2009 equation (Levey AS et al. Ann Intern Med 2009;150:604-612): 141 × min(SCr/0.9,1)^−0.411 × max(SCr/0.9,1)^−1.209
× 0.993^age = 99.27 for case 3 (hand). Ours is the 2021 refit without race (Inker LA et al. N Engl J Med 2021;385:1737-1749), MED-RENAL-001
Source. Ours stays.
Submissions: 1 fields + 3 runs.
