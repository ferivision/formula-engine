package date

import (
	"time"

	"github.com/ferivision/formula-engine/internal/registry"
)

func init() {
	registry.Register(dateDiffFunction{})
}

type dateDiffFunction struct{}

func (dateDiffFunction) Name() string { return "DATE_DIFF" }
func (dateDiffFunction) MinArgs() int { return 2 }
func (dateDiffFunction) MaxArgs() int { return 2 }

func (f dateDiffFunction) ValidateArgTypes(args []registry.Value) error {
	if err := checkArgCount(f.Name(), f.MinArgs(), f.MaxArgs(), args); err != nil {
		return err
	}
	if _, err := toTime(f.Name(), args[0]); err != nil {
		return err
	}
	_, err := toTime(f.Name(), args[1])
	return err
}

// Evaluate returns the whole number of calendar days from date1 to
// date2 (date2 - date1), computed from each date's own year/month/day
// components rather than raw duration. Raw duration/24h would be
// wrong across a DST transition, where a "day" can be 23 or 25 real
// hours; comparing calendar dates normalized into UTC (which has no
// DST) sidesteps that entirely, and correctly accounts for leap years
// since it's built on time.Date's own calendar arithmetic.
func (f dateDiffFunction) Evaluate(args []registry.Value) (registry.Value, error) {
	if err := f.ValidateArgTypes(args); err != nil {
		return nil, err
	}
	t1, _ := toTime(f.Name(), args[0])
	t2, _ := toTime(f.Name(), args[1])
	return float64(daysSinceEpoch(t2) - daysSinceEpoch(t1)), nil
}

func daysSinceEpoch(t time.Time) int64 {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC).Unix() / 86400
}
