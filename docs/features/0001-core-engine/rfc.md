# RFC: Formula Engine — Core Evaluation — Technical Design

**Status:** Proposed
**Author:** _TBD_
**Feature folder:** `docs/features/0001-core-engine/`
**Related:** `prd.md` (same folder), `master-prompt.md` (same folder — execution instructions for Claude Code; this document is the technical design, not an instruction set)

---

## 1. Context

The PRD defines Formula Engine as a stateless, importable Go library
that evaluates spreadsheet-style formulas against supplied data and
returns results, with no storage, no network service, and no UI.

This RFC proposes **how** to build that: package layout, public API
shape, core algorithms, and concrete answers to the open questions
left in the PRD. Where the PRD said "TBD," this document commits to a
specific approach — that's the point of an RFC.

---

## 2. Public API Surface

The library exposes a small number of top-level functions. Everything
else (tokenizer, parser, evaluator, dependency resolution) lives under
`internal/` and is not importable by consumers — this satisfies
PRD NFR-6 (small, clean public interface).

```go
package formulaengine

// FormulaInput is one named formula to evaluate.
type FormulaInput struct {
	Name       string // identifies this formula in the output map
	Expression string // e.g. "MAX(subtotal - discount, 0)"
}

// Evaluate computes one or more formulas against the given data and
// returns a result per formula name, or an error if the call as a
// whole cannot proceed (e.g. a circular reference across the inputs).
func Evaluate(formulas []FormulaInput, data map[string]any) (map[string]Result, error)

// Result is the outcome of evaluating a single formula.
type Result struct {
	Value any   // the computed value, nil if Err is set
	Err   error // per-formula error, nil on success
}
```

Design intent: `Evaluate` takes the whole batch at once (this is how
chaining without persistence works — see §4). It returns a map keyed
by formula name so callers can match results back to their inputs
without relying on order.

---

## 3. Package Structure

```
formula-engine/
├── go.mod
├── go.sum
├── formula.go                    # public API: Evaluate()
├── formula_test.go                # top-level integration-style tests
│
├── internal/
│   ├── parser/
│   │   ├── tokenizer.go           # formula string -> token slice
│   │   ├── token.go               # Token, TokenType
│   │   ├── ast_parser.go          # tokens -> AST (precedence climbing)
│   │   ├── ast.go                 # Node, NodeType
│   │   └── syntax_error.go        # error with line/column position
│   │
│   ├── registry/
│   │   ├── registry.go            # map[string]Function of built-ins
│   │   ├── function.go            # Function interface
│   │   ├── math/                  # max.go, min.go, sum.go, avg.go, round.go, floor.go, ceil.go, abs.go
│   │   ├── logic/                 # if.go, and.go, or.go, not.go
│   │   ├── text/                  # concat.go, upper.go, lower.go, trim.go, length.go
│   │   ├── date/                  # now.go, date_diff.go, date_add.go
│   │   └── comparison/            # equals.go, between.go
│   │
│   ├── evaluator/
│   │   ├── evaluator.go           # walks the AST, produces a value
│   │   ├── context.go             # field data available during evaluation
│   │   ├── type_coercion.go       # conversion rules between operand types
│   │   └── evaluation_error.go
│   │
│   ├── dependency/
│   │   ├── extractor.go           # find referenced fields/formulas from an AST
│   │   ├── graph.go                # build the dependency graph for one call
│   │   ├── topological_sort.go     # iterative Kahn's-algorithm ordering
│   │   └── cycle_detector.go       # detect circular references
│   │
│   └── apperror/
│       └── errors.go               # typed errors (SyntaxError, CircularReferenceError, UndefinedReferenceError, TypeError)
│
└── test/
    └── integration/                 # black-box tests calling only the public API
```

### Conventions

- `snake_case.go` file names, per Go convention.
- No `class` — `struct` + methods, or plain functions for stateless
  logic (tokenizer, evaluator core).
