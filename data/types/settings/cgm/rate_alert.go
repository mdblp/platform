package cgm

import (
	"math"

	"github.com/tidepool-org/platform/data"
	"github.com/tidepool-org/platform/structure"
)

const (
	RateAlertUnitsMgdLMinute  = "mg/dL/minute"
	RateAlertUnitsMmolLMinute = "mmol/L/minute"

	FallAlertRateMgdLMinuteMaximum  = 10.0
	FallAlertRateMgdLMinuteMinimum  = 1.0
	FallAlertRateMmolLMinuteMaximum = 0.55507
	FallAlertRateMmolLMinuteMinimum = 0.05551
	RiseAlertRateMgdLMinuteMaximum  = 10.0
	RiseAlertRateMgdLMinuteMinimum  = 1.0
	RiseAlertRateMmolLMinuteMaximum = 0.55507
	RiseAlertRateMmolLMinuteMinimum = 0.05551
)

func RateAlertUnits() []string {
	return []string{
		RateAlertUnitsMgdLMinute,
		RateAlertUnitsMmolLMinute,
	}
}

type RateAlert struct {
	Alert `bson:",inline"`
	Rate  *float64 `json:"rate,omitempty" bson:"rate,omitempty"`
	Units *string  `json:"units,omitempty" bson:"units,omitempty"`
}

func (r *RateAlert) Parse(parser data.ObjectParser) {
	r.Alert.Parse(parser)
	r.Rate = parser.ParseFloat("rate")
	r.Units = parser.ParseString("units")
}

func (r *RateAlert) Validate(validator structure.Validator) {
	r.Alert.Validate(validator)
	if unitsValidator := validator.String("units", r.Units); r.Rate != nil {
		unitsValidator.Exists().OneOf(RateAlertUnits()...)
	} else {
		unitsValidator.NotExists()
	}
}

type FallAlert struct {
	RateAlert `bson:",inline"`
}

func ParseFallAlert(parser data.ObjectParser) *FallAlert {
	if parser.Object() == nil {
		return nil
	}
	datum := NewFallAlert()
	datum.Parse(parser)
	parser.ProcessNotParsed()
	return datum
}

func NewFallAlert() *FallAlert {
	return &FallAlert{}
}

func (f *FallAlert) Validate(validator structure.Validator) {
	f.RateAlert.Validate(validator)
	validator.Float64("rate", f.Rate).InRange(FallAlertRateRangeForUnits(f.Units))
}

func FallAlertRateRangeForUnits(units *string) (float64, float64) {
	if units != nil {
		switch *units {
		case RateAlertUnitsMgdLMinute:
			return FallAlertRateMgdLMinuteMinimum, FallAlertRateMgdLMinuteMaximum
		case RateAlertUnitsMmolLMinute:
			return FallAlertRateMmolLMinuteMinimum, FallAlertRateMmolLMinuteMaximum
		}
	}
	return -math.MaxFloat64, math.MaxFloat64
}

type RiseAlert struct {
	RateAlert `bson:",inline"`
}

func ParseRiseAlert(parser data.ObjectParser) *RiseAlert {
	if parser.Object() == nil {
		return nil
	}
	datum := NewRiseAlert()
	datum.Parse(parser)
	parser.ProcessNotParsed()
	return datum
}

func NewRiseAlert() *RiseAlert {
	return &RiseAlert{}
}

func (r *RiseAlert) Validate(validator structure.Validator) {
	r.RateAlert.Validate(validator)
	validator.Float64("rate", r.Rate).InRange(RiseAlertRateRangeForUnits(r.Units))
}

func RiseAlertRateRangeForUnits(units *string) (float64, float64) {
	if units != nil {
		switch *units {
		case RateAlertUnitsMgdLMinute:
			return RiseAlertRateMgdLMinuteMinimum, RiseAlertRateMgdLMinuteMaximum
		case RateAlertUnitsMmolLMinute:
			return RiseAlertRateMmolLMinuteMinimum, RiseAlertRateMmolLMinuteMaximum
		}
	}
	return -math.MaxFloat64, math.MaxFloat64
}
