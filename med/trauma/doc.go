// Package trauma computes injury scores and applies trauma decision rules from their primary sources: the
// ABC massive-transfusion score, the Canadian C-Spine Rule, NEXUS (cervical spine, chest, chest CT), the New
// Trauma Score, the Ottawa ankle (1992 derivation) and knee rules, the PECARN and PedSRC paediatric abdominal
// rules, the Pittsburgh knee rule, RibScore, the TASH score and the Wisconsin criteria.
//
// ID short form: TRAUMA (IDs are MED-TRAUMA-NNN).
//
// Units: every numeric input names its unit (ageYears, systolicMmHg, heartRatePerMin, spo2Percent,
// hemoglobinGDL, baseExcessMmolL, astUL); gcs is the Glasgow Coma Scale total, 3–15. A finding that a rule
// words as an absence ("no abdominal tenderness") is passed as the presence of that finding, and the field's
// comment quotes the source. NaN or ±Inf returns ErrNotFinite before any range check; an input outside its
// physical range returns ErrOutOfRange. Decision rules return the rule's own result (a count of criteria or a
// boolean), named as the source names it.
//
// Results are numbers and the source's own category labels; nothing here diagnoses, interprets or advises.
//
// Functions:
//
//	MED-TRAUMA-001  ABCScore                               ABC score (0–4) for massive transfusion
//	MED-TRAUMA-002  CanadianCSpineRule                     Canadian C-Spine Rule's three questions and outcome
//	MED-TRAUMA-003  NEXUSChest                             count of the seven NEXUS Chest criteria present
//	MED-TRAUMA-004  NEXUSChestCT                           NEXUS Chest CT-All and Chest CT-Major: any criterion present
//	MED-TRAUMA-005  NEXUSCSpineLowProbability              NEXUS cervical spine: all five low-probability criteria met
//	MED-TRAUMA-006  NewTraumaScore                         New Trauma Score, triage sum and probability of survival
//	MED-TRAUMA-007  OttawaAnkleRule1992                    1992 Ottawa ankle malleolar and midfoot finding sets
//	MED-TRAUMA-008  OttawaKneeRule                         count of the five Ottawa knee rule variables present
//	MED-TRAUMA-009  PECARNIntraAbdominalInjuryVeryLowRisk  PECARN blunt abdominal rule: every item met
//	MED-TRAUMA-010  PedSRCVeryLowRisk                      PedSRC rule: count of its five variables and none present
//	MED-TRAUMA-011  PittsburghKneeRule                     Pittsburgh knee rule combination
//	MED-TRAUMA-012  RibScore                               RibScore (0–6) from six radiographic variables
//	MED-TRAUMA-013  TASHScore                              TASH score (0–28)
//	MED-TRAUMA-014  WisconsinCriteria                      count of the five Wisconsin maxillofacial CT criteria present
//	MED-TRAUMA-015  TLICS                         Thoracolumbar Injury Classification and Severity Score, 1–10
//	MED-TRAUMA-016  InjurySeverityScore           Injury Severity Score from six AIS region maxima, 0–75
//	MED-TRAUMA-017  SINSScore                     Spinal Instability Neoplastic Score, 0–18
//	MED-TRAUMA-018  PediatricTraumaScore          Pediatric Trauma Score, −6 to 12
//	MED-TRAUMA-019  MESSScore                     Mangled Extremity Severity Score
//	MED-TRAUMA-020  RevisedTraumaScore            Revised Trauma Score: GCS/SBP/RR codes, T-RTS 0–12 and weighted RTS
//	MED-TRAUMA-021  NottinghamHipFractureScore    Nottingham Hip Fracture Score (0–10), haemoglobin in g/dL
//	MED-TRAUMA-022  NottinghamHipFractureScoreSI  Nottingham Hip Fracture Score (0–10), haemoglobin in g/L
//	MED-TRAUMA-023  PECARNCSpineRule              PECARN cervical spine rule group: high-risk, CART-identified or neither factor
package trauma
