package prob

import "testing"

// Ross §2.4 / Casella–Berger §1.2: one card, A = heart (13/52), B = king (4/52), A∩B = 1/52
// → 16/52 = 4/13. Hand-checked; tolerance a few ulps for the rounded inputs.
func TestUnionProbability_sourceExample(t *testing.T) {
	if got, err := UnionProbability(13.0/52, 4.0/52, 1.0/52); err != nil || !within(got, 4.0/13, 1e-15) {
		t.Errorf("UnionProbability = (%v, %v), want (%v, nil)", got, err, 4.0/13)
	}
}

func TestUnionProbability_cases(t *testing.T) {
	tests := []struct {
		name           string
		a, b, ab, want float64
	}{
		{"mutually exclusive", 0.25, 0.5, 0, 0.75},
		{"B inside A", 0.5, 0.25, 0.25, 0.5},
		{"all zero", 0, 0, 0, 0},
		{"certain union", 1, 1, 1, 1},
		{"complementary", 0.25, 0.75, 0, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := UnionProbability(tt.a, tt.b, tt.ab); err != nil || got != tt.want {
				t.Errorf("UnionProbability(%v, %v, %v) = (%v, %v), want %v", tt.a, tt.b, tt.ab, got, err, tt.want)
			}
		})
	}
}

func TestUnionProbability_errors(t *testing.T) {
	_, err := UnionProbability(nan, 0.5, 0)
	requireErr(t, "pA NaN", err, ErrNotFinite)
	_, err = UnionProbability(0.5, 0.5, -inf)
	requireErr(t, "pAAndB -Inf", err, ErrNotFinite)
	_, err = UnionProbability(0.5, 1.5, 0)
	requireErr(t, "pB > 1", err, ErrProbability)
	_, err = UnionProbability(0.3, 0.5, 0.4)
	requireErr(t, "pAAndB > pA", err, ErrInconsistent)
	_, err = UnionProbability(0.75, 0.5, 0)
	requireErr(t, "union > 1", err, ErrInconsistent)
}

func TestUnionProbability_noAlloc(t *testing.T) {
	requireNoAlloc(t, "UnionProbability", func() { sinkFloat, _ = UnionProbability(0.25, 0.5, 0.125) })
}

func BenchmarkUnionProbability(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = UnionProbability(0.25, 0.5, 0.125)
	}
}
