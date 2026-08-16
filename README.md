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

### Arrays

Arrays are a new value type on top of the scalars above: an ordered
list of either **records** (`[]map[string]any`, or a `[]any` where
every element is a `map[string]any` — e.g. what `json.Unmarshal`
produces for a JSON array of objects) or plain scalars (`[]any` of
numbers/strings/bools). Pass one in via the data map; there's no array
literal syntax in formulas themselves.

```go
formulaengine.Evaluate(
	[]formulaengine.FormulaInput{{Name: "result", Expression: `SUMIF(orders, "status", "shipped", "qty")`}},
	map[string]any{
		"orders": []map[string]any{
			{"status": "shipped", "qty": 5.0},
			{"status": "pending", "qty": 2.0},
			{"status": "shipped", "qty": 3.0},
		},
	},
)
// result: 8
```

A `nil` element in the array (a blank slot) is skipped by aggregate
functions, not counted or summed. Mixing records and scalars in the
same array is a `TypeError` as soon as the array is looked up, not
deferred until a function tries to use it.

#### SUMIF / COUNTIF

`SUMIF(array, conditionField, conditionValue, sumField)` sums
`sumField` across records where `conditionField` equals
`conditionValue`. `COUNTIF(array, conditionField, conditionValue)`
counts matching records instead of summing. Both require an array of
records (not scalars); an **empty** array returns `0` for both, not
an error.

```go
formulaengine.FormulaInput{Expression: `SUMIF(orders, "status", "shipped", "qty")`}
// result: 8 (5 + 3, the two "shipped" orders' qty)

formulaengine.FormulaInput{Expression: `COUNTIF(orders, "status", "shipped")`}
// result: 2
```

#### AVERAGEIF / MINIF / MAXIF

Same shape as `SUMIF`: `AVERAGEIF(array, conditionField, conditionValue,
avgField)`, `MINIF(..., minField)`, `MAXIF(..., maxField)`. **Unlike**
`SUMIF`/`COUNTIF`, an empty array — or an array with zero matching
records — is a formula-level `TypeError` for all three, not `0`:
averaging, min, or max over nothing is undefined, so there's no
sensible zero-value result to return instead.

```go
formulaengine.FormulaInput{Expression: `AVERAGEIF(orders, "status", "shipped", "qty")`}
// result: 5 -- average of the two "shipped" orders' qty (4 and 6)

formulaengine.FormulaInput{Expression: `MINIF(orders, "status", "shipped", "qty")`}
// result: 4

formulaengine.FormulaInput{Expression: `MAXIF(orders, "status", "shipped", "qty")`}
// result: 6
```

#### FILTER

`FILTER(array, conditionField, conditionValue)` returns a **new**
array containing only the records where `conditionField` equals
`conditionValue` — the input array is never modified. A `nil` element
can't be meaningfully matched against a condition, so it's silently
excluded from the result rather than causing an error.

```go
formulaengine.FormulaInput{Expression: `FILTER(orders, "status", "shipped")`}
// result: an Array of just the "shipped" records
```

#### UNIQUE

`UNIQUE(array)` deduplicates a **scalar** array by direct value
equality, preserving first-seen order — `1` (number) and `"1"`
(string) are treated as different values here, unlike `FILTER`'s
condition matching, which does coerce across types. `UNIQUE(array,
field)` instead deduplicates a **record** array by a field's value,
keeping the first record seen for each distinct value.

```go
formulaengine.FormulaInput{Expression: "UNIQUE(nums)"}
// data: map[string]any{"nums": []any{1.0, 2.0, 1.0, 3.0}}
// result: an Array of [1, 2, 3]

formulaengine.FormulaInput{Expression: `UNIQUE(orders, "sku")`}
// result: an Array with one record per distinct "sku", first-seen order
```

#### SORT

`SORT(array, sortField, direction)` returns a **new**, sorted array —
`direction` must be `"asc"` or `"desc"`, anything else is a formula-
level error. Uses Go's `sort.SliceStable` internally, so records
sharing an equal sort key keep their original relative order rather
than being reshuffled arbitrarily.