- Interfaces defined on the consumer side (e.g. `evaluator` package
  defines what it needs from the registry, not the other way around).
- Everything under `internal/` is invisible outside this module —
  this is the mechanism that enforces NFR-6, not just a convention.

### Development Environment

- **Go version:** `1.25`. `go.mod` should declare `go 1.25`.
- **Container base image (for CI/build/test environments):**
  `golang:1.25-bookworm`. Don't substitute a different Go version or
  a different base image tag without updating this section first —
  the pinned version is a fixed constraint, not just a suggestion.
- **Local development runs inside this same image**, not whatever Go
  version (if any) is installed on a contributor's machine — this
  keeps everyone, plus CI, on identical behavior. In practice this
  means using the repo's `Makefile` (`make test`, `make build`, etc.),
  which wraps `docker run` against `golang:1.25-bookworm` with the
  project directory mounted in, rather than invoking `go` directly on
  the host.

---

## 4. Chaining Without Persistence

Since the library holds no state (PRD NFR-1), formula chaining works
entirely within a single `Evaluate` call:

1. Parse every `FormulaInput.Expression` into an AST.
2. `extractor.go` walks each AST and records which field names or
   other formula names it references.
3. `graph.go` builds a dependency graph across all formulas passed in
   this call (not across calls).
4. `cycle_detector.go` checks the graph for cycles before anything is
   evaluated. If found, return a `CircularReferenceError` immediately
   — no partial evaluation happens (see §7 for why circularity is
   different from a runtime error).
5. `topological_sort.go` produces a valid evaluation order.
6. `evaluator.go` evaluates formulas in that order, making
   already-computed formula results available as data for formulas
   later in the order.

If the caller wants formula B (from a previous request) available
when evaluating formula A now, they must pass B's definition again in
the current call — the library will not remember it. This is a direct
consequence of the stateless design and is called out explicitly here
so it isn't rediscovered as a surprise during implementation.

---

## 5. Resolving PRD Open Question: Chaining Depth

**Decision:** No enforced maximum depth. `topological_sort.go` and
`cycle_detector.go` are implemented **iteratively** (explicit stack/
queue, e.g. Kahn's algorithm for the sort and iterative DFS with a
visited-set for cycle detection) — not recursively. This means chain
depth is bounded only by available memory, not by Go's call stack, so
correctness holds regardless of depth as required by PRD FR-4.

Performance at large depths is a benchmarking concern for Phase 15,
not a correctness concern to design around now.

---

## 6. Resolving PRD Open Question: Custom Functions

**Decision for this version:** built-in functions only. The
`registry.Function` interface is defined so that adding a custom
function later would be additive (new registry entries), but no
public API for consumer-registered functions ships in this version.
Revisit if/when real usage shows a need.

```go
// internal/registry/function.go
type Function interface {
    Name() string
    MinArgs() int
    MaxArgs() int // -1 = unbounded
    ValidateArgTypes(args []Value) error
    Evaluate(args []Value) (Value, error)
}
```

Each function lives in its own file and registers itself in
`registry.go` — this is the Registry Pattern (+ Strategy Pattern for
the contract), which is what makes "add one more built-in function" a
self-contained, independently testable unit of work (relevant later
for daily tickets).

---

## 7. Resolving PRD Open Question: Partial vs. Fail-All on Runtime Error

**Decision:** **partial success.** If formula A fails at runtime
(e.g. division by zero) but formula B in the same call does not
depend on A, B still returns its computed value. A's entry in the
result map carries its own `Err`; B's carries `nil`.

