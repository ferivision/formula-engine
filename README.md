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

Every argument must be a number (a numeric string or a `bool` also
coerces, per `rfc.md` §10) — anything else is a `TypeError`.

#### MAX / MIN

Largest / smallest of one or more arguments.

```go
formulaengine.FormulaInput{Expression: "MAX(1, 5, 3)"}
// data: none needed
// result: 5

formulaengine.FormulaInput{Expression: "MIN(1, 5, 3)"}
// data: none needed
// result: 1
```

#### SUM / AVG

Sum / average of one or more arguments.

```go
formulaengine.FormulaInput{Expression: "SUM(1, 2, 3.5)"}
// data: none needed
// result: 6.5

formulaengine.FormulaInput{Expression: "AVG(2, 4, 6)"}
// data: none needed
// result: 4
```

#### ROUND

Rounds to the given number of decimal places (defaults to `0` if the
second argument is omitted).

```go
formulaengine.FormulaInput{Expression: "ROUND(3.14159, 2)"}
// data: none needed
// result: 3.14
```

#### FLOOR / CEIL

Rounds down / up to the nearest integer.

```go
formulaengine.FormulaInput{Expression: "FLOOR(3.7)"}
// data: none needed
// result: 3

formulaengine.FormulaInput{Expression: "CEIL(3.2)"}
// data: none needed
// result: 4
```

#### ABS

Absolute value. (There's no unary minus yet — see "not yet
supported" — so a negative input has to come from subtraction, not a
`-5` literal.)

```go
formulaengine.FormulaInput{Expression: "ABS(0 - 5)"}
// data: none needed
// result: 5
```

### Built-in functions — logic

There's no boolean literal syntax yet, so every example below passes
the condition in via a data field rather than a literal `true`/
`false` in the expression itself. All arguments must be actual
booleans — there's no numeric-truthiness coercion (`AND(1, 0)` is a
`TypeError`, not `false`), since `rfc.md` §10 doesn't define one.

#### AND / OR

`AND` is `true` only if every argument is `true`; `OR` is `true` if
any argument is. Both take one or more arguments.

```go
formulaengine.FormulaInput{Expression: "AND(a, b)"}
// data: map[string]any{"a": true, "b": true}
// result: true

formulaengine.FormulaInput{Expression: "OR(a, b)"}
// data: map[string]any{"a": false, "b": true}
// result: true
```

#### NOT

Negates a single boolean.

```go
formulaengine.FormulaInput{Expression: "NOT(a)"}
// data: map[string]any{"a": true}
// result: false
```

#### IF

`IF(condition, ifTrue, ifFalse)` is short-circuiting: only the branch
that's actually taken gets evaluated, so a runtime error in the
*other* branch never surfaces.

```go
formulaengine.FormulaInput{Expression: "IF(cond, 1, 1/0)"}
// data: map[string]any{"cond": true}
// result: 1 -- the "1/0" branch never runs
```

### Built-in functions — text

Every argument must be a string, or `nil` (treated as `""` per
`rfc.md` §10) — a number or bool is a `TypeError`, since the table
doesn't define text-context coercion for those.

#### CONCAT

Concatenates one or more strings.

```go
formulaengine.FormulaInput{Expression: `CONCAT("foo", "bar")`}
// data: none needed
// result: "foobar"
```

#### UPPER / LOWER

Converts case.

```go
formulaengine.FormulaInput{Expression: `UPPER("hello")`}
// data: none needed
// result: "HELLO"

formulaengine.FormulaInput{Expression: `LOWER("HELLO")`}
// data: none needed
// result: "hello"
```

#### TRIM

Removes leading/trailing whitespace.

```go
formulaengine.FormulaInput{Expression: `TRIM("  hello  ")`}
// data: none needed
// result: "hello"
```

#### LENGTH

Character count — Unicode runes, not bytes, so multi-byte characters
still count as one each.

```go
formulaengine.FormulaInput{Expression: `LENGTH("hello")`}
// data: none needed
// result: 5
```

### Built-in functions — date

Dates are represented as Go `time.Time` values — there's no date
literal syntax, so pass them in via the data map (as shown below) or
build them with `NOW`.

#### NOW

Current time, with no arguments.

```go
formulaengine.FormulaInput{Expression: "NOW()"}
// data: none needed
// result: the current time, as a time.Time
```

#### DATE_ADD

`DATE_ADD(date, amount, unit)` adds `amount` of `unit` (`"days"`,
`"months"`, or `"years"`; a negative `amount` subtracts) to `date`.
`amount` follows the same numeric coercion as arithmetic (`rfc.md`
§10) — a bool, `nil`, or numeric string all work, e.g.
`DATE_ADD(start, "3", "days")` is the same as passing `3`. Month/year
arithmetic uses Go's `time.Time.AddDate`, which rolls a day that
doesn't exist in the target month into the following month (e.g.
Jan 31 + 1 month lands in early March) rather than clamping to the
month's last day — that's Go's documented behavior, not a bug in this
library.

```go
formulaengine.FormulaInput{Expression: `DATE_ADD(start, 10, "days")`}
// data: map[string]any{"start": time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)}
// result: 2024-01-11

formulaengine.FormulaInput{Expression: `DATE_ADD(start, 2, "months")`}
// data: map[string]any{"start": time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)}
// result: 2024-03-01
```

#### DATE_DIFF

`DATE_DIFF(date1, date2)` returns the whole number of calendar days
from `date1` to `date2` (negative if `date2` is earlier), computed
from each date's own year/month/day rather than raw duration — so
it's correct across both a leap year and a DST transition, where a
"day" can otherwise be 23 or 25 real hours.

```go
formulaengine.FormulaInput{Expression: "DATE_DIFF(d1, d2)"}
// data: map[string]any{
//     "d1": time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
//     "d2": time.Date(2024, 1, 11, 0, 0, 0, 0, time.UTC),
// }
// result: 10
```

### Built-in functions — comparison

#### EQUALS

Reports whether two values are equal. Same-typed strings or bools
compare directly; a number compared against a bool or a *numeric*
string is coerced per `rfc.md` §10 (`true` → `1`, `nil` → `0`, numeric
strings are parsed); a non-numeric string compared against a number
is a `TypeError`, same as arithmetic's "number + string" rule.

```go
formulaengine.FormulaInput{Expression: "EQUALS(5, 5)"}
// data: none needed
// result: true

formulaengine.FormulaInput{Expression: `EQUALS(5, "5")`}
// data: none needed
// result: true -- the numeric string is parsed and compared as 5
```

#### BETWEEN

Reports whether a value falls within a range, **inclusive of both
bounds**.

```go
formulaengine.FormulaInput{Expression: "BETWEEN(10, 1, 10)"}
// data: none needed
// result: true -- 10 is the upper bound itself, and BETWEEN includes it
```

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

- **Unary minus.** `ABS(-5)` and `-price` are not valid syntax yet —
  write `0 - 5` / `0 - price` instead. Negative number literals will
  be added in a future ticket.

All planned Phase 1-12 built-in functions are now implemented. See
[`rfc.md` §15](docs/features/0001-core-engine/rfc.md) for what's
still ahead (type-coercion hardening, full partial-success error
semantics, concurrency/benchmark verification), or the repo's open
issues for what's in progress right now.
