package test

import (
	"github.com/tidepool-org/platform/data/types/activity/physical"
	"github.com/tidepool-org/platform/pointer"
	"github.com/tidepool-org/platform/test"
)

func NewDuration() *physical.Duration {
	datum := physical.NewDuration()
	datum.Units = pointer.FromString(physical.DurationUnitsHours)
	datum.Value = pointer.FromFloat64(test.RandomFloat64FromRange(physical.DurationValueHoursMinimum, physical.DurationValueHoursMaximum))
	return datum
}

func CloneDuration(datum *physical.Duration) *physical.Duration {
	if datum == nil {
		return nil
	}
	clone := physical.NewDuration()
	clone.Units = pointer.CloneString(datum.Units)
	clone.Value = pointer.CloneFloat64(datum.Value)
	return clone
}
