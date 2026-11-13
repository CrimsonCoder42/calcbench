# calculator.net map, 2026-11-13

A map of every calculator listed on calculator.net's four index pages (index page, section, name, full URL),
so later runs can check calcbench's numbers against theirs one calculator at a time. This run read only the
four index pages. It opened no calculator page, submitted nothing, read no JavaScript and copied no
descriptive text. The "calcbench match" column was filled from `INDEX.md` without visiting the site; it is a
first pass for Devin to read, not a verdict.

## Files

| File | What |
|---|---|
| `MAP.md` | the table: # · index page · section · calculator name · URL · calcbench match |
| `map.json` | the same rows as `[{index, section, name, url, match}]` |
| `robots.txt` | calculator.net's robots.txt, verbatim, fetched first |
| `out/` (gitignored) | the script's raw output and the install logs |

Script: `testing/web/scripts/calculator_net_map.mjs`, run once from `testing/web/` with Playwright 1.65.0
(headless Chromium 1254, default user agent). It was not changed after the run.

## Conduct on their site

- robots.txt loaded first and saved verbatim. It has `User-agent: *` with no `Disallow` lines (plus a
  sitemap line), so no index path was skipped.
- One headless browser, one page at a time, no parallel loads, no retries. A random 4-8 s wait between page
  loads. Five page loads in total: robots.txt plus the four indexes, all HTTP 200, no block page.
- Only anchor text and href were collected, from the main content column. A link's section is the nearest
  heading above it ("Math Calculators" is recorded as "Math").
- Dropped: the two breadcrumb links above the lists, and the site-wide sidebar. The sidebar starts at the
  first link back to an index page, and it repeats calculators already listed.

## Counts (fetched 2026-11-13)

| Index page | Rows | Expected | HAVE | PARTIAL | NONE | OUT OF SCOPE |
|---|---|---|---|---|---|---|
| financial-calculator.html | 71 | ~75 | 22 | 25 | 10 | 14 |
| fitness-and-health-calculator.html | 28 | ~29 | 5 | 12 | 10 | 1 |
| math-calculator.html | 38 | ~38 | 18 | 5 | 13 | 2 |
| other-calculator.html | 54 | ~56 | 0 | 8 | 38 | 8 |
| **total** | **191** | ~198 | 45 | 50 | 71 | 25 |

The small shortfalls against the expected counts are what the pages list today. All 191 URLs are unique,
the lists are static HTML (nothing is lazy-loaded), and no list was cut short.

## Match key

`HAVE <IDs>`: calcbench computes what the page computes. `PARTIAL <IDs> (what differs)`: the core formula
exists but an input, method or variant on their page does not. `NONE`: nothing in calcbench yet.
`OUT OF SCOPE`: tax tables, jurisdiction rates, random generators, lookup-table converters, entertainment.
