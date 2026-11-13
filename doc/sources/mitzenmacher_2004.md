# Broder and Mitzenmacher

**Citations** (as written in the doc comments):
- Broder and Mitzenmacher, "Network applications of Bloom filters: a survey", Internet Mathematics 1(4):485–509, 2004, doi:10.1080/15427951.2004.10129096, §2.1 (f = (1 − e^(−kn/m))^k)
- Broder and Mitzenmacher, "Network applications of Bloom filters: a survey", Internet Mathematics 1(4):485–509, 2004, doi:10.1080/15427951.2004.10129096, §2.1 (k = ln 2 · m/n minimises the rate) and §2.2 (rate ≤ ε requires m ≥ n·log₂(1/ε)/ln 2)
- Broder and Mitzenmacher, "Network applications of Bloom filters: a survey", Internet Mathematics 1(4):485–509, 2004, doi:10.1080/15427951.2004.10129096, §2.4 (the OR of two filters with the same bits and hash functions represents the union)

**Used by**: ALG-BLOOM-002, ALG-BLOOM-003, ALG-BLOOM-004.

**Checked**: the builder's catalog section in doc/algorithms_core_2026-11-12/catalog/ states how each citation was confirmed (Crossref DOI, table of contents, or the text read); citations not read first-hand are listed in doc/algorithms_core_2026-11-12/REPORT.md. The worked examples are reproduced in the functions' tests.
