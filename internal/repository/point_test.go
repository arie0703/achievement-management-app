package repository

import (
	"fmt"
	"testing"
	"time"

	"achievement-management/internal/config"
	"achievement-management/internal/errors"
	"achievement-management/internal/models"
)

func TestPointRepository_GetCurrentPoints(t *testing.T) {
	testPoints := &models.CurrentPoints{
		ID:        "current",
		Point:     100,
		UpdatedAt: time.Now(),
	}

	mockRepo := &MockRepository{
		getItemFunc: func(tableName string, key map[string]interface{}, result interface{}) error {
			if points, ok := result.(*models.CurrentPoints); ok {
				*points = *testPoints
			}
			return nil
		},
	}

	config := &config.Config{CurrentPointsTable: "test-current-points"}
	repo := NewPointRepository(mockRepo, config)

	result, err := repo.GetCurrentPoints()
	if err != nil {
		t.Errorf("GetCurrentPoints failed: %v", err)
	}

	if result.Point != testPoints.Point {
		t.Errorf("Expected point %d, got %d", testPoints.Point, result.Point)
	}
}

func TestPointRepository_GetCurrentPoints_NotFound(t *testing.T) {
	mockRepo := &MockRepository{
		getItemFunc: func(tableName string, key map[string]interface{}, result interface{}) error {
			return fmt.Errorf("item not found in table test-current-points")
		},
	}

	config := &config.Config{CurrentPointsTable: "test-current-points"}
	repo := NewPointRepository(mockRepo, config)

	result, err := repo.GetCurrentPoints()
	if err != nil {
		t.Errorf("GetCurrentPoints should not fail when item not found: %v", err)
	}

	// 初回の場合は0ポイントで初期化される
	if result.Point != 0 {
		t.Errorf("Expected point 0 for initial state, got %d", result.Point)
	}
	if result.ID != "current" {
		t.Errorf("Expected ID 'current', got %s", result.ID)
	}
}

func TestPointRepository_UpdateCurrentPoints(t *testing.T) {
	mockRepo := &MockRepository{}
	config := &config.Config{CurrentPointsTable: "test-current-points"}
	repo := NewPointRepository(mockRepo, config)

	points := &models.CurrentPoints{
		Point: 150,
	}

	err := repo.UpdateCurrentPoints(points)
	if err != nil {
		t.Errorf("UpdateCurrentPoints failed: %v", err)
	}

	// IDが自動設定されることを確認
	if points.ID != "current" {
		t.Errorf("Expected ID 'current', got %s", points.ID)
	}

	// 更新日時が設定されることを確認
	if points.UpdatedAt.IsZero() {
		t.Error("UpdatedAt should be set")
	}
}

func TestPointRepository_UpdateCurrentPoints_ValidationError(t *testing.T) {
	mockRepo := &MockRepository{}
	config := &config.Config{CurrentPointsTable: "test-current-points"}
	repo := NewPointRepository(mockRepo, config)

	tests := []struct {
		name        string
		points      *models.CurrentPoints
		expectedErr string
	}{
		{
			name:        "nil points",
			points:      nil,
			expectedErr: "validation error for field 'points': points cannot be nil",
		},
		{
			name: "negative point",
			points: &models.CurrentPoints{
				Point: -10,
			},
			expectedErr: "validation error for field 'point': point cannot be negative",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.UpdateCurrentPoints(tt.points)
			if err == nil {
				t.Error("Expected validation error")
			}
			if err.Error() != tt.expectedErr {
				t.Errorf("Expected error '%s', got '%s'", tt.expectedErr, err.Error())
			}
		})
	}
}

func TestPointRepository_CreateRewardHistory(t *testing.T) {
	mockRepo := &MockRepository{}
	config := &config.Config{RewardHistoryTable: "test-reward-history"}
	repo := NewPointRepository(mockRepo, config)

	history := &models.RewardHistory{
		RewardID:    "reward-123",
		RewardTitle: "Test Reward",
		PointCost:   50,
	}

	err := repo.CreateRewardHistory(history)
	if err != nil {
		t.Errorf("CreateRewardHistory failed: %v", err)
	}

	// IDが生成されることを確認
	if history.ID == "" {
		t.Error("ID should be generated")
	}

	// 獲得日時が設定されることを確認
	if history.RedeemedAt.IsZero() {
		t.Error("RedeemedAt should be set")
	}
}

