package services

import (
	"testing"
	"time"

	"achievement-management/internal/errors"
	"achievement-management/internal/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAchievementRepository モック達成目録リポジトリ
type MockAchievementRepository struct {
	mock.Mock
}

func (m *MockAchievementRepository) Create(achievement *models.Achievement) error {
	args := m.Called(achievement)
	return args.Error(0)
}

func (m *MockAchievementRepository) Update(achievement *models.Achievement) error {
	args := m.Called(achievement)
	return args.Error(0)
}

func (m *MockAchievementRepository) GetByID(id string) (*models.Achievement, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Achievement), args.Error(1)
}

func (m *MockAchievementRepository) List() ([]*models.Achievement, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Achievement), args.Error(1)
}

func (m *MockAchievementRepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

// MockPointRepository モックポイントリポジトリ
type MockPointRepository struct {
	mock.Mock
}

func (m *MockPointRepository) GetCurrentPoints() (*models.CurrentPoints, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.CurrentPoints), args.Error(1)
}

func (m *MockPointRepository) UpdateCurrentPoints(points *models.CurrentPoints) error {
	args := m.Called(points)
	return args.Error(0)
}

func (m *MockPointRepository) CreateRewardHistory(history *models.RewardHistory) error {
	args := m.Called(history)
	return args.Error(0)
}

func (m *MockPointRepository) GetRewardHistory() ([]*models.RewardHistory, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.RewardHistory), args.Error(1)
}

func (m *MockPointRepository) TransactPointsAndHistory(pointsUpdate *models.CurrentPoints, history *models.RewardHistory) error {
	args := m.Called(pointsUpdate, history)
	return args.Error(0)
}

func (m *MockPointRepository) AddPoints(points int) error {
	args := m.Called(points)
	return args.Error(0)
}

func (m *MockPointRepository) SubtractPoints(points int) error {
	args := m.Called(points)
	return args.Error(0)
}

