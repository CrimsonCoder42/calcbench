# cash-back-or-low-interest-calculator

FIELDS (form `calform`): ccashback (1,000), cinterestrate high (5 %), cinterestratelow (2 %), cloanamount = auto price
(50,000), cloanterm (60 months), cdownpayment (10,000), ctradeinvalue (0), cstate (blank|AL…), csaletax (7 %),
ctitlereg (2,000), cttrinloan (false).

MATCH: PARTIAL → composed. Each offer is a loan: FIN-INT-018 payment, FIN-LOAN-001 interest; the cash back reduces
the amount financed. The verdict is a subtraction of total costs; no comparison function is built.

COMPARE (scratch: checks/scratch/cash-back-or-low-interest-calculator.go)

| inputs | ours | theirs | diff | verdict |
|---|---|---|---|---|
| cash back: 39,000 at 5 %, 60 mo | pmt 735.98; interest 5,158.66 | 735.98; 5,158.69 | −0.03 | pmt MATCH; interest convention (see mortgage-calculator) |
| low rate: 40,000 at 2 %, 60 mo | pmt 701.11; interest 2,066.67 | 701.11; 2,066.62 | +0.05 | pmt MATCH; interest convention |
| cash back 3,000: 27,000 at 4 %, 36 mo | pmt 797.15; interest 1,697.30 | 797.15; 1,697.31 | −0.01 | MATCH |
| low rate: 30,000 at 2.9 %, 36 mo | pmt 871.11; interest 1,360.17 | 871.11; 1,360.13 | +0.04 | pmt MATCH; interest convention |

Both verdicts agree (low rate better in case 1; cash back better in case 2).
Submissions: 2.
