# salary-calculator

FIELDS (form `calform`): camount (50) + cunit (Hourly|Daily|Weekly|Bi-Weekly|Semi-Monthly|Monthly|Quarterly),
chours (40/week), cdays (5/week), cholidays (10/yr), cvacation (15/yr).

MATCH: NONE (D-021 plans `fin/personal`, SHORT PERSON, taxonomy row "planned").
RESEARCH: Fair Labor Standards Act regulation 29 CFR 778.113 (read 2026-11-13 at law.cornell.edu/cfr/text/29/778.113):
(a) the hourly rate of a weekly salary = salary ÷ hours it is meant to cover; (b) monthly × 12 ÷ 52 and semimonthly
× 24 ÷ 52 give the equivalent weekly wage. Their "holidays & vacation adjusted" column treats days off as unpaid
(annual − days off × daily pay); no source found for that convention, not built.
BUILD: new package fin/personal (doc.go copies fin/loan's shape, SHORT PERSON from the taxonomy):
FIN-PERSON-001 WeeklyPayEquivalent, fin/personal/weekly_pay_equivalent.go:38;
FIN-PERSON-002 HourlyRateFromWeeklyPay, fin/personal/hourly_rate_from_weekly_pay.go:22.

COMPARE (scratch: checks/scratch/salary-calculator.go)

| inputs | ours | theirs | diff | verdict |
|---|---|---|---|---|
| 50/h, 40 h/wk (unadjusted) | monthly 8,666.67 → weekly 2,000.00; semimonthly 4,333.33 → 2,000.00; 2,000/40 h → 50.00/h | weekly 2,000; monthly 8,667; semimonthly 4,333; hourly 50.00 | 0 at their whole dollars | MATCH |
| 6,500/month, 37.5 h/wk | weekly 1,500.00; hourly 40.00; biweekly 3,000.00 | 1,500; 40.00; 3,000 | 0 | MATCH |

Submissions: 2.
