# Changelog

## Unreleased

### Phase 1 — Foundation & Public API Skeleton

- Initialized Go module (`go 1.25`) and added `internal/apperror` with
  the `FormulaError` type and `ErrorCode` constants from `rfc.md` §11
  (closes #1).
- Added the public API skeleton — `Evaluate`, `FormulaInput`, `Result`
  per `rfc.md` §2, currently a stub (closes #2). Phase 1 complete.
