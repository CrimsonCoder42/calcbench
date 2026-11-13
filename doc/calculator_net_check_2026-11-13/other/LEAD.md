# other: notes for the lead (A-OTHER)

## New packages created (all planned in doc/taxonomy.md; doc.go headers to check at merge)

| path | SHORT | IDs |
|---|---|---|
| sci/astronomy | ASTRO | SCI-ASTRO-001..006 |
| sci/physics | PHYS | SCI-PHYS-001..007 |
| sci/chemistry | CHEM | SCI-CHEM-001..004 |
| sci/meteorology | METEO | SCI-METEO-001..004 |
| eng/electrical | ELEC | ENG-ELEC-001..011 |
| eng/computing | COMPUT | ENG-COMPUT-001..002 |
| unit/area | AREA | UNT-AREA-001..002 |
| sci/geoscience | GEO | SCI-GEO-001..003 (session 3, distance-calculator handoff from A-MATH) |

Taxonomy status for these seven rows: planned → built. Every function has its doc comment (ID, Time, Space,
Source, Tags, Units, Example; Range where the source states one), a _test.go with oracle, edge and error cases
and a benchmark. Not run: go test (per brief). gofmt -l, go build ./..., go vet on the seven packages: clean.
eng/electrical imports algo/reduce (NeumaierAdd, ALG-RED-002) for the series/parallel sums: allowed by the
import rule.

## INDEX rows to regenerate (D-070): 36 new functions

SCI-ASTRO-001 GregorianToJulianDayNumber · 002 JulianDayNumberToGregorian · 003 GregorianDayOfWeek ·
004 GregorianDayOfYear · 005 DaysBetweenDates · 006 GregorianAddDays (Meeus, Astronomical Algorithms 2e, ch. 7).
SCI-PHYS-001 Density · 002 MassFromDensity · 003 VolumeFromDensity · 004 AverageSpeed · 005 DistanceFromSpeed ·
006 DurationFromSpeed · 007 PowerFromWork (OpenStax University Physics).
SCI-CHEM-001 AmountFromMass · 002 Molarity · 003 MolarityFromMass · 004 MassForMolarity (OpenStax Chemistry 2e;
IUPAC Gold Book).
SCI-METEO-001 WindChillF (OFCM FCM-R19-2003) · 002 HeatIndexF (Rothfusz 1990; NWS WPC) · 003 DewPointC ·
004 RelativeHumidityFromDewPoint (Alduchov & Eskridge 1996; Lawrence 2005).
ENG-ELEC-001..006 OhmsLawFrom{VoltageCurrent, VoltageResistance, VoltagePower, CurrentResistance, CurrentPower,
ResistancePower} · 007 VoltageDropTwoWire · 008 VoltageDropThreePhase (Glover et al. 5e §2.5) ·
009 SeriesResistance · 010 ParallelResistance · 011 WireResistance.
ENG-COMPUT-001 IPv4SubnetOf (RFC 4632, RFC 950, RFC 3021) · 002 TransmissionTime (Kurose & Ross 7e §1.4.1).
UNT-AREA-001 SquareFeetToSquareMeters · 002 SquareMetersToSquareFeet (NIST SP 811 B.8).

## doc/sources notes for the lead (sources used; worked examples reproduced in tests)
- Meeus, Astronomical Algorithms 2e (1998) ch. 7: Examples 7.a, 7.c, 7.d, 7.e, 7.f in sci/astronomy tests; plus
  a 500-point round-trip oracle against Go's time package (proleptic Gregorian).
- OFCM FCM-R19-2003 wind chill chart values (−7, −57, 36, −26 °F) in wind_chill_test.go.
- Several sources were cited from memory of their section numbers (OpenStax sections, Glover §2.5, Kurose §1.4.1,
  Lawrence eq. numbers, Alduchov eq. 21): the reviewer should open each and confirm the locator.

## Decisions for Devin
1. (Observed in session 2: calculator.net counts a 29 Feb birthday's common-year anniversary as 28 Feb.) Age in "Y years M months D days" and date + N months/years: a convention (borrowing days; 29 Feb anniversary
   in a common year; end-of-month clamping). Recommendation: anniversary of 29 Feb = 1 Mar in common years
   (the legal reading in e.g. UK Family Law Reform Act 1969 s.9 is a jurisdiction; needs a chosen source).
   Not built.
