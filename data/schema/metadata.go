package schema

import "time"

type (
	Metadata struct {
		Id                 string    `bson:"_id,omitempty"`
		CreationTimestamp  time.Time `bson:"creationTimestamp,omitempty"`
		UserId             string    `bson:"userId,omitempty"`
		OldestCbgTimestamp time.Time `bson:"oldestCbgTimestamp,omitempty"`
		NewestCbgTimestamp time.Time `bson:"newestCbgTimestamp,omitempty"`
	}
)
