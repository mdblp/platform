package service

import (
	"github.com/tidepool-org/platform/auth"
	dataClient "github.com/tidepool-org/platform/data/client"
	"github.com/tidepool-org/platform/data/deduplicator"
	dataStore "github.com/tidepool-org/platform/data/store"
	"github.com/tidepool-org/platform/permission"
	"github.com/tidepool-org/platform/service"
	syncTaskStore "github.com/tidepool-org/platform/synctask/store"
)

type Context interface {
	service.Context

	AuthClient() auth.Client
	PermissionClient() permission.Client

	DataDeduplicatorFactory() deduplicator.Factory

	DataRepository() dataStore.DataRepository
	SyncTaskRepository() syncTaskStore.SyncTaskRepository

	DataClient() dataClient.Client
}

type HandlerFunc func(context Context)
