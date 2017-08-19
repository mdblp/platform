package service

import (
	"github.com/tidepool-org/platform/auth"
	"github.com/tidepool-org/platform/config"
	"github.com/tidepool-org/platform/log"
	"github.com/tidepool-org/platform/version"
)

type Service interface {
	VersionReporter() version.Reporter
	ConfigReporter() config.Reporter
	Logger() log.Logger

	AuthClient() auth.Client
}