func TestAchievementService_Create(t *testing.T) {
	tests := []struct {
		name                string
		achievement         *models.Achievement
		setupMocks          func(*MockAchievementRepository, *MockPointRepository)
		expectedError       error
		expectedErrorType   interface{}
		expectedErrorField  string
	}{
		{
			name: "正常な達成目録作成",
			achievement: &models.Achievement{
				Title:       "テスト達成目録",
				Description: "テスト用の達成目録です",
				Point:       100,
			},
			setupMocks: func(achievementRepo *MockAchievementRepository, pointRepo *MockPointRepository) {
				achievementRepo.On("Create", mock.AnythingOfType("*models.Achievement")).Return(nil)
				pointRepo.On("AddPoints", 100).Return(nil)
			},
			expectedError: nil,
		},
		{
			name:        "nilの達成目録",
			achievement: nil,
			setupMocks: func(achievementRepo *MockAchievementRepository, pointRepo *MockPointRepository) {
				// モックの設定は不要
			},
			expectedError:      &errors.ValidationError{},
			expectedErrorType:  &errors.ValidationError{},
			expectedErrorField: "achievement",
		},
		{
			name: "タイトルが空の達成目録",
			achievement: &models.Achievement{
				Title:       "",
				Description: "テスト用の達成目録です",
				Point:       100,
			},
			setupMocks: func(achievementRepo *MockAchievementRepository, pointRepo *MockPointRepository) {
				// モックの設定は不要
			},
			expectedError:      &errors.ValidationError{},
			expectedErrorType:  &errors.ValidationError{},
			expectedErrorField: "title",
		},
		{
			name: "ポイントが0の達成目録",
			achievement: &models.Achievement{
				Title:       "テスト達成目録",
				Description: "テスト用の達成目録です",
				Point:       0,
			},
			setupMocks: func(achievementRepo *MockAchievementRepository, pointRepo *MockPointRepository) {
				// モックの設定は不要
			},
			expectedError:      &errors.ValidationError{},
			expectedErrorType:  &errors.ValidationError{},
			expectedErrorField: "point",
		},
		{
			name: "ポイントが負の達成目録",
			achievement: &models.Achievement{
				Title:       "テスト達成目録",
				Description: "テスト用の達成目録です",
				Point:       -10,
			},
			setupMocks: func(achievementRepo *MockAchievementRepository, pointRepo *MockPointRepository) {
				// モックの設定は不要
			},
			expectedError:      &errors.ValidationError{},
			expectedErrorType:  &errors.ValidationError{},
			expectedErrorField: "point",
		},
		{
			name: "達成目録作成エラー",
			achievement: &models.Achievement{
				Title:       "テスト達成目録",
				Description: "テスト用の達成目録です",
				Point:       100,
			},
			setupMocks: func(achievementRepo *MockAchievementRepository, pointRepo *MockPointRepository) {
				achievementRepo.On("Create", mock.AnythingOfType("*models.Achievement")).Return(&errors.DatabaseError{})
			},
			expectedError:     &errors.DatabaseError{},
			expectedErrorType: &errors.DatabaseError{},
		},
		{
			name: "ポイント加算エラー（ロールバック成功）",
			achievement: &models.Achievement{
				ID:          "test-id",
				Title:       "テスト達成目録",
				Description: "テスト用の達成目録です",
				Point:       100,
			},
			setupMocks: func(achievementRepo *MockAchievementRepository, pointRepo *MockPointRepository) {
				achievementRepo.On("Create", mock.AnythingOfType("*models.Achievement")).Return(nil)
				pointRepo.On("AddPoints", 100).Return(&errors.DatabaseError{})
				achievementRepo.On("Delete", "test-id").Return(nil)
			},
			expectedError:     &errors.DatabaseError{},
			expectedErrorType: &errors.DatabaseError{},
		},
		{
			name: "ポイント加算エラー（ロールバック失敗）",
			achievement: &models.Achievement{
				ID:          "test-id",
				Title:       "テスト達成目録",
				Description: "テスト用の達成目録です",
				Point:       100,
			},
			setupMocks: func(achievementRepo *MockAchievementRepository, pointRepo *MockPointRepository) {
				achievementRepo.On("Create", mock.AnythingOfType("*models.Achievement")).Return(nil)
				pointRepo.On("AddPoints", 100).Return(&errors.DatabaseError{})
				achievementRepo.On("Delete", "test-id").Return(&errors.DatabaseError{})
			},
			expectedError:     &errors.DatabaseError{},
			expectedErrorType: &errors.DatabaseError{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			achievementRepo := new(MockAchievementRepository)
			pointRepo := new(MockPointRepository)
			
			tt.setupMocks(achievementRepo, pointRepo)
			
			service := NewAchievementService(achievementRepo, pointRepo)
			err := service.Create(tt.achievement)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.IsType(t, tt.expectedErrorType, err)
				
				if tt.expectedErrorField != "" {
					if validationErr, ok := err.(*errors.ValidationError); ok {
						assert.Equal(t, tt.expectedErrorField, validationErr.Field)
					}
				}
			} else {
				assert.NoError(t, err)
			}

			achievementRepo.AssertExpectations(t)
			pointRepo.AssertExpectations(t)
		})
	}
}

func TestAchievementService_Update(t *testing.T) {
	tests := []struct {
		name                string
		id                  string
		achievement         *models.Achievement
		setupMocks          func(*MockAchievementRepository, *MockPointRepository)
		expectedError       error
		expectedErrorType   interface{}
		expectedErrorField  string
	}{
		{
			name: "正常な達成目録更新",
			id:   "test-id",
			achievement: &models.Achievement{
				Title:       "更新されたテスト達成目録",
				Description: "更新されたテスト用の達成目録です",
				Point:       150,
			},
			setupMocks: func(achievementRepo *MockAchievementRepository, pointRepo *MockPointRepository) {
				achievementRepo.On("Update", mock.MatchedBy(func(a *models.Achievement) bool {
					return a.ID == "test-id" && a.Title == "更新されたテスト達成目録"
				})).Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "IDが空",
			id:   "",
			achievement: &models.Achievement{
				Title:       "テスト達成目録",
				Description: "テスト用の達成目録です",
				Point:       100,
			},
			setupMocks: func(achievementRepo *MockAchievementRepository, pointRepo *MockPointRepository) {
				// モックの設定は不要
			},
			expectedError:      &errors.ValidationError{},
			expectedErrorType:  &errors.ValidationError{},
			expectedErrorField: "id",
		},
		{
			name:        "nilの達成目録",
			id:          "test-id",
			achievement: nil,
			setupMocks: func(achievementRepo *MockAchievementRepository, pointRepo *MockPointRepository) {
				// モックの設定は不要
			},
			expectedError:      &errors.ValidationError{},
			expectedErrorType:  &errors.ValidationError{},
			expectedErrorField: "achievement",
		},
		{
			name: "タイトルが空",
			id:   "test-id",
			achievement: &models.Achievement{
				Title:       "",
				Description: "テスト用の達成目録です",
				Point:       100,
			},
			setupMocks: func(achievementRepo *MockAchievementRepository, pointRepo *MockPointRepository) {
				// モックの設定は不要
			},
			expectedError:      &errors.ValidationError{},
			expectedErrorType:  &errors.ValidationError{},
			expectedErrorField: "title",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			achievementRepo := new(MockAchievementRepository)
			pointRepo := new(MockPointRepository)
			
			tt.setupMocks(achievementRepo, pointRepo)
			
			service := NewAchievementService(achievementRepo, pointRepo)
			err := service.Update(tt.id, tt.achievement)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.IsType(t, tt.expectedErrorType, err)
				
				if tt.expectedErrorField != "" {
					if validationErr, ok := err.(*errors.ValidationError); ok {
						assert.Equal(t, tt.expectedErrorField, validationErr.Field)
					}
				}
			} else {
				assert.NoError(t, err)
			}

			achievementRepo.AssertExpectations(t)
			pointRepo.AssertExpectations(t)
		})
	}
}

