package repository

import (
	"fmt"
	"testing"
	"time"

	"achievement-management/internal/config"
	"achievement-management/internal/errors"
	"achievement-management/internal/models"
)

// MockRepository リポジトリのモック
type MockRepository struct {
	putItemFunc    func(tableName string, item interface{}) error
	getItemFunc    func(tableName string, key map[string]interface{}, result interface{}) error
	scanFunc       func(tableName string, result interface{}) error
	deleteItemFunc func(tableName string, key map[string]interface{}) error
}

func (m *MockRepository) PutItem(tableName string, item interface{}) error {
	if m.putItemFunc != nil {
		return m.putItemFunc(tableName, item)
	}
	return nil
}

func (m *MockRepository) GetItem(tableName string, key map[string]interface{}, result interface{}) error {
	if m.getItemFunc != nil {
		return m.getItemFunc(tableName, key, result)
	}
	return nil
}

func (m *MockRepository) UpdateItem(tableName string, key map[string]interface{}, updateExpression string, expressionAttributeValues map[string]interface{}) error {
	return nil
}

func (m *MockRepository) Scan(tableName string, result interface{}) error {
	if m.scanFunc != nil {
		return m.scanFunc(tableName, result)
	}
	return nil
}

func (m *MockRepository) DeleteItem(tableName string, key map[string]interface{}) error {
	if m.deleteItemFunc != nil {
		return m.deleteItemFunc(tableName, key)
	}
	return nil
}

func (m *MockRepository) TransactWrite(items []TransactWriteItem) error {
	return nil
}

func TestAchievementRepository_Create(t *testing.T) {
	mockRepo := &MockRepository{}
	config := &config.Config{
		Tables: config.TableConfig{
			Achievements: "test-achievements",
		},
	}
	repo := NewAchievementRepository(mockRepo, config)

	achievement := &models.Achievement{
		Title:       "Test Achievement",
		Description: "Test Description",
		Point:       100,
	}

	err := repo.Create(achievement)
	if err != nil {
		t.Errorf("Create failed: %v", err)
	}

	// IDが生成されていることを確認
	if achievement.ID == "" {
		t.Error("ID should be generated")
	}

	// 作成日時が設定されていることを確認
	if achievement.CreatedAt.IsZero() {
		t.Error("CreatedAt should be set")
	}
}

func TestAchievementRepository_Create_ValidationError(t *testing.T) {
	mockRepo := &MockRepository{}
	config := &config.Config{
Tables: config.TableConfig{
Achievements: "test-achievements",
},
}
	repo := NewAchievementRepository(mockRepo, config)

	tests := []struct {
		name        string
		achievement *models.Achievement
		expectedErr string
	}{
		{
			name:        "nil achievement",
			achievement: nil,
			expectedErr: "validation error for field 'achievement': achievement cannot be nil",
		},
		{
			name: "empty title",
			achievement: &models.Achievement{
				Title: "",
				Point: 100,
			},
			expectedErr: "validation error for field 'title': title is required",
		},
		{
			name: "zero point",
			achievement: &models.Achievement{
				Title: "Test",
				Point: 0,
			},
			expectedErr: "validation error for field 'point': point must be positive",
		},
		{
			name: "negative point",
			achievement: &models.Achievement{
				Title: "Test",
				Point: -10,
			},
			expectedErr: "validation error for field 'point': point must be positive",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Create(tt.achievement)
			if err == nil {
				t.Error("Expected validation error")
			}
			if err.Error() != tt.expectedErr {
				t.Errorf("Expected error '%s', got '%s'", tt.expectedErr, err.Error())
			}
		})
	}
}

func TestAchievementRepository_GetByID(t *testing.T) {
	testAchievement := &models.Achievement{
		ID:          "test-id",
		Title:       "Test Achievement",
		Description: "Test Description",
		Point:       100,
		CreatedAt:   time.Now(),
	}

	mockRepo := &MockRepository{
		getItemFunc: func(tableName string, key map[string]interface{}, result interface{}) error {
			if achievement, ok := result.(*models.Achievement); ok {
				*achievement = *testAchievement
			}
			return nil
		},
	}

	config := &config.Config{
Tables: config.TableConfig{
Achievements: "test-achievements",
},
}
	repo := NewAchievementRepository(mockRepo, config)

	result, err := repo.GetByID("test-id")
	if err != nil {
		t.Errorf("GetByID failed: %v", err)
	}

	if result.ID != testAchievement.ID {
		t.Errorf("Expected ID %s, got %s", testAchievement.ID, result.ID)
	}
}

