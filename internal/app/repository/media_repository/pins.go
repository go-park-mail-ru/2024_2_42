package mediarepository

import (
	"database/sql"
	"errors"
	"fmt"
	"pinset/internal/app/models"
	userRepository "pinset/internal/app/repository/user_repository"
	"time"

	internal_errors "pinset/internal/errors"
)

func (mrc *MediaRepositoryController) CreatePin(pin *models.Pin) error {
	// var boardID uint64
	// err := mrc.db.QueryRow(GetFirstUserBoardByID, pin.AuthorID).Scan(&boardID)
	// if err != nil {
	// 	return fmt.Errorf("psql GetFirstUserBoardByID: %w", err)
	// }

	err := mrc.db.QueryRow(CreatePin, pin.AuthorID, pin.Title, pin.Description, pin.MediaUrl, pin.RelatedLink).Scan(&pin.PinID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			pin.PinID = 0
		}
		return fmt.Errorf("psql CreatePin: %w", err)
	}

	if pin.PinID == 0 {
		return internal_errors.ErrBadPinInputData
	}

	mrc.logger.WithField("pin was succesfully created with pinID", pin.PinID).Info("createPin func")
	return nil
}

func (mrc *MediaRepositoryController) GetAllPins(userID uint64) ([]*models.Pin, error) {
	rows, err := mrc.db.Query(GetAllPins)
	if err != nil {
		return nil, fmt.Errorf("getAllPins: %w", err)
	}
	defer rows.Close()

	var pins []*models.Pin
	for rows.Next() {
		pin := &models.Pin{}

		err := rows.Scan(
			&pin.PinID,
			&pin.AuthorID,
			&pin.MediaUrl,
			&pin.Title,
			&pin.Description)
		if err != nil {
			return nil, fmt.Errorf("getAllPins rows.Next: %w", err)
		}

		authorInfo := &models.UserPin{}
		authorInfo.UserID = pin.AuthorID

		err = mrc.db.QueryRow(GetUserInfoForPin, &pin.AuthorID).
			Scan(&authorInfo.NickName, &authorInfo.AvatarUrl)
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				return nil, fmt.Errorf("getAllPins GetUserInfoForPin: %w", err)
			}
		}

		err = mrc.db.QueryRow(userRepository.GetFollowingsCount, &pin.AuthorID).
			Scan(&authorInfo.FollowingsCount)
		if err != nil {
			return nil, fmt.Errorf("getAllPins GetFollowingsCount: %w", err)
		}

		if userID != uint64(0) {
			var availableBoards []*models.Board
			availableBoards, err = mrc.GetAllBoardsByOwnerID(userID)
			if err != nil {
				return nil, fmt.Errorf("getAllPins GetAllBoardsByOwnerID: %w", err)
			}
			pin.Boards = availableBoards
		}

		pin.AuthorInfo = authorInfo

		pins = append(pins, pin)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("getAllPins rows.Err: %w", err)
	}

	return pins, nil
}

func (mrc *MediaRepositoryController) GetPinPreviewInfoByPinID(pinID uint64) (*models.Pin, error) {
	var pinPreviewInfo models.Pin

	err := mrc.db.QueryRow(GetPinPreviewInfoByPinID, pinID).Scan(&pinPreviewInfo.PinID, &pinPreviewInfo.AuthorID, &pinPreviewInfo.MediaUrl, &pinPreviewInfo.Views)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, internal_errors.ErrPinDoesntExists
		}
		return nil, fmt.Errorf("psql getPinPreviewInfoByPinID: %w", err)
	}

	return &pinPreviewInfo, nil
}

func (mrc *MediaRepositoryController) GetPinPageInfoByPinID(pinID uint64) (*models.Pin, error) {
	var pinPreviewInfo models.Pin

	err := mrc.db.QueryRow(GetPinPageInfoByPinID, pinID).Scan(&pinPreviewInfo.PinID, &pinPreviewInfo.AuthorID, &pinPreviewInfo.Title,
		&pinPreviewInfo.Description, &pinPreviewInfo.RelatedLink, &pinPreviewInfo.Geolocation, &pinPreviewInfo.CreationTime)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, internal_errors.ErrPinDoesntExists
		}
		return nil, fmt.Errorf("psql getPinPageInfoByPinID: %w", err)
	}

	return &pinPreviewInfo, nil
}

func (mrc *MediaRepositoryController) GetPinAuthorNameByUserID(userID uint64) (*models.User, error) {
	var author models.User
	err := mrc.db.QueryRow(GetPinAuthorByUserID, userID).Scan(&author.UserName, &author.AvatarUrl)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, internal_errors.ErrBadUserID
		}
		return nil, fmt.Errorf("psql getPinAuthorNameByUserID: %w", err)
	}

	return &author, nil
}

