package invest

import "testing"

// OpenStax Principles of Finance §16.1, Sam's Sporting Goods: −$16,000, then $2,000, $4,000, $5,000,
// $5,000, $5,000 → 4 years exactly.
func TestPaybackPeriod_openStax(t *testing.T) {
	got, err := PaybackPeriod([]int64{-1600000, 200000, 400000, 500000, 500000, 500000})
	if err != nil || got != 4 {
		t.Errorf("PaybackPeriod = %v, %v; want 4", got, err)
	}
}

func TestPaybackPeriod_cases(t *testing.T) {
	for _, c := range []struct {
		name  string
		flows []int64
		want  float64
	}{
		{"fraction", []int64{-1000, 400, 400, 400}, 2.5},
		{"one period", []int64{-1000, 1000}, 1},
		{"first flow positive", []int64{100, -50}, 0},
		{"first flow zero", []int64{0}, 0},
		{"later outflow ignored after recovery", []int64{-100, 200, -1000}, 0.5},
		{"dip then recover", []int64{-100, 50, -50, 200}, 2.5},
		{"extremes", []int64{-maxCents, maxCents}, 1},
	} {
		got, err := PaybackPeriod(c.flows)
		if err != nil || got != c.want {
			t.Errorf("%s: %v, %v; want %v", c.name, got, err, c.want)
		}
	}
}

func TestPaybackPeriod_errors(t *testing.T) {
	_, err := PaybackPeriod(nil)
	wantErr(t, "empty", err, ErrDomain)
	_, err = PaybackPeriod([]int64{-1000, 100, 100})
	wantErr(t, "never recovered", err, ErrNoSolution)
	_, err = PaybackPeriod([]int64{-1000, maxCents + 1})
	wantErr(t, "flow > 2^53", err, ErrOutOfRange)
	big := make([]int64, 1100)
	for i := range big {
		big[i] = -maxCents
	}
	_, err = PaybackPeriod(big) // 1100·2^53 > 2^63
	wantErr(t, "cumulative overflow", err, ErrOverflow)
}

func TestPaybackPeriod_noAlloc(t *testing.T) {
	flows := []int64{-1600000, 200000, 400000, 500000, 500000, 500000}
	if a := testing.AllocsPerRun(100, func() { sinkFloat, sinkErr = PaybackPeriod(flows) }); a != 0 {
		t.Errorf("allocs = %v, want 0", a)
	}
}

func BenchmarkPaybackPeriod(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		flows := benchFlows(n)
		flows[0] = -100_000 * int64(n-1) // recovered only at the last period: the worst case
		b.Run(itoa(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				sinkFloat, sinkErr = PaybackPeriod(flows)
			}
		})
	}
}
