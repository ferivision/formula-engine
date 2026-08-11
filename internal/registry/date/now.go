package date

import (
	"time"

	"github.com/ferivision/formula-engine/internal/registry"
)

// Now is the seam NOW() calls through, so tests (and, if ever
// needed, callers embedding this library) can make its result
// deterministic instead of depending on wall-clock time -- required
// for NFR-2 (identical output for identical input).
var Now = time.Now

func init() {
	registry.Register(nowFunction{})
}

type nowFunction struct{}

func (nowFunction) Name() string { return "NOW" }
func (nowFunction) MinArgs() int { return 0 }
func (nowFunction) MaxArgs() int { return 0 }

func (nowFunction) ValidateArgTypes(args []registry.Value) error {
	return checkArgCount("NOW", 0, 0, args)
}

func (f nowFunction) Evaluate(args []registry.Value) (registry.Value, error) {
	if err := f.ValidateArgTypes(args); err != nil {
		return nil, err
	}
	return Now(), nil
}
