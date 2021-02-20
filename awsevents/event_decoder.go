package awsevents

import (
	"errors"
	"os"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/mitchellh/mapstructure"
)

var (
	ErrCannotDecodeAWSEvent = errors.New("cannot decode aws event")
)

func decode(raw map[string]interface{}) (interface{}, error) {
	switch raw["detail-type"] {
	case "AWS API Call via CloudTrail":
		var out CloudTrailAPICall

		err := decodeMap(raw, &out)

		return out, err
	case "EC2 Instance Terminate Unsuccessful":
		var out EC2InstanceLaunchUnsuccessful

		err := decodeMap(raw, &out)

		return out, err
	case "EC2 Instance Launch Successful":
		var out EC2InstanceLaunchSuccessful

		err := decodeMap(raw, &out)

		return out, err
	case "EC2 Instance State-change Notification":
		var out EC2InstanceStateChangeNotification

		err := decodeMap(raw, &out)

		return out, err
	case "AWS Service Event via CloudTrail":
		var out CloudTrailServiceEvent

		err := decodeMap(raw, &out)

		return out, err
	case "EC2 Instance-launch Lifecycle Action":
		var out EC2InstanceLaunchLifecycleAction

		err := decodeMap(raw, &out)

		return out, err
	case "EC2 Instance Terminate Successful":
		var out EC2InstanceTerminateSuccessful

		err := decodeMap(raw, &out)

		return out, err
	case "EBS Volume Notification":
		var out EBSVolumeEventNotification

		err := decodeMap(raw, &out)

		return out, err
	case "EC2 Spot Instance Request Fulfillment":
		var out EC2SpotInstanceRequestFulfillment

		err := decodeMap(raw, &out)

		return out, err
	case "EC2 Spot Instance Interruption Warning":
		var out EC2SpotInstanceInterruptionWarning

		err := decodeMap(raw, &out)

		return out, err
	case "EC2 Instance Rebalance Recommendation":
		var out EC2InstanceRebalanceRecommendation

		err := decodeMap(raw, &out)

		return out, err
	default:
		spew.Dump(raw)
		os.Exit(1)
	}

	return nil, ErrCannotDecodeAWSEvent
}

func decodeMap(raw map[string]interface{}, out interface{}) error {
	decoder, _ := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			mapstructure.StringToTimeHookFunc(time.RFC3339),
			mapstructure.StringToIPHookFunc(),
		),
		Result: out,
	})

	return decoder.Decode(raw)
}
