# Kirsch and Mitzenmacher

**Citations** (as written in the doc comments):
- Kirsch and Mitzenmacher, "Less hashing, same performance: building a better Bloom filter", Random Structures & Algorithms 33(2):187–218, 2008, doi:10.1002/rsa.20208, §2 (standard Bloom filter; FPR (1 − ρ)^k given the fraction ρ of 0 bits) and §5.2 Theorem 5.2 ((enhanced) double hashing h₁ + i·h₂ + f(i) mod m, any m, any f)
- Kirsch and Mitzenmacher, Random Structures & Algorithms 33(2):187–218, 2008, doi:10.1002/rsa.20208, §2

**Used by**: ALG-BLOOM-001, ALG-BLOOM-002.

**Checked**: the builder's catalog section in doc/algorithms_core_2026-11-12/catalog/ states how each citation was confirmed (Crossref DOI, table of contents, or the text read); citations not read first-hand are listed in doc/algorithms_core_2026-11-12/REPORT.md. The worked examples are reproduced in the functions' tests.
