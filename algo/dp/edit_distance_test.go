package dp

import (
	"math"
	"math/rand/v2"
	"testing"
)

// editDistanceSlow fills Wagner and Fischer's full (m+1)×(n+1) cost table with unit costs.
func editDistanceSlow(a, b []byte) int {
	d := make([][]int, len(a)+1)
	for i := range d {
		d[i] = make([]int, len(b)+1)
		d[i][0] = i
	}
	for j := range d[0] {
		d[0][j] = j
	}
	for i := 1; i <= len(a); i++ {
		for j := 1; j <= len(b); j++ {
			sub := d[i-1][j-1]
			if a[i-1] != b[j-1] {
				sub++
			}
			d[i][j] = min(sub, d[i-1][j]+1, d[i][j-1]+1)
		}
	}
	return d[len(a)][len(b)]
}

// Jurafsky, Martin, Speech and Language Processing, 3rd ed. draft (Dec 2021), §2.5 "Minimum Edit
// Distance", p. 22 and Fig. 2.14: with Levenshtein unit costs (insert, delete, substitute each 1, the
// cost model EditDistance uses) intention → execution is 5. (J&M's eq. (2.8) variant, substitution
// cost 2, gives 8 in Fig. 2.18; that variant is not this function.) Our own pair, by hand: harbor →
// arbour deletes h and inserts u (2); no single edit suffices, as every position differs.
func TestEditDistance_workedExamples(t *testing.T) {
	tests := []struct {
		from, to string
		want     int
	}{
		{"harbor", "arbour", 2},
		{"intention", "execution", 5},
	}
	for _, tt := range tests {
		if got := EditDistance([]byte(tt.from), []byte(tt.to)); got != tt.want {
			t.Errorf("EditDistance(%q, %q) = %d, want %d", tt.from, tt.to, got, tt.want)
		}
	}
}

func TestEditDistance_cases(t *testing.T) {
	tests := []struct {
		name     string
		from, to string
		want     int
	}{
		{"both empty", "", "", 0},
		{"from empty", "", "abc", 3},
		{"to empty", "abcd", "", 4},
		{"identical", "same", "same", 0},
		{"one substitution", "cat", "cut", 1},
		{"one insertion", "cat", "cart", 1},
		{"one deletion", "cart", "cat", 1},
		{"reversed", "abc", "cba", 2},
		{"duplicates", "aaaa", "aa", 2},
		{"disjoint", "abc", "xyz", 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EditDistance([]byte(tt.from), []byte(tt.to)); got != tt.want {
				t.Errorf("EditDistance(%q, %q) = %d, want %d", tt.from, tt.to, got, tt.want)
			}
			if got := EditDistance([]byte(tt.to), []byte(tt.from)); got != tt.want {
				t.Errorf("EditDistance(%q, %q) = %d, want %d (symmetric)", tt.to, tt.from, got, tt.want)
			}
		})
	}
}

func TestEditDistance_runesVersusBytes(t *testing.T) {
	// "café" → "cafe": one character substituted, but é is two bytes in UTF-8.
	if got := EditDistance([]rune("café"), []rune("cafe")); got != 1 {
		t.Errorf("runes: got %d, want 1", got)
	}
	if got := EditDistance([]byte("café"), []byte("cafe")); got != 2 {
		t.Errorf("bytes: got %d, want 2", got)
	}
}

func TestEditDistance_nan(t *testing.T) {
	a := []float64{1, math.NaN()}
	if got := EditDistance(a, a); got != 1 {
		t.Errorf("EditDistance with NaN = %d, want 1 (NaN != NaN)", got)
	}
}

func TestEditDistance_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 400 {
		a := randomBytes(rng, rng.IntN(30), "abc")
		b := randomBytes(rng, rng.IntN(30), "abc")
		if got, want := EditDistance(a, b), editDistanceSlow(a, b); got != want {
			t.Fatalf("trial %d: EditDistance(%q, %q) = %d, want %d", trial, a, b, got, want)
		}
	}
}

func TestOneEditDistance_cases(t *testing.T) {
	tests := []struct {
		name     string
		from, to string
		want     bool
	}{
		{"both empty", "", "", false},
		{"insert into empty", "", "a", true},
		{"two inserts", "", "ab", false},
		{"identical", "abc", "abc", false},
		{"substitution first", "abc", "xbc", true},
		{"substitution last", "abc", "abx", true},
		{"two substitutions", "abc", "xbx", false},
		{"deletion middle", "abc", "ac", true},
		{"deletion last", "abc", "ab", true},
		{"insertion first", "bc", "abc", true},
		{"length differs by 2", "a", "abc", false},
		{"same length, swap", "ab", "ba", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := OneEditDistance([]byte(tt.from), []byte(tt.to)); got != tt.want {
				t.Errorf("OneEditDistance(%q, %q) = %v, want %v", tt.from, tt.to, got, tt.want)
			}
			if got := OneEditDistance([]byte(tt.to), []byte(tt.from)); got != tt.want {
				t.Errorf("OneEditDistance(%q, %q) = %v, want %v (symmetric)", tt.to, tt.from, got, tt.want)
			}
		})
	}
}

func TestOneEditDistance_allocations(t *testing.T) {
	a, b := []byte("abcdef"), []byte("abcxef")
	if n := testing.AllocsPerRun(100, func() { OneEditDistance(a, b) }); n != 0 {
		t.Errorf("OneEditDistance allocates %v times, want 0", n)
	}
}

// Oracle: exactly one edit apart iff the full distance is 1.
func TestOneEditDistance_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 600 {
		a := randomBytes(rng, rng.IntN(6), "ab")
		b := randomBytes(rng, rng.IntN(6), "ab")
		if got, want := OneEditDistance(a, b), editDistanceSlow(a, b) == 1; got != want {
			t.Fatalf("trial %d: OneEditDistance(%q, %q) = %v, want %v", trial, a, b, got, want)
		}
	}
}

func BenchmarkEditDistance(b *testing.B) {
	for _, size := range quadraticSizes {
		rng := rand.New(rand.NewPCG(1, 2))
		x, y := randomBytes(rng, size.n, "abcd"), randomBytes(rng, size.n, "abcd")
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_ = EditDistance(x, y)
			}
		})
	}
}

func BenchmarkOneEditDistance(b *testing.B) {
	for _, size := range linearSizes {
		rng := rand.New(rand.NewPCG(1, 2))
		x := randomBytes(rng, size.n, "abcd")
		y := append([]byte(nil), x...)
		y[size.n-1] ^= 1 // the one difference is at the end: the whole prefix is compared
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_ = OneEditDistance(x, y)
			}
		})
	}
}
