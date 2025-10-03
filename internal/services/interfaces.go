package services

import "achievement-management/internal/models"

// AchievementService 達成目録サービス
type AchievementService interface {
	Create(achievement *models.Achievement) error
	Update(id string, achievement *models.Achievement) error
	GetByID(id string) (*models.Achievement, error)
	List() ([]*models.Achievement, error)
	Delete(id string) error
}

// RewardService 報酬サービス
type RewardService interface {
	Create(reward *models.Reward) error
	Update(id string, reward *models.Reward) error
	GetByID(id string) (*models.Reward, error)
	List() ([]*models.Reward, error)
	Delete(id string) error
	Redeem(rewardID string) error
}

// PointService ポイントサービス
type PointService interface {
	GetCurrentPoints() (*models.CurrentPoints, error)
	AddPoints(points int) error
	SubtractPoints(points int) error
	AggregatePoints() (*models.PointSummary, error)
	GetRewardHistory() ([]*models.RewardHistory, error)
}