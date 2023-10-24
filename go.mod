module github.com/tidepool-org/platform

go 1.15

require (
	github.com/ant0ine/go-json-rest v3.3.2+incompatible // indirect
	github.com/auth0/go-jwt-middleware/v2 v2.0.1
	github.com/blang/semver v3.5.1+incompatible
	github.com/fatih/color v1.10.0 // indirect
	github.com/githubnemo/CompileDaemon v1.2.1
	github.com/google/go-cmp v0.5.9 // indirect
	github.com/google/uuid v1.3.0
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/kr/pretty v0.3.0 // indirect
	github.com/mattn/go-isatty v0.0.16 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.4 // indirect
	github.com/mdblp/go-db v1.0.1
	github.com/mdblp/go-json-rest v3.3.3+incompatible
	github.com/mjibson/esc v0.2.0
	github.com/onsi/ginkgo v1.16.5
	github.com/onsi/gomega v1.10.5
	github.com/prometheus/client_golang v1.14.0
	github.com/rogpeppe/go-internal v1.8.0 // indirect
	github.com/sirupsen/logrus v1.9.2
	github.com/stretchr/testify v1.8.1 // indirect
	go.mongodb.org/mongo-driver v1.11.6
	go.uber.org/automaxprocs v1.5.3
	go.uber.org/fx v1.17.1
	golang.org/x/lint v0.0.0-20201208152925-83fdc39ff7b5
	golang.org/x/tools v0.6.0
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/tylerb/graceful.v1 v1.2.15
)

replace gopkg.in/fsnotify.v1 v1.4.7 => gopkg.in/fsnotify/fsnotify.v1 v1.4.7

replace github.com/ugorji/go v1.1.5-pre => github.com/ugorji/go v1.1.7
