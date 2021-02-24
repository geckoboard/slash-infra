package awsutil

import (
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
)

const EnvVarPrefixForAwsRoles = "AWS_ROLE_"

type AWSCredentials struct {
	// SDKSession contains all of the config an SDK client needs to work
	// It is safe for concurrent use
	SDKSession *session.Session
	// AccountID is the numerical ID of the AWS account that owns the role the credentials assume
	AccountID string
	// AccountAlias is an underscored alias for this set of credentials
	AccountAlias string
}

// DetectAWSCredentials uses environment variables to build instances
// of the EC2 client library for each AWS account it should discover resources
// within.
//
// The main variables for configuration are:
//
// `AWS_ROLE_{account alias}` - The role slash-infra should assume to gain access
// to the account known as {account alias}
//
// `AWS_REGION_{account alias}` - If the account's resources are in a region
// other than us-east-1, specify it here.
//
// If an account uses several regions, then you can specify role several times
// under different aliases. e.g.
//
// ```
// AWS_ROLE_DEV_US_EAST=...
// AWS_ROLE_DEV_EU=...
// ```
func DetectAWSCredentials() ([]AWSCredentials, error) {
	detected := []AWSCredentials{}
	environ := os.Environ()

	for _, pair := range environ {
		if !strings.HasPrefix(pair, EnvVarPrefixForAwsRoles) {
			continue
		}
		parts := strings.SplitN(pair, "=", 2)
		key := parts[0]
		roleArn := parts[1]

		awsAccountAlias := key[len(EnvVarPrefixForAwsRoles):]

		// Some of our infra is not in us-east-1 (e.g. dev-vpc)
		// Allow slash-infra to create clients that will discover resources in those regions
		region := os.Getenv(fmt.Sprintf("AWS_REGION_%s", awsAccountAlias))
		if region == "" {
			region = "us-east-1"
		}

		// Session provides configuration for the SDK's service
		// clients. Sessions can be shared across service clients that
		// share the same base configuration.
		sess, err := session.NewSession(&aws.Config{
			Credentials: credentials.NewEnvCredentials(),
			Region:      aws.String(region),
		})
		if err != nil {
			return nil, err
		}

		role, err := arn.Parse(roleArn)
		if err != nil {
			return nil, err
		}

		creds := AWSCredentials{
			SDKSession: sess.Copy(&aws.Config{
				Credentials: stscreds.NewCredentials(sess, roleArn),
			}),
			AccountID:    role.AccountID,
			AccountAlias: awsAccountAlias,
		}
		detected = append(detected, creds)
	}

	return detected, nil
}
