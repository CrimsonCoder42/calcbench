# Day Counter (row 8)

CLASSIFY: in scope (count of days between two dates). Business-day/holiday options out of scope (holiday table).
FIELDS (calcf): today_Month_ID/Day_ID/Year_ID (start), ageat_Month_ID/Day_ID/Year_ID (end), addendday checkbox
(include end day, default false), useholiday radio, holiday checkboxes hdny..hdne.
MATCH: HAVE SCI-ASTRO-005 DaysBetweenDates (end day excluded; +1 to include it, as their addendday).
COMPARE (session 2, hidden today/ageat fields; ours: checks/scratch/day-counter.go and day-counter-weekday.go):

| inputs | ours | theirs | diff | verdict |
|---|---|---|---|---|
| 1910-04-20 → 1986-02-09 | 27689 d; Wed → Sun | 27689 d; Wednesday → Sunday | 0 | MATCH |
| 1899-12-31 → 1900-03-01 (1900 not leap) | 60 d; Sun → Thu | 60 d; Sunday → Thursday | 0 | MATCH |
| 2024-01-01 → 2024-12-31, end day included | 365 + 1 = 366 d; Mon → Tue | 366 d "including both days" | 0 | MATCH |

Submissions: 3 (session 2).
CLOSE: MATCH (SCI-ASTRO-005, SCI-ASTRO-003); reopened from BLOCKED.
