package cgm

import (
	"math"

	"github.com/tidepool-org/platform/structure"
)

const (
	SnoozeUnitsHours             = "hours"
	SnoozeUnitsMinutes           = "minutes"
	SnoozeUnitsSeconds           = "seconds"
	SnoozeDurationHoursMaximum   = 10.0
	SnoozeDurationHoursMinimum   = 0.0
	SnoozeDurationMinutesMaximum = SnoozeDurationHoursMaximum * 60.0
	SnoozeDurationMinutesMinimum = SnoozeDurationHoursMinimum * 60.0
	SnoozeDurationSecondsMaximum = SnoozeDurationMinutesMaximum * 60.0
	SnoozeDurationSecondsMinimum = SnoozeDurationMinutesMinimum * 60.0
)

func SnoozeUnits() []string {
	return []string{
		SnoozeUnitsHours,
		SnoozeUnitsMinutes,
		SnoozeUnitsSeconds,
	}
}

type Snooze struct {
	Duration *float64 `json:"duration,omitempty" bson:"duration,omitempty"`
	Units    *string  `json:"units,omitempty" bson:"units,omitempty"`
}

func ParseSnooze(parser structure.ObjectParser) *Snooze {
	if !parser.Exists() {
		return nil
	}
	datum := NewSnooze()
	parser.Parse(datum)
	return datum
}

func NewSnooze() *Snooze {
	return &Snooze{}
}

func (s *Snooze) Parse(parser structure.ObjectParser) {
	s.Duration = parser.Float64("duration")
	s.Units = parser.String("units")
}

func (s *Snooze) Validate(validator structure.Validator) {
	validator.Float64("duration", s.Duration).Exists().InRange(SnoozeDurationRangeForUnits(s.Units))
	validator.String("units", s.Units).Exists().OneOf(SnoozeUnits()...)
}

func SnoozeDurationRangeForUnits(units *string) (float64, float64) {
	if units != nil {
		switch *units {
		case SnoozeUnitsHours:
			return SnoozeDurationHoursMinimum, SnoozeDurationHoursMaximum
		case SnoozeUnitsMinutes:
			return SnoozeDurationMinutesMinimum, SnoozeDurationMinutesMaximum
		case SnoozeUnitsSeconds:
			return SnoozeDurationSecondsMinimum, SnoozeDurationSecondsMaximum
		}
	}
	return -math.MaxFloat64, math.MaxFloat64
}
