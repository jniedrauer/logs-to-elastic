// S3 functions
package awsapi

import (
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	log "github.com/sirupsen/logrus"
)

func GetFromS3(file *os.File, s3Metadata *events.S3Entity, awsRegion string) error {
	s, err := GetSession(awsRegion)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	downloader := s3manager.NewDownloader(s)

	_, err = downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(s3Metadata.Bucket.Name),
			Key:    aws.String(s3Metadata.Object.Key),
		})
	if err != nil {
		return err
	}

	return nil
}
