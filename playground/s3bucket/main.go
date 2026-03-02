package main

import (
	"context"
	"crypto/md5"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	client := &Client{
		bucket:  os.Getenv("BUCKET_NAME"),
		timeout: 5 * time.Minute,
		region:  os.Getenv("AWS_REGION"),
	}

	conf, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(client.region))
	if err != nil {
		log.Fatal(err)
	}

	s3Client := s3.NewFromConfig(conf)
	s3Presign := s3.NewPresignClient(s3Client)

	client.client = s3Client
	client.presigned = s3Presign

	//

	var filePath, keyResult = "./file/sample_file_upload.txt", ""

	//
	{

		file, err := os.Open(filePath)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		input := &Input{
			Filepath:    filePath,
			Body:        file,
			ContentType: "",
		}

		output, err := client.Upload(context.Background(), input)
		if err != nil {
			log.Fatal(err)
		}

		keyResult = output.Key

		fmt.Println("Key: ", output.Key)
		fmt.Println("Etag: ", output.ETag)
		fmt.Println("Location: ", output.Location)
	}

	{
		output, err := client.PresignURL(context.Background(), keyResult)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Presign: ", output.Location)
	}
}

type Client struct {
	client    *s3.Client
	presigned *s3.PresignClient
	bucket    string
	region    string
	timeout   time.Duration
}

func (c *Client) Upload(ctx context.Context, input *Input) (*Output, error) {
	var err error
	key := formatKey(input.Filepath)

	key, err = sanitizeKey(key)
	if err != nil {
		return nil, err
	}

	key = formatKey(key)

	body := input.Body
	var contentMD5 *string

	if seeker, ok := body.(io.ReadSeeker); ok {
		hash := md5.New()
		if _, err := io.Copy(hash, seeker); err != nil {
			return nil, err
		}

		md5Sum := base64.StdEncoding.EncodeToString(hash.Sum(nil))
		contentMD5 = &md5Sum

		if _, err := seeker.Seek(0, io.SeekStart); err != nil {
			return nil, err
		}
	}

	uploader := manager.NewUploader(c.client, func(u *manager.Uploader) {
		u.PartSize = 10 * 1024 * 1024
		u.Concurrency = 5
	})

	output, err := uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket:               aws.String(c.bucket),
		Key:                  aws.String(key),
		Body:                 body,
		ContentType:          aws.String(input.ContentType),
		ServerSideEncryption: types.ServerSideEncryptionAes256,
		ContentMD5:           contentMD5,
	})
	if err != nil {
		return nil, err
	}

	return &Output{
		Location: output.Location,
		Key:      key,
		ETag:     *output.ETag,
	}, nil
}

func (c *Client) PresignURL(ctx context.Context, key string) (*Output, error) {
	req, err := c.presigned.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(key),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = c.timeout
	})
	if err != nil {
		return nil, err
	}

	return &Output{
		Location: req.URL,
		Key:      key,
	}, nil
}

type Output struct {
	ETag     string
	Location string
	Key      string
}

type Input struct {
	Filepath    string
	Body        io.Reader
	ContentType string
}

func formatKey(filePath string) string {
	fileName := filePath[strings.LastIndex(filePath, "/")+1:]
	return fmt.Sprintf("%d-%s", time.Now().Unix(), fileName)
}

func sanitizeKey(key string) (string, error) {
	if key == "" ||
		strings.Contains(key, "\u0000") ||
		strings.Contains(key, "..") ||
		strings.HasPrefix(key, "/") ||
		strings.HasPrefix(key, "\\") {
		return "", errors.New("invalid key")
	}

	cleaned := path.Clean(key)

	if cleaned == "." || cleaned == ".." {
		return "", errors.New("invalid key")
	}

	normalized := strings.ReplaceAll(cleaned, "\\", "/")

	return normalized, nil
}
