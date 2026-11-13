package rheumatology

// AcuteGoutDiagnosticRuleScore returns the score of the primary-care diagnostic rule for acute gouty arthritis
// without joint fluid analysis (Janssens 2010), 0–13 in steps of 0.5: male sex 2; previous self-reported
// arthritis attack 2; onset within 1 day 0.5; joint redness 1; first metatarsophalangeal joint involvement 2.5;
// hypertension or one or more cardiovascular diseases 1.5; serum uric acid > 5.88 mg/dL (0.35 mmol/L) 3.5. The
// validation reports "The positive predictive value of a score of ≥8 points was 0.87, the negative predictive
// value of a score of ≤4 points was 0.95"; those are its population's figures and are not returned.
// Errors: none; every combination of the items is defined.
//
// ID:      MED-RHEUM-006
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Janssens HJ, Fransen J, van de Lisdonk EH, van Riel PL, van Weel C, Janssen M. A diagnostic rule for acute gouty arthritis in primary care without joint fluid analysis. Arch Intern Med. 2010;170(13):1120-1126. doi:10.1001/archinternmed.2010.196, PMID 20625017
// Source:  Kienhorst LB, Janssens HJ, Fransen J, Janssen M. The validation of a diagnostic rule for gout without joint fluid analysis: a prospective study. Rheumatology (Oxford). 2015;54(4):609-614. doi:10.1093/rheumatology/keu378, PMID 25231179, abstract (cut-offs ≥8 and ≤4)
// Source:  Lee KH, Choi ST, Lee SK, Lee JH, Yoon BY. Application of a Novel Diagnostic Rule in the Differential Diagnosis between Acute Gouty Arthritis and Septic Arthritis. J Korean Med Sci. 2015;30(6):700-704. doi:10.3346/jkms.2015.30.6.700, PMID 26028920, PMC4444468, Methods (the seven variables and points)
// Source:  McConville CI, Cox JM. Spinal gout diagnosis in chiropractic practice: narrative review. J Can Chiropr Assoc. 2023;67(1):50-66. PMID 37250460, PMC10211406, Table 3 (the seven variables and points)
// Verified: secondary (Lee et al. 2015, PMC4444468; McConville & Cox 2023, PMC10211406)
// Tags:    rheumatology, gout, diagnostic rule, janssens, score
// Units:   each item present or absent; result points 0–13 (multiples of 0.5, exact in binary)
// Range:   any combination of the items (monoarthritis in primary care, the derivation population)
// Example: man, previous attack, MTP1, uric acid > 5.88 → 10; redness and onset within 1 day only → 1.5 (HAND-DERIVED from the points)
func AcuteGoutDiagnosticRuleScore(male, previousArthritisAttack, onsetWithinOneDay, jointRedness, firstMTPJoint, hypertensionOrCardiovascularDisease, serumUricAcidOver588MgDL bool) float64 {
	s := 0.0
	if male {
		s += 2
	}
	if previousArthritisAttack {
		s += 2
	}
	if onsetWithinOneDay {
		s += 0.5
	}
	if jointRedness {
		s++
	}
	if firstMTPJoint {
		s += 2.5
	}
	if hypertensionOrCardiovascularDisease {
		s += 1.5
	}
	if serumUricAcidOver588MgDL {
		s += 3.5
	}
	return s
}
