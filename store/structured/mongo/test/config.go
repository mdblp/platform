package test

import (
	"time"

	storeStructuredMongo "github.com/tidepool-org/platform/store/structured/mongo"
)

//NewConfig creates a test Mongo configuration
func NewConfig() *storeStructuredMongo.Config {
	conf := &storeStructuredMongo.Config{
		Database:               Database(),
		CollectionPrefix:       NewCollectionPrefix(),
		Timeout:                5 * time.Second,
		WaitConnectionInterval: 1 * time.Second,
		MaxConnectionAttempts:  1,
	}
	conf.Addresses = []string{Address()}

	return conf
}
