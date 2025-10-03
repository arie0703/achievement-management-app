package services

import (
	"achievement-management/internal/errors"
	"achievement-management/internal/models"
	"achievement-management/internal/repository"
)

// PointServiceImpl ポイントサービスの実装
type PointServiceImpl struct {
	pointRepo       repository.PointRepository
	achievementRepo repository.AchievementRepository
}

// NewPointService ポイントサービスを作成
func NewPointService(pointRepo repository.PointRepository, achievementRepo repository.AchievementRepository) PointService {
	return &PointServiceImpl{
		pointRepo:       pointRepo,
		achievementRepo: achievementRepo,
	}
}

// GetCurrentPoints 現在のポイントを取得
func (s *PointServiceImpl) GetCurrentPoints() (*models.CurrentPoints, error) {
	return s.pointRepo.GetCurrentPoints()
}

// AddPoints ポイントを加算
func (s *PointServiceImpl) AddPoints(points int) error {
	if points <= 0 {
		return &errors.ValidationError{Field: "points", Message: "points must be positive"}
	}

	return s.pointRepo.AddPoints(points)
}

// SubtractPoints ポイントを減算
func (s *PointServiceImpl) SubtractPoints(points int) error {
	if points <= 0 {
		return &errors.ValidationError{Field: "points", Message: "points must be positive"}
	}

	return s.pointRepo.SubtractPoints(points)
}

// AggregatePoints 全達成目録のポイントを集計し、現在のポイントと比較
func (s *PointServiceImpl) AggregatePoints() (*models.PointSummary, error) {
	// 全達成目録を取得
	achievements, err := s.achievementRepo.List()
	if err != nil {
		return nil, &errors.ServiceError{
			Operation: "AggregatePoints",
			Message:   "failed to get achievements list",
			Cause:     err,
		}
	}

	// 達成目録のポイント合計を計算
	totalPoints := 0
	for _, achievement := range achievements {
		if achievement != nil {
			totalPoints += achievement.Point
		}
	}

	// 現在のポイントを取得
	currentPoints, err := s.pointRepo.GetCurrentPoints()
	if err != nil {
		return nil, &errors.ServiceError{
			Operation: "AggregatePoints",
			Message:   "failed to get current points",
			Cause:     err,
		}
	}

	// 差異を計算（達成目録の合計 - 現在のポイント）
	difference := totalPoints - currentPoints.Point

	// 集計結果を作成
	summary := &models.PointSummary{
		TotalAchievements: len(achievements),
		TotalPoints:       totalPoints,
		CurrentBalance:    currentPoints.Point,
		Difference:        difference,
	}

	return summary, nil
}

// GetRewardHistory 報酬獲得履歴を取得
func (s *PointServiceImpl) GetRewardHistory() ([]*models.RewardHistory, error) {
	return s.pointRepo.GetRewardHistory()
}