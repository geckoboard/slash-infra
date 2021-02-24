package slackbot

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/geckoboard/slash-infra/awsutil"
	"github.com/geckoboard/slash-infra/search"
	"github.com/geckoboard/slash-infra/slackutil"
	"github.com/hashicorp/go-cleanhttp"
	"github.com/spf13/cobra"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type SlashInfraSearcher interface {
	Search(ctx context.Context, query string) []search.ResultSet
}

func Command() *cobra.Command {
	return &cobra.Command{
		Use:   "slackbot",
		Short: "Run the full slack bot (including http server, and eventbridge listener)",
		Run: func(cmd *cobra.Command, args []string) {
			creds, err := awsutil.DetectAWSCredentials()
			if err != nil {
				log.Fatal("unable to detect AWS credentials", err)
			}

			awsHTTPClient := cleanhttp.DefaultPooledClient()
			awsHTTPClient.Transport = otelhttp.NewTransport(
				awsHTTPClient.Transport,
				otelhttp.WithSpanNameFormatter(func(operation string, r *http.Request) string {
					return "aws sdk"
				}),
			)

			for i, _ := range creds {
				// Add telemetry to all the SDK http clients
				creds[i].SDKSession = creds[i].SDKSession.Copy(&aws.Config{
					HTTPClient: awsHTTPClient,
				})
			}

			searchers := []SlashInfraSearcher{
				search.NewEC2Resolver(creds),
			}

			server := makeHttpHandler(searchers)

			handler := slackutil.VerifyRequestSignature(os.Getenv("SLACK_SIGNING_SECRET"))(server)

			port := os.Getenv("PORT")
			if port == "" {
				port = "8090"
			}

			log.Fatal(http.ListenAndServe(":"+port, handler))
		},
	}

}
