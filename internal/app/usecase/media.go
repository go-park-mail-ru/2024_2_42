package usecase

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"net/http"
	delivery "pinset/internal/app/delivery/http"
)

func NewMediaUsecase(repo MediaRepository) delivery.MediaUsecase {
	return &MediaUsecaseController{
		repo: repo,
	}
}

func (muc *MediaUsecaseController) GetMedia() error {
	return nil
}

func (muc *MediaUsecaseController) UploadMedia(files []*multipart.FileHeader) error {
	for _, fileHeader := range files {
		// Open the file
		file, err := fileHeader.Open()
		if err != nil {
			return err
		}
		defer file.Close()

		fmt.Printf("the uploaded file: name[%s], size[%d], header[%#v]\n", fileHeader.Filename, fileHeader.Size, fileHeader.Header)

		mediaBytes := make([]byte, fileHeader.Size)
		_, err = file.Read(mediaBytes)
		if err != nil {
			return err
		}
		r := bytes.NewReader(mediaBytes)

		// Checking the content type
		fileType := http.DetectContentType(mediaBytes)
		if !muc.repo.HasCorrectContentType(fileType) {
			return err
		}

		bucketName := muc.repo.GetBucketNameForContentType(fileType)
		_, err = muc.repo.UploadMedia(bucketName, r, fileHeader.Size)
		if err != nil {
			return err
		}
	}

	return nil
}
