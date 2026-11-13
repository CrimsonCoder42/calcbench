// Package trie holds prefix trees (tries) over string keys: a map with prefix queries, ranked
// autocompletion, wildcard matching, shortest-root replacement, and a permuterm index for
// prefix-and-suffix queries.
//
// ID short form: TRIE (IDs are ALG-TRIE-NNN).
//
// Keys are bytes: a key of m bytes is a path of m edges, each labelled by one byte, and every function
// that matches or counts characters matches or counts bytes, not runes (a UTF-8 character of 2–4 bytes
// is 2–4 positions). Children are kept in ascending byte order, so every iteration is in byte-wise
// lexicographic order, the order of Go's < on strings, and every result is deterministic.
// The zero Trie is empty and ready to use; methods on a nil *Trie act as on an empty one (Insert stores
// nothing). No function modifies its arguments. Walks are iterative, so long keys never exhaust the
// stack. Not safe for concurrent writes.
//
// Functions:
//
//	ALG-TRIE-001  Trie.Insert                 store a value under a key; reports a replacement
//	ALG-TRIE-002  Trie.Search                 the value under a key
//	ALG-TRIE-003  Trie.HasPrefix              whether any key starts with a prefix
//	ALG-TRIE-004  Trie.Delete                 remove a key, pruning branches that lead to no key
//	ALG-TRIE-005  Trie.Len                    the number of keys
//	ALG-TRIE-006  Trie.WithPrefix             iterator over the keys with a prefix, in byte order
//	ALG-TRIE-007  Trie.AutocompleteTopK       the k best keys with a prefix by a caller's rank, ties by byte order
//	ALG-TRIE-008  Trie.WildcardSearch         iterator over the keys matching a pattern where '.' is any one byte
//	ALG-TRIE-009  Trie.ShortestPrefixOf       the shortest key that is a prefix of a word
//	ALG-TRIE-010  ReplaceWordsByRoot          each word replaced by its shortest root prefix
//	ALG-TRIE-011  NewPrefixSuffixIndex        a permuterm index of words for prefix-and-suffix queries
//	ALG-TRIE-012  PrefixSuffixIndex.Lookup    the indexes of the words with a prefix and a suffix, ascending
//	ALG-TRIE-013  PrefixSuffixIndex.LookupLast  the largest such index
package trie
