# discount-calculator

FIELDS (form `calform`): pricebefore (59.99), discount (15), priceafter, savedamount (solve for any two), format
(p = percent off | f = fixed amount off).

MATCH: NONE before this session. RESEARCH: Kirk, Contemporary Mathematics, OpenStax 2023, §6.2, Example 6.10
(sale price = original × (1 − percent discount)). BUILD: FIN-BIZ-054 SalePriceAfterDiscount,
fin/business/sale_price_after_discount.go:21 (built with the margin row). A fixed amount off is a subtraction.

COMPARE (scratch: checks/scratch/discount-calculator.go)

| inputs | ours | theirs | diff | verdict |
|---|---|---|---|---|
| 59.99, 15 % off | 50.99; saved 9.00 | 50.99; 9.00 | 0 | MATCH |
| 75.80, 25 % off (source Example 6.10) | 56.85; saved 18.95 | 56.85; 18.95 | 0 | MATCH |

Submissions: 2.
