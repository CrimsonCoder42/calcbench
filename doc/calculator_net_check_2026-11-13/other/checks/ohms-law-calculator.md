# Ohm's Law Calculator (row 30)

CLASSIFY: in scope.
FIELDS (calform): v 110 (kV/V/mV), c 0.2 (A/mA), r empty (Ω/kΩ/MΩ), p empty (W/kW/MW/hp/btu). Any two given.
MATCH: NONE. BUILD eng/electrical ENG-ELEC-001..006 (one function per known pair, result Circuit{Volts, Amperes,
Ohms, Watts}). Source: OpenStax University Physics Vol. 2 (2016) §9.4 Ohm's law, §9.5 power.
COMPARE (ours: checks/scratch/ohms-law-calculator.go):

| inputs | ours | theirs | diff | verdict |
|---|---|---|---|---|
| 110 V, 0.2 A | R 550 Ω (P 22 W) | R 550 Ω | 0 | MATCH |
| 47 Ω, 5 W | I 0.32616403652672 A, V 15.329709716755891 V | I 0.32616403652672 A, V 15.329709716756 V | 0 at 12 dp | MATCH |

Correction to session 1: the "DRIVER: only the first result line" note was my error (my jq kept result.text[0:2]);
result.text carries every solved quantity. fieldsAfter shows the page does not write answers back into v/c.
Submissions: session 1: 2; session 2: 2. Total 4.
CLOSE: BUILT ENG-ELEC-001..006 (MATCH).
