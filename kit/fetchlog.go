package kit

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"io"
	"os"
)

func FetchFromS3(file io.WriterAt, s3Bucket, s3ObjKey string) error {
	sess, _ := session.NewSession(&aws.Config{
		Credentials: credentials.NewEnvCredentials(),
	})
	downloader := s3manager.NewDownloader(sess)
	_, err := downloader.Download(file, &s3.GetObjectInput{
		Bucket: aws.String(s3Bucket),
		Key:    aws.String(s3ObjKey),
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return nil
}
