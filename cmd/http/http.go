package http

import (
	"log"
	"net/http"
	"os"

	"github.com/geckoboard/slash-infra/slackutil"
	"github.com/spf13/cobra"
)

func Command() *cobra.Command {
	return &cobra.Command{
		Use:   "http",
		Short: "Run the slack bot HTTP server",
		Run: func(cmd *cobra.Command, args []string) {
			server := makeHttpHandler()

			handler := slackutil.VerifyRequestSignature(os.Getenv("SLACK_SIGNING_SECRET"))(server)

			port := os.Getenv("PORT")
			if port == "" {
				port = "8090"
			}

			log.Fatal(http.ListenAndServe(":"+port, handler))
		},
	}

}
