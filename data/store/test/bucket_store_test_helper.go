package test

import (
	goComMgo "github.com/mdblp/go-db/mongo"
)

// NewConfig creates a test Mongo configuration
func NewConfig() *goComMgo.Config {
	conf := &goComMgo.Config{}
	conf.FromEnv()
	conf.Database = "data_read_test"

	return conf
}
