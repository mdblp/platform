package commontypes_test

import (
	"time"

	commontypes "github.com/tidepool-org/platform/data/types/common"
	"github.com/tidepool-org/platform/test"

	"github.com/tidepool-org/platform/pointer"
)

func NewInputTime() *commontypes.InputTime {
	datum := commontypes.NewInputTime()
	timeReference := test.RandomTime()
	datum.InputTime = pointer.FromString(timeReference.Format(time.RFC3339Nano))
	return datum
}

func CloneInputTime(datum *commontypes.InputTime) *commontypes.InputTime {
	if datum == nil {
		return nil
	}
	clone := commontypes.NewInputTime()
	clone.InputTime = pointer.CloneString(datum.InputTime)
	return clone
}
