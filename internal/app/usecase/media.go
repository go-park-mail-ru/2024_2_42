package usecase

import (
	"bytes"
	"mime/multipart"
	"net/http"
	delivery "pinset/internal/app/delivery/http"

	internal_errors "pinset/internal/errors"
)

func NewMediaUsecase(repo MediaRepository) delivery.MediaUsecase {
	return &MediaUsecaseController{
		repo: repo,
	}
}

func (muc *MediaUsecaseController) GetMedia() error {
	// not implemented
	return nil
}

func (muc *MediaUsecaseController) UploadMedia(files []*multipart.FileHeader) ([]string, error) {
	var uploadedMediaIds []string

	for _, fileHeader := range files {
		// Open the file
		file, err := fileHeader.Open()
		if err != nil {
			return []string{}, err
		}
		defer file.Close()

		mediaBytes := make([]byte, fileHeader.Size)
		_, err = file.Read(mediaBytes)
		if err != nil {
			return []string{}, err
		}
		r := bytes.NewReader(mediaBytes)

		// Checking the content type
		fileType := http.DetectContentType(mediaBytes)
		if !muc.repo.HasCorrectContentType(fileType) {
			return []string{}, internal_errors.ErrWrongMediaContentType
		}

		bucketName := muc.repo.GetBucketNameForContentType(fileType)
		uploadedMediaId, err := muc.repo.UploadMedia(bucketName, fileHeader.Filename, r, fileHeader.Size)
		if err != nil {
			return []string{}, err
		}

		uploadedMediaIds = append(uploadedMediaIds, uploadedMediaId)
	}

	return uploadedMediaIds, nil
}
