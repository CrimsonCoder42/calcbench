# Brief: calculator.net check, 2026-11-13

Every agent in this run works from this file. `.claude/CLAUDE.md` and `.claude/rules/*.md` apply to every line.

## The rule

calculator.net is a reading list, not a source. Its list tells us which calculators people expect. We research
each one from primary sources (`.claude/rules/sources.md`), write our own function from those sources, then
compare our output with theirs on the same inputs as an independent check. When the two differ, the primary
source decides. We never change our code to match their number without a citation that says so. We never copy
their text, their JavaScript or their descriptions. From their page we read two things only: the input fields
(names, units, defaults) and the result values. A formula name on their page (e.g. "Mifflin-St Jeor") is a pointer
to go and look up, never a citation. No calculator.net text or URL ever appears in a .go file.

## Conduct on their site

- One agent per index page, so at most four browser sessions in all, each behaving like one person.
- One calculator page open at a time. Wait 3-6 s before every driver run (`sleep 4 && node ...` in the same
  command). At most 8 submissions per calculator. No retries in a loop: a failed run is looked at, not repeated.
- A 403, a 429 or a block page (driver exit 3): stop all site work, record it under "Site log" in PROGRESS.md,
  carry on with research-only tasks (CLASSIFY, MATCH, RESEARCH, BUILD).
- Store only field names, the inputs we sent and the result values. Never the page text.

## The two drivers (agents never write browser code)

Run from `testing/web/`, bounded:

```
cd /Users/devinanderson/Desktop/calcbench/testing/web
sleep 4 && perl -e 'alarm 90; exec @ARGV' -- node scripts/calculator_net_fields.mjs <url>
sleep 4 && perl -e 'alarm 90; exec @ARGV' -- node scripts/calculator_net_run.mjs <url> '<JSON name->value>' [form name]
```

- `fields`: one page load. Prints JSON `{fields: [{form, name, type, label, value, unit, options}]}`.
- `run`: fills the named fields (selects, radios and checkboxes first, then text), clicks Calculate and prints
  `{submit, sent, missing, result: {text, tables}}`. Fields left out keep the page default.
  - `sent` is what the form held at submit time. Always check it against what you meant to send: a page script may
    reformat or convert a value.
  - Checkbox values are `true`/`false`. A radio takes its option value.
  - Pass the form name (e.g. `calform`) when the page has more than one form.
- Proven on the Mortgage Calculator: $240,000 at 6 %/yr over 30 years gives $1,438.92/month, the same as a hand
  computation (`out/run_mortgage.json`).
- Driver update (lead, after A-OTHER session 1):
  - A field may be named by its id when it has no name.
  - A page with no Calculate button runs in live mode (`submit: "live"`): typing events are fired, then the page
    is read.
  - `fieldsAfter` lists every form text box and select after the submit. Solvers such as Ohm's law write their
    answers there.
  - Date pickers: set the hidden field the form submits, in MM/DD/YYYY. For example the Age Calculator takes
    `{"today":"05/15/1990","ageat":"11/13/2026"}` and returns 13,331 days (`out/run_age_driverfix.json`). The
    month, day and year widgets beside it are display only.
- Driver update 2 (lead, after A-FIN and A-MATH session 1):
  - Both drivers now wait for the full page load plus 1.5 s. Pages that showed "only the site search form" draw
    their form late; Auto Loan now shows 16 fields.
  - Fields are found anywhere on the page, in or out of a `<form>`.
  - The Calculate button is chosen in priority order, never the site search button.
  - An optional 4th argument picks a button by its label, e.g. `... '' Multiply` for one matrix operation (pass
    `''` as the form name for the whole page).
- Shared scratchpad: give every helper script a name with your agent prefix (`scratchpad/afin/...`,
  `scratchpad/aother/...`). A helper with a shared name was overwritten by another agent and wrote into the
  wrong PROGRESS.md.
- Exit codes: 3 = blocked, 4 = no button or no result block. A driver gap (a result table not captured, a widget it
  cannot fill) goes into PROGRESS.md as `DRIVER: <what>` for the lead. Do not work around it in your own code.

