package mongo

import (
	"context"
	"reflect"
	"testing"
	"time"

	goComMgo "github.com/mdblp/go-db/mongo"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/tidepool-org/platform/data/schema"
)

func TestMongoBucketStoreClient_BuildUserMetadata(t *testing.T) {
	type args struct {
		incomingUserMetadata *schema.Metadata
		dbUserMetadata       *schema.Metadata
	}
	testTime := time.Now()
	beforeTestTime := testTime.Add(-24 * time.Hour)
	veryOldTime := testTime.Add(-30 * 365 * 24 * time.Hour)
	tests := []struct {
		name string
		args args
		want *schema.Metadata
	}{
		{
			name: "given empty user metadata should create a new one with passed params",
			args: args{
				dbUserMetadata: nil,
				incomingUserMetadata: &schema.Metadata{
					Id:                  "test1234",
					CreationTimestamp:   testTime,
					UserId:              "123456789",
					OldestDataTimestamp: testTime,
					NewestDataTimestamp: testTime,
				},
			},
			want: &schema.Metadata{
				Id:                  "test1234",
				CreationTimestamp:   testTime,
				UserId:              "123456789",
				OldestDataTimestamp: testTime,
				NewestDataTimestamp: testTime,
			},
		},
		{
			name: "given a 70's data timestamp should not update oldest metadata",
			args: args{
				dbUserMetadata: &schema.Metadata{
					Id:                  "metadata1234",
					CreationTimestamp:   testTime,
					UserId:              "123456789",
					OldestDataTimestamp: testTime,
					NewestDataTimestamp: testTime,
				},
				incomingUserMetadata: &schema.Metadata{
					Id:                  "metadata1234",
					CreationTimestamp:   veryOldTime,
					UserId:              "123456789",
					OldestDataTimestamp: veryOldTime,
					NewestDataTimestamp: veryOldTime,
				},
			},
			want: &schema.Metadata{
				Id:                  "metadata1234",
				CreationTimestamp:   testTime,
				UserId:              "123456789",
				OldestDataTimestamp: testTime,
				NewestDataTimestamp: testTime,
			},
		},
		{
			name: "given a normal data timestamp should update oldest metadata",
			args: args{
				dbUserMetadata: &schema.Metadata{
					Id:                  "metadata1234",
					CreationTimestamp:   testTime,
					UserId:              "123456789",
					OldestDataTimestamp: testTime,
					NewestDataTimestamp: testTime,
				},
				incomingUserMetadata: &schema.Metadata{
					Id:                  "metadata1234",
					CreationTimestamp:   beforeTestTime,
					UserId:              "123456789",
					OldestDataTimestamp: beforeTestTime,
					NewestDataTimestamp: beforeTestTime,
				},
			},
			want: &schema.Metadata{
				Id:                  "metadata1234",
				CreationTimestamp:   testTime,
				UserId:              "123456789",
				OldestDataTimestamp: beforeTestTime,
				NewestDataTimestamp: testTime,
			},
		},
		{
			name: "given a normal data timestamp should update newest metadata",
			args: args{
				dbUserMetadata: &schema.Metadata{
					Id:                  "metadata1234",
					CreationTimestamp:   beforeTestTime,
					UserId:              "123456789",
					OldestDataTimestamp: beforeTestTime,
					NewestDataTimestamp: beforeTestTime,
				},
				incomingUserMetadata: &schema.Metadata{
					Id:                  "metadata1234",
					CreationTimestamp:   testTime,
					UserId:              "123456789",
					OldestDataTimestamp: testTime,
					NewestDataTimestamp: testTime,
				},
			},
			want: &schema.Metadata{
				Id:                  "metadata1234",
				CreationTimestamp:   beforeTestTime,
				UserId:              "123456789",
				OldestDataTimestamp: beforeTestTime,
				NewestDataTimestamp: testTime,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := goComMgo.Config{}
			config.FromEnv()
			c, _ := NewMongoBucketStoreClient(&config, &log.Logger{}, 2015)
			if got, _ := c.RefreshUserMetadata(tt.args.dbUserMetadata, tt.args.incomingUserMetadata); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BuildUserMetadata() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMongoBucketStoreClient_upsertPhysicalActivities_Error(t *testing.T) {
	type args struct {
		date              string
		dataType          string
		userId            *string
		sample            schema.ISample
		creationTimestamp time.Time
	}
	userId := "user1"
	tests := []struct {
		name     string
		given    args
		expected assert.ErrorAssertionFunc
	}{
		{
			name: "should return error when date format invalid",
			given: args{
				date:              "01-01-2024",
				dataType:          "PhysicalActivity",
				userId:            &userId,
				sample:            schema.PhysicalActivity{},
				creationTimestamp: time.Now(),
			},
			expected: assert.Error,
		},
		{
			name: "should return error when sample is not correct type",
			given: args{
				date:              "2024-01-01",
				dataType:          "PhysicalActivity",
				userId:            &userId,
				sample:            schema.AlarmSample{},
				creationTimestamp: time.Now(),
			},
			expected: assert.Error,
		},
		{
			name: "should return error when userId is nil",
			given: args{
				date:              "2024-01-01",
				dataType:          "PhysicalActivity",
				userId:            nil,
				sample:            schema.PhysicalActivity{},
				creationTimestamp: time.Now(),
			},
			expected: assert.Error,
		},
		{
			name: "should return error when activity is nil",
			given: args{
				date:              "2024-01-01",
				dataType:          "PhysicalActivity",
				userId:            &userId,
				sample:            nil,
				creationTimestamp: time.Now(),
			},
			expected: assert.Error,
		},
		{
			name: "should return error when dataType is wrong",
			given: args{
				date:              "2024-01-01",
				dataType:          "physical-activity",
				userId:            &userId,
				sample:            schema.PhysicalActivity{},
				creationTimestamp: time.Now(),
			},
			expected: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := goComMgo.Config{}
			config.FromEnv()
			c, _ := NewMongoBucketStoreClient(&config, &log.Logger{}, 2015)
			err := c.upsertPhysicalActivities(
				t.Context(),
				tt.given.sample,
				tt.given.userId,
				tt.given.date,
				tt.given.creationTimestamp,
				tt.given.dataType)
			tt.expected(t, err)
		})
	}
}

func TestMongoBucketStoreClient_upsertPhysicalActivities(t *testing.T) {
	userId := "user2"
	date := "2024-02-01"
	dataType := "PhysicalActivity"
	activity1 := schema.PhysicalActivity{
		Sample: schema.Sample{
			Timestamp:      time.Now().UTC().Truncate(time.Millisecond),
			Timezone:       "UTC",
			TimezoneOffset: 0,
		},
		Uuid:     "",
		Guid:     "ap1",
		DeviceId: "device123456",
		Duration: schema.Duration{
			Units: "hours",
			Value: 1,
		},
		ReportedIntensity: "medium",
		UpdateTimestamp:   "07-02-2026",
		CreateTimestamp:   "07-02-2026",
	}

	activity1Updated := activity1
	activity1Updated.UpdateTimestamp = "08-02-2026"
	activity1Updated.ReportedIntensity = "low"
	activity1Updated.Timestamp = time.Now().UTC().Truncate(time.Millisecond)

	activity2 := schema.PhysicalActivity{
		Sample: schema.Sample{
			Timestamp:      time.Now().Add(-12 * time.Hour).UTC().Truncate(time.Millisecond),
			Timezone:       "UTC",
			TimezoneOffset: 0,
		},
		Uuid:     "",
		Guid:     "ap2",
		DeviceId: "device123456",
		Duration: schema.Duration{
			Units: "hours",
			Value: 1,
		},
		ReportedIntensity: "medium",
		UpdateTimestamp:   "07-01-2026",
		CreateTimestamp:   "07-01-2026",
	}

	tests := []struct {
		name     string
		samples  []schema.PhysicalActivity
		expected func(*testing.T, *mongo.Cursor)
	}{
		{
			name:    "should create bucket and insert first sample",
			samples: []schema.PhysicalActivity{activity1},
			expected: func(t *testing.T, cursor *mongo.Cursor) {
				var buckets []PhysicalActivityBucket
				cursor.All(context.Background(), &buckets)
				assert.Equal(t, 1, len(buckets))
				assert.Equal(t, 1, len(buckets[0].Samples))
				assert.Equal(t, activity1, buckets[0].Samples[0])
			},
		},
		{
			name:    "should update existing sample if same deviceId and guid",
			samples: []schema.PhysicalActivity{activity1Updated},
			expected: func(t *testing.T, cursor *mongo.Cursor) {
				var buckets []PhysicalActivityBucket
				cursor.All(context.Background(), &buckets)

				assert.Equal(t, 1, len(buckets))
				assert.Equal(t, 1, len(buckets[0].Samples))
				assert.Equal(t, activity1Updated, buckets[0].Samples[0])
			},
		},
		{
			name:    "should insert new sample when different guid",
			samples: []schema.PhysicalActivity{activity2},
			expected: func(t *testing.T, cursor *mongo.Cursor) {
				var buckets []PhysicalActivityBucket
				cursor.All(context.Background(), &buckets)

				assert.Equal(t, 1, len(buckets))
				assert.Equal(t, 2, len(buckets[0].Samples))
				assert.ElementsMatch(t,
					[]schema.PhysicalActivity{activity1Updated, activity2},
					buckets[0].Samples,
				)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := goComMgo.Config{}
			config.FromEnv()
			c, _ := NewMongoBucketStoreClient(&config, &log.Logger{}, 2015)
			for _, sample := range tt.samples {
				err := c.upsertPhysicalActivities(
					context.Background(),
					sample,
					&userId,
					date,
					time.Now(),
					dataType,
				)
				assert.NoError(t, err)
			}

			collectionName := "coldDailyPhysicalActivity"
			cursor, err := c.Collection(collectionName).Find(context.Background(), bson.D{})
			assert.NoError(t, err)
			tt.expected(t, cursor)
		})
	}
}
