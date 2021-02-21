package debug

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

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
		if strings.Contains(x.Detail.Cause, "the scheduled action") {
			fmt.Printf("instance %q scheduled start\n", x.Detail.EC2InstanceID)
			return
		}
		fmt.Printf("instance %q launched successfully at %s\n", x.Detail.EC2InstanceID, x.Detail.StartTime)
	case awsevents.EC2InstanceLaunchUnsuccessful:
		fmt.Printf("instance %q did not launch successfully\n", x.Detail.EC2InstanceID)
	case awsevents.EC2InstanceTerminateSuccessful:
		// This change is in response to a scheduled scale-down
		if strings.Contains(x.Detail.Cause, "the scheduled action") {
			fmt.Printf("instance %q scheduled scale down termination\n", x.Detail.EC2InstanceID)
			return
		}
		if strings.Contains(x.Detail.Cause, "taken out of service in response to an EC2 health check") {
			fmt.Printf("instance %q failed health check\n", x.Detail.EC2InstanceID)
			return
		}
		fmt.Printf("instance %q was terminated: %q\n", x.Detail.EC2InstanceID, x.Detail.Cause)
		spew.Dump(x)
		os.Exit(1)
	default:
		//genericEvent := x.(genericevent).GenericAWSEventInfo()

		//fmt.Println("unable to deal with event", genericEvent.DetailType)
	}
}
