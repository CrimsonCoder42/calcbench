package invest

import "fmt"

// BidAskSpread returns the relative quoted spread, the quoted spread (ask − bid) scaled by the
// bid-ask midpoint (ask + bid)/2, as a decimal. A locked or crossed quote (ask ≤ bid) is rejected, as
// the source excludes it. Computed as 2(ask − bid)/(ask + bid) from exact int64 sums.
// Errors: ErrOutOfRange if a price is above 2⁵³ cents; ErrNegative if a price is negative; ErrDomain
// if askCents ≤ bidCents.
//
// ID:      FIN-INVEST-030
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Collver, "A characterization of market quality for small capitalization US equities", US Securities and Exchange Commission, Division of Trading and Markets staff paper, September 2014, p. 2, https://www.sec.gov/marketstructure/research/small_cap_liquidity.pdf
// Tags:    invest, market, liquidity, bid-ask-spread, quoted-spread, relative-spread, trading-cost
// Units:   bidCents and askCents in cents; result decimal (fraction of the midpoint)
// Example: hand-checked: bid $24.19, ask $24.21 → 2·2/4840 = 0.000826446…
func BidAskSpread(bidCents, askCents int64) (float64, error) {
	if err := checkNonNegativeCents("bidCents", bidCents); err != nil {
		return 0, err
	}
	if err := checkNonNegativeCents("askCents", askCents); err != nil {
		return 0, err
	}
	if askCents <= bidCents {
		return 0, fmt.Errorf("%w: askCents = %d ≤ bidCents = %d (locked or crossed quote)", ErrDomain, askCents, bidCents)
	}
	return 2 * float64(askCents-bidCents) / float64(askCents+bidCents), nil
}
