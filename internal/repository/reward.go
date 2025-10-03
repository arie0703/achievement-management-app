package repository

import (
	"fmt"
	"time"

	"achievement-management/internal/config"
	"achievement-management/internal/errors"
	"achievement-management/internal/models"

	"github.com/oklog/ulid/v2"
)

// RewardRepositoryImpl 報酬リポジトリの実装
type RewardRepositoryImpl struct {
	repo   Repository
	config *config.Config
}

// NewRewardRepository 報酬リポジトリを作成
func NewRewardRepository(repo Repository, config *config.Config) RewardRepository {
	return &RewardRepositoryImpl{
		repo:   repo,
		config: config,
	}
}

// Create 報酬を作成
func (r *RewardRepositoryImpl) Create(reward *models.Reward) error {
	if reward == nil {
		return &errors.ValidationError{Field: "reward", Message: "reward cannot be nil"}
	}

	// バリデーション
	if err := r.validateReward(reward); err != nil {
		return err
	}

	// IDが空の場合はULIDを生成
	if reward.ID == "" {
		reward.ID = ulid.Make().String()
	}

	// 作成日時を設定
	if reward.CreatedAt.IsZero() {
		reward.CreatedAt = time.Now()
	}

	err := r.repo.PutItem(r.config.Tables.Rewards, reward)
	if err != nil {
		return &errors.DatabaseError{
			Operation: "Create",
			Table:     r.config.Tables.Rewards,
			Cause:     err,
		}
	}

	return nil
}

// Update 報酬を更新
func (r *RewardRepositoryImpl) Update(reward *models.Reward) error {
	if reward == nil {
		return &errors.ValidationError{Field: "reward", Message: "reward cannot be nil"}
	}

	if reward.ID == "" {
		return &errors.ValidationError{Field: "id", Message: "id is required for update"}
	}

	// バリデーション
	if err := r.validateReward(reward); err != nil {
		return err
	}

	// 既存のアイテムが存在するかチェック
	existing, err := r.GetByID(reward.ID)
	if err != nil {
		return err
	}

	// 作成日時は元の値を保持
	reward.CreatedAt = existing.CreatedAt

	err = r.repo.PutItem(r.config.Tables.Rewards, reward)
	if err != nil {
		return &errors.DatabaseError{
			Operation: "Update",
			Table:     r.config.Tables.Rewards,
			Cause:     err,
		}
	}

	return nil
}

// GetByID IDで報酬を取得
func (r *RewardRepositoryImpl) GetByID(id string) (*models.Reward, error) {
	if id == "" {
		return nil, &errors.ValidationError{Field: "id", Message: "id is required"}
	}

	key := map[string]interface{}{
		"id": id,
	}

	var reward models.Reward
	err := r.repo.GetItem(r.config.Tables.Rewards, key, &reward)
	if err != nil {
		if err.Error() == fmt.Sprintf("item not found in table %s", r.config.Tables.Rewards) {
			return nil, errors.ErrNotFound
		}
		return nil, &errors.DatabaseError{
			Operation: "GetByID",
			Table:     r.config.Tables.Rewards,
			Cause:     err,
		}
	}

	return &reward, nil
}

// List すべての報酬を取得
func (r *RewardRepositoryImpl) List() ([]*models.Reward, error) {
	var rewards []*models.Reward
	err := r.repo.Scan(r.config.Tables.Rewards, &rewards)
	if err != nil {
		return nil, &errors.DatabaseError{
			Operation: "List",
			Table:     r.config.Tables.Rewards,
			Cause:     err,
		}
	}

	return rewards, nil
}

// Delete 報酬を削除
func (r *RewardRepositoryImpl) Delete(id string) error {
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

	err = r.repo.DeleteItem(r.config.Tables.Rewards, key)
	if err != nil {
		return &errors.DatabaseError{
			Operation: "Delete",
			Table:     r.config.Tables.Rewards,
			Cause:     err,
		}
	}

	return nil
}

// validateReward 報酬のバリデーション
func (r *RewardRepositoryImpl) validateReward(reward *models.Reward) error {
	if reward.Title == "" {
		return &errors.ValidationError{Field: "title", Message: "title is required"}
	}

	if reward.Point <= 0 {
		return &errors.ValidationError{Field: "point", Message: "point must be positive"}
	}

	return nil
}