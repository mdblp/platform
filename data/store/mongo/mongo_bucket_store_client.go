package mongo

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	goComMgo "github.com/mdblp/go-db/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/tidepool-org/platform/data/schema"
)

var ErrIncorrectTimestamp = errors.New("impossible to bulk upsert samples having a incorrect timestamp")
var ErrEmptyOrNilUserId = errors.New("impossible to upsert an array of sample for an empty or nil user id")
var ErrUnableToParseBucketDayTime = errors.New("unable to parse cbg day time")
var ErrInvalidDataType = errors.New("invalid empty data type")

var dailyPrefixCollections = []string{"coldDaily", "hotDaily"}

type MongoBucketStoreClient struct {
	*goComMgo.StoreClient
	log                         *log.Logger
	minimalYearSupportedForData int
}

// Create a new bucket store client for a mongo DB if active is set to true, nil otherwise
func NewMongoBucketStoreClient(config *goComMgo.Config, logger *log.Logger, minimalYearSupportedForData int) (*MongoBucketStoreClient, error) {
	if config == nil {
		return nil, errors.New("bucket store mongo configuration is missing")
	}

	if logger == nil {
		return nil, errors.New("logger is missing for bucket store client")
	}

	client := MongoBucketStoreClient{}
	client.log = logger
	store, err := goComMgo.NewStoreClient(config, logger)
	client.StoreClient = store
	client.minimalYearSupportedForData = minimalYearSupportedForData
	return &client, err
}

/* bucket methods */

