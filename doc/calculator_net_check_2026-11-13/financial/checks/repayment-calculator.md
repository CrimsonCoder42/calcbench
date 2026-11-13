# repayment-calculator

FIELDS (form `calform`): cloanamount (10,000), cinterestrate (10 %), ccompound (annually…biweekly; monthly),
cpayback (daily…year; month), cpaybackway (t = fixed term | a = fixed amount), cyears (5) + cmonths (0),
cpaybackwayamount (200).

MATCH: HAVE FIN-INT-018 (fixed term), FIN-INT-043 AnnuityPeriods (fixed amount, built this session). Compounding
different from the payment frequency is not covered directly (convert the rate first).

COMPARE (scratch: checks/scratch/repayment-calculator.go)

| inputs | ours | theirs | diff | verdict |
|---|---|---|---|---|
| 10,000, 10 % monthly, 5 y | 212.47 | 212.47 | 0 | MATCH |
| 10,000, 10 %, 200/month | N = 64.949 → 65 months | not captured | — | DRIVER |

DRIVER: `run` with cpaybackwayamount fails: the locator resolves to `<table id="cpaybackwayamount">` before the
input of the same name ("Element is not an <input>"). Not retried.
Submissions: 2 (second one failed in the driver before submit).
