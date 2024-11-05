package mediarepository

import (
	"database/sql"
	"fmt"
	"io"
	"path/filepath"
	"pinset/configs/s3"
	"pinset/internal/app/usecase"

	"github.com/google/uuid"
	"github.com/minio/minio-go"
	"github.com/sirupsen/logrus"
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

type MediaRepositoryController struct {
	db              *sql.DB
	logger          *logrus.Logger
	client          *minio.Client
	ImageBucketName string
	VideoBucketName string
	AudioBucketName string
}

func NewMediaRepository(db *sql.DB, logger *logrus.Logger) (usecase.MediaRepository, error) {
	config := s3.NewMinioParams()
	client, err := NewMinioClient(config)
	if err != nil {
		return nil, fmt.Errorf("minio client: %w", err)
	}

	return &MediaRepositoryController{
		db: db,
		logger: logger,
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

func (mrc *MediaRepositoryController) UploadMedia(bucketName, fileName string, media io.Reader, mediaSize int64) (string, error) {
	objectName := uuid.New().String() + filepath.Ext(fileName)
	_, err := mrc.client.PutObject(bucketName, objectName, media, mediaSize, minio.PutObjectOptions{
		ContentType: minioUploadFileType,
	})
	if err != nil {
		return "", err
	}

	return mrc.GeneratePublicMediaUrl(bucketName, objectName), nil
}

func (mrc *MediaRepositoryController) GeneratePublicMediaUrl(bucketName, objectName string) string {
	config := s3.NewMinioParams()
	publicUrl := "http://" + config.Endpoint + "/" + bucketName + "/" + objectName
	return publicUrl
}
