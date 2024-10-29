package usecase

import (
	delivery "pinset/internal/app/delivery/http"
	"pinset/internal/app/models"
)

func NewFeedUsecase(repo FeedRepository) delivery.FeedUsecase {
	return &FeedUsecaseController{
		repo: repo,
	}
}

func (fuc *FeedUsecaseController) Feed() models.Feed {
	pins := fuc.repo.GetPins()
	return models.NewFeed(pins)
}
