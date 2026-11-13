package scales

import (
	"errors"
	"testing"
)

// HAND-DERIVED from the GCS Assessment Aid (Teasdale 2015) ratings: eye 1–4, verbal 1–5, motor 1–6, summed.
// E3 (To sound) + V4 (Confused) + M5 (Localising) = 12; E4 V5 M6 = 15 (the maximum); E1 V1 M1 = 3 (the minimum).
func TestGlasgowComaScore_example(t *testing.T) {
	tests := []struct{ eye, verbal, motor, want int }{
		{3, 4, 5, 12},
		{4, 5, 6, 15},
		{1, 1, 1, 3},
		{GCSEyeToPressure, GCSVerbalWords, GCSMotorNormalFlexion, 2 + 3 + 4},
	}
	for _, tc := range tests {
		if got, err := GlasgowComaScore(tc.eye, tc.verbal, tc.motor); err != nil || got != tc.want {
			t.Errorf("GlasgowComaScore(%d, %d, %d) = (%v, %v), want %d", tc.eye, tc.verbal, tc.motor, got, err, tc.want)
		}
	}
}

// The exported rating constants carry the Assessment Aid's numbers.
func TestGlasgowComaScore_constants(t *testing.T) {
	got := []int{GCSNotTestable,
		GCSEyeSpontaneous, GCSEyeToSound, GCSEyeToPressure, GCSEyeNone,
		GCSVerbalOrientated, GCSVerbalConfused, GCSVerbalWords, GCSVerbalSounds, GCSVerbalNone,
		GCSMotorObeysCommands, GCSMotorLocalising, GCSMotorNormalFlexion, GCSMotorAbnormalFlexion, GCSMotorExtension, GCSMotorNone}
	want := []int{0, 4, 3, 2, 1, 5, 4, 3, 2, 1, 6, 5, 4, 3, 2, 1}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("constant %d = %d, want %d", i, got[i], want[i])
		}
	}
}

func TestGlasgowComaScore_errors(t *testing.T) {
	call := func(e, v, m int) func() error {
		return func() error { _, err := GlasgowComaScore(e, v, m); return err }
	}
	runErrCases(t, []errCase{
		{"eye NT", call(GCSNotTestable, 4, 5), ErrNotTestable},
		{"verbal NT", call(3, GCSNotTestable, 5), ErrNotTestable},
		{"motor NT", call(3, 4, GCSNotTestable), ErrNotTestable},
		{"all NT", call(0, 0, 0), ErrNotTestable},
		{"eye -1", call(-1, 4, 5), ErrOutOfRange},
		{"eye 5", call(5, 4, 5), ErrOutOfRange},
		{"verbal -1", call(3, -1, 5), ErrOutOfRange},
		{"verbal 6", call(3, 6, 5), ErrOutOfRange},
		{"motor -1", call(3, 4, -1), ErrOutOfRange},
		{"motor 7", call(3, 4, 7), ErrOutOfRange},
		// Out of range wins over NT, whichever component is which.
		{"eye NT, motor 7", call(0, 4, 7), ErrOutOfRange},
		{"eye 9, verbal NT", call(9, 0, 5), ErrOutOfRange},
	})
	if _, err := GlasgowComaScore(0, 4, 7); errors.Is(err, ErrNotTestable) {
		t.Errorf("eye NT, motor 7: err = %v, want only ErrOutOfRange", err)
	}
}

func TestGlasgowComaScore_noAllocation(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkInt, _ = GlasgowComaScore(3, 4, 5) }); n != 0 {
		t.Errorf("GlasgowComaScore allocates %v times, want 0", n)
	}
}

func BenchmarkGlasgowComaScore(b *testing.B) {
	b.ReportAllocs()
	e, v, m := 3, 4, 5
	for b.Loop() {
		sinkInt, _ = GlasgowComaScore(e, v, m)
	}
}
