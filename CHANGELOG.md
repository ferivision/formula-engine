# Changelog

## Unreleased

### Phase 1 — Foundation & Public API Skeleton

- Initialized Go module (`go 1.25`) and added `internal/apperror` with
  the `FormulaError` type and `ErrorCode` constants from `rfc.md` §11
  (closes #1).
- Added the public API skeleton — `Evaluate`, `FormulaInput`, `Result`
  per `rfc.md` §2, currently a stub (closes #2). Phase 1 complete.

### Phase 2 — Tokenizer

- Added `internal/parser` with `Tokenize`, covering numbers,
  identifiers, operators, and punctuation, with line/column tracking
  (closes #3).
- Added string literal tokenizing (with escaped quotes) and
  `SyntaxError`-backed reporting for unterminated strings and
  unrecognized characters, satisfying FR-7 (closes #4). Phase 2
  complete.

### Phase 3 — AST Parser

- Added `internal/parser` AST (`Node`, `NodeType`) and a
  precedence-climbing parser for literals, identifiers, arithmetic
  operators, and parenthesized expressions, satisfying FR-3 for the
  non-function-call case (closes #5).
- Added function-call parsing (`NodeFunctionCall`), including
  arbitrarily nested calls and 0/1/many argument lists, completing
  FR-3 (closes #6).
- Added regression tests locking in the parser's syntax-error paths
  (unbalanced parens, trailing operators, trailing comma in a call,
  empty input) with correct `ErrSyntax` code and position — covers
  the PRD's "invalid formula" use case (closes #7). Phase 3 complete.

### Phase 4 — Function Registry & Math Functions

- Added `internal/registry` with the `Function` interface and a
  self-registering `Register`/`Lookup` registry per `rfc.md` §6
  (closes #8).
- Added `internal/registry/math`: `MAX`, `MIN`, `SUM`, `AVG`, `ROUND`,
  `FLOOR`, `CEIL`, `ABS`, each self-registering via `init()`, with
  arg-count and numeric-type validation (closes #9). Phase 4 complete.

### Phase 5 — Evaluator Core

- Added `internal/evaluator` with `Evaluate`, walking literal and
  arithmetic (`+ - * /`) AST nodes, with division-by-zero surfaced as
  a formula-level `ErrRuntime` per `rfc.md` §7 (closes #10).
- Wired function-call nodes into the evaluator via the registry
  (`registry.Lookup` + `Function.Evaluate`), supporting arbitrary
  nesting and clear errors for unknown functions (closes #11). Phase
  5 complete.

### Phase 6 — Field References & Evaluation Context

- Added `internal/evaluator.Context` and wired identifier nodes
  through it; `Evaluate` now takes a `*Context` parameter. A field
  absent from the data map surfaces a clear `ErrUndefinedReference`
  per FR-6, distinct from a field present with a `nil` value (closes
  #12).
- Added `internal/evaluator/type_coercion.go` implementing every row
  of `rfc.md` §10: string+string concatenation via `+`, number+string
  numeric parsing, bool-as-1/0 in arithmetic, null-as-0 in `+`/`-`
  (and the resulting divide-by-zero error for `/` by null), null-as-
  empty-string in text context, and `ErrTypeMismatch` for anything
  else (closes #13). Phase 6 complete.

### Phase 7 — Dependency Graph & Iterative Topological Sort

- Added `internal/dependency` with `ExtractReferences`, an
  explicit-stack (non-recursive) AST walk collecting referenced
  field/formula names, including from nested function-call arguments
  (closes #14).
- Added `BuildGraph`, building a dependency graph across a batch of
  formulas: references to other formulas in the batch become graph
  edges, plain data fields are resolved directly (not edges), and
  anything else is flagged as undefined, feeding FR-6 (closes #15).
- Added `TopologicalSort` (iterative Kahn's algorithm) producing a
  dependency-respecting evaluation order regardless of input order;
  cycles are safely omitted from the result rather than panicking,
  with formal rejection deferred to Phase 8 (closes #16). Phase 7
  complete.

### Phase 8 — Cycle Detection

- Added `HasCycle`, an iterative (three-color, explicit-stack) cycle
  detector, correct on direct and indirect cycles at 1000+ nodes
  without recursion-depth concerns, per PRD use case 5 (closes #17).
- `Evaluate` now runs the real pipeline end-to-end: parse every
  formula, build the dependency graph, reject circular references and
  undefined references as call-level errors with no partial results
  (rfc.md §7), then topologically sort and evaluate in order, feeding
  each formula's result to later formulas as data. Covers PRD use
  cases 1, 3, 5, and 6 through the public API (closes #18). Phase 8
  complete.

### Phase 9 — Logic Functions

- Added `internal/registry/logic`: `AND`, `OR`, `NOT`. Only real
  booleans are accepted -- rfc.md §10 has no defined coercion from
  number/string/null into a boolean, so anything else hits the
  table's existing `ErrTypeMismatch` catch-all rather than inventing
  a new rule (closes #19).
- Added `IF` with true short-circuit evaluation: the evaluator now
  special-cases it at the AST level (only the taken branch is
  evaluated), since `registry.Function`'s `Evaluate(args []Value)`
  signature can't support lazy branches. `IF(cond, 1, 1/0)` returns
  `1` without error when `cond` is `true`. Completes Phase 9 (closes
  #20).

### Phase 10 — Text Functions

- Added `internal/registry/text`: `CONCAT`, `UPPER`, `LOWER`, `TRIM`,
  `LENGTH`. Arguments must be a string or `nil` (treated as `""` per
  `rfc.md` §10); anything else is a formula-level `ErrTypeMismatch`,
  since the table defines no text-context coercion for
  numbers/bools. `LENGTH` counts Unicode runes, not bytes (closes
  #21).

### Phase 11 — Date Functions

- Added `internal/registry/date`: `NOW`, `DATE_ADD`. Dates are
  represented as `time.Time` values. `NOW` reads through an
  injectable `Now` var (a seam for deterministic tests, per NFR-2)
  instead of calling `time.Now` directly. `DATE_ADD(date, amount,
  unit)` supports `"days"`/`"months"`/`"years"` via `time.AddDate`,
  including its documented month/year-rollover behavior (closes
  #22).
- Added `DATE_DIFF`, returning whole calendar days between two dates.
  Computed from each date's own year/month/day (normalized into UTC)
  rather than raw duration, so it's correct across a leap year (Feb
  28 -> Mar 1 in 2024 is 2 days) and a real DST transition (Mar 8 ->
  Mar 9, 2025 in America/New_York is exactly 1 day despite only 23
  real hours elapsing). Completes Phase 11 (closes #23).

### Phase 12 — Comparison Functions

- Added `internal/registry/comparison`: `EQUALS`, `BETWEEN`. Mixed
  number/numeric-string comparisons resolve via the same coercion
  rules as arithmetic (`rfc.md` §10); a non-numeric string compared
  against a number is a formula-level `ErrTypeMismatch`. `BETWEEN` is
  inclusive of both bounds. Completes Phase 12 -- all planned Phase
  1-12 built-in functions are now implemented (closes #24).

### Phase 13 — Type Coercion Hardening

- Fixed a real gap the audit found: math functions (`MAX MIN SUM AVG
  ROUND FLOOR CEIL ABS`) only coerced `float64`/`int`, silently
  disagreeing with arithmetic operators on `bool`, `nil`, and numeric
  strings (`MAX(true, 0)` errored while `true + 0` didn't). Comparison
  functions were missing the `int` case in the other direction. Fixed
  both to the one canonical coercion rule already correct in
  `internal/evaluator`, which now also accepts `int` for consistency.
  README's math section already (accidentally) claimed this coercion
  worked -- it does now (closes #25).
- Audited text, logic, and date functions the same way. Found and
  fixed the same gap in `DATE_ADD`'s `amount` argument (only
  `float64` was accepted, not even `int`); fixed to the same
  canonical rule. Text (`CONCAT UPPER LOWER TRIM LENGTH`) and logic
  (`AND OR NOT`) were already correct -- confirmed with audit tests
  rather than changed. Completes Phase 13 (closes #26).

### Phase 14 — Error Model Refinement

- Confirmed (no production change needed): a formula-level runtime
  error doesn't block an independent sibling in the same `Evaluate`
  call, per `rfc.md` §7 -- already correct by construction, since
  each formula's result is independent and only a successful result
  feeds later formulas. Locked in with a test and a README example
  (closes #27).
- Fixed a real bug the acceptance test caught: a formula depending on
  one that failed at runtime was misreported as
  `ErrUndefinedReference` (the failed dependency's value never landed
  in the context map, so it looked exactly like a missing field).
  `Evaluate` now checks each formula's dependencies for a prior
  failure before evaluating it and reports a distinct "depends on
  formula X, which failed: ..." error instead -- which cascades
  transitively through any chain of dependents for free, since a
  dependency's own inherited failure is itself a failure. Completes
  Phase 14 (closes #28).

### Phase 15 — Concurrency Safety & Benchmarks

- Added a concurrency test: 200 goroutines call `Evaluate`
  simultaneously with distinct inputs, each asserting its own correct
  result. Clean under `make test-race` -- no race found, no
  production change needed, since the registry's function map is
  populated once at `init()` and never mutated afterward (closes
  #29).
- Added a benchmark suite at 10/100/1,000/10,000 chained formulas in
  one `Evaluate` call, feeding real data into PRD §9's open question
  on call-size performance (no budget is set here, per that section).
  Measured on the CI container (arm64, `golang:1.25-bookworm`):

  | Chain length | Time/op | Memory/op | Allocs/op |
  |---|---|---|---|
  | 10 | 10.9 µs | 13.8 KB | 180 |
  | 100 | 235.8 µs | 414.6 KB | 1,713 |
  | 1,000 | 19.4 ms | 33.8 MB | 17,242 |
  | 10,000 | 2.47 s | 3.5 GB | 322,203 |

  Scaling is clearly superlinear (1,000 -> 10,000 is a 10x chain but
  ~127x the time and ~103x the memory), not the roughly-linear result
  an O(n) pipeline would produce. The likely cause: `Evaluate`'s
  per-formula loop rebuilds a fresh context map from `data` +
  `computed` on every iteration, and `computed` grows by one entry
  each time -- that's O(n²) total map-copy work across a chain of
  length n. Recorded as a finding for a future optimization ticket,
  not fixed here, per this ticket's explicit scope (closes #30).
  **This completes Phase 15 and the core engine (rfc.md §15,
  Phases 1-15).**
