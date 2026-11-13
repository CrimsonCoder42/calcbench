// Package scales holds clinical scoring scales: the Epworth Sleepiness Scale total and its published bands,
// the Glasgow Coma Scale total and Rimel's head-injury groups, and the Apgar score and Apgar's 1953 groups.
//
// ID short form: SCALE (IDs are MED-SCALE-NNN).
//
// Units: every item and total is an integer count of points, with the range its source allows; an item
// outside that range returns ErrOutOfRange naming the item. Category functions return the source's own
// label, verbatim, as a named string type. No item wording is embedded for the Epworth Sleepiness Scale
// (the scale is © Murray W Johns and its text is licensed); only the sum and the published bands are here.
//
// Results are numbers and the source's own category labels; nothing here diagnoses, interprets or advises.
//
// Functions:
//
//	MED-SCALE-001  EpworthSleepinessScore    Epworth Sleepiness Scale total (0–24) from the 8 item ratings
//	MED-SCALE-002  EpworthCategory           Epworth band for a total, the official site's own label
//	MED-SCALE-003  GlasgowComaScore          Glasgow Coma Scale total (3–15) from eye, verbal and motor ratings
//	MED-SCALE-004  GlasgowComaSeverityRimel  Rimel 1982 head-injury group for a GCS total: severe, moderate, minor
//	MED-SCALE-005  ApgarScore                Apgar score (0–10) from the five signs, each rated 0–2
//	MED-SCALE-006  ApgarCondition1953        Apgar 1953 group for a score: good, fair or poor condition
//	MED-SCALE-007  HARKScore                 HARK score (0–4): one point per yes to the four HARK questions
//	MED-SCALE-008  DangerAssessment5Count    count (0–5) of yes answers to the five DA-5 questions
//	MED-SCALE-011  NeckDisabilityIndex  Neck Disability Index (0–50), Vernon & Mior 1991
//	MED-SCALE-012  NDIDisability        Neck Disability Index band, the source's label
//	MED-SCALE-013  WASTScore            Woman Abuse Screening Tool sum (8–24), Brown 1996
//	MED-SCALE-014  FLACCScore           FLACC behavioural pain score (0–10), arithmetic only
//	MED-SCALE-015  PAINADScore          PAINAD pain score (0–10), arithmetic only
package scales
