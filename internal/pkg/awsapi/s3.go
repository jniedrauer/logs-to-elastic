// S3 functions
package awsapi

import (
	"io/ioutil"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	log "github.com/sirupsen/logrus"
)

// Create a temp file and download an S3 key to it
func GetFromS3(s3Metadata *events.S3Entity, awsRegion string) (string, error) {
	file, err := ioutil.TempFile("", "s3logs")
	defer file.Close()
	if err != nil {
		return file.Name(), err
	}

	s, err := GetSession(awsRegion)
	if err != nil {
		return file.Name(), err
	}

	downloader := s3manager.NewDownloader(s)

	log.Debug("downloading file: s3://", s3Metadata.Bucket.Name, "/", s3Metadata.Object.Key)

	_, err = downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(s3Metadata.Bucket.Name),
			Key:    aws.String(s3Metadata.Object.Key),
		})

	return file.Name(), err
}
