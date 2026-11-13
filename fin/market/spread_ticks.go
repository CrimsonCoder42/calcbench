package market

// SpreadTicks returns the quoted spread ask − bid in ticks. A locked quote gives 0.
// For the relative spread 2(ask − bid)/(ask + bid) use fin/invest BidAskSpread (FIN-INVEST-030): the
// ratio is unit-free, so ticks may be passed as its cents.
// Errors: ErrNegative if a price is negative; ErrOutOfRange if a price is above 2⁵³ ticks; ErrDomain if
// askTicks < bidTicks (crossed).
//
// ID:      FIN-MKT-003
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Collver, "A characterization of market quality for small capitalization US equities", US Securities and Exchange Commission, Division of Trading and Markets staff paper, September 2014, p. 2 (quoted spread = ask − bid), https://www.sec.gov/marketstructure/research/small_cap_liquidity.pdf
// Tags:    market, microstructure, spread, quoted-spread, ticks, liquidity
// Units:   bidTicks, askTicks and the result in ticks
// Example: hand-checked: bid 2419, ask 2421 → 2
func SpreadTicks(bidTicks, askTicks int64) (int64, error) {
	if err := checkQuote(bidTicks, askTicks); err != nil {
		return 0, err
	}
	return askTicks - bidTicks, nil
}