// Perform a bulk of operations on bucket records based on the operation argument, update a record if found overwhise created it.
// The bucket is searched by its id.
func (c *MongoBucketStoreClient) UpsertMany(ctx context.Context, userId *string, creationTimestamp time.Time, samples []schema.ISample, dataType string) error {

	if userId == nil {
		return ErrEmptyOrNilUserId
	}

	if creationTimestamp.IsZero() {
		return ErrIncorrectTimestamp
	}

	if len(samples) == 0 {
		c.log.Debugf("no %v sample to write, nothing to add in bucket", dataType)
		return nil
	}

	if dataType == "" {
		return ErrInvalidDataType
	}

	var operations []mongo.WriteModel

	// transform as mongo operations
	// no data validation is done here as it is done in above layer in the Validate function
	for _, sample := range samples {
		ts := sample.GetTimestamp().Format("2006-01-02")
		switch dataType {
		case "Cbg":
			ops, _ := buildCbgUpdateOneModel(sample, userId, ts, creationTimestamp)
			operations = append(operations, ops...)
		case "Basal":
			ops, _ := buildBasalUpdateOneModel(sample, userId, ts, creationTimestamp)
			operations = append(operations, ops...)
		case "Bolus":
			ops, _ := buildBolusUpdateOneModel(sample, userId, ts, creationTimestamp)
			operations = append(operations, ops...)
		case "Alarm":
			ops, _ := buildAlarmUpdateOneModel(sample, userId, ts, creationTimestamp)
			operations = append(operations, ops...)
		case "Mode":
			ops, _ := buildModeUpdateOneModel(sample, userId, ts, creationTimestamp)
			operations = append(operations, ops...)
		case "loopMode":
			ops, _ := buildLoopModeWriteModel(sample, userId)
			operations = append(operations, ops...)
		case "Calibration":
			ops, _ := buildCalibrationUpdateOneModel(sample, userId, ts, creationTimestamp)
			operations = append(operations, ops...)
		case "Flush":
			ops, _ := buildFlushUpdateOneModel(sample, userId, ts, creationTimestamp)
			operations = append(operations, ops...)
		case "Prime":
			ops, _ := buildPrimeUpdateOneModel(sample, userId, ts, creationTimestamp)
			operations = append(operations, ops...)
		case "ReservoirChange":
			ops, _ := buildReservoirChangeUpdateOneModel(sample, userId, ts, creationTimestamp)
			operations = append(operations, ops...)
		case "Wizard":
			ops, _ := buildWizardUpdateOneModel(sample, userId, ts, creationTimestamp)
			operations = append(operations, ops...)
		case "Food":
			ops, _ := buildFoodUpdateOneModel(sample, userId, ts, creationTimestamp)
			operations = append(operations, ops...)
		case "PhysicalActivity":
			err := c.upsertPhysicalActivities(ctx, sample, userId, ts, creationTimestamp, dataType)
			if err != nil {
				return fmt.Errorf("cannot upsert physical activity %v : %w", sample, err)
			}
		case "SecurityBasal":
			ops, _ := buildSecurityBasalUpdateOneModel(sample, userId, ts, creationTimestamp)
			operations = append(operations, ops...)
		}
	}

	// Specify an option to turn the bulk insertion with no order of operation
	bulkOption := options.BulkWriteOptions{}
	bulkOption.SetOrdered(false)

	switch dataType {
	/* Do nothing in case of physical activity since operation on the DB are already done */
	case "PhysicalActivity":
	case "SecurityBasal":
		// SecurityBasal event is recorded without hot/cold collection
		_, err := c.Collection("currentSettings").BulkWrite(ctx, operations, &bulkOption)
		if err != nil {
			return err
		}
	case "loopMode":
		// loop mode event is recorded without hot/cold collection
		_, err := c.Collection("loopMode").BulkWrite(ctx, operations, &bulkOption)
		if err != nil {
			return err
		}
	default:
		// update or insert in Hot Daily and Cold Daily
		for _, collectionPrefix := range dailyPrefixCollections {
			// by default
			collectionName := collectionPrefix + dataType

			if dataType == "Alarm" || dataType == "Mode" ||
				dataType == "Calibration" || dataType == "Flush" ||
				dataType == "Prime" || dataType == "ReservoirChange" {
				collectionName = collectionPrefix + "DeviceEvent"
			}
			// All meals (represented by the Wizards object) and the rescue carbs (represented by the Food object)
			// will go to the same collection
			if dataType == "Wizard" || dataType == "Food" {
				collectionName = collectionPrefix + "Food"
			}

			_, err := c.Collection(collectionName).BulkWrite(ctx, operations, &bulkOption)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func buildSecurityBasalUpdateOneModel(sample schema.ISample, userId *string, date string, creationTimestamp time.Time) ([]mongo.WriteModel, error) {

	strUserId := *userId
	var updates []mongo.WriteModel

	op := mongo.NewUpdateOneModel()
	op.SetFilter(bson.D{{Key: "_id", Value: strUserId}})
	op.SetUpdate(bson.D{ // update
		{Key: "$set", Value: bson.D{
			{Key: "securityBasals", Value: sample}}},
		{Key: "$setOnInsert", Value: bson.D{
			{Key: "creationTimestamp", Value: creationTimestamp}}},
	})
	op.SetUpsert(true)
	updates = append(updates, op)

	return updates, nil
}

func buildCbgUpdateOneModel(sample schema.ISample, userId *string, date string, creationTimestamp time.Time) ([]mongo.WriteModel, error) {
	day, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil, ErrUnableToParseBucketDayTime
	}

	strUserId := *userId
	var updates []mongo.WriteModel

	op := mongo.NewUpdateOneModel()
	op.SetFilter(bson.D{{Key: "_id", Value: strUserId + "_" + date}})
	op.SetUpdate(bson.D{ // update
		{Key: "$addToSet", Value: bson.D{
			{Key: "samples", Value: sample}}},
		{Key: "$setOnInsert", Value: bson.D{
			{Key: "_id", Value: strUserId + "_" + date},
			{Key: "creationTimestamp", Value: creationTimestamp},
			{Key: "day", Value: day},
			{Key: "userId", Value: strUserId}}},
	})
	op.SetUpsert(true)
	updates = append(updates, op)
	return updates, nil
}

func buildBasalUpdateOneModel(sample schema.ISample, userId *string, date string, creationTimestamp time.Time) ([]mongo.WriteModel, error) {
	day, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil, ErrUnableToParseBucketDayTime
	}

	strUserId := *userId
	var updates []mongo.WriteModel

	// Insert the bucket if not exist and then insert the sample in it
	basalFirstOp := mongo.NewUpdateOneModel()
	var array []schema.ISample
	basalFirstOp.SetFilter(bson.D{{Key: "_id", Value: strUserId + "_" + date}})
	basalFirstOp.SetUpdate(bson.D{ // update
		{Key: "$setOnInsert", Value: bson.D{
			{Key: "_id", Value: strUserId + "_" + date},
			{Key: "creationTimestamp", Value: creationTimestamp},
			{Key: "day", Value: day},
			{Key: "userId", Value: strUserId},
			{Key: "samples", Value: append(array, sample)},
		},
		},
	})
	basalFirstOp.SetUpsert(true)

	// Update the basal
	basalSecondOp := mongo.NewUpdateOneModel()
	elemfilter := sample.(schema.BasalSample)
	if elemfilter.Guid != "" {
		// All fields update based on guid
		basalSecondOp.SetFilter(bson.D{
			{Key: "_id", Value: strUserId + "_" + date},
			{Key: "samples", Value: bson.D{
				{Key: "$elemMatch", Value: bson.D{
					{Key: "guid", Value: elemfilter.Guid},
				},
				},
			},
			},
		})
		basalSecondOp.SetUpdate(bson.D{ // update
			{Key: "$set", Value: bson.D{
				{Key: "samples.$.internalId", Value: elemfilter.InternalID},
				{Key: "samples.$.duration", Value: elemfilter.Duration},
				{Key: "samples.$.rate", Value: elemfilter.Rate},
				{Key: "samples.$.deliveryType", Value: elemfilter.DeliveryType},
				{Key: "samples.$.timestamp", Value: elemfilter.Timestamp},
			},
			},
		})
	} else {
		// Duration update based on rate/deliveryType/timestamp (nil guid)
		basalSecondOp.SetFilter(bson.D{
			{Key: "_id", Value: strUserId + "_" + date},
			{Key: "samples", Value: bson.D{
				{Key: "$elemMatch", Value: bson.D{
					{Key: "guid", Value: nil},
					{Key: "rate", Value: elemfilter.Rate},
					{Key: "deliveryType", Value: elemfilter.DeliveryType},
					{Key: "timestamp", Value: elemfilter.Timestamp},
				},
				},
			},
			},
		})
		basalSecondOp.SetUpdate(bson.D{ // update
			{Key: "$set", Value: bson.D{
				{Key: "samples.$.internalId", Value: elemfilter.InternalID},
				{Key: "samples.$.duration", Value: elemfilter.Duration},
			},
			},
		})
	}

	// Otherwise we know that we did not update the basal so we guarantee an insertion
	// in the array
	basalThirdOp := mongo.NewUpdateOneModel()
	basalThirdOp.SetFilter(bson.D{{Key: "_id", Value: strUserId + "_" + date}})
	basalThirdOp.SetUpdate(bson.D{ // update
		{Key: "$addToSet", Value: bson.D{
			{Key: "samples", Value: sample}}},
	})
	updates = append(updates, basalFirstOp, basalSecondOp, basalThirdOp)

	return updates, nil
}

func buildBolusUpdateOneModel(sample schema.ISample, userId *string, date string, creationTimestamp time.Time) ([]mongo.WriteModel, error) {
	day, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil, ErrUnableToParseBucketDayTime
	}

	strUserId := *userId
	var updates []mongo.WriteModel

	// Insert the bucket if not exist and then insert the sample in it
	bolusFirstOp := mongo.NewUpdateOneModel()
	var array []schema.ISample
	bolusFirstOp.SetFilter(bson.D{{Key: "_id", Value: strUserId + "_" + date}})
	bolusFirstOp.SetUpdate(bson.D{ // update
		{Key: "$setOnInsert", Value: bson.D{
			{Key: "_id", Value: strUserId + "_" + date},
			{Key: "creationTimestamp", Value: creationTimestamp},
			{Key: "day", Value: day},
			{Key: "userId", Value: strUserId},
			{Key: "samples", Value: append(array, sample)},
		},
		},
	})
	bolusFirstOp.SetUpsert(true)
	updates = append(updates, bolusFirstOp)

	// Update the bolus
	elemfilter := sample.(schema.BolusSample)
	if elemfilter.Guid != "" && elemfilter.DeviceId != "" {
		bolusSecondOp := mongo.NewUpdateOneModel()
		bolusSecondOp.SetFilter(bson.D{
			{Key: "_id", Value: strUserId + "_" + date},
			{Key: "samples", Value: bson.D{
				{Key: "$elemMatch", Value: bson.D{
					{Key: "guid", Value: elemfilter.Guid},
					{Key: "deviceId", Value: elemfilter.DeviceId},
				},
				},
			},
			},
		})
		bolusSecondOp.SetUpdate(bson.D{ // update
			{Key: "$set", Value: bson.D{
				{Key: "samples.$.normal", Value: elemfilter.Normal},
				{Key: "samples.$.uuid", Value: elemfilter.Uuid},
			},
			},
		})
		updates = append(updates, bolusSecondOp)
	}
	// Otherwise we know that we did not update, so we guarantee an insertion
	// in the array
	bolusThirdOp := mongo.NewUpdateOneModel()
	bolusThirdOp.SetFilter(bson.D{{Key: "_id", Value: strUserId + "_" + date}})
	bolusThirdOp.SetUpdate(bson.D{ // update
		{Key: "$addToSet", Value: bson.D{
			{Key: "samples", Value: sample}}},
	})
	updates = append(updates, bolusThirdOp)

	return updates, nil
}

