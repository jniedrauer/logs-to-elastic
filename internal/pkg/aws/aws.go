/*
Cache AWS services between Lambda invocations
*/
package aws

import (
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/jniedrauer/logs-to-elastic/internal/pkg/env"
)

var defaultRegion string = "us-east-1"
var sess *session.Session
var sessErr error
var once sync.Once

func GetSession() (*session.Session, error) {
	once.Do(func() {
		region := env.GetEnvOrDefault("AWS_REGION", defaultRegion)
		sess, sessErr = getNewSession(region)
	})

	return sess, sessErr
}

func getNewSession(region string) (*session.Session, error) {
	return session.NewSession(&aws.Config{
		Region: aws.String(region)},
	)
}
