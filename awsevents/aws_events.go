package awsevents

import (
	"time"
)

// AWSEvent contains fields that are standard in all AWS events from eventbridge
// https://docs.aws.amazon.com/eventbridge/latest/userguide/aws-events.html
type AWSEvent struct {
	// By default, this is set to 0 (zero) in all events.
	Version string

	// A unique value is generated for every event. This can be helpful in
	// tracing events as they move through rules to targets, and are
	// processed.
	ID string

	// Identifies, in combination with the source field, the fields and values that
	// appear in the detail field.
	//
	// All events that are delivered via CloudTrail have AWS API Call via
	// CloudTrail as the value for detail-type. For more information, see Events
	// Delivered Via CloudTrail.
	DetailType string `mapstructure:"detail-type"`

	// Identifies the service that sourced the event. All events sourced
	// from within AWS begin with "aws."
	Source string

	// The 12-digit number identifying an AWS account.
	AWSAccountID string `mapstructure:"account"`

	// The event timestamp, which can be specified by the service
	// originating the event. If the event spans a time interval, the
	// service might choose to report the start time, so this value can be
	// noticeably before the time the event is actually received.
	EventTimestamp time.Time `mapstructure:"time"`

	// Identifies the AWS region where the event originated.
	AWSRegion string `mapstructure:"region"`

	// This mapstructure array contains ARNs that identify resources that are
	// involved in the event. Inclusion of these ARNs is at the discretion
	// of the service. For example, Amazon EC2 instance state-changes
	// include Amazon EC2 instance ARNs, Auto Scaling events include ARNs
	// for both instances and Auto Scaling groups, but API calls with AWS
	// CloudTrail do not include resource ARNs.
	Resources []string
}

// Code that switches on the concrete type of events listed in this file can
// cast unknown types to an interface that implements this method in order to
// access the generic AWS event data.
//
//    ev interface{}
//    info := ev.(interface{GenericAWSEventInfo()AWSEvent}).GenericAWSEventInfo()
//
// This only works if the concrete types embed the `AWSEvent` type
func (ev AWSEvent) GenericAWSEventInfo() AWSEvent {
	return ev
}

// https://docs.aws.amazon.com/eventbridge/latest/userguide/event-types.html#events-for-services-not-listed
type CloudTrailAPICall struct {
	AWSEvent `mapstructure:",squash"`

	Detail CloudTrailAPICallDetail
}

type CloudTrailAPICallDetail struct {
	EventVersion string                 `mapstructure:"version"`
	UserIdentity map[string]interface{} `mapstructure:"userIdentity"`
	EventTime    time.Time              `mapstructure:"eventTime"`
	EventSource  string                 `mapstructure:"eventSource"`

	EventName string `mapstructure:"eventName"`
	AWSRegion string `mapstructure:"awsRegion"`
	// Some events contain values that _may_ be ip addresses, or may be the name
	// of an AWS service.
	//
	// see "sourceIPAddress" documentation in https://docs.aws.amazon.com/awscloudtrail/latest/userguide/cloudtrail-event-reference-record-contents.html
	SourceIPAddress string `mapstructure:"sourceIPAddress"`
	UserAgent       string `mapstructure:"userAgent"`

	RequestParameters map[string]interface{} `mapstructure:"requestParameters"`

	ResponseElements map[string]interface{} `mapstructure:"responseElements"`

	RequestID string `mapstructure:"requestID"`
	EventID   string `mapstructure:"eventID"`
	EventType string `mapstructure:"eventType"`
}

// https://docs.aws.amazon.com/awscloudtrail/latest/userguide/non-api-aws-service-events.html
type CloudTrailServiceEvent struct {
	AWSEvent `mapstructure:",squash"`

	Detail CloudTrailServiceEventDetail
}

type CloudTrailServiceEventDetail struct {
	EventSource         string                 `mapstructure:"eventSource"`
	ReadOnly            bool                   `mapstructure:"readOnly"`
	UserIdentity        map[string]interface{} `mapstructure:"userIdentity"`
	AWSRegion           string                 `mapstructure:"awsRegion"`
	UserAgent           string                 `mapstructure:"userAgent"`
	EventCategory       string                 `mapstructure:"eventCategory"`
	EventVersion        string                 `mapstructure:"eventVersion"`
	EventID             string                 `mapstructure:"eventID"`
	EventTime           time.Time              `mapstructure:"eventTime"`
	EventType           string                 `mapstructure:"eventType"`
	ManagementEvent     bool                   `mapstructure:"managementEvent"`
	ServiceEventDetails map[string]interface{} `mapstructure:"serviceEventDetails"`
	EventName           string                 `mapstructure:"eventName"`
	// Some events contain values that _may_ be ip addresses, or may be the name
	// of an AWS service.
	//
	// see "sourceIPAddress" documentation in https://docs.aws.amazon.com/awscloudtrail/latest/userguide/cloudtrail-event-reference-record-contents.html
	SourceIPAddress string `mapstructure:"sourceIPAddress"`
}

