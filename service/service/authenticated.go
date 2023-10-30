package service

import (
	"github.com/tidepool-org/platform/application"
	authClient "github.com/tidepool-org/platform/auth/client"
	"github.com/tidepool-org/platform/errors"
	"github.com/tidepool-org/platform/platform"
)

type Authenticated struct {
	*Service
	authClient *authClient.Client
}

func NewAuthenticated() *Authenticated {
	return &Authenticated{
		Service: New(),
	}
}

func (a *Authenticated) Initialize(provider application.Provider) error {
	if err := a.Service.Initialize(provider); err != nil {
		return err
	}

	return a.initializeAuthClient()
}

func (a *Authenticated) Terminate() {
	a.terminateAuthClient()

	a.Service.Terminate()
}

func (a *Authenticated) initializeAuthClient() error {
	a.Logger().Debug("Loading auth client config")

	userAgent := a.UserAgent()
	cfg := authClient.NewConfig()
	cfg.ExternalConfig.AuthenticationConfig.UserAgent = userAgent
	cfg.ExternalConfig.AuthorizationConfig.UserAgent = userAgent
	if err := cfg.Load(a.ConfigReporter().WithScopes("auth", "client")); err != nil {
		return errors.Wrap(err, "unable to load auth client config")
	}

	a.Logger().Debug("Creating auth client")

	clnt, err := authClient.NewClient(cfg, platform.AuthorizeAsService, a.Name(), a.Logger())
	if err != nil {
		return errors.Wrap(err, "unable to create auth client")
	}
	a.authClient = clnt

	a.SetAuthClient(a.authClient)

	return nil
}

func (a *Authenticated) terminateAuthClient() {
	if a.authClient != nil {
		a.Logger().Debug("Destroying auth client")
		a.authClient = nil

		a.SetAuthClient(nil)
	}
}
