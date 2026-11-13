// Package seq transforms range-over-func iterators (iter.Seq and iter.Seq2): mapping, filtering,
// taking a prefix, grouping into chunks and batches, zipping, flattening and sliding windows.
//
// ID short form: SEQ (IDs are ALG-SEQ-NNN).
//
// Laziness: every function returns an iterator and does no work until it is ranged over. Each element is
// pulled from the upstream only when the consumer asks for the next one; nothing is read ahead.
//
// Early stop: when the consumer's loop body stops (break, return), the helper's yield to the upstream
// returns false at once, so the upstream stops too and is never asked for another element (Go
// specification, "For statements with range clause"). Zip, the one helper that holds a second upstream
// by iter.Pull, always calls that upstream's stop function (deferred), on early stop, on exhaustion and
// on a panic in the loop body.
//
// No goroutines, no channels: these are pure functions over iter.Seq. Channel pipelines (generators,
// or-done, tee, fan-out/fan-in, errgroup) own goroutines and are not library material; see the catalog.
//
// Buffer reuse: Chunk, Batch and Window yield ONE slice, rewritten between steps (no allocation per
// group once the buffer has grown). The slice is valid until the loop body returns; copy it to keep it
// (slices.Clone) and never modify it.
//
// Sizes: a function taking a bare count or size returns (iter.Seq, error), so an invalid count is an
// error, not silence. A nil iterator is the empty sequence. A nil function makes the empty sequence
// (Map, Filter, Batch's endsBatch excepted: a nil endsBatch never ends a batch early).
//
// Functions:
//
//	ALG-SEQ-001  Map      each element transformed by a function
//	ALG-SEQ-002  Filter   the elements a predicate keeps
//	ALG-SEQ-003  Take     the first n elements, then stop the upstream
//	ALG-SEQ-004  Chunk    consecutive groups of exactly size elements (the last may be shorter)
//	ALG-SEQ-005  Zip      pairs from two iterators, stopping at the shorter (iter.Pull, stop always called)
//	ALG-SEQ-006  Flatten  the elements of a sequence of sequences, in order
//	ALG-SEQ-007  Batch    groups closed at a size limit OR when a predicate marks an element as the last
//	ALG-SEQ-008  Window   every run of width consecutive elements (sliding window), reused ring buffer
package seq
