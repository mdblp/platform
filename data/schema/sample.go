package schema

import "time"

type Timestamper interface {
	GetTimestamp() time.Time
}

type (
	Sample struct {
		Timestamp      time.Time `bson:"timestamp,omitempty"`
		Timezone       string    `bson:"timezone,omitempty"`
		TimezoneOffset int       `bson:"timezoneOffset,omitempty"`
	}
)