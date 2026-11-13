# Wave 2 conventions (binding on every builder; read after the wave 1 files named below)

Read, in order: `.claude/CLAUDE.md`; `.claude/rules/*.md`; `doc/medicine_2026-11-12/CONVENTIONS.md` (every `med/*`
function: errors, units, Range, category types, tests); `doc/port_advomedix_2026-11-12/PORT_CONVENTIONS.md` (the unit
of work, identity rules, tools that work, policy outcomes); THIS file, which overrides wave 1 where they differ; your
brief and entry list.

## 1. Every entry ends in exactly ONE final state (no bare "not ported")
- **PORTED**: built; formula, constants, categories and validated population verified in the primary, read in full or
  in the part that states them (abstract only when it states the complete formula with every constant).
- **PORTED-SECONDARY**: the primary could not be opened, but TWO independent authoritative sources state the same
  formula, constants and categories AND name the primary. Authoritative = a guideline body's own document, a standard
  textbook, an NIH / CDC / WHO page, or a peer-reviewed article in PMC that reproduces the instrument. Never a calculator
  site, Wikipedia, a blog or the reference itself. Built, with both sources as `Source:` lines after the primary, a
  `Verified: secondary (<source 1>; <source 2>)` line in the doc comment (after Source, before Tags), and a row in
  your section of CONFIRM (see §5).
- **REUSED**: the same calculation already has a calcbench ID (grep `^// ID:` and INDEX.md).
- **NEEDS PDF**: neither route worked. The row carries the exact citation: authors, journal, year, volume, pages, DOI,
  PMID, so Devin can drop `<PMID>.pdf` into `doc/port_advomedix_wave2_2026-11-12/papers/`. (That folder is empty at the
  start of this run.)
- **HELD BY POLICY**: outputs a drug dose or an amount to administer, or reproduces a licensed questionnaire's text.
  Record the label or licence. For a licensed questionnaire whose scoring arithmetic is free to use, PORT the
  arithmetic only (state "items not embedded; licence: <holder>" in the doc) and the row is PORTED.
- **NOT A CALCULATION**: an identity, an administrative or billing rule, or a narrative pathway with NO published rule
  table. Final; it stays with the consumer.
- **SOURCE DOES NOT SUPPORT**: the entry's formula, constants or categories match no source that can be found, or
  contradict the paper cited. Not ported. Your row says what the entry computes, what the paper says, and the locator
  (it goes to CONSUMER_DEFECTS.md: the consumer shows an unsupported number).

## 2. How to reach a paper, in this order; record the route that worked
1. PubMed Central / Europe PMC full text (`https://www.ebi.ac.uk/europepmc/webservices/rest/<PMCID>/fullTextXML`;
   find the PMCID with `https://www.ebi.ac.uk/europepmc/webservices/rest/search?query=EXT_ID:<PMID>%20AND%20SRC:MED&resultType=core&format=json`).
2. The publisher's open-access page. 3. An author manuscript found through the DOI. 4. The guideline body's own PDF.
5. The papers/ folder. 6. The abstract, only when it states the complete formula with every constant
   (`https://eutils.ncbi.nlm.nih.gov/entrez/eutils/efetch.fcgi?db=pubmed&id=<PMIDs>&rettype=abstract&retmode=text`).
Then the PORTED-SECONDARY route. Do not retry a refused site with another tool. You fetch pages with the `web-fetch`
agent type (Agent tool) or WebFetch if you have it: ask it for VERBATIM table cells / sentences and the URL that
worked; large pages are summarised past ~39k characters, so ask for the specific table or section only. A PDF it saves
can be read with `pdftotext -layout <file> -` in Bash.

## 3. What each kind of leftover needs
- Formulas and scores whose paper was not opened: §2, then the unit of work (PORT_CONVENTIONS §1).
- Lookup tables and percentile charts (growth, spirometry reference values, blood-pressure percentiles): only from the
  official data file of the publishing body (CDC, WHO, NHLBI, GLI). Save the file under the package's `data/` folder,
  with a `data/README.md` giving its URL, the date read and its SHA-256 (`shasum -a 256`); load it with `go:embed`; test
  at least five published values per table. No official file: NEEDS PDF.
- Decision rules: PORT when the source publishes an explicit rule or point table (criteria counts such as "2 major or
  1 major + 2 minor" are explicit rules); the function returns the source's own category label and nothing else.
  Otherwise NOT A CALCULATION.
- Wrong or mismatched citations: find the instrument the entry really computes and cite that; if none fits, SOURCE
  DOES NOT SUPPORT.

## 4. Code shape (unchanged from wave 1, restated)
- Package per doc/taxonomy.md; only the packages in your brief are yours. Nothing ported in wave 1 is rewritten: never
  edit an existing file. New files only; IDs continue from the "next free" numbers in your brief.
- An existing package's doc.go is NOT edited by you: put the lines to add in your report. A NEW package (only those your
  brief names) gets a doc.go you write (package comment, "ID short form: X (IDs are MED-X-NNN).", the list).
- Doc comment: line 1 names the function; an Errors: line; ID, Time, Space, Source (one line each, to the equation /
  table / page, with DOI or PMID), [Verified: secondary (...)], Tags, Units, Range, Example. Both unit systems where
  clinicians use both. Errors, never panics. Numbers and the source's category names only; never advice words
  ("normal", "should", "dangerous", "recommend"); never "patient", "profile", "tenant", "user" or any product name.
- Tests (written, never run by you): the paper's worked example or table values; `Test<Func>_parity` reproducing the
  source entry's `_test.go` cases (or its JSON example) as data, `// parity: ledger <entry id>, source case N`; a
  disagreement is a finding (assert the PRIMARY, note the source value, ledger says which side the primary supports);
  both sides of every category edge; out-of-range refusal; NaN/±Inf; AllocsPerRun == 0 for scalars;
  `Benchmark<Func>` with `b.ReportAllocs()`. Every oracle or loop capped; no test over 60 s.
- Allowed commands: `gofmt -l <your dirs>`, `go build ./...`, `go vet ./<your packages>/...`. No `go test`. Every
  command bounded; nothing in the background. Never `rm`; a file you want gone is `mv -n` to `_delete/<same path>`.

## 5. What you write (your owned doc paths; no shared files)
- `doc/port_advomedix_wave2_2026-11-12/ledger/<wave>_<builder>.md`: one row per entry, columns exactly
  `entry id | source path:line | calcbench ID and file:line | FINAL STATE | route that worked | citation (exact: authors, journal, year;vol:pages, DOI, PMID, locator) | parity | consumer change`.
  Append as you go. For SOURCE DOES NOT SUPPORT the citation cell also says "entry computes …; paper says …".
- `doc/port_advomedix_wave2_2026-11-12/sources/<entry id>.md`: one note per PORTED / PORTED-SECONDARY entry, 3–6
  lines: the citation, the locator, the values transcribed (every constant, cut-off and label), the route, the date.
  The reviewer checks code against it without the network.
- At the bottom of your ledger file, three short lists the lead copies: `CONFIRM` (PORTED-SECONDARY rows: entry, both
  sources), `NEEDS PDF` (entry, full citation, PMID), `CONSUMER DEFECTS` (entry, what it computes, what the paper says,
  locator, severity high/medium/low).
- Report to the lead in at most 10 lines: counts per final state; IDs with file ranges; doc.go lines for existing
  packages; anything another builder's package needs; anything unfinished and why.
