package neuro

import (
	"errors"
	"testing"
)

// Brennan 2018 abstract: GCS (3–15) minus non-reacting pupils (0–2), range 1–15.
func TestGCSPupilsScore_example(t *testing.T) {
	tests := []struct{ gcs, pupils, want int }{
		{3, 2, 1}, {15, 0, 15}, {8, 1, 7}, {3, 0, 3}, {15, 2, 13},
	}
	for _, tc := range tests {
		if got, err := GCSPupilsScore(tc.gcs, tc.pupils); err != nil || got != tc.want {
			t.Errorf("GCSPupilsScore(%d, %d) = (%d, %v), want %d", tc.gcs, tc.pupils, got, err, tc.want)
		}
	}
	for _, in := range [][2]int{{2, 0}, {16, 0}, {10, -1}, {10, 3}} {
		if _, err := GCSPupilsScore(in[0], in[1]); !errors.Is(err, ErrOutOfRange) {
			t.Errorf("GCSPupilsScore(%v) error = %v, want ErrOutOfRange", in, err)
		}
	}
}

func TestGCSPupilsScore_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkInt, _ = GCSPupilsScore(8, 1) }); n != 0 {
		t.Errorf("allocates %v times, want 0", n)
	}
}

func BenchmarkGCSPupilsScore(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkInt, _ = GCSPupilsScore(8, 1)
	}
}
