package series

// The reference's series kernel also computed statistics that calcbench already has. Those were REUSED,
// not rewritten (the ledger rows say which function replaces which); these tests run the reference's own
// cases through the calcbench functions, so the reuse claims are checked, not asserted.

import (
	"math"
	"testing"
	"time"

	"github.com/CrimsonCoder42/calcbench/algo/interval"
	"github.com/CrimsonCoder42/calcbench/num/stats"
)

func near(t *testing.T, name string, got, want, tol float64) {
	t.Helper()
	if math.Abs(got-want) > tol {
		t.Errorf("%s = %v, want %v (±%v)", name, got, want, tol)
	}
}

// parity: ledger Median, Quantile, MAD, MomentsOf, source case 1 (resting pulse with one spike: median
// 63.5, unscaled MAD 1.5, type-7 Q1 62.25 and Q3 64.75, sample SD 22.5 ±0.005 as printed).
func TestReuse_robustParity(t *testing.T) {
	v := []float64{62, 64, 63, 61, 65, 118}
	med, err := stats.Median(v)
	if err != nil {
		t.Fatal(err)
	}
	near(t, "median", med, 63.5, 1e-12)
	mad, _ := stats.MedianAbsoluteDeviation(v)
	near(t, "mad", mad, 1.5, 1e-12)
	q1, _ := stats.Percentile(v, 0.25, stats.HyndmanFan7)
	q3, _ := stats.Percentile(v, 0.75, stats.HyndmanFan7)
	near(t, "q1", q1, 62.25, 1e-12)
	near(t, "q3", q3, 64.75, 1e-12)
	sd, _ := stats.SampleStandardDeviation(v)
	near(t, "sample sd", sd, 22.5, 0.005)
}

// parity: ledger MomentsOf, source case 2 (sample variance 30 on 1e9 + {4, 7, 13, 16}, where the one-pass
// formula fails).
func TestReuse_varianceParity(t *testing.T) {
	v, err := stats.SampleVariance([]float64{1e9 + 4, 1e9 + 7, 1e9 + 13, 1e9 + 16})
	if err != nil {
		t.Fatal(err)
	}
	near(t, "variance", v, 30, 1e-6)
}

// parity: ledger OLSPerDay, source case 1 (least-squares slope 0.2963/day on the irregular readings, x in
// elapsed days).
func TestReuse_olsPerDayParity(t *testing.T) {
	x := []float64{0, 0.25, 0.5, 1, 1.25, 1.5, 7, 13}
	y := []float64{93, 94, 93, 94, 95, 95, 97, 97}
	_, slope, err := stats.LinearRegression(x, y)
	if err != nil {
		t.Fatal(err)
	}
	near(t, "ols", slope, 0.2963, 0.00005)
}

// parity: ledger TimeInRange, source case 1 (caller edges 70 and just above 180: counts 0, 9, 1). The
// reference's open-ended bands are NUM-STATS-008 with the edges padded by ±MaxFloat64.
func TestReuse_bandCountsParity(t *testing.T) {
	g := []float64{110, 180, 95, 150, 152, 90, 200, 120, 130, 85}
	edges := []float64{-math.MaxFloat64, 70, math.Nextafter(180, math.Inf(1)), math.MaxFloat64}
	c, err := stats.FrequencyDistribution(g, edges)
	if err != nil {
		t.Fatal(err)
	}
	if c[0] != 0 || c[1] != 9 || c[2] != 1 {
		t.Errorf("counts %v, want [0 9 1]", c)
	}
	sd, _ := stats.PopulationStandardDeviation(g)
	near(t, "population sd", sd, 36.9, 0.05)
}

// parity: ledger MergeIntervals, source case 1 (overlap merges: longest 420 min; with a 30-min gap one
// block of 600 min). The gap form is ALG-INTVL-001 on [start, end + gap] in Unix nanoseconds, then
// end − gap.
func TestReuse_mergeWithGapParity(t *testing.T) {
	base := time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC)
	m := func(a int) int64 { return base.Add(time.Duration(a) * time.Minute).UnixNano() }
	merge := func(gap time.Duration) []interval.Interval[int64] {
		in := []interval.Interval[int64]{
			{Start: m(1140), End: m(1380)}, {Start: m(1145), End: m(1560)}, {Start: m(1580), End: m(1740)}}
		for i := range in {
			in[i].End += int64(gap)
		}
		out, err := interval.MergeOverlapping(in)
		if err != nil {
			t.Fatal(err)
		}
		for i := range out {
			out[i].End -= int64(gap)
		}
		return out
	}
	out := merge(0)
	if len(out) != 2 || time.Duration(out[0].End-out[0].Start) != 420*time.Minute {
		t.Errorf("gap 0: %v", out)
	}
	out = merge(30 * time.Minute)
	if len(out) != 1 || time.Duration(out[0].End-out[0].Start) != 600*time.Minute {
		t.Errorf("gap 30: %v", out)
	}
}