2. Molecular weight: adopt CIAAW 2021 abridged standard atomic weights as a versioned constant table? (row 26).
3. NEC Chapter 9 Tables 8/9 (conductor R and X) are NFPA copyrighted data: keep conductor resistance a caller
   input (current design) or license the table? (row 28).
4. unit/power (W, hp 745.69987158227022 W, metric hp, BTU/h) is not in doc/taxonomy.md: proposed path
   unit/power, SHORT POWER, Holds "power conversions (NIST SP 811 B.8)". Needed by rows 31 and 49.
5. Clock arithmetic rows (time, hours, time-duration, time-card) were SKIPPED as stdlib time.Duration work;
   confirm, or ask for unit/duration (planned DUR) conversions.
6. Final-grade solve ("score needed on the final") home: num/stats or life/everyday? Not proposed as a
   package; no life/everyday package was needed this session.
7. Unit slice: rows 10, 17, 19, 23, 24, 34 want more conversions (ft, yd, mi, km in unit/length; ft³, yd³, L,
   gal in unit/volume; mph, km/h, knot in unit/speed; B/bit and SI/IEC prefixes in unit/datasize). All planned
   packages; a build slice of their own.

## Handoffs
- A-FIT: tile (life/home), mulch, gravel (life/garden), height (med/peds), tip (life/shopping), golf-handicap
  (life/sports), sleep (med), fuel-cost, gas-mileage, mileage, tire-size (life/automotive).
- A-FIN: gdp (fin/economics), electricity cost (fin; energy part noted).
- A-MATH: roman-numeral-converter (algo/text).

## Driver gaps (DRIVER)
Session 2: date pickers and live mode are fixed. Still open: (1) voltage-drop has no named or id'd element that
selects the resistance mode (NEC table / wire size / custom); (2) conversion-calculator's mass, temperature and area
tabs are not listed by the fields script. Withdrawn: the ohms-law 'first result line only' gap (my jq slice).

## _delete moves (path | replaced by | date)
- unit/area/area_test.go → _delete/unit/area/area_test.go | split into square_feet_to_square_meters_test.go and
  square_meters_to_square_feet_test.go (one test file per concept) | 2026-11-13 (file created and moved in this
  session; never shipped)

## FOUND, NOT FIXED
None.

## Session 3
- PROGRESS.md repaired: another agent's same-named helper (setrow.py in the shared scratchpad) had written
  financial closes into about 35 of my rows. Rebuilt from checks/*.md CLOSE lines with
  scratchpad/aother/aother_repair_progress.py; 0 FIN- IDs left in rows. My old setrow.py is not used again.
- New INDEX rows: SCI-GEO-001 GreatCircleDistance, SCI-GEO-002 VincentyInverseDistance, SCI-GEO-003
  WGS84Distance; exported constant MeanEarthRadiusMeters (GRS 80 R1, Moritz). sci/geoscience imports num/trig
  (CentralAngle, NUM-TRIG-031).
- The Vincenty line (a) oracle (14 110 526.170 m) was recalled from the paper, not read this session; our
  value 14 110 526.1696 m agrees to 0.4 mm, but the reviewer should open the paper and confirm.
- Possible later build: Karney 2013 geodesic solver to replace Vincenty near antipodes (ErrNoConvergence today).

## Reviewer follow-ups (source citations)
- Done in code: RFC 4632, RFC 950, RFC 3021 Source lines now carry doi:10.17487/RFC4632, 10.17487/RFC0950,
  10.17487/RFC3021 (eng/computing/ipv4_subnet.go). OpenStax Source lines now name the authors and URL:
  Ling, Sanny, Moebs, University Physics Vols 1-2 (2016); Flowers, Theopold, Langley, Robinson, Chemistry 2e
  (2019) (sci/physics, sci/chemistry, eng/electrical). INDEX rows for those IDs must be regenerated.
- Still to confirm by opening the sources: the OpenStax section numbers, Glover §2.5, Kurose §1.4.1, the Lawrence
  and Alduchov equation numbers, Vincenty line (a).
