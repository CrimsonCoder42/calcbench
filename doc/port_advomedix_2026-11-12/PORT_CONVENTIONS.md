# Port conventions (binding on every builder in this run)

Read, in this order: `.claude/CLAUDE.md`; `.claude/rules/*.md`; `doc/medicine_2026-11-12/CONVENTIONS.md` (binding for
every `med/*` function: errors, units, Range, category types, tests); this file; your brief and entry list.

## 0. Identity rules (hard)
- No product or system name in calcbench code, tests, doc comments or doc.go: never "AdvoMedix", "advomedix", "Ava",
  "patient", "profile", "tenant", "user". Use "person", "subject", "the caller". The source path appears ONLY in your
  ledger section (`ledger/<wave>_<builder>.md`). Parity tests cite the ledger row, never the source path.
- The reference is READ-ONLY. Never open `.env`, `.env.*`, `app.env`, `vendor/`, `node_modules/`. Never edit there.

## 1. The unit of work, per entry (no skipping; record the outcome even when it is NOT PORTED)
1. READ the JSON, its `.go` and `_test.go` if present, and the citation it carries (label + url).
2. DUPLICATE: `grep -rn "^// ID:" /Users/devinanderson/Desktop/calcbench --include=*.go`, `INDEX.md`, every
   `doc/catalog/*.md` (medicine.md lists HELD rows too), and your own ledger. Same formula → REUSED (ledger row only).
   A HELD row in medicine.md stays HELD: write NOT PORTED with "held: doc/catalog/medicine.md §4/§5".
3. CITATION: the source's citation is a hypothesis. Find the PRIMARY (paper by DOI/PMID, guideline, official
   instrument). The formula, coefficients, item weights, cut-offs and labels must be SEEN in the primary
   (abstract, open full text, official PDF) — D-M2 of the medicine run. Seen only in a secondary → NOT PORTED,
   with what you found (citation confirmed or not, where the formula was seen). Tools that work:
   - PubMed abstracts: `https://eutils.ncbi.nlm.nih.gov/entrez/eutils/efetch.fcgi?db=pubmed&id=<PMID[,PMID…]>&rettype=abstract&retmode=text`
     (batch several PMIDs per call); search: `…/esearch.fcgi?db=pubmed&term=<terms>&retmode=json`.
   - Europe PMC full text: `https://www.ebi.ac.uk/europepmc/webservices/rest/<PMCID>/fullTextXML`; search
     `…/rest/search?query=<q>&resultType=core&format=json`.
   - Crossref metadata: `https://api.crossref.org/works/<DOI>`.
   - Most publisher pages (ahajournals, jacc, sciencedirect, wiley, oup) return 403 to fetches; do not loop on them.
   - Never vcalc, MDCalc, Wikipedia, blogs or calculator sites as the source (pointers only).
4. OPTIMIZE: state best time/space. Scalar formulas: O(1), no allocation. Note anything worse in the source (alloc per
   call, recomputed constants, string parsing, O(n²)) and write the better form; ledger "old → new".
5. WRITE in calcbench shape (rules/go-functions.md + medicine CONVENTIONS.md): one file per concept, doc comment with
   Errors line, ID, Time, Space, Source (one line each, to the table/equation/page, DOI or PMID), Tags, Units, Range,
   Example. Units in names; both unit systems where the source code converts (twins, CONVENTIONS §Units). Category
   types carry the PRIMARY source's labels verbatim (not the app's labels). No advice words.
6. TESTS (rules/go-testing.md): the primary's worked example/table; `Test<Func>_parity` reproducing EVERY case of the
   source `_test.go` (and the JSON's documented example for non-runnable entries), expected values copied as data, a
   comment `// parity: ledger <entry id>, source case N`. A case where the source disagrees with the primary is kept
   as a test asserting the PRIMARY's value, with the source's value in the comment and a ledger note — never deleted.
   Both sides of every Range and category edge; NaN/±Inf; AllocsPerRun == 0; `Benchmark<Func>` with ReportAllocs.
7. LEDGER ROW (your section file), columns exactly:
   `entry id | source path:line | calcbench ID and file:line | REUSED / PORTED / NOT PORTED / REQUESTED | citation correction | complexity old → new | parity result | consumer change`
   Write rows to the file as you go (every few entries), so nothing is lost.

## 2. Policy outcomes (decided; do not research these further)
- `doses = true` (marked D in your list): NOT PORTED, "held by policy: dose / amount to administer (D-M5)". Still note the
  citation the JSON carries.
- Pure clinical decision rules / pathways / guideline flowcharts with no computed number or published score
  (e.g. "ACEP clinical policy", "management", "indications", "workup"): NOT PORTED, "decision pathway, not a calculation".
- Coding/billing/administrative (E/M coding, VA rating): NOT PORTED, "administrative rules, not a clinical calculation".
- Instruments whose items are copyrighted and licensed (the item text): port only the arithmetic and the published
  bands; never embed item text (D-M8).

## 3. IDs and packages
- New packages (SHORT; create doc.go yourself, IDs from 001): `med/gastro` GASTRO · `med/endocrine` ENDO ·
  `med/psychiatry` PSYCH · `med/critical` CRIT · `med/hematology` HEMAT · `med/oncology` ONCO · `med/rheumatology` RHEUM ·
  `med/infectious` INFECT · `med/pediatrics` PEDS · `med/thermoregulation` THERMO · `med/emergency` EMERG ·
  `unit/temperature` TEMP · `med/neuro` NEURO · `med/trauma` TRAUMA · `med/obstetric` OB · `life/fitness` FIT ·
  `life/nutrition` NUTR.
- Existing packages (from the medicine run) are ADDITIVE ONLY: new files, IDs continue after the last; never edit
  an existing file, including its doc.go (report the doc.go lines to add; the lead places them). Next free IDs:
  ANTHRO 010 · CARDIO 004 · CHEM 011 · HEPAT 005 · RENAL 012 · RESP 002 · SCALE 007. A new file in such a package may
  not redeclare its sentinels (ErrNotFinite, ErrOutOfRange, ErrDomain, ErrNotTestable where present) or helpers:
  read the package first and reuse them.
- Only the package(s) your brief lists are yours. An entry whose home is someone else's package: REQUESTED, with the
  full research in the ledger row (citation seen, formula), for its owner.

## 4. Builders never run tests or benchmarks
Allowed: `gofmt -l <your dirs>`, `go build ./...`, `go vet ./<your packages>/...`. One single named test only if you
cannot continue without the answer (say so in your report).

## 5. Report (≤10 lines)
Counts per outcome; packages and IDs with file:line (range is fine: "MED-GASTRO-001…014 med/gastro/*.go"); tests
written; doc.go lines for the lead (existing packages); REQUESTED rows; anything unfinished and why.