func buildAlarmUpdateOneModel(sample schema.ISample, userId *string, date string, creationTimestamp time.Time) ([]mongo.WriteModel, error) {
	day, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil, ErrUnableToParseBucketDayTime
	}

	strUserId := *userId
	var updates []mongo.WriteModel

	// Insert the bucket if not exist and then insert the sample in it
	firstOp := mongo.NewUpdateOneModel()
	var array []schema.ISample
	firstOp.SetFilter(bson.D{{Key: "_id", Value: strUserId + "_" + date}})
	firstOp.SetUpdate(bson.D{ // update
		{Key: "$setOnInsert", Value: bson.D{
			{Key: "_id", Value: strUserId + "_" + date},
			{Key: "creationTimestamp", Value: creationTimestamp},
			{Key: "day", Value: day},
			{Key: "userId", Value: strUserId},
			{Key: "alarms", Value: append(array, sample)},
		},
		},
	})
	firstOp.SetUpsert(true)
	updates = append(updates, firstOp)

	// Update the bolus
	elemfilter := sample.(schema.AlarmSample)
	if elemfilter.Guid != "" && elemfilter.DeviceId != "" {
		secondOp := mongo.NewUpdateOneModel()
		secondOp.SetFilter(bson.D{
			{Key: "_id", Value: strUserId + "_" + date},
			{Key: "alarms", Value: bson.D{
				{Key: "$elemMatch", Value: bson.D{
					{Key: "guid", Value: elemfilter.Guid},
					{Key: "deviceId", Value: elemfilter.DeviceId},
				},
				},
			},
			},
		})
		secondOp.SetUpdate(bson.D{ // update
			{Key: "$set", Value: bson.D{
				{Key: "alarms.$.level", Value: elemfilter.Level},
				{Key: "alarms.$.ackStatus", Value: elemfilter.AckStatus},
				{Key: "alarms.$.updateTimestamp", Value: elemfilter.UpdateTimestamp},
			},
			},
		})
		updates = append(updates, secondOp)
	}
	// Otherwise we know that we did not update, so we guarantee an insertion
	// in the array
	thirdOp := mongo.NewUpdateOneModel()
	thirdOp.SetFilter(bson.D{{Key: "_id", Value: strUserId + "_" + date}})
	thirdOp.SetUpdate(bson.D{ // update
		{Key: "$addToSet", Value: bson.D{
			{Key: "alarms", Value: sample}}},
	})
	updates = append(updates, thirdOp)

	return updates, nil
}

