# due-date-calculator

## FIELDS (form calform, 2026-11-13)
ccalcmethod last period (default) | ultrasound | conception date | ivf · cLastPeriodDay (default 2026-09-14) · cCycleLength 28 ·
cConceptionDay (2026-09-28) · cUltrasoundDate (2026-11-13), cUltrasoundWeek 22, cUltrasoundDay 6 · cIVFDay (2026-09-19), civfage 3|5|6.
Results: due date, current week/weeks+days (uses today), trimester, likely conception date.

## MATCH
HAVE MED-OB-003 DueDateFromLMP, MED-OB-004 GestationalAgeFromLMP (built in row 13 this session).
Not built: conception / ultrasound / IVF methods and the cycle-length adjustment — the ACOG CO 700 copy read does not state their day
counts (NEEDS SOURCE, see pregnancy-calculator.md).

## COMPARE (scratch: checks/scratch/due-date-calculator.go)
| inputs | ours | theirs | verdict |
|---|---|---|---|
| LMP 2026-09-14 (their default), cycle 28, today 2026-11-13 | EDD 2027-06-21; GA 8 w 4 d | Jun 21, 2027; 8 weeks 4 days | MATCH |
| LMP 2023-12-31, cycle 28 (session 2, hidden field cLastPeriodDay) | EDD 2024-10-06 | Oct 6, 2024 | MATCH |
| LMP 2026-07-04, cycle 32, today 2026-11-13 | EDD 2027-04-10; GA 18 w 6 d | Apr 14, 2027; 18 weeks 2 days | THEIRS DIFFERS |

## DISCREPANCY: THEIRS DIFFERS (cycle length)
With a 32-day cycle they move the due date 4 days later (+ (cycle − 28) days) and the gestational age 4 days earlier. ACOG CO 700
(Obstet Gynecol 2017;129:e150, Background) gives LMP + 280 days and says the LMP method does not account for cycle irregularity; it
states no cycle-length adjustment. Ours (MED-OB-003/004) takes no cycle length and stays. At 28 days the two agree.
Submissions: session 1: 1 fields + 1 run; session 2: 2 runs (total 3 runs).
