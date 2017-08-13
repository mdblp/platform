package mongo_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"time"

	"github.com/tidepool-org/platform/config/test"
	"github.com/tidepool-org/platform/pointer"
	baseConfig "github.com/tidepool-org/platform/store/mongo"
	"github.com/tidepool-org/platform/user/store/mongo"
)

var _ = Describe("Config", func() {
	Context("NewConfig", func() {
		It("returns a new config with default values", func() {
			config := mongo.NewConfig()
			Expect(config).ToNot(BeNil())
			Expect(config.PasswordSalt).To(Equal(""))
		})
	})

	Context("with new config", func() {
		var config *mongo.Config

		BeforeEach(func() {
			config = mongo.NewConfig()
			Expect(config).ToNot(BeNil())
		})

		Context("Load", func() {
			var configReporter *test.Reporter

			BeforeEach(func() {
				configReporter = test.NewReporter()
				configReporter.Config["password_salt"] = "pepper"
			})

			It("returns an error if config reporter is missing", func() {
				Expect(config.Load(nil)).To(MatchError("mongo: config reporter is missing"))
			})

			It("returns an error if base config is missing", func() {
				config.Config = nil
				Expect(config.Load(configReporter)).To(MatchError("mongo: config is missing"))
			})

			It("returns an error if base config returns an error", func() {
				configReporter.Config["tls"] = "abc"
				Expect(config.Load(configReporter)).To(MatchError("mongo: tls is invalid"))
			})

			It("uses default password salt if not set", func() {
				delete(configReporter.Config, "password_salt")
				Expect(config.Load(configReporter)).To(Succeed())
				Expect(config.PasswordSalt).To(BeEmpty())
			})

			It("returns successfully and uses values from config reporter", func() {
				Expect(config.Load(configReporter)).To(Succeed())
				Expect(config.PasswordSalt).To(Equal("pepper"))
			})
		})

		Context("with valid values", func() {
			BeforeEach(func() {
				config.Config = baseConfig.NewConfig()
				config.Config.Addresses = []string{"1.2.3.4", "5.6.7.8"}
				config.Config.TLS = false
				config.Config.Database = "database"
				config.Config.CollectionPrefix = "collection_prefix"
				config.Config.Username = pointer.String("username")
				config.Config.Password = pointer.String("password")
				config.Config.Timeout = 5 * time.Second
				config.PasswordSalt = "pepper"
			})

			Context("Validate", func() {
				It("return success if all are valid", func() {
					Expect(config.Validate()).To(Succeed())
				})

				It("returns an error if the base config is missing", func() {
					config.Config = nil
					Expect(config.Validate()).To(MatchError("mongo: config is missing"))
				})

				It("returns an error if the base config is not valid", func() {
					config.Config.Addresses = nil
					Expect(config.Validate()).To(MatchError("mongo: addresses is missing"))
				})

				It("returns an error if the password salt is missing", func() {
					config.PasswordSalt = ""
					Expect(config.Validate()).To(MatchError("mongo: password salt is missing"))
				})
			})
		})
	})
})
