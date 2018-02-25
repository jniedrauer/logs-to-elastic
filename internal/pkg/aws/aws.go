/*
Cache AWS services between Lambda invocations
*/
package aws

import (
	"os"

	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

var defaultRegion string = "us-east-1"
var sess *session.Session
var sessErr error
var once sync.Once

func GetSession() (*session.Session, error) {
	once.Do(func() {
		region := getRegion()
		sess, sessErr = getNewSession(region)
	})

	return sess, sessErr
}

func getNewSession(region string) (*session.Session, error) {
	return session.NewSession(&aws.Config{
		Region: aws.String(region)},
	)
}

func getRegion() string {
	region, set := os.LookupEnv("AWS_REGION")
	if !set {
		region = defaultRegion
	}

	return region
}
