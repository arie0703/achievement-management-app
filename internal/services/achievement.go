package services

import (
	"achievement-management/internal/errors"
	"achievement-management/internal/models"
	"achievement-management/internal/repository"
)

// AchievementServiceImpl 達成目録サービスの実装
type AchievementServiceImpl struct {
	achievementRepo repository.AchievementRepository
	pointRepo       repository.PointRepository
}

// NewAchievementService 達成目録サービスを作成
func NewAchievementService(achievementRepo repository.AchievementRepository, pointRepo repository.PointRepository) AchievementService {
	return &AchievementServiceImpl{
		achievementRepo: achievementRepo,
		pointRepo:       pointRepo,
	}
}

// Create 達成目録を作成し、ポイントを自動加算
func (s *AchievementServiceImpl) Create(achievement *models.Achievement) error {
	if achievement == nil {
		return &errors.ValidationError{Field: "achievement", Message: "achievement cannot be nil"}
	}

	// バリデーション
	if err := s.validateAchievement(achievement); err != nil {
		return err
	}

	// 達成目録を作成
	if err := s.achievementRepo.Create(achievement); err != nil {
		return err
	}

	// ポイントを自動加算
	if err := s.pointRepo.AddPoints(achievement.Point); err != nil {
		// ポイント加算に失敗した場合、作成した達成目録を削除してロールバック
		if deleteErr := s.achievementRepo.Delete(achievement.ID); deleteErr != nil {
			// ロールバックも失敗した場合は、両方のエラーを含む複合エラーを返す
			return &errors.DatabaseError{
				Operation: "Create",
				Table:     "achievements and current_points",
				Cause:     err,
			}
		}
		return err
	}

	return nil
}

// Update 達成目録を更新
func (s *AchievementServiceImpl) Update(id string, achievement *models.Achievement) error {
	if id == "" {
		return &errors.ValidationError{Field: "id", Message: "id is required"}
	}

	if achievement == nil {
		return &errors.ValidationError{Field: "achievement", Message: "achievement cannot be nil"}
	}

	// バリデーション
	if err := s.validateAchievement(achievement); err != nil {
		return err
	}

	// IDを設定
	achievement.ID = id

	// 更新実行
	return s.achievementRepo.Update(achievement)
}

// GetByID IDで達成目録を取得
func (s *AchievementServiceImpl) GetByID(id string) (*models.Achievement, error) {
	if id == "" {
		return nil, &errors.ValidationError{Field: "id", Message: "id is required"}
	}

	return s.achievementRepo.GetByID(id)
}

// List すべての達成目録を取得
func (s *AchievementServiceImpl) List() ([]*models.Achievement, error) {
	return s.achievementRepo.List()
}

// Delete 達成目録を削除
func (s *AchievementServiceImpl) Delete(id string) error {
	if id == "" {
		return &errors.ValidationError{Field: "id", Message: "id is required"}
	}

	return s.achievementRepo.Delete(id)
}

// validateAchievement 達成目録のバリデーション
func (s *AchievementServiceImpl) validateAchievement(achievement *models.Achievement) error {
	if achievement.Title == "" {
		return &errors.ValidationError{Field: "title", Message: "title is required"}
	}

	if achievement.Point <= 0 {
		return &errors.ValidationError{Field: "point", Message: "point must be positive"}
	}

	return nil
}