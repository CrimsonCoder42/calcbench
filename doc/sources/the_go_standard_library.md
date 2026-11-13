# The Go standard library

Citation (as in the doc comment): The Go standard library, package time, func Date ("normalizes" out-of-range day values; proleptic Gregorian calendar; the zone offset is the one in effect in loc at that civil time), https://pkg.go.dev/time#Date (read 2026-11-12)

Used by: NUM-SERIES-004

How checked: read in the primary as the builder's ledger row records (the S14 port ledger (LEDGER.md in the S14 work folder listed in doc/README.md), the row of each ID above); tests reproduce its values; reviewer sample 2026-11-12.
