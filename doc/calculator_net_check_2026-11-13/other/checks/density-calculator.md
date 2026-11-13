# Density Calculator (row 21)

CLASSIFY: in scope (ρ = m/V, solve for any one).
FIELDS (calcf): cdensity (empty), cdensityunit (kg/m³, kg/cm³, g/m³, g/cm³, kg/L, g/L, lb/in³, lb/ft³ ...),
cvolume 1, cvolumeunit (m³, L, mL, gal US, gal UK, ft³, yd³, in³, cm³), cmass 8900, cmassunit (kg, g, mg, t, lb,
oz, ct). The blank field is the one solved.
MATCH: NONE before; BUILD sci/physics (planned PHYS): SCI-PHYS-001 Density, -002 MassFromDensity,
-003 VolumeFromDensity. Source: OpenStax University Physics Vol. 1 (2016) §14.1; SI Brochure 9th ed. Table 5.
COMPARE (ours: checks/scratch/density-calculator.go; lb via UNT-MASS PoundsToKilograms, 1 ft³ = 0.028316846592 m³
exact per NIST SP 811 B.8, written in the scratch file since unit/volume does not exist yet):

| inputs | ours | theirs | diff | verdict |
|---|---|---|---|---|
| 8900 kg, 1 m³ | 8900 kg/m³ | 8900 | 0 | MATCH |
| 100 lb, 1 ft³ | 1601.846337396 kg/m³ | 1601.846337396 | 0 at 9 dp | MATCH |

Submissions: 2.
CLOSE: BUILT SCI-PHYS-001..003 (MATCH).