func buildModeUpdateOneModel(sample schema.ISample, userId *string, date string, creationTimestamp time.Time) ([]mongo.WriteModel, error) {
	day, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil, ErrUnableToParseBucketDayTime
	}

	strUserId := *userId
	var updates []mongo.WriteModel

	// Insert the bucket if not exist and then insert the sample in it
	firstOp := mongo.NewUpdateOneModel()
	var array []schema.ISample
	firstOp.SetFilter(bson.D{{Key: "_id", Value: strUserId + "_" + date}})
	firstOp.SetUpdate(bson.D{ // update
		{Key: "$setOnInsert", Value: bson.D{
			{Key: "_id", Value: strUserId + "_" + date},
			{Key: "creationTimestamp", Value: creationTimestamp},
			{Key: "day", Value: day},
			{Key: "userId", Value: strUserId},
			{Key: "modes", Value: append(array, sample)},
		},
		},
	})
	firstOp.SetUpsert(true)
	updates = append(updates, firstOp)

	// Update the bolus
	elemfilter := sample.(schema.Mode)
	if elemfilter.Guid != "" && elemfilter.DeviceId != "" {
		secondOp := mongo.NewUpdateOneModel()
		secondOp.SetFilter(bson.D{
			{Key: "_id", Value: strUserId + "_" + date},
			{Key: "modes", Value: bson.D{
				{Key: "$elemMatch", Value: bson.D{
					{Key: "guid", Value: elemfilter.Guid},
					{Key: "deviceId", Value: elemfilter.DeviceId},
				},
				},
			},
			},
		})
		secondOp.SetUpdate(bson.D{ // update
			{Key: "$set", Value: bson.D{
				{Key: "modes.$.duration", Value: elemfilter.Duration},
				{Key: "modes.$.inputTimestamp", Value: elemfilter.InputTimestamp},
			},
			},
		})
		updates = append(updates, secondOp)
	}
	// Otherwise we know that we did not update, so we guarantee an insertion
	// in the array
	thirdOp := mongo.NewUpdateOneModel()
	thirdOp.SetFilter(bson.D{{Key: "_id", Value: strUserId + "_" + date}})
	thirdOp.SetUpdate(bson.D{ // update
		{Key: "$addToSet", Value: bson.D{
			{Key: "modes", Value: sample}}},
	})
	updates = append(updates, thirdOp)

	return updates, nil
}

