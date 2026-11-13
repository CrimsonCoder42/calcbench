package text

import (
	"encoding/json"
	"math/rand/v2"
	"regexp"
	"testing"
)

// jsonNumberRE transcribes the RFC 8259 §6 ABNF as a regular expression; the oracle.
var jsonNumberRE = regexp.MustCompile(`^-?(0|[1-9][0-9]*)(\.[0-9]+)?([eE][-+]?[0-9]+)?$`)

// RFC 8259 §6: each production of the grammar, accepted, and the forms it excludes, rejected.
func TestIsJSONNumber_rfc8259Grammar(t *testing.T) {
	valid := []string{
		"0", "-0", "7", "-7", "10", "1234567890", "0.5", "-0.25", "3.14159", "1e5", "1E5", "1e+5", "1e-5",
		"-1.5e-10", "0e0", "0.0E00", "1e999", "123456789012345678901234567890",
	}
	invalid := []string{
		"", "-", "+1", "01", "-01", "00", ".5", "5.", "1.e5", "1e", "1e+", "e5", "1.5.2", "1e5.0", "--1",
		" 1", "1 ", "0x1F", "Infinity", "NaN", "1_000", "1,5", "١", "1e5e5", "-.", "0.",
	}
	for _, s := range valid {
		if !IsJSONNumber(s) {
			t.Errorf("IsJSONNumber(%q) = false, want true", s)
		}
	}
	for _, s := range invalid {
		if IsJSONNumber(s) {
			t.Errorf("IsJSONNumber(%q) = true, want false", s)
		}
	}
}

func TestIsJSONNumber_noAllocation(t *testing.T) {
	allocs := testing.AllocsPerRun(100, func() {
		if IsJSONNumber("-1.5e-10") {
			benchSink++
		}
	})
	if allocs != 0 {
		t.Errorf("IsJSONNumber allocates %v times, want 0", allocs)
	}
}

// Two oracles on strings over the grammar's own alphabet: the regular expression above, and
// encoding/json (a JSON text built only from these bytes, with no white space, can only be a number).
func TestIsJSONNumber_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 2000 {
		s := randomText(rng, rng.IntN(8), "-+.eE0019")
		got := IsJSONNumber(s)
		if want := jsonNumberRE.MatchString(s); got != want {
			t.Fatalf("trial %d: IsJSONNumber(%q) = %v, regexp says %v", trial, s, got, want)
		}
		if want := json.Valid([]byte(s)); got != want {
			t.Fatalf("trial %d: IsJSONNumber(%q) = %v, json.Valid says %v", trial, s, got, want)
		}
	}
}

func BenchmarkIsJSONNumber(b *testing.B) {
	b.Run("scalar", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			if IsJSONNumber("-12345.678901e-23") {
				benchSink++
			}
		}
	})
}
