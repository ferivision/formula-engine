# Formula Engine

A stateless Go library for evaluating spreadsheet-style formulas —
things like `MAX(subtotal - discount, 0)` or `price * quantity` —
against data you supply, and getting the computed results back. It's
a pure computation library, not a service: no network layer, no
storage, no UI. Think "a small expression evaluator you import and
call," similar in spirit to how you'd pull in a utility library, but
for formula evaluation instead of general data helpers.

Every call is self-contained — nothing persists between calls, and
nothing carries over from one `Evaluate` call to the next. Whatever
you want to keep (formula definitions, records, results) is entirely
your own program's responsibility.

For the full product rationale and technical design, see
[`docs/features/0001-core-engine/prd.md`](docs/features/0001-core-engine/prd.md)
and [`rfc.md`](docs/features/0001-core-engine/rfc.md) in the same
folder.

## Installation

```sh
go get github.com/ferivision/formula-engine
```

## Usage

```go
package main

import (
	"fmt"

	formulaengine "github.com/ferivision/formula-engine"
)

func main() {
	results, err := formulaengine.Evaluate(
		[]formulaengine.FormulaInput{
			{Name: "subtotal", Expression: "price * quantity"},
			{Name: "discount", Expression: "subtotal * 0.1"},
			{Name: "total", Expression: "MAX(subtotal - discount, 0)"},
		},
		map[string]any{
			"price":    100.0,
			"quantity": 2.0,
		},
	)
	if err != nil {
		// Call-level error: a circular or undefined reference among
		// the formulas passed in. Nothing was evaluated.
		panic(err)
	}

	for _, name := range []string{"subtotal", "discount", "total"} {
		r := results[name]
		if r.Err != nil {
			// Formula-level error: this formula (and anything that
			// depends on it) failed, but unrelated formulas in the
			// same call still returned a result.
			fmt.Printf("%s: error: %v\n", name, r.Err)
			continue
		}
		fmt.Printf("%s: %v\n", name, r.Value)
	}
	// subtotal: 200
	// discount: 20
	// total: 180
}
```

`total` here references `subtotal` and `discount` — two other
formulas passed in the *same* call — rather than raw data. The engine
figures out the correct evaluation order on its own; formulas don't
need to be listed in dependency order.

## Currently supported

Every row below is something you can put directly in a
`FormulaInput.Expression` string. `Expression` → `Result` shows what
`results["name"].Value` comes back as for that expression (assuming
no error).

### Arithmetic & precedence

Standard `+ - * /` with correct precedence, and parentheses to
override it.

| Expression | Result |
|---|---|
| `2 + 3 * 4` | `14` |
| `(2 + 3) * 4` | `20` |
| `10 / 4` | `2.5` |

### Field references

Reference any key from the `data` map you pass to `Evaluate` by name.

```go
formulaengine.Evaluate(
	[]formulaengine.FormulaInput{{Name: "total", Expression: "price * quantity"}},
	map[string]any{"price": 10.0, "quantity": 3.0},
)
// total: 30
```

Referencing a key that isn't in `data` (and isn't another formula in
the same call) is a clear `ErrUndefinedReference`, not a silent zero.

### Chained formulas

A formula can reference another formula from the *same call* by its
`Name`, instead of a data field. You don't need to list them in
dependency order — the engine figures that out:

```go
formulaengine.Evaluate(
	[]formulaengine.FormulaInput{
		{Name: "total", Expression: "subtotal - discount"},   // listed before its dependencies
		{Name: "subtotal", Expression: "price * quantity"},
		{Name: "discount", Expression: "subtotal * 0.1"},
	},
	map[string]any{"price": 100.0, "quantity": 2.0},
)
// subtotal: 200, discount: 20, total: 180
```

A cycle (e.g. `a` depends on `b` and `b` depends on `a`) is rejected
as a clear `ErrCircularReference` before anything is evaluated.

### Type coercion

Mixed-type operands follow fixed rules (see `rfc.md` §10 for the full
table) instead of behaving inconsistently per function:

| Expression | Data | Result | Why |
|---|---|---|---|
| `1 + "2"` | — | `3` | numeric string is parsed |
| `"foo" + "bar"` | — | `"foobar"` | string + string concatenates |
| `active + 1` | `{"active": true}` | `2` | `true` → `1` |
| `x + 1` | `{"x": nil}` | `1` | missing/null → `0` in arithmetic |

### Built-in functions — math

