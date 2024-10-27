package s3

import "pinset/configs"

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
		Endpoint:        configs.LookUpStringEnvVar("MINIO_S3_ENDPOINT", "localhost:9000"),
		AccessKeyID:     configs.LookUpStringEnvVar("MINIO_S3_ACCESS_KEY", "minioadmin"),
		SecretAccessKey: configs.LookUpStringEnvVar("MINIO_S3_SECRET_ACCESS_KEY", "minioadmin"),
		UseSSL:          configs.LookUpBoolEnvVar("MINIO_S3_USE_SSL", false),
		ImageBucketName: configs.LookUpStringEnvVar("MINIO_IMG_BUCKET_NAME", "images"),
		VideoBucketName: configs.LookUpStringEnvVar("MINIO_VID_BUCKET_NAME", "videos"),
		AudioBucketName: configs.LookUpStringEnvVar("MINIO_AUD_BUCKET_NAME", "audios"),
	}
}
