# Wheatcroft 2021

**Citation**: Wheatcroft, "Forecasting football matches by predicting match statistics", *Journal of Sports Analytics*, 2021, 7(2):77–97, doi:10.3233/JSA-200462. Author preprint arXiv:2001.09097v1.

| Taken | Locator | Used by |
|---|---|---|
| Odds-implied probability rᵢ = 1/Oᵢ, so fair decimal odds are 1/p | §2.1 "Betting odds" | NUM-PROB-014 |
| Odds-implied probability rᵢ = 1/Oᵢ; example O = 3, 1.4 → 0.33, 0.71 | §2.1 "Betting odds" | NUM-PROB-015 |
| Overround π = Σᵢ 1/Oᵢ − 1 | §2.1 "Betting odds", eq. (1) | NUM-PROB-016 |

How checked: the §2.1 example is reproduced in `TestImpliedProbability_wheatcroftExample` (implied_probability_test.go:19) and `TestOverround_wheatcroftExample` (overround_test.go:18; π = 1/3 + 1/1.4 − 1 = 1/21).

Status: published version (journal, volume, issue, pages, DOI) confirmed on Crossref 2026-11-12. Locators read first-hand in the author preprint arXiv:2001.09097v1 (text extracted with pdftotext), 2026-11-12; the code says the locators are from the preprint. Note: Wheatcroft's "overround" is the booksum minus 1; Wunderlich–Memmert use the word for the booksum itself, so they are not cited for NUM-PROB-016 (question for Devin in fix_prob.md).
