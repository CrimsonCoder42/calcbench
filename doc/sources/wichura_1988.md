# Wichura 1988 (AS 241)

**Citation**: Wichura, "Algorithm AS 241: The Percentage Points of the Normal Distribution", *Journal of the Royal Statistical Society, Series C (Applied Statistics)*, 1988, 37(3):477–484, doi:10.2307/2347330.

Taken: PPND16: the split points 0.425 and 5, the constants 0.180625 and 1.6, the coefficients A0–A7, B1–B7, C0–C7, D1–D7, E0–E7, F1–F7. Used by NUM-PROB-061 NormalQuantile, NUM-PROB-062 LognormalQuantile.

Checked: Crossref (vol. 37, no. 3, first page 477; the end page 484 is not in Crossref). The coefficients came from the journal's published listing at StatLib (lib.stat.cmu.edu/apstat/241, header "ALGORITHM AS241 APPL. STATIST. (1988) VOL. 37, NO. 3"); two fetches agreed digit for digit, and the copy in normal_quantile.go was compared digit by digit. Tests: NIST critical values and a 2000-point round trip through NormalCDF and NormalSurvival in all three regions. Status: confirmed.
