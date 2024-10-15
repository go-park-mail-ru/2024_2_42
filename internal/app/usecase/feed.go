package usecase

import (
	"pinset/internal/app/repository"
	"pinset/internal/models"
)

func NewFeedUsecase(repo repository.FeedRepository) FeedUsecase {
	return &feedUsecaseController{
		repo: repo,
	}
}

func (fuc *feedUsecaseController) Feed() models.Feed {
	pins := fuc.repo.GetPins()
	return models.NewFeed(pins)
}
