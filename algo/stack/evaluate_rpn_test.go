package stack

import (
	"errors"
	"math"
	"math/rand/v2"
	"strconv"
	"strings"
	"testing"
)

// Hamblin (1962) reverse Polish order on our own expression (5 + 1) × (12 − 4) / 3 = 16. Hand-derived.
func TestEvaluateRPN_example(t *testing.T) {
	if got, err := EvaluateRPN(strings.Fields("5 1 + 12 4 - * 3 /")); err != nil || got != 16 {
		t.Errorf("EvaluateRPN = (%d, %v), want (16, nil)", got, err)
	}
}

func TestEvaluateRPN_cases(t *testing.T) {
	tests := []struct {
		name string
		expr string
		want int64
	}{
		{"one literal", "42", 42},
		{"signed literals", "-7 +3 +", -4},
		{"subtraction order", "10 3 -", 7},
		{"division order", "20 4 /", 5},
		{"truncates toward zero", "-7 2 /", -3},
		{"truncates positive", "7 -2 /", -3},
		{"zero product", "0 9223372036854775807 *", 0},
		{"nested", "2 3 4 * +", 14},
		{"max literal", "9223372036854775807", math.MaxInt64},
		{"min literal", "-9223372036854775808", math.MinInt64},
		{"reach max exactly", "9223372036854775806 1 +", math.MaxInt64},
		{"reach min exactly", "-9223372036854775807 1 -", math.MinInt64},
		{"min times one", "-9223372036854775808 1 *", math.MinInt64},
		{"min divided by one", "-9223372036854775808 1 /", math.MinInt64},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := EvaluateRPN(strings.Fields(tt.expr)); err != nil || got != tt.want {
				t.Errorf("EvaluateRPN(%q) = (%d, %v), want (%d, nil)", tt.expr, got, err, tt.want)
			}
		})
	}
}

func TestEvaluateRPN_errors(t *testing.T) {
	tests := []struct {
		name   string
		tokens []string
		want   error
	}{
		{"no tokens", nil, ErrMalformed},
		{"operator alone", []string{"+"}, ErrMalformed},
		{"one operand", []string{"1", "-"}, ErrMalformed},
		{"left over", []string{"1", "2"}, ErrMalformed},
		{"unknown operator", []string{"1", "2", "%"}, ErrBadToken},
		{"empty token", []string{""}, ErrBadToken},
		{"float literal", []string{"1.5"}, ErrBadToken},
		{"spaces in token", []string{" 1"}, ErrBadToken},
		{"literal too large", []string{"9223372036854775808"}, ErrOverflow},
		{"add overflow", strings.Fields("9223372036854775807 1 +"), ErrOverflow},
		{"add underflow", strings.Fields("-9223372036854775808 -1 +"), ErrOverflow},
		{"sub overflow", strings.Fields("9223372036854775807 -1 -"), ErrOverflow},
		{"sub underflow", strings.Fields("-9223372036854775808 1 -"), ErrOverflow},
		{"mul overflow", strings.Fields("4611686018427387904 2 *"), ErrOverflow},
		{"mul min by -1", strings.Fields("-9223372036854775808 -1 *"), ErrOverflow},
		{"mul -1 by min", strings.Fields("-1 -9223372036854775808 *"), ErrOverflow},
		{"div min by -1", strings.Fields("-9223372036854775808 -1 /"), ErrOverflow},
		{"divide by zero", strings.Fields("1 0 /"), ErrDivideByZero},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := EvaluateRPN(tt.tokens); !errors.Is(err, tt.want) || got != 0 {
				t.Errorf("EvaluateRPN(%q) = (%d, %v), want (0, %v)", tt.tokens, got, err, tt.want)
			}
		})
	}
}

// Oracle: random expression trees of depth ≤ 4 over operands in [-6, 6], evaluated directly while built.
// At most 16 operands of magnitude ≤ 6 keep every value below 6¹⁶ ≈ 2.8e12, far from the int64 limits.
func TestEvaluateRPN_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	var build func(depth int) ([]string, int64, bool)
	build = func(depth int) ([]string, int64, bool) {
		if depth == 0 || rng.IntN(3) == 0 {
			v := int64(rng.IntN(13) - 6)
			return []string{strconv.FormatInt(v, 10)}, v, true
		}
		lt, lv, lok := build(depth - 1)
		rt, rv, rok := build(depth - 1)
		op := "+-*/"[rng.IntN(4)]
		var v int64
		ok := lok && rok
		switch op {
		case '+':
			v = lv + rv
		case '-':
			v = lv - rv
		case '*':
			v = lv * rv
		case '/':
			if rv == 0 {
				ok = false
			} else {
				v = lv / rv
			}
		}
		return append(append(lt, rt...), string(op)), v, ok
	}
	for trial := range 400 {
		tokens, want, ok := build(4)
		got, err := EvaluateRPN(tokens)
		if ok && (err != nil || got != want) {
			t.Fatalf("trial %d: EvaluateRPN(%q) = (%d, %v), want %d", trial, tokens, got, err, want)
		}
		if !ok && !errors.Is(err, ErrDivideByZero) {
			t.Fatalf("trial %d: EvaluateRPN(%q) = (%d, %v), want ErrDivideByZero", trial, tokens, got, err)
		}
	}
}

func BenchmarkEvaluateRPN(b *testing.B) {
	for _, size := range benchSizes {
		// 1 2 + 3 + 4 + … : n tokens alternating literal and operator after the first.
		tokens := []string{"1"}
		for i := 1; len(tokens) < size.n; i++ {
			tokens = append(tokens, strconv.Itoa(i%1000), "+")
		}
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _ = EvaluateRPN(tokens)
			}
		})
	}
}