func TestAchievementService_GetByID(t *testing.T) {
	tests := []struct {
		name                string
		id                  string
		setupMocks          func(*MockAchievementRepository, *MockPointRepository)
		expectedAchievement *models.Achievement
		expectedError       error
		expectedErrorType   interface{}
		expectedErrorField  string
	}{
		{
			name: "正常な達成目録取得",
			id:   "test-id",
			setupMocks: func(achievementRepo *MockAchievementRepository, pointRepo *MockPointRepository) {
				achievement := &models.Achievement{
					ID:          "test-id",
					Title:       "テスト達成目録",
					Description: "テスト用の達成目録です",
					Point:       100,
					CreatedAt:   time.Now(),
				}
				achievementRepo.On("GetByID", "test-id").Return(achievement, nil)
			},
			expectedAchievement: &models.Achievement{
				ID:          "test-id",
				Title:       "テスト達成目録",
				Description: "テスト用の達成目録です",
				Point:       100,
			},
			expectedError: nil,
		},
		{
			name: "IDが空",
			id:   "",
			setupMocks: func(achievementRepo *MockAchievementRepository, pointRepo *MockPointRepository) {
				// モックの設定は不要
			},
			expectedAchievement: nil,
			expectedError:       &errors.ValidationError{},
			expectedErrorType:   &errors.ValidationError{},
			expectedErrorField:  "id",
		},
		{
			name: "存在しない達成目録",
			id:   "non-existent-id",
			setupMocks: func(achievementRepo *MockAchievementRepository, pointRepo *MockPointRepository) {
				achievementRepo.On("GetByID", "non-existent-id").Return(nil, errors.ErrNotFound)
			},
			expectedAchievement: nil,
			expectedError:       errors.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			achievementRepo := new(MockAchievementRepository)
			pointRepo := new(MockPointRepository)
			
			tt.setupMocks(achievementRepo, pointRepo)
			
			service := NewAchievementService(achievementRepo, pointRepo)
			achievement, err := service.GetByID(tt.id)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Nil(t, achievement)
				
				if tt.expectedErrorType != nil {
					assert.IsType(t, tt.expectedErrorType, err)
				}
				
				if tt.expectedErrorField != "" {
					if validationErr, ok := err.(*errors.ValidationError); ok {
						assert.Equal(t, tt.expectedErrorField, validationErr.Field)
					}
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, achievement)
				assert.Equal(t, tt.expectedAchievement.ID, achievement.ID)
				assert.Equal(t, tt.expectedAchievement.Title, achievement.Title)
				assert.Equal(t, tt.expectedAchievement.Description, achievement.Description)
				assert.Equal(t, tt.expectedAchievement.Point, achievement.Point)
			}

			achievementRepo.AssertExpectations(t)
			pointRepo.AssertExpectations(t)
		})
	}
}

