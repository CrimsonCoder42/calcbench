# Knuth, The Art of Computer Programming, Vol. 1

**Citation**: Knuth, *The Art of Computer Programming*, Vol. 1: *Fundamental Algorithms*, 3rd ed., Addison-Wesley, 1997, §1.2.8 (Fibonacci numbers; φ = (1+√5)/2).

Taken: the golden ratio φ = (1+√5)/2. Used by NUM-LOGS-001 GoldenRatio. Checked: golden_ratio_test.go: `TestGoldenRatio_identity` (φ² = φ + 1 in big.Float), `TestGoldenRatio_closedForm`, `TestGoldenRatio_matchesMathPhi` (equals math.Phi). Second source: OEIS A001622 (see oeis.md). Status: section locator from the lead's build, 2026-11-12.


## Used by S10 (algorithms core, 2026-11-12)

IDs: ALG-DP-010, ALG-INTMTH-001, ALG-INTMTH-003, ALG-INTMTH-004, ALG-INTMTH-005, ALG-INTMTH-006, ALG-INTMTH-013, ALG-INTMTH-014, ALG-INTMTH-015, ALG-INTMTH-016, ALG-LIST-001, ALG-LIST-002, ALG-LIST-003, ALG-LIST-004, ALG-LIST-005, ALG-LIST-006, ALG-LIST-007, ALG-LIST-008, ALG-LIST-009, ALG-LIST-010, ALG-LIST-012, ALG-LIST-013, ALG-LIST-014, ALG-LIST-017, ALG-MATRIX-001, ALG-QUEUE-002, ALG-SAMPLE-001, ALG-SAMPLE-002, ALG-SAMPLE-007, ALG-SEQ-008, ALG-STACK-001, ALG-STACK-002, ALG-STACK-003, ALG-STACK-004, ALG-STACK-008, ALG-STACK-009, ALG-STACK-010, ALG-STREAM-001, ALG-STREAM-002, ALG-TREE-006, ALG-TREE-007, ALG-TREE-008, ALG-TREE-014, ALG-TREE-015, ALG-TREE-017.

Citations as written in the doc comments:
- Knuth, The Art of Computer Programming, Vol. 1, 3rd ed., Addison-Wesley, 1997, §2.2.1 (stacks)
- Knuth, The Art of Computer Programming, Vol. 1, 3rd ed., Addison-Wesley, 1997, §2.2.2 (a circular array indexed modulo M)
- Knuth, The Art of Computer Programming, Vol. 1, 3rd ed., Addison-Wesley, 1997, §2.2.2 (a queue in a circular array), via ALG-QUEUE-002
- Knuth, The Art of Computer Programming, Vol. 1, 3rd ed., Addison-Wesley, 1997, §2.2.2 (a queue in a circular array, indices taken modulo M)
- Knuth, The Art of Computer Programming, Vol. 1, 3rd ed., Addison-Wesley, 1997, §2.2.2 (a queue in a circular array, pointers advanced modulo M)
- Knuth, The Art of Computer Programming, Vol. 1, 3rd ed., Addison-Wesley, 1997, §2.2.3 (a stack as a linked list)
- Knuth, The Art of Computer Programming, Vol. 1, 3rd ed., Addison-Wesley, 1997, §2.2.3 (deletion from a linked list by relinking the predecessor)
- Knuth, The Art of Computer Programming, Vol. 1, 3rd ed., Addison-Wesley, 1997, §2.2.3 (linked allocation)
- Knuth, The Art of Computer Programming, Vol. 1, 3rd ed., Addison-Wesley, 1997, §2.2.3 (linked allocation: INFO and LINK fields, Λ the null link)
- Knuth, The Art of Computer Programming, Vol. 1, 3rd ed., Addison-Wesley, 1997, §2.2.3 (linked allocation: building two lists by relinking, then joining them)
- Knuth, The Art of Computer Programming, Vol. 1, 3rd ed., Addison-Wesley, 1997, §2.2.3 (traversal by LINK)
- Knuth, The Art of Computer Programming, Vol. 1, 3rd ed., Addison-Wesley, 1997, §2.2.3, Exercise 7 (inversion of a linked list, applied to each block)
- Knuth, The Art of Computer Programming, Vol. 1, 3rd ed., Addison-Wesley, 1997, §2.2.3, Exercise 7 (inversion of the second half) and §2.2.3 (relinking); the split is ALG-LIST-010's walk. Proof of the order: catalog/list.md ALG-LIST-013
- Knuth, The Art of Computer Programming, Vol. 1, 3rd ed., Addison-Wesley, 1997, §2.2.3, Exercise 7 (invert a linked linear list by changing its links)
- Knuth, The Art of Computer Programming, Vol. 1, 3rd ed., Addison-Wesley, 1997, §2.2.5 (doubly linked lists with a list head: LLINK, RLINK; O(1) deletion and insertion)
- Knuth, The Art of Computer Programming, Vol. 1: Fundamental Algorithms, 3rd ed., Addison-Wesley, 1997, §1.2.8, eq. (1) and the table of F(n) (definition)
- Knuth, The Art of Computer Programming, Vol. 1: Fundamental Algorithms, 3rd ed., Addison-Wesley, 1997, §1.3.3 (a permutation applied in place cycle by cycle); the rotation's cycles have length 4, or 1 at the centre
- Knuth, The Art of Computer Programming, Vol. 1: Fundamental Algorithms, 3rd ed., Addison-Wesley, 1997, §2.3.1, Algorithm T
- Knuth, The Art of Computer Programming, Vol. 1: Fundamental Algorithms, 3rd ed., Addison-Wesley, 1997, §2.3.1, Algorithm T (traverse binary tree in inorder, with a stack)
- Knuth, The Art of Computer Programming, Vol. 1: Fundamental Algorithms, 3rd ed., Addison-Wesley, 1997, §2.3.1, exercise 12 (preorder analogue of Algorithm T)
- Knuth, The Art of Computer Programming, Vol. 1: Fundamental Algorithms, 3rd ed., Addison-Wesley, 1997, §2.3.1, exercise 13 (postorder with a stack: remember whether the right subtree is done)
- Knuth, The Art of Computer Programming, Vol. 1: Fundamental Algorithms, 3rd ed., Addison-Wesley, 1997, §2.3.1, exercise 7 (preorder and inorder determine the binary tree)
- Knuth, The Art of Computer Programming, Vol. 1: Fundamental Algorithms, 3rd ed., Addison-Wesley, 1997, §2.3.1, exercise 7 (the preorder case; postorder is its mirror read backwards)

Checked: see the builder's catalog section in doc/algorithms_core_2026-11-12/catalog/ (confirmation status per citation; not first-hand items are listed in PROGRESS.md and REPORT.md).