func buildLoopModeWriteModel(sample schema.ISample, userId *string) ([]mongo.WriteModel, error) {

	strUserId := *userId

	elem := sample.(schema.Mode)
	//hack: a mapping here is required to add the user id in to the saved document
	// as the mode struct is a shared model with the bucket elements
	doc := bson.M{}
	doc["timestamp"] = elem.Timestamp
	doc["timezone"] = elem.Timezone
	doc["timezoneOffset"] = elem.TimezoneOffset
	doc["subType"] = elem.SubType
	doc["deviceId"] = elem.DeviceId
	doc["guid"] = elem.Guid
	doc["duration"] = elem.Duration
	doc["inputTimestamp"] = elem.InputTimestamp
	doc["userId"] = strUserId

	var updates []mongo.WriteModel
	var writeOp mongo.WriteModel
	if elem.Guid != "" && elem.DeviceId != "" {
		writeOp = mongo.NewReplaceOneModel().SetFilter(bson.M{"guid": elem.Guid, "userId": strUserId, "deviceId": elem.DeviceId}).SetReplacement(doc).SetUpsert(true)
	} else {
		writeOp = mongo.NewInsertOneModel().SetDocument(doc)
	}
	updates = append(updates, writeOp)
	return updates, nil
}

func buildCalibrationUpdateOneModel(sample schema.ISample, userId *string, date string, creationTimestamp time.Time) ([]mongo.WriteModel, error) {
	day, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil, ErrUnableToParseBucketDayTime
	}

	strUserId := *userId
	var updates []mongo.WriteModel

	op := mongo.NewUpdateOneModel()
	op.SetFilter(bson.D{{Key: "_id", Value: strUserId + "_" + date}})
	op.SetUpdate(bson.D{ // update
		{Key: "$addToSet", Value: bson.D{
			{Key: "calibrations", Value: sample}}},
		{Key: "$setOnInsert", Value: bson.D{
			{Key: "_id", Value: strUserId + "_" + date},
			{Key: "creationTimestamp", Value: creationTimestamp},
			{Key: "day", Value: day},
			{Key: "userId", Value: strUserId}}},
	})
	op.SetUpsert(true)
	updates = append(updates, op)

	return updates, nil
}

func buildFlushUpdateOneModel(sample schema.ISample, userId *string, date string, creationTimestamp time.Time) ([]mongo.WriteModel, error) {
	day, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil, ErrUnableToParseBucketDayTime
	}

	strUserId := *userId
	var updates []mongo.WriteModel

	op := mongo.NewUpdateOneModel()
	op.SetFilter(bson.D{{Key: "_id", Value: strUserId + "_" + date}})
	op.SetUpdate(bson.D{ // update
		{Key: "$addToSet", Value: bson.D{
			{Key: "flushs", Value: sample}}},
		{Key: "$setOnInsert", Value: bson.D{
			{Key: "_id", Value: strUserId + "_" + date},
			{Key: "creationTimestamp", Value: creationTimestamp},
			{Key: "day", Value: day},
			{Key: "userId", Value: strUserId}}},
	})
	op.SetUpsert(true)
	updates = append(updates, op)

	return updates, nil
}

func buildPrimeUpdateOneModel(sample schema.ISample, userId *string, date string, creationTimestamp time.Time) ([]mongo.WriteModel, error) {
	day, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil, ErrUnableToParseBucketDayTime
	}

	strUserId := *userId
	var updates []mongo.WriteModel

	op := mongo.NewUpdateOneModel()
	op.SetFilter(bson.D{{Key: "_id", Value: strUserId + "_" + date}})
	op.SetUpdate(bson.D{ // update
		{Key: "$addToSet", Value: bson.D{
			{Key: "primes", Value: sample}}},
		{Key: "$setOnInsert", Value: bson.D{
			{Key: "_id", Value: strUserId + "_" + date},
			{Key: "creationTimestamp", Value: creationTimestamp},
			{Key: "day", Value: day},
			{Key: "userId", Value: strUserId}}},
	})
	op.SetUpsert(true)
	updates = append(updates, op)

	return updates, nil
}

