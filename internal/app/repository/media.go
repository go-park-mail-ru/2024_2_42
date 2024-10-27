package repository

import (
	"fmt"
	"io"
	"path/filepath"
	"pinset/configs/s3"
	"pinset/internal/app/usecase"

	"github.com/google/uuid"
	"github.com/minio/minio-go"
)

const (
	mimeImgJpegType = "image/jpeg"
	mimeImgJpgType  = "image/jpg"
	mimeImgPngType  = "image/png"
	mimeImgGifType  = "image/gif"

	mimeVidMp4Type = "video/mp4"

	mimeAudMp3Type = "audio/mpeg"
	mimeAudAacType = "audio/aac"
	mimeAudWavType = "audio/wav"
)

const (
	minioUploadFileType = "application/octet-stream"
)

func NewMediaRepository() (usecase.MediaRepository, error) {
	config := s3.NewMinioParams()
	client, err := NewMinioClient(config)
	if err != nil {
		return nil, fmt.Errorf("minio client: %w", err)
	}

	return &MediaRepositoryController{
		client:          client,
		ImageBucketName: config.ImageBucketName,
		VideoBucketName: config.VideoBucketName,
		AudioBucketName: config.AudioBucketName,
	}, nil
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

func (mrc *MediaRepositoryController) GetMedia(bucketName, objectName string) ([]byte, error) {
	// not implemented
	return []byte{}, nil
}

func (mrc *MediaRepositoryController) UploadMedia(bucketName, fileName string, media io.Reader, mediaSize int64) (string, error) {
	objectName := uuid.New().String() + filepath.Ext(fileName)
	_, err := mrc.client.PutObject(bucketName, objectName, media, mediaSize, minio.PutObjectOptions{
		ContentType: minioUploadFileType,
	})
	if err != nil {
		return "", err
	}

	return objectName, nil
}
