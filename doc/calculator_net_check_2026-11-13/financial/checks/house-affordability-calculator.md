# house-affordability-calculator

FIELDS: form `calform` (by income): cannualincome (120,000), cloanterm (30 y), cinterestrate (7.263 %),
cmonthlydebt (0), cdownpercent (20) + cdownpaymentunit (p|d), cpropertytax (1.5 %/yr of price) + unit,
choa (0) + unit, cinsurance (0.5 %/yr) + unit, cdtiunit (cv = 28/36 | fha | va | 10…50).
Form `calform2` (by budget): c2budget (3,500/month), c2loanterm, c2interestrate, c2downpercent + unit,
c2coverfees (true), c2propertytax, c2hoa, c2insurance, c2maintenance (1.5 %/yr) + units.

MATCH: PARTIAL FIN-LOAN-010 AffordableLoanAmount (housing ratio, fixed monthly taxes and insurance, PV of the
allowance). Their costs are % of the unknown price, so the price solves P = R·F/(d + c·F/12), F = (P/A,i,n),
d = loan share of price, c = annual cost rate. Not built: that wrapper, and the ratio limits (28/36, FHA, VA are
lender or agency policy; ours takes the ratio from the caller). They also add PMI at an assumed 0.5 %/yr of the
loan when down payment < 20 % (no input field; a lender quote, not a formula).

COMPARE (scratch: checks/scratch/house-affordability-calculator.go; monthly taxes and insurance taken from their price)

| inputs | ours | theirs | diff | verdict |
|---|---|---|---|---|
| 120k/yr, 28 %, 7.263 %, 30 y, 20 % down, tax 1.5 % + ins 0.5 % | loan 314,115.23 (T&I 654.41 rounded to the cent); self-consistent hand solve 314,115.60 | loan 314,116; price 392,644 | < $1 (T&I cent rounding) | MATCH |
| 80k/yr, 28 %, 5 %, 15 y, 10 % down, tax 1 % + ins 0.5 % + their PMI 80/mo | loan 192,172.46; hand solve 192,172.57 | 192,173; price 213,525 | < $1 | MATCH |

Submissions: 3 (the first printed only part of the result; repeated once with the default form to read it).
