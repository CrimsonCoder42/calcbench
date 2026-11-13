# Date Calculator (row 2)

CLASSIFY: in scope for the calendar arithmetic. The business-day / holiday options (useholiday, hd* checkboxes:
US federal holidays) are OUT OF SCOPE: a jurisdiction's holiday table.

FIELDS: form calcf (between two dates): today_*/ageat_* date selects + year text, addendday checkbox (include end
day, default false), useholiday radio, holiday checkboxes hdny..hdne, custom holidays n0..n19/m0..m19/d0..d19.
Form calcf2 (add/subtract): today2_Month_ID/Day_ID/Year_ID (default today), c2op select +/- (default +),
c2year/c2month/c2week/c2day text (default 0), usebusinessday checkbox (false), holiday options as calcf.

MATCH: NONE before this run; now HAVE SCI-ASTRO-005 DaysBetweenDates, SCI-ASTRO-006 GregorianAddDays,
SCI-ASTRO-003 GregorianDayOfWeek (built under row 1). Adding years/months (calendar months, end-of-month
clamping) is a convention not built; see LEAD.md.

COMPARE (ours: checks/scratch/date-calculator.go; start date is their default 2026-11-13 because the driver
cannot set the today2_* date selects):

| inputs | ours | theirs | diff | verdict |
|---|---|---|---|---|
| 2026-11-13 + 27689 d | 2102-09-05 Tuesday | Sep 5, 2102, Tuesday | 0 | MATCH |
| 2026-11-13 − 1000 d | 2024-02-17 Saturday | Feb 17, 2024, Saturday | 0 | MATCH |
| 2026-11-13 + 36525 d | 2126-11-14 Thursday | Nov 14, 2126, Thursday | 0 | MATCH |

DRIVER: today2_Month_ID, today2_Day_ID, today2_Year_ID reported missing (date-picker fields not fillable).
Submissions: 3.

CLOSE: MATCH (SCI-ASTRO-003, -005, -006).
