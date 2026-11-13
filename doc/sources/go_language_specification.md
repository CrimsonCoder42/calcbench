# The Go Programming Language Specification

**Citations** (as written in the doc comments):
- The Go Programming Language Specification, "For statements with range clause" (range over function iterators); package iter (iter.Seq)
- The Go Programming Language Specification, "For statements with range clause" (range over function iterators); package iter (iter.Seq2)
- The Go Programming Language Specification, Go 1.23+, "For statements with range clause" (range over a function: yield returns false when the loop body stops and must not be called again), https://go.dev/ref/spec#For_range (read 2026-11-12)
- The Go Programming Language Specification, go1.27, "For statements with range clause"
- The Go Programming Language Specification, go1.27, "For statements with range clause" (range over a function; yield returns false after a break and is not called again)

**Used by**: ALG-GRAPH-007, ALG-GRAPH-009, ALG-GRAPH-010, ALG-GRAPH-022, ALG-GRAPH-023, ALG-SEQ-001, ALG-SEQ-002, ALG-SEQ-003, ALG-SEQ-004, ALG-SEQ-005, ALG-SEQ-006, ALG-SEQ-007, ALG-SEQ-008, ALG-TREE-006, ALG-TREE-007, ALG-TREE-008, ALG-TREE-009.

**Checked**: the builder's catalog section in doc/algorithms_core_2026-11-12/catalog/ states how each citation was confirmed (Crossref DOI, table of contents, or the text read); citations not read first-hand are listed in doc/algorithms_core_2026-11-12/REPORT.md. The worked examples are reproduced in the functions' tests.
