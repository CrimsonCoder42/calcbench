# pace-calculator

## FIELDS (2026-11-13)
calform: ctime 00:50:25 (hh:mm:ss) · cdistance 5 · cdistanceunit Miles|Kilometers|Meters|Yards · cpickevent (Marathon, Half, 1K, 5K, 10K, 1M, 5M...) ·
cpace 00:08:10 · cpaceunit tpm|tpk|mph|kph|mpm|mps|ypm|yps.
calform2: multipoint pace (cmpdistanceN, cmpdistanceunitN, cmptimeN, N = 1..12) · calform3: pace converter (ccfrom, ccpaceunitfrom/to) ·
calform4: finish-time projection (cftdistancenow, cfttime, cftdistancefull).
Results: pace in eight units, race-distance times.

## MATCH
HAVE LIF-FIT-009 Pace, LIF-FIT-010 Splits. Unit conversion mi→km (× 1.609344) done in scratch only; the converter and finish-time
projection forms are plain rate arithmetic (Pace × distance) and were not compared.

## COMPARE (scratch: checks/scratch/pace-calculator.go)
| inputs | ours | theirs | diff | verdict |
|---|---|---|---|---|
| 50:25 over 5 mi | 605.000 s/mi; 375.930 s/km | 10:05 /mi; 6:15.93 /km | 0 | MATCH |
| 45:00 over 10 km | 270.000 s/km | 4:30 /km | 0 | MATCH |
| 3:30:15 over 42.195 km | 298.969 s/km | 4:58.97 /km | 0 at 0.01 s | MATCH |
Submissions: 1 fields + 3 runs.
