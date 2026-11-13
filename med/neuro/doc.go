// Package neuro computes neurology scores and published bands from their primary sources: stroke and TIA
// scores, intracerebral hemorrhage volume and grading, coma and head-injury indices, seizure risk on EEG,
// delirium screening bands and neuro-oncology groups.
//
// ID short form: NEURO (IDs are MED-NEURO-NNN).
//
// Units: every input and result names its unit (centimetres, mL, years, minutes, mmol/L, points). Item scores
// are integer points checked against the range the source gives; NaN and ±Inf return ErrNotFinite; inputs
// outside the source's range return ErrOutOfRange. The Glasgow Coma Scale total itself is MED-SCALE-003.
//
// Results are numbers and the source's own category labels; nothing here diagnoses, interprets or advises.
//
// Functions:
//
//	MED-NEURO-001  ABC2HemorrhageVolumeML  intracerebral hemorrhage volume by ABC/2 (Kothari 1996)
//	MED-NEURO-002  ABCD2Score              ABCD² score (0–7) after TIA (Johnston 2007)
//	MED-NEURO-003  ABCD2RiskGroupOf        ABCD² stratum: low, moderate or high risk
//	MED-NEURO-004  CPSSSScore              Cincinnati Prehospital Stroke Severity Scale (0–4) from NIHSS items
//	MED-NEURO-005  DRAGONScore             DRAGON score (0–10) after IV thrombolysis
//	MED-NEURO-006  ICHScore                ICH Score (0–6) (Hemphill 2001)
//	MED-NEURO-007  SEDANScore              SEDAN score (0–6) for post-thrombolysis hemorrhage
//	MED-NEURO-008  SICHScore               secondary intracerebral hemorrhage score (0–6)
//	MED-NEURO-009  FASTEDScore             FAST-ED score (0–9)
//	MED-NEURO-010  RACEScore               RACE scale (0–9)
//	MED-NEURO-011  PREScore                Pittsburgh Response to Endovascular therapy score
//	MED-NEURO-012  TwoHELPS2BScore         2HELPS2B EEG seizure-risk score (0–7)
//	MED-NEURO-013  GCSPupilsScore          GCS-Pupils score (1–15)
//	MED-NEURO-014  TokuhashiRevisedGroup   revised Tokuhashi group for a total (0–15)
//	MED-NEURO-015  FourATCategoryOf        4AT band for a total (0–12)
//	MED-NEURO-016  CanadianCTHeadRule                Canadian CT Head Rule: high- and medium-risk factor groups
//	MED-NEURO-017  CATCHRule                         CATCH rule: high- and medium-risk factor groups (children)
//	MED-NEURO-018  NewOrleansCriteriaCount           New Orleans criteria present (0–7)
//	MED-NEURO-019  NEXUSHeadCTLowRisk                NEXUS Head CT low-risk profile: under 65, none of seven findings
//	MED-NEURO-020  OttawaSAHRuleCount                Ottawa SAH Rule findings present (0–6)
//	MED-NEURO-021  UCDavisHeadTraumaLowRisk          UC Davis rule: none of five predictors (children)
//	MED-NEURO-022  PECARNHeadInjuryVeryLowRisk       PECARN head-trauma very-low-risk rule for the child's age
//	MED-NEURO-023  HINTSAnyFinding                   any of the three HINTS findings present
//	MED-NEURO-024  ASPECTSScore                      ASPECTS (0–10): 10 minus regions with early ischemic change
//	MED-NEURO-025  LAMSScore                         Los Angeles Motor Scale (0–5)
//	MED-NEURO-026  FOURScore                         FOUR score (0–16)
//	MED-NEURO-027  CAMICUPositive                    CAM-ICU algorithm: features 1 and 2 plus 3 or 4
//	MED-NEURO-028  ModifiedFisherGrade               modified Fisher grade (0–4) from SAH thickness and IVH
//	MED-NEURO-030  CanadianTIAScore                  Canadian TIA Score (−3 to 23)
//	MED-NEURO-032  AcuteStrokeDysphagiaScreenFailed  Acute Stroke Dysphagia Screen: any item present
//	MED-NEURO-033  MTOQ4Score                        mTOQ-4 sum (0–8)
//	MED-NEURO-034  MTOQ4CategoryOf                   mTOQ-4 treatment-efficacy category for a sum (0–8)
//	MED-NEURO-035  MoCATotal                         MoCA total (0–30) with the education point
package neuro
