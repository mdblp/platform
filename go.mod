module github.com/tidepool-org/platform

go 1.15

require (
	github.com/ant0ine/go-json-rest v3.3.2+incompatible
	github.com/blang/semver v3.5.1+incompatible
	github.com/fatih/color v1.10.0 // indirect
	github.com/githubnemo/CompileDaemon v1.2.1
	github.com/google/uuid v1.3.0
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.10.0 // indirect
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/kr/text v0.2.0 // indirect
	github.com/mdblp/go-common v1.1.0
	github.com/mdblp/go-json-rest v3.3.3+incompatible
	github.com/mjibson/esc v0.2.0
	github.com/onsi/ginkgo v1.16.5
	github.com/onsi/gomega v1.10.5
	github.com/prometheus/client_golang v1.11.0
	github.com/sirupsen/logrus v1.8.1
	github.com/tidepool-org/devices/api v0.0.0-20210517133954-8f12767986b5
	go.mongodb.org/mongo-driver v1.7.2
	go.uber.org/fx v1.17.1
	golang.org/x/crypto v0.0.0-20210220033148-5ea612d1eb83 // indirect
	golang.org/x/lint v0.0.0-20201208152925-83fdc39ff7b5
	golang.org/x/tools v0.1.5
	google.golang.org/grpc v1.45.0
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/tylerb/graceful.v1 v1.2.15
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)

replace gopkg.in/fsnotify.v1 v1.4.7 => gopkg.in/fsnotify/fsnotify.v1 v1.4.7

replace github.com/ugorji/go v1.1.5-pre => github.com/ugorji/go v1.1.7
