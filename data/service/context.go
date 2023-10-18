package service

import (
	"github.com/tidepool-org/platform/auth"
	dataClient "github.com/tidepool-org/platform/data/client"
	dataStore "github.com/tidepool-org/platform/data/store"
	"github.com/tidepool-org/platform/permission"
	"github.com/tidepool-org/platform/service"
)

type Context interface {
	service.Context

	AuthClient() auth.Client
	PermissionClient() permission.Client

	DataRepository() dataStore.DataRepository

	DataClient() dataClient.Client
}

type HandlerFunc func(context Context)
