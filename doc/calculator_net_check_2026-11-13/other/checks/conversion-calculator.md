# Conversion Calculator (row 19)

CLASSIFY: in scope (unit conversions, NIST SP 811 Appendix B).
FIELDS (calForm): fromVal text (empty), toVal text, calFrom select 1..11, calTo select 1..11; length tab labels
Meter, Kilometer, Centimeter, Millimeter, Micrometer, Nanometer, Mile, Yard, Foot, Inch, Light Year (default
Meter → Kilometer). Other quantities are other tabs.
MATCH: PARTIAL UNT-LEN-001/002 (cm↔in), UNT-MASS-001/002, UNT-TEMP-001..004, UNT-AREA-001/002 (new, row 12).
Not covered: m/km/mm/mile/yard/foot/light-year lengths, volume, speed, pressure, energy, duration.
COMPARE (session 2, live mode; ours: checks/scratch/conversion-calculator.go):

| inputs | ours | theirs | diff | verdict |
|---|---|---|---|---|
| 12 in → cm | 30.48 (UNT-LEN-002) | 30.48 | 0 | MATCH |
| 100 cm → in | 39.3700787402 (UNT-LEN-001) | 39.37007874 | 0 at 8 dp | MATCH |

Only the length tab's selects (calFrom/calTo) are fields on this page; mass, temperature and area are other
tabs the fields script does not list, so UNT-MASS/TEMP/AREA were not compared here (mass was compared on the
weight-calculator row, area on square-footage).
Session 1: 1 submission (exit 4). Session 2: 2. Total 3.
CLOSE: MATCH (UNT-LEN-001/002; other quantities PARTIAL). The remaining unit packages (unit/volume, unit/speed, unit/pressure, unit/energy,
more unit/length) are a build slice of their own: LEAD.md.
