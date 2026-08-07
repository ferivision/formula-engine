# Product Requirements Document: Formula Engine

**Status:** Draft
**Owner:** _TBD_

---

## 1. Problem Statement

Teams building data tools (spreadsheet-like apps, database/record
systems, dashboards) frequently need to let users define computed
fields — e.g. `total = price × quantity` or
`status = IF(stock > 0, "in stock", "out of stock")` — similar to
formulas in Excel, Airtable, or Notion.

Today, this kind of logic usually ends up hardcoded per use case,
duplicated across services, and hard to extend. There's no shared,
reusable way to define a formula, hand it data, and get back a
computed result — with correct handling of formulas that depend on
other formulas.

## 2. Product Concept

Formula Engine is a **pure computation library** — code that other
programs import and call directly, not a standalone service with its
own API or UI. The consuming program passes in one or more formula
definitions plus the data they need, and gets the computed result(s)
back as a normal return value.

The library holds no state of its own: every call is self-contained,
and nothing carries over between calls. Whatever the consuming
program wants to persist — formula definitions, records, results — is
entirely its own responsibility; this library only computes.

That design keeps the library:
- **Reusable** — any program can adopt it without adopting a
  particular storage or deployment model.
- **Predictable** — the same input always produces the same output.
- **Cheap to integrate** — no network hop, no separate process to
  run or operate.

### Intended usage

- As an **importable library** inside a Go program: call an
  evaluation function with formulas and data, get results back
  in-process.
- As the **foundation for a future "formula field" feature** in other
  products, where the product itself owns storing formula definitions
  and records, and simply calls this library whenever it needs to
  (re)compute values.

---

## 3. Goals

- Evaluate one or more formula expressions against supplied data and
  return the correct result for each.
- Provide a library of built-in functions across categories: math,
  logic, text, date, and comparison.
- Support formulas that reference other formulas in the same call
  (chaining), and resolve the correct computation order automatically
  — correctness must hold regardless of how deep the chain goes.
- Detect and reject circular references among the formulas in a call,
  with a clear error.
- Detect and reject references to fields/formulas that don't exist in
  the call, with a clear error.
- Report syntax errors clearly, including where in the formula the
  error occurred.
- Apply predictable, documented rules when a function receives
  arguments of mixed or unexpected types.
- Behave correctly when called concurrently from multiple goroutines,
  with no interference between calls.

## 4. Non-Goals

- Persisting formula definitions, input data, or results between
  calls — that stays with the consuming program.
- A standalone network service, HTTP API, or any deployment/runtime
  of its own.
- Authentication, authorization, or multi-tenant access control.
- Any UI, formula builder, or front-end.
- User-defined/custom functions beyond the built-in set (open
  question, see §9).
- An enforced maximum depth on formula chaining (see §9 — this is a
  performance question to revisit once real usage is known, not a
  product limitation to design in up front).

---

## 5. Use Cases

1. **Simple calculation** — call the library with one formula
   (`price * quantity`) and the relevant data; get back the computed
   number.
2. **Conditional logic** — call the library with a formula using `IF`
   to derive a status label from a numeric field.
3. **Chained formulas** — pass several related formulas in one call
   (e.g. `subtotal`, `discount`, `total`, where `total` depends on the
   other two); get all of them computed in the correct order,
   returned together.
4. **Invalid formula** — pass a formula with a typo or unbalanced
   parentheses; get a clear, specific error describing what's wrong
   and roughly where.
5. **Circular formula** — accidentally pass formulas that reference
   each other in a loop; get a clear error back rather than a hang or
   crash, no matter how long the cycle is.
6. **Missing reference** — pass a formula referencing a field that
   isn't included in the data; get a clear error.

---

## 6. Functional Requirements

| ID | Requirement |
|---|---|
| FR-1 | The library shall evaluate one or more formula expressions against the data supplied in the same call. |
| FR-2 | The library shall provide built-in functions across math, logic, text, date, and comparison categories. |
| FR-3 | The library shall support nested function calls with correct operator precedence. |
| FR-4 | The library shall support formulas referencing other formulas passed in the same call (chaining) and compute them in the correct dependency order, with correctness independent of chain depth. |
| FR-5 | The library shall detect circular references among the formulas passed in and return a clear, specific error. |
| FR-6 | The library shall detect references to undefined fields/formulas and return a clear, specific error. |
| FR-7 | The library shall report syntax errors with enough detail to locate the problem in the formula text. |
| FR-8 | The library shall apply documented type-coercion rules when function arguments have mixed types. |
| FR-9 | The library shall return results for all formulas passed in together, preserving which result belongs to which formula. |

## 7. Non-Functional Requirements

| ID | Requirement |
|---|---|
| NFR-1 | The library must be stateless — no data persists between calls, and no call affects the outcome of another. |
| NFR-2 | The library must produce identical output for identical input, every time. |
| NFR-3 | The library must be safe to call concurrently from multiple goroutines within the same consuming program. |
| NFR-4 | The library must execute fast enough that it doesn't become a bottleneck in the consuming program's normal workflow. |
| NFR-5 | Adding a new built-in function must not require changes to unrelated parts of the codebase. |
| NFR-6 | The library must expose a small, clean public interface, so it's easy to integrate and doesn't leak internal implementation details to consumers. |

---

## 8. Success Metrics

- % of evaluation calls that succeed without error, for well-formed
  input.
- Execution time per evaluation call (p50/p95), measured in-process.
- Number of built-in functions available at launch vs. roadmap.
- Test coverage of edge cases per function and per core module
  (parser, evaluator, dependency resolution).

_(To be refined with concrete targets once real usage is known.)_

---

## 9. Open Questions

- **Chaining depth**: should there be any practical/enforced limit on
  how many formulas can chain together in one call, or is this purely
  a performance characteristic to measure and optimize later rather
  than a product-level constraint? Current assumption: no artificial
  limit — the algorithm should work iteratively so depth is a
  performance question, not a correctness one.
- Should the library support user-defined/custom functions in the
  future, or stay limited to the built-in set?
- What's the expected typical call size (number of formulas per call)
  in real usage, for performance planning?
- Should a runtime error in one formula (e.g. division by zero) fail
  the whole call, or just that one formula's result while the others
  still return?
- What's the target execution-time budget per call?
- Is a single Go package enough, or will other languages eventually
  need this logic too (which would push toward a service later, not
  just a library)?

---

## 10. Next Steps

Once this PRD is agreed on, the technical approach (public function
signatures, evaluation algorithm, concurrency model, error format,
etc.) will be proposed in an **RFC**, key technical decisions from
that RFC will be logged as **ADRs**, and implementation work will be
broken down into tickets from there.