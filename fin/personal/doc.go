// Package personal holds household and payroll conventions: pay-period conversions now, and later the
// retirement and pension planning rows (decision D-021).
//
// ID short form: PERSON (IDs are FIN-PERSON-NNN).
//
// Money is int64 cents; every money result is rounded half-even to the cent. Errors wrap the package
// sentinels with the offending parameter; test with errors.Is.
//
// Functions:
//
//	FIN-PERSON-001  WeeklyPayEquivalent        pay for a period × periods a year ÷ 52 (29 CFR 778.113(b))
//	FIN-PERSON-002  HourlyRateFromWeeklyPay    weekly pay ÷ hours it is meant to cover (29 CFR 778.113(a))
package personal
