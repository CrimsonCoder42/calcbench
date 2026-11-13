// Package text matches, hashes, splits and validates strings.
//
// ID short form: TEXT (IDs are ALG-TEXT-NNN).
//
// Bytes, not runes, unless stated: the matching, hashing and palindrome functions are generic over
// S ~string | ~[]byte, compare bytes, and report byte offsets, so one implementation serves string and
// []byte callers with no conversion. On UTF-8 text a byte offset is where the match starts in the
// encoded string, ready to slice with. Rune semantics are used only where bytes would change the
// answer's meaning: IsAnagram compares rune multisets, and WildcardMatch's '?' and RegexMatch's '.'
// match one rune; there, bytes that are not valid UTF-8 each read as U+FFFD (the Go conversion rule).
// ReverseWords and StringToInt split on ASCII white space (space, \t, \n, \v, \f, \r), which is UTF-8
// safe: no byte of a multi-byte rune is ASCII.
//
// Allocation: searches allocate only their pattern tables (KMPSearch: the prefix function, one int per
// pattern byte; RabinKarpSearch: nothing). Functions that produce many indices append them to dst and
// return the extended slice, like strconv.AppendInt: pass nil to allocate, or buf[:0] to reuse a buffer.
// On error dst is returned with its length unchanged. No function modifies its input.
//
// Errors are package sentinels wrapped with the offending parameter; test with errors.Is.
//
// Not built: StrStr (the first match only) is strings.Index / bytes.Index. Roman numerals and plain
// dotted version numbers are held until a primary source for their rules is confirmed.
//
// Functions:
//
//	ALG-TEXT-001  PrefixFunction               longest proper border of every prefix (KMP failure function)
//	ALG-TEXT-002  KMPSearch                    every start of a pattern in a text, Knuth–Morris–Pratt
//	ALG-TEXT-003  ZFunction                    longest common prefix of the string and each suffix
//	ALG-TEXT-004  NewRollingHash               polynomial hash of a sliding window (RollingHash: Append, Roll, Sum, Len, Reset)
//	ALG-TEXT-005  RabinKarpSearch              every start of a pattern, by fingerprints, each hit verified
//	ALG-TEXT-006  RepeatedSubstrings           first starts of the length-L substrings that occur twice or more
//	ALG-TEXT-007  LongestCommonPrefix          length of the prefix shared by every string
//	ALG-TEXT-008  PalindromeRadii              longest palindrome at every centre, Manacher
//	ALG-TEXT-009  LongestPalindromicSubstring  the longest palindromic substring
//	ALG-TEXT-010  IsAnagram                    whether two strings hold the same runes
//	ALG-TEXT-011  RunLengthEncode              PackBits run-length encoding
//	ALG-TEXT-012  RunLengthDecode              PackBits run-length decoding, with an output cap
//	ALG-TEXT-013  ReverseWords                 the words in reverse order, single-spaced
//	ALG-TEXT-014  MinimumWindowSubstring       shortest window holding every byte of a target
//	ALG-TEXT-015  StringToInt                  strtol-style base-10 prefix parse, overflow is an error
//	ALG-TEXT-016  IsJSONNumber                 whether a string is a JSON number (RFC 8259 §6)
//	ALG-TEXT-017  ColumnLabel                  spreadsheet column label of a number (bijective base 26)
//	ALG-TEXT-018  ColumnNumber                 number of a spreadsheet column label
//	ALG-TEXT-019  CompareSemVer                Semantic Versioning 2.0.0 precedence of two versions
//	ALG-TEXT-020  WildcardMatch                whole-text match with '?' and '*'
//	ALG-TEXT-021  RegexMatch                   whole-text match with '.' and '*'
package text