func (mrc *MediaRepositoryController) GetPinBookmarksNumberByPinID(pinID uint64) (uint64, error) {
	var bookmarksNumber uint64
	err := mrc.db.QueryRow(GetPinBookmarksNumberByPinID, pinID).Scan(&bookmarksNumber)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, internal_errors.ErrBadPinID
		}
		return 0, fmt.Errorf("psql getPinBookmarksNumberByPinID: %w", err)
	}

	return bookmarksNumber, nil
}

func (mrc *MediaRepositoryController) UpdatePinInfoByPinID(pin *models.Pin) error {
	var pinID uint64

	_, err := mrc.db.Exec(UpdatePinInfoByPinID, pin.Title, pin.Description, pin.BoardID, pin.MediaUrl, pin.RelatedLink, pin.Geolocation, pin.PinID)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return internal_errors.ErrPinDoesntExists
		}
		return fmt.Errorf("psql updatePinInfoByPinID: %w", err)
	}

	mrc.logger.WithField("updatePinInfoByPinID with pinID:", pinID).Info()
	return nil
}

func (mrc *MediaRepositoryController) UpdatePinViewsByPinID(pinID uint64) error {
	_, err := mrc.db.Exec(UpdatePinViewsByPinID, pinID)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return internal_errors.ErrPinDoesntExists
		}
		return fmt.Errorf("psql updatePinViewsByPinID: %w", err)
	}

	mrc.logger.WithField("updatePinViewsByPinID successfull with pinID:", pinID).Info()
	return nil
}

// UpdatePinUpdateTimeByPinID = `UPDATE "pin" SET update_time = $1 WHERE pin_id = $2;`
func (mrc *MediaRepositoryController) UpdatePinUpdateTimeByPinID() error {
	// not implemented
	return nil
}

func (mrc *MediaRepositoryController) DeletePinByPinID(pinID uint64) error {
	_, err := mrc.db.Exec(DeletePinByPinID, pinID)
	if err != nil {
		return internal_errors.ErrPinDoesntExists
	}

	mrc.logger.WithField("pin was succesfully deleted with pinID", pinID).Info()
	return nil
}

func (mrc *MediaRepositoryController) GetAllCommentariesByPinID(pinID uint64) ([]*models.Comment, error) {
	rows, err := mrc.db.Query(GetAllCommentariesByPinID, pinID)
	if err != nil {
		return nil, fmt.Errorf("getAllCommentariesByPinID: %w", err)
	}
	defer rows.Close()

	var commentsList []*models.Comment
	for rows.Next() {
		var commentID, pinID, authorID uint64
		var body string
		var creationTime, updateTime time.Time
		if err := rows.Scan(&commentID, &pinID, &authorID, &body, &creationTime, &updateTime); err != nil {
			return nil, fmt.Errorf("getAllCommentariesByPinID rows.Next: %w", err)
		}
		commentsList = append(commentsList, &models.Comment{
			CommentID:    commentID,
			PinID:        pinID,
			AuthorID:     authorID,
			Body:         body,
			CreationTime: creationTime,
			UpdateTime:   updateTime,
		})
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("getAllCommentariesByPinID rows.Err: %w", err)
	}

	return commentsList, nil
}

func (mrc *MediaRepositoryController) GetBookmarkOnUserPin(ownerID, pinID uint64) (uint64, error) {
	var bookmarkID uint64

	err := mrc.db.QueryRow(GetBookmarkOnUserPin, ownerID, pinID).Scan(&bookmarkID)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return 0, internal_errors.ErrBookmarkDoesntExists
		}
		return 0, fmt.Errorf("psql getBookmarkOnUserPin: %w", err)
	}

	mrc.logger.WithField("getBookmarkOnUserPin successfull", pinID).Info()
	return bookmarkID, nil
}

func (mrc *MediaRepositoryController) CreatePinBookmark(bookmark *models.Bookmark) error {
	err := mrc.db.QueryRow(CreatePinBookmark, bookmark.OwnerID, bookmark.PinID, bookmark.BookmarkTime).Scan(&bookmark.BookmarkID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			bookmark.BookmarkID = 0
		}
		return fmt.Errorf("psql createPinBookmark: %w", err)
	}

	if bookmark.BookmarkID == 0 {
		return internal_errors.ErrBadBookmarkInputData
	}

	mrc.logger.WithField("bookmark was succesfully created with bookmarkID", bookmark.BookmarkID).Info("createPin func")
	return nil
}

func (mrc *MediaRepositoryController) DeletePinBookmarkByBookmarkID(bookmarkID uint64) error {
	_, err := mrc.db.Exec(DeletePinBookmarkByBookmarkID, bookmarkID)
	if err != nil {
		return internal_errors.ErrBookmarkDoesntExists
	}

	mrc.logger.WithField("bookmark was succesfully deleted with bookmarkID", bookmarkID).Info()
	return nil
}