func buildReservoirChangeUpdateOneModel(sample schema.ISample, userId *string, date string, creationTimestamp time.Time) ([]mongo.WriteModel, error) {
	day, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil, ErrUnableToParseBucketDayTime
	}

	strUserId := *userId
	var updates []mongo.WriteModel

	// one operation because normally the event is sent once
	op := mongo.NewUpdateOneModel()
	op.SetFilter(bson.D{{Key: "_id", Value: strUserId + "_" + date}})
	op.SetUpdate(bson.D{ // update
		{Key: "$addToSet", Value: bson.D{
			{Key: "reservoirChanges", Value: sample}}},
		{Key: "$setOnInsert", Value: bson.D{
			{Key: "_id", Value: strUserId + "_" + date},
			{Key: "creationTimestamp", Value: creationTimestamp},
			{Key: "day", Value: day},
			{Key: "userId", Value: strUserId}}},
	})
	op.SetUpsert(true)
	updates = append(updates, op)

	return updates, nil
}

func buildWizardUpdateOneModel(sample schema.ISample, userId *string, date string, creationTimestamp time.Time) ([]mongo.WriteModel, error) {
	day, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil, ErrUnableToParseBucketDayTime
	}

	strUserId := *userId
	var updates []mongo.WriteModel

	// Insert the bucket if not exist and then insert the sample in it
	firstOp := mongo.NewUpdateOneModel()
	firstOp.SetFilter(bson.D{{Key: "_id", Value: strUserId + "_" + date}})
	firstOp.SetUpdate(bson.D{ // update
		{Key: "$setOnInsert", Value: bson.D{
			{Key: "_id", Value: strUserId + "_" + date},
			{Key: "creationTimestamp", Value: creationTimestamp},
			{Key: "day", Value: day},
			{Key: "userId", Value: strUserId},
			{Key: "meals", Value: []schema.ISample{sample}},
		},
		},
	})
	firstOp.SetUpsert(true)
	updates = append(updates, firstOp)

	// Update
	elemfilter := sample.(schema.Wizard)
	if elemfilter.Guid != "" && elemfilter.DeviceId != "" {
		secondOp := mongo.NewUpdateOneModel()
		secondOp.SetFilter(bson.D{
			{Key: "_id", Value: strUserId + "_" + date},
			{Key: "meals", Value: bson.D{
				{Key: "$elemMatch", Value: bson.D{
					{Key: "guid", Value: elemfilter.Guid},
					{Key: "deviceId", Value: elemfilter.DeviceId},
				},
				},
			},
			},
		})
		secondOp.SetUpdate(bson.D{ // update
			{Key: "$set", Value: bson.D{
				{Key: "meals.$.carbInput", Value: elemfilter.CarbInput},
				{Key: "meals.$.bolus", Value: elemfilter.BolusId},
				{Key: "meals.$.inputTimestamp", Value: elemfilter.InputTimestamp},
				{Key: "meals.$.inputMeal", Value: elemfilter.InputMeal},
			},
			},
		})

		updates = append(updates, secondOp)
	}
	// Otherwise we know that we did not update, so we guarantee an insertion
	// in the array
	thirdOp := mongo.NewUpdateOneModel()
	thirdOp.SetFilter(bson.D{{Key: "_id", Value: strUserId + "_" + date}})
	thirdOp.SetUpdate(bson.D{ // update
		{Key: "$addToSet", Value: bson.D{
			{Key: "meals", Value: sample}}},
	})
	updates = append(updates, thirdOp)

	return updates, nil
}

func buildFoodUpdateOneModel(sample schema.ISample, userId *string, date string, creationTimestamp time.Time) ([]mongo.WriteModel, error) {
	day, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil, ErrUnableToParseBucketDayTime
	}

	strUserId := *userId
	var updates []mongo.WriteModel

	// one operation because normally the event is sent once
	op := mongo.NewUpdateOneModel()
	op.SetFilter(bson.D{{Key: "_id", Value: strUserId + "_" + date}})
	op.SetUpdate(bson.D{ // update
		{Key: "$addToSet", Value: bson.D{
			{Key: "rescueCarbs", Value: sample}}},
		{Key: "$setOnInsert", Value: bson.D{
			{Key: "_id", Value: strUserId + "_" + date},
			{Key: "creationTimestamp", Value: creationTimestamp},
			{Key: "day", Value: day},
			{Key: "userId", Value: strUserId}}},
	})
	op.SetUpsert(true)
	updates = append(updates, op)

	return updates, nil
}

