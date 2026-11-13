# NIST Dataplot Reference Manual

**Citation**: NIST, *Dataplot Reference Manual*, https://www.itl.nist.gov/div898/software/dataplot/refman2/ (pages read 2026-11-12). Accepted as a source by D-010.

| Page | Taken | Used by |
|---|---|---|
| auxillar/coefvari.htm "COEFFICIENT OF VARIATION" | cv = s/x̄; relsd = 100(s/\|x̄\|) | NUM-STATS-003, -030 |
| ch2/relat_sd.pdf "RELATIVE STANDARD DEVIATION" | relsd = 100·(sd/\|mean\|) (cites Snedecor & Cochran 6th ed. pp. 62–64, not checked) | NUM-STATS-030 |
| auxillar/groupper.htm "GROUPED PERCENTILE" | p_i = l + ((iN/100) − F)/f·C, uniform-within-class note (cites Selby, CRC Standard Mathematical Tables, 22nd ed., 1974, p. 572, not checked) | NUM-STATS-006 |
| auxillar/correlat.htm "CORRELATION" | r = S_xy/(√S_xx·√S_yy) | NUM-STATS-014 |
| auxillar/ssqmean.htm "SUM OF SQUARES FROM MEAN" | SSQ = Σ(X − X̄)² | NUM-STATS-025 |
| ch2/weigmean.pdf "WEIGHTED MEAN" | x̄w = Σwx/Σw | NUM-STATS-026 |
| auxillar/weigtmea.htm "WEIGHTED TRIMMED MEAN" | "none of the weights can be negative", at least one positive | NUM-STATS-026 (weight rule only) |
| auxillar/zscore.htm "ISO 13528 ZSCORE" | Z = (X − Xref)/sigma, and solved for X (cites ISO 13528:2005 pp. 25–26) | NUM-STATS-028, -029 |
| auxillar/mad.htm "MEDIAN ABSOLUTE DEVIATION" | MAD = median\|x − x̃\|, MADN = MAD/0.6745 | NUM-STATS-031 |
| ch2/sd_mean.pdf "STANDARD DEVIATION OF THE MEAN" | s/√N | NUM-STATS-023 |

How checked: each page fetched and its formula compared with the code; worked examples in the tests. Status: confirmed. Note: weigtmea.htm is the weighted TRIMMED mean page; the plain weighted mean is ch2/weigmean.pdf (the stats catalog names weigtmea.htm for the plain mean, which is wrong).

## Matrix gaps (2026-11-12)

"CORRELATION": correlation = covariance/(s_x·s_y). Used by NUM-LINALG-167, 168.