// https://docs.aws.amazon.com/autoscaling/ec2/userguide/cloud-watch-events.html#launch-unsuccessful
type EC2InstanceLaunchUnsuccessful struct {
	AWSEvent `mapstructure:",squash"`

	Detail EC2InstanceLaunchUnsuccessfulDetail
}

type EC2InstanceLaunchUnsuccessfulDetail struct {
	StatusCode           string
	AutoScalingGroupName string
	ActivityID           string `mapstructure:"ActivityId"`
	Details              map[string]string
	RequestID            string `mapstructure:"requestId"`
	StatusMessage        string
	EndTime              time.Time
	EC2InstanceID        string `mapstructure:"EC2InstanceId"`
	StartTime            time.Time
	Cause                string
}

// https://docs.aws.amazon.com/autoscaling/ec2/userguide/cloud-watch-events.html#launch-successful
type EC2InstanceLaunchSuccessful struct {
	AWSEvent `mapstructure:",squash"`

	Detail EC2InstanceLaunchSuccessfulDetail
}

type EC2InstanceLaunchSuccessfulDetail struct {
	StatusCode           string
	Description          string
	AutoScalingGroupName string
	ActivityID           string `mapstructure:"ActivityId"`
	Details              map[string]string
	RequestID            string `mapstructure:"requestId"`
	StatusMessage        string
	EndTime              time.Time
	EC2InstanceID        string `mapstructure:"EC2InstanceId"`
	StartTime            time.Time
	Cause                string
}

// https://docs.aws.amazon.com/eventbridge/latest/userguide/event-types.html#ec2-event-type
type EC2InstanceStateChangeNotification struct {
	AWSEvent `mapstructure:",squash"`

	Detail EC2InstanceStateChangeNotificationDetail
}

type EC2InstanceStateChangeNotificationDetail struct {
	InstanceID string
	State      string
}

// https://docs.aws.amazon.com/autoscaling/ec2/userguide/cloud-watch-events.html#launch-lifecycle-action
type EC2InstanceLaunchLifecycleAction struct {
	AWSEvent `mapstructure:",squash"`

	Detail EC2InstanceLaunchLifecycleActionDetail
}

type EC2InstanceLaunchLifecycleActionDetail struct {
	LifecycleActionToken string
	AutoScalingGroupName string
	LicecycleHookName    string
	EC2InstanceID        string `mapstructure:"EC2InstanceId"`
	LifecycleTransition  string
	NotificationMetadata string
}

// https://docs.aws.amazon.com/autoscaling/ec2/userguide/cloud-watch-events.html#terminate-lifecycle-action
type EC2InstanceTerminateSuccessful struct {
	AWSEvent `mapstructure:",squash"`

	Detail EC2InstanceTerminateSuccessfulDetail
}

type EC2InstanceTerminateSuccessfulDetail struct {
	StatusCode           string
	Description          string
	AutoScalingGroupName string
	ActivityID           string `mapstructure:"ActivityId"`
	Details              map[string]string
	RequestID            string `mapstructure:"RequestId"`
	StatusMessage        string
	EndTime              time.Time
	EC2InstanceID        string `mapstructure:"EC2InstanceId"`
	StartTime            time.Time
	Cause                string
}

// https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ebs-cloud-watch-events.html
type EBSVolumeEventNotification struct {
	AWSEvent `mapstructure:",squash"`

	Detail EBSVolumeEventNotificationDetail
}

type EBSVolumeEventNotificationDetail struct {
	Result    string `mapstructure:"result"`
	Cause     string `mapstructure:"cause"`
	Event     string `mapstructure:"event"`
	RequestID string `mapstructure:"request-id"`
}

// This is undocumented in the AWS docs!!!
type EC2SpotInstanceRequestFulfillment struct {
	AWSEvent `mapstructure:",squash"`

	Detail EC2SpotInstanceRequestFulfillmentDetail
}

type EC2SpotInstanceRequestFulfillmentDetail struct {
	Description           string `mapstructure:"description"`
	SpotInstanceRequestID string `mapstructure:"spot-instance-request-id"`
	InstanceID            string `mapstructure:"instance-id"`
}

// https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/spot-interruptions.html#spot-instance-termination-notices
type EC2SpotInstanceInterruptionWarning struct {
	AWSEvent `mapstructure:",squash"`

	Detail EC2SpotInstanceInterruptionWarningDetail
}

type EC2SpotInstanceInterruptionWarningDetail struct {
	InstanceID     string `mapstructure:"instance-id"`
	InstanceAction string `mapstructure:"instance-action"`
}

// https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/rebalance-recommendations.html#monitor-rebalance-recommendations
type EC2InstanceRebalanceRecommendation struct {
	AWSEvent `mapstructure:",squash"`

	Detail EC2InstanceRebalanceRecommendationDetail
}

type EC2InstanceRebalanceRecommendationDetail struct {
	InstanceID string `mapstructure:"instance-id"`
}
