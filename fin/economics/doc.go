// Package economics holds macroeconomic and market-structure measures: price indexes, inflation, the real
// interest rate and market concentration.
//
// ID short form: ECON (IDs are FIN-ECON-NNN).
//
// Money is int64 cents; rates are float64 decimals per period. Index levels (CPI, HHI) are returned in the
// points their defining source uses (base 100; squared percent), not as decimals.
//
// Functions:
//
//	FIN-ECON-001  ConsumerPriceIndex        fixed-basket price index, base period = 100
//	FIN-ECON-002  InflationRate             decimal change between two index levels
//	FIN-ECON-003  RealInterestRate          exact Fisher relation (1+i)/(1+π) − 1
//	FIN-ECON-004  HerfindahlHirschmanIndex  sum of squared percent market shares, 0–10 000
//	FIN-ECON-005  GDPExpenditure            C + I + G + (X − M), cents
package economics
