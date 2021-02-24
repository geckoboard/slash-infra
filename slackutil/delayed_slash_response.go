package slackutil

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"

	bugsnag "github.com/bugsnag/bugsnag-go"
	"go.opentelemetry.io/otel/trace"
)

var forceShowSlashCommandInChannelResponse = Response{ResponseType: ResponseInChannel}

type DelayedSlashResponse struct {
	// A mesage to send the user while we're preparing a response to
	PendingResponse Response

	// Should the command be visible to all other users in the channel?
	// Doing this changes how we respond to the slash command webhook
	ShowSlashCommandInChannel bool

	Handler func(context.Context, SlashCommandRequest, MessageResponder)
}

func (d DelayedSlashResponse) Run(ctx context.Context, slackClient *Client, w http.ResponseWriter, command SlashCommandRequest) {
	ctx = trace.ContextWithSpan(context.Background(), trace.SpanFromContext(ctx))

	if d.ShowSlashCommandInChannel {
		// By default slack treats responses to slash commands as
		// "ephemeral", and will prevent the slash command from showing
		// up in chat history.
		// If we signal in our response that the slash command is not
		// ephemeral, then the original command will appear in the chat
		// history.
		// https://api.slack.com/slash-commands#responding_immediate_response
		RespondWith(w, forceShowSlashCommandInChannelResponse)
	}

	// We run the handler in a goroutine so that we can confirm receipt of slack's
	// slash command webhook (by returning 200 OK) as soon as possible
	go d.runHandler(ctx, slackClient, command)
}

func (d DelayedSlashResponse) runHandler(ctx context.Context, slackClient *Client, command SlashCommandRequest) {
	responder := MessageResponder{command, slackClient}

	done := make(chan struct{})

	// Not using a waitgroup here as we don't really care about cleaning up this goroutine
	go func() {
		defer bugsnag.AutoNotify(ctx)
		d.Handler(ctx, command, responder)
		close(done)
	}()

	notifyUserTimeout := time.After(700 * time.Millisecond)

	for {
		select {
		case <-done:
			return
		case <-notifyUserTimeout:
			responder.EphemeralResponse(ctx, d.PendingResponse)
		}
	}
}

type MessageResponder struct {
	command     SlashCommandRequest
	slackClient *Client
}

func (m MessageResponder) EphemeralResponse(ctx context.Context, resp Response) {
	resp.ResponseType = ResponseEphemeral

	b, err := json.Marshal(&resp)
	if err != nil {
		panic(err)
	}

	r, err := http.NewRequest("POST", m.command.ResponseURL, bytes.NewReader(b))
	if err != nil {
		panic(err)
	}
	r = r.WithContext(ctx)

	httpresp, err := m.slackClient.httpClient.Do(r)
	if err != nil {
		panic(err)
	}
	httpresp.Body.Close()
}

func (m MessageResponder) PublicResponse(ctx context.Context, resp Response) {
	resp.ResponseType = ResponseInChannel

	b, err := json.Marshal(&resp)
	if err != nil {
		panic(err)
	}

	r, err := http.NewRequest("POST", m.command.ResponseURL, bytes.NewReader(b))
	if err != nil {
		panic(err)
	}
	r = r.WithContext(ctx)

	httpresp, err := m.slackClient.httpClient.Do(r)
	if err != nil {
		panic(err)
	}
	httpresp.Body.Close()
}

type SlashCommandResponder interface {
	PublicResponse(Response)
}
