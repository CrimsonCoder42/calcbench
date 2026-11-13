# Age Calculator (row 1)

CLASSIFY: in scope (calendar arithmetic, proleptic Gregorian).

FIELDS (form calcf): today_Month_ID select 0..11 (default 10), today_Day_ID select 1..31 (default 13),
today_Year_ID text (default 2026) = date of birth; ageat_Month_ID / ageat_Day_ID / ageat_Year_ID = "age at" date
(default today). Results: years/months/days, and totals in months, weeks, days, hours, minutes, seconds.

MATCH: NONE in INDEX.md (only num/series local-day windows, NUM-SERIES-004..008, which take time.Time instants).

RESEARCH: Meeus J. Astronomical Algorithms, 2nd ed., Willmann-Bell, 1998, ch. 7 "Julian Day": eq. 7.1 (JD of a
date), the inverse method, day of week ((JD+1.5) mod 7), day of year, interval = difference of JDs.
Worked examples reproduced: 7.a (1957-10-04 → JD 2436116.31), 7.c (inverse), 7.d (1910-04-20 → 1986-02-09 =
27689 d), 7.e (1954-06-30 Wednesday), 7.f (1978-11-14 day 318; 1988-04-22 day 113).
Not built: the "X years Y months Z days" split. It is a convention (how to borrow days from the month before;
what a 29 February birthday's anniversary is in a common year) with no primary source found in this session;
listed for Devin in LEAD.md.

BUILD: sci/astronomy, SCI-ASTRO-001..006 (GregorianToJulianDayNumber, JulianDayNumberToGregorian,
GregorianDayOfWeek, GregorianDayOfYear, DaysBetweenDates, GregorianAddDays). The total-days figure of an age is
DaysBetweenDates.

COMPARE (ours: checks/scratch/age-calculator.go; round trip against Go's time package every 97 days over
JDN 38..366,963,925: 0 mismatches):

| inputs | ours | theirs | diff | verdict |
|---|---|---|---|---|
| 1910-04-20 → 1986-02-09 (Meeus 7.d) | 27689 d | 27,689 d (75 y 9 m 20 d) | 0 | MATCH |
| born 1990-05-15, at 2026-11-13 | 13331 d | 13,331 d (lead's run, out/run_age_driverfix.json) | 0 | MATCH |
| born 2000-02-29, at 2025-03-01 (common year) | 9132 d | 9,132 d (25 y 0 m 1 d) | 0 | MATCH |
| born 2000-02-29, at 2028-02-29 | 10227 d | 10,227 d (28 y 0 m 0 d) | 0 | MATCH |

Their Y/M/D split, observed: a 29 Feb birthday's anniversary in a common year is 28 Feb (2025-03-01 is
"25 years 0 months 1 day"). That is one of the conventions in the LEAD.md decision; ours does not build the split.
Session 1: 1 submission (driver gap, since fixed). Session 2 (driver: hidden today/ageat MM/DD/YYYY fields):
lead 1 + 3. Total 5.

CLOSE: BUILT SCI-ASTRO-001..006, compare MATCH (the Y/M/D split waits on Devin).
