// Package shopping holds everyday purchase calculators: the tip on a bill.
//
// ID short form: SHOP (IDs are LIF-SHOP-NNN).
//
// Units: money is int64 cents in and out; a percentage rate is int64 basis points (1 bp = 0.01 %), so every
// result is exact integer arithmetic with one named rounding step. Overflow is checked (ErrOverflow); a
// negative amount or rate is ErrOutOfRange.
//
// Functions:
//
//	LIF-SHOP-001  TipCents  tip on a bill, bill × rate, rounded half to even to the cent
package shopping
