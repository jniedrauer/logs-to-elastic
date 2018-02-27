/*
Cache AWS services between Lambda invocations
*/
package aws

import (
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

var defaultRegion string = "us-east-1"
var sess *session.Session
var sessErr error
var once sync.Once

func GetSession(region string) (*session.Session, error) {
	once.Do(func() {
		sess, sessErr = getNewSession(region)
	})

	return sess, sessErr
}

func getNewSession(region string) (*session.Session, error) {
	return session.NewSession(&aws.Config{
		Region: aws.String(region)},
	)
}
