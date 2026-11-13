# Consolidation (W2): private helpers replaced by shared code

## Work list (lead's inventory, 2026-11-13, before W2)

These are the named private helpers that define a summation, a compensated sum, a NaN/Inf check or a gcd (`file:line: signature`).
Files with an inline compensated sum (`comp += (...)`) are listed after them. Area owners: C1 algo/** · C2 num/** · C3 fin/**, life/**, unit/** · C4 med/**.
Inline `math.IsNaN || math.IsInf` expressions are not helpers and are out of scope (BRIEF R5).

```
algo/tree/segment_tree.go:129:func checkFinite[T number](v T) error 
algo/geometry2d/max_points_on_line.go:62:func binaryGCD(u, v uint64) uint64 
algo/stream/moving_average_stream.go:134:func neumaierAdd(sum, comp, x float64) (float64, float64) 
algo/graph/graph.go:286:func checkedSum(a, b int64) (int64, bool) 
algo/intmath/happy_number.go:4:func digitSquareSum(n uint64) uint64 
algo/window/rolling_sum.go:86:func neumaierAdd(sum, comp, x float64) (float64, float64) 
num/series/trailing_time_window_sum.go:72:func neumaierAdd(sum, comp, x float64) (float64, float64) 
num/series/theil_sen_slope.go:97:func isFinite(v float64) bool 
num/trig/deg_to_rad.go:20:func isFinite(x float64) bool 
num/trig/law_of_sines_angle.go:87:func kahanBS(a, b, angleA float64) float64 
num/trig/law_of_sines_angle.go:95:func kahanBT(a, b, angleA float64) float64 
num/trig/law_of_sines_angle.go:111:func kahanBC(a, b, angleA float64) float64 
num/prob/probability_from_counts.go:26:func isFinite(x float64) bool 
num/prob/probability_from_counts.go:31:func checkFinite(name string, x float64) error 
num/linalg/vec_norm.go:38:func notFiniteAt(name string, i int) error 
num/linalg/mat_norms.go:92:func neumaierAdd(s, c, x float64) (float64, float64) 
num/linalg/mat.go:28:func firstNotFinite(a []float64) int 
num/linalg/vec2.go:37:func isFinite(x float64) bool 
num/linalg/vec2.go:41:func errNotFinite(name string) error 
num/combinatorics/binomial.go:60:func gcd(a, b uint64) uint64 
num/linalg/vec_dot.go:10:func twoSum(a, b float64) (x, y float64) 
num/integer/modular_arithmetic.go:38:func gcd(a, b uint64) uint64 
num/geometry/point_distance.go:22:func isFinite(x float64) bool 
num/geometry/point_distance.go:27:func checkFinite(x float64, name string) error 
num/geometry/polygon_area.go:10:func neumaierAdd(sum, comp, term float64) (newSum, newComp float64) 
num/stats/mean_absolute_deviation.go:34:func absoluteDeviationSum(values []float64, mean, scale float64) float64 
num/stats/compensated_sum.go:34:func isFinite(x float64) bool 
num/stats/compensated_sum.go:39:func checkFinite(name string, values []float64) error 
num/algebra/polynomial_roots.go:245:func compensatedHorner(poly []float64, z complex128, reversed bool) (p, dp complex128, bound float64) 
num/algebra/polynomial_roots.go:296:func twoSum(a, b float64) (s, e float64) 
num/algebra/polynomial_evaluate.go:27:func isFinite(v float64) bool 
fin/interest/future_value.go:64:func checkFinite(name string, v float64) error 
fin/economics/consumer_price_index.go:44:func checkFinite(name string, v float64) error 
fin/invest/current_yield.go:68:func checkFinite(name string, v float64) error 
life/fitness/load_progression_acsm.go:19:func checkFinite(name string, v float64) error 
life/nutrition/resting_energy_expenditure_mifflin_st_jeor.go:20:func isFinite(v float64) bool 
life/nutrition/resting_energy_expenditure_mifflin_st_jeor.go:23:func checkFinite(name string, v float64) error 
life/nutrition/resting_energy_expenditure_mifflin_st_jeor.go:31:func checkFinitePositive(name string, v float64) error 
med/hepatic/fib4.go:19:func isFinite(v float64) bool 
med/hepatic/fib4.go:22:func checkFinite(name string, v float64) error 
med/chemistry/corrected_calcium_payne.go:26:func isFinite(v float64) bool 
med/chemistry/corrected_calcium_payne.go:29:func checkFinite(name string, v float64) error 
med/cardio/max_heart_rate_tanaka.go:21:func checkFinite(name string, v float64) error 
med/renal/egfr_ckd_epi_2021.go:23:func isFinite(v float64) bool 
med/renal/egfr_ckd_epi_2021.go:26:func checkFinite(name string, v float64) error 
med/gastro/bisap_score.go:17:func isFinite(v float64) bool 
med/gastro/bisap_score.go:20:func checkFinite(name string, v float64) error 
med/endocrine/glucose_management_indicator.go:19:func isFinite(v float64) bool 
med/endocrine/glucose_management_indicator.go:22:func checkFinite(name string, v float64) error 
med/endocrine/glucose_management_indicator.go:38:func checkFinitePositive(name string, v float64) error 
med/anthropometry/body_mass_index.go:19:func isFinite(v float64) bool 
med/anthropometry/body_mass_index.go:22:func checkFinite(name string, v float64) error 
med/respiratory/diffusing_capacity_roughton_forster.go:21:func checkPositiveFinite(name string, v float64) error 
med/respiratory/peak_flow_zone_nhlbi.go:66:func checkFiniteNonNegative(name string, v float64) error 
med/thermoregulation/hypothermia_stage_swiss.go:17:func checkFinite(name string, v float64) error 
```

Not summation or finiteness helpers, so they stay (the lead's reading; a C agent may disagree with a reason):
- algo/graph checkedSum: int64 overflow-checked add.
- algo/intmath digitSquareSum: digits.
- num/trig kahanBS/BT/BC: Kahan's triangle formulas, not Kahan summation.
- num/algebra compensatedHorner: polynomial evaluation. Its twoSum is replaced.
- num/stats absoluteDeviationSum: a statistic. Its inner sum may call reduce.

Inline compensated sums: algo/stream/moving_average_stream.go · algo/window/rolling_sum.go · num/series/trailing_time_window_sum.go · num/prob/total_probability.go · num/linalg/mat_norms.go · num/geometry/polygon_area.go · num/stats/compensated_sum.go · life/fitness/total_volume_load.go · fin/economics/herfindahl_hirschman_index.go · med/respiratory/pef_diurnal_variability.go

Import-order limit: algo/** may not import num/**, so algo/geometry2d binaryGCD cannot call num/integer.GCD (D6 vs the import rule; C1 lists it, does not change it).

## Done (C agents append: file | private helper | replaced by ID | _delete copy)

| file | private helper | replaced by | _delete copy |
|---|---|---|---|
| med/hepatic/fib4.go | isFinite (removed; checkFinite wrapper kept, body calls AllFinite) | ALG-RED-010 AllFinite | _delete/med/hepatic/fib4.go |
| med/chemistry/corrected_calcium_payne.go | isFinite (removed; checkFinite wrapper kept, body calls AllFinite) | ALG-RED-010 AllFinite | _delete/med/chemistry/corrected_calcium_payne.go |
| med/renal/egfr_ckd_epi_2021.go | isFinite (removed; checkFinite wrapper kept, body calls AllFinite) | ALG-RED-010 AllFinite | _delete/med/renal/egfr_ckd_epi_2021.go |
| med/gastro/bisap_score.go | isFinite (removed; checkFinite wrapper kept, body calls AllFinite) | ALG-RED-010 AllFinite | _delete/med/gastro/bisap_score.go |
| med/endocrine/glucose_management_indicator.go | isFinite (removed; checkFinite wrapper kept, body calls AllFinite) | ALG-RED-010 AllFinite | _delete/med/endocrine/glucose_management_indicator.go |
| med/anthropometry/body_mass_index.go | isFinite (removed; checkFinite wrapper kept, body calls AllFinite) | ALG-RED-010 AllFinite | _delete/med/anthropometry/body_mass_index.go |
| med/hepatic/maddrey_discriminant_function.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/med/hepatic/maddrey_discriminant_function.go |
| med/renal/fractional_excretion_sodium.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/med/renal/fractional_excretion_sodium.go |
| med/endocrine/cgm_time_in_ranges.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/med/endocrine/cgm_time_in_ranges.go |
| med/endocrine/glucose_management_indicator_test.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/med/endocrine/glucose_management_indicator_test.go |
| med/cardio/max_heart_rate_tanaka.go | checkFinite predicate (wrapper kept) | ALG-RED-010 AllFinite | _delete/med/cardio/max_heart_rate_tanaka.go |
| med/thermoregulation/hypothermia_stage_swiss.go | checkFinite predicate (wrapper kept) | ALG-RED-010 AllFinite | _delete/med/thermoregulation/hypothermia_stage_swiss.go |
| med/critical/news2.go | finite predicate (wrapper kept) | ALG-RED-010 AllFinite | _delete/med/critical/news2.go |
| med/trauma/abc_score.go | finite predicate (wrapper kept) | ALG-RED-010 AllFinite | _delete/med/trauma/abc_score.go |
| med/neuro/abc2_hemorrhage_volume.go | checkRange NaN/Inf predicate (wrapper kept) | ALG-RED-010 AllFinite | _delete/med/neuro/abc2_hemorrhage_volume.go |
| med/pediatrics/kocher_criteria_count.go | checkFloat NaN/Inf predicate (wrapper kept) | ALG-RED-010 AllFinite | _delete/med/pediatrics/kocher_criteria_count.go |
| med/hematology/nccn_ipi_score.go | checkFloat NaN/Inf predicate (wrapper kept) | ALG-RED-010 AllFinite | _delete/med/hematology/nccn_ipi_score.go |
| med/respiratory/diffusing_capacity_roughton_forster.go | checkPositiveFinite NaN/Inf predicate (wrapper kept) | ALG-RED-010 AllFinite | _delete/med/respiratory/diffusing_capacity_roughton_forster.go |
| med/respiratory/six_minute_walk_distance_enright.go | checkClosedRange NaN/Inf predicate (wrapper kept; also serves checkFiniteNonNegative) | ALG-RED-010 AllFinite | _delete/med/respiratory/six_minute_walk_distance_enright.go |
| med/respiratory/pef_diurnal_variability.go | inline Neumaier sum (identical algorithm and order) | ALG-RED-002 NeumaierAdd | _delete/med/respiratory/pef_diurnal_variability.go |
| fin/interest/future_value.go | checkRate, checkFinite NaN/Inf predicates (wrappers kept) | ALG-RED-010 AllFinite | _delete/fin/interest/future_value.go |
| fin/interest/single_payment_compound_amount.go | factorResult NaN/Inf predicate (wrapper kept) | ALG-RED-010 AllFinite | _delete/fin/interest/single_payment_compound_amount.go |
| fin/interest/net_present_value.go | npvFloat private Horner loop (now FMA Horner via a 64-slot stack chunk; d7_equivalence_test.go) | NUM-ALGEB-001 PolynomialEvaluate | _delete/fin/interest/net_present_value.go |
| fin/loan/amortization_schedule.go | checkRate NaN/Inf predicate (wrapper kept) | ALG-RED-010 AllFinite | _delete/fin/loan/amortization_schedule.go |
| fin/economics/consumer_price_index.go | checkFinite predicate (wrapper kept) | ALG-RED-010 AllFinite | _delete/fin/economics/consumer_price_index.go |
| fin/economics/inflation_rate.go | checkResult NaN/Inf predicate (wrapper kept) | ALG-RED-010 AllFinite | _delete/fin/economics/inflation_rate.go |
| fin/economics/real_interest_rate.go | checkRate NaN/Inf predicate (wrapper kept) | ALG-RED-010 AllFinite | _delete/fin/economics/real_interest_rate.go |
| fin/economics/herfindahl_hirschman_index.go | inline Neumaier sum (identical algorithm and order) | ALG-RED-002 NeumaierAdd | _delete/fin/economics/herfindahl_hirschman_index.go |
| fin/business/free_cash_flow_to_firm.go | checkTaxRate NaN/Inf predicate (wrapper kept) | ALG-RED-010 AllFinite | _delete/fin/business/free_cash_flow_to_firm.go |
| fin/invest/current_yield.go | checkRate, checkFinite, finiteResult NaN/Inf predicates (wrappers kept) | ALG-RED-010 AllFinite | _delete/fin/invest/current_yield.go |
| fin/invest/black_scholes_call.go | normalCDF body (wrapper kept for the ±Inf limits; bit-identical, d7_equivalence_test.go) | NUM-PROB-048 NormalCDF | _delete/fin/invest/black_scholes_call.go |
| life/fitness/load_progression_acsm.go | checkFinite, checkResult NaN/Inf predicates (wrappers kept) | ALG-RED-010 AllFinite | _delete/life/fitness/load_progression_acsm.go |
| life/fitness/total_volume_load.go | inline Neumaier sum (identical algorithm and order) | ALG-RED-002 NeumaierAdd | _delete/life/fitness/total_volume_load.go |
| life/nutrition/resting_energy_expenditure_mifflin_st_jeor.go | isFinite (removed; checkFinite wrapper kept, body calls AllFinite) | ALG-RED-010 AllFinite | _delete/life/nutrition/resting_energy_expenditure_mifflin_st_jeor.go |
| life/nutrition/percent_daily_value.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/life/nutrition/percent_daily_value.go |
| life/nutrition/amdr_grams.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/life/nutrition/amdr_grams.go |
| life/nutrition/protein_range_issn.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/life/nutrition/protein_range_issn.go |
| unit/length/centimeters_to_inches.go | checkLength NaN/Inf predicate (wrapper kept) | ALG-RED-010 AllFinite | _delete/unit/length/centimeters_to_inches.go |
| unit/mass/kilograms_to_pounds.go | checkMass NaN/Inf predicate (wrapper kept) | ALG-RED-010 AllFinite | _delete/unit/mass/kilograms_to_pounds.go |
| fin/interest/internal_rate_of_return.go, fin/invest/{modified_internal_rate_of_return,profitability_index,macaulay_duration}.go | private Horner loops | NOT CHANGED (C3): see C3 hand-back | — |
| algo/stream/moving_average_stream.go | neumaierAdd (removed; identical branch condition and order) | ALG-RED-002 NeumaierAdd | _delete/algo/stream/moving_average_stream.go |
| algo/window/rolling_sum.go | neumaierAdd (removed; identical branch condition and order) | ALG-RED-002 NeumaierAdd | _delete/algo/window/rolling_sum.go |
| algo/tree/segment_tree.go | checkFinite[T number] (not done: generic over ints and float32/float64 via finite[T]; reduce.AllFinite takes float64 only, a conversion would add a type switch per update) | none | none |
| algo/geometry2d/max_points_on_line.go | binaryGCD (not done: import order; algo may not import num/integer.GCD) | none | none |
| num/algebra/cubic_roots.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/algebra/cubic_roots.go |
| num/algebra/floor_div_mod.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/algebra/floor_div_mod.go |
| num/algebra/linear_evaluate.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/algebra/linear_evaluate.go |
| num/algebra/linear_solve.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/algebra/linear_solve.go |
| num/algebra/polynomial_add.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/algebra/polynomial_add.go |
| num/algebra/polynomial_companion.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/algebra/polynomial_companion.go |
| num/algebra/polynomial_derivative.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/algebra/polynomial_derivative.go |
| num/algebra/polynomial_divide.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/algebra/polynomial_divide.go |
| num/algebra/polynomial_evaluate.go | isFinite (removed; identical predicate/kernel) | ALG-RED-010 AllFinite | _delete/num/algebra/polynomial_evaluate.go |
| num/algebra/polynomial_evaluate_2d.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/algebra/polynomial_evaluate_2d.go |
| num/algebra/polynomial_evaluate_grid_2d.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/algebra/polynomial_evaluate_grid_2d.go |
| num/algebra/polynomial_from_roots.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/algebra/polynomial_from_roots.go |
| num/algebra/polynomial_integral.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/algebra/polynomial_integral.go |
| num/algebra/polynomial_integral_order.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/algebra/polynomial_integral_order.go |
| num/algebra/polynomial_multiply.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/algebra/polynomial_multiply.go |
| num/algebra/polynomial_roots.go | twoSum (removed; same operations as reduce.TwoSum) | ALG-RED-003 TwoSum | _delete/num/algebra/polynomial_roots.go |
| num/algebra/polynomial_trim.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/algebra/polynomial_trim.go |
| num/algebra/polynomial_value_and_derivative.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/algebra/polynomial_value_and_derivative.go |
| num/algebra/quadratic_discriminant.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/algebra/quadratic_discriminant.go |
| num/algebra/quadratic_discriminant_test.go | caller of package isFinite (test: call only) | ALG-RED-010 AllFinite | _delete/num/algebra/quadratic_discriminant_test.go |
| num/algebra/quadratic_real_roots.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/algebra/quadratic_real_roots.go |
| num/algebra/quadratic_roots.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/algebra/quadratic_roots.go |
| num/algebra/round_decimals.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/algebra/round_decimals.go |
| num/algebra/sign.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/algebra/sign.go |
| num/calculus/gradient.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/calculus/gradient.go |
| num/calculus/interp.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/calculus/interp.go |
| num/calculus/trapezoid.go | isFinite (removed; checkFinite wrapper kept, body calls AllFinite); neumaierSum.add body (struct kept; identical branch and order) | ALG-RED-010 AllFinite, ALG-RED-002 NeumaierAdd | _delete/num/calculus/trapezoid.go |
| num/combinatorics/binomial.go | gcd (removed; same Algorithm B; combinatorics may import integer per D-009) | NUM-INTEG-004 GCD | _delete/num/combinatorics/binomial.go |
| num/geometry/ellipse_perimeter.go | caller of package neumaierAdd | ALG-RED-002 NeumaierAdd | _delete/num/geometry/ellipse_perimeter.go |
| num/geometry/ellipse_perimeter_test.go | caller of package neumaierAdd (test: call only) | ALG-RED-002 NeumaierAdd | _delete/num/geometry/ellipse_perimeter_test.go |
| num/geometry/point_distance.go | isFinite (removed; identical predicate/kernel) | ALG-RED-010 AllFinite | _delete/num/geometry/point_distance.go |
| num/geometry/polygon_area.go | neumaierAdd (removed; identical predicate/kernel) | ALG-RED-002 NeumaierAdd, ALG-RED-010 AllFinite | _delete/num/geometry/polygon_area.go |
| num/geometry/polygon_perimeter.go | caller of package neumaierAdd | ALG-RED-002 NeumaierAdd | _delete/num/geometry/polygon_perimeter.go |
| num/integer/modular_arithmetic.go | gcd (removed; same Algorithm B line for line) | NUM-INTEG-004 GCD | _delete/num/integer/modular_arithmetic.go |
| num/integer/modular_arithmetic_test.go | caller of package gcd (test: call only) | NUM-INTEG-004 GCD | _delete/num/integer/modular_arithmetic_test.go |
| num/integer/pollard_rho.go | caller of package gcd | NUM-INTEG-004 GCD | _delete/num/integer/pollard_rho.go |
| num/linalg/cholesky.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/linalg/cholesky.go |
| num/linalg/covariance.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/linalg/covariance.go |
| num/linalg/cylindrical.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/linalg/cylindrical.go |
| num/linalg/eigen.go | caller of package finiteN | ALG-RED-010 AllFinite | _delete/num/linalg/eigen.go |
| num/linalg/eigenvectors_general.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/linalg/eigenvectors_general.go |
| num/linalg/lu.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/linalg/lu.go |
| num/linalg/mat.go | isFinite caller; firstNotFinite kept (returns the index), predicate is AllFinite | ALG-RED-010 AllFinite | _delete/num/linalg/mat.go |
| num/linalg/mat2.go | finiteN (removed; callers use reduce.AllFinite(m[:]...)) | ALG-RED-010 AllFinite | _delete/num/linalg/mat2.go |
| num/linalg/mat3.go | caller of package isFinite, finiteN | ALG-RED-010 AllFinite | _delete/num/linalg/mat3.go |
| num/linalg/mat3_rotation.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/linalg/mat3_rotation.go |
| num/linalg/mat4.go | caller of package isFinite, finiteN | ALG-RED-010 AllFinite | _delete/num/linalg/mat4.go |
| num/linalg/mat4_transform.go | caller of package isFinite, finiteN | ALG-RED-010 AllFinite | _delete/num/linalg/mat4_transform.go |
| num/linalg/mat_characteristics.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/linalg/mat_characteristics.go |
| num/linalg/mat_cofactor.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/linalg/mat_cofactor.go |
| num/linalg/mat_cofactor_fixed.go | caller of package isFinite, finiteN | ALG-RED-010 AllFinite | _delete/num/linalg/mat_cofactor_fixed.go |
| num/linalg/mat_condition.go | caller of package isFinite, twoSum | ALG-RED-003 TwoSum, ALG-RED-010 AllFinite | _delete/num/linalg/mat_condition.go |
| num/linalg/mat_condition_exact.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/linalg/mat_condition_exact.go |
| num/linalg/mat_elementwise.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/linalg/mat_elementwise.go |
| num/linalg/mat_log_det.go | caller of package neumaierAdd | ALG-RED-002 NeumaierAdd | _delete/num/linalg/mat_log_det.go |
| num/linalg/mat_norms.go | neumaierAdd (removed; identical predicate/kernel) | ALG-RED-002 NeumaierAdd, ALG-RED-010 AllFinite | _delete/num/linalg/mat_norms.go |
| num/linalg/mat_norms_extra.go | caller of package isFinite, neumaierAdd | ALG-RED-002 NeumaierAdd, ALG-RED-010 AllFinite | _delete/num/linalg/mat_norms_extra.go |
| num/linalg/mat_rank.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/linalg/mat_rank.go |
| num/linalg/polynomial_fit.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/linalg/polynomial_fit.go |
| num/linalg/polynomial_fit_covariance.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/linalg/polynomial_fit_covariance.go |
| num/linalg/quat.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/linalg/quat.go |
| num/linalg/quat_axis_angle.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/linalg/quat_axis_angle.go |
| num/linalg/quat_euler.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/linalg/quat_euler.go |
| num/linalg/quat_matrix.go | caller of package finiteN | ALG-RED-010 AllFinite | _delete/num/linalg/quat_matrix.go |
| num/linalg/quat_slerp.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/linalg/quat_slerp.go |
| num/linalg/spherical.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/linalg/spherical.go |
| num/linalg/svd.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/linalg/svd.go |
| num/linalg/symmetric_eigen.go | caller of package isFinite, finiteN | ALG-RED-010 AllFinite | _delete/num/linalg/symmetric_eigen.go |
| num/linalg/triangular.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/linalg/triangular.go |
| num/linalg/vandermonde.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/linalg/vandermonde.go |
| num/linalg/vec.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/linalg/vec.go |
| num/linalg/vec2.go | isFinite (removed); finite2 wrapper kept, body reduce.AllFinite(v[:]...) | ALG-RED-010 AllFinite | _delete/num/linalg/vec2.go |
| num/linalg/vec3.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/linalg/vec3.go |
| num/linalg/vec3_area.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/linalg/vec3_area.go |
| num/linalg/vec3_norm.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/linalg/vec3_norm.go |
| num/linalg/vec3_product.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/linalg/vec3_product.go |
| num/linalg/vec3_projection.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/linalg/vec3_projection.go |
| num/linalg/vec3_rotate.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/linalg/vec3_rotate.go |
| num/linalg/vec_angle.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/linalg/vec_angle.go |
| num/linalg/vec_dot.go | twoSum (removed; same operations as reduce.TwoSum, error-free so exact) | ALG-RED-003 TwoSum | _delete/num/linalg/vec_dot.go |
| num/linalg/vec_elementwise.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/linalg/vec_elementwise.go |
| num/linalg/vec_norm.go | caller of package isFinite, twoSum | ALG-RED-003 TwoSum, ALG-RED-010 AllFinite | _delete/num/linalg/vec_norm.go |
| num/series/average_rate_of_change_per_day.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/series/average_rate_of_change_per_day.go |
| num/series/average_real_variability.go | caller of package isFinite, neumaierAdd | ALG-RED-002 NeumaierAdd, ALG-RED-010 AllFinite | _delete/num/series/average_real_variability.go |
| num/series/fixed_stride_buckets.go | caller of package isFinite, neumaierAdd | ALG-RED-002 NeumaierAdd, ALG-RED-010 AllFinite | _delete/num/series/fixed_stride_buckets.go |
| num/series/local_hour_window_mean.go | caller of package isFinite, neumaierAdd | ALG-RED-002 NeumaierAdd, ALG-RED-010 AllFinite | _delete/num/series/local_hour_window_mean.go |
| num/series/theil_sen_slope.go | isFinite (removed; identical predicate/kernel) | ALG-RED-010 AllFinite | _delete/num/series/theil_sen_slope.go |
| num/series/trailing_time_window_sum.go | neumaierAdd (removed; identical predicate/kernel) | ALG-RED-002 NeumaierAdd, ALG-RED-010 AllFinite | _delete/num/series/trailing_time_window_sum.go |
| num/series/weekly_sums.go | caller of package isFinite, neumaierAdd | ALG-RED-002 NeumaierAdd, ALG-RED-010 AllFinite | _delete/num/series/weekly_sums.go |
| num/stats/chi_square_goodness_of_fit.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/stats/chi_square_goodness_of_fit.go |
| num/stats/chi_square_goodness_of_fit_uniform.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/stats/chi_square_goodness_of_fit_uniform.go |
| num/stats/coefficient_of_variation.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/stats/coefficient_of_variation.go |
| num/stats/compensated_sum.go | isFinite (removed; checkFinite/checkScalar wrappers kept); compensatedSum.add body (struct kept; scaling, two-pass and floor logic in meanOf/correctedSquares/correctedProducts kept, only the add step calls reduce) | ALG-RED-010 AllFinite, ALG-RED-002 NeumaierAdd | _delete/num/stats/compensated_sum.go |
| num/stats/covariance.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/stats/covariance.go |
| num/stats/frequency_density.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/stats/frequency_density.go |
| num/stats/grouped_mean.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/stats/grouped_mean.go |
| num/stats/linear_regression.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/stats/linear_regression.go |
| num/stats/mean_absolute_deviation.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/stats/mean_absolute_deviation.go |
| num/stats/median_absolute_deviation.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/stats/median_absolute_deviation.go |
| num/stats/percent_relative_standard_deviation.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/stats/percent_relative_standard_deviation.go |
| num/stats/range.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/stats/range.go |
| num/stats/raw_score_from_z_score.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/stats/raw_score_from_z_score.go |
| num/stats/sum_of_squares.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/stats/sum_of_squares.go |
| num/stats/summary_statistics.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/stats/summary_statistics.go |
| num/stats/weighted_mean.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/stats/weighted_mean.go |
| num/stats/z_score.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/stats/z_score.go |
| num/trig/acot.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/trig/acot.go |
| num/trig/acoth.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/trig/acoth.go |
| num/trig/acsc.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/trig/acsc.go |
| num/trig/acsch.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/trig/acsch.go |
| num/trig/angle_difference.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/trig/angle_difference.go |
| num/trig/asec.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/trig/asec.go |
| num/trig/asech.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/trig/asech.go |
| num/trig/central_angle.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/trig/central_angle.go |
| num/trig/cot.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/trig/cot.go |
| num/trig/coth.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/trig/coth.go |
| num/trig/csc.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/trig/csc.go |
| num/trig/csch.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/trig/csch.go |
| num/trig/deg_to_rad.go | isFinite (removed; identical predicate/kernel) | ALG-RED-010 AllFinite | _delete/num/trig/deg_to_rad.go |
| num/trig/law_of_sines_side.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/trig/law_of_sines_side.go |
| num/trig/normalize_angle.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/trig/normalize_angle.go |
| num/trig/normalize_angle_signed.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/trig/normalize_angle_signed.go |
| num/trig/rad_to_deg.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/trig/rad_to_deg.go |
| num/trig/sec.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/trig/sec.go |
| num/trig/sech.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/trig/sech.go |
| num/trig/unwrap_angles.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/trig/unwrap_angles.go |
| num/logs/log_sum_exp.go | inline Neumaier sum (sum >= x on non-negative terms = |sum| >= |x|; identical) | ALG-RED-002 NeumaierAdd | _delete/num/logs/log_sum_exp.go |
| num/signal/fft.go | cmplxFinite predicate (wrapper kept) | ALG-RED-010 AllFinite | _delete/num/signal/fft.go |
| num/trig kahanBS/BT/BC, num/algebra compensatedHorner, num/linalg dot2/twoProduct, Kahan 2×2 determinants | not done: not summation helpers (formulas / FMA error-free products); only their twoSum step now calls reduce.TwoSum | none | none |
| num/special, num/float, num/logs inline math.IsNaN||math.IsInf checks | not done: inline expressions, not named helpers (R5, LATER) | none | none |
| num/signal/convolve.go firstNotFinite | already calls reduce.AllFinite (index/error wrapper, R5) | none | none |
| num/prob/chi_square_cdf.go regularizedGamma/gammaPrefactor (used by ChiSquareCDF, ChiSquareSurvival, ChiSquareQuantile) | SUPERSEDED (see the B2 pass-4 row below: replaced after B5b moved Loader's prefactor into num/special) — earlier note: NOT replaced (B2, 2026-11-13): the 1e-14 gate cannot hold. special.gammaPrefactor (regularized_gamma.go:109) forms exp(a·ln x − x − lgamma(a)) directly, whose rounding grows like ε·a·|ln x| (its own doc says it loses about log10(a) digits); the private one uses Loader's saddle-point form above a = 15 (chi_square_cdf.go:31), and ChiSquareCDF accepts df ≤ 10⁶ (a = 5·10⁵). Recommendation: B5 moves Loader's prefactor into num/special, then the swap is a straight replacement | none (file unchanged) | none |
| num/prob/beta_pdf.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/prob/beta_pdf.go |
| num/prob/binomial_pmf.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/prob/binomial_pmf.go |
| num/prob/cauchy_pdf.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/prob/cauchy_pdf.go |
| num/prob/cauchy_quantile.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/prob/cauchy_quantile.go |
| num/prob/chi_square_mean_variance.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/prob/chi_square_mean_variance.go |
| num/prob/decimal_odds_from_probability.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/prob/decimal_odds_from_probability.go |
| num/prob/discrete_variance.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/prob/discrete_variance.go |
| num/prob/expected_value.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/prob/expected_value.go |
| num/prob/exponential_mean_variance.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/prob/exponential_mean_variance.go |
| num/prob/exponential_quantile.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/prob/exponential_quantile.go |
| num/prob/f_distribution_mean_variance.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/prob/f_distribution_mean_variance.go |
| num/prob/f_distribution_pdf.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/prob/f_distribution_pdf.go |
| num/prob/gamma_mean_variance.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/prob/gamma_mean_variance.go |
| num/prob/gamma_pdf.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/prob/gamma_pdf.go |
| num/prob/geometric_trials_mean_variance.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/prob/geometric_trials_mean_variance.go |
| num/prob/gumbel_mean_variance.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/prob/gumbel_mean_variance.go |
| num/prob/gumbel_pdf.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/prob/gumbel_pdf.go |
| num/prob/gumbel_quantile.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/prob/gumbel_quantile.go |
| num/prob/laplace_mean_variance.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/prob/laplace_mean_variance.go |
| num/prob/laplace_pdf.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/prob/laplace_pdf.go |
| num/prob/laplace_quantile.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/prob/laplace_quantile.go |
| num/prob/logistic_mean_variance.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/prob/logistic_mean_variance.go |
| num/prob/logistic_pdf.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/prob/logistic_pdf.go |
| num/prob/logistic_quantile.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/prob/logistic_quantile.go |
| num/prob/lognormal_mean_variance.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/prob/lognormal_mean_variance.go |
| num/prob/lognormal_pdf.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/prob/lognormal_pdf.go |
| num/prob/lognormal_quantile.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/prob/lognormal_quantile.go |
| num/prob/multivariate_normal_pdf.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/prob/multivariate_normal_pdf.go |
| num/prob/negative_binomial_failures_mean_variance.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/prob/negative_binomial_failures_mean_variance.go |
| num/prob/negative_binomial_trials_mean_variance.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/prob/negative_binomial_trials_mean_variance.go |
| num/prob/normal_pdf.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/prob/normal_pdf.go |
| num/prob/normal_quantile.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/prob/normal_quantile.go |
| num/prob/overround.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/prob/overround.go |
| num/prob/pareto_mean_variance.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/prob/pareto_mean_variance.go |
| num/prob/pareto_pdf.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/prob/pareto_pdf.go |
| num/prob/pareto_quantile.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/prob/pareto_quantile.go |
| num/prob/power_function_pdf.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/prob/power_function_pdf.go |
| num/prob/probability_from_counts.go | isFinite (removed; identical predicate, |x| <= MaxFloat64); checkFinite/checkProbability/checkPositive kept as error wrappers | ALG-RED-010 AllFinite | _delete/num/prob/probability_from_counts.go |
| num/prob/rayleigh_mean_variance.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/prob/rayleigh_mean_variance.go |
| num/prob/rayleigh_pdf.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/prob/rayleigh_pdf.go |
| num/prob/rayleigh_quantile.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/prob/rayleigh_quantile.go |
| num/prob/student_t_pdf.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/prob/student_t_pdf.go |
| num/prob/total_probability.go | compensatedSum.add inline Neumaier step (same branch and order; identical) + caller of package isFinite; addProduct FMA correction kept | ALG-RED-002 NeumaierAdd, ALG-RED-010 AllFinite | _delete/num/prob/total_probability.go |
| num/prob/triangular_mean_variance.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/prob/triangular_mean_variance.go |
| num/prob/triangular_pdf.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/prob/triangular_pdf.go |
| num/prob/uniform_cdf.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/prob/uniform_cdf.go.2026-11-13 |
| num/prob/uniform_mean_variance.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/prob/uniform_mean_variance.go |
| num/prob/uniform_pdf.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/prob/uniform_pdf.go |
| num/prob/uniform_quantile.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/prob/uniform_quantile.go |
| num/prob/weibull_mean_variance.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/prob/weibull_mean_variance.go |
| num/prob/weibull_pdf.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/prob/weibull_pdf.go |
| num/prob/weibull_quantile.go | caller of package isFinite | ALG-RED-010 AllFinite | _delete/num/prob/weibull_quantile.go |
| num/prob other compensated sums (discrete_variance.go:39, binomial/poisson/hypergeometric CDF tails, overround, expected_value, bayes_posterior) | not changed: they go through compensatedSum.add, which now calls reduce.NeumaierAdd; discrete_variance's FMA rounding term on comp is an error-free product correction, kept | ALG-RED-002 NeumaierAdd (via add) | none |
| num/prob/chi_square_cdf.go (regularizedGamma, gammaPrefactor, gammaMaxTerms) + callers chi_square_survival.go, chi_square_quantile.go | private incomplete gamma retired after B5b LOADER READY; old code kept verbatim as test-only oracle in chi_square_cdf_private_gamma_test.go (TestChiSquareCDF_matchesPrivateGamma: 1e-14 relative in the bulk, wider documented allowance in a ≤ 15 far tails and √a·2⁻⁵² near the bulk for large df); quantile's Newton density now ChiSquarePDF | NUM-SPEC-001/002 RegularizedGammaP/Q | _delete/num/prob/chi_square_cdf.go, chi_square_survival.go, chi_square_quantile.go, doc.go |
