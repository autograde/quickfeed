module github.com/autograde/quickfeed

go 1.14

require (
	github.com/360EntSecGroup-Skylar/excelize v1.4.1
	github.com/Microsoft/go-winio v0.5.0 // indirect
	github.com/alta/protopatch v0.3.4
	github.com/autograde/quickfeed/kit v0.1.0
	github.com/containerd/containerd v1.5.3 // indirect
	github.com/docker/docker v20.10.7+incompatible
	github.com/docker/go-connections v0.4.0 // indirect
	github.com/fatih/color v1.9.0 // indirect
	github.com/golang/protobuf v1.5.2
	github.com/google/go-cmp v0.5.6
	github.com/google/go-github/v35 v35.3.0
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/gorilla/mux v1.8.0 // indirect
	github.com/gorilla/sessions v1.2.1
	github.com/gosimple/slug v1.9.0
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-hclog v0.12.2 // indirect
	github.com/hashicorp/go-retryablehttp v0.7.0 // indirect
	github.com/labstack/echo-contrib v0.11.0
	github.com/labstack/echo/v4 v4.4.0
	github.com/markbates/goth v1.67.1
	github.com/mattn/go-isatty v0.0.13 // indirect
	github.com/mattn/go-sqlite3 v2.0.3+incompatible // indirect
	github.com/moby/term v0.0.0-20210619224110-3f7ff695adc6 // indirect
	github.com/morikuni/aec v1.0.0 // indirect
	github.com/prometheus/client_golang v1.11.0
	github.com/prometheus/common v0.29.0 // indirect
	github.com/prometheus/procfs v0.7.0 // indirect
	github.com/sirupsen/logrus v1.8.1 // indirect
	github.com/uber-go/atomic v1.4.0 // indirect
	github.com/urfave/cli v1.22.4
	github.com/xanzy/go-gitlab v0.50.1
	go.uber.org/atomic v1.8.0 // indirect
	go.uber.org/multierr v1.7.0 // indirect
	go.uber.org/zap v1.18.1
	golang.org/x/crypto v0.0.0-20210711020723-a769d52b0f97 // indirect
	golang.org/x/net v0.0.0-20210614182718-04defd469f4e // indirect
	golang.org/x/oauth2 v0.0.0-20210628180205-a41e5a781914
	golang.org/x/sys v0.0.0-20210630005230-0f9fa26af87c // indirect
	golang.org/x/time v0.0.0-20210611083556-38a9dc6acbc6 // indirect
	golang.org/x/tools v0.1.5 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto v0.0.0-20210714021259-044028024a4f // indirect
	google.golang.org/grpc v1.39.0
	google.golang.org/grpc/cmd/protoc-gen-go-grpc v1.1.0
	google.golang.org/protobuf v1.27.1
	gopkg.in/yaml.v2 v2.4.0
	gorm.io/driver/sqlite v1.1.4
	gorm.io/gorm v1.21.12
)

replace github.com/autograde/quickfeed/kit => ./kit
