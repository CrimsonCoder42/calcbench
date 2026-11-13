// Package hematology computes haematology and haematological-malignancy prognostic scores from their primary
// sources: NCCN-IPI, GIPSS, MAPS, MALT-IPI, MYSEC-PM and the RIETE bleeding score.
//
// ID short form: HEMAT (IDs are MED-HEMAT-NNN).
//
// Units: every input names its unit (ageYears, ldhRatioToULN, hemoglobinGDL, plateletsE9PerL, blastsPercent,
// creatinineMgDL); findings are booleans. An input outside its range returns ErrOutOfRange naming it; NaN and
// ±Inf return ErrNotFinite. Category functions return the source's own group label, verbatim, as a named
// string type.
//
// Results are numbers and the source's own category labels; nothing here diagnoses, interprets or advises.
//
// Functions:
//
//	MED-HEMAT-001  NCCNIPIScore        NCCN-IPI score (0–8) for diffuse large B-cell lymphoma
//	MED-HEMAT-002  NCCNIPIRiskGroup    NCCN-IPI risk group for a score
//	MED-HEMAT-003  GIPSSScore          GIPSS genetic score for primary myelofibrosis
//	MED-HEMAT-004  GIPSSRiskCategory   GIPSS risk category for a score
//	MED-HEMAT-005  MAPSScore           MAPS clinical-molecular score (0–6) for systemic mastocytosis
//	MED-HEMAT-006  MALTIPIScore        MALT-IPI, the number of its three risk factors present
//	MED-HEMAT-007  MALTIPIRiskGroup    MALT-IPI risk group for a score
//	MED-HEMAT-008  MYSECPMScore        MYSEC-PM score for secondary myelofibrosis (points, with age × 0.15)
//	MED-HEMAT-009  RIETEBleedingScore  RIETE major-bleeding score (0–8) during anticoagulation for VTE
//	MED-HEMAT-010  RIETERiskGroup      RIETE risk group for a score
//	MED-HEMAT-011  DuvalCIBMTRScore                 Duval 2010 pre-transplant score (0–6) for AML in relapse or induction failure
//	MED-HEMAT-012  ICCPrimaryMyelofibrosis          ICC 2022 PMF stage whose criteria are met: pre-PMF or overt PMF
//	MED-HEMAT-013  ICCSystemicMastocytosis          whether the ICC 2022 systemic mastocytosis criteria are met
//	MED-HEMAT-016  IMPROVEBleedingScore             IMPROVE bleeding risk score (0–30.5) for medical inpatients
//	MED-HEMAT-017  CRPUpperReferenceLimitWenerMgDL  age- and sex-adjusted CRP upper reference limit, mg/dL (Wener 2000)
//	MED-HEMAT-019  PLASMICScore                     PLASMIC score (0–7), Bendapudi 2017
//	MED-HEMAT-020  PLASMICRiskGroup                 PLASMIC risk group, the source's label
package hematology
