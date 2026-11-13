package stack

import (
	"errors"
	"math/rand/v2"
	"strings"
	"testing"
)

const standardPairs = "()[]{}"

// Sedgewick & Wayne, Algorithms, 4th ed., §1.3, Exercise 1.3.4: the book's own two strings.
func TestBalancedBrackets_sedgewickWayne(t *testing.T) {
	if got, err := BalancedBrackets("[()]{}{[()()]()}", standardPairs); err != nil || !got {
		t.Errorf("BalancedBrackets(balanced) = (%v, %v), want (true, nil)", got, err)
	}
	if got, err := BalancedBrackets("[(])", standardPairs); err != nil || got {
		t.Errorf("BalancedBrackets([(])) = (%v, %v), want (false, nil)", got, err)
	}
}

func TestBalancedBrackets_cases(t *testing.T) {
	tests := []struct {
		name  string
		text  string
		pairs string
		want  bool
	}{
		{"empty", "", standardPairs, true},
		{"no pairs", "(((", "", true},
		{"one opener", "(", standardPairs, false},
		{"one closer", ")", standardPairs, false},
		{"pair", "()", standardPairs, true},
		{"reversed pair", ")(", standardPairs, false},
		{"other runes ignored", "f(x) = [a, {b}]", standardPairs, true},
		{"crossed", "([)]", standardPairs, false},
		{"only parentheses checked", "[(])", "()", true},
		{"unicode pairs", "«‹›»", "«»‹›", true},
		{"unicode crossed", "«‹»›", "«»‹›", false},
		{"invalid bytes ignored", "(\xff)", standardPairs, true},
		{"deep", strings.Repeat("(", 1000) + strings.Repeat(")", 1000), standardPairs, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := BalancedBrackets(tt.text, tt.pairs); err != nil || got != tt.want {
				t.Errorf("BalancedBrackets(%q, %q) = (%v, %v), want (%v, nil)", tt.text, tt.pairs, got, err, tt.want)
			}
		})
	}
}

func TestBalancedBrackets_errors(t *testing.T) {
	for _, pairs := range []string{"(", "()[", "||", "()(]", "(())", "«»«]", "a\xff", "(�"} {
		if got, err := BalancedBrackets("()", pairs); !errors.Is(err, ErrPairs) || got {
			t.Errorf("BalancedBrackets(pairs %q) = (%v, %v), want (false, ErrPairs)", pairs, got, err)
		}
	}
}

// Oracle: repeatedly delete adjacent matched pairs; balanced exactly when nothing is left.
func TestBalancedBrackets_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 500 {
		var sb strings.Builder
		for range rng.IntN(16) {
			sb.WriteByte("()[]"[rng.IntN(4)])
		}
		text := sb.String()
		reduced := text
		for {
			next := strings.ReplaceAll(strings.ReplaceAll(reduced, "()", ""), "[]", "")
			if next == reduced {
				break
			}
			reduced = next
		}
		want := reduced == ""
		if got, err := BalancedBrackets(text, "()[]"); err != nil || got != want {
			t.Fatalf("trial %d: BalancedBrackets(%q) = (%v, %v), want %v", trial, text, got, err, want)
		}
	}
}

func BenchmarkBalancedBrackets(b *testing.B) {
	for _, size := range benchSizes {
		text := strings.Repeat("({[", size.n/6) + strings.Repeat("]})", size.n/6)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _ = BalancedBrackets(text, standardPairs)
			}
		})
	}
}
