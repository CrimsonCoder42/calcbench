# margin-calculator

FIELDS: `calform` (profit margin): ccost (120), csale revenue (160), cmargin, cprofit (solve for any two).
`calform2` (stock margin): cprice (18.3), cshares (100), cmarginratio (30 %). `calform3` (forex margin): cexchange
(1.3), cratio leverage select (1…50; 20), cunits (100).

MATCH: PARTIAL FIN-BIZ-018 GrossProfitMargin; markup on cost was missing. Forms 2 and 3 are one product each
(price × shares × requirement; rate × units ÷ leverage); the requirement and leverage are broker or Reg T policy,
supplied by the caller. Not built.
RESEARCH: Kirk, Contemporary Mathematics, OpenStax 2023, §6.2 (markup = percent markup × cost; retail = cost × (1 +
markup)), Examples 6.10 and 6.16 (read 2026-11-13).
BUILD: FIN-BIZ-052 MarkupRate, fin/business/markup_rate.go:17; FIN-BIZ-053 RetailPriceFromMarkup,
fin/business/retail_price_from_markup.go:33; FIN-BIZ-054 SalePriceAfterDiscount (for discount-calculator),
fin/business/sale_price_after_discount.go:21.

COMPARE (scratch: checks/scratch/margin-calculator.go)

| inputs | ours | theirs | diff | verdict |
|---|---|---|---|---|
| cost 120, revenue 160 | margin 0.2500; markup 0.3333 | 25.00 %; 33.33 % | 0 | MATCH |
| cost 37.25, revenue 49.99 | margin 0.2549; markup 0.3420 | 25.49 %; 34.20 % | 0 | MATCH |
| stock 18.3 × 100 × 30 % (hand) | 549.00 | 549.00 | 0 | MATCH (arithmetic, not built) |

Submissions: 2.
