package s3

import (
	"pinset/configs"
)

type MinioParams struct {
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	UseSSL          bool
	ImageBucketName string
	VideoBucketName string
	AudioBucketName string
}

func NewMinioParams() MinioParams {
	return MinioParams{
		Endpoint:        configs.LookUpStringEnvVar("MINIO_S3_ENDPOINT", ""),
		AccessKeyID:     configs.LookUpStringEnvVar("MINIO_S3_ACCESS_KEY", ""),
		SecretAccessKey: configs.LookUpStringEnvVar("MINIO_S3_SECRET_ACCESS_KEY", ""),
		UseSSL:          configs.LookUpBoolEnvVar("MINIO_S3_USE_SSL", true),
		ImageBucketName: configs.LookUpStringEnvVar("MINIO_IMG_BUCKET_NAME", ""),
		VideoBucketName: configs.LookUpStringEnvVar("MINIO_VID_BUCKET_NAME", ""),
		AudioBucketName: configs.LookUpStringEnvVar("MINIO_AUD_BUCKET_NAME", ""),
	}
}
