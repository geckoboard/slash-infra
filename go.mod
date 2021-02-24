module github.com/geckoboard/slash-infra

go 1.15

require (
	github.com/aws/aws-sdk-go v1.16.8-0.20181217231416-7c690b7a4c41
	github.com/bugsnag/bugsnag-go v1.4.0
	github.com/bugsnag/panicwrap v1.2.1-0.20180510051541-1d162ee1264c // indirect
	github.com/davecgh/go-spew v1.1.1
	github.com/gofrs/uuid v3.3.0+incompatible // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2
	github.com/jmespath/go-jmespath v0.0.0-20180206201540-c2b33e8439af // indirect
	github.com/joho/godotenv v1.3.1-0.20181120194748-69ed1d913aa8
	github.com/julienschmidt/httprouter v1.2.1-0.20181021223831-26a05976f9bf
	github.com/kardianos/osext v0.0.0-20170510131534-ae77be60afb1 // indirect
	github.com/mitchellh/mapstructure v1.1.2
	github.com/spf13/cobra v1.1.1
	go.opentelemetry.io/contrib/instrumentation/net/http v0.11.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.17.0
	go.opentelemetry.io/otel v0.17.0
	go.opentelemetry.io/otel/exporters/otlp v0.17.0
	go.opentelemetry.io/otel/exporters/stdout v0.17.0
	go.opentelemetry.io/otel/sdk v0.17.0
	google.golang.org/grpc v1.35.0
)
