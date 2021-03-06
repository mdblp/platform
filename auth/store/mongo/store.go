package mongo

import (
	mgo "github.com/globalsign/mgo"

	"github.com/tidepool-org/platform/auth/store"
	"github.com/tidepool-org/platform/log"
	storeStructuredMongo "github.com/tidepool-org/platform/store/structured/mongo"
)

type Store struct {
	*storeStructuredMongo.Store
}

var (
	authIndexes = map[string][]mgo.Index{
		"provider_sessions": {
			{Key: []string{"id"}, Unique: true, Background: true},
			{Key: []string{"userId"}, Background: true},
			{Key: []string{"userId", "type", "name"}, Unique: true, Background: true},
		},
		"restricted_tokens": {
			{Key: []string{"id"}, Unique: true, Background: true},
			{Key: []string{"userId"}, Background: true},
		},
	}
)

func NewStore(cfg *storeStructuredMongo.Config, lgr log.Logger) (*Store, error) {
	if cfg != nil {
		cfg.Indexes = authIndexes
	}
	str, err := storeStructuredMongo.NewStore(cfg, lgr)
	if err != nil {
		return nil, err
	}

	return &Store{
		Store: str,
	}, nil
}

func (s *Store) NewProviderSessionSession() store.ProviderSessionSession {
	return s.providerSessionSession()
}

func (s *Store) NewRestrictedTokenSession() store.RestrictedTokenSession {
	return s.restrictedTokenSession()
}

func (s *Store) providerSessionSession() *ProviderSessionSession {
	return &ProviderSessionSession{
		Session: s.Store.NewSession("provider_sessions"),
	}
}

func (s *Store) restrictedTokenSession() *RestrictedTokenSession {
	return &RestrictedTokenSession{
		Session: s.Store.NewSession("restricted_tokens"),
	}
}