Circular references and undefined references are treated differently
— those are **call-level** errors (§4 step 4), because they mean the
dependency graph itself is invalid, so nothing in that call can be
safely evaluated. Runtime errors (bad math, type mismatch a coercion
rule can't resolve) are **formula-level** errors and don't block
unrelated formulas.

If a formula depends on another formula that failed, it also fails,
with an error indicating the failure is inherited from its dependency
— not re-attempting evaluation with a missing value.

---

## 8. Resolving PRD Open Question: Typical Call Size / Performance Budget

**Decision:** no hard target yet — insufficient real usage data (PRD
§9 acknowledges this). Phase 15 adds a benchmark suite (10, 100, 1000,
10000 chained formulas) so a concrete number can be set from measured
data rather than guessed. This RFC commits to *measuring*, not to a
number.

---

## 9. Resolving PRD Open Question: Single Package vs. Future Service

**Decision:** ship as a Go package only. If cross-language use becomes
a real requirement later, that's a new RFC (likely: wrap this package
in a thin HTTP/gRPC service) — not a redesign of the core logic, since
parser/evaluator/dependency logic is decoupled from any transport
concern already (see §3, nothing in `internal/` knows about HTTP).

---

## 10. Type Coercion Rules

| Operand types | Behavior |
|---|---|
| number + number | Standard arithmetic |
| string + string (via `CONCAT` or `+`) | Concatenation |
| number + string | Attempt numeric parse of string; error if not parseable |
| bool in numeric context | `true` → 1, `false` → 0 |
| null/missing field in arithmetic | Treated as 0 for `+`/`-`, error for `/` by null |
| null/missing field in text context | Treated as empty string |
| mismatched types with no defined rule | `TypeError`, formula-level (not call-level) |

This table is the authoritative source during implementation — any
function's `ValidateArgTypes` implementation should match it, and
disagreements should update this table first, not be resolved
ad hoc per function.

---

## 11. Error Model

```go
// internal/apperror/errors.go
type ErrorCode string

const (
    ErrSyntax             ErrorCode = "syntax_error"
    ErrCircularReference  ErrorCode = "circular_dependency_detected"
    ErrUndefinedReference ErrorCode = "undefined_reference"
    ErrTypeMismatch       ErrorCode = "type_mismatch"
    ErrRuntime            ErrorCode = "runtime_error"
)

type FormulaError struct {
    Code    ErrorCode
    Message string
    Line    int // 0 if not applicable
    Column  int // 0 if not applicable
}

func (e *FormulaError) Error() string { /* ... */ }
```

Call-level errors (`ErrCircularReference`, `ErrUndefinedReference`)
are returned as `Evaluate`'s second return value. Formula-level errors
(`ErrTypeMismatch`, `ErrRuntime`, `ErrSyntax` for that one formula) are
attached to that formula's `Result.Err`.

---

## 12. Alternatives Considered

- **Separate `Parse()` / `Evaluate()` calls instead of one `Evaluate`
  call** — rejected for v1: adds API surface and a "did you parse
  first" footgun for consumers, with no benefit while there's no
  caching/reuse of parsed formulas across calls (there's nowhere to
  store them, per the stateless design).
- **Recursive dependency resolution** — rejected; see §5, iterative
  is required for FR-4's depth-independence guarantee.
- **Fail-all on any runtime error** — rejected; see §7, partial
  success gives more useful results per call and matches how
  spreadsheet tools like Excel/Airtable behave (one bad cell doesn't
  blank the whole sheet).
- **Lisp-style S-expression syntax instead of infix** (as seen in some
  existing Go expression evaluators) — rejected; infix
  (`MAX(a, b)`, `a + b`) matches the Excel/Airtable/Notion mental model
  the PRD explicitly targets.

---

## 13. Testing Strategy

- **Unit tests** per package (`parser`, `registry/*`, `evaluator`,
  `dependency`) — these are the bulk of test volume and map directly
  to daily tickets.
- **Integration tests** (`test/integration/`) exercise only the public
  `Evaluate` function, covering the use cases from the PRD (§5):
  simple calculation, conditional logic, chained formulas, invalid
  formula, circular formula, missing reference.
- **Benchmarks** (Phase 15) for call size / chain depth, feeding into
  §8's open performance question.

---

## 14. Next Steps

