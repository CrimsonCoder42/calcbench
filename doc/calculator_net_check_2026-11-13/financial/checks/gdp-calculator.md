# gdp-calculator (row 72, HANDOFF from A-OTHER; their note: other/checks/gdp-calculator.md)

FIELDS: form `calcf` (expenditure): personalconsumption, grossinvestment, governmentconsumption, exportv, importv (all
blank by default). Form `calcf2` (income approach): employeecompensation, proprietorsincome, rents,
corporateprofits, interestincome, indirectbusinesstaxes, depreciation, netincomeofforeigners. The income approach is
a plain sum and was not built.

MATCH: NONE (FIN-ECON-001..004 hold no GDP identity).
RESEARCH: Greenlaw, Shapiro and MacDonald, Principles of Macroeconomics 3e, OpenStax 2022, §6.1: "GDP = C + I + G +
(X – M)"; Table 6.1 (U.S. 2020, $ trillion): 14.0 + 3.6 + 3.9 + 2.1 − 2.7 = 20.9 (D-018).
BUILD: FIN-ECON-005 GDPExpenditure, fin/economics/gdp_expenditure.go:19.

COMPARE (scratch: checks/scratch/gdp-calculator.go)

| inputs | ours | theirs | diff | verdict |
|---|---|---|---|---|
| Table 6.1: 14.0, 3.6, 3.9, 2.1, 2.7 | 20.9 | 20.9 | 0 | MATCH |
| 1,250,000.55; −30,000; 400,000; 150,000; 220,000.25 | 1,550,000.30 | 1,550,000.3 | 0 | MATCH |

Submissions: 2.
