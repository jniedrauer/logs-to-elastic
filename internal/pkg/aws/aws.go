/*
Cache AWS services between Lambda invocations
*/
package aws

import (
	"os"

	log "github.com/sirupsen/logrus"

	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

var defaultRegion string = "us-east-1"
var sess *session.Session
var err error
var once sync.Once

func GetSession() (*session.Session, error) {
	once.Do(func() {
		region := getRegion()

		sess, err = session.NewSession(&aws.Config{
			Region: aws.String(region)},
		)
		if err != nil {
			log.Fatal("Error getting session: %v", err)
			os.Exit(1)
		}
	})
	return sess, nil
}

func getRegion() string {
	region, set := os.LookupEnv("AWS_REGION")
	if !set {
		region = defaultRegion
	}

	return region
}
