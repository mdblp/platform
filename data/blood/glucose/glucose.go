package glucose

import (
	"math"

	"github.com/tidepool-org/platform/pointer"
)

const (
	MmolL = "mmol/L"
	Mmoll = "mmol/l"

	MgdL = "mg/dL"
	Mgdl = "mg/dl"

	MmolLMinimum float64 = 0.0
	MmolLMaximum float64 = 55.0

	MgdLMinimum float64 = 0.0
	MgdLMaximum float64 = 1000.0

	MmolLToMgdLConversionFactor float64 = 18.01577
	MmolLToMgdLPrecisionFactor  float64 = 10.0
)

func Units() []string {
	return []string{MmolL, Mmoll, MgdL, Mgdl}
}

func ValueRangeForUnits(units *string) (float64, float64) {
	if units != nil {
		switch *units {
		case MmolL, Mmoll:
			return MmolLMinimum, MmolLMaximum
		case MgdL, Mgdl:
			return MgdLMinimum, MgdLMaximum
		}
	}
	return -math.MaxFloat64, math.MaxFloat64
}

func NormalizeUnits(units *string) *string {
	if units != nil {
		switch *units {
		case MmolL, Mmoll, MgdL, Mgdl:
			return pointer.FromString(MmolL)
		}
	}
	return units
}

func NormalizeValueForUnits(value *float64, units *string) *float64 {
	if value != nil && units != nil {
		switch *units {
		case MgdL, Mgdl:
			roundedValue := math.Round(*value / MmolLToMgdLConversionFactor * MmolLToMgdLPrecisionFactor)
			floatValue := roundedValue / MmolLToMgdLPrecisionFactor
			return &floatValue
		}
	}
	return value
}
