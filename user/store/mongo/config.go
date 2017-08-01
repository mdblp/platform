package mongo

import (
	"github.com/tidepool-org/platform/config"
	"github.com/tidepool-org/platform/errors"
	"github.com/tidepool-org/platform/store/mongo"
)

type Config struct {
	*mongo.Config
	PasswordSalt string
}

func NewConfig() *Config {
	return &Config{
		Config: mongo.NewConfig(),
	}
}

func (c *Config) Load(configReporter config.Reporter) error {
	if c.Config == nil {
		return errors.New("mongo", "config is missing")
	}
	if err := c.Config.Load(configReporter); err != nil {
		return err
	}

	c.PasswordSalt = configReporter.StringOrDefault("password_salt", "")

	return nil
}

func (c *Config) Validate() error {
	if c.Config == nil {
		return errors.New("mongo", "config is missing")
	}
	if err := c.Config.Validate(); err != nil {
		return err
	}

	if c.PasswordSalt == "" {
		return errors.New("mongo", "password salt is missing")
	}

	return nil
}

func (c *Config) Clone() *Config {
	return &Config{
		Config:       c.Config.Clone(),
		PasswordSalt: c.PasswordSalt,
	}
}