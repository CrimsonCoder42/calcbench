# Dew Point Calculator (row 46)

CLASSIFY: in scope.
FIELDS (calform): airtemperature 70, airtemperatureunit fahrenheit (celsius, kelvin), humidity 65, dewpoint
(empty), dewpointunit fahrenheit. The blank one of humidity/dewpoint is solved.
MATCH: NONE. BUILD sci/meteorology SCI-METEO-003 DewPointC, SCI-METEO-004 RelativeHumidityFromDewPoint (Magnus
form, Alduchov & Eskridge 1996 coefficients 17.625 / 243.04 °C; Lawrence, BAMS 2005). °F via UNT-TEMP-001/002.
COMPARE (ours: checks/scratch/meteorology.go):

| inputs | ours | theirs | diff | verdict |
|---|---|---|---|---|
| 70 °F, 65 % | 14.2756 °C / 57.6961 °F | 14.3 °C / 57.7 °F | 0 at 1 dp | MATCH |
| 32 °F, 80 % | −3.0386 °C / 26.5306 °F | −3.0 °C / 26.5 °F | 0 at 1 dp | MATCH |

Submissions: 2.
CLOSE: BUILT SCI-METEO-003, SCI-METEO-004 (MATCH).
