# ADR 0001: Detect Record arrays via element inspection, not just concrete Go slice type

**Status:** Accepted
**Date:** 2026-08-12
**Related:** `docs/features/0002-array-aggregate-lookup-function/feature-0002-rfc.md`
§2, §9; issue F2-FASE-1.2

## Context

Feature-0002's `Array`/`Record` value type needs to convert a Go value
passed in `data` into an internal `Array`. The RFC (§2) describes two
native Go shapes a consumer can pass: `[]map[string]any` (Records) and
`[]any` (scalars).

In practice, `[]any` is not exclusively for scalars. It's also the
shape `encoding/json.Unmarshal` produces when decoding a JSON array
into `any`: every element of that `[]any` is itself a `map[string]any`
for each JSON object — never `[]map[string]any` directly. A consumer
building `data` from JSON therefore cannot produce `[]map[string]any`
without an extra manual re-typing step, which works against the PRD's
own goal (§9) of keeping the Record shape "flexible... consistent with
how `data` already works."

## Decision

`NewArray` inspects the elements of a `[]any` input rather than only
trusting the slice's static Go type:

- if every element is `map[string]any`, treat the array as a Record
  array (`IsRecord: true`), converting each element to `Record`
- if every element is a non-map scalar, treat it as a scalar array
  (`IsRecord: false`)
- if elements are a mix of the two, reject with a formula-level
  `TypeError` at construction time, per rfc.md §9's explicit "not
  deferred to first use" rule

`[]map[string]any` remains supported directly as an explicit path —
Go's type system already guarantees homogeneity for that concrete
slice type, so no per-element inspection is needed there; mixing is
only possible (and only checked) on the `[]any` path.

## Consequences

- JSON-decoded data (`[]any` of `map[string]any`) works without
  requiring callers to manually re-type it as `[]map[string]any`.
- The mixed-type check only has real teeth on the `[]any` path. Worth
  remembering if this logic is ever refactored, so the check isn't
  duplicated onto (or accidentally dropped from) the wrong branch.
- A `[]any` containing a mix of Record-shaped and non-Record elements
  always errors, even if the caller intended something else (e.g. a
  scalar array that happens to contain one stray map) — there's no
  partial or best-effort interpretation.
