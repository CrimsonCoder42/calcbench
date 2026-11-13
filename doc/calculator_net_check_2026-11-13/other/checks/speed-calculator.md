# Speed Calculator (row 24)

CLASSIFY: in scope (v = d/t, solve for any one).
FIELDS (calcf): cspeed (empty), cspeedunit ms (kmh, mph, kn, fts, ... c), cdistance 100, cdistanceunit m (mm, cm,
m, km, in, ft, yd, mi, nmi, ly), ctimeh / ctimem / ctimes (default 20 s). calcf2: speed unit converter fromVal 1,
calFrom mph, calTo ms.
MATCH: NONE before; BUILD sci/physics SCI-PHYS-004 AverageSpeed, -005 DistanceFromSpeed, -006 DurationFromSpeed.
Source: OpenStax University Physics Vol. 1 (2016) §3.1. Speed-unit conversions (unit/speed, planned) not built.
COMPARE (ours: 100000/3600 and 100/20 by SCI-PHYS-004):

| inputs | ours | theirs | diff | verdict |
|---|---|---|---|---|
| 100 km in 1 h | 27.777777777778 m/s | 27.777777777778 | 0 at 12 dp | MATCH |
| 100 m in 20 s | 5 m/s | 5 | 0 | MATCH |

Submissions: 2.
CLOSE: BUILT SCI-PHYS-004..006 (MATCH).