## Agents and owned paths (code you may create or edit; read anything)

| Agent | Index | Folder | Owned code paths |
|---|---|---|---|
| A-FIN | financial | `financial/` | `fin/**` |
| A-FIT | fitness-and-health | `fitness-and-health/` | `med/**`, `life/**` (except `life/everyday`) |
| A-MATH | math | `math/` | `num/**`, `algo/**` |
| A-OTHER | other | `other/` | `unit/**`, `sci/**`, `eng/**`, `life/everyday` (new, only if needed) |

Each agent also owns its own folder here: `PROGRESS.md`, `DISCREPANCIES.md`, `checks/<slug>.md`,
`checks/scratch/<slug>.go`, `LEAD.md`. The lead merges each `DISCREPANCIES.md` into the root `DISCREPANCIES.md`.
A calculator whose home is another agent's path:
- write the research into your `checks/<slug>.md`;
- mark the row `HANDOFF <agent>`;
- the lead adds it to that agent's next relaunch list.

**The LEAD owns**: `INDEX.md` (regenerated from code, D-070), `doc/taxonomy.md`, `_delete/LEDGER.md`, every
shared ledger, and the header of each new package's `doc.go`. Agents never edit those files.
- **A new package that `doc/taxonomy.md` already plans** (e.g. `fin/personal` PERSON, `sci/meteorology` METEO,
  `eng/electrical` ELEC, `life/automotive` AUTO, `sci/astronomy` ASTRO for Julian-day calendar maths): the agent
  creates it with a `doc.go` that copies the shape of an existing one (e.g. `fin/loan/doc.go`) and uses the SHORT
  from the taxonomy row. The lead checks the header at merge.
- **A package not in the taxonomy** (`life/everyday` included): write the proposed path, SHORT and one-line
  "Holds" in `LEAD.md` and build into it only once the lead has added the row (the lead answers at relaunch).
  Until then the row stays open as `WAITING SHORT`.
- **`LEAD.md` also carries**: INDEX/taxonomy notes, `_delete` moves made (path | replaced by | date) for the ledger,
  driver gaps, and handoffs.

## Next free ID per prefix (INDEX.md 2026-11-13, REMOVED rows counted; new topics start at 001)

FIN-BIZ 051 · FIN-ECON 005 · FIN-INT 043 · FIN-INVEST 053 · FIN-LOAN 011 · FIN-MKT 018 ·
LIF-FIT 015 · LIF-NUTR 008 ·
MED-ANTHRO 016 · MED-CARDIO 069 · MED-CHEM 014 · MED-ENDO 032 · MED-OB 003 · MED-PEDS 009 · MED-RENAL 028 ·
MED-SCALE 016 · MED-THERMO 005 (other MED topics: highest in INDEX.md plus one) ·
NUM-ALGEB 029 · NUM-CALC 005 · NUM-COMBI 023 · NUM-FLT 007 · NUM-GEOM 083 · NUM-INTEG 007 · NUM-LINALG 192 ·
NUM-LOGS 005 · NUM-PROB 156 · NUM-SERIES 020 · NUM-SIG 017 · NUM-SPEC 021 · NUM-STATS 047 · NUM-TRIG 033 ·
ALG-INTMTH 017 (other ALG topics: highest plus one) ·
UNT-LEN 003 · UNT-MASS 003 · UNT-TEMP 005 · new UNT/SCI/ENG/LIF topics from 001.
Only the owning agent numbers a topic, so two agents never take the same ID.

## Baseline (2026-11-13, before this run)

`go build ./...` OK · `go vet ./...` clean (`out/vet_baseline.txt`, empty) · `gofmt -l .` clean ·
INDEX.md: 1,564 live functions + 29 REMOVED rows.

## Per calculator, in order

1. **CLASSIFY.** OUT OF SCOPE = tax brackets and benefit tables, currency rates, time-zone databases, random or
   password generators, entertainment, encoders, anything that needs a data table that changes by year or
   jurisdiction. The stable formula inside such a calculator is still in scope, e.g. a progressive tax on a
   supplied bracket slice. Write the one-line reason, mark `SKIPPED (reason)`, next row.
