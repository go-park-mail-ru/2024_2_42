package mediarepository

import (
	"database/sql"
	"fmt"
	"io"
	"path/filepath"
	"pinset/configs/s3"
	"pinset/internal/app/models"
	"pinset/internal/app/usecase"
	"pinset/mailer-service/mailer"

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
	mailerManager   mailer.ChatServiceClient
	ImageBucketName string
	VideoBucketName string
	AudioBucketName string
}

func NewMediaRepository(db *sql.DB, logger *logrus.Logger, manager mailer.ChatServiceClient) (usecase.MediaRepository, error) {
	config := s3.NewMinioParams()
	client, err := NewMinioClient(config)
	if err != nil {
		return nil, fmt.Errorf("minio client: %w", err)
	}

	logger.Info("MinioRepo created succesful!")
	return &MediaRepositoryController{
		db:              db,
		logger:          logger,
		client:          client,
		mailerManager:   manager,
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
	publicUrl := "http://" + s3.MinioEndPoint + "/" + bucketName + "/" + objectName
	return publicUrl
}

func (mrc *MediaRepositoryController) SetMark(markReq *models.Mark) error {
	var createdMarkID uint64
	err := mrc.db.QueryRow(`INSERT INTO mark (user_id, survey_id, question_id, score) VALUES ($1, $2, $3, $4) 
	RETURNING mark_id`, markReq.UserID, markReq.SurveyID, markReq.QuestionID, markReq.Score).Scan(&createdMarkID)

	if err != nil {
		return fmt.Errorf("psql SetMark %w", err)
	}

	mrc.logger.WithField("mark was successfully set with ID", createdMarkID).Info("setmark func")
	return nil
}

func (mrc *MediaRepositoryController) GetRandomSurvey() (*models.Survey, error) {
	surveyAnswer := &models.Survey{}
	err := mrc.db.QueryRow(`SELECT survey_id, title FROM survey ORDER BY RANDOM() LIMIT 1`).Scan(
		&surveyAnswer.SurveyID,
		&surveyAnswer.Title,
	)
	if err != nil {
		return nil, fmt.Errorf("psql GetRandomSurvey %w", err)
	}
	return surveyAnswer, nil
}

func (mrc *MediaRepositoryController) GetSurvey(surveyID uint64) (*models.Survey, error) {
	surveyAnswer := &models.Survey{}
	err := mrc.db.QueryRow(`SELECT survey_id, title FROM survey WHERE survey_id=$1`, surveyID).Scan(
		&surveyAnswer.SurveyID,
		&surveyAnswer.Title,
	)
	if err != nil {
		return nil, fmt.Errorf("psql GetSurvey %w", err)
	}

	return surveyAnswer, nil
}

func (mrc *MediaRepositoryController) GetSurveyQuestions(surveyID uint64) ([]*models.Question, error) {
	rows, err := mrc.db.Query(`SELECT question_id, content FROM question WHERE survey_id=$1`, surveyID)
	if err != nil {
		return nil, fmt.Errorf("getSurveyQuestions: %w", err)
	}
	defer rows.Close()

	var questionList []*models.Question
	for rows.Next() {
		question := &models.Question{}
		if err := rows.Scan(&question.QuestionID,
			&question.Content); err != nil {
			return nil, fmt.Errorf("getSurveyQuestions rows.Next: %w", err)
		}
		questionList = append(questionList, question)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("getSurveyQuestions rows.Err: %w", err)
	}
	return questionList, nil
}
