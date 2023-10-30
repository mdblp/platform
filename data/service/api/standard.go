package api

import (
	"net/http"

	"github.com/mdblp/go-json-rest/rest"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	dataClient "github.com/tidepool-org/platform/data/client"
	dataService "github.com/tidepool-org/platform/data/service"
	dataContext "github.com/tidepool-org/platform/data/service/context"
	dataStore "github.com/tidepool-org/platform/data/store"
	"github.com/tidepool-org/platform/errors"
	"github.com/tidepool-org/platform/permission"
	"github.com/tidepool-org/platform/service"
	"github.com/tidepool-org/platform/service/api"
)

type Standard struct {
	*api.API
	permissionClient permission.Client
	dataStore        dataStore.Store
	dataClient       dataClient.Client
}

func NewStandard(svc service.Service, permissionClient permission.Client,
	store dataStore.Store, dataClient dataClient.Client) (*Standard, error) {
	if permissionClient == nil {
		return nil, errors.New("permission client is missing")
	}
	if store == nil {
		return nil, errors.New("data store DEPRECATED is missing")
	}
	if dataClient == nil {
		return nil, errors.New("data client is missing")
	}

	a, err := api.New(svc)
	if err != nil {
		return nil, err
	}

	return &Standard{
		API:              a,
		permissionClient: permissionClient,
		dataStore:        store,
		dataClient:       dataClient,
	}, nil
}

func (s *Standard) DEPRECATEDInitializeRouter(routes []dataService.Route) error {
	baseRoutes := []dataService.Route{
		dataService.MakeRoute("GET", "/status", s.StatusGet),
		dataService.MakeRoute("GET", "/version", s.VersionGet),
	}

	routes = append(baseRoutes, routes...)

	var contextRoutes []*rest.Route
	for _, route := range routes {
		contextRoutes = append(contextRoutes, &rest.Route{
			HttpMethod: route.Method,
			PathExp:    route.Path,
			Func:       s.withContext(route.Handler),
		})
	}
	metricRoute := rest.Get("/metrics", func(w rest.ResponseWriter, r *rest.Request) {
		promhttp.Handler().ServeHTTP(w.(http.ResponseWriter), r.Request)
	})

	contextRoutes = append(contextRoutes, metricRoute)

	router, err := rest.MakeRouter(contextRoutes...)
	if err != nil {
		return errors.Wrap(err, "unable to create router")
	}

	s.DEPRECATEDAPI().SetApp(router)

	return nil
}

func (s *Standard) withContext(handler dataService.HandlerFunc) rest.HandlerFunc {
	return dataContext.WithContext(s.AuthClient(), s.permissionClient, s.dataStore, s.dataClient, handler)
}
