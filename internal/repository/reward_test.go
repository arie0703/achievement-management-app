package repository

import (
	"fmt"
	"testing"
	"time"

	"achievement-management/internal/config"
	"achievement-management/internal/errors"
	"achievement-management/internal/models"
)

func TestRewardRepository_Create(t *testing.T) {
	mockRepo := &MockRepository{}
	config := &config.Config{RewardsTable: "test-rewards"}
	repo := NewRewardRepository(mockRepo, config)

	reward := &models.Reward{
		Title:       "Test Reward",
		Description: "Test Description",
		Point:       50,
	}

	err := repo.Create(reward)
	if err != nil {
		t.Errorf("Create failed: %v", err)
	}

	// IDが生成されていることを確認
	if reward.ID == "" {
		t.Error("ID should be generated")
	}

	// 作成日時が設定されていることを確認
	if reward.CreatedAt.IsZero() {
		t.Error("CreatedAt should be set")
	}
}

func TestRewardRepository_Create_ValidationError(t *testing.T) {
	mockRepo := &MockRepository{}
	config := &config.Config{RewardsTable: "test-rewards"}
	repo := NewRewardRepository(mockRepo, config)

	tests := []struct {
		name        string
		reward      *models.Reward
		expectedErr string
	}{
		{
			name:        "nil reward",
			reward:      nil,
			expectedErr: "validation error for field 'reward': reward cannot be nil",
		},
		{
			name: "empty title",
			reward: &models.Reward{
				Title: "",
				Point: 50,
			},
			expectedErr: "validation error for field 'title': title is required",
		},
		{
			name: "zero point",
			reward: &models.Reward{
				Title: "Test",
				Point: 0,
			},
			expectedErr: "validation error for field 'point': point must be positive",
		},
		{
			name: "negative point",
			reward: &models.Reward{
				Title: "Test",
				Point: -10,
			},
			expectedErr: "validation error for field 'point': point must be positive",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Create(tt.reward)
			if err == nil {
				t.Error("Expected validation error")
			}
			if err.Error() != tt.expectedErr {
				t.Errorf("Expected error '%s', got '%s'", tt.expectedErr, err.Error())
			}
		})
	}
}

func TestRewardRepository_GetByID(t *testing.T) {
	testReward := &models.Reward{
		ID:          "test-id",
		Title:       "Test Reward",
		Description: "Test Description",
		Point:       50,
		CreatedAt:   time.Now(),
	}

	mockRepo := &MockRepository{
		getItemFunc: func(tableName string, key map[string]interface{}, result interface{}) error {
			if reward, ok := result.(*models.Reward); ok {
				*reward = *testReward
			}
			return nil
		},
	}

	config := &config.Config{RewardsTable: "test-rewards"}
	repo := NewRewardRepository(mockRepo, config)

	result, err := repo.GetByID("test-id")
	if err != nil {
		t.Errorf("GetByID failed: %v", err)
	}

	if result.ID != testReward.ID {
		t.Errorf("Expected ID %s, got %s", testReward.ID, result.ID)
	}
}

func TestRewardRepository_GetByID_NotFound(t *testing.T) {
	mockRepo := &MockRepository{
		getItemFunc: func(tableName string, key map[string]interface{}, result interface{}) error {
			return fmt.Errorf("item not found in table test-rewards")
		},
	}

	config := &config.Config{RewardsTable: "test-rewards"}
	repo := NewRewardRepository(mockRepo, config)

	_, err := repo.GetByID("non-existent-id")
	if err != errors.ErrNotFound {
		t.Errorf("Expected ErrNotFound, got %v", err)
	}
}

func TestRewardRepository_GetByID_EmptyID(t *testing.T) {
	mockRepo := &MockRepository{}
	config := &config.Config{RewardsTable: "test-rewards"}
	repo := NewRewardRepository(mockRepo, config)

	_, err := repo.GetByID("")
	if err == nil {
		t.Error("Expected validation error for empty ID")
	}
}

