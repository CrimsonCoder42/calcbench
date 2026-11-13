# auto-lease-calculator

FIELDS (form `calform`): cloanamount auto price (50,000) / cmonthlypay (600) (solve either way), cloanterm (36 months),
cratetype (interestrate|moneyfactor), cinterestrate (0.00208 as money factor), cdownpayment (10,000),
ctradeinvalue (0), csaletax (7 %), cresidualvalue (24,000).

MATCH: NONE.
RESEARCH: Regulation M, 12 CFR 1013.4(f) "Payment calculation" (read 2026-11-13 at consumerfinance.gov/rules-policy/regulations/1013/4/)
lists the disclosure lines: gross capitalized cost, capitalized cost reduction, adjusted capitalized cost, residual
value, depreciation and amortized amounts (adjusted cap cost − residual), rent charge, total of base periodic
payments, lease payments, base periodic payment. The rent charge there is a disclosed amount (total of payments −
depreciation), not a formula; the regulation never mentions a "money factor". The usual lease-payment formula
rent = (adjusted cap cost + residual) × money factor, and money factor = APR/2400, are industry conventions; the
Federal Reserve's "Keys to Vehicle Leasing" pages (likely source) returned 403/404. Model form A-2 is an image only.
NEEDS SOURCE: the money-factor rent charge. Nothing built. No submissions.
