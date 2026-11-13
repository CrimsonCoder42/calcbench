// Package cron parses POSIX crontab time specifications and finds the next instant a specification
// fires, in any time zone, with daylight-saving transitions handled by a stated rule.
//
// ID short form: CRON (IDs are ALG-CRON-NNN).
//
// Syntax (IEEE Std 1003.1-2024, "crontab", INPUT FILES): exactly five fields separated by blanks (space or
// tab): minute 0–59, hour 0–23, day of month 1–31, month 1–12, day of week 0–6 (0 = Sunday). Each field is
// `*` or a comma-separated list of elements, an element being a number or a range lo-hi with lo ≤ hi. Steps
// (`/`), names, 7 for Sunday, `?`, `L`, `@hourly` and a seconds field are not POSIX and are ErrSyntax; a
// number outside its field's range or lo > hi is ErrRange. A specification no day can ever satisfy (a day of
// month that no allowed month has, such as "0 0 30 2 *") is ErrNeverFires.
//
// Day matching is the common-cron rule: the month always restricts. When day of month and day of week are
// both not literally `*`, a day matches when either matches; when only one is restricted it alone decides;
// when both are `*`, every day matches. This agrees with POSIX's own example "0 0 1,15 * 1" (the 1st and 15th
// plus every Monday) but not with POSIX's literal wording when the month is also restricted (POSIX would
// then OR the day of week with month-and-day-of-month); see doc/ASSUMPTIONS.md A-068.
//
// Daylight saving (CronSpec.Next): a matching wall time that does not exist (spring-forward gap) fires at
// the first instant after the gap (the transition instant, e.g. 03:00 EDT); several matching times inside
// one gap fire once. A matching wall time that occurs twice (fall-back overlap) fires once, at its first
// (earlier-offset) occurrence; the second occurrence never fires.
//
// Clock: nothing here reads the clock; Next takes the instant to search after.
//
// Functions:
//
//	ALG-CRON-001  ParseCron     a five-field POSIX crontab time specification parsed into a CronSpec
//	ALG-CRON-002  CronSpec.Next the earliest instant strictly after a time at which the spec fires
package cron
