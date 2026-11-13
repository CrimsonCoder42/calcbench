# Heat Index Calculator (row 45)

CLASSIFY: in scope.
FIELDS: calform: airtemperature 85, airtemperatureunit fahrenheit (celsius, kelvin), humidity 70 (%).
calform2: airtemperature2 85, dewpoint 70 (°F/°C/K).
MATCH: NONE. BUILD sci/meteorology SCI-METEO-002 HeatIndexF (Rothfusz, NWS SR 90-23, 1990; NWS WPC heat index
equation page: simple form, 80 °F switch, low/high-RH adjustments). The dew-point form uses SCI-METEO-004.
COMPARE (ours: checks/scratch/meteorology.go):

| inputs | ours | theirs | diff | verdict |
|---|---|---|---|---|
| 85 °F, 70 % | 92.7021 °F | 93 °F | 0 at 0 dp | MATCH |
| 100 °F, 10 % (low-RH adjustment) | 94.1225 °F | 94 °F | 0 at 0 dp | MATCH |
| 84 °F, 90 % (high-RH adjustment) | 98.3425 °F | 98 °F | 0 at 0 dp | MATCH |
| 85 °F, dew point 70 °F (RH 60.8997 % by SCI-METEO-004) | 89.5347 °F | 90 °F | 0 at 0 dp | MATCH |

Submissions: 4.
CLOSE: BUILT SCI-METEO-002 (MATCH).
