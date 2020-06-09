package test

import (
	dataTypesActivityPhysiscalTest "github.com/tidepool-org/platform/data/types/activity/physical/test"
	"github.com/tidepool-org/platform/data/types/bolus/biphasic"
	"github.com/tidepool-org/platform/data/types/bolus/normal"
	"github.com/tidepool-org/platform/pointer"
	"github.com/tidepool-org/platform/test"
)

func NewLinkedBolus() *biphasic.LinkedBolus {
	datum := biphasic.NewLinkedBolus()
	datum.Normal = pointer.FromFloat64(test.RandomFloat64FromRange(normal.NormalMinimum, normal.NormalMaximum))
	datum.Duration = dataTypesActivityPhysiscalTest.NewDuration()
	return datum
}

func CloneLinkedBolus(datum *biphasic.LinkedBolus) *biphasic.LinkedBolus {
	if datum == nil {
		return nil
	}
	clone := biphasic.NewLinkedBolus()
	clone.Normal = pointer.CloneFloat64(datum.Normal)
	clone.Duration = dataTypesActivityPhysiscalTest.CloneDuration(datum.Duration)
	return clone
}
