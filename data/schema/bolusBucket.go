package schema

import "time"

type (
	BolusBucket struct {
		Id                string        `bson:"_id,omitempty"`
		CreationTimestamp time.Time     `bson:"creationTimestamp,omitempty"`
		UserId            string        `bson:"userId,omitempty" `
		Day               time.Time     `bson:"day,omitempty"` // ie: 2021-09-28
		Samples           []BolusSample `bson:"samples"`
	}

	BolusSample struct {
		Sample         `bson:",inline"`
		Uuid           string   `bson:"uuid,omitempty"`
		DeviceId       string   `bson:"deviceId,omitempty"`
		Guid           string   `bson:"guid,omitempty"`
		BolusType      string   `bson:"bolusType,omitempty" json:"bolusType,omitempty"`
		Normal         float64  `bson:"normal,omitempty" json:"normal,omitempty"`
		ExpectedNormal *float64 `bson:"expectedNormal,omitempty" json:"expectedNormal,omitempty"`
		InsulinOnBoard *float64 `bson:"insulinOnBoard,omitempty" json:"insulinOnBoard,omitempty"`
		Prescriptor    *string  `bson:"prescriptor,omitempty" json:"prescriptor,omitempty"`
		BiphasicId     *string  `bson:"biphasicId,omitempty" json:"biphasicId,omitempty"`
		Part           int64    `bson:"part,omitempty" json:"part,omitempty"`
	}
)

func (b BolusBucket) GetId() string {
	return b.Id
}

func (s BolusSample) GetTimestamp() time.Time {
	return s.Timestamp
}
