package stack

import (
	"errors"
	"fmt"
	"math"
	"strconv"
)

// Errors of EvaluateRPN (ErrOverflow is shared with the rectangle functions). Wrapped with the token's
// position; test with errors.Is.
var (
	// ErrMalformed: the expression is empty, an operator lacks two operands, or values are left over.
	ErrMalformed = errors.New("stack: malformed reverse Polish expression")
	// ErrBadToken: a token is neither an integer literal nor one of + - * /.
	ErrBadToken = errors.New("stack: token is not an integer or an operator")
	// ErrDivideByZero: a division has a zero divisor.
	ErrDivideByZero = errors.New("stack: division by zero")
)

// EvaluateRPN returns the value of the reverse Polish (postfix) expression tokens. Each token is a
// base-10 int64 literal with an optional sign ("-7", "+3") or one of the operators + - * /, which pop
// the right operand, then the left, and push the result. Division truncates toward zero, as Go's / does.
// Every result is exact or an error. tokens is read, never modified.
// Errors: ErrMalformed when tokens is empty, an operator has fewer than two operands, or more than one
// value is left; ErrBadToken for any other token; ErrOverflow when a literal or a result is outside
// int64 (including math.MinInt64 / -1); ErrDivideByZero when a divisor is 0.
//
// ID:      ALG-STACK-009
// Time:    O(n + Σ Lᵢ)   n = len(tokens), Lᵢ = bytes of token i
// Space:   O(n) extra for the operand stack (one allocation)
// Source:  Hamblin, "Translation to and from Polish notation", The Computer Journal, 1962, 5(3):210–213, doi:10.1093/comjnl/5.3.210
// Source:  Sedgewick, Wayne, Algorithms, 4th ed., Addison-Wesley, 2011, §1.3 (Dijkstra's two-stack evaluator; Exercise 1.3.11 EvaluatePostfix)
// Source:  Knuth, The Art of Computer Programming, Vol. 1, 3rd ed., Addison-Wesley, 1997, §2.2.1 (stacks)
// Tags:    stack, postfix, reverse-polish, rpn, expression, evaluate, calculator, overflow-checked
func EvaluateRPN(tokens []string) (int64, error) {
	if len(tokens) == 0 {
		return 0, fmt.Errorf("%w: no tokens", ErrMalformed)
	}
	operands := make([]int64, 0, len(tokens)/2+1)
	for i, tok := range tokens {
		if len(tok) == 1 && (tok[0] == '+' || tok[0] == '-' || tok[0] == '*' || tok[0] == '/') {
			n := len(operands)
			if n < 2 {
				return 0, fmt.Errorf("%w: operator %q at token %d has %d operand(s)", ErrMalformed, tok, i, n)
			}
			left, right := operands[n-2], operands[n-1]
			r, err := applyOperator(tok[0], left, right)
			if err != nil {
				return 0, fmt.Errorf("%w: %d %s %d at token %d", err, left, tok, right, i)
			}
			operands = append(operands[:n-2], r)
			continue
		}
		v, err := strconv.ParseInt(tok, 10, 64)
		if err != nil {
			if errors.Is(err, strconv.ErrRange) {
				return 0, fmt.Errorf("%w: literal %q at token %d", ErrOverflow, tok, i)
			}
			return 0, fmt.Errorf("%w: %q at token %d", ErrBadToken, tok, i)
		}
		operands = append(operands, v)
	}
	if len(operands) != 1 {
		return 0, fmt.Errorf("%w: %d values left", ErrMalformed, len(operands))
	}
	return operands[0], nil
}

// applyOperator returns left op right, exactly, or ErrOverflow / ErrDivideByZero.
func applyOperator(op byte, left, right int64) (int64, error) {
	switch op {
	case '+':
		if (right > 0 && left > math.MaxInt64-right) || (right < 0 && left < math.MinInt64-right) {
			return 0, ErrOverflow
		}
		return left + right, nil
	case '-':
		if (right < 0 && left > math.MaxInt64+right) || (right > 0 && left < math.MinInt64+right) {
			return 0, ErrOverflow
		}
		return left - right, nil
	case '*':
		if left == 0 || right == 0 {
			return 0, nil
		}
		if (left == -1 && right == math.MinInt64) || (right == -1 && left == math.MinInt64) {
			return 0, ErrOverflow
		}
		p := left * right
		if p/right != left {
			return 0, ErrOverflow
		}
		return p, nil
	default: // '/'
		if right == 0 {
			return 0, ErrDivideByZero
		}
		if left == math.MinInt64 && right == -1 {
			return 0, ErrOverflow
		}
		return left / right, nil
	}
}
