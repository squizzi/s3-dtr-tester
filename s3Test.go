package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
)

var config struct{ AccessKeyID, SecretAccessKey string }

var s *session.Session

func init() {
	// Load Configuration
}

func main() {
	// Environment
	os.Setenv("AWS_ACCESS_KEY_ID", config.AccessKeyID)
	os.Setenv("AWS_SECRET_ACCESS_KEY", config.SecretAccessKey)
	s = session.New(&aws.Config{
		Region: aws.String(endpoints.UsEast1RegionID),
	})
	// OR
	s = session.New(&aws.Config{
		Region: aws.String(endpoints.UsEast1RegionID),
		Credentials: credentials.NewStaticCredentials(
			config.AccessKeyID,
			config.SecretAccessKey,
			"",
		),
	})
}
