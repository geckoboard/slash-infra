package slackbot

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/geckoboard/slash-infra/awsutil"
	"github.com/geckoboard/slash-infra/search"
	"github.com/geckoboard/slash-infra/slackutil"
	"github.com/spf13/cobra"
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
