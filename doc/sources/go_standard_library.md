# Go standard library

**Citations** (as written in the doc comments):
- Go standard library, package iter documentation, "Pulling Values" (stop must be called when the caller is no longer interested in next values and next has not yet signaled that the sequence is over), https://pkg.go.dev/iter#Pull (read 2026-11-12)
- Go standard library, package iter documentation, https://pkg.go.dev/iter (read 2026-11-12)
- The Go standard library, package slices, func Compact (replaces consecutive runs of equal elements with a single copy), https://pkg.go.dev/slices#Compact (read 2026-11-12)

**Used by**: ALG-SEQ-001, ALG-SEQ-002, ALG-SEQ-003, ALG-SEQ-004, ALG-SEQ-005, ALG-SEQ-006, ALG-SEQ-007, ALG-TWOPTR-003.

**Checked**: the builder's catalog section in doc/algorithms_core_2026-11-12/catalog/ states how each citation was confirmed (Crossref DOI, table of contents, or the text read); citations not read first-hand are listed in doc/algorithms_core_2026-11-12/REPORT.md. The worked examples are reproduced in the functions' tests.
