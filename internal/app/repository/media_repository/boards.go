package mediarepository

import (
	"database/sql"
	"errors"
	"fmt"
	"pinset/internal/app/models"
	internal_errors "pinset/internal/errors"
	"time"
)

func (mrc *MediaRepositoryController) GetBoardPinsByBoardID(boardID uint64) ([]uint64, error) {
	rows, err := mrc.db.Query(GetPinsIDByBoardID, boardID)
	if err != nil {
		return nil, fmt.Errorf("GetBoardPinsByBoardID: %w", err)
	}
	defer rows.Close()

	var pinIDs []uint64

	for rows.Next() {
		var pinID uint64

		if err := rows.Scan(&pinID); err != nil {
			return nil, fmt.Errorf("GetBoardPinsByBoardID rows.Next: %w", err)
		}

		pinIDs = append(pinIDs, pinID)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetBoardPinsByBoardID rows.Err: %w", err)
	}

	return pinIDs, nil
}

func (mrc *MediaRepositoryController) AddPinToBoard(boardID uint64, pinID uint64) error {
	var updatedBoardID uint64
	err := mrc.db.QueryRow(AddPinToBoard, boardID, pinID).Scan(&updatedBoardID)
	if err != nil {
		return fmt.Errorf("AddPinToBoard %w", err)
	}

	mrc.logger.WithField("Pin successfully added to board", updatedBoardID).Info("addPinToBoard func")
	return nil
}

func (mrc *MediaRepositoryController) GetAllBoardsByOwnerID(ownerID uint64) ([]*models.Board, error) {
	rows, err := mrc.db.Query(GetAllBoardsByOwnerID, ownerID)
	if err != nil {
		return nil, fmt.Errorf("getAllBoardsByOwnerID: %w", err)
	}
	defer rows.Close()

	var boards []*models.Board
	for rows.Next() {
		var boardID, ownerID uint64
		var title, description string
		var public bool
		var creationTime, updateTime time.Time

		if err := rows.Scan(&boardID, &ownerID, &title, &description, &public, &creationTime, &updateTime); err != nil {
			return nil, fmt.Errorf("getAllBoardsByOwnerID rows.Next: %w", err)
		}

		boards = append(boards, &models.Board{
			BoardID:      boardID,
			OwnerID:      ownerID,
			Name:         title,
			Description:  description,
			Public:       public,
			CreationTime: creationTime,
			UpdateTime:   updateTime,
		})
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("getAllBoardsByOwnerID rows.Err: %w", err)
	}

	return boards, nil
}

func (mrc *MediaRepositoryController) GetBoardByBoardID(boardID uint64) (*models.Board, error) {
	row := mrc.db.QueryRow(GetBoardByBoardID, boardID)
	if row.Err() != nil {
		if errors.Is(row.Err(), sql.ErrNoRows) {
			return nil, internal_errors.ErrBoardDoesntExists
		}
		return nil, fmt.Errorf("psql getBoardByBoardID: %w", row.Err())
	}

	var id, ownerID uint64
	var name, description string
	var public bool
	var creationTime, updateTime time.Time

	if err := row.Scan(&id, &ownerID, &name, &description, &public, &creationTime, &updateTime); err != nil {
		return nil, fmt.Errorf("getAllBoardsByOwnerID rows.Next: %w", err)
	}

	return &models.Board{
		BoardID:      boardID,
		OwnerID:      ownerID,
		Name:         name,
		Description:  description,
		Public:       public,
		CreationTime: creationTime,
		UpdateTime:   updateTime,
	}, nil
}

func (mrc *MediaRepositoryController) CreateBoard(board *models.Board) error {
	err := mrc.db.QueryRow(CreateBoard, board.OwnerID, board.Name, board.Description, board.Public).Scan(&board.BoardID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			board.BoardID = 0
		}
		return fmt.Errorf("psql createBoard: %w", err)
	}

	if board.BoardID == 0 {
		return internal_errors.ErrBadBoardID
	}

	mrc.logger.WithField("board was succesful created with boardID", board.BoardID).Info("createBoard func")
	return nil
}

func (mrc *MediaRepositoryController) UpdateBoardByBoardID(board *models.Board) error {
	var boardID uint64

	err := mrc.db.QueryRow(UpdateBoardByBoardID, board.Name, board.Description, board.Public).Scan(&boardID)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return internal_errors.ErrBoardDoesntExists
		}
		return fmt.Errorf("psql updateBoardByBoardID: %w", err)
	}

	mrc.logger.WithField("updateBoardByBoardID with boardID:", boardID).Info()
	return nil
}

func (mrc *MediaRepositoryController) DeleteBoardByBoardID(boardID uint64) error {
	_, err := mrc.db.Exec(DeleteBoardByBoardID, boardID)
	if err != nil {
		return internal_errors.ErrBoardDoesntExists
	}

	mrc.logger.WithField("board was succesfil deleted with boardID", boardID).Info()
	return nil
}
