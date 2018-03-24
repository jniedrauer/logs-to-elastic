// S3 functions
package awsapi

import (
	"io"
	"io/ioutil"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/service/s3/s3manager/s3manageriface"
	log "github.com/sirupsen/logrus"
)

type S3Client struct {
	s3manageriface.DownloaderAPI
}

// Create a temp file and download an S3 key to it
func (c *S3Client) Get(bucket string, key string) (string, error) {
	file, err := ioutil.TempFile("", "s3download")
	defer file.Close()
	if err != nil {
		return file.Name(), err
	}

	log.Debug("downloading file: s3://", bucket, "/", key)

	_, err = c.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		})

	return file.Name(), err
}

type MockS3Iface struct {
	TestData []byte
}

func (c MockS3Iface) Download(w io.WriterAt, input *s3.GetObjectInput, options ...func(*s3manager.Downloader)) (n int64, err error) {
	w.WriteAt([]byte(c.TestData), 0)
	return 0, nil
}

func (c MockS3Iface) DownloadWithContext(ctx aws.Context, w io.WriterAt, input *s3.GetObjectInput, options ...func(*s3manager.Downloader)) (int64, error) {
	return 0, nil
}
