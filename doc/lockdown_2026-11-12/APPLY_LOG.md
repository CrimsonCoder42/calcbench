# Step 1: how each PENDING_SHARED.md line was applied (2026-11-12)

Order of edits: taxonomy · INDEX · SCOPE · DECISIONS (and ASSUMPTIONS) · CHANGELOG · requests → TODO or SATISFIED ·
sources moved. INDEX.md is now generated from the code (D-070), so every §1 INDEX section is applied by the
regeneration; the check of each section against the generated rows is in D-070.

| File (lines) | Section | Applied as |
|---|---|---|
| finance (286) | §1 F-1…F-13 | DECISIONS D-018…D-030 (header keeps "finance F-n") |
| | §2 A-fin-1…9 | ASSUMPTIONS A-038…A-046 |
| | §3 SCOPE row | SCOPE S11 |
| | §4 taxonomy | `fin/personal` PERSON planned row; F-7 import sentence under the table; vCalc Finance row |
| | §5 CHANGELOG | entry "Finance (shipped)", now with file:line for all 122 IDs |
| | §6 sources README rows | 32 rows added to doc/sources/README.md (files already in doc/sources/) |
| | §7 INDEX (124 rows) | identical to the generated rows (124/124) |
| linalg (210) | §1 INDEX (115) | 113 identical; 2 changed by later fixes, code wins |
| | §2 SCOPE | S12 |
| | §3 taxonomy | `num/linalg` Holds gains "3D coordinate conversions (cylindrical, spherical)" |
| | §4 decisions (15 rows) | D-031…D-045 |
| | §5 sources | 13 README rows (12 new; jeannerod already listed); higham_2002.md §6.1 note appended |
| | §6 requests | sci/physics, sci/astronomy, sci/geoscience → TODO; num/geometry polar rows SATISFIED (NUM-GEOM-007/008); FMA advisory → rule sync (LEFT FOR DEVIN) |
| | §7 CHANGELOG | entry "`num/linalg` (shipped)" with file:line |
| | §8 rule sync | .claude/ is Devin's: LEFT FOR DEVIN in the report |
| geometry (151) | §1 INDEX (82) | superseded: the math build's rows (80 built, 018/019 REMOVED rows kept) |
| | §2 SCOPE | already S7 (math build) |
| | §3a taxonomy | `num/logs` Holds widened to named constants (D-011) |
| | §3b elliptic integrals | already in TODO |
| | §4 requests | already in TODO (math build); golden ratio SATISFIED (NUM-LOGS-001) |
| | §5 decisions | answered by Devin as D-011 |
| | §6 sources | created by the math build (zwillinger_crc_31.md replaces crc_32; bronshtein dropped); adlaj_2012 not used |
| prob (111) | §1 INDEX (62) | superseded by the math build's rows |
| | §2 SCOPE | already S9 |
| | §3 taxonomy | `num/special` stays a TODO question (no answer recorded) |
| | §4 requests | already in TODO; med offer SATISFIED (Bayes used by med) |
| | §5 sampling note | already in TODO |
| | §6 sources | created by the math build |
| | §7 decisions | answered by Devin as D-013 |
| medicine (108) | §1 INDEX (44) | 42 identical; 2 changed by later fixes, code wins |
| | §2 SCOPE | S13 |
| | §3 taxonomy | 9 MED rows (7 built, fluids and pharmacology planned); vCalc Medicine row |
| | §4 D-M1…D-M9, A-M1…A-M10 | D-046…D-054, A-047…A-056 |
| | §5 CHANGELOG | entry "Medicine (shipped)" with file:line |
| | §6 TODO | HELD-SOURCE and dose-family rows, unit/* note → TODO; Epworth licence and Traub–Kichen → resolved in Step 4 (CITATIONS.md) |
| | §7 sources | 27 notes moved with mv -n into doc/sources/ and indexed |
| | §8 moves | sci/physics, sci/chemistry, fin/business → TODO; num/prob SATISFIED (NUM-PROB-008/010/011) |
| algebra (74) | §1 INDEX | superseded by the math build's rows |
| | §2 SCOPE | already S5 |
| | §3 taxonomy | num/integer row and vCalc row already present |
| | §4 requests | already in TODO; vector addition SATISFIED (NUM-LINALG-001) |
| | §5 sources | created by the math build |
| | §6 decisions | answered by Devin as D-005…D-009 |
| trig (74) | §1 INDEX | superseded by the math build's rows |
| | §2 SCOPE | already S6 |
| | §3 taxonomy | `num/trig` Holds replaced; import sentence under the table; great-circle row kept as a TODO (no such vCalc row) |
| | §4 offers to geometry | taken by the build (018/019 REMOVED → NUM-TRIG-020/021) |
| | §5 requests | already in TODO |
| port (63) | §1 INDEX (212) | 209 identical; 3 changed by review fixes, code wins; 3 REMOVED rows kept with function names |
| | §2 taxonomy | 13 med rows + oncology (planned) + life/nutrition; Holds for temperature, fitness, series |
| | §3 SCOPE | S14 |
| | §4 D-P1…D-P15, A-P1 | D-055…D-069, A-057 |
| | §5 CHANGELOG | entry "Port … (shipped)" with file:line |
| | §6 requests | TODO; algo/heap TopK-by-count SATISFIED (ALG-HASH-003) |
| | §7 sources | 161 notes moved with mv -n into doc/sources/ and indexed |
| combinatorics (62) | §1 INDEX | superseded by the math build's rows (22 built, big twins included) |
| | §2 SCOPE | already S4 |
| | §3 taxonomy | `algo/enumerate` already present (S10, built) |
| | §4 offers | already in TODO |
| | §5 sources | created by the math build (brualdi not used: A-012/-013 cite OEIS) |

Totals: 9 files, 1,139 lines; source notes moved 188 (27 + 161); README rows added 232.
