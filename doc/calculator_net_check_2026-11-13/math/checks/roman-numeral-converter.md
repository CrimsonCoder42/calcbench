# roman-numeral-converter (A-MATH, 2026-11-13, HANDOFF from A-OTHER)

CLASSIFY: in scope (integer ↔ Roman numeral conversion; home would be algo/intmath next to the other digit-string conversions, ALG-INTMTH).
FIELDS: not loaded (research first; no point loading the page with no source).
MATCH: NONE (no Roman numeral function in INDEX.md).
RESEARCH (primary or standard reference for the composition rules: the seven symbol values, the six subtractive pairs IV IX XL XC CD CM, at most three repeats, range 1–3999):
- The Unicode Standard, Version 15.0, ch. 22 "Symbols", §22.3 "Numerals", "Roman Numerals" (read first-hand 2026-11-13 from unicode.org/versions/Unicode15.0.0/ch22.pdf): says Roman numerals are preferably composed from Latin letters and that U+2150..U+218F hold compatibility forms; it gives NO composition rules. The Unicode Character Database numeric values for U+2160..U+216F would source the single-symbol values (Ⅰ 1, Ⅴ 5, Ⅹ 10, Ⅼ 50, Ⅽ 100, Ⅾ 500, Ⅿ 1000, and the precomposed Ⅳ 4, Ⅸ 9) but not XL, XC, CD, CM or the repeat limit.
- Cajori, A History of Mathematical Notations, Vol. 1, Open Court, 1928, pp. 30–37 (Roman notation): located by search only, not read first-hand (full-text download not permitted in this session). It is a historical survey of many variant usages, and its claim about the late origin of subtractive forms is disputed, so it may not fix one "standard form" anyway.
- Encyclopedia of Mathematics, "Roman numerals": server error (502), not read.
Outcome: NEEDS SOURCE. Nothing built. For Devin: (a) accept Cajori Vol. 1 §§ on Roman notation (or Ifrah, The Universal History of Numbers, Wiley 2000, the Roman numerals chapter) once someone reads the page and confirms it states the six subtractive pairs, or (b) decide the "standard form" convention (1–3999, six subtractive pairs, ≤ 3 repeats, strict parsing) as a D-id with the UCD symbol values as the cited part.
