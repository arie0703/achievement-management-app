package repository

import (
	"fmt"
	"time"

	"achievement-management/internal/config"
	"achievement-management/internal/errors"
	"achievement-management/internal/models"

	"github.com/oklog/ulid/v2"
)

// PointRepositoryImpl ポイントリポジトリの実装
type PointRepositoryImpl struct {
	repo   Repository
	config *config.Config
}

// NewPointRepository ポイントリポジトリを作成
func NewPointRepository(repo Repository, config *config.Config) PointRepository {
	return &PointRepositoryImpl{
		repo:   repo,
		config: config,
	}
}

// GetCurrentPoints 現在のポイントを取得
func (r *PointRepositoryImpl) GetCurrentPoints() (*models.CurrentPoints, error) {
	key := map[string]interface{}{
		"id": "current",
	}

	var currentPoints models.CurrentPoints
	err := r.repo.GetItem(r.config.Tables.CurrentPoints, key, &currentPoints)
	if err != nil {
		if err.Error() == fmt.Sprintf("item not found in table %s", r.config.Tables.CurrentPoints) {
			// 初回の場合は0ポイントで初期化
			return &models.CurrentPoints{
				ID:        "current",
				Point:     0,
				UpdatedAt: time.Now(),
			}, nil
		}
		return nil, &errors.DatabaseError{
			Operation: "GetCurrentPoints",
			Table:     r.config.Tables.CurrentPoints,
			Cause:     err,
		}
	}

	return &currentPoints, nil
}

// UpdateCurrentPoints 現在のポイントを更新
func (r *PointRepositoryImpl) UpdateCurrentPoints(points *models.CurrentPoints) error {
	if points == nil {
		return &errors.ValidationError{Field: "points", Message: "points cannot be nil"}
	}

	// IDを固定値に設定
	points.ID = "current"
	
	// 更新日時を設定
	points.UpdatedAt = time.Now()

	// ポイントが負の値にならないようにチェック
	if points.Point < 0 {
		return &errors.ValidationError{Field: "point", Message: "point cannot be negative"}
	}

	err := r.repo.PutItem(r.config.Tables.CurrentPoints, points)
	if err != nil {
		return &errors.DatabaseError{
			Operation: "UpdateCurrentPoints",
			Table:     r.config.Tables.CurrentPoints,
			Cause:     err,
		}
	}

	return nil
}

// CreateRewardHistory 報酬獲得履歴を作成
func (r *PointRepositoryImpl) CreateRewardHistory(history *models.RewardHistory) error {
	if history == nil {
		return &errors.ValidationError{Field: "history", Message: "history cannot be nil"}
	}

	// バリデーション
	if err := r.validateRewardHistory(history); err != nil {
		return err
	}

	// IDが空の場合はULIDを生成
	if history.ID == "" {
		history.ID = ulid.Make().String()
	}

	// 獲得日時を設定
	if history.RedeemedAt.IsZero() {
		history.RedeemedAt = time.Now()
	}

	err := r.repo.PutItem(r.config.Tables.RewardHistory, history)
	if err != nil {
		return &errors.DatabaseError{
			Operation: "CreateRewardHistory",
			Table:     r.config.Tables.RewardHistory,
			Cause:     err,
		}
	}

	return nil
}

// GetRewardHistory 報酬獲得履歴を取得
func (r *PointRepositoryImpl) GetRewardHistory() ([]*models.RewardHistory, error) {
	var history []*models.RewardHistory
	err := r.repo.Scan(r.config.Tables.RewardHistory, &history)
	if err != nil {
		return nil, &errors.DatabaseError{
			Operation: "GetRewardHistory",
			Table:     r.config.Tables.RewardHistory,
			Cause:     err,
		}
	}

	return history, nil
}

// TransactPointsAndHistory ポイント更新と履歴記録をトランザクションで実行
func (r *PointRepositoryImpl) TransactPointsAndHistory(pointsUpdate *models.CurrentPoints, history *models.RewardHistory) error {
	if pointsUpdate == nil {
		return &errors.ValidationError{Field: "pointsUpdate", Message: "pointsUpdate cannot be nil"}
	}
	if history == nil {
		return &errors.ValidationError{Field: "history", Message: "history cannot be nil"}
	}

	// バリデーション
	if err := r.validateRewardHistory(history); err != nil {
		return err
	}

	// ポイントが負の値にならないようにチェック
	if pointsUpdate.Point < 0 {
		return &errors.ValidationError{Field: "point", Message: "point cannot be negative"}
	}

	// IDと日時を設定
	pointsUpdate.ID = "current"
	pointsUpdate.UpdatedAt = time.Now()

	if history.ID == "" {
		history.ID = ulid.Make().String()
	}
	if history.RedeemedAt.IsZero() {
		history.RedeemedAt = time.Now()
	}

	// トランザクションアイテムを準備
	transactItems := []TransactWriteItem{
		{
			TableName: r.config.Tables.CurrentPoints,
			Item:      pointsUpdate,
			Operation: "PUT",
		},
		{
			TableName: r.config.Tables.RewardHistory,
			Item:      history,
			Operation: "PUT",
		},
	}

	err := r.repo.TransactWrite(transactItems)
	if err != nil {
		return &errors.DatabaseError{
			Operation: "TransactPointsAndHistory",
			Table:     fmt.Sprintf("%s,%s", r.config.Tables.CurrentPoints, r.config.Tables.RewardHistory),
			Cause:     err,
		}
	}

	return nil
}

// validateRewardHistory 報酬獲得履歴のバリデーション
func (r *PointRepositoryImpl) validateRewardHistory(history *models.RewardHistory) error {
	if history.RewardID == "" {
		return &errors.ValidationError{Field: "reward_id", Message: "reward_id is required"}
	}

	if history.RewardTitle == "" {
		return &errors.ValidationError{Field: "reward_title", Message: "reward_title is required"}
	}

	if history.PointCost <= 0 {
		return &errors.ValidationError{Field: "point_cost", Message: "point_cost must be positive"}
	}

	return nil
}

// AddPoints ポイントを加算（達成目録追加時に使用）
func (r *PointRepositoryImpl) AddPoints(points int) error {
	if points <= 0 {
		return &errors.ValidationError{Field: "points", Message: "points must be positive"}
	}

	// 現在のポイントを取得
	currentPoints, err := r.GetCurrentPoints()
	if err != nil {
		return err
	}

	// ポイントを加算
	currentPoints.Point += points

	// 更新
	return r.UpdateCurrentPoints(currentPoints)
}

// SubtractPoints ポイントを減算（報酬獲得時に使用）
func (r *PointRepositoryImpl) SubtractPoints(points int) error {
	if points <= 0 {
		return &errors.ValidationError{Field: "points", Message: "points must be positive"}
	}

	// 現在のポイントを取得
	currentPoints, err := r.GetCurrentPoints()
	if err != nil {
		return err
	}

	// ポイントが不足していないかチェック
	if currentPoints.Point < points {
		return errors.ErrInsufficientPoints
	}

	// ポイントを減算
	currentPoints.Point -= points

	// 更新
	return r.UpdateCurrentPoints(currentPoints)
}