// Cache AWS services between Lambda invocations using a singleton model
package awsapi

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

var defaultRegion string = "us-east-1"
var sess, err = session.NewSession(&aws.Config{
	Region: aws.String(defaultRegion)},
)

// Return an AWS session using a global variable to cache connections in
// between invocations
func Session() (*session.Session, error) {
	return sess, err
}
