# Grade Calculator (row 38)

CLASSIFY: in scope (weighted average of scores). The letter for a score (A/B/C cut-offs) and letter-to-percent
inputs are an institution's scale: out of scope. The "what do I need on the final" solve (required =
(target − current·(1 − w_final)) / w_final) is algebra on the same weighted mean; not built (home is a lead
question: num/stats variant or life/everyday; LEAD.md).
FIELDS (gradeform): rows i = 1..n: d<i> name, s<i> score (number or letter), l<i> letter select (a+..f), w<i>
weight; defaults Homework 1 90 w5, Project B w20, Midterm exam 88 w20.
MATCH: HAVE NUM-STATS-026 WeightedMean.
COMPARE (ours: checks/scratch/grade-calculator.go):

| inputs | ours | theirs | diff | verdict |
|---|---|---|---|---|
| 90 w5, 80 w20, 88 w20 | 84.6666666667 | 84.7 | 0 at 1 dp | MATCH |

Submissions: 1.
CLOSE: MATCH (NUM-STATS-026; letter scale out of scope).