```go
formulaengine.FormulaInput{Expression: `SORT(orders, "qty", "asc")`}
// result: orders ordered by qty ascending, ties broken by original order
```

#### FLATTEN

`FLATTEN(arrayOfArrays)` concatenates an array of nested arrays into
one flat array — the one place nested arrays are allowed in this
library, since it's consuming pre-existing nesting (e.g. from a
consumer that passed `[]any{[]any{...}, []any{...}}`) rather than
producing or navigating it elsewhere.

```go
formulaengine.FormulaInput{Expression: "FLATTEN(groups)"}
// data: map[string]any{"groups": []any{[]any{1.0, 2.0}, []any{3.0, 4.0}}}
// result: an Array of [1, 2, 3, 4]
```

A concretely-typed Go `[][]any` isn't recognized here (or anywhere
`NewArray` converts data) — wrap nested slices as `[]any{...}` instead.

#### VLOOKUP

`VLOOKUP(key, table, keyField, returnField)` finds the first record in
`table` where `keyField` equals `key`, and returns that record's
`returnField` value. A key that isn't found is a **formula-level**
error — it doesn't block an unrelated formula in the same `Evaluate`
call.

```go
formulaengine.FormulaInput{Expression: `VLOOKUP("B2", prices, "sku", "price")`}
// result: 20
```

#### MATCH

`MATCH(key, array, field)` returns the 1-based position of the first
record where `field` equals `key` — same not-found treatment as
`VLOOKUP`.

```go
formulaengine.FormulaInput{Expression: `MATCH("B2", prices, "sku")`}
// result: 2
```

### Error handling

| Situation | Where it shows up |
|---|---|
| Invalid syntax (e.g. unbalanced parens) | that formula's `Result.Err` |
| Runtime error (e.g. division by zero) | that formula's `Result.Err` |
| Circular reference among the formulas in a call | `Evaluate`'s second return value — no results at all |
| Reference to a field/formula not in the call | `Evaluate`'s second return value — no results at all |

Circular and undefined references are call-level: since the whole
dependency graph is invalid, nothing in that call is evaluated.
Everything else (bad syntax, a type mismatch, a runtime error) is
scoped to that one formula's `Result.Err` — an unrelated formula in
the same call still gets its result:

```go
formulaengine.Evaluate(
	[]formulaengine.FormulaInput{
		{Name: "a", Expression: "1 / 0"},
		{Name: "b", Expression: "5 + 5"},
	},
	nil,
)
// a: Err is set (division by zero), Value is nil
// b: Err is nil, Value is 10 -- unaffected by a's failure
```

If a formula instead *depends on* one that failed, it inherits that
failure with a distinct message rather than being evaluated with a
missing value or misreported as an undefined reference — `a` is
defined, it just failed at runtime, which is a different situation
from `a` not existing at all. This cascades transitively through any
chain of dependents.

```go
formulaengine.Evaluate(
	[]formulaengine.FormulaInput{
		{Name: "a", Expression: "1 / 0"},
		{Name: "b", Expression: "a + 1"},
	},
	nil,
)
// a.Err: runtime_error: division by zero
// b.Err: runtime_error: depends on formula "a", which failed: runtime_error: division by zero
```

## Performance

A single, non-chained formula evaluates in low single-digit
microseconds. Chained formulas currently scale **superlinearly**, not
linearly, with chain length — going from a 1,000-formula chain to a
10,000-formula chain (10x the formulas) takes about 127x longer, not
~10x. See `CHANGELOG.md`'s Phase 15 entry for the measured numbers and
the likely cause. No performance budget is set (per the PRD's open
question on this) and no fix has been made yet — if you're chaining
more than a few hundred formulas in one call, benchmark your own
workload rather than assuming linear scaling.

## Not yet supported

- **Unary minus.** `ABS(-5)` and `-price` are not valid syntax yet —
  write `0 - 5` / `0 - price` instead. Negative number literals will
  be added in a future ticket.
- **More array functions.** `INDEX` and `FIND` (the remaining lookup
  functions) are still being built out. See
  [feature-0002's rfc.md §14](docs/features/0002-array-aggregate-lookup-function/feature-0002-rfc.md)
  for the roadmap, or the repo's open issues for what's in progress
  right now.
