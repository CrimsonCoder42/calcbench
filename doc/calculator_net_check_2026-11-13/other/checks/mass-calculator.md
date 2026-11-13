# Mass Calculator (row 22)

CLASSIFY: in scope (m = ρ·V).
FIELDS (calcf): cdensity 8900, cdensityunit (kg/m³ and others), cvolume 1, cvolumeunit (m³, ft³, yd³, in³, km³,
mi³, cm³, mm³, L, mL, pint, quart, gallon US ...).
MATCH: HAVE SCI-PHYS-002 MassFromDensity (built under row 21).
COMPARE (ours: checks/scratch/density-calculator.go):

| inputs | ours | theirs | diff | verdict |
|---|---|---|---|---|
| 1000 kg/m³, 0.25 m³ | 250 kg | 250 kg | 0 | MATCH |

Submissions: 1.
CLOSE: MATCH (SCI-PHYS-002).
