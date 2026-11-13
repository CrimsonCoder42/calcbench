# retirement-calculator

FIELDS: form `calform` (how much do I need): cagenow (35), cretireage (67), clifeexp (85), cincomenow (70,000/yr),
cincomeinc (3 %/yr), cretinclevel (75) + crilunit (p|d), cinvreturn (6 %), cinflation (3 %), cotherincome (0/month),
csavingnow (30,000), cannualsave (10) + csaveunit (p|d).
Form `calform3` (how much can I withdraw): cagenow, cretireage, clifeexp, csavingnow (30,000), cannualaddition (0),
cmonthlyaddition (500), cinvreturn (6 %), cinflation (3 %).
Form `calform4` (how long will it last): cneededamount (600,000), cmonthlywithdraw (5,000/month), cinvreturn (6 %).

MATCH: form 4 HAVE by composition: FIN-INT-044 (monthly rate equivalent to 6 % a year), FIN-INT-043 AnnuityPeriods,
FIN-INT-018. Forms 1 and 3 are multi-stage household plans (income growth, replacement ratio, inflation); D-021
plans them for `fin/personal` (PERSON) and no single primary source defines the plan. Not built; not compared.

COMPARE (scratch: checks/scratch/retirement-calculator.go)

| inputs | ours | theirs | diff | verdict |
|---|---|---|---|---|
| 600,000, withdraw 5,000/month, 6 %/yr | j = 0.486755 %/mo; N = 180.678 months = 15 y 0.68 mo | 15 years 0.7 months | 0 | MATCH |
| 600,000 over 10 years, 6 % | 6,613.44/month | 6,613.44 | 0 | MATCH |

Convention: their 6 % is an effective annual return converted to an equivalent monthly rate (the nominal 0.5 %/mo
would give 183.7 months and 6,661.23), withdrawals at month end.
Submissions: 1.
