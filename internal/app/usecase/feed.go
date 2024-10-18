package usecase

import (
	delivery "pinset/internal/app/delivery/http"
	"pinset/internal/app/models"
)

func NewFeedUsecase(repo FeedRepository) delivery.FeedUsecase {
	return &feedUsecaseController{
		repo: repo,
	}
}

func (fuc *feedUsecaseController) Feed() models.Feed {
	pins := fuc.repo.GetPins()
	return models.NewFeed(pins)
}
