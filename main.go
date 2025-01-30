package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/joho/godotenv"
	"io"
	"os"
)

const (
	AWS_BUCKET = "golearnbucket"
)

var pl = fmt.Printf

var awsConfig aws.Config
var s3Client *s3.Client

func init() {

	// env

	err := godotenv.Load()
	if err != nil {
		panic("Env load err ")
	}

	// aws config
	awsC, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ap-south-1"))
	if err != nil {
		panic("AWS Config failed")
	}
	awsConfig = awsC

	s3Client = s3.NewFromConfig(awsConfig)

}

func main() {
	file, err := os.Open(".gitignore")

	if err != nil {
		pl("%s", err)
	}
	defer file.Close()

	var buf bytes.Buffer

	if _, err := io.Copy(&buf, file); err != nil {
		pl("File copying error %s\n", err)
		return
	}

	uploadFile("Hello", buf.Bytes())
	downLoadFile("Hello")
	deleteFile("Hello")
}

func uploadFile(name string, file []byte) {
	_, err := s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(AWS_BUCKET),
		Key:    aws.String(name),
		Body:   bytes.NewReader(file),
		Metadata: map[string]string{
			"Name": name,
		},
	},
	)
	if err != nil {
		pl("Failed to Upload file %s\n", name)
		return
	}
	pl("uploaded Successfully %s\n", name)
}

func downLoadFile(name string) {
	res, err := s3Client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(AWS_BUCKET),
		Key:    aws.String(name),
	})
	if err != nil {
		pl("Error while downloading File %s %s\n", name, err)
		return
	}
	pl("Downloaded Successfully %v\n", res.Metadata["name"])
}

func deleteFile(name string) {
	_, err := s3Client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Key:    aws.String(name),
		Bucket: aws.String(AWS_BUCKET),
	})

	if err != nil {
		pl("Error while deleting file %s %s\n", name, err)
		return
	}

	pl("Deleted Successfully %s\n", name)
}
