package dp

import "fmt"

// CanJumpToEnd reports whether the last index of maxJump can be reached from index 0 when, from index
// i, one move goes forward to any index in i+1 … i+maxJump[i]. A single element is already at the end
// (true); empty input has no end (false). maxJump is read, never modified.
// Errors: ErrNegative when an entry is negative.
//
// ID:      ALG-DP-020
// Time:    O(n)   n = len(maxJump); one pass keeping the farthest index reached
// Space:   O(1), no allocation
// Source:  Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §22.2 (breadth-first search; on the implicit graph i → i+1 … i+maxJump[i] the reachable set from 0 is a prefix [0, f], so reachability reduces to the farthest index f; proof in catalog/dp.md)
// Tags:    dp, greedy, bfs, reachability, jump-game, implicit-graph
func CanJumpToEnd(maxJump []int) (bool, error) {
	n := len(maxJump)
	if err := checkJumps(maxJump); err != nil || n == 0 {
		return false, err
	}
	farthest := 0
	for i := 0; i <= farthest && farthest < n-1; i++ {
		farthest = max(farthest, reach(i, maxJump[i], n))
	}
	return farthest >= n-1, nil
}

// MinJumpsToEnd returns the fewest moves from index 0 to the last index of maxJump, where one move
// from index i goes forward to any index in i+1 … i+maxJump[i]; ok is false when the end cannot be
// reached. A single element needs 0 moves; empty input has no end (ok = false). maxJump is read, never
// modified.
// Errors: ErrNegative when an entry is negative.
//
// ID:      ALG-DP-021
// Time:    O(n)   n = len(maxJump)
// Space:   O(1), no allocation
// Source:  Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §22.2 (breadth-first search gives fewest-edge paths; on this implicit graph each BFS layer is a contiguous run of indices, so a layer is two numbers and the search is one pass; proof in catalog/dp.md)
// Tags:    dp, greedy, bfs, shortest-path, jump-game, implicit-graph
func MinJumpsToEnd(maxJump []int) (jumps int, ok bool, err error) {
	n := len(maxJump)
	if err := checkJumps(maxJump); err != nil || n == 0 {
		return 0, false, err
	}
	// Layer d is the run (layerEnd of d−1, layerEnd]; next is the farthest index layer d reaches.
	layerEnd, next := 0, 0
	for i := 0; layerEnd < n-1; i++ {
		if i > layerEnd {
			if next <= layerEnd {
				return 0, false, nil // the next layer is empty
			}
			jumps++
			layerEnd = next
			if layerEnd >= n-1 {
				break
			}
		}
		next = max(next, reach(i, maxJump[i], n))
	}
	if layerEnd < n-1 {
		return 0, false, nil
	}
	return jumps, true, nil
}

// reach returns min(i+jump, n−1) without overflowing int.
func reach(i, jump, n int) int {
	if jump >= n-1-i {
		return n - 1
	}
	return i + jump
}

// checkJumps reports a negative entry.
func checkJumps(maxJump []int) error {
	for i, j := range maxJump {
		if j < 0 {
			return fmt.Errorf("%w: maxJump[%d] = %d", ErrNegative, i, j)
		}
	}
	return nil
}
