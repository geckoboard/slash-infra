package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"

	"github.com/geckoboard/slash-infra/cmd/debug"
	"github.com/geckoboard/slash-infra/cmd/slackbot"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp"
	"go.opentelemetry.io/otel/exporters/otlp/otlpgrpc"
	"go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc/credentials"
)

var (
	rootCmd = &cobra.Command{
		Use:   "slash-infra",
		Short: "A slackbot for working with AWS infrastructure",
	}
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Llongfile)

	// In development it's easier to store environment variables in a .env folder
	godotenv.Load()

	// initialize trace provider.
	cleanup := initTracer()
	defer cleanup()

	rootCmd.AddCommand(slackbot.Command())
	rootCmd.AddCommand(debug.Command())

	rootCmd.Execute()
}

// initTracer creates and registers trace provider instance.
func initTracer() func() {
	apikey := os.Getenv("HONEYCOMB_KEY")
	dataset := os.Getenv("HONEYCOMB_DATASET")

	exporter, err := otlp.NewExporter(
		context.Background(),
		otlpgrpc.NewDriver(
			otlpgrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, "")),
			otlpgrpc.WithEndpoint("api.honeycomb.io:443"),
			otlpgrpc.WithHeaders(map[string]string{
				"x-honeycomb-team":    apikey,
				"x-honeycomb-dataset": dataset,
			}),
		),
	)
	if err != nil {
		log.Fatal(err)
	}

	otel.SetTracerProvider(
		trace.NewTracerProvider(
			trace.WithConfig(trace.Config{DefaultSampler: trace.AlwaysSample()}),
			trace.WithSpanProcessor(trace.NewBatchSpanProcessor(exporter)),
		),
	)
	return func() {
		exporter.Shutdown(context.Background())
	}
}
