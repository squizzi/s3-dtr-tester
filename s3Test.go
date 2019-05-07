package main

import (
	"crypto/rand"
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var awsConfig struct {
	AccessKeyID     *string
	SecretAccessKey *string
	RegionID        *string
	Endpoint        *string
	Bucket          *string
}

func main() {
	// Assign awsConfig to flags
	help := flag.Bool("-h", false, "print help menu")
	awsConfig.AccessKeyID = flag.String("a", "", "AWS Access Key")
	awsConfig.SecretAccessKey = flag.String("s", "", "AWS Secret Key")
	awsConfig.RegionID = flag.String("r", "", "AWS Region")
	awsConfig.Endpoint = flag.String("e", "", "AWS Endpoint URL")
	awsConfig.Bucket = flag.String("b", "", "AWS Bucket")
	flag.Parse()
	if *help {
		flag.PrintDefaults()
		return
	}
	if *awsConfig.AccessKeyID == "" || *awsConfig.SecretAccessKey == "" || *awsConfig.RegionID == "" || *awsConfig.Endpoint == "" || *awsConfig.Bucket == "" {
		fmt.Println("Error: All AWS config variables must be provided!")
		flag.PrintDefaults()
		return
	}

	// Validate the awsConfig.Endpoint
	if *awsConfig.Endpoint != "" {
		_, err := url.ParseRequestURI(*awsConfig.Endpoint)
		if err != nil {
			fmt.Printf("Error: The given AWS endpoint is not a valid URL: %s", err)
			return
		}
	}

	// Configure InsecureSkipVerify to ignore TLS Cert errors
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	// Create a session
	bucket := aws.String(*awsConfig.Bucket)
	key := aws.String("testobject")

	s3Config := &aws.Config{
		Credentials:      credentials.NewStaticCredentials(*awsConfig.AccessKeyID, *awsConfig.SecretAccessKey, ""),
		Endpoint:         aws.String(*awsConfig.Endpoint),
		Region:           aws.String(*awsConfig.RegionID),
		DisableSSL:       aws.Bool(true),
		S3ForcePathStyle: aws.Bool(true),
		HTTPClient:       client,
	}
	newSession, err := session.NewSession(s3Config)
	if err != nil {
		fmt.Printf("Failed to establish a new session with aws config: %s", err)
		return
	}

	s3Client := s3.New(newSession)

	// Upload a new object "testobject" with a uuid string
	uuidString, err := newUUID()
	if err != nil {
		fmt.Printf("Unable to generate UUID for testobject: %s", err)
		return
	}
	_, err = s3Client.PutObject(&s3.PutObjectInput{
		Body:   strings.NewReader(uuidString),
		Bucket: bucket,
		Key:    key,
	})
	if err != nil {
		fmt.Printf("Failed to PutObject: %s", err)
		return
	}
	fmt.Println("Successfully PutObject")
	// Cleanup our PutObject if we succeed
	fmt.Println("Cleaning up...")
	_, err = s3Client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: bucket,
		Key:    key,
	})
	if err != nil {
		fmt.Printf("Failed to cleanup object via DeleteObject: %s", err)
		return
	}
}

// newUUID generates a random UUID according to RFC 4122
func newUUID() (string, error) {
	uuid := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, uuid)
	if n != len(uuid) || err != nil {
		return "", err
	}
	// variant bits; see section 4.1.1
	uuid[8] = uuid[8]&^0xc0 | 0x80
	// version 4 (pseudo-random); see section 4.1.3
	uuid[6] = uuid[6]&^0xf0 | 0x40
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
}
