package dependency

import (
	"sort"
	"testing"

	"github.com/ferivision/formula-engine/internal/parser"
)

func extractFromSource(t *testing.T, input string) []string {
	t.Helper()
	tokens, err := parser.Tokenize(input)
	if err != nil {
		t.Fatalf("Tokenize(%q) error = %v", input, err)
	}
	node, err := parser.Parse(tokens)
	if err != nil {
		t.Fatalf("Parse(%q) error = %v", input, err)
	}
	refs := ExtractReferences(node)
	sort.Strings(refs)
	return refs
}

func assertRefs(t *testing.T, got []string, want []string) {
	t.Helper()
	sort.Strings(want)
	if len(got) != len(want) {
		t.Fatalf("refs = %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("refs = %v, want %v", got, want)
		}
	}
}

func TestExtractReferences_SimpleIdentifier(t *testing.T) {
	assertRefs(t, extractFromSource(t, "price"), []string{"price"})
}

func TestExtractReferences_LiteralsOnlyHaveNoReferences(t *testing.T) {
	assertRefs(t, extractFromSource(t, "1 + 2"), []string{})
}

func TestExtractReferences_BinaryExpression(t *testing.T) {
	assertRefs(t, extractFromSource(t, "price * quantity"), []string{"price", "quantity"})
}

func TestExtractReferences_NestedFunctionCallArguments(t *testing.T) {
	assertRefs(t, extractFromSource(t, "MAX(price, MIN(discount, total))"), []string{"price", "discount", "total"})
}

func TestExtractReferences_FunctionNameItselfNotIncluded(t *testing.T) {
	assertRefs(t, extractFromSource(t, "MAX(price, 1)"), []string{"price"})
}

func TestExtractReferences_DuplicatesAreDeduplicated(t *testing.T) {
	assertRefs(t, extractFromSource(t, "price + price"), []string{"price"})
}
