// Command testserver is a manual-testing convenience, not part of the
// formulaengine library. It wraps the public Evaluate() function
// behind one HTTP endpoint so formulas can be tried by sending a
// request instead of writing a Go test file each time.
package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"

	formulaengine "github.com/ferivision/formula-engine"
)

type evaluateRequest struct {
	Formulas []formulaengine.FormulaInput `json:"formulas"`
	Data     map[string]any               `json:"data"`
}

type resultJSON struct {
	Value any     `json:"value"`
	Error *string `json:"error"`
}

type errorJSON struct {
	Error string `json:"error"`
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

func handleEvaluate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, errorJSON{Error: "only POST is supported"})
		return
	}

	var req evaluateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, errorJSON{Error: "invalid request body: " + err.Error()})
		return
	}

	results, err := formulaengine.Evaluate(req.Formulas, req.Data)
	if err != nil {
		writeJSON(w, http.StatusUnprocessableEntity, errorJSON{Error: err.Error()})
		return
	}

	body := make(map[string]resultJSON, len(results))
	for name, res := range results {
		rj := resultJSON{Value: res.Value}
		if res.Err != nil {
			msg := res.Err.Error()
			rj.Error = &msg
		}
		body[name] = rj
	}
	writeJSON(w, http.StatusOK, body)
}

func main() {
	addr := flag.String("addr", ":8080", "address to listen on")
	flag.Parse()

	http.HandleFunc("/evaluate", handleEvaluate)
	log.Printf("testserver listening on %s", *addr)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
