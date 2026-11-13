// Package gastro holds gastroenterology scores: BISAP for acute pancreatitis, the Rockall scores for upper
// gastrointestinal haemorrhage, the Pediatric Appendicitis Score, the RIPASA appendicitis score and the I-SEE
// severity index for eosinophilic esophagitis with its categories.
//
// ID short form: GASTRO (IDs are MED-GASTRO-NNN).
//
// Units: every numeric input names its unit (BUN mg/dL, blood pressure mmHg, pulse beats/min, eosinophils per
// high-power field, age in completed decimal years); clinical findings are booleans or named levels whose
// constants quote the source's table cells. NaN or ±Inf in any input returns ErrNotFinite before any range
// check; a negative measurement or an unknown level returns ErrOutOfRange. Totals are integer points, except
// RIPASA, whose half-points are exact in float64.
//
// Results are numbers and the source's own category labels; nothing here diagnoses, interprets or advises.
//
// Functions:
//
//	MED-GASTRO-001  BISAPScore                  BISAP (0–5), Wu 2008
//	MED-GASTRO-002  RockallScorePreEndoscopy    Rockall score before endoscopy (0–7), Rockall 1996 Table III
//	MED-GASTRO-003  RockallScoreComplete        complete Rockall score (0–11), Rockall 1996 Table III
//	MED-GASTRO-004  PediatricAppendicitisScore  Pediatric Appendicitis Score (0–10), Samuel 2002
//	MED-GASTRO-005  RIPASAScore                 RIPASA score (1.5–16), Chong 2010
//	MED-GASTRO-006  ISEEScore                   I-SEE total for eosinophilic esophagitis, Dellon 2022 Table 1
//	MED-GASTRO-007  ISEESeverity                I-SEE severity category, the source's own label
//	MED-GASTRO-008  AcutePancreatitisDiagnosisScore  Jin acute pancreatitis model (−3 to 14), Jin 2024 Table 1
//	MED-GASTRO-009  RansonScore                      Ranson's 11 early signs, non-gallstone (0–11)
//	MED-GASTRO-010  GlasgowBlatchfordScore           Glasgow-Blatchford score (0–23), Blatchford 2000 via Banister 2018 Table 1
//	MED-GASTRO-011  AlvaradoScore                    Alvarado score (0–10), Alvarado 1986 via Alvarado 2016
//	MED-GASTRO-013  OaklandScore                     Oakland score for lower GI bleeding (0–33), Oakland 2017 Table 2
//	MED-GASTRO-014  PUCAIScore                       Pediatric Ulcerative Colitis Activity Index (0–85), Turner 2007
package gastro
