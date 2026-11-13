# pregnancy-calculator

## FIELDS (form calform, 2026-11-13)
ccalcmethod: due date (default) | last period | ultrasound | conception date | ivf · date pickers (Month_ID select 0-11, Day_ID select,
Year_ID text) for dueDateDay (default 2027-03-13), cLastPeriodDay, cConceptionDay, cUltrasoundDate, cIVFDay · cCycleLength 22-44 (28) ·
cUltrasoundWeek 22 / cUltrasoundDay 6 · civfage 3|5|6 (5).
Results: current week, weeks+days, due date, trimester, conception date. Uses today's date.

## MATCH
Was PARTIAL MED-OB-001. Built the LMP due date and gestational age.

## RESEARCH
ACOG Committee Opinion No 700: Methods for Estimating the Due Date. Obstet Gynecol 2017;129(5):e150-e154, doi:10.1097/AOG.0000000000002046,
PMID 28426621. Background: "By convention, the EDD is 280 days after the first day of the LMP." (read in an archived copy of the
committee opinion's text). Primary (the guideline itself). ART day counts (266 / 261 / 263) were NOT visible in the copy read — not built.
Cycle-length adjustment: the source says the LMP method does not account for cycle irregularity; no adjustment rule given — not built.

## BUILD
MED-OB-003 DueDateFromLMP, MED-OB-004 GestationalAgeFromLMP — med/obstetric/due_date.go; tests due_date_test.go.

## COMPARE (scratch: checks/scratch/pregnancy-calculator.go)
| inputs | ours | theirs | verdict |
|---|---|---|---|
| their default: EDD 2027-03-13, today 2026-11-13 | LMP 2026-06-06 → EDD 2027-03-13; GA on 2026-11-13 = 22 w 6 d | 22 weeks 6 days, due Mar 13, 2027 | MATCH |
| last period 2026-01-01, cycle 28 (session 2, hidden field cLastPeriodDay) | EDD 2026-10-08 | Oct 8, 2026 | MATCH |
| last period 2024-02-29, cycle 28 | EDD 2024-12-05 | Dec 5, 2024 | MATCH |
| last period 2026-08-01, cycle 28, today 2026-11-13 | EDD 2027-05-08; GA 14 w 6 d | May 8, 2027; 14 weeks 6 days | MATCH |
Session 1 DRIVER gap fixed by the lead: dates go in the hidden MM/DD/YYYY field (here cLastPeriodDay; the fields driver does not list it).
Submissions: session 1: 1 fields + 1 run; session 2: 1 fields + 3 runs (total 4 runs).
