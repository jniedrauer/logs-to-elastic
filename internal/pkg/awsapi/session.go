// Cache AWS services between Lambda invocations using a singleton model
package awsapi

import (
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

var defaultRegion string = "us-east-1"
var sess *session.Session
var sessErr error
var once sync.Once

// Return an AWS session, creating one if required
func GetSession(region string) (*session.Session, error) {
	once.Do(func() {
		sess, sessErr = getNewSession(region)
	})

	return sess, sessErr
}

// Create a new AWS session to region
func getNewSession(region string) (*session.Session, error) {
	return session.NewSession(&aws.Config{
		Region: aws.String(region)},
	)
}
