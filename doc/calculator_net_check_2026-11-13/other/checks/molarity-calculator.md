# Molarity Calculator (row 25)

CLASSIFY: in scope (c = n/V, n = m/M).
FIELDS (calcf): cmass 100, cmassunit g (t, kg, g, mg, ug, Da, lb, oz), cmw 58.44, cmwunit g/mol (kg/mol),
cvolume 2, cvolumeunit L (m³, cm³, mm³, L, mL, gal, qt, pt, cup, fl oz, yd³, ft³, in³), cconcentration (empty),
cconcentrationunit M (mM, uM, nM, mol/m³, kg/L, g/L, mg/L, ppm, ppb, kg/m³, g/m³). The blank field is solved.
MATCH: NONE before; BUILD sci/chemistry (planned CHEM): SCI-CHEM-001 AmountFromMass, -002 Molarity,
-003 MolarityFromMass, -004 MassForMolarity. Source: OpenStax Chemistry 2e (2019) §3.1, §3.3; IUPAC Gold Book
"amount concentration", "molar mass".
COMPARE (ours: hand values in molarity_test.go, same functions):

| inputs | ours | theirs | diff | verdict |
|---|---|---|---|---|
| 100 g, 58.44 g/mol, 2 L → c | 0.8555783709787816 M; n 1.7111567419575633 mol | 0.85557837097878 M; 1.7111567419576 mol | 0 at 14 dp | MATCH |
| 0.5 M, 2 L, 58.44 g/mol → m | 58.44 g | 58.44 g | 0 | MATCH |

Mass concentration (their 50 g/L = m/V) is not a separate function (one division).
Submissions: 2.
CLOSE: BUILT SCI-CHEM-001..004 (MATCH).
