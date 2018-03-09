package factory

import (
	"github.com/tidepool-org/platform/data"
	dataTypesActivityPhysical "github.com/tidepool-org/platform/data/types/activity/physical"
	dataTypesBasal "github.com/tidepool-org/platform/data/types/basal"
	dataTypesBasalFactory "github.com/tidepool-org/platform/data/types/basal/factory"
	dataTypesBloodGlucoseContinuous "github.com/tidepool-org/platform/data/types/blood/glucose/continuous"
	dataTypesBloodGlucoseSelfmonitored "github.com/tidepool-org/platform/data/types/blood/glucose/selfmonitored"
	dataTypesBloodKetone "github.com/tidepool-org/platform/data/types/blood/ketone"
	dataTypesBolus "github.com/tidepool-org/platform/data/types/bolus"
	dataTypesBolusFactory "github.com/tidepool-org/platform/data/types/bolus/factory"
	dataTypesCalculator "github.com/tidepool-org/platform/data/types/calculator"
	dataTypesDevice "github.com/tidepool-org/platform/data/types/device"
	dataTypesDeviceFactory "github.com/tidepool-org/platform/data/types/device/factory"
	dataTypesFood "github.com/tidepool-org/platform/data/types/food"
	dataTypesInsulin "github.com/tidepool-org/platform/data/types/insulin"
	dataTypesSettingsCGM "github.com/tidepool-org/platform/data/types/settings/cgm"
	dataTypesSettingsPump "github.com/tidepool-org/platform/data/types/settings/pump"
	dataTypesStateReported "github.com/tidepool-org/platform/data/types/state/reported"
	dataTypesUpload "github.com/tidepool-org/platform/data/types/upload"
	"github.com/tidepool-org/platform/service"
)

var types = []string{
	dataTypesActivityPhysical.Type,
	dataTypesBasal.Type,
	dataTypesBloodGlucoseContinuous.Type,
	dataTypesBloodGlucoseSelfmonitored.Type,
	dataTypesBloodKetone.Type,
	dataTypesBolus.Type,
	dataTypesCalculator.Type,
	dataTypesDevice.Type,
	dataTypesFood.Type,
	dataTypesInsulin.Type,
	dataTypesSettingsCGM.Type,
	dataTypesSettingsPump.Type,
	dataTypesStateReported.Type,
	dataTypesUpload.Type,
}

func NewDatum(parser data.ObjectParser) data.Datum {
	if parser.Object() == nil {
		return nil
	}

	value := parser.ParseString("type")
	if value == nil {
		parser.AppendError("type", service.ErrorValueNotExists())
		return nil
	}

	switch *value {
	case dataTypesActivityPhysical.Type:
		return dataTypesActivityPhysical.Init()
	case dataTypesBasal.Type:
		return dataTypesBasalFactory.NewBasalDatum(parser)
	case dataTypesBloodGlucoseContinuous.Type:
		return dataTypesBloodGlucoseContinuous.Init()
	case dataTypesBloodGlucoseSelfmonitored.Type:
		return dataTypesBloodGlucoseSelfmonitored.Init()
	case dataTypesBloodKetone.Type:
		return dataTypesBloodKetone.Init()
	case dataTypesBolus.Type:
		return dataTypesBolusFactory.NewBolusDatum(parser)
	case dataTypesCalculator.Type:
		return dataTypesCalculator.Init()
	case dataTypesDevice.Type:
		return dataTypesDeviceFactory.NewDeviceDatum(parser)
	case dataTypesFood.Type:
		return dataTypesFood.Init()
	case dataTypesInsulin.Type:
		return dataTypesInsulin.Init()
	case dataTypesSettingsCGM.Type:
		return dataTypesSettingsCGM.Init()
	case dataTypesSettingsPump.Type:
		return dataTypesSettingsPump.Init()
	case dataTypesStateReported.Type:
		return dataTypesStateReported.Init()
	case dataTypesUpload.Type:
		return dataTypesUpload.Init()
	}

	parser.AppendError("type", service.ErrorValueStringNotOneOf(*value, types))
	return nil
}

func ParseDatum(parser data.ObjectParser) *data.Datum {
	datum := NewDatum(parser)
	if datum == nil {
		return nil
	}

	datum.Parse(parser)
	return &datum
}