| Function | Expression | Result |
|---|---|---|
| `MAX` | `MAX(1, 5, 3)` | `5` |
| `MIN` | `MIN(1, 5, 3)` | `1` |
| `SUM` | `SUM(1, 2, 3.5)` | `6.5` |
| `AVG` | `AVG(2, 4, 6)` | `4` |
| `ROUND` | `ROUND(3.14159, 2)` | `3.14` |
| `FLOOR` | `FLOOR(3.7)` | `3` |
| `CEIL` | `CEIL(3.2)` | `4` |
| `ABS` | `ABS(0 - 5)` | `5` |

`MAX`/`MIN`/`SUM`/`AVG` take one or more arguments; `ROUND` takes an
optional second argument for decimal places (defaults to `0`).

### Built-in functions — logic

There's no boolean literal syntax yet, so examples below pass the
condition in via a data field rather than a literal `true`/`false` in
the expression itself.

| Function | Expression | Data | Result |
|---|---|---|---|
| `AND` | `AND(a, b)` | `{"a": true, "b": true}` | `true` |
| `OR` | `OR(a, b)` | `{"a": false, "b": true}` | `true` |
| `NOT` | `NOT(a)` | `{"a": true}` | `false` |

`AND`/`OR` take one or more arguments. All arguments must be actual
booleans — there's no numeric-truthiness coercion (`AND(1, 0)` is a
`TypeError`, not `false`), since `rfc.md` §10 doesn't define one.

`IF(condition, ifTrue, ifFalse)` is short-circuiting: only the branch
that's actually taken gets evaluated. `IF(cond, 1, 1/0)` returns `1`
without error when `cond` is `true` — the `1/0` branch never runs.
The condition follows the same real-boolean-only rule as `AND`/`OR`.

### Built-in functions — text

| Function | Expression | Result |
|---|---|---|
| `CONCAT` | `CONCAT("foo", "bar")` | `"foobar"` |
| `UPPER` | `UPPER("hello")` | `"HELLO"` |
| `LOWER` | `LOWER("HELLO")` | `"hello"` |
| `TRIM` | `TRIM("  hello  ")` | `"hello"` |
| `LENGTH` | `LENGTH("hello")` | `5` |

`CONCAT` takes one or more arguments. All of them (and `UPPER`/
`LOWER`/`TRIM`/`LENGTH`'s single argument) must be a string or `nil`
(treated as `""`, per `rfc.md` §10) — a number or bool is a
`TypeError`, since the table doesn't define text-context coercion for
those. `LENGTH` counts Unicode runes, not bytes.

### Built-in functions — date

Dates are represented as Go `time.Time` values — there's no date
literal syntax, so pass them in via the data map (as shown below) or
build them with `NOW`.

| Function | Expression | Data | Result |
|---|---|---|---|
| `NOW` | `NOW()` | — | current time as `time.Time` |
| `DATE_ADD` | `DATE_ADD(start, 10, "days")` | `{"start": time.Date(2024,1,1,...)}` | `2024-01-11` |
| `DATE_ADD` | `DATE_ADD(start, 2, "months")` | `{"start": time.Date(2024,1,1,...)}` | `2024-03-01` |

`DATE_ADD(date, amount, unit)` supports `"days"`, `"months"`, and
`"years"` (a negative `amount` subtracts). Month/year arithmetic uses
Go's `time.Time.AddDate`, which rolls a day that doesn't exist in the
target month into the following month (e.g. Jan 31 + 1 month lands in
early March) rather than clamping to the month's last day — that's
Go's documented behavior, not a bug in this library.

### Error handling

| Situation | Where it shows up |
|---|---|
| Invalid syntax (e.g. unbalanced parens) | that formula's `Result.Err` |
| Circular reference among the formulas in a call | `Evaluate`'s second return value — no results at all |
| Reference to a field/formula not in the call | `Evaluate`'s second return value — no results at all |

Circular and undefined references are call-level: since the whole
dependency graph is invalid, nothing in that call is evaluated.
Everything else (bad syntax, a type mismatch) is scoped to that one
formula's `Result.Err`.

## Not yet supported

- Comparison functions, and `DATE_DIFF` — the built-in function set
  is actively growing. See
  [`rfc.md` §15](docs/features/0001-core-engine/rfc.md) for the
  implementation roadmap, or the repo's open issues for what's in
  progress right now.
- **Unary minus.** `ABS(-5)` and `-price` are not valid syntax yet —
  write `0 - 5` / `0 - price` instead. Negative number literals will
  be added in a future ticket.
