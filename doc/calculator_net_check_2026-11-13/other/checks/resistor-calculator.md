# Resistor Calculator (row 29)

CLASSIFY: in scope: colour-code value, series and parallel networks, wire resistance.
FIELDS: (no form name) bandnum 3..6 (default 4), band1/band2/band3 radios black..white (green, red, blue),
multiplier radio black..white, gold, silver (blue), tolerance radio (gold), temperatureCoefficient radio (brown).
calform: pnumberinputs textarea "10,2,38,23,38,23,21" (parallel). calform2: snumberinputs (series, same default).
calform3: length 100 foot, diameter 0.05 inch, conductivity 59600000 S/m (material select 6.3e7, 5.96e7, 5.8e7,
4.52e7, 3.78e7).
MATCH: NONE. BUILD eng/electrical: ENG-ELEC-009 SeriesResistance, ENG-ELEC-010 ParallelResistance (OpenStax
University Physics Vol. 2 §10.2 eqs. 10.2-10.3, Neumaier-compensated via ALG-RED-002), ENG-ELEC-011
WireResistance (§9.3, R = ρL/A).
Colour code NOT built: the digit/multiplier code is IEC 60062:2016 "Marking codes for resistors and capacitors"
§3 and the tolerance/temperature-coefficient letters changed between editions; the standard was not in hand to
cite clause and table (NEEDS SOURCE for that part).
COMPARE (ours: hand values in the tests; 100 ft = 30.48 m, 0.05 in = 0.00127 m exact, NIST SP 811 B.8):

| inputs | ours | theirs | diff | verdict |
|---|---|---|---|---|
| series 10,2,38,23,38,23,21 Ω | 155 Ω | 155 | 0 | MATCH |
| parallel same | 1.2703136679493923 Ω | 1.2703136679494 | 0 at 13 dp | MATCH |
| wire 100 ft, Ø 0.05 in, 5.96e7 S/m | 0.4037117406548103 Ω | 0.40371174065481 | 0 at 14 dp | MATCH |

Submissions: 3.
CLOSE: BUILT ENG-ELEC-009..011 (MATCH; colour code NEEDS SOURCE).
