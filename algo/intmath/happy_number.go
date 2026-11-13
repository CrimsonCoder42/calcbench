package intmath

// digitSquareSum returns the sum of the squares of the decimal digits of n (≤ 20·81 = 1620).
func digitSquareSum(n uint64) uint64 {
	var s uint64
	for ; n > 0; n /= 10 {
		d := n % 10
		s += d * d
	}
	return s
}

// IsHappyNumber reports whether repeatedly replacing n by the sum of the squares of its decimal digits
// reaches 1 (7 → 49 → 97 → 130 → 10 → 1 is happy; 4 → 16 → 37 → … → 4 is not). The sequence is
// eventually periodic, and Floyd's cycle-finding method decides which cycle it enters with two
// iterates (one moving twice as fast) and no table. 0 is not happy (it maps to itself).
//
// ID:      ALG-INTMTH-006
// Time:    O(log n)   the first step reads the ≤ 20 digits of n; every later value is ≤ 1620 and the cycle is found in O(1) further steps
// Space:   O(1), no allocation
// Source:  Knuth, The Art of Computer Programming, Vol. 2, 3rd ed., Addison-Wesley, 1997, §3.1, exercise 6 (Floyd's cycle-finding method: find n > 0 with X_n = X_2n using O(1) memory)
// Source:  Guy, Unsolved Problems in Number Theory, 3rd ed., Springer, 2004, problem E34 "Happy numbers" (definition; every unhappy number reaches the cycle 4, 16, 37, 58, 89, 145, 42, 20)
// Tags:    integer, digits, happy-number, cycle-detection, floyd, tortoise-hare
func IsHappyNumber(n uint64) bool {
	slow, fast := n, digitSquareSum(n)
	for fast != 1 && slow != fast {
		slow = digitSquareSum(slow)
		fast = digitSquareSum(digitSquareSum(fast))
	}
	return fast == 1
}
