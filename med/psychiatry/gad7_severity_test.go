package psychiatry

import (
	"errors"
	"testing"
)

// Manual p. 6 cutpoints 5, 10, 15; each edge on both sides.
func TestGAD7SeverityOf_example(t *testing.T) {
	tests := []struct {
		total int
		want  GAD7Severity
	}{
		{0, GAD7BelowMild}, {4, GAD7BelowMild}, {5, GAD7Mild}, {9, GAD7Mild},
		{10, GAD7Moderate}, {14, GAD7Moderate}, {15, GAD7Severe}, {21, GAD7Severe},
	}
	for _, tc := range tests {
		if got, err := GAD7SeverityOf(tc.total); err != nil || got != tc.want {
			t.Errorf("GAD7SeverityOf(%d) = (%q, %v), want %q", tc.total, got, err, tc.want)
		}
	}
	for _, bad := range []int{-1, 22} {
		if _, err := GAD7SeverityOf(bad); !errors.Is(err, ErrOutOfRange) {
			t.Errorf("GAD7SeverityOf(%d) error = %v, want ErrOutOfRange", bad, err)
		}
	}
}

// parity: ledger gad7, source cases 1–7. The band edges agree with the source; the labels are the manual's
// words, not the source's ("Minimal/Mild/Moderate/Severe Anxiety"); for 0–4 the manual names no band, where
// the source says "Minimal Anxiety".
func TestGAD7SeverityOf_parity(t *testing.T) {
	tests := []struct {
		total int
		want  GAD7Severity
	}{
		{0, GAD7BelowMild}, // source case 1: source "Minimal Anxiety"
		{4, GAD7BelowMild}, // source case 2: source "Minimal Anxiety"
		{5, GAD7Mild},      // source case 3: source "Mild Anxiety"
		{9, GAD7Mild},      // source case 4: source "Mild Anxiety"
		{12, GAD7Moderate}, // source case 5: source "Moderate Anxiety"
		{15, GAD7Severe},   // source case 6: source "Severe Anxiety"
		{21, GAD7Severe},   // source case 7: source "Severe Anxiety"
	}
	for i, tc := range tests {
		if got, err := GAD7SeverityOf(tc.total); err != nil || got != tc.want {
			t.Errorf("source case %d: GAD7SeverityOf(%d) = (%q, %v), want %q", i+1, tc.total, got, err, tc.want)
		}
	}
}

func TestGAD7SeverityOf_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { c, _ := GAD7SeverityOf(12); sinkString = string(c) }); n != 0 {
		t.Errorf("allocates %v times, want 0", n)
	}
}

func BenchmarkGAD7SeverityOf(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		c, _ := GAD7SeverityOf(12)
		sinkString = string(c)
	}
}
