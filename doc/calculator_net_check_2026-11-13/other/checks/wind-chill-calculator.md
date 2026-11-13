# Wind Chill Calculator (row 44)

CLASSIFY: in scope.
FIELDS (calform): windspeed 15, windspeedunit mph (km/h, m/s, knots), airtemperature 10, airtemperatureunit
fahrenheit (celsius, kelvin).
MATCH: NONE. BUILD sci/meteorology (planned METEO): SCI-METEO-001 WindChillF. Source: OFCM FCM-R19-2003 (WCT
equation, range T ≤ 50 °F, V ≥ 3 mph); Osczevski & Bluestein, BAMS 2005;86(10):1453-1458.
COMPARE (ours: checks/scratch/meteorology.go):

| inputs | ours | theirs | diff | verdict |
|---|---|---|---|---|
| 10 °F, 15 mph | −6.5895 °F | −7 °F | 0 at 0 dp | MATCH |
| −20 °F, 40 mph | −56.6243 °F | −57 °F | 0 at 0 dp | MATCH |

Submissions: 2.
CLOSE: BUILT SCI-METEO-001 (MATCH).
