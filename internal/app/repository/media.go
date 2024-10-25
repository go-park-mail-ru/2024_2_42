package repository

import (
	"fmt"
	"io"
	"os"
	"pinset/configs/s3"
	"pinset/internal/app/usecase"

	"github.com/google/uuid"
	"github.com/minio/minio-go"
)

const (
	// image
	mimeImgJpegType = "image/jpeg"
	mimeImgJpgType  = "image/jpg"
	mimeImgPngType  = "image/png"
	mimeImgGifType  = "image/gif"

	// video
	mimeVidMp4Type = "video/mp4"

	// audio
	mimeAudMp3Type = "audio/mpeg"
	mimeAudAacType = "audio/aac"
	mimeAudWavType = "audio/wav"
)

const (
	minioUploadFileType = "application/octet-stream"
)

func NewMediaRepository() usecase.MediaRepository {
	config := s3.NewMinioParams()
	client, err := NewMinioClient(config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can't create Minio client: %s\n", err.Error())
		return nil
	}

	return &MediaRepositoryController{
		client:          client,
		ImageBucketName: config.ImageBucketName,
		VideoBucketName: config.VideoBucketName,
		AudioBucketName: config.AudioBucketName,
	}
}

func (mrc *MediaRepositoryController) GetBucketNameForContentType(fileType string) string {
	switch fileType {
	case mimeImgJpegType, mimeImgJpgType, mimeImgPngType, mimeImgGifType:
		return mrc.ImageBucketName
	case mimeVidMp4Type:
		return mrc.VideoBucketName
	case mimeAudMp3Type, mimeAudAacType, mimeAudWavType:
		return mrc.AudioBucketName
	default:
		return ""
	}
}

func (mrc *MediaRepositoryController) HasCorrectContentType(fileType string) bool {
	return fileType == mimeImgJpegType ||
		fileType == mimeImgJpgType ||
		fileType == mimeImgPngType ||
		fileType == mimeImgGifType ||
		fileType == mimeVidMp4Type ||
		fileType == mimeAudMp3Type ||
		fileType == mimeAudAacType ||
		fileType == mimeAudWavType
}

func (mrc *MediaRepositoryController) GetMedia(bucketName, objectName string) error {
	object, err := mrc.client.GetObject(bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return err
	}
	defer object.Close()

	return nil
}

func (mrc *MediaRepositoryController) UploadMedia(bucketName string, media io.Reader, mediaSize int64) (string, error) {
	objectName := uuid.New().String()
	_, err := mrc.client.PutObject(bucketName, objectName, media, mediaSize, minio.PutObjectOptions{
		ContentType: minioUploadFileType,
	})
	if err != nil {
		return "", err
	}

	return objectName, nil
}
