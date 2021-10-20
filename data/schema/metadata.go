package schema

import "time"

type (
	Metadata struct {
		Id                       string    `bson:"_id,omitempty"`
		CreationTimestamp        time.Time `bson:"creationTimestamp,omitempty"`
		UserId                   string    `bson:"userId,omitempty"`
		OldestCbgSampleTimestamp time.Time `bson:"oldestCbgSampleTimestamp,omitempty"`
		NewestCbgSampleTimestamp time.Time `bson:"newestCbgSampleTimestamp,omitempty"`
	}
)
