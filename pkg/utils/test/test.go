//go:build test_all

package testutils

import (
	"bufio"
	"bytes"
	"context"
	"crypto/tls"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/falcosecurity/dbg-go/pkg/root"
	s3utils "github.com/falcosecurity/dbg-go/pkg/utils/s3"
	"github.com/johannesboyne/gofakes3"
	"github.com/johannesboyne/gofakes3/backend/s3mem"
	json "github.com/json-iterator/go"
	"github.com/stretchr/testify/assert"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func RunTestParsingLogs(t *testing.T, runTest func() error, parsedMsg interface{}, parsingCB func() bool) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(io.Writer(&buf), nil))
	slog.SetDefault(logger)

	err := runTest()
	assert.NoError(t, err)

	scanner := bufio.NewScanner(&buf)
	for scanner.Scan() {
		err = json.Unmarshal(scanner.Bytes(), parsedMsg)
		assert.NoError(t, err)
		if parsingCB() == false {
			break
		}
	}
}

func PreCreateFolders(opts root.Options, driverVersionsToBeCreated []string) error {
	for _, driverVersion := range driverVersionsToBeCreated {
		configPath := root.BuildConfigPath(opts, driverVersion, "")
		err := os.MkdirAll(configPath, 0700)
		if err != nil {
			return err
		}
	}
	return nil
}

func S3CreateTestBucket(t *testing.T, objectKeys []string) *s3utils.Client {
	backend := s3mem.New()
	faker := gofakes3.New(backend)
	ts := httptest.NewServer(faker.Server())
	t.Cleanup(func() {
		ts.Close()
	})

	// Difference in configuring the client

	// Setup a new config
	cfg, _ := config.LoadDefaultConfig(
		context.Background(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("TESTKEY", "TESTSECRET", "TESTSESSION")),
		config.WithHTTPClient(&http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		}),
		config.WithEndpointResolverWithOptions(
			aws.EndpointResolverWithOptionsFunc(func(_, _ string, _ ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{URL: ts.URL}, nil
			}),
		),
	)

	// Create an Amazon S3 v2 client, important to use o.UsePathStyle
	// alternatively change local DNS settings, e.g., in /etc/hosts
	// to support requests to http://<bucketname>.127.0.0.1:32947/...
	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = true
	})

	// Create bucket
	_, err := client.CreateBucket(context.Background(), &s3.CreateBucketInput{
		Bucket: aws.String(s3utils.S3Bucket),
	})
	assert.NoError(t, err)
	t.Cleanup(func() {
		_, _ = client.DeleteBucket(context.Background(), &s3.DeleteBucketInput{
			Bucket: aws.String(s3utils.S3Bucket),
		})
	})

	// Create requested test keys
	for _, key := range objectKeys {
		_, err = client.PutObject(context.Background(), &s3.PutObjectInput{
			Bucket: aws.String(s3utils.S3Bucket),
			Key:    aws.String(key),
		})
		assert.NoError(t, err)
	}
	return &s3utils.Client{Client: client}
}
