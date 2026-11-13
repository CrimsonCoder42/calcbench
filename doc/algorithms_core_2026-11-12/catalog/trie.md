# algo/trie (SHORT = TRIE) — new package (builder W2-B7)

Planned in doc/taxonomy.md:20 (`algo/trie`, TRIE, "prefix trees"). Files: doc.go, trie.go, autocomplete_top_k.go,
wildcard_search.go, replace_words_by_root.go, prefix_suffix_index.go, each with its _test.go. Imports algo/heap
(AutocompleteTopK). New error: `ErrNegativeK` (autocomplete_top_k.go:12). Exported types: `Trie[V]`
(trie.go:17), `Entry[V]` (autocomplete_top_k.go:15), `PrefixSuffixIndex` (prefix_suffix_index.go:18), constant
`Wildcard = '.'` (wildcard_search.go:6).

## Sources (confirmed)
- Fredkin, "Trie memory", Communications of the ACM, 1960, 3(9):490–499, doi:10.1145/367390.367400 —
  metadata confirmed on Crossref 2026-11-12 (paper itself not read).
- Knuth, TAOCP Vol. 3: Sorting and Searching, 2nd ed., Addison-Wesley, 1998, §6.3 "Digital Searching",
  pp. 492–512 (section and pages confirmed via a search listing; book not opened in this run).
- Manning, Raghavan, Schütze, *Introduction to Information Retrieval*, Cambridge UP, 2008, §3.2.1 "General
  wildcard queries", Permuterm indexes — READ FIRST-HAND in the 2007 preliminary draft (§3.1.1, pp. 48–49):
  "the dictionary consists of all rotations of each term (with the $ terminating symbol appended)"; the query
  m*n is rotated to n$m* and looked up as a prefix, finding man, men and (via n$moro) moron. The same page
  describes the losing alternative: a normal and a reverse B-tree, intersecting the prefix set W and the suffix
  set R. Section number in the 2008 print edition is §3.2.1 (TOC); **lead: confirm the print page**.
- CLRS 3rd ed. §6.5 for the bounded priority queue in AutocompleteTopK.

## Design decisions (each with the alternative that lost)
- **Children**: sorted byte labels + child pointers per node, binary search (≤ 8 comparisons, σ ≤ 256).
  Lost: `[256]*node` per node (O(1) step but 2 KiB per node: 2 GiB for 1e6 nodes); a map per node (no order,
  so every walk would sort). Sorted children give byte-order iteration for free → every output deterministic.
- **Byte semantics** everywhere (keys, '.', prefix/suffix). Lost: rune semantics (needs UTF-8 decoding on
  every step and an invalid-UTF-8 policy); a caller wanting runes can store one key per normalised string.
- **Iteration**: `iter.Seq2[string, V]` for WithPrefix and WildcardSearch (lazy, stops early); iterative walks
  with an explicit stack, so a 1 MB key never overflows the goroutine stack.
- **Zero value usable**; nil *Trie acts empty (Insert stores nothing), as Heap.
- **ReplaceWordsByRoot ignores empty roots** (ASSUMPTION: "" is a prefix of everything and would erase every
  word). **Lead: log in doc/ASSUMPTIONS.md** or tell me the other policy.

### ALG-TRIE-001…006 Trie.Insert / Search / HasPrefix / Delete / Len / WithPrefix (trie.go:65, 96, 112, 131, 173, 192)
- **Best known**: Θ(m) per operation for a key of m bytes (must read the key); our O(m·log σ) = O(m) with
  σ ≤ 256. WithPrefix O(m + s), s = nodes below the prefix (output-sensitive up to key lengths).
  Lost: sorted slice + binary search (O(m log n) search, O(n) insert); hash map (no prefix queries).
- **Delete**: one walk down remembering the deepest ancestor that must stay (root, holds a key, or has
  another child); cut that one edge. O(m), O(1) extra; no parent pointers, no recursion. Invariant: every
  non-root node leads to a key, so HasPrefix is "node exists". Tests check node counts after pruning
  (trie_test.go TestTrie_deletePrunes, and the oracle compares node count with the number of distinct
  prefixes of the remaining keys).
- **Worked example**: Knuth's 31 common English words (§6.2.2/§6.3); expected results derived by hand
  (the list was not re-checked against the book; noted in trie_test.go).
- **Errors**: none. Input not modified. **Tags**: trie, prefix-tree. LeetCode 208.

### ALG-TRIE-007 Trie.AutocompleteTopK (autocomplete_top_k.go:38)
`func (t *Trie[V]) AutocompleteTopK(dst []Entry[V], prefix string, k int, rank func(a, b V) int) ([]Entry[V], error)`
- **Decision**: ranked by a caller comparator on the payload (cmp.Compare shape), ties broken by byte order
  of the key (smaller first); rank nil = byte order alone with early stop after k. Lost: a `weight func(V)
  float64` (forces floats, NaN policy) and storing per-node precomputed top-k lists (O(1) query, but
  O(k·depth) extra per insert and fixed k).