func TestAchievementRepository_GetByID_NotFound(t *testing.T) {
	mockRepo := &MockRepository{
		getItemFunc: func(tableName string, key map[string]interface{}, result interface{}) error {
			return fmt.Errorf("item not found in table test-achievements")
		},
	}

	config := &config.Config{
Tables: config.TableConfig{
Achievements: "test-achievements",
},
}
	repo := NewAchievementRepository(mockRepo, config)

	_, err := repo.GetByID("non-existent-id")
	if err != errors.ErrNotFound {
		t.Errorf("Expected ErrNotFound, got %v", err)
	}
}

func TestAchievementRepository_GetByID_EmptyID(t *testing.T) {
	mockRepo := &MockRepository{}
	config := &config.Config{
Tables: config.TableConfig{
Achievements: "test-achievements",
},
}
	repo := NewAchievementRepository(mockRepo, config)

	_, err := repo.GetByID("")
	if err == nil {
		t.Error("Expected validation error for empty ID")
	}
}

func TestAchievementRepository_List(t *testing.T) {
	testAchievements := []*models.Achievement{
		{
			ID:          "test-id-1",
			Title:       "Test Achievement 1",
			Description: "Test Description 1",
			Point:       100,
			CreatedAt:   time.Now(),
		},
		{
			ID:          "test-id-2",
			Title:       "Test Achievement 2",
			Description: "Test Description 2",
			Point:       200,
			CreatedAt:   time.Now(),
		},
	}

	mockRepo := &MockRepository{
		scanFunc: func(tableName string, result interface{}) error {
			if achievements, ok := result.(*[]*models.Achievement); ok {
				*achievements = testAchievements
			}
			return nil
		},
	}

	config := &config.Config{
Tables: config.TableConfig{
Achievements: "test-achievements",
},
}
	repo := NewAchievementRepository(mockRepo, config)

	results, err := repo.List()
	if err != nil {
		t.Errorf("List failed: %v", err)
	}

	if len(results) != 2 {
		t.Errorf("Expected 2 achievements, got %d", len(results))
	}
}

func TestAchievementRepository_Update(t *testing.T) {
	existingAchievement := &models.Achievement{
		ID:          "test-id",
		Title:       "Original Title",
		Description: "Original Description",
		Point:       100,
		CreatedAt:   time.Now().Add(-time.Hour),
	}

	mockRepo := &MockRepository{
		getItemFunc: func(tableName string, key map[string]interface{}, result interface{}) error {
			if achievement, ok := result.(*models.Achievement); ok {
				*achievement = *existingAchievement
			}
			return nil
		},
	}

	config := &config.Config{
Tables: config.TableConfig{
Achievements: "test-achievements",
},
}
	repo := NewAchievementRepository(mockRepo, config)

	updatedAchievement := &models.Achievement{
		ID:          "test-id",
		Title:       "Updated Title",
		Description: "Updated Description",
		Point:       200,
	}

	err := repo.Update(updatedAchievement)
	if err != nil {
		t.Errorf("Update failed: %v", err)
	}

	// 作成日時が保持されていることを確認
	if !updatedAchievement.CreatedAt.Equal(existingAchievement.CreatedAt) {
		t.Error("CreatedAt should be preserved during update")
	}
}

func TestAchievementRepository_Delete(t *testing.T) {
	existingAchievement := &models.Achievement{
		ID:          "test-id",
		Title:       "Test Achievement",
		Description: "Test Description",
		Point:       100,
		CreatedAt:   time.Now(),
	}

	mockRepo := &MockRepository{
		getItemFunc: func(tableName string, key map[string]interface{}, result interface{}) error {
			if achievement, ok := result.(*models.Achievement); ok {
				*achievement = *existingAchievement
			}
			return nil
		},
	}

	config := &config.Config{
Tables: config.TableConfig{
Achievements: "test-achievements",
},
}
	repo := NewAchievementRepository(mockRepo, config)

	err := repo.Delete("test-id")
	if err != nil {
		t.Errorf("Delete failed: %v", err)
	}
}

func TestAchievementRepository_Delete_NotFound(t *testing.T) {
	mockRepo := &MockRepository{
		getItemFunc: func(tableName string, key map[string]interface{}, result interface{}) error {
			return fmt.Errorf("item not found in table test-achievements")
		},
	}

	config := &config.Config{
Tables: config.TableConfig{
Achievements: "test-achievements",
},
}
	repo := NewAchievementRepository(mockRepo, config)

	err := repo.Delete("non-existent-id")
	if err != errors.ErrNotFound {
		t.Errorf("Expected ErrNotFound, got %v", err)
	}
}