package mediarepository

import (
	"pinset/configs/s3"

	"github.com/minio/minio-go"
	"github.com/minio/minio-go/pkg/credentials"
)

const (
	minioBucketLocation = "eu-central-1"
)

func NewMinioClient(config s3.MinioParams) (*minio.Client, error) {
	minioClient, err := minio.NewWithOptions(config.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.AccessKeyID, config.SecretAccessKey, ""),
		Secure: config.UseSSL,
	})
	if err != nil {
		return nil, err
	}

	buckets := []string{config.ImageBucketName, config.VideoBucketName, config.AudioBucketName}
	for _, bucketName := range buckets {
		err = NewMinioBucket(minioClient, bucketName)
		if err != nil {
			return nil, err
		}
	}

	return minioClient, nil
}

func NewMinioBucket(minioClient *minio.Client, bucketName string) error {
	exists, _ := minioClient.BucketExists(bucketName)

	// If exists then just return
	if exists {
		return nil
	}

	err := minioClient.MakeBucket(bucketName, minioBucketLocation)
	if err != nil {
		return err
	}

	return nil
}
