package test

import (
	"github.com/tidepool-org/platform/data/types/common"
	dataTypesCommonTest "github.com/tidepool-org/platform/data/types/common/test"
	"github.com/tidepool-org/platform/pointer"
	"github.com/tidepool-org/platform/test"
)

func NewDuration() *common.Duration {
	datum := dataTypesCommonTest.NewDuration()
	datum.Units = pointer.FromString(common.DurationUnitsHours)
	datum.Value = pointer.FromFloat64(test.RandomFloat64FromRange(common.DurationValueHoursMinimum, common.DurationValueHoursMaximum))
	return datum
}

func CloneDuration(datum *common.Duration) *common.Duration {
	if datum == nil {
		return nil
	}
	clone := dataTypesCommonTest.NewDuration()
	clone.Units = pointer.CloneString(datum.Units)
	clone.Value = pointer.CloneFloat64(datum.Value)
	return clone
}
