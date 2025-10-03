package repository

import (
	"fmt"
	"time"

	"achievement-management/internal/config"
	"achievement-management/internal/errors"
	"achievement-management/internal/models"

	"github.com/oklog/ulid/v2"
)

// AchievementRepositoryImpl 達成目録リポジトリの実装
type AchievementRepositoryImpl struct {
	repo   Repository
	config *config.Config
}

// NewAchievementRepository 達成目録リポジトリを作成
func NewAchievementRepository(repo Repository, config *config.Config) AchievementRepository {
	return &AchievementRepositoryImpl{
		repo:   repo,
		config: config,
	}
}

// Create 達成目録を作成
func (r *AchievementRepositoryImpl) Create(achievement *models.Achievement) error {
	if achievement == nil {
		return &errors.ValidationError{Field: "achievement", Message: "achievement cannot be nil"}
	}

	// バリデーション
	if err := r.validateAchievement(achievement); err != nil {
		return err
	}

	// IDが空の場合はULIDを生成
	if achievement.ID == "" {
		achievement.ID = ulid.Make().String()
	}

	// 作成日時を設定
	if achievement.CreatedAt.IsZero() {
		achievement.CreatedAt = time.Now()
	}

	err := r.repo.PutItem(r.config.Tables.Achievements, achievement)
	if err != nil {
		return &errors.DatabaseError{
			Operation: "Create",
			Table:     r.config.Tables.Achievements,
			Cause:     err,
		}
	}

	return nil
}

// Update 達成目録を更新
func (r *AchievementRepositoryImpl) Update(achievement *models.Achievement) error {
	if achievement == nil {
		return &errors.ValidationError{Field: "achievement", Message: "achievement cannot be nil"}
	}

	if achievement.ID == "" {
		return &errors.ValidationError{Field: "id", Message: "id is required for update"}
	}

	// バリデーション
	if err := r.validateAchievement(achievement); err != nil {
		return err
	}

	// 既存のアイテムが存在するかチェック
	existing, err := r.GetByID(achievement.ID)
	if err != nil {
		return err
	}

	// 作成日時は元の値を保持
	achievement.CreatedAt = existing.CreatedAt

	err = r.repo.PutItem(r.config.Tables.Achievements, achievement)
	if err != nil {
		return &errors.DatabaseError{
			Operation: "Update",
			Table:     r.config.Tables.Achievements,
			Cause:     err,
		}
	}

	return nil
}

// GetByID IDで達成目録を取得
func (r *AchievementRepositoryImpl) GetByID(id string) (*models.Achievement, error) {
	if id == "" {
		return nil, &errors.ValidationError{Field: "id", Message: "id is required"}
	}

	key := map[string]interface{}{
		"id": id,
	}

	var achievement models.Achievement
	err := r.repo.GetItem(r.config.Tables.Achievements, key, &achievement)
	if err != nil {
		if err.Error() == fmt.Sprintf("item not found in table %s", r.config.Tables.Achievements) {
			return nil, errors.ErrNotFound
		}
		return nil, &errors.DatabaseError{
			Operation: "GetByID",
			Table:     r.config.Tables.Achievements,
			Cause:     err,
		}
	}

	return &achievement, nil
}

// List すべての達成目録を取得
func (r *AchievementRepositoryImpl) List() ([]*models.Achievement, error) {
	var achievements []*models.Achievement
	err := r.repo.Scan(r.config.Tables.Achievements, &achievements)
	if err != nil {
		return nil, &errors.DatabaseError{
			Operation: "List",
			Table:     r.config.Tables.Achievements,
			Cause:     err,
		}
	}

	return achievements, nil
}

// Delete 達成目録を削除
func (r *AchievementRepositoryImpl) Delete(id string) error {
	if id == "" {
		return &errors.ValidationError{Field: "id", Message: "id is required"}
	}

	// 存在確認
	_, err := r.GetByID(id)
	if err != nil {
		return err
	}

	key := map[string]interface{}{
		"id": id,
	}

	err = r.repo.DeleteItem(r.config.Tables.Achievements, key)
	if err != nil {
		return &errors.DatabaseError{
			Operation: "Delete",
			Table:     r.config.Tables.Achievements,
			Cause:     err,
		}
	}

	return nil
}

// validateAchievement 達成目録のバリデーション
func (r *AchievementRepositoryImpl) validateAchievement(achievement *models.Achievement) error {
	if achievement.Title == "" {
		return &errors.ValidationError{Field: "title", Message: "title is required"}
	}

	if achievement.Point <= 0 {
		return &errors.ValidationError{Field: "point", Message: "point must be positive"}
	}

	return nil
}