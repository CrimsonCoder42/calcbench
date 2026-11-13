package prob

import "testing"

// simpson integrates f over [a, b] by the composite Simpson rule with n (even) panels (Burden & Faires 9th
// ed. §4.4); the identity tests use it to check that a density integrates to its CDF difference.
func simpson(f func(float64) float64, a, b float64, n int) float64 {
	h := (b - a) / float64(n)
	s := f(a) + f(b)
	for i := 1; i < n; i++ {
		w := 2.0
		if i%2 == 1 {
			w = 4
		}
		s += w * f(a+float64(i)*h)
	}
	return s * h / 3
}

// The bisection solves an increasing h exactly to 4 ulps: the uniform CDF on [0, 1] inverts to u itself.
func TestQuantileByBisection_uniform(t *testing.T) {
	for _, u := range []float64{1e-300, 1e-10, 0.25, 0.5, 0.75, 1 - 1e-10} {
		got, err := quantileByBisection(u, 0, 1, func(x float64) (float64, float64, error) {
			return min(x, 1), max(1-x, 0), nil
		})
		if err != nil || !within(got, u, 1e-15) {
			t.Errorf("uniform u = %v: (%v, %v)", u, got, err)
		}
	}
}

func TestQuantileByBisection_overflow(t *testing.T) {
	_, err := quantileByBisection(0.5, 0, 1, func(float64) (float64, float64, error) { return 0, 1, nil })
	requireErr(t, "never reaches u", err, ErrOverflow)
}
