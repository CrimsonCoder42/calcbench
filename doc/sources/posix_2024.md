# POSIX fmod() and remainder()

**Citation**: IEEE Std 1003.1-2024 (The Open Group Base Specifications Issue 8), fmod() and remainder(), DESCRIPTION, https://pubs.opengroup.org/onlinepubs/9799919799/functions/fmod.html and .../remainder.html (read 2026-11-12).

Taken: the definitions of fmod (result has the sign of x, |r| < |y|) and remainder (round-to-nearest quotient). Used by NUM-TRIG-003 (fmod), NUM-TRIG-004, -005 (remainder). Checked: both pages fetched and the quoted definitions matched; interval tests on the normalised angles. Status: confirmed. Kahan 1987 ("Branch Cuts…") was dropped: page range seen only in a search summary, and POSIX covers NUM-TRIG-004.
