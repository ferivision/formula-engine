# ADR 0002: Context.Lookup returns (Value, error), not (Value, bool)

**Status:** Accepted
**Date:** 2026-08-12
**Related:** `docs/features/0002-array-aggregate-lookup-function/feature-0002-rfc.md`
§9; issue F2-FASE-2.2

## Context

feature-0001's `Context.Lookup(name string) (registry.Value, bool)`
only needed to distinguish two states: the field is present, or it
isn't (`bool`). Feature-0002 adds a third state: the field is
present, but converting it into an `Array` fails (a mixed-type
`[]any`, per ADR 0001) — and per rfc.md §9, that failure must surface
"at the point the Array is constructed... not deferred to first use,"
i.e. at lookup time, not later when some aggregate/transform function
happens to touch it.

A `bool` return can't represent three states. Options considered:

1. Add a third return value: `(Value, bool, error)`.
2. Keep `(Value, bool)` and swallow conversion failures, deferring the
   error to whatever function first consumes the malformed Array.
3. Change to `(Value, error)`, with a sentinel `errFieldNotFound` for
   the "absent" case and any other error meaning "present but
   invalid."

Option 2 directly contradicts rfc.md §9's "not deferred to first use"
requirement, so it was rejected outright (it's also what F2-FASE-2.1
did as a deliberately temporary stopgap, which this ticket replaces).
Option 1 works but adds a third return value to every call site for a
distinction (`bool` vs. `error`) that's really just "which kind of
not-okay is this," which `error` alone already expresses once there's
a way to tell the two failure kinds apart.

## Decision

`Context.Lookup` returns `(registry.Value, error)`. A field absent
from the data map returns the sentinel `errFieldNotFound`; a field
present but failing Array construction returns that construction
error directly (already a well-formed `*apperror.FormulaError` with
`ErrTypeMismatch`, from `NewArray`); a field present and valid returns
`(value, nil)`.

Callers distinguish the two failure modes with `errors.Is(err,
errFieldNotFound)`. `evaluator.go`'s `NodeIdentifier` case is the only
call site: on `errFieldNotFound` it produces `ErrUndefinedReference`
(FR-6, unchanged behavior); on any other error it returns that error
as-is (the new formula-level `TypeError` behavior); otherwise it
returns the value.

An identifier that already resolves to an `Array` (e.g. a prior
formula's computed result, such as a future `FILTER` output) passes
through `Lookup` unchanged — `NewArray` is never re-invoked on an
already-converted `Array`, since it only recognizes raw
`[]map[string]any`/`[]any` Go slices.

## Consequences

- `Context.Lookup` is used from exactly one place
  (`evaluator.go`'s `NodeIdentifier` case), so this contract change
  had a single call site to update — a much larger surface would
  make this kind of return-type change riskier to reason about.
- Any future direct caller of `Lookup` must use `errors.Is` against
  `errFieldNotFound` rather than treating a non-nil error as
  automatically meaning "undefined reference" — conflating the two
  would silently turn a real `TypeError` into a misleading
  `ErrUndefinedReference`, the same class of bug fixed in
  feature-0001's inherited-failure work (issue #28).
- `errFieldNotFound` is unexported — it's an internal signaling detail
  of this package, not something `internal/evaluator`'s own callers
  need to construct themselves, only compare against.
