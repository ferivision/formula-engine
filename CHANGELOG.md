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
