# Report: HFT and execution maths, S19 (2026-11-13)
**Start check:** passed. S18's reviewer had finished, and build and vet were clean. This Mac has no `timeout` binary,
so every command was bounded by go's own -timeout.
**Builders** (all done in one pass): · **B1 fin/market (new), FIN-MKT-001…017:** Mid, WeightedMid, SpreadTicks,
Level1Imbalance, OrderFlowImbalance, TickRuleSign, LeeReadySign, BVC (normal and t), VPIN, VWAP,
ImplementationShortfall, KyleLambda, RollSpread, Amihud, Almgren–Chriss trajectory and cost. · **B2 fin/invest,
FIN-INVEST-032…052:** returns, annualisation, realized and EWMA variance, Parkinson, Garman–Klass, Rogers–Satchell,
Beta, TrackingError, IR, Sortino, MaxDrawdown, Calmar, historical and parametric VaR, ES, Kelly. 031 was skipped:
doc.go already reserved it. · **B3:** NUM-SERIES-012…019 (Autocorrelation, ACF, PACF, LjungBox, VarianceRatio,
HurstExponentRS, AR1Fit, MeanReversionHalfLife) and NUM-SIG-016 KalmanFilter1D. · **B4 algo/stream,
ALG-STREAM-005…016:** windowed variance, covariance and z-score; online covariance; the EWMA stream; streaming and
windowed VWAP; the tick and bp helpers.
**Tally:** BUILT 59 · REUSED 4 · NEEDS SOURCE 3. · REUSED: RelativeSpread → FIN-INVEST-030; SharpeRatio → 025;
one-pair return → 010/012; SpreadZScore → NUM-STATS-028. · NEEDS SOURCE: DepthWeightedImbalance, TWAP, and the
continuous Kelly (μ − r)/σ².
**Count:** 1,505 → 1,564 functions; 75 → 76 packages. INDEX == code == doc.go: yes (1564 | 1564 (+28 REMOVED) | 1564).
**Tests:** · `go test ./... -count=1 -timeout 15m`: 75 ok, 2 FAIL, and `-race` failed the same 2 with no races. Both
were code bugs: Almgren–Chriss x₀ ≠ X: an arm64 FMA fused a cancelling subtraction. Autocorrelation returned NaN on
subnormal input: its scale factor overflowed. · Owners fixed both. The rerun of fin/market, num/series, algo/stream
and testing: all ok, plain and -race.
**The reviewer's two numeric checks:** · **WindowedVariance vs a two-pass sample variance** (the NUM-STATS-021
formula), 1,000 fixed-seed values, every step: 1e-9 relative for w ≥ 7; for w ≤ 3, 1e-9 of the series variance. The
reviewer judged this a real limit of the sliding downdate (D-086). · **EWMAStream vs ALG-WINDOW-014:** 1,000 values,
five α, every step, within 1e-15. Passed.
**Other checks:** · No hidden constants: 252, 0.94 and 1.645 appear only as stated example parameters. · Every Add,
Update and Predict is O(1) by reading the code. · No existing signature changed. · Allocations match the documented
Space.
**FOUND, NOT FIXED:** · **Domain words in algo:** the tick and bp helpers and the two VWAP streams are in algo/stream,
per your owner list, but CLAUDE.md says algo has no domain words (D-086). They could move to fin/market with new IDs.
· **Unchecked citation locators:** the builders flagged locators they could not check against the documents: CLM
(1.4.x), RiskMetrics [5.3], BJRL sections, Hamilton [1.1.10], Amihud eq. (1), Stoikov and Cont–Kukanov–Stoikov
equations. Young 1991 "Futures 20(15):40" came from my brief and is unverified. · **RunningStats:** ALG-STREAM-004 has
no Reset (D3 shape). It is released, so I left it. · **RollSpread** allocates O(n) where an allocation-free pass
exists.
**DELETED:** nothing removed. Edited files were copied to `_delete/` first, with LEDGER lines.
**LEFT FOR DEVIN:** · whether to move the algo/stream domain helpers; · sources for the 3 NEEDS SOURCE rows; · the
citation check; · `/tmp/x`; · emptying `_delete/`.