type PhysicalActivityBucket struct {
	Id                string                    `bson:"_id,omitempty"`
	CreationTimestamp time.Time                 `bson:"creationTimestamp,omitempty"`
	UserId            string                    `bson:"userId,omitempty"`
	Day               time.Time                 `bson:"day,omitempty"` // ie: 2021-09-28
	Samples           []schema.PhysicalActivity `bson:"samples"`
}

func findPhysicalActivityBucket(ctx context.Context, collection *mongo.Collection, _id string) (*PhysicalActivityBucket, error) {
	filter := bson.D{
		{Key: "_id", Value: _id},
	}
	var result *PhysicalActivityBucket
	err := collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return result, fmt.Errorf("error while fetching physicalActivity bucket in hot collection"+
			" with _id=[%s]: %s", _id, err.Error())
	}
	return result, err
}

func (c *MongoBucketStoreClient) runInsertOperation(ctx context.Context, document bson.D, dataType string) error {
	for _, collectionPrefix := range dailyPrefixCollections {
		collectionName := collectionPrefix + dataType
		_, err := c.Collection(collectionName).InsertOne(ctx, document)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *MongoBucketStoreClient) runUpdateOperation(ctx context.Context, filter bson.D, document bson.D, dataType string) error {
	for _, collectionPrefix := range dailyPrefixCollections {
		collectionName := collectionPrefix + dataType
		_, err := c.Collection(collectionName).UpdateOne(ctx, filter, document)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *MongoBucketStoreClient) upsertPhysicalActivities(ctx context.Context, sample schema.ISample, userId *string, date string, creationTimestamp time.Time, dataType string) error {
	coldCollectionName := ""
	for _, collectionPrefix := range dailyPrefixCollections {
		if strings.Contains(collectionPrefix, "cold") {
			coldCollectionName = collectionPrefix + dataType
		}
	}

	day, err := time.Parse("2006-01-02", date)
	if err != nil {
		return ErrUnableToParseBucketDayTime
	}
	if userId == nil {
		return fmt.Errorf("userId cannot be nil")
	}
	if sample == nil {
		return fmt.Errorf("sample cannot be nil")
	}
	pa, ok := sample.(schema.PhysicalActivity)
	if !ok {
		return fmt.Errorf("invalid sample type, expecting PhysicalActivity")
	}
	strUserId := *userId
	id := strUserId + "_" + date

	paBucket, err := findPhysicalActivityBucket(ctx, c.Collection(coldCollectionName), id)
	if err != nil {
		return err
	}

	/* Create bucket if it does not exist */
	if paBucket == nil {
		return c.createBucketAndInsertActivity(ctx, pa, id, creationTimestamp, day, userId, dataType)
	}

	/*Bucket existing, we add the sample to it or update it if already existing*/
	array := paBucket.Samples
	sampleFound := false
	for _, sample := range array {
		if sample.DeviceId == pa.DeviceId && sample.Guid == pa.Guid {
			sampleFound = true
			break
		}
	}

	if sampleFound {
		fmt.Printf("Updating activity")
		return c.updateActivity(ctx, id, pa, dataType)
	}

	// Otherwise we know that we did not update, so we guarantee an insertion
	// in the array
	return c.insertActivity(ctx, id, pa, dataType)
}

func (c *MongoBucketStoreClient) insertActivity(ctx context.Context, id string, pa schema.PhysicalActivity, dataType string) error {
	filter := bson.D{{Key: "_id", Value: id}}
	pa.CreateTimestamp = pa.UpdateTimestamp
	update := bson.D{
		{Key: "$addToSet", Value: bson.D{
			{Key: "samples", Value: pa},
		}}}
	return c.runUpdateOperation(ctx, filter, update, dataType)
}

func (c *MongoBucketStoreClient) updateActivity(ctx context.Context, id string, pa schema.PhysicalActivity, dataType string) error {
	filter := bson.D{
		{Key: "_id", Value: id},
		{Key: "samples.deviceId", Value: pa.DeviceId},
		{Key: "samples.guid", Value: pa.Guid},
	}
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "samples.$.reportedIntensity", Value: pa.ReportedIntensity},
			{Key: "samples.$.duration", Value: pa.Duration},
			{Key: "samples.$.updateTimestamp", Value: pa.UpdateTimestamp},
			{Key: "samples.$.timestamp", Value: pa.Timestamp},
		}},
	}

	fmt.Printf("update operation with %v and %v", filter, update)
	return c.runUpdateOperation(ctx, filter, update, dataType)
}

