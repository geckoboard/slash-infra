package debug

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/davecgh/go-spew/spew"
	"github.com/geckoboard/slash-infra/awsevents"
	"github.com/spf13/cobra"
)

var (
	eventbridgeCmd = &cobra.Command{
		Use:   "eventbridge",
		Short: "Tail the eventbridge sqs queue",
		Run: func(cmd *cobra.Command, args []string) {
			sess, _ := session.NewSession(&aws.Config{
				Credentials: credentials.NewEnvCredentials(),
				Region:      aws.String("eu-west-2"),
			})

			consumer := &awsevents.SQSConsumer{
				SQSQueueURL: "https://sqs.eu-west-2.amazonaws.com/691393038123/matt-eventbridge-tests",
				SQSService:  sqs.New(sess),
			}

			log.Print("starting poll loop")
			err := consumer.PollLoop(context.TODO(), &spewer{})

			log.Fatal(err)
		},
	}
)

type spewer struct{}

type genericevent interface {
	GenericAWSEventInfo() awsevents.AWSEvent
}

func (spewer) HandleAWSEvent(ev interface{}) {
	switch x := ev.(type) {
	case awsevents.EC2InstanceLaunchSuccessful:
		fmt.Printf("instance %q launched successfully at %s\n", x.Detail.EC2InstanceID, x.Detail.StartTime)
	case awsevents.EC2InstanceLaunchUnsuccessful:
		fmt.Printf("instance %q did not launch successfully\n", x.Detail.EC2InstanceID)
	case awsevents.EC2InstanceTerminateSuccessful:
		fmt.Printf("instance %q was terminated: %q\n", x.Detail.EC2InstanceID, x.Detail.Cause)
		spew.Dump(x)
		os.Exit(1)
	case nil:
		fmt.Println("WTFFFF")
	default:
		//genericEvent := x.(genericevent).GenericAWSEventInfo()

		//fmt.Println("unable to deal with event", genericEvent.DetailType)
	}
}
