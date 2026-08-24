# testserver

A local HTTP wrapper around the library's public `Evaluate()`
function, for trying out formulas by sending a request instead of
writing a Go test file each time. This is a manual-testing
convenience, not part of the `formulaengine` library itself, and it
follows none of the library's stateless/no-network design — see
`CLAUDE.md`'s "What NOT to do" for why that guardrail doesn't apply
here.

## Run it

```
go run ./cmd/testserver
```

Listens on `:8080` by default; override with `-addr`:

```
go run ./cmd/testserver -addr :9090
```

## Use it

`POST /evaluate` with a JSON body of `formulas` and `data`:

```
curl -X POST http://localhost:8080/evaluate \
  -d '{
    "formulas": [
      {"name": "total", "expression": "price * qty"},
      {"name": "bad", "expression": "1 / 0"}
    ],
    "data": {"price": 10, "qty": 2}
  }'
```

```json
{
  "total": {"value": 20, "error": null},
  "bad": {"value": null, "error": "runtime_error: division by zero"}
}
```

A formula-level error (like the division by zero above) shows up in
that formula's own `error` field with an HTTP 200 — the rest of the
batch still succeeds, matching the library's partial-success model.

A call-level error (e.g. a circular reference across the formulas)
means there are no per-formula results at all, so it's reported as a
top-level error instead, with HTTP 422:

```
curl -X POST http://localhost:8080/evaluate \
  -d '{"formulas": [{"name": "a", "expression": "b + 1"}, {"name": "b", "expression": "a + 1"}], "data": {}}'
```

```json
{"error": "circular_dependency_detected: circular reference detected among the given formulas"}
```

A malformed request body returns HTTP 400.
