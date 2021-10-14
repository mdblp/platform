package mongo

import (
	"context"
	"errors"
	"time"

	log "github.com/sirupsen/logrus"

	goComMgo "github.com/mdblp/go-common/clients/mongo"
	"github.com/tidepool-org/platform/data/schema"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoBucketStoreClient struct {
	*goComMgo.StoreClient
	log *log.Logger
}

// Create a new bucket store client for a mongo DB
func NewMongoBucketStoreClient(config *goComMgo.Config, logger *log.Logger) (*MongoBucketStoreClient, error) {
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
	return &client, err
}

/* bucket methods */

// Look for a single bucket based on its Id
func (c *MongoBucketStoreClient) Find(ctx context.Context, bucket *schema.CbgBucket) (result *schema.CbgBucket, err error) {
	
	if bucket.Id != "" {
		var query bson.M = bson.M{}
		tid, _ := primitive.ObjectIDFromHex(bucket.Id)
		query["_id"] = tid
		opts := options.FindOne()
		opts.SetSort(bson.D{primitive.E{Key: "_id", Value: -1}})
		if err = c.Collection("hotDailyCbg").FindOne(ctx, query, opts).Decode(&result); err != nil && err != mongo.ErrNoDocuments {
			c.log.WithError(err)
			return result, err
		}
	
		return result, nil
	}

	return nil, errors.New("Find called with an empty bucket.Id")
}

// Update a bucket record if found overwhise it will be created. The bucket is searched by its id.
func (c *MongoBucketStoreClient)Upsert(ctx context.Context, userId *string, creationTimestamp time.Time, sample *schema.CbgSample) error {
	
	if sample == nil {
		return errors.New("impossible to upsert a nil sample")
	}

	if sample.TimeStamp.IsZero() {
		return errors.New("impossible to upsert a sample having a incorrect timestamp")
	}

	if userId == nil {
		return errors.New("impossible to upsert a sample for an empty user id")
	}

	// Extrat ISODate from sample timestamp
	ts := sample.TimeStamp.Format("02-01-2006")
	valTrue := true
	strUserId := *userId

	c.log.Info("upsert cbg sample for: " + strUserId + "_" + ts)
	// save in hotDailyCbg
	_, err := c.Collection("hotDailyCbg").UpdateOne(
		ctx,
		bson.D{{Key: "_id", Value: strUserId + "_" + ts }}, // filter
		bson.D{ // update
			{Key: "$addToSet", Value: bson.D{{Key: "samples", Value: sample }}},
			{Key: "$setOnInsert", Value: bson.D{{Key: "_id", Value: strUserId + "_" + ts }}},
			{Key: "$setOnInsert", Value: bson.D{{Key: "creationTimestamp", Value: creationTimestamp }}},
			{Key: "$setOnInsert", Value: bson.D{{Key: "day", Value: ts }}},
		},
		&options.UpdateOptions{Upsert: &valTrue}, //options
	)

	if err != nil {
		return err
	}

	// save in coldDailyCbg
	_, err = c.Collection("coldDailyCbg").UpdateOne(
		ctx,
		bson.D{{Key: "_id", Value: strUserId + "_" + ts }}, // filter
		bson.D{ // update
			{Key: "$addToSet", Value: bson.D{{Key: "samples", Value: sample }}},
			{Key: "$setOnInsert", Value: bson.D{{Key: "_id", Value: strUserId + "_" + ts }}},
			{Key: "$setOnInsert", Value: bson.D{{Key: "creationTimestamp", Value: creationTimestamp }}},
			{Key: "$setOnInsert", Value: bson.D{{Key: "day", Value: ts }}},
		},
		&options.UpdateOptions{Upsert: &valTrue}, //options
	)

	return err
}

// Perform a bulk of operations on bucket records based on the operation argument, update a record if found overwhise created ot. 
// The bucket is searched by its id.
func (c *MongoBucketStoreClient)UpsertMany(ctx context.Context, userId *string, operations []mongo.WriteModel) error {
	
	// Specify an option to turn the bulk insertion in order of operation
	bulkOption := options.BulkWriteOptions{}
	bulkOption.SetOrdered(true)

	_, err := c.Collection("hotDailyCbg").BulkWrite(ctx, operations, &bulkOption)
	if err != nil {
        return err
	}

	_, err = c.Collection("coldDailyCbg").BulkWrite(ctx, operations, &bulkOption)

	return err
}

// Deletes a bucket record from the DB
func (c *MongoBucketStoreClient) Remove(ctx context.Context, bucket *schema.CbgBucket) error {

	if bucket.Id != "" {
		if _, err := c.Collection("hotDailyCbg").DeleteOne(ctx, bson.M{"_id": bucket.Id}); err != nil {
			return err
		}

		if _, err := c.Collection("coldDailyCbg").DeleteOne(ctx, bson.M{"_id": bucket.Id}); err != nil {
			return err
		}
	}

	return errors.New("Remove called with an empty bucket.Id")
}
