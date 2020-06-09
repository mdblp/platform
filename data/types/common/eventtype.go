package common

import (
	"github.com/tidepool-org/platform/data"
	"github.com/tidepool-org/platform/structure"
)

const (
	StartEvent = "start"
	StopEvent  = "stop"
)

func Events() []string {
	return []string{
		StartEvent,
		StopEvent,
	}
}

type EventType struct {
	EventType *string `json:"eventType,omitempty" bson:"eventType,omitempty"`
}

func NewEventType() *EventType {
	return &EventType{}
}

func (e *EventType) Parse(parser structure.ObjectParser) {
	e.EventType = parser.String("eventType")
}

func (e *EventType) Validate(validator structure.Validator) {
	validator.String("eventType", e.EventType).OneOf(Events()...)
}

func (e *EventType) Normalize(normalizer data.Normalizer) {
}
