# Horsepower Calculator (row 49)

CLASSIFY: in scope (P = F·d/t; power-unit conversion).
FIELDS: calform: force 10, forceunit newton (kN, lbf, kgf), distance 3, distanceunit meter (km, mile, yard),
time 2, timeunit second (minute, hour, day). calform2: power converter convertamount 1, convertfrom/convertto
(W, kW, mechanical hp, metric hp, electrical hp, boiler hp, BTU/h).
MATCH: NONE. BUILD sci/physics SCI-PHYS-007 PowerFromWork (OpenStax University Physics Vol. 1 §7.1, §7.4).
Power-unit conversions (hp = 550 ft·lbf/s = 745.69987158227022 W, NIST SP 811 B.8) need unit/power, which is
not in the taxonomy (LEAD.md).
COMPARE (1 lbf = 4.4482216152605 N, 1 mi = 1609.344 m, both exact, NIST SP 811 B.8; converted by hand):

| inputs | ours | theirs | diff | verdict |
|---|---|---|---|---|
| 10 N, 3 m, 2 s | 15 W | 15 W | 0 | MATCH |
| 1000 lbf, 1 mile, 1 min | 119311.97945 W (160 hp) | 119,311.9794 W | 0 at 4 dp | MATCH |

Submissions: 2.
CLOSE: BUILT SCI-PHYS-007 (MATCH).