Key decisions above become **ADRs** (one per significant decision:
iterative dependency resolution, partial-success error model, no
custom functions in v1, package-only distribution) — written as they
happen, not all upfront. Implementation proceeds phase by phase per
§15.

---

## 15. Implementation Phases

Sequential — each phase assumes prior phases are done. Estimates are
working days; use them for pacing, not as a hard deadline.

| # | Phase | Goal | Modules touched | Definition of done | Est. |
|---|---|---|---|---|---|
| 1 | Foundation & Public API Skeleton | Compiling module with the final public API shape | `go.mod`, `formula.go`, `internal/apperror/` | `Evaluate(...)` exists with §2's exact signature, compiles, trivial passing test | 1d |
| 2 | Tokenizer | Formula string → token stream | `internal/parser/tokenizer.go`, `token.go` | Numbers, strings, identifiers, operators, parens, commas all tokenize; malformed input reports `SyntaxError` with line/column | 2d |
| 3 | AST Parser | Tokens → AST with correct precedence | `internal/parser/ast.go`, `ast_parser.go` | Binary expressions, nested function calls, and parentheses all parse into a correct tree | 3d |
| 4 | Function Registry & Math Functions | Registry pattern + math category end-to-end | `internal/registry/`, `registry/math/` | `MAX MIN SUM AVG ROUND FLOOR CEIL ABS` implemented, registered, independently tested | 2d |
| 5 | Evaluator Core (no field refs) | Evaluate literals + math functions | `internal/evaluator/evaluator.go` | `"MAX(1 + 2, 4 * 5)"` evaluates correctly with no external data | 2d |
| 6 | Field References & Evaluation Context | Formulas read from the data map | `internal/evaluator/context.go`, `type_coercion.go` | `"price * quantity"` evaluates given data; missing field → clear error; §10 coercion table implemented for nulls | 2d |
| 7 | Dependency Graph & Iterative Topological Sort | Correct evaluation order for chained formulas | `internal/dependency/` | 3+ chained formulas evaluate correctly regardless of input order; sort is iterative (Kahn's algorithm), not recursive | 3d |
| 8 | Cycle Detection | Reject circular references before evaluating | `internal/dependency/cycle_detector.go`, `formula.go` | Direct and indirect cycles rejected as a call-level `ErrCircularReference`, no partial results | 2d |
| 9 | Logic Functions | `IF AND OR NOT` with short-circuit | `internal/registry/logic/` | `IF(TRUE, 1, 1/0)` proves short-circuit; all four implemented and tested | 2d |
| 10 | Text Functions | `CONCAT UPPER LOWER TRIM LENGTH`, null-safe | `internal/registry/text/` | All five implemented, unicode-correct, null handled per §10 | 1-2d |
| 11 | Date Functions | `NOW DATE_DIFF DATE_ADD`, timezone-safe | `internal/registry/date/` | All three implemented; DST and leap-year edge cases tested | 2d |
| 12 | Comparison Functions | `EQUALS BETWEEN` | `internal/registry/comparison/` | Both implemented with type-aware comparison per §10 | 1d |
| 13 | Type Coercion Hardening | Close gaps across all function categories | `internal/evaluator/type_coercion.go` | Cross-cutting test suite verifies every row of §10 against every relevant category | 2d |
| 14 | Error Model Refinement | Full partial-success model across chains | `internal/evaluator/`, `formula.go` | Independent formula failures don't block unrelated results; dependent failures propagate with a clear "inherited" reason | 2d |
| 15 | Concurrency Safety & Benchmarks | Confirm thread-safety, establish perf baseline | `formula_test.go`, `formula_bench_test.go` | `go test -race` clean under concurrent load; benchmarks at 10/100/1,000/10,000-formula batches documented | 2d |

**Total: ~29 working days.** At a realistic pace of a few sessions a
week (not necessarily every single day), this comfortably spans
several weeks of genuine, incremental commit history.