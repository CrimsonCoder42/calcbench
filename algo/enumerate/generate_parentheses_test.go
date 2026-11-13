package enumerate

import (
	"errors"
	"fmt"
	"math"
	"testing"
)

func parens(t *testing.T, n int) []string {
	t.Helper()
	seq, err := GenerateParentheses(n)
	if err != nil {
		t.Fatalf("GenerateParentheses(%d): %v", n, err)
	}
	var out []string
	for b := range seq {
		out = append(out, string(b))
	}
	return out
}

// Knuth, TAOCP Vol. 4A, §7.2.1.6, Algorithm P for n = 3 (hand-run of steps P1–P5): the five strings
// from ()()() to ((())).
func TestGenerateParentheses_knuthExample(t *testing.T) {
	got := parens(t, 3)
	want := []string{"()()()", "()(())", "(())()", "(()())", "((()))"}
	if fmt.Sprint(got) != fmt.Sprint(want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestGenerateParentheses_cases(t *testing.T) {
	for n, want := range map[int][]string{0: {""}, 1: {"()"}, 2: {"()()", "(())"}} {
		if got := parens(t, n); fmt.Sprint(got) != fmt.Sprint(want) {
			t.Errorf("n=%d: got %q, want %q", n, got, want)
		}
	}
}

func TestGenerateParentheses_errors(t *testing.T) {
	if _, err := GenerateParentheses(-1); !errors.Is(err, ErrNegative) {
		t.Errorf("-1: err = %v, want ErrNegative", err)
	}
	for _, n := range []int{math.MaxInt, maxAllocBytes / 2} {
		if _, err := GenerateParentheses(n); !errors.Is(err, ErrTooLarge) {
			t.Errorf("%d: err = %v, want ErrTooLarge (no panic)", n, err)
		}
	}
	if _, err := GenerateParentheses((maxAllocBytes - 1) / 2); err != nil {
		t.Errorf("largest allowed: err = %v (the iterator is lazy, nothing is allocated)", err)
	}
}

func TestGenerateParentheses_earlyBreak(t *testing.T) {
	seq, _ := GenerateParentheses(40)
	count := 0
	for range seq {
		count++
		if count == 7 {
			break
		}
	}
	if count != 7 {
		t.Errorf("count = %d, want 7", count)
	}
}

// balanced reports whether s is a Dyck word.
func balanced(s string) bool {
	depth := 0
	for _, c := range s {
		if c == '(' {
			depth++
		} else {
			depth--
		}
		if depth < 0 {
			return false
		}
	}
	return depth == 0
}

// TestGenerateParentheses_oracle: for n ≤ 10, every string is balanced, distinct, strictly increasing
// with ')' before '(' and the count is the Catalan number (computed by C(n+1) = C(n)·2(2n+1)/(n+2)).
func TestGenerateParentheses_oracle(t *testing.T) {
	catalan := 1
	for n := 0; n <= 10; n++ {
		got := parens(t, n)
		if len(got) != catalan {
			t.Fatalf("n=%d: %d strings, want Catalan %d", n, len(got), catalan)
		}
		key := func(s string) string { // map ')' < '(' onto byte order
			b := []byte(s)
			for i := range b {
				if b[i] == ')' {
					b[i] = '0'
				} else {
					b[i] = '1'
				}
			}
			return string(b)
		}
		for i, s := range got {
			if len(s) != 2*n || !balanced(s) {
				t.Fatalf("n=%d: %q is not a balanced string of length %d", n, s, 2*n)
			}
			if i > 0 && key(got[i-1]) >= key(s) {
				t.Fatalf("n=%d: %q does not follow %q", n, s, got[i-1])
			}
		}
		catalan = catalan * 2 * (2*n + 1) / (n + 2)
	}
}

func BenchmarkGenerateParentheses(b *testing.B) {
	for _, n := range []int{8, 10, 12} {
		b.Run(fmt.Sprint(n), func(b *testing.B) {
			b.ReportAllocs()
			sink := 0
			for b.Loop() {
				seq, _ := GenerateParentheses(n)
				for s := range seq {
					sink += int(s[0])
				}
			}
			benchSink = sink
		})
	}
}
