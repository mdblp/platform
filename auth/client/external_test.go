package client_test

import (
	"context"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/ghttp"

	authClient "github.com/tidepool-org/platform/auth/client"
	authTest "github.com/tidepool-org/platform/auth/test"
	"github.com/tidepool-org/platform/errors"
	errorsTest "github.com/tidepool-org/platform/errors/test"
	"github.com/tidepool-org/platform/log"
	logTest "github.com/tidepool-org/platform/log/test"
	"github.com/tidepool-org/platform/request"
	"github.com/tidepool-org/platform/test"
	testHttp "github.com/tidepool-org/platform/test/http"
)

var _ = Describe("External", func() {
	var config *authClient.ExternalConfig
	var name string
	var logger *logTest.Logger

	BeforeEach(func() {
		config = authClient.NewExternalConfig()
		config.AuthenticationConfig.UserAgent = testHttp.NewUserAgent()
		config.AuthorizationConfig.UserAgent = testHttp.NewUserAgent()
		config.ServerSessionTokenSecret = authTest.NewServiceSecret()
		name = test.RandomString()
		logger = logTest.NewLogger()
	})

	Context("NewExternal", func() {
		BeforeEach(func() {
			config.AuthenticationConfig.Address = testHttp.NewAddress()
			config.AuthorizationConfig.Address = testHttp.NewAddress()
		})

		It("returns an error when the config is missing", func() {
			config = nil
			client, err := authClient.NewExternal(config, name, logger)
			errorsTest.ExpectEqual(err, errors.New("config is missing"))
			Expect(client).To(BeNil())
		})

		It("returns an error when the name is missing", func() {
			name = ""
			client, err := authClient.NewExternal(config, name, logger)
			errorsTest.ExpectEqual(err, errors.New("name is missing"))
			Expect(client).To(BeNil())
		})

		It("returns an error when the logger is missing", func() {
			logger = nil
			client, err := authClient.NewExternal(config, name, nil)
			errorsTest.ExpectEqual(err, errors.New("logger is missing"))
			Expect(client).To(BeNil())
		})

		It("returns success", func() {
			Expect(authClient.NewExternal(config, name, logger)).ToNot(BeNil())
		})
	})

	Context("with server and new client", func() {
		var server *Server
		var requestHandlers []http.HandlerFunc
		var responseHeaders http.Header
		var client *authClient.External
		var sessionToken string
		var details request.Details
		var ctx context.Context

		BeforeEach(func() {
			server = NewServer()
			requestHandlers = nil
			responseHeaders = http.Header{"Content-Type": []string{"application/json; charset=utf-8"}}
			sessionToken = authTest.NewSessionToken()
			details = request.NewDetails(request.MethodSessionToken, "", sessionToken, "patient")
			ctx = context.Background()
			ctx = log.NewContextWithLogger(ctx, logger)
		})

		JustBeforeEach(func() {
			server.AppendHandlers(CombineHandlers(requestHandlers...))
			var err error
			config.AuthenticationConfig.Address = server.URL()
			config.AuthorizationConfig.Address = server.URL()
			client, err = authClient.NewExternal(config, name, logger)
			Expect(err).ToNot(HaveOccurred())
			Expect(client).ToNot(BeNil())
			ctx = request.NewContextWithDetails(ctx, details)
		})

		AfterEach(func() {
			if server != nil {
				server.Close()
			}
		})

		Context("EnsureAuthorized", func() {
			Context("without server response", func() {
				AfterEach(func() {
					Expect(server.ReceivedRequests()).To(BeEmpty())
				})

				It("returns an error when the context is missing", func() {
					ctx = nil
					errorsTest.ExpectEqual(client.EnsureAuthorized(ctx), errors.New("context is missing"))
				})

				It("returns an error when the details are missing", func() {
					ctx = request.NewContextWithDetails(ctx, nil)
					errorsTest.ExpectEqual(client.EnsureAuthorized(ctx), request.ErrorUnauthorized())
				})

				It("returns successfully when the details are for a user", func() {
					ctx = request.NewContextWithDetails(ctx, request.NewDetails(request.MethodSessionToken, authTest.RandomUserID(), sessionToken, "patient"))
					Expect(client.EnsureAuthorized(ctx)).To(Succeed())
				})

				It("returns successfully when the details are for a service", func() {
					Expect(client.EnsureAuthorized(ctx)).To(Succeed())
				})
			})
		})

		Context("EnsureAuthorizedService", func() {
			Context("without server response", func() {
				AfterEach(func() {
					Expect(server.ReceivedRequests()).To(BeEmpty())
				})

				It("returns an error when the context is missing", func() {
					ctx = nil
					errorsTest.ExpectEqual(client.EnsureAuthorizedService(ctx), errors.New("context is missing"))
				})

				It("returns an error when the details are missing", func() {
					ctx = request.NewContextWithDetails(ctx, nil)
					errorsTest.ExpectEqual(client.EnsureAuthorizedService(ctx), request.ErrorUnauthorized())
				})

				It("returns an error when the details are for not a service", func() {
					ctx = request.NewContextWithDetails(ctx, request.NewDetails(request.MethodSessionToken, authTest.RandomUserID(), sessionToken, "patient"))
					errorsTest.ExpectEqual(client.EnsureAuthorizedService(ctx), request.ErrorUnauthorized())
				})

				It("returns successfully when the details are for a service", func() {
					Expect(client.EnsureAuthorizedService(ctx)).To(Succeed())
				})
			})
		})

		Context("EnsureAuthorizedUser", func() {
			var requestUserID string
			var targetUserID string

			BeforeEach(func() {
				requestUserID = authTest.RandomUserID()
				targetUserID = authTest.RandomUserID()
				details = request.NewDetails(request.MethodSessionToken, requestUserID, sessionToken, "patient")
			})

			Context("without server response", func() {
				AfterEach(func() {
					Expect(server.ReceivedRequests()).To(BeEmpty())
				})

				It("returns an error when the context is missing", func() {
					ctx = nil
					userID, err := client.EnsureAuthorizedUser(ctx, targetUserID)
					errorsTest.ExpectEqual(err, errors.New("context is missing"))
					Expect(userID).To(BeEmpty())
				})

				It("returns an error when the target user id is missing", func() {
					targetUserID = ""
					userID, err := client.EnsureAuthorizedUser(ctx, targetUserID)
					errorsTest.ExpectEqual(err, errors.New("target user id is missing"))
					Expect(userID).To(BeEmpty())
				})

				It("returns an error when the details are missing", func() {
					ctx = request.NewContextWithDetails(ctx, nil)
					userID, err := client.EnsureAuthorizedUser(ctx, targetUserID)
					errorsTest.ExpectEqual(err, request.ErrorUnauthorized())
					Expect(userID).To(BeEmpty())
				})

				It("returns successfully when the details are for a service and authorized permission is custodian", func() {
					ctx = request.NewContextWithDetails(ctx, request.NewDetails(request.MethodSessionToken, "", sessionToken, "patient"))
					Expect(client.EnsureAuthorizedUser(ctx, targetUserID)).To(Equal(""))
				})

				It("returns successfully when the details are for the target user", func() {
					ctx = request.NewContextWithDetails(ctx, request.NewDetails(request.MethodSessionToken, targetUserID, sessionToken, "patient"))
					Expect(client.EnsureAuthorizedUser(ctx, targetUserID)).To(Equal(targetUserID))
				})
			})

			Context("with server response when the details are not for the target user", func() {
				Context("with a successful response with incorrect permissions", func() {
					BeforeEach(func() {
						requestHandlers = append(requestHandlers, RespondWith(http.StatusOK, `{"view": {}}`, responseHeaders))
					})

					It("returns an error", func() {
						ctx = request.NewContextWithDetails(ctx, request.NewDetails(request.MethodSessionToken, "unknownId", sessionToken, "patient"))
						userID, err := client.EnsureAuthorizedUser(ctx, targetUserID)
						errorsTest.ExpectEqual(err, request.ErrorUnauthorized())
						Expect(userID).To(Equal(""))
					})
				})
			})
		})
	})
})