func TestPointRepository_CreateRewardHistory_ValidationError(t *testing.T) {
	mockRepo := &MockRepository{}
	config := &config.Config{RewardHistoryTable: "test-reward-history"}
	repo := NewPointRepository(mockRepo, config)

	tests := []struct {
		name        string
		history     *models.RewardHistory
		expectedErr string
	}{
		{
			name:        "nil history",
			history:     nil,
			expectedErr: "validation error for field 'history': history cannot be nil",
		},
		{
			name: "empty reward_id",
			history: &models.RewardHistory{
				RewardID:    "",
				RewardTitle: "Test",
				PointCost:   50,
			},
			expectedErr: "validation error for field 'reward_id': reward_id is required",
		},
		{
			name: "empty reward_title",
			history: &models.RewardHistory{
				RewardID:    "reward-123",
				RewardTitle: "",
				PointCost:   50,
			},
			expectedErr: "validation error for field 'reward_title': reward_title is required",
		},
		{
			name: "zero point_cost",
			history: &models.RewardHistory{
				RewardID:    "reward-123",
				RewardTitle: "Test",
				PointCost:   0,
			},
			expectedErr: "validation error for field 'point_cost': point_cost must be positive",
		},
		{
			name: "negative point_cost",
			history: &models.RewardHistory{
				RewardID:    "reward-123",
				RewardTitle: "Test",
				PointCost:   -10,
			},
			expectedErr: "validation error for field 'point_cost': point_cost must be positive",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.CreateRewardHistory(tt.history)
			if err == nil {
				t.Error("Expected validation error")
			}
			if err.Error() != tt.expectedErr {
				t.Errorf("Expected error '%s', got '%s'", tt.expectedErr, err.Error())
			}
		})
	}
}

func TestPointRepository_GetRewardHistory(t *testing.T) {
	testHistory := []*models.RewardHistory{
		{
			ID:          "history-1",
			RewardID:    "reward-1",
			RewardTitle: "Test Reward 1",
			PointCost:   50,
			RedeemedAt:  time.Now(),
		},
		{
			ID:          "history-2",
			RewardID:    "reward-2",
			RewardTitle: "Test Reward 2",
			PointCost:   100,
			RedeemedAt:  time.Now(),
		},
	}

	mockRepo := &MockRepository{
		scanFunc: func(tableName string, result interface{}) error {
			if history, ok := result.(*[]*models.RewardHistory); ok {
				*history = testHistory
			}
			return nil
		},
	}

	config := &config.Config{RewardHistoryTable: "test-reward-history"}
	repo := NewPointRepository(mockRepo, config)

	results, err := repo.GetRewardHistory()
	if err != nil {
		t.Errorf("GetRewardHistory failed: %v", err)
	}

	if len(results) != 2 {
		t.Errorf("Expected 2 history records, got %d", len(results))
	}
}

func TestPointRepository_TransactPointsAndHistory(t *testing.T) {
	mockRepo := &MockRepository{}
	config := &config.Config{
		CurrentPointsTable:  "test-current-points",
		RewardHistoryTable: "test-reward-history",
	}
	repo := NewPointRepository(mockRepo, config)

	pointsUpdate := &models.CurrentPoints{
		Point: 50,
	}

	history := &models.RewardHistory{
		RewardID:    "reward-123",
		RewardTitle: "Test Reward",
		PointCost:   50,
	}

	err := repo.TransactPointsAndHistory(pointsUpdate, history)
	if err != nil {
		t.Errorf("TransactPointsAndHistory failed: %v", err)
	}

	// IDと日時が設定されることを確認
	if pointsUpdate.ID != "current" {
		t.Errorf("Expected points ID 'current', got %s", pointsUpdate.ID)
	}
	if pointsUpdate.UpdatedAt.IsZero() {
		t.Error("UpdatedAt should be set")
	}
	if history.ID == "" {
		t.Error("History ID should be generated")
	}
	if history.RedeemedAt.IsZero() {
		t.Error("RedeemedAt should be set")
	}
}