func TestAchievementService_List(t *testing.T) {
	tests := []struct {
		name                 string
		setupMocks           func(*MockAchievementRepository, *MockPointRepository)
		expectedAchievements []*models.Achievement
		expectedError        error
	}{
		{
			name: "正常な達成目録一覧取得",
			setupMocks: func(achievementRepo *MockAchievementRepository, pointRepo *MockPointRepository) {
				achievements := []*models.Achievement{
					{
						ID:          "test-id-1",
						Title:       "テスト達成目録1",
						Description: "テスト用の達成目録です1",
						Point:       100,
						CreatedAt:   time.Now(),
					},
					{
						ID:          "test-id-2",
						Title:       "テスト達成目録2",
						Description: "テスト用の達成目録です2",
						Point:       200,
						CreatedAt:   time.Now(),
					},
				}
				achievementRepo.On("List").Return(achievements, nil)
			},
			expectedAchievements: []*models.Achievement{
				{
					ID:          "test-id-1",
					Title:       "テスト達成目録1",
					Description: "テスト用の達成目録です1",
					Point:       100,
				},
				{
					ID:          "test-id-2",
					Title:       "テスト達成目録2",
					Description: "テスト用の達成目録です2",
					Point:       200,
				},
			},
			expectedError: nil,
		},
		{
			name: "空の達成目録一覧",
			setupMocks: func(achievementRepo *MockAchievementRepository, pointRepo *MockPointRepository) {
				achievementRepo.On("List").Return([]*models.Achievement{}, nil)
			},
			expectedAchievements: []*models.Achievement{},
			expectedError:        nil,
		},
		{
			name: "データベースエラー",
			setupMocks: func(achievementRepo *MockAchievementRepository, pointRepo *MockPointRepository) {
				achievementRepo.On("List").Return(nil, &errors.DatabaseError{})
			},
			expectedAchievements: nil,
			expectedError:        &errors.DatabaseError{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			achievementRepo := new(MockAchievementRepository)
			pointRepo := new(MockPointRepository)
			
			tt.setupMocks(achievementRepo, pointRepo)
			
			service := NewAchievementService(achievementRepo, pointRepo)
			achievements, err := service.List()

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Nil(t, achievements)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, len(tt.expectedAchievements), len(achievements))
				
				for i, expected := range tt.expectedAchievements {
					assert.Equal(t, expected.ID, achievements[i].ID)
					assert.Equal(t, expected.Title, achievements[i].Title)
					assert.Equal(t, expected.Description, achievements[i].Description)
					assert.Equal(t, expected.Point, achievements[i].Point)
				}
			}

			achievementRepo.AssertExpectations(t)
			pointRepo.AssertExpectations(t)
		})
	}
}

func TestAchievementService_Delete(t *testing.T) {
	tests := []struct {
		name               string
		id                 string
		setupMocks         func(*MockAchievementRepository, *MockPointRepository)
		expectedError      error
		expectedErrorType  interface{}
		expectedErrorField string
	}{
		{
			name: "正常な達成目録削除",
			id:   "test-id",
			setupMocks: func(achievementRepo *MockAchievementRepository, pointRepo *MockPointRepository) {
				achievementRepo.On("Delete", "test-id").Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "IDが空",
			id:   "",
			setupMocks: func(achievementRepo *MockAchievementRepository, pointRepo *MockPointRepository) {
				// モックの設定は不要
			},
			expectedError:      &errors.ValidationError{},
			expectedErrorType:  &errors.ValidationError{},
			expectedErrorField: "id",
		},
		{
			name: "存在しない達成目録の削除",
			id:   "non-existent-id",
			setupMocks: func(achievementRepo *MockAchievementRepository, pointRepo *MockPointRepository) {
				achievementRepo.On("Delete", "non-existent-id").Return(errors.ErrNotFound)
			},
			expectedError: errors.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			achievementRepo := new(MockAchievementRepository)
			pointRepo := new(MockPointRepository)
			
			tt.setupMocks(achievementRepo, pointRepo)
			
			service := NewAchievementService(achievementRepo, pointRepo)
			err := service.Delete(tt.id)

			if tt.expectedError != nil {
				assert.Error(t, err)
				
				if tt.expectedErrorType != nil {
					assert.IsType(t, tt.expectedErrorType, err)
				}
				
				if tt.expectedErrorField != "" {
					if validationErr, ok := err.(*errors.ValidationError); ok {
						assert.Equal(t, tt.expectedErrorField, validationErr.Field)
					}
				}
			} else {
				assert.NoError(t, err)
			}

			achievementRepo.AssertExpectations(t)
			pointRepo.AssertExpectations(t)
		})
	}
}