func TestRewardRepository_List(t *testing.T) {
	testRewards := []*models.Reward{
		{
			ID:          "test-id-1",
			Title:       "Test Reward 1",
			Description: "Test Description 1",
			Point:       50,
			CreatedAt:   time.Now(),
		},
		{
			ID:          "test-id-2",
			Title:       "Test Reward 2",
			Description: "Test Description 2",
			Point:       100,
			CreatedAt:   time.Now(),
		},
	}

	mockRepo := &MockRepository{
		scanFunc: func(tableName string, result interface{}) error {
			if rewards, ok := result.(*[]*models.Reward); ok {
				*rewards = testRewards
			}
			return nil
		},
	}

	config := &config.Config{RewardsTable: "test-rewards"}
	repo := NewRewardRepository(mockRepo, config)

	results, err := repo.List()
	if err != nil {
		t.Errorf("List failed: %v", err)
	}

	if len(results) != 2 {
		t.Errorf("Expected 2 rewards, got %d", len(results))
	}
}

func TestRewardRepository_Update(t *testing.T) {
	existingReward := &models.Reward{
		ID:          "test-id",
		Title:       "Original Title",
		Description: "Original Description",
		Point:       50,
		CreatedAt:   time.Now().Add(-time.Hour),
	}

	mockRepo := &MockRepository{
		getItemFunc: func(tableName string, key map[string]interface{}, result interface{}) error {
			if reward, ok := result.(*models.Reward); ok {
				*reward = *existingReward
			}
			return nil
		},
	}

	config := &config.Config{RewardsTable: "test-rewards"}
	repo := NewRewardRepository(mockRepo, config)

	updatedReward := &models.Reward{
		ID:          "test-id",
		Title:       "Updated Title",
		Description: "Updated Description",
		Point:       100,
	}

	err := repo.Update(updatedReward)
	if err != nil {
		t.Errorf("Update failed: %v", err)
	}

	// 作成日時が保持されていることを確認
	if !updatedReward.CreatedAt.Equal(existingReward.CreatedAt) {
		t.Error("CreatedAt should be preserved during update")
	}
}

func TestRewardRepository_Update_ValidationError(t *testing.T) {
	mockRepo := &MockRepository{}
	config := &config.Config{RewardsTable: "test-rewards"}
	repo := NewRewardRepository(mockRepo, config)

	tests := []struct {
		name        string
		reward      *models.Reward
		expectedErr string
	}{
		{
			name:        "nil reward",
			reward:      nil,
			expectedErr: "validation error for field 'reward': reward cannot be nil",
		},
		{
			name: "empty id",
			reward: &models.Reward{
				ID:    "",
				Title: "Test",
				Point: 50,
			},
			expectedErr: "validation error for field 'id': id is required for update",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Update(tt.reward)
			if err == nil {
				t.Error("Expected validation error")
			}
			if err.Error() != tt.expectedErr {
				t.Errorf("Expected error '%s', got '%s'", tt.expectedErr, err.Error())
			}
		})
	}
}

func TestRewardRepository_Delete(t *testing.T) {
	existingReward := &models.Reward{
		ID:          "test-id",
		Title:       "Test Reward",
		Description: "Test Description",
		Point:       50,
		CreatedAt:   time.Now(),
	}

	mockRepo := &MockRepository{
		getItemFunc: func(tableName string, key map[string]interface{}, result interface{}) error {
			if reward, ok := result.(*models.Reward); ok {
				*reward = *existingReward
			}
			return nil
		},
	}

	config := &config.Config{RewardsTable: "test-rewards"}
	repo := NewRewardRepository(mockRepo, config)

	err := repo.Delete("test-id")
	if err != nil {
		t.Errorf("Delete failed: %v", err)
	}
}

func TestRewardRepository_Delete_NotFound(t *testing.T) {
	mockRepo := &MockRepository{
		getItemFunc: func(tableName string, key map[string]interface{}, result interface{}) error {
			return fmt.Errorf("item not found in table test-rewards")
		},
	}

	config := &config.Config{RewardsTable: "test-rewards"}
	repo := NewRewardRepository(mockRepo, config)

	err := repo.Delete("non-existent-id")
	if err != errors.ErrNotFound {
		t.Errorf("Expected ErrNotFound, got %v", err)
	}
}

func TestRewardRepository_Delete_EmptyID(t *testing.T) {
	mockRepo := &MockRepository{}
	config := &config.Config{RewardsTable: "test-rewards"}
	repo := NewRewardRepository(mockRepo, config)

	err := repo.Delete("")
	if err == nil {
		t.Error("Expected validation error for empty ID")
	}
}