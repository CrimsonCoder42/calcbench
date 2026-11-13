# Bentley, Programming Pearls

**Citation**: Bentley, *Programming Pearls*, 2nd ed., Addison-Wesley, 2000.

| Taken | Locator | Used by | How checked |
|---|---|---|---|
| Binary search whose invariant brackets the first occurrence (x[l] < t <= x[u]) | Column 9.3 | ALG-SEARCH-001…004 | Oracle tests; locator NOT verified against a copy here (TODO.md) |


## Used by S10 (algorithms core, 2026-11-12)

IDs: ALG-DP-022, ALG-DP-023, ALG-HASH-010, ALG-HASH-011, ALG-MATRIX-006, ALG-PREFIX-001, ALG-PREFIX-002, ALG-PREFIX-009, ALG-PREFIX-011, ALG-SEARCH-005, ALG-SEARCH-006, ALG-SEARCH-008, ALG-TEXT-010, ALG-TEXT-013.

Citations as written in the doc comments:
- Bentley, Programming Pearls, 2nd ed., Addison-Wesley, 2000, Column 2, problem C (anagram signatures, §2.8, pp. 18–20)
- Bentley, Programming Pearls, 2nd ed., Addison-Wesley, 2000, Column 2, problem C (anagrams share a signature, §2.4 and §2.8)
- Bentley, Programming Pearls, 2nd ed., Addison-Wesley, 2000, Column 2, problem C (signatures, §2.8, pp. 18–20)
- Bentley, Programming Pearls, 2nd ed., Addison-Wesley, 2000, Column 2.3 (reversal: reverse the whole, then each piece, so the pieces swap order and each reads forwards)
- Bentley, Programming Pearls, 2nd ed., Addison-Wesley, 2000, Column 8 "Algorithm Design Techniques", §8.2 (sum of x[i..j] = cumarr[j] − cumarr[i−1])
- Bentley, Programming Pearls, 2nd ed., Addison-Wesley, 2000, Column 8 "Algorithm Design Techniques", §8.2 (the cumulative array cumarr)
- Bentley, Programming Pearls, 2nd ed., Addison-Wesley, 2000, Column 8 "Algorithm Design Techniques", §8.4 "A Scanning Algorithm" (Column 7 in the 1st ed., 1986)
- Bentley, Programming Pearls, 2nd ed., Addison-Wesley, 2000, Column 8, Problems (the two-dimensional maximum-sum subarray, reduced to one dimension over pairs of rows)
- Bentley, Programming Pearls, 2nd ed., Addison-Wesley, 2000, Column 8, §8.4 (a trade's profit as a subvector sum of day-to-day changes; the unlimited case is the sum of the positive changes, proved above and in catalog/dp.md)
- Bentley, Programming Pearls, 2nd ed., Addison-Wesley, 2000, Column 8, §8.4 (the scanning algorithm for the maximum-sum subvector; prices[s] − prices[b] is the sum of the day-to-day changes b+1 … s, so the best trade is the maximum subvector of the changes, and the scan's running state reduces to the lowest price so far)
- Bentley, Programming Pearls, 2nd ed., Addison-Wesley, 2000, Column 9.3 (first position where a monotone test holds)
- Bentley, Programming Pearls, 2nd ed., Addison-Wesley, 2000, Column 9.3 (first position where a monotone test holds: here rotated[i] < rotated[0])
- Bentley, Programming Pearls, 2nd ed., Addison-Wesley, 2000, Column 9.3 (the invariant that brackets the first position where the test holds)

Checked: see the builder's catalog section in doc/algorithms_core_2026-11-12/catalog/ (confirmation status per citation; not first-hand items are listed in PROGRESS.md and REPORT.md).
