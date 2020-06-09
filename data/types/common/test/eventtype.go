package test

import (
	"github.com/tidepool-org/platform/data/types/common"
	"github.com/tidepool-org/platform/pointer"
	"github.com/tidepool-org/platform/test"
)

func NewEventType() *common.EventType {
	datum := common.NewEventType()
	datum.EventType = pointer.FromString(test.RandomStringFromArray(common.Events()))
	return datum
}

func CloneEventType(datum *common.EventType) *common.EventType {
	if datum == nil {
		return nil
	}
	clone := common.NewEventType()
	clone.EventType = pointer.CloneString(datum.EventType)
	return clone
}
