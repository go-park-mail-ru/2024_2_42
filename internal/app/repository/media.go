package repository

import (
	"fmt"
	"log"
	"os"
	"pinset/configs/s3"
	"pinset/internal/app/usecase"

	"github.com/minio/minio-go"
	"github.com/minio/minio-go/pkg/credentials"
)

// Добавить ЛОГГЕР //

func NewMinioClient(config s3.MinioParams) *minio.Client {
	minioClient, err := minio.NewWithOptions(config.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.AccessKeyID, config.SecretAccessKey, ""),
		Secure: config.UseSSL,
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return nil
	}

	return minioClient
}

func NewMediaRepository() usecase.MediaRepository {
	config := s3.NewMinioParams()
	return &MediaRepositoryController{
		client: NewMinioClient(config),
	}
}

func (mrc *MediaRepositoryController) GetMedia(bucket, name string) error {
	object, err := mrc.client.GetObject("mybucket", "myobject", minio.GetObjectOptions{})
	if err != nil {
		return err
	}
	defer object.Close()

	return nil
	// localFile, err := os.Create("/tmp/local-file.jpg")
	// if err != nil {
	// 	return err
	// }
	// defer localFile.Close()

	// if _, err = io.Copy(localFile, object); err != nil {
	// 	return err
	// }
}

func (mrc *MediaRepositoryController) UploadMedia(bucket, location string) error {
	err := mrc.client.MakeBucket(bucket, location)
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, errBucketExists := mrc.client.BucketExists(bucket)
		if errBucketExists == nil && exists {
			log.Printf("We already own %s\n", bucket)
		} else {
			log.Fatalln(err)
		}
	} else {
		log.Printf("Successfully created %s\n", bucket)
	}

	// Upload the test file
	// Change the value of filePath if the file is in another location
	objectName := "testdata"
	filePath := "/tmp/testdata"
	contentType := "application/octet-stream"

	// Upload the test file with FPutObject
	info, err := mrc.client.FPutObject(bucket, objectName, filePath, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Successfully uploaded %s of size %d\n", objectName, info)

	return nil
}
