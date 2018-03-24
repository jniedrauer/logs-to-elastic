// Cache AWS services between Lambda invocations using a singleton model
package awsapi

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

var defaultRegion string = "us-east-1"
var Session, _ = session.NewSession(&aws.Config{
	Region: aws.String(defaultRegion)},
)
