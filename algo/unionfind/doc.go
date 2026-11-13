// Package unionfind holds disjoint-set (union-find) structures: a DisjointSet over the integer elements
// 0..n-1 with union by rank and path compression, and the online count of islands on a grid built on it.
//
// ID short form: UNION (IDs are ALG-UNION-NNN).
//
// Conventions: elements are ints 0..Len()-1; an element outside that range is ErrOutOfRange, never a
// panic. Find compresses paths, so a DisjointSet is not safe for concurrent use. Grid positions are Cell
// values (Row, Col from the top-left, counting from 0); algo/graph uses the same type. No function modifies
// its input slices.
//
// Functions:
//
//	ALG-UNION-001  NewDisjointSet  n singleton sets {0} … {n-1}
//	ALG-UNION-002  DisjointSet.MakeSet  add one new singleton element
//	ALG-UNION-003  DisjointSet.Find  the representative of an element's set (with path compression)
//	ALG-UNION-004  DisjointSet.Union  merge two sets (union by rank), reporting whether they were separate
//	ALG-UNION-005  DisjointSet.Connected  whether two elements share a set
//	ALG-UNION-006  DisjointSet.Count  the number of sets
//	ALG-UNION-007  DisjointSet.Len  the number of elements
//	ALG-UNION-008  IslandsOnline  the island count after each land addition on a grid
package unionfind
