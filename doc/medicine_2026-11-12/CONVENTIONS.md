# Conventions for every `med/*` package (medicine build, 2026-11-12)

Binding on every builder in this run. Decided by the lead under the request's "no stops" instruction; each
decision is in DECISIONS_DRAFT.md (for doc/DECISIONS.md via PENDING_SHARED.md §4) with the alternative that lost.

## Shape
- One package per specialty (`med/anthropometry`, `med/cardio`, `med/respiratory`, `med/renal`, `med/chemistry`,
  `med/fluids`, `med/hepatic`, `med/scales`, `med/pharmacology`). `med/*` packages do not import each other and
  import nothing from calcbench (standard library only: `errors`, `fmt`, `math`).
- `doc.go`: package comment, `ID short form: <SHORT> (IDs are MED-<SHORT>-NNN).`, the units policy for the
  package, the statement "Results are numbers and the source's own category labels; nothing here diagnoses,
  interprets or advises.", and the function list `MED-<SHORT>-NNN  Name  one line`.
- One file per concept (snake_case), its `_test.go` beside it. Sentinels live in the file of the first function
  (ASSUMPTIONS A-008), in a `var (...)` block, each with a comment:
  - `ErrNotFinite`  = errors.New("<pkg>: value is NaN or infinite")
  - `ErrOutOfRange` = errors.New("<pkg>: input outside the range the source validated")
  - `ErrDomain`     = errors.New("<pkg>: inputs valid one by one but the formula is undefined for them together")
    (only where needed, e.g. a zero denominator from two in-range inputs)
  Errors are wrapped with the parameter: `fmt.Errorf("%w: ageYears = %v, want 18 ≤ age ≤ 90", ErrOutOfRange, ageYears)`.
  A NaN or ±Inf input returns ErrNotFinite before any range check.

## Doc comment (`.claude/rules/go-functions.md`), medical form
```
// EGFRCKDEPI2021 returns the estimated glomerular filtration rate, in mL/min/1.73 m², from serum creatinine
// by the 2021 CKD-EPI creatinine equation (race-free).
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range.
//
// ID:      MED-RENAL-001
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Inker LA, Eneanya ND, Coresh J, et al. New creatinine- and cystatin C-based equations to estimate GFR without race. N Engl J Med. 2021;385(19):1737-1749. doi:10.1056/NEJMoa2102953, <locator>
// Tags:    renal, egfr, ckd-epi, creatinine
// Units:   creatinineMgDL mg/dL; ageYears years (decimal, completed age as the caller has it); result mL/min/1.73 m²
// Range:   ageYears 18–… ; creatinineMgDL … (the population the source validated)
// Example: … (the source's own numbers, or "hand-derived from eq. X" with the arithmetic)
```
- Line 1: what it returns, with the unit. No word that judges a value ("normal", "abnormal", "dangerous",
  "healthy", "should", "risk of"): only the source's category labels, quoted, with the citation.
- `Source:` cited to the equation, table or page; a DOI for every article; one `Source:` line per source, each on ONE line (the index is built from it; see num/trig/coth.go).
- Errors in the doc comment: a line `// Errors: ...` before the blank line, as num/trig/coth.go does.
- `Range:` is the input range the source validated (ages, weights, heights, lab ranges, population). Outside
  it: ErrOutOfRange. Where the paper reports a population but not hard limits, the Range states the reported
  limits (e.g. "ages 18–92, the study's range") and the function enforces those numbers.

## Units (numeric policy, decided by Devin)
- Every input and result says its unit in its name: `weightKg`, `heightCm`, `creatinineMgDL`,
  `creatinineUmolL`, `glucoseMgDL`, `glucoseMmolL`, result names in the function name or doc line 1.
- Where clinicians use both systems, one function per system: `XxxMgDL` / `XxxSI` (or `XxxImperial`). The SI
  function applies the documented conversion factor ONCE at entry, as a named unexported constant with its
  source on the line above, then calls the same core; never a hidden conversion. Tests assert the two agree to
  the source's precision.
- Conversion constants (until `unit/*` exists; FOUND, NOT FIXED for Devin): `cmPerInch = 2.54` (exact, by the
  1959 international yard and pound; NIST SP 811 (2008) §B.8), `kgPerPound = 0.45359237` (exact, same),
  creatinine 88.4 µmol/L per mg/dL (KDIGO 2012 CKD guideline, Kidney Int Suppl 3:1, "Notes on units"), glucose
  18.016 mg/dL per mmol/L only where the source states its own factor (otherwise use the source's).
- Ages: `ageYears float64`, the caller's completed age in years (decimal allowed); children's formulas that
  are stated in months take `ageMonths float64`.
- Sex: `female bool` where the source's coefficient depends on sex (the source's own binary).
- Rounding: never, unless the source rounds; then to the source's precision, documented.

## Categories
- A category is a named string type with constants whose values are the source's labels VERBATIM, e.g.
  `type EpworthCategory string; const EpworthLowerNormal EpworthCategory = "Lower Normal Daytime Sleepiness"`.
  The function returns the constant; its doc comment cites the table. Edges are tested on both sides.

## Scores
- Item inputs are `int` with their allowed range checked (ErrOutOfRange outside), returned total is `int`.

## Tests (`.claude/rules/go-testing.md`)
1. The source's worked example or table values, with the citation in a comment. If the source has none,
   a hand-derived value with the arithmetic in the comment (labelled HAND-DERIVED).
2. A boundary value on each side of every Range edge and every category edge.
3. NaN, +Inf, −Inf → ErrNotFinite; out of range → ErrOutOfRange (errors.Is).
4. Unit-system twins agree within a stated tolerance (the source's precision).
5. `BenchmarkXxx` with `b.ReportAllocs()`, one scalar case; and a `TestXxx_noAlloc` using
   `testing.AllocsPerRun` asserting 0.
Floats compared with a tolerance stated in a comment (e.g. 1e-9 relative for exact arithmetic, the source's
printed precision for its tables).

## Builders
Write code and tests; run only `gofmt -l <pkg>`, `go build ./med/...`, `go vet ./med/<pkg>/`. Never run tests.
Touch only your own package folder(s). Report every citation you could not confirm.