func (c *MongoBucketStoreClient) createBucketAndInsertActivity(ctx context.Context, pa schema.PhysicalActivity, id string, creationTimestamp time.Time, day time.Time, userId *string, dataType string) error {
	pa.CreateTimestamp = pa.UpdateTimestamp
	document := bson.D{
		{Key: "_id", Value: id},
		{Key: "creationTimestamp", Value: creationTimestamp},
		{Key: "day", Value: day},
		{Key: "userId", Value: userId},
		{Key: "samples", Value: []schema.PhysicalActivity{pa}},
	}
	return c.runInsertOperation(ctx, document, dataType)
}

// UpsertMetaData update or insert in MetaData
func (c *MongoBucketStoreClient) UpsertMetaData(ctx context.Context, userId *string, incomingUserMetadata *schema.Metadata) error {

	var dbUserMetadata *schema.Metadata
	var performUpdate bool

	opts := options.FindOne()
	if err := c.Collection("metadata").FindOne(ctx, bson.M{"userId": userId}, opts).Decode(&dbUserMetadata); err != nil && err != mongo.ErrNoDocuments {
		c.log.WithError(err)
		return err
	}

	dbUserMetadata, performUpdate = c.RefreshUserMetadata(dbUserMetadata, incomingUserMetadata)
	valTrue := true

	if performUpdate {
		c.log.Debug("perform update on metadata collection in data_read db")
		_, err := c.Collection("metadata").UpdateOne(ctx,
			bson.M{"userId": userId},
			bson.D{
				{Key: "$set", Value: bson.D{
					{Key: "oldestDataTimestamp", Value: dbUserMetadata.OldestDataTimestamp},
					{Key: "newestDataTimestamp", Value: dbUserMetadata.NewestDataTimestamp}}},
				{Key: "$setOnInsert", Value: bson.D{
					{Key: "creationTimestamp", Value: dbUserMetadata.CreationTimestamp},
					{Key: "userId", Value: dbUserMetadata.UserId}}},
			},
			&options.UpdateOptions{Upsert: &valTrue},
		)
		return err
	}

	return nil
}

func (c *MongoBucketStoreClient) BuildUserMetadata(incomingUserMetadata *schema.Metadata, creationTimestamp time.Time, strUserId string, dataTimestamp time.Time) *schema.Metadata {
	if incomingUserMetadata == nil {
		incomingUserMetadata = &schema.Metadata{
			CreationTimestamp:   creationTimestamp,
			UserId:              strUserId,
			OldestDataTimestamp: dataTimestamp,
			NewestDataTimestamp: dataTimestamp,
		}
	} else {
		if incomingUserMetadata.OldestDataTimestamp.After(dataTimestamp) {
			incomingUserMetadata.OldestDataTimestamp = dataTimestamp
		} else if incomingUserMetadata.NewestDataTimestamp.Before(dataTimestamp) {
			incomingUserMetadata.NewestDataTimestamp = dataTimestamp
		}
	}
	return incomingUserMetadata
}

func (c *MongoBucketStoreClient) RefreshUserMetadata(dbUserMetadata *schema.Metadata, incomingUserMetadata *schema.Metadata) (*schema.Metadata, bool) {
	if dbUserMetadata != nil {
		var performUpdate = false
		//Linked to YLP-1981, in some situation the DBLG1 is sending a data with a timestamp in the near 1970's ...
		//We do not want to update our metadata with this value. The CBG will be recorded with the 1970's date to keep
		//a trace of it, but it won't be displayed since the data is erroneous.
		if dbUserMetadata.OldestDataTimestamp.After(incomingUserMetadata.OldestDataTimestamp) && incomingUserMetadata.OldestDataTimestamp.Year() > c.minimalYearSupportedForData {
			c.log.WithField("oldestDataTimestamp", incomingUserMetadata.OldestDataTimestamp).Debug("set perform update to true and update OldestDataTimestamp db value")
			performUpdate = true
			dbUserMetadata.OldestDataTimestamp = incomingUserMetadata.OldestDataTimestamp
		}
		if dbUserMetadata.NewestDataTimestamp.Before(incomingUserMetadata.NewestDataTimestamp) {
			c.log.WithField("newestDataTimestamp", incomingUserMetadata.NewestDataTimestamp).Debug("set perform update to true and update NewestDataTimestamp db value")
			performUpdate = true
			dbUserMetadata.NewestDataTimestamp = incomingUserMetadata.NewestDataTimestamp
		}
		return dbUserMetadata, performUpdate
	} else {
		return incomingUserMetadata, true
	}
}
