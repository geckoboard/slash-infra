package awsevents

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/bugsnag/bugsnag-go"
)

const (
	// The number of seconds to "reserve" messages from the queue before other consumers can access them
	visibilityTimeoutSeconds = 10

	// How long should each poll to the SQS API last when receiving messages
	waitTimeSeconds = 10

	// How many messages can we accept in one go from each `ReceiveMessage` API call
	maxNumberOfMessages = 10
)

type EventHandler interface {
	HandleAWSEvent(interface{})
}

// Consumer polls an SQS queue for eventbridge events
type SQSConsumer struct {
	SQSQueueURL string

	SQSService *sqs.SQS
}

func (s *SQSConsumer) PollLoop(ctx context.Context, handler EventHandler) error {
	for {
		select {
		case <-ctx.Done():
			// cancelllllled
			return nil
		default:
			// All's good

		}

		req := &sqs.ReceiveMessageInput{
			QueueUrl: aws.String(s.SQSQueueURL),
			AttributeNames: []*string{
				aws.String("QueueAttributeName"), // Required
				aws.String("SentTimestamp"),
				aws.String("ApproximateReceiveCount"),
				aws.String("ApproximateFirstReceiveTimestamp"),
			},
			VisibilityTimeout:   aws.Int64(int64(visibilityTimeoutSeconds)),
			WaitTimeSeconds:     aws.Int64(int64(waitTimeSeconds)),
			MaxNumberOfMessages: aws.Int64(maxNumberOfMessages),
		}

		msgList, err := s.SQSService.ReceiveMessageWithContext(aws.Context(ctx), req)
		if err != nil {
			// The AWS-SDK wraps errors so we need to peek into the message
			if errors.Is(err, context.Canceled) {
				return context.Canceled
			}

			bugsnag.Notify(err, ctx)
			return err
		}

		for _, msg := range msgList.Messages {
			var raw map[string]interface{}

			json.NewDecoder(strings.NewReader(*msg.Body)).Decode(&raw)

			ev, err := decode(raw)

			if err != nil {
				bugsnag.Notify(err, ctx, bugsnag.MetaData{
					"Event": {
						"UntypedStructure": raw,
					},
				})
			}

			if ev != nil {
				handler.HandleAWSEvent(ev)
			}
			//s.SQSService.DeleteMessageWithContext(aws.Context(ctx), &sqs.DeleteMessageInput{
			//	QueueUrl:      aws.String(s.SQSQueueURL),
			//	ReceiptHandle: msg.ReceiptHandle,
			//})
		}

	}
}
