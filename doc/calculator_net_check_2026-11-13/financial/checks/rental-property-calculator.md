# rental-property-calculator

FIELDS (form `calform`): cprice (200,000), cuseloan (yes|no) + cdownpayment (20 %), cinterest (6 %), cloanterm (30),
cothercost closing (6,000), cneedrepair (no) + crepaircost, cafterrepairvalue; annual ctax (3,000), cinsurance (1,200),
choa (0), cmaintenance (2,000), cother (500), each with a % annual increase (3); crent (2,000/mo) + increase,
cotherincome (0) + increase, cvacancy (5 %), cmanagement (0 %), cknowsellprice (no) + cappreciation (3 %/yr) or
csellprice, cholding (20 y), csellcost (8 %).

MATCH: PARTIAL. IRR: FIN-INT-039. Mortgage payment: FIN-INT-018. Cap rate: NONE before this session.
RESEARCH: OCC Comptroller's Handbook, Commercial Real Estate Lending v2.0 (March 2022): Glossary p. 138 "Capitalization
rate: the ratio between a property's stabilized NOI and the property's sales price"; p. 43 value = NOI ÷ cap rate;
NOI = gross income less operating expenses, excluding interest, principal and income taxes (p. 140).
Cash-on-cash return: NEEDS SOURCE (the appraisal text that defines "equity dividend rate", Appraisal Institute, The
Appraisal of Real Estate, could not be opened; the OCC handbook does not define it). Theirs is cumulative over the
holding period (874.57 %), not the usual annual ratio, so it was not compared.

BUILD: FIN-INVEST-053 CapitalizationRate, fin/invest/capitalization_rate.go:19.

COMPARE (scratch: checks/scratch/rental-property-calculator.go)

| inputs | ours | theirs | diff | verdict |
|---|---|---|---|---|
| default: 200k, NOI 16,100 (first year) | cap 0.080500 | 8.05 % | 0 | MATCH |
| 350k cash, rent 2,800 + 100, vacancy 8 %, mgmt 10 %, tax 4,200, ins 1,500, HOA 600, maint 3,000 → NOI 19,514.40 | cap 0.055755 | 5.58 % | 0 | MATCH |
| same, 20 y hold, all lines +3 %/yr, sale at price·1.03²⁰ less 8 % | IRR (FIN-INT-039) 8.3309 % | 8.33 % | 0 | MATCH |

Submissions: 2.
