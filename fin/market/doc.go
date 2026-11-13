// Package market holds market-microstructure and execution measures: the mid and weighted mid,
// the quoted spread, order-book and order-flow imbalance, trade signing (tick rule, Lee–Ready, bulk
// volume classification), VPIN, VWAP, implementation shortfall, Kyle's λ, the Roll spread, Amihud
// illiquidity and the Almgren–Chriss execution schedule and its cost.
//
// ID short form: MKT (IDs are FIN-MKT-NNN).
//
// Prices are int64 ticks (0 ≤ ticks ≤ 2⁵³) with tickSize, the price of one tick, a parameter; a price
// result is ticks × tickSize. Sizes, volumes and statistics are float64. Returns are fractions
// (0.01 = 1 %). No constant (tick size, bucket count, risk aversion…) is hidden: each is a parameter.
// Errors wrap the package sentinels (ErrDomain, ErrNegative, …); NaN or ±Inf in is an error, never a
// silent NaN out. Sums are Neumaier-compensated (algo/reduce). Package market imports no other fin
// package; the relative spread is fin/invest BidAskSpread (FIN-INVEST-030).
//
// Functions:
//
//	FIN-MKT-001  Mid                         (bid + ask)/2 in price units
//	FIN-MKT-002  WeightedMid                 Stoikov's imbalance-weighted mid I·ask + (1 − I)·bid
//	FIN-MKT-003  SpreadTicks                 quoted spread ask − bid in ticks
//	FIN-MKT-004  Level1Imbalance             queue imbalance Q^b/(Q^b + Q^a)
//	FIN-MKT-005  OrderFlowImbalance          Cont–Kukanov–Stoikov OFI = Σ e_n
//	FIN-MKT-006  TickRuleSign                trade signs by the tick test
//	FIN-MKT-007  LeeReadySign                trade signs by the Lee–Ready quote-then-tick rule
//	FIN-MKT-008  BulkVolumeClassification    buy/sell split V·Z(ΔP/σ), normal CDF
//	FIN-MKT-009  BulkVolumeClassificationT   buy/sell split V·t_ν(ΔP/σ), Student-t CDF
//	FIN-MKT-010  VPIN                        Σ|V^S − V^B|/(nV) over volume buckets
//	FIN-MKT-011  VWAP                        Σ p·v / Σ v in price units
//	FIN-MKT-012  ImplementationShortfall     Perold's execution and opportunity cost
//	FIN-MKT-013  KyleLambda                  slope of price change on signed volume
//	FIN-MKT-014  RollSpread                  2√(−cov(Δpₜ, Δpₜ₋₁))
//	FIN-MKT-015  AmihudIlliquidity           mean |R|/dollar volume
//	FIN-MKT-016  AlmgrenChrissTrajectory     optimal holdings sinh(κ(T − t))/sinh(κT)·X
//	FIN-MKT-017  AlmgrenChrissExpectedCost   E and V of a schedule's cost, linear impact
package market
