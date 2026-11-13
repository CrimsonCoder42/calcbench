// Package fitness holds resistance-training, powerlifting and running calculators: one-repetition-maximum
// estimates from reps to fatigue (Brzycki, Wathen), the load at a percent of 1 RM, a set's and a session's
// volume load, the ACSM 2009 load-increase range, reps in reserve from RPE, the IPF GL coefficient and points,
// and pace and splits.
//
// ID short form: FIT (IDs are LIF-FIT-NNN).
//
// Units: a load whose formula is a ratio or a product (1 RM estimates, percent of 1 RM, volume load, load
// progression) is taken in kg or lb and returned in the same unit; the parameter name ends in KgOrLb to say
// so. A formula with a body weight inside a non-linear term (IPF GL) takes kg only, as its source does.
// Repetitions and scale points are int. Every function checks NaN and ±Inf (ErrNotFinite) and its range
// (ErrOutOfRange); a result beyond float64 is ErrDomain.
//
// Results are numbers and the source's own labels; nothing here advises.
//
// Functions:
//
//	LIF-FIT-001  LoadProgressionACSM     ACSM 2009 2–10% load-increase range and whether its rep condition is met
//	LIF-FIT-002  IPFGLCoefficient        IPF GL coefficient (2020) for body weight, sex and competition type
//	LIF-FIT-003  IPFGLPoints             IPF GL points (2020) for a result, rounded per the IPF procedure
//	LIF-FIT-004  (withdrawn before release: Mayhew 1992 primary not seen; ID not reused)
//	LIF-FIT-005  (withdrawn before release: Mayhew 1992 primary not seen; ID not reused)
//	LIF-FIT-006  RepsInReserve           repetitions in reserve from an RPE of 7–10 (10 − RPE)
//	LIF-FIT-007  LoadAtPercentOneRepMax  training load at a percent of 1 RM
//	LIF-FIT-008  VolumeLoad              one set's volume load, repetitions × external load
//	LIF-FIT-009  Pace                    pace, in seconds per unit of distance
//	LIF-FIT-010  Splits                  split times of a cumulative distance stream for whole units of unitKmOrMi
//	LIF-FIT-011  TotalVolumeLoad         volume load of a session or a week
//	LIF-FIT-012  OneRepMaxBrzycki  1 RM from reps to fatigue, Brzycki: load / (1.0278 − 0.0278 × reps)
//	LIF-FIT-014  OneRepMaxWathen   1 RM from reps to fatigue, Wathen: 100 × load / (48.8 + 53.8 e^(−0.075 reps))
//	LIF-FIT-015  EnergyExpenditureMET  energy expended (kcal) from MET, body mass (kg) and duration (h), 1 MET ≈ 1.0 kcal/kg/h
package fitness