2. **FIELDS.** Run `fields` once. Save the fields and defaults to `checks/<slug>.md`.
3. **MATCH.** grep `INDEX.md` and the package by the mathematics. Record `HAVE <IDs>`, `PARTIAL <IDs>` (what their
   inputs cover that ours do not) or `NONE`.
4. **RESEARCH** (PARTIAL and NONE).
   - Primary sources only, cited to the chapter or equation: textbooks, standards (ASHRAE, ACI, NIST, FDA, WHO,
     ACSM, USGA...), journal papers, regulations (Reg Z, Reg DD).
   - OpenStax finance and economics texts count (D-018).
   - No primary source after a real search: `NEEDS SOURCE` with what you looked for, and nothing is built.
5. **BUILD** (PARTIAL and NONE with a source).
   - Our own function from the source, in the owned package.
   - The doc-comment template complete (`.claude/rules/go-functions.md`): ID, Time, Space, Source, Tags, Units,
     Example; Range for medical.
   - Money is int64 cents, rates are per stated period, and NaN/Inf in returns an error (`.claude/rules/numeric.md`).
   - `_test.go`:
     - oracle values from the source's worked example or a hand computation (written in the test comment);
     - edge and error cases;
     - a benchmark with `b.ReportAllocs()`.
   - Add the function to the package `doc.go` list.
   - `gofmt -l`, `go build ./...` and `go vet ./<pkg>/` on the package. **Never `go test`.**
   - Existing exported signatures never change. A different contract is a new function with a new ID; the old one
     moves with `mv -n` to `_delete/<same path>` only when it is replaced, logged in `LEAD.md`.
6. **COMPARE.** Pick 3-5 input sets: the source's worked example first, then one ordinary case, then one edge their
   form accepts.
   - Compute ours with `checks/scratch/<slug>.go`, `//go:build ignore` on line 1, run from the repo root with
     `perl -e 'alarm 60; exec @ARGV' -- go run doc/calculator_net_check_2026-11-13/<folder>/checks/scratch/<slug>.go`.
   - Submit the same inputs through `run`.
   - Table in `checks/<slug>.md`: inputs | ours | theirs | diff | verdict.
   - Match = equal at their displayed precision, once the obvious conventions are aligned and written down: cents,
     day count, payment timing, rounding mode, unit.
7. **DISCREPANCY.** Ours differs beyond their displayed precision: go back to the source. Record the citation in
   `checks/<slug>.md` and in your `DISCREPANCIES.md`, with one of:
   - `OURS WRONG`: fix the code and the test oracle, and update Source if a better source was found.
   - `THEIRS DIFFERS`: ours stays. Write the convention or formula they appear to use, and the source backing ours.
   - `UNRESOLVED`: ours stays. Record both numbers and the two candidate sources for Devin.

   Never align to their number by changing a constant without a source.
8. **CLOSE** the row: `MATCH` / `FIXED (what, source)` / `BUILT <IDs>` / `SKIPPED (reason)` / `NEEDS SOURCE` /
   `HANDOFF <agent>` / `BLOCKED` / `UNRESOLVED`.
   - A row is closed as soon as its outcome is known.
   - Update PROGRESS.md after every row, so a relaunch loses nothing.

## Limits

- **Session cap.** Each agent session is capped at 45 minutes of wall clock: note `date` at the start and stop at
  45 min after closing the row in hand. Hand back a summary of PROGRESS.md (rows closed this session by outcome,
  rows still open, LEAD.md items).
- **Commands.** Every command is bounded with `perl -e 'alarm N; exec @ARGV' -- <cmd>` (macOS has no `timeout`).
  No background processes. Stop any PID you start.
- **Files.** Never `rm`; deletes and replacements are one `mv -n` each into `_delete/<same path>`. Git is read-only.
  Never open `.env` files, anything under advo_production or red_lite, or `testing/web/node_modules/`.
  `.claude/**` and `go.mod` are protected.
- **Other sessions.** Other Claude sessions may be live in this tree. Touch only your owned paths. A red file
  elsewhere goes under FOUND, NOT FIXED in `LEAD.md`.
