// Package backoff holds retry-delay schedules: capped exponential backoff and Brooker's three jittered
// variants (full, equal and decorrelated jitter).
//
// ID short form: BACKOF (IDs are ALG-BACKOF-NNN).
//
// Randomness: no function draws a random number. The caller passes the draw u ∈ [0, 1) (for example
// rand.Float64 from its own seeded source), so every function is pure: same inputs, same output. No
// function reads the clock, sleeps, starts a goroutine, keeps state or allocates. Durations are exact
// integer nanoseconds; u is scaled in integers (⌊d·⌊u·2⁵³⌋ / 2⁵³⌋ through a 128-bit product), never by a
// float multiply of d.
//
// Functions:
//
//	ALG-BACKOF-001  Exponential         min(max, base·2^attempt), saturating, no float
//	ALG-BACKOF-002  FullJitter          a draw in [0, sleep), sleep = Exponential
//	ALG-BACKOF-003  EqualJitter         sleep/2 plus a draw in [0, sleep/2): a value in [⌈sleep/2⌉, sleep]
//	ALG-BACKOF-004  DecorrelatedJitter  min(max, a draw in [base, 3·prev)), from the previous sleep
package backoff
