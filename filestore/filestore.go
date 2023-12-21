package filestore

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"gopkg.in/ini.v1"
)

type s3operator struct {
	client *s3.Client
	bucket string
}

func newS3operator() *s3operator {
	accessKey, secretKey, err := loadFilestoreCredentials()
	region := "nl-ams"

	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL: "https://s3.nl-ams.scw.cloud",
		}, nil
	})

	// Create a new AWS configuration
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
		config.WithRegion(region),
		config.WithEndpointResolverWithOptions(customResolver),
	)

	if err != nil {
		fmt.Println("Error loading AWS config:", err)
	}

	client := s3.NewFromConfig(cfg)

	return &s3operator{
		client: client,
		bucket: "quiz3",
	}
}

func (s *s3operator) download(fileName string) {
	ctx := context.Background()
	output, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &s.bucket,
		Key:    aws.String(fileName),
	})
	if err != nil {
		log.Println("Unable to locate file ", fileName, ": ", err)
		return
	}
	defer output.Body.Close()
	file, err := os.Create("public/quiz-images/" + fileName)
	if err != nil {
		log.Printf("Couldn't create file %v. Here's why: %v\n", fileName, err)
		return
	}
	defer file.Close()
	body, err := io.ReadAll(output.Body)
	if err != nil {
		log.Printf("Couldn't read object body from %v. Here's why: %v\n", fileName, err)
		return
	}
	_, err = file.Write(body)
}

func (s *s3operator) listFiles() {
	output, err := s.client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: &s.bucket,
	})

	if err != nil {
		fmt.Println("Error listing objects:", err)
		return
	}

	fmt.Println("Objects in the bucket:")
	for _, item := range output.Contents {
		fmt.Println("Name:", *item.Key)
	}
}

func ProcessImages(images []string) {
	s3client := newS3operator()
	dir := "public/quiz-images/"
	for _, image := range images {
		if _, err := os.Stat(dir + image); err != nil {
			s3client.download(image)
		}
	}
}

func GetImageStorageUrl() string {
	cfile, err := ini.Load("conf.ini")
	if err != nil {
		log.Fatal(err)
	}
	return cfile.Section("imagestorage").Key("url").String()
}

func loadFilestoreCredentials() (string, string, error) {
	cfile, err := ini.Load("conf.ini")
	if err != nil {
		log.Fatal(err)
	}
	key := cfile.Section("imagestorage").Key("key").String()
	secret := cfile.Section("imagestorage").Key("secret").String()
	return key, secret, nil
}
