# credit-card-payoff-calculator

FIELDS (form `calform`): camountcanpay (500/month total budget), per card k: cnm_k name, cbal_k balance (4,600; 3,900;
6,000), cmpy_k minimum payment (100; 90; 120), cint_k rate (18.99 %; 19.99 %; 15.99 %), up to 10+ cards.

MATCH: NONE (multi-debt payoff: every card gets its minimum, the extra goes to one card chosen by a rule —
highest rate first ("avalanche") or smallest balance first ("snowball") — and freed minimums roll over).
RESEARCH: the per-card balance recursion is the amortization step we have (FIN-LOAN-003). The allocation rules are
consumer-education conventions; no regulation, textbook or standard defining them was found this session (the
CFPB's consumer pages on debt repayment are the likely place). NEEDS SOURCE; nothing built, no submissions.
