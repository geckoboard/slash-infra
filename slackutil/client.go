package slackutil

import (
	"net/http"

	"github.com/hashicorp/go-cleanhttp"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type Client struct {
	httpClient *http.Client
}

func New() *Client {
	h := cleanhttp.DefaultPooledClient()
	h.Transport = otelhttp.NewTransport(
		h.Transport,
		otelhttp.WithSpanNameFormatter(func(operation string, r *http.Request) string {
			return "slack client"
		}),
	)

	return &Client{
		httpClient: h,
	}
}
