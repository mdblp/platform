package mongo

import (
	"context"
	"fmt"
	"time"

	mongoDriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"github.com/tidepool-org/platform/errors"
)

type Store struct {
	client *mongoDriver.Client
	config *Config
}

type Status struct {
	Ping string
}

func NewStore(c *Config) (*Store, error) {
	if c == nil {
		return nil, errors.New("database config is empty")
	}

	store := &Store{
		config: c,
	}

	cs := c.AsConnectionString()
	clientOptions := options.Client().
		ApplyURI(cs).
		SetConnectTimeout(store.config.Timeout).
		SetServerSelectionTimeout(store.config.Timeout)

	var attempts int64 = 1
	var err error
	var closingChannel = make(chan bool, 1)
	for {
		var timer <-chan time.Time
		if attempts == int64(0) {
			timer = time.After(0)
		} else {
			timer = time.After(store.config.WaitConnectionInterval)
		}
		select {
		case <-closingChannel:
			close(closingChannel)
			if err != nil {
				return nil, err
			}
			return store, nil
		case <-timer:
			err = initConnexion(store, clientOptions)
			if err == nil {
				closingChannel <- true
			} else {
				if store.config.MaxConnectionAttempts > 0 && store.config.MaxConnectionAttempts <= attempts {
					closingChannel <- true
				} else {
					attempts++
				}
			}
		}
	}
}

func initConnexion(store *Store, clientOptions *options.ClientOptions) error {
	var err error
	store.client, err = mongoDriver.Connect(context.Background(), clientOptions)
	if err != nil {
		fmt.Println("connection options are invalid")
		return errors.Wrap(err, "connection options are invalid")
	}
	ctx, cancel := context.WithTimeout(context.Background(), store.config.Timeout)
	defer cancel()
	err = store.client.Ping(ctx, readpref.PrimaryPreferred())
	if err != nil {
		fmt.Println("cannot ping store")
		return errors.Wrap(err, "cannot ping store")
	}
	return nil
}

func (o *Store) GetRepository(collection string) *Repository {
	return NewRepository(o.GetCollectionWithArchive(collection))
}

func (o *Store) GetCollectionWithArchive(collection string) (*mongoDriver.Collection, *mongoDriver.Collection) {
	db := o.client.Database(o.config.Database)
	prefixed := fmt.Sprintf("%s%s", o.config.CollectionPrefix, collection)
	prefixedArchive := fmt.Sprintf("%s%s_archive", o.config.CollectionPrefix, collection)
	return db.Collection(prefixed), db.Collection(prefixedArchive)
}

func (o *Store) GetCollection(collection string) *mongoDriver.Collection {
	db := o.client.Database(o.config.Database)
	prefixed := fmt.Sprintf("%s%s", o.config.CollectionPrefix, collection)
	return db.Collection(prefixed)
}

func (o *Store) Ping(ctx context.Context) error {
	if o.client == nil {
		return errors.New("store has not been initialized")
	}

	return o.client.Ping(ctx, readpref.Primary())
}

func (o *Store) Status(ctx context.Context) *Status {
	status := &Status{
		Ping: "FAILED",
	}

	if o.Ping(ctx) == nil {
		status.Ping = "OK"
	}

	return status
}

func (o *Store) Terminate(ctx context.Context) error {
	if o.client == nil {
		return nil
	}

	return o.client.Disconnect(ctx)
}
