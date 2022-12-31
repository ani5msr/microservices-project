module github.com/ani5msr/microservices-project

go 1.16

require (
	github.com/Masterminds/squirrel v1.5.3
	github.com/go-kit/kit v0.12.0
	github.com/gorilla/mux v1.8.0
	github.com/lib/pq v1.10.7
	github.com/nats-io/nats.go v1.13.0
	github.com/prometheus/client_golang v1.11.0
	github.com/uber/jaeger-client-go v2.30.0+incompatible
	gopkg.in/yaml.v2 v2.4.0

)

require github.com/uber/jaeger-lib v2.4.1+incompatible // indirect

require (
	github.com/go-kit/log v0.2.1
	github.com/google/go-cmp v0.5.9 // indirect
	github.com/onsi/ginkgo v1.16.5
	github.com/onsi/gomega v1.13.0
	github.com/opentracing/opentracing-go v1.2.0
	github.com/stretchr/testify v1.8.0 // indirect
	golang.org/x/crypto v0.1.0 // indirect
	golang.org/x/sys v0.3.0 // indirect
	google.golang.org/protobuf v1.28.0 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
)
