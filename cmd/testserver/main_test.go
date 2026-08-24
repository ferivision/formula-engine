package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func doEvaluate(t *testing.T, body string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, "/evaluate", strings.NewReader(body))
	rec := httptest.NewRecorder()
	handleEvaluate(rec, req)
	return rec
}

func TestHandleEvaluate_ReturnsPerFormulaResults(t *testing.T) {
	body := `{
		"formulas": [{"name": "total", "expression": "price * qty"}],
		"data": {"price": 10, "qty": 2}
	}`

	rec := doEvaluate(t, body)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200, body = %s", rec.Code, rec.Body.String())
	}
	var got map[string]resultJSON
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("json.Unmarshal() error = %v, body = %s", err, rec.Body.String())
	}
	if got["total"].Error != nil {
		t.Errorf("total.Error = %v, want nil", *got["total"].Error)
	}
	if got["total"].Value != 20.0 {
		t.Errorf("total.Value = %v, want 20", got["total"].Value)
	}
}

// A formula-level error (e.g. division by zero) must be reflected in
// that formula's own error field, not a 500 -- consistent with the
// library's partial-success model. An unrelated sibling formula must
// still succeed.
func TestHandleEvaluate_FormulaLevelErrorDoesNotFailTheRequest(t *testing.T) {
	body := `{
		"formulas": [
			{"name": "bad", "expression": "1 / 0"},
			{"name": "good", "expression": "5 + 5"}
		],
		"data": {}
	}`

	rec := doEvaluate(t, body)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200, body = %s", rec.Code, rec.Body.String())
	}
	var got map[string]resultJSON
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("json.Unmarshal() error = %v, body = %s", err, rec.Body.String())
	}
	if got["bad"].Error == nil {
		t.Error("bad.Error = nil, want a division-by-zero error")
	}
	if got["good"].Error != nil || got["good"].Value != 10.0 {
		t.Errorf("good = %+v, want {Value: 10, Error: nil}", got["good"])
	}
}

// A call-level error (e.g. a circular reference) means Evaluate()
// itself returns an error with no results -- this must surface as a
// clear error response, not a panic or an empty 200.
func TestHandleEvaluate_CallLevelErrorReturnsErrorResponse(t *testing.T) {
	body := `{
		"formulas": [
			{"name": "a", "expression": "b + 1"},
			{"name": "b", "expression": "a + 1"}
		],
		"data": {}
	}`

	rec := doEvaluate(t, body)

	if rec.Code == http.StatusOK {
		t.Fatalf("status = 200, want a non-200 error response, body = %s", rec.Body.String())
	}
	var got errorJSON
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("json.Unmarshal() error = %v, body = %s", err, rec.Body.String())
	}
	if got.Error == "" {
		t.Error("Error = \"\", want a non-empty message")
	}
}

func TestHandleEvaluate_MalformedJSONReturns400(t *testing.T) {
	rec := doEvaluate(t, `{not valid json`)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400, body = %s", rec.Code, rec.Body.String())
	}
}

func TestHandleEvaluate_WrongMethodReturns405(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/evaluate", nil)
	rec := httptest.NewRecorder()
	handleEvaluate(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("status = %d, want 405", rec.Code)
	}
}
