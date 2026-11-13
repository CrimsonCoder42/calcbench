# refinance-calculator

FIELDS (form `calform`): cloandata select (b = know balance | t = know original loan), crbalance (250,000),
cmpay (1,800), cloanamount (300,000), cloanterm (30), cyears (20) + cmonths (0) remaining, cinterestrate (7 %),
cnyears (20), cninterestrate (6 %), cpoints (2 % of new loan), ccosts (1,500), ccashout (0).

MATCH: PARTIAL. Current loan's remaining term: FIN-INT-043 AnnuityPeriods (built this session for row 3). New
payment: FIN-INT-018. New loan's APR with points and fees as prepaid finance charges: FIN-INT-032 (Reg Z
Appendix J, amount financed = loan − points − fees). Not built: their "break-even" month (30) — the page gives no
definition that reproduces it (6,500 / 8.92 a month would be 729 months), and we have no source for theirs.

COMPARE (scratch: checks/scratch/refinance-calculator.go)

| inputs | ours | theirs | diff | verdict |
|---|---|---|---|---|
| new loan 250,000, 20 y, 6 % | pmt 1,791.08; total 429,859.20 (pmt × 240) | 1,791.08; 429,858.64 | pmt 0; total +0.56 | pmt MATCH; total THEIRS DIFFERS (unrounded payment × n, see mortgage-calculator) |
| same, 2 points + 1,500 fees | APR 6.329 % | 6.329 % | 0 | MATCH |
| current: 250,000, 1,800/mo, 7 % | N = 285.69 → 286 payments; balance after 285 payments 1,242.95 | 285 months; total 513,000 = 285 × 1,800 | 1 month | THEIRS DIFFERS |

Current-loan length: after 285 payments of 1,800 the balance is still 1,242.95 (FIN-LOAN-003 BalloonBalance), so a
286th, partial payment is due; theirs truncates N and understates total payments by about 1,250. Ours is the N
that solves (P/A,i,N) = PV/PMT (NCEES FE Reference Handbook 10.0.1, p. 230), confirmed by the balance recursion.
Submissions: 1.
