package services

import (
	"achievement-management/internal/errors"
	"achievement-management/internal/models"
	"achievement-management/internal/repository"
)

// RewardServiceImpl 報酬サービスの実装
type RewardServiceImpl struct {
	rewardRepo repository.RewardRepository
	pointRepo  repository.PointRepository
}

// NewRewardService 報酬サービスを作成
func NewRewardService(rewardRepo repository.RewardRepository, pointRepo repository.PointRepository) RewardService {
	return &RewardServiceImpl{
		rewardRepo: rewardRepo,
		pointRepo:  pointRepo,
	}
}

// Create 報酬を作成
func (s *RewardServiceImpl) Create(reward *models.Reward) error {
	if reward == nil {
		return &errors.ValidationError{Field: "reward", Message: "reward cannot be nil"}
	}

	// バリデーション
	if err := s.validateReward(reward); err != nil {
		return err
	}

	// 報酬を作成
	return s.rewardRepo.Create(reward)
}

// Update 報酬を更新
func (s *RewardServiceImpl) Update(id string, reward *models.Reward) error {
	if id == "" {
		return &errors.ValidationError{Field: "id", Message: "id is required"}
	}

	if reward == nil {
		return &errors.ValidationError{Field: "reward", Message: "reward cannot be nil"}
	}

	// バリデーション
	if err := s.validateReward(reward); err != nil {
		return err
	}

	// IDを設定
	reward.ID = id

	// 更新実行
	return s.rewardRepo.Update(reward)
}

// GetByID IDで報酬を取得
func (s *RewardServiceImpl) GetByID(id string) (*models.Reward, error) {
	if id == "" {
		return nil, &errors.ValidationError{Field: "id", Message: "id is required"}
	}

	return s.rewardRepo.GetByID(id)
}

// List すべての報酬を取得
func (s *RewardServiceImpl) List() ([]*models.Reward, error) {
	return s.rewardRepo.List()
}

// Delete 報酬を削除
func (s *RewardServiceImpl) Delete(id string) error {
	if id == "" {
		return &errors.ValidationError{Field: "id", Message: "id is required"}
	}

	return s.rewardRepo.Delete(id)
}

// Redeem 報酬を獲得（ポイント減算と履歴記録）
func (s *RewardServiceImpl) Redeem(rewardID string) error {
	if rewardID == "" {
		return &errors.ValidationError{Field: "rewardID", Message: "rewardID is required"}
	}

	// 報酬を取得
	reward, err := s.rewardRepo.GetByID(rewardID)
	if err != nil {
		return err
	}

	// 現在のポイントを取得
	currentPoints, err := s.pointRepo.GetCurrentPoints()
	if err != nil {
		return err
	}

	// ポイントが十分かチェック
	if currentPoints.Point < reward.Point {
		return &errors.BusinessLogicError{
			Operation: "Redeem",
			Reason:    "insufficient points",
		}
	}

	// ポイント減算後の値を計算
	updatedPoints := &models.CurrentPoints{
		ID:    "current",
		Point: currentPoints.Point - reward.Point,
	}

	// 報酬獲得履歴を作成
	rewardHistory := &models.RewardHistory{
		RewardID:    reward.ID,
		RewardTitle: reward.Title,
		PointCost:   reward.Point,
	}

	// トランザクションでポイント減算と履歴記録を実行
	if err := s.pointRepo.TransactPointsAndHistory(updatedPoints, rewardHistory); err != nil {
		return err
	}

	return nil
}

// validateReward 報酬のバリデーション
func (s *RewardServiceImpl) validateReward(reward *models.Reward) error {
	if reward.Title == "" {
		return &errors.ValidationError{Field: "title", Message: "title is required"}
	}

	if reward.Point <= 0 {
		return &errors.ValidationError{Field: "point", Message: "point must be positive"}
	}

	return nil
}