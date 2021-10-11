package mongo

import (
	"context"
	"errors"
	log "github.com/sirupsen/logrus"

	goComMgo "github.com/mdblp/go-common/clients/mongo"
	"github.com/tidepool-org/platform/data/schema"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoStoreBucketClient struct {
	*goComMgo.StoreClient
	log *log.Logger
}

// Create a new store client for a mongo DB
func NewMongoStoreBucketClient(config *goComMgo.Config, logger *log.Logger) (*MongoStoreBucketClient, error) {
	if config == nil {
		return nil, errors.New("bucket store mongo configuration is missing")
	}

	if logger == nil {
		return nil, errors.New("logger is missing for bucket store client")
	}
	
	client := MongoStoreBucketClient{}
	client.log = logger
	store, err := goComMgo.NewStoreClient(config, logger)
	client.StoreClient = store
	return &client, err
}

/* bucket methods */

// Look for a single bucket based on its internal ID or its public code
func (c *MongoStoreBucketClient) Find(ctx context.Context, bucket *schema.CbgBucket) (result *schema.CbgBucket, err error) {
	
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

// Update a bucket record. The bucket is searched by its internal id
func (c *MongoStoreBucketClient)Upsert(ctx context.Context, userId string, sample *schema.Sample) error {
	
	// Extrat ISODate from sample timestamp
	ts := sample.TimeStamp.Format("02-01-2006")
	valTrue := true

	c.log.Info("upsert cbg sample for: " + userId + "_" + ts)
	
	result, err := c.Collection("hotDailyCbg").UpdateOne(
		ctx,
		bson.D{{Key: "_id", Value: userId + "_" + ts }}, // filter
		bson.D{ // update
			{Key: "$addToSet", Value: bson.D{{Key: "measurements", Value: sample }}},
			{Key: "$setOnInsert", Value: bson.D{{Key: "_id", Value: userId + "_" + ts }}},
			{Key: "$setOnInsert", Value: bson.D{{Key: "createdTimestamp", Value: ts }}},
			{Key: "$setOnInsert", Value: bson.D{{Key: "day", Value: ts }}},
		},
		&options.UpdateOptions{Upsert: &valTrue}, //options
	)
	c.log.Info("error: ", err)
	c.log.Debug(result)

	return err
}

// Deletes a bucket record from the DB
func (c *MongoStoreBucketClient) Remove(ctx context.Context, bucket *schema.CbgBucket) error {

	if bucket.Id != "" {
		if _, err := c.Collection("hotDailyCbg").DeleteOne(ctx, bson.M{"_id": bucket.Id}); err != nil {
			return err
		}
	}

	return errors.New("Remove called with an empty bucket.Id")
}