- **Best known**: O(m + s log k) with a k-bounded heap over the subtree; O(m + k·L) is possible only with
  per-node cached lists (lost, above). Because the walk is in byte order, a later key loses every tie, so it
  enters the heap only if strictly better: one rank call per key and a string allocation only on entry.
- **Errors**: ErrNegativeK (k < 0). k = 0 → nothing. Huge k never sizes the heap (hint min(k, Len)).
- **Tags**: autocomplete, typeahead, top-k. LeetCode 1268, 642.

### ALG-TRIE-008 Trie.WildcardSearch (wildcard_search.go:23)
- '.' (exported `Wildcard`) matches exactly one BYTE; no escape (a literal '.' in a key is matched by '.'
  anyway); match is exact length. Iterator in byte order; existence = range once and stop.
- **Best known**: O(m) without '.', O(nodes at depth ≤ m) worst — the trie walk is optimal among
  index-free methods. Lost: scanning all keys (O(total bytes)); recursion (depth = m).
- **Tags**: wildcard, dont-care, partial-match. LeetCode 211.

### ALG-TRIE-009 Trie.ShortestPrefixOf, ALG-TRIE-010 ReplaceWordsByRoot (replace_words_by_root.go:15, 47)
- ShortestPrefixOf stops at the first node holding a key: O(p), returns word[:p] (no allocation, tested).
- ReplaceWordsByRoot(dst, roots, words) builds a Trie[struct{}] of roots: O(R + W). Lost: for each word,
  test every root (O(W·|roots|)); hash set of roots probing each prefix length (O(Σ L²) hashing).
- **Tags**: replace-words, stemming. LeetCode 648.

### ALG-TRIE-011…013 NewPrefixSuffixIndex / Lookup / LookupLast (prefix_suffix_index.go:49, 117, 145)
- **Method**: permuterm (Manning et al. §3.2.1) with one change, proved here: the book stores rotations
  w[j:] $ w[:j] and answers X*Y by the prefix Y$X, which requires len(w) ≥ |X|+|Y| (no overlap). We store
  w[j:] $ w (the whole word after $), so for a query (p, s): a key s' $ w has prefix s $ p ⇔ s' = s (the
  separator, symbol 256, is not a byte, so it pins the end of s) and w starts with p; s' = w[j:] = s ⇔ s is a
  suffix of w. Hence matches = {w : HasPrefix(w,p) ∧ HasSuffix(w,s)}, overlap allowed, and each match has
  exactly one key below the query node (j = len(w) − len(s)), so no deduplication.
  Same O(Σ L²) space as the rotations. The separator is a 257th symbol (uint16 labels), not '#' or '$', so
  words containing any byte are safe (tested).
- **Best known / lost**: query O(p + q) for LookupLast (each node keeps the largest word index through it);
  Lookup O(p + q + s + r log r). Lost: two tries (prefix and reversed) and intersecting W ∩ R — the book's
  own alternative — costs O(|W| + |R|) per query even when the answer is empty; brute force O(Σ L) per query.
- **Worked example**: the book's m*n → man, men, moron (prefix_suffix_index_test.go).
- **Errors**: none; nil/zero index is empty; LookupLast returns -1 when nothing matches.
- **Tags**: permuterm, prefix-and-suffix, word-filter. LeetCode 745.

## Consumer evidence (read-only)
- /Users/devinanderson/Desktop/red_lite_cache/internal/cache/ordered.go:209-231 `appendPrefix`/`cutPrefix`:
  a prefix group walk over a sorted B+-tree of keys. Trie.WithPrefix gives the same ordered prefix
  enumeration, but **does not replace it**: the cache's ordered tree also serves range scans, ordered
  eviction and cut-in-place of a prefix run; a trie would be a second index. A Trie-backed DeletePrefix
  would need a "cut subtree" method (not built; would be a new ID if asked).
- /Users/devinanderson/Desktop/red_lite_cache/internal/sheet/sweeplog.go:17-24 `covers()`: checks one key
  against one logged prefix or suffix (strings.HasPrefix / HasSuffix). One check per entry is already O(len);
  **not replaced**. If the log grows large and the question becomes "which logged suffix-deletes cover this
  key", a Trie of reversed suffixes with ShortestPrefixOf (or a walk) answers it in O(len(key)).
- a consumer's cache key-prefix constants (cache/keys.go:60),
  :159 and :311 (two more prefix constants): a ':'-delimited key grammar swept by DeletePrefix. Motivates
  HasPrefix/WithPrefix (find what a sweep would cover before running it); no code change proposed.
