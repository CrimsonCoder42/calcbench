package text

import (
	"errors"
	"testing"
)

// Semantic Versioning 2.0.0, item 11: the two precedence chains the specification gives as examples.
func TestCompareSemVer_specChains(t *testing.T) {
	chains := [][]string{
		{"1.0.0", "2.0.0", "2.1.0", "2.1.1"},
		{"1.0.0-alpha", "1.0.0-alpha.1", "1.0.0-alpha.beta", "1.0.0-beta", "1.0.0-beta.2", "1.0.0-beta.11", "1.0.0-rc.1", "1.0.0"},
	}
	for _, chain := range chains {
		for i := range chain {
			for j := range chain {
				want := 0
				if i < j {
					want = -1
				} else if i > j {
					want = 1
				}
				if got, err := CompareSemVer(chain[i], chain[j]); got != want || err != nil {
					t.Errorf("CompareSemVer(%q, %q) = %d, %v; want %d, nil", chain[i], chain[j], got, err, want)
				}
			}
		}
	}
}

func TestCompareSemVer_cases(t *testing.T) {
	tests := []struct {
		name string
		a, b string
		want int
	}{
		{"equal", "1.2.3", "1.2.3", 0},
		{"zeros", "0.0.0", "0.0.0", 0},
		{"numeric, not ASCII", "1.10.0", "1.9.0", 1},
		{"patch decides", "1.2.3", "1.2.4", -1},
		{"build metadata ignored", "1.0.0+20130313144700", "1.0.0+exp.sha.5114f85", 0},
		{"build with pre-release ignored", "1.0.0-beta+exp", "1.0.0-beta", 0},
		{"numbers beyond int64", "99999999999999999999.0.0", "9.0.0", 1},
		{"numeric identifier below alphanumeric", "1.0.0-1", "1.0.0-a", -1},
		{"hyphens inside identifiers", "1.0.0-x-y-z.--", "1.0.0-x-y-z.-", 1},
		{"identifier with digits and letters is alphanumeric", "1.0.0-1a", "1.0.0-2", 1},
		{"pre-release 0", "1.0.0-0", "1.0.0-0.0", -1},
		{"build may have leading zeros", "1.0.0+001", "1.0.0", 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := CompareSemVer(tt.a, tt.b); got != tt.want || err != nil {
				t.Errorf("CompareSemVer(%q, %q) = %d, %v; want %d, nil", tt.a, tt.b, got, err, tt.want)
			}
			if got, err := CompareSemVer(tt.b, tt.a); got != -tt.want || err != nil {
				t.Errorf("CompareSemVer(%q, %q) = %d, %v; want %d, nil", tt.b, tt.a, got, err, -tt.want)
			}
		})
	}
}

func TestCompareSemVer_errors(t *testing.T) {
	invalid := []string{
		"", "1", "1.2", "1.2.3.4", "v1.2.3", "01.2.3", "1.02.3", "1.2.03", "1..3", "1.2.", ".1.2",
		"1.2.3-", "1.2.3+", "1.2.3-01", "1.2.3-a..b", "1.2.3-a.", "1.2.3-.a", "1.2.3+.a", "1.2.3+a..b",
		"1.2.3-a_b", "1.2.3+a b", "1.2.3-é", "-1.2.3", "1.2.-3", " 1.2.3", "1.2.3 ", "a.b.c",
	}
	for _, s := range invalid {
		if _, err := CompareSemVer(s, "1.0.0"); !errors.Is(err, ErrSemVer) {
			t.Errorf("CompareSemVer(%q, valid) error = %v, want ErrSemVer", s, err)
		}
		if _, err := CompareSemVer("1.0.0", s); !errors.Is(err, ErrSemVer) {
			t.Errorf("CompareSemVer(valid, %q) error = %v, want ErrSemVer", s, err)
		}
	}
}

func TestCompareSemVer_noAllocation(t *testing.T) {
	if allocs := testing.AllocsPerRun(100, func() { benchSink, _ = CompareSemVer("1.0.0-beta.11", "1.0.0-beta.2") }); allocs != 0 {
		t.Errorf("CompareSemVer allocates %v times, want 0", allocs)
	}
}

// Transitivity and antisymmetry over a mixed set: the comparison must be a total preorder.
func TestCompareSemVer_oracleOrderProperties(t *testing.T) {
	set := []string{
		"0.0.1", "0.1.0", "1.0.0-0", "1.0.0-0.a", "1.0.0-alpha", "1.0.0-alpha.1", "1.0.0-alpha.beta",
		"1.0.0-beta.11", "1.0.0-beta.2", "1.0.0", "1.0.0+b", "1.0.1-rc", "10.0.0", "2.0.0",
	}
	for _, a := range set {
		for _, b := range set {
			ab, _ := CompareSemVer(a, b)
			ba, _ := CompareSemVer(b, a)
			if ab != -ba {
				t.Fatalf("antisymmetry: (%q, %q) = %d but reversed = %d", a, b, ab, ba)
			}
			for _, c := range set {
				bc, _ := CompareSemVer(b, c)
				ac, _ := CompareSemVer(a, c)
				if ab <= 0 && bc <= 0 && ac > 0 {
					t.Fatalf("transitivity: %q <= %q <= %q but %q > %q", a, b, c, a, c)
				}
			}
		}
	}
}

func BenchmarkCompareSemVer(b *testing.B) {
	b.Run("scalar", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			benchSink, _ = CompareSemVer("1.0.0-alpha.beta.11+build.7", "1.0.0-alpha.beta.2")
		}
	})
}
