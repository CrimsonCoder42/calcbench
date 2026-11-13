# calcbench

A Go library of calculators and algorithms shared by several projects. Module
`github.com/CrimsonCoder42/calcbench`, Go 1.27, standard library only. No server, no database, no network, no
configuration, no logging: pure functions with typed inputs, named units, a cited primary source per formula,
and tests that reproduce the source's worked example.

## The groups

Folders are by subject (what the maths is), never by audience. One function has one home.

| Group | Folder | What lives there | ID prefix |
|---|---|---|---|
| Algorithms | `algo/<topic>/` | data structures and algorithm patterns with no domain words | `ALG` |
| Numbers | `num/<topic>/` | pure mathematics: algebra, calculus, combinatorics, float, geometry, integer, linalg, logs, prob, series, signal, special, stats, trig | `NUM` |
| Units | `unit/<quantity>/` | unit conversions, one package per physical quantity | `UNT` |
| Science | `sci/<field>/` | physics, chemistry, astronomy and the other sciences (not started) | `SCI` |
| Engineering | `eng/<field>/` | engineering formulas (not started) | `ENG` |
| Money | `fin/<field>/` | finance, business, economics | `FIN` |
| Medical | `med/<specialty>/` | clinical calculators and reference ranges, each with its primary source | `MED` |
| Everyday | `life/<field>/` | fitness, nutrition and other everyday calculators | `LIF` |

Imports run top to bottom: `algo` imports nothing here; `num` may import `algo`; `unit` may import `num`; the
domain groups may import `algo`, `num` and `unit`, never each other. The full map of packages, built and planned,
is [doc/taxonomy.md](doc/taxonomy.md).

## Finding a function

- [INDEX.md](INDEX.md) has one row per exported function, sorted by ID:
  `ID | Function | Package | Does | Time | Space | Source | Tags`. Search it by tag or by words.
- IDs are `<PREFIX>-<SHORT>-<NNN>`, for example `ALG-HEAP-002`, `FIN-LOAN-001`, `MED-RENAL-012`. An ID is never
  reused; a removed function keeps a REMOVED row.
- To find the code: `grep -rn "ID:      MED-RENAL-012" .`
- Each package's `doc.go` lists its functions with their IDs.

## Using it

```go
import "github.com/CrimsonCoder42/calcbench/med/renal"

stage, err := renal.AKINStage(1.0, 2.1, 1.0, 24, false) // creatinine mg/dL, urine mL/kg/h, hours
```

Every function returns an error for invalid input (NaN, ±Inf, out of the source's range) instead of panicking;
the sentinel errors are exported per package.

## The rules, in one paragraph

Every function is pure, safe on any input, documented (ID, Time, Space, Source to the line, Tags; calculators
also Units and Example; medical functions also Range), tested (the source's worked example, edge cases, an
oracle where a slow version exists, a benchmark with allocations), and written at the best known time and space.
Formulas come from primary sources only; catalog and problem sites are never copied. Medical results are numbers
and the source's own category labels, never advice. The full rules are in
[.claude/rules/](.claude/rules/): [go-functions.md](.claude/rules/go-functions.md),
[go-testing.md](.claude/rules/go-testing.md), [go-complexity.md](.claude/rules/go-complexity.md),
[numeric.md](.claude/rules/numeric.md), [sources.md](.claude/rules/sources.md).
[testing/index_test.go](testing/index_test.go) fails when the code, the doc.go lists, INDEX.md and the taxonomy
disagree.

## Status (2026-11-13)

| Group | Packages | Functions (IDs) |
|---|---|---|
| ALG | 30 | 375 |
| NUM | 14 | 627 |
| UNT | 3 | 8 |
| FIN | 6 | 160 |
| MED | 21 | 376 |
| LIF | 2 | 18 |
| **Total** | **76** | **1,564** |

`go test ./... -count=1`: 77 packages ok, 0 failures (the 76 above plus `testing`; S19 review: 75 ok in the full run, the 2 fixed packages ok on rerun). Slices and what is not started: [doc/scope/SCOPE.md](doc/scope/SCOPE.md). Docs index:
[doc/README.md](doc/README.md).
