# Taxonomy: where everything lives

The map of calcbench. Folders are by SUBJECT, one home per function. Proposed 2026-11-12; Devin approves
changes. A package is created when its first function lands; until then its row here reserves the name.
ID = `<GROUP>-<SHORT>-<NNN>`; SHORT is fixed here and in the package's `doc.go`.

## Groups and packages

One row per package (D-072). Status: **built** (the folder holds code) or **planned** (the name and SHORT are reserved;
the folder is created with its first function). testing/index_test.go checks this table against the tree.

| Group | Package (import path under calcbench/) | SHORT | Status | Holds |
|---|---|---|---|---|
| ALG | `algo/bits` | BITS | built | bit manipulation |
| ALG | `algo/backoff` | BACKOF | built | retry backoff: capped exponential, full, equal and decorrelated jitter; the caller passes the random draw (S20, D-087) |
| ALG | `algo/bloom` | BLOOM | built | Bloom filter with enhanced double hashing, false-positive rate, optimal parameters (S10) |
| ALG | `algo/checksum` | CHKSUM | planned | check digits and checksums (D-016) |
| ALG | `algo/cron` | CRON | built | POSIX five-field cron specs: parse, next fire time with stated DST rules (S20, D-087) |
| ALG | `algo/dp` | DP | built | dynamic-programming patterns |
| ALG | `algo/enumerate` | ENUM | built | permutations, combinations, subsets, parentheses, Gray codes as iterators (S10; proposed by the combinatorics pass) |
| ALG | `algo/geometry2d` | GEO2D | built | computational geometry on exact int64 points: hull, closest pair, point in polygon, intersections (S10; measurement formulas are num/geometry) |
| ALG | `algo/graph` | GRAPH | built | traversal, shortest path, topological order |
| ALG | `algo/hashing` | HASH | built | hash-table patterns: frequency, top-k by count, grouping, set operations (S10, D-017: not `hash`) |
| ALG | `algo/heap` | HEAP | built | binary heap, top-k |
| ALG | `algo/interval` | INTVL | built | merging, intersecting and covering ranges |
| ALG | `algo/intmath` | INTMTH | built | integer arithmetic algorithms: power by squaring, integer sqrt, digits, Hamming numbers, decimal-string arithmetic (S10; gcd, modular arithmetic and primes stay in num/integer, D-005) |
| ALG | `algo/list` | LIST | built | linked lists, LRU and LFU caches |
| ALG | `algo/matrix` | MATRIX | built | grid algorithms over [][]T: rotation, transpose, sparse multiply, Game of Life (S10; float64 linear algebra is num/linalg) |
| ALG | `algo/merge` | MERGE | built | merging sorted sequences |
| ALG | `algo/prefix` | PREFIX | built | prefix sums and difference arrays |
| ALG | `algo/queue` | QUEUE | built | queue from two stacks, circular buffer (S10; the monotonic queue is algo/window's) |
| ALG | `algo/ratelimit` | RATELM | built | token bucket, GCRA leaky bucket, fixed and sliding-log windows; the caller passes the clock (S10) |
| ALG | `algo/reduce` | RED | built | slice reductions every package shares: Neumaier-compensated Sum and its step, Prod, CumSum, CumProd, Diff, ArgMin, ArgMax, AllFinite (D-081) |
| ALG | `algo/sample` | SAMPLE | built | shuffling, reservoir and weighted sampling with a caller-supplied random source (S10) |
| ALG | `algo/search` | SEARCH | built | binary search and its bounds on sorted data |
| ALG | `algo/selection` | SELECT | built | k-th element, median without a full sort |
| ALG | `algo/seq` | SEQ | built | range-over-func helpers over iter.Seq (S10, D-017: not `iter`) |
| ALG | `algo/sorting` | SORT | built | stable, counting, radix, bucket and partial sorts; argsort; three-way partition (S10, D-017: not `sort`) |
| ALG | `algo/stack` | STACK | built | monotonic stacks, min stack, RPN, bracket matching (S10) |
| ALG | `algo/stream` | STREAM | built | streaming aggregates the caller drives: moving average, hit counter, disjoint intervals, running statistics (S10) |
| ALG | `algo/text` | TEXT | built | string matching and distance |
| ALG | `algo/tree` | TREE | built | binary search trees, traversals, Fenwick and segment trees |
| ALG | `algo/trie` | TRIE | built | prefix trees |
| ALG | `algo/twopointer` | TWOPTR | built | two-pointer scans over sorted or in-place slices (S10) |
| ALG | `algo/unionfind` | UNION | built | disjoint sets |
| ALG | `algo/window` | WINDOW | built | sliding-window sums, means, minima and maxima |
| NUM | `num/algebra` | ALGEB | built | roots, polynomials, basic operations, boolean |
| NUM | `num/calculus` | CALC | built | numeric differentiation and integration on tabulated data: trapezoid, gradient, table interpolation (S18); differential equations not started |
| NUM | `num/combinatorics` | COMBI | built | permutations, combinations |
| NUM | `num/float` | FLT | built | floating-point helpers Go's math lacks: isclose/allclose (symmetric), spacing (ulp), nan_to_num, machine constants (D-081) |
| NUM | `num/geometry` | GEOM | built | lengths, areas, volumes, angles |
| NUM | `num/integer` | INTEG | built | primes, factorisation, gcd, lcm, modular arithmetic (D-005) |
| NUM | `num/linalg` | LINALG | built | vectors, matrices, quaternions, 3D coordinate conversions (cylindrical, spherical); the matrix toolkit: echelon forms and subspaces, cofactors, checks, Cholesky, QR, least squares, symmetric and 2×2 eigenproblems, SVD, norms (D-078) |
| NUM | `num/logs` | LOGS | built | logarithms, exponentials, and named mathematical constants (π-family, e, φ the golden ratio, D-011) |
| NUM | `num/prob` | PROB | built | probability, odds, expected value |
| NUM | `num/series` | SERIES | built | sequences and series, trend and smoothing of a number series, time windows and buckets (D-062) |
| NUM | `num/signal` | SIG | built | discrete Fourier transform (radix-2, Bluestein), frequency bins and shifts, window functions, convolution and correlation (D-081) |
| NUM | `num/special` | SPEC | built | special functions Go's math lacks: regularized incomplete gamma and beta, beta, Bessel I0, sinc, orthogonal polynomials (Legendre, Chebyshev, Hermite, Laguerre) and their roots (D-081) |
| NUM | `num/stats` | STATS | built | descriptive statistics, percentiles, regression, distributions |
| NUM | `num/trig` | TRIG | built | reciprocal and inverse functions math lacks, degree/radian, angle normalisation, triangle solving (plane, right), central angle; never wraps math (D-012) |
| UNT | `unit/area` | AREA | built | square foot and square meter (NIST SP 811 B.8; 2026-11-13) |
| UNT | `unit/datasize` | DATA | planned | conversions for one quantity |
| UNT | `unit/duration` | DUR | planned | conversions for one quantity |
| UNT | `unit/energy` | ENERGY | planned | conversions for one quantity |
| UNT | `unit/length` | LEN | built | conversions for one quantity |
| UNT | `unit/mass` | MASS | built | conversions for one quantity |
| UNT | `unit/pressure` | PRESS | planned | conversions for one quantity |
| UNT | `unit/rf` | RF | planned | conversions for one quantity |
| UNT | `unit/speed` | SPEED | planned | conversions for one quantity |
| UNT | `unit/temperature` | TEMP | built | °C, °F and K conversions (NIST SP 811) |
| UNT | `unit/volume` | VOL | planned | conversions for one quantity |
| SCI | `sci/astronomy` | ASTRO | built | Julian Day Number calendar arithmetic: day of week, day of year, days between, add days (Meeus ch. 7; 2026-11-13) |
| SCI | `sci/biology` | BIO | planned | the science formulas |
| SCI | `sci/chemistry` | CHEM | built | amount of substance from mass, molarity (2026-11-13) |
| SCI | `sci/ecology` | ECO | planned | the science formulas |
| SCI | `sci/geoscience` | GEO | built | great-circle and Vincenty WGS 84 surface distance (2026-11-13) |
| SCI | `sci/meteorology` | METEO | built | NWS wind chill and heat index, Magnus dew point and relative humidity (2026-11-13) |
| SCI | `sci/physics` | PHYS | built | density, average speed, average power (2026-11-13) |
| ENG | `eng/aerospace` | AERO | planned | the engineering formulas |
| ENG | `eng/chemical` | CHEME | planned | the engineering formulas |
| ENG | `eng/civil` | CIVIL | planned | the engineering formulas |
| ENG | `eng/computing` | COMPUT | built | IPv4 subnetting by CIDR prefix, transmission time over a link (calculator.net check, 2026-11-13) |
| ENG | `eng/controls` | CTRL | planned | the engineering formulas |
| ENG | `eng/electrical` | ELEC | built | Ohm's law with power, series and parallel resistors, wire resistance, resistive voltage drop (2026-11-13) |
| ENG | `eng/fluid` | FLUID | planned | the engineering formulas |
| ENG | `eng/industrial` | INDUS | planned | the engineering formulas |
| ENG | `eng/mechanical` | MECH | planned | the engineering formulas |
| ENG | `eng/telecom` | TELE | planned | the engineering formulas |
| ENG | `eng/thermal` | THERM | planned | the engineering formulas |
| FIN | `fin/business` | BIZ | built | money |
| FIN | `fin/commodity` | COMMOD | planned | money |
| FIN | `fin/economics` | ECON | built | money |
| FIN | `fin/interest` | INT | built | money |
| FIN | `fin/invest` | INVEST | built | money |
| FIN | `fin/loan` | LOAN | built | money |
| FIN | `fin/market` | MKT | built | market microstructure and execution: mid and microprice, spreads, order-book imbalance and order-flow imbalance, trade signing, VPIN, VWAP/TWAP, implementation shortfall, Kyle's λ, Roll spread, Amihud illiquidity, Almgren–Chriss schedule (S19) |
| FIN | `fin/personal` | PERSON | built | salary conversions, retirement and pension savings (D-021); pay-period conversions built 2026-11-13 |
| MED | `med/anthropometry` | ANTHRO | built | body size: BMI, BSA, weight estimates, body fat (D-046) |
| MED | `med/cardio` | CARDIO | built | heart rate, ECG intervals, haemodynamics, vascular physics, cardiovascular and VTE risk scores (D-046, D-060) |
| MED | `med/chemistry` | CHEM | built | derived clinical-chemistry values: corrected calcium, osmolality, eAG, LDL, anion gap (D-046) |
| MED | `med/critical` | CRIT | built | early-warning, sepsis and critical-care scores (D-058) |
| MED | `med/emergency` | EMERG | built | emergency-department scores not owned by a specialty (D-058) |
| MED | `med/endocrine` | ENDO | built | glucose, diabetes and thyroid indices, ADA categories, CGM metrics (D-058) |
| MED | `med/fluids` | FLUID | planned | blood volume, burn area, fluid volumes; dose-type rows held by policy (D-046, D-050) |
| MED | `med/gastro` | GASTRO | built | gastrointestinal and pancreatic scores (D-058) |
| MED | `med/hematology` | HEMAT | built | haematologic prognostic and bleeding scores (D-058) |
| MED | `med/hepatic` | HEPAT | built | liver indices: FIB-4, MELD family, APRI (D-046) |
| MED | `med/infectious` | INFECT | built | infectious-disease scores (D-058) |
| MED | `med/neuro` | NEURO | built | stroke, coma, head-injury and seizure scores (D-058) |
| MED | `med/obstetric` | OB | built | obstetric categories and scores (D-058) |
| MED | `med/oncology` | ONCO | built | oncology scores (D-058) |
| MED | `med/pediatrics` | PEDS | built | paediatric scores and reference values (D-058) |
| MED | `med/pharmacology` | PHARM | planned | pharmacokinetic and dose arithmetic; held by policy (D-046, D-050) |
| MED | `med/psychiatry` | PSYCH | built | psychiatric and substance-use instruments: arithmetic and published bands (D-058) |
| MED | `med/renal` | RENAL | built | GFR estimation, clearance, fractional excretion, KDIGO and AKI categories (D-046) |
| MED | `med/respiratory` | RESP | built | gas exchange, ventilation, diffusing capacity, spirometry and oxygenation indices (D-046) |
| MED | `med/rheumatology` | RHEUM | built | classification-criteria scores in rheumatology (D-058) |
| MED | `med/scales` | SCALE | built | scored clinical scales with their published bands (D-046) |
| MED | `med/thermoregulation` | THERMO | built | body-temperature staging (D-058) |
| MED | `med/trauma` | TRAUMA | built | injury scores and clinical decision rules (D-058) |
| LIF | `life/automotive` | AUTO | built | fuel economy, fuel consumption per 100 units, trip fuel cost (40 CFR 600.002; 2026-11-13) |
| LIF | `life/brewing` | BREW | planned | everyday calculators |
| LIF | `life/cooking` | COOK | planned | everyday calculators |
| LIF | `life/fitness` | FIT | built | strength, load and pace (D-063, D-067) |
| LIF | `life/games` | GAMES | planned | everyday calculators |
| LIF | `life/garden` | GARDEN | planned | everyday calculators |
| LIF | `life/home` | HOME | planned | everyday calculators |
| LIF | `life/music` | MUSIC | planned | everyday calculators |
| LIF | `life/nutrition` | NUTR | built | energy expenditure, DRI tables and %DV (D-058) |
| LIF | `life/photography` | PHOTO | planned | everyday calculators |
| LIF | `life/shopping` | SHOP | built | tip in cents and basis points (2026-11-13) |
| LIF | `life/sports` | SPORT | built | World Handicap System Handicap Index (Rules of Handicapping 2024; 2026-11-13) |
| LIF | `life/travel` | TRAVEL | planned | everyday calculators |

Import directions inside a group: `num/geometry` may import `num/trig`; `num/trig` imports nothing from `num/` (D-009, D-012).
Inside `fin`, `fin/loan` and `fin/invest` may import `fin/interest`; nothing else imports another `fin` package (D-024).
`med/*` packages do not import each other in production code (D-068).

## Where an outside catalog's categories land (vCalc's list, 2026-11-12)

vCalc groups by audience and repeats categories (Music, Sports, Medicine, Business, Automotive and Unit
Conversions each appear more than once). We use it to see what exists, never as a source (see
`.claude/rules/sources.md`). Each category maps to one home; the audience becomes a tag.

| vCalc category | Our home | Tag |
|---|---|---|
| Education, K-12, Textbooks, vCalc Training | no folder: the function lives by its subject | `education` |
| Hobbies, Industries | no folder: the function lives by its subject | `hobby`, `industry` |
| Mathematics: Algebra, Basic Operations, Boolean | `num/algebra` | |
| Calculus, Differential Equations | `num/calculus` | |
| Combinatorics | `num/combinatorics` | |
| Geometry | `num/geometry` | |
| Linear Algebra, Matrix, Vectors, Quaternion | `num/linalg` | |
| Logarithms, Transcendental Numbers | `num/logs` | |
| Probability, Games of Chance (the maths) | `num/prob` | |
| Series and Sequence | `num/series` | |
| Statistics | `num/stats` | |
| Trigonometry | `num/trig` | |
| Group Theory | not planned; ask Devin if needed | |
| Number Theory (primes, factors, divisibility) | `num/integer` | |
| Computer Science | `algo/*`; data size → `unit/datasize`; performance laws and software metrics → `eng/computing`; check digits → `algo/checksum` (D-016) | |
| Unit Conversions (Area, Mass, Volume, RF and the rest) | `unit/<quantity>` | |
| Physics, Physical Properties | `sci/physics` | |
| Chemistry | `sci/chemistry` | |
| Astronomy | `sci/astronomy` | |
| Biology, Life Science, Zoology, Forestry | `sci/biology` | |
| Geoscience | `sci/geoscience` | |
| Meteorology | `sci/meteorology` | |
| Ecology, Environmental Engineering | `sci/ecology` | |
| Psychology | not planned; ask Devin if needed | |
| Mechanical Engineering, Manufacturing | `eng/mechanical` | |
| Electrical Engineering | `eng/electrical` | |
| Civil Engineering, Construction, Architecture and Design, Drilling | `eng/civil` | |
| Fluid Mechanics, Hydraulics | `eng/fluid` | |
| Heat Transfer | `eng/thermal` | |
| Chemical Engineering | `eng/chemical` | |
| Aerospace Engineering, Aviation | `eng/aerospace` | |
| Telecommunication | `eng/telecom` | |
| Measurement and Controls, Systems Engineering | `eng/controls` | |
| Industrial Engineering | `eng/industrial` | |
| Finance | `fin/loan`, `fin/interest`, `fin/invest`; salary and retirement rows in `fin/personal` (D-021) | |
| Business, Hospitality Industry, Food Services (costing) | `fin/business` | |
| Economics | `fin/economics` | |
| Commodities, Coins and Bullion, Precious Metals, Jewelers, Fuel | `fin/commodity` | |
| Medicine, Physiology | `med/<specialty>` per the MED rows above | |
| Health and Fitness | `life/fitness`; energy and nutrient intake in `life/nutrition` | |
| Sports | `life/sports` | |
| Games of Chance (payouts), Tabletop / Board Gaming, On-line Gaming | `life/games` | |
| Cooking, Entertaining | `life/cooking` | |
| Home Improvement, Quilting | `life/home` | |
| Gardening - Landscaping, Lawn and Garden, Outdoor Living, Agriculture | `life/garden` | |
| Travel, Transportation | `life/travel` | |
| Automotive | `life/automotive` | |
| Shopping | `life/shopping` | |
| Photography, Art | `life/photography` | |
| Music | `life/music` | |
| Beer Making | `life/brewing` | |
| Firearms / Shooting, Ammunition Management | not planned; ask Devin | |
| Religion, Community, Off-Beat | not planned; ask Devin | |
