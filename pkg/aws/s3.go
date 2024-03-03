package aws

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"io"
	"os"
	"path/filepath"
)

func UploadJsonToS3(stage string, s3Bucket, s3Path string, data interface{}) error {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("us-east-1"),
		Credentials: credentials.NewEnvCredentials(),
	})
	if err != nil {
		UploadLogToS3(stage, fmt.Sprintf("error creating AWS session: %v", err))
		return fmt.Errorf("error creating AWS session: %w", err)
	}

	// Serialization of the Markets slice in JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		UploadLogToS3(stage, fmt.Sprintf("error marshalling data to JSON: %v", err))
		return fmt.Errorf("error marshalling data to JSON: %w", err)
	}

	// Uploading JSON to S3
	s3Client := s3.New(sess)
	_, err = s3Client.PutObject(&s3.PutObjectInput{
		Bucket:               aws.String(s3Bucket),
		Key:                  aws.String(s3Path),
		ACL:                  aws.String("private"),
		Body:                 bytes.NewReader(jsonData),
		ContentLength:        aws.Int64(int64(len(jsonData))),
		ContentType:          aws.String("application/json"),
		ServerSideEncryption: aws.String("AES256"),
	})

	if err != nil {
		UploadLogToS3(stage, fmt.Sprintf("error uploading JSON to S3: %v", err))
		return fmt.Errorf("error uploading JSON to S3: %w", err)
	}

	return nil
}

func DownloadJsonFromS3(stage string, s3Bucket, s3Path string) ([]byte, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("us-east-1"),
		Credentials: credentials.NewEnvCredentials(),
	})
	if err != nil {
		UploadLogToS3(stage, fmt.Sprintf("error creating AWS session: %w", err))
		return nil, fmt.Errorf("error creating AWS session: %w", err)
	}

	s3Client := s3.New(sess)
	output, err := s3Client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(s3Bucket),
		Key:    aws.String(s3Path),
	})
	if err != nil {
		UploadLogToS3(stage, fmt.Sprintf("error downloading JSON from S3: %w", err))
		return nil, fmt.Errorf("error downloading JSON from S3: %w", err)
	}
	defer output.Body.Close()

	bodyBytes, err := io.ReadAll(output.Body)
	if err != nil {
		UploadLogToS3(stage, fmt.Sprintf("error reading S3 object body: %w", err))
		return nil, fmt.Errorf("error reading S3 object body: %w", err)
	}

	return bodyBytes, nil
}

func UploadLogToS3(stage, logString string) {
	s3Bucket := "wt--logs"
	s3Path := "parser.log"
	logString = fmt.Sprintf("%s: %s", stage, logString)

	// Creating a new AWS session
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("us-east-1"),
		Credentials: credentials.NewEnvCredentials(),
	})
	if err != nil {
		_ = fmt.Errorf("error creating AWS session: %w", err)
	}

	// Loading to S3
	s3Client := s3.New(sess)
	_, err = s3Client.PutObject(&s3.PutObjectInput{
		Bucket:               aws.String(s3Bucket),
		Key:                  aws.String(s3Path),
		ACL:                  aws.String("private"),
		Body:                 bytes.NewReader([]byte(logString)),
		ContentLength:        aws.Int64(int64(len(logString))),
		ContentType:          aws.String("text/plain"), // The content type is set as plain text
		ServerSideEncryption: aws.String("AES256"),     // Server Encryption
	})

	if err != nil {
		_ = fmt.Errorf("error uploading log to S3: %w", err)
	}

}

func UploadFileToS3(s3Bucket, filePath, s3Path string) error {

	// Creating a new AWS session
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("us-east-1"),
		Credentials: credentials.NewEnvCredentials(),
	})
	if err != nil {
		return err
	}

	// Open file
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Getting information about the file to get the size
	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}
	var size = fileInfo.Size()

	buffer := make([]byte, size)
	_, err = file.Read(buffer)
	if err != nil {
		return err
	}

	_, fileName := filepath.Split(filePath)
	s3Key := s3Path + "/" + fileName

	// Creating an S3 object
	s3Client := s3.New(sess)

	// Uploading the file
	_, err = s3Client.PutObject(&s3.PutObjectInput{
		Bucket:               aws.String(s3Bucket),
		Key:                  aws.String(s3Key),
		ACL:                  aws.String("private"),
		Body:                 bytes.NewReader(buffer),
		ContentLength:        aws.Int64(size),
		ContentType:          aws.String("application/octet-stream"),
		ContentDisposition:   aws.String("attachment"),
		ServerSideEncryption: aws.String("AES256"),
	})

	return err
}