func TestPointRepository_TransactPointsAndHistory_ValidationError(t *testing.T) {
	mockRepo := &MockRepository{}
	config := &config.Config{
		CurrentPointsTable:  "test-current-points",
		RewardHistoryTable: "test-reward-history",
	}
	repo := NewPointRepository(mockRepo, config)

	tests := []struct {
		name         string
		pointsUpdate *models.CurrentPoints
		history      *models.RewardHistory
		expectedErr  string
	}{
		{
			name:         "nil pointsUpdate",
			pointsUpdate: nil,
			history: &models.RewardHistory{
				RewardID:    "reward-123",
				RewardTitle: "Test",
				PointCost:   50,
			},
			expectedErr: "validation error for field 'pointsUpdate': pointsUpdate cannot be nil",
		},
		{
			name: "nil history",
			pointsUpdate: &models.CurrentPoints{
				Point: 50,
			},
			history:     nil,
			expectedErr: "validation error for field 'history': history cannot be nil",
		},
		{
			name: "negative points",
			pointsUpdate: &models.CurrentPoints{
				Point: -10,
			},
			history: &models.RewardHistory{
				RewardID:    "reward-123",
				RewardTitle: "Test",
				PointCost:   50,
			},
			expectedErr: "validation error for field 'point': point cannot be negative",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.TransactPointsAndHistory(tt.pointsUpdate, tt.history)
			if err == nil {
				t.Error("Expected validation error")
			}
			if err.Error() != tt.expectedErr {
				t.Errorf("Expected error '%s', got '%s'", tt.expectedErr, err.Error())
			}
		})
	}
}

func TestPointRepository_AddPoints(t *testing.T) {
	currentPoints := &models.CurrentPoints{
		ID:        "current",
		Point:     100,
		UpdatedAt: time.Now(),
	}

	mockRepo := &MockRepository{
		getItemFunc: func(tableName string, key map[string]interface{}, result interface{}) error {
			if points, ok := result.(*models.CurrentPoints); ok {
				*points = *currentPoints
			}
			return nil
		},
	}

	config := &config.Config{CurrentPointsTable: "test-current-points"}
	repo := NewPointRepository(mockRepo, config)

	err := repo.AddPoints(50)
	if err != nil {
		t.Errorf("AddPoints failed: %v", err)
	}
}

func TestPointRepository_AddPoints_ValidationError(t *testing.T) {
	mockRepo := &MockRepository{}
	config := &config.Config{CurrentPointsTable: "test-current-points"}
	repo := NewPointRepository(mockRepo, config)

	tests := []struct {
		name        string
		points      int
		expectedErr string
	}{
		{
			name:        "zero points",
			points:      0,
			expectedErr: "validation error for field 'points': points must be positive",
		},
		{
			name:        "negative points",
			points:      -10,
			expectedErr: "validation error for field 'points': points must be positive",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.AddPoints(tt.points)
			if err == nil {
				t.Error("Expected validation error")
			}
			if err.Error() != tt.expectedErr {
				t.Errorf("Expected error '%s', got '%s'", tt.expectedErr, err.Error())
			}
		})
	}
}

func TestPointRepository_SubtractPoints(t *testing.T) {
	currentPoints := &models.CurrentPoints{
		ID:        "current",
		Point:     100,
		UpdatedAt: time.Now(),
	}

	mockRepo := &MockRepository{
		getItemFunc: func(tableName string, key map[string]interface{}, result interface{}) error {
			if points, ok := result.(*models.CurrentPoints); ok {
				*points = *currentPoints
			}
			return nil
		},
	}

	config := &config.Config{CurrentPointsTable: "test-current-points"}
	repo := NewPointRepository(mockRepo, config)

	err := repo.SubtractPoints(50)
	if err != nil {
		t.Errorf("SubtractPoints failed: %v", err)
	}
}

func TestPointRepository_SubtractPoints_InsufficientPoints(t *testing.T) {
	currentPoints := &models.CurrentPoints{
		ID:        "current",
		Point:     30, // 不足
		UpdatedAt: time.Now(),
	}

	mockRepo := &MockRepository{
		getItemFunc: func(tableName string, key map[string]interface{}, result interface{}) error {
			if points, ok := result.(*models.CurrentPoints); ok {
				*points = *currentPoints
			}
			return nil
		},
	}

	config := &config.Config{CurrentPointsTable: "test-current-points"}
	repo := NewPointRepository(mockRepo, config)

	err := repo.SubtractPoints(50)
	if err != errors.ErrInsufficientPoints {
		t.Errorf("Expected ErrInsufficientPoints, got %v", err)
	}
}