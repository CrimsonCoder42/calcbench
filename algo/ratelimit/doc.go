// Package ratelimit holds rate meters that decide whether a request conforms to a rate: a token bucket,
// the GCRA leaky bucket, a fixed (jumping) window and a sliding (moving) window log.
//
// ID short form: RATELM (IDs are ALG-RATELM-NNN).
//
// Clock: every call takes now time.Time; no meter reads the clock, starts a goroutine or keeps global
// state. Pass times from one clock (time.Now readings carry the monotonic clock, which Sub and Before use).
// A now earlier than the latest now a meter has seen is treated as that latest time: no time has passed,
// nothing refills or expires, and the call is still decided (doc/algorithms_core_2026-11-12/catalog/ratelimit.md,
// decision R1).
//
// Results: Allow returns (ok, retryAfter). When ok, retryAfter is 0 and the request has been counted.
// When refused, nothing is counted and retryAfter is the shortest wait after which the same request
// would conform if nothing else arrives; it is Never when the request can never conform (n greater than
// the burst or limit). n <= 0 returns (true, 0) and changes nothing. A zero-value meter refuses every
// request with Never; build meters with their New function.
//
// Concurrency: no meter is safe for concurrent use and none locks. A server keeps one meter per key in
// its own map and guards it with its own mutex; the map, the idle sweep and the lock stay the caller's.
//
// Functions:
//
//	ALG-RATELM-001  TokenBucket       token bucket, exact integer refill, n tokens per request (RFC 2697 §3)
//	ALG-RATELM-002  LeakyBucket       GCRA virtual scheduling, one cell per request (ITU-T I.371 Annex A)
//	ALG-RATELM-003  FixedWindow       jumping window: at most limit per consecutive window of length W
//	ALG-RATELM-004  SlidingWindowLog  moving window: at most limit in any span of length W, exact log
//	ALG-RATELM-005  TokenBucket.Restore  give back n tokens taken earlier, capped at the burst
package ratelimit
