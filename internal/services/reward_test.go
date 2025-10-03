package services

import (
	"testing"
	"time"

	"achievement-management/internal/errors"
	"achievement-management/internal/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRewardRepository モック報酬リポジトリ
type MockRewardRepository struct {
	mock.Mock
}

func (m *MockRewardRepository) Create(reward *models.Reward) error {
	args := m.Called(reward)
	return args.Error(0)
}

func (m *MockRewardRepository) Update(reward *models.Reward) error {
	args := m.Called(reward)
	return args.Error(0)
}

func (m *MockRewardRepository) GetByID(id string) (*models.Reward, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Reward), args.Error(1)
}

func (m *MockRewardRepository) List() ([]*models.Reward, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Reward), args.Error(1)
}

func (m *MockRewardRepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestRewardService_Create(t *testing.T) {
	tests := []struct {
		name               string
		reward             *models.Reward
		setupMocks         func(*MockRewardRepository, *MockPointRepository)
		expectedError      error
		expectedErrorType  interface{}
		expectedErrorField string
	}{
		{
			name: "正常な報酬作成",
			reward: &models.Reward{
				Title:       "テスト報酬",
				Description: "テスト用の報酬です",
				Point:       50,
			},
			setupMocks: func(rewardRepo *MockRewardRepository, pointRepo *MockPointRepository) {
				rewardRepo.On("Create", mock.AnythingOfType("*models.Reward")).Return(nil)
			},
			expectedError: nil,
		},
		{
			name:   "nilの報酬",
			reward: nil,
			setupMocks: func(rewardRepo *MockRewardRepository, pointRepo *MockPointRepository) {
				// モックの設定は不要
			},
			expectedError:      &errors.ValidationError{},
			expectedErrorType:  &errors.ValidationError{},
			expectedErrorField: "reward",
		},
		{
			name: "タイトルが空の報酬",
			reward: &models.Reward{
				Title:       "",
				Description: "テスト用の報酬です",
				Point:       50,
			},
			setupMocks: func(rewardRepo *MockRewardRepository, pointRepo *MockPointRepository) {
				// モックの設定は不要
			},
			expectedError:      &errors.ValidationError{},
			expectedErrorType:  &errors.ValidationError{},
			expectedErrorField: "title",
		},
		{
			name: "ポイントが0の報酬",
			reward: &models.Reward{
				Title:       "テスト報酬",
				Description: "テスト用の報酬です",
				Point:       0,
			},
			setupMocks: func(rewardRepo *MockRewardRepository, pointRepo *MockPointRepository) {
				// モックの設定は不要
			},
			expectedError:      &errors.ValidationError{},
			expectedErrorType:  &errors.ValidationError{},
			expectedErrorField: "point",
		},
		{
			name: "ポイントが負の報酬",
			reward: &models.Reward{
				Title:       "テスト報酬",
				Description: "テスト用の報酬です",
				Point:       -10,
			},
			setupMocks: func(rewardRepo *MockRewardRepository, pointRepo *MockPointRepository) {
				// モックの設定は不要
			},
			expectedError:      &errors.ValidationError{},
			expectedErrorType:  &errors.ValidationError{},
			expectedErrorField: "point",
		},
		{
			name: "報酬作成エラー",
			reward: &models.Reward{
				Title:       "テスト報酬",
				Description: "テスト用の報酬です",
				Point:       50,
			},
			setupMocks: func(rewardRepo *MockRewardRepository, pointRepo *MockPointRepository) {
				rewardRepo.On("Create", mock.AnythingOfType("*models.Reward")).Return(&errors.DatabaseError{})
			},
			expectedError:     &errors.DatabaseError{},
			expectedErrorType: &errors.DatabaseError{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rewardRepo := new(MockRewardRepository)
			pointRepo := new(MockPointRepository)

			tt.setupMocks(rewardRepo, pointRepo)

			service := NewRewardService(rewardRepo, pointRepo)
			err := service.Create(tt.reward)

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

			rewardRepo.AssertExpectations(t)
			pointRepo.AssertExpectations(t)
		})
	}
}

func TestRewardService_Update(t *testing.T) {
	tests := []struct {
		name               string
		id                 string
		reward             *models.Reward
		setupMocks         func(*MockRewardRepository, *MockPointRepository)
		expectedError      error
		expectedErrorType  interface{}
		expectedErrorField string
	}{
		{
			name: "正常な報酬更新",
			id:   "test-id",
			reward: &models.Reward{
				Title:       "更新されたテスト報酬",
				Description: "更新されたテスト用の報酬です",
				Point:       75,
			},
			setupMocks: func(rewardRepo *MockRewardRepository, pointRepo *MockPointRepository) {
				rewardRepo.On("Update", mock.MatchedBy(func(r *models.Reward) bool {
					return r.ID == "test-id" && r.Title == "更新されたテスト報酬"
				})).Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "IDが空",
			id:   "",
			reward: &models.Reward{
				Title:       "テスト報酬",
				Description: "テスト用の報酬です",
				Point:       50,
			},
			setupMocks: func(rewardRepo *MockRewardRepository, pointRepo *MockPointRepository) {
				// モックの設定は不要
			},
			expectedError:      &errors.ValidationError{},
			expectedErrorType:  &errors.ValidationError{},
			expectedErrorField: "id",
		},
		{
			name:   "nilの報酬",
			id:     "test-id",
			reward: nil,
			setupMocks: func(rewardRepo *MockRewardRepository, pointRepo *MockPointRepository) {
				// モックの設定は不要
			},
			expectedError:      &errors.ValidationError{},
			expectedErrorType:  &errors.ValidationError{},
			expectedErrorField: "reward",
		},
		{
			name: "タイトルが空",
			id:   "test-id",
			reward: &models.Reward{
				Title:       "",
				Description: "テスト用の報酬です",
				Point:       50,
			},
			setupMocks: func(rewardRepo *MockRewardRepository, pointRepo *MockPointRepository) {
				// モックの設定は不要
			},
			expectedError:      &errors.ValidationError{},
			expectedErrorType:  &errors.ValidationError{},
			expectedErrorField: "title",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rewardRepo := new(MockRewardRepository)
			pointRepo := new(MockPointRepository)

			tt.setupMocks(rewardRepo, pointRepo)

			service := NewRewardService(rewardRepo, pointRepo)
			err := service.Update(tt.id, tt.reward)

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

			rewardRepo.AssertExpectations(t)
			pointRepo.AssertExpectations(t)
		})
	}
}

func TestRewardService_GetByID(t *testing.T) {
	tests := []struct {
		name               string
		id                 string
		setupMocks         func(*MockRewardRepository, *MockPointRepository)
		expectedReward     *models.Reward
		expectedError      error
		expectedErrorType  interface{}
		expectedErrorField string
	}{
		{
			name: "正常な報酬取得",
			id:   "test-id",
			setupMocks: func(rewardRepo *MockRewardRepository, pointRepo *MockPointRepository) {
				reward := &models.Reward{
					ID:          "test-id",
					Title:       "テスト報酬",
					Description: "テスト用の報酬です",
					Point:       50,
					CreatedAt:   time.Now(),
				}
				rewardRepo.On("GetByID", "test-id").Return(reward, nil)
			},
			expectedReward: &models.Reward{
				ID:          "test-id",
				Title:       "テスト報酬",
				Description: "テスト用の報酬です",
				Point:       50,
			},
			expectedError: nil,
		},
		{
			name: "IDが空",
			id:   "",
			setupMocks: func(rewardRepo *MockRewardRepository, pointRepo *MockPointRepository) {
				// モックの設定は不要
			},
			expectedReward:     nil,
			expectedError:      &errors.ValidationError{},
			expectedErrorType:  &errors.ValidationError{},
			expectedErrorField: "id",
		},
		{
			name: "存在しない報酬",
			id:   "non-existent-id",
			setupMocks: func(rewardRepo *MockRewardRepository, pointRepo *MockPointRepository) {
				rewardRepo.On("GetByID", "non-existent-id").Return(nil, errors.ErrNotFound)
			},
			expectedReward: nil,
			expectedError:  errors.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rewardRepo := new(MockRewardRepository)
			pointRepo := new(MockPointRepository)

			tt.setupMocks(rewardRepo, pointRepo)

			service := NewRewardService(rewardRepo, pointRepo)
			reward, err := service.GetByID(tt.id)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Nil(t, reward)

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
				assert.NotNil(t, reward)
				assert.Equal(t, tt.expectedReward.ID, reward.ID)
				assert.Equal(t, tt.expectedReward.Title, reward.Title)
				assert.Equal(t, tt.expectedReward.Description, reward.Description)
				assert.Equal(t, tt.expectedReward.Point, reward.Point)
			}

			rewardRepo.AssertExpectations(t)
			pointRepo.AssertExpectations(t)
		})
	}
}

func TestRewardService_List(t *testing.T) {
	tests := []struct {
		name            string
		setupMocks      func(*MockRewardRepository, *MockPointRepository)
		expectedRewards []*models.Reward
		expectedError   error
	}{
		{
			name: "正常な報酬一覧取得",
			setupMocks: func(rewardRepo *MockRewardRepository, pointRepo *MockPointRepository) {
				rewards := []*models.Reward{
					{
						ID:          "test-id-1",
						Title:       "テスト報酬1",
						Description: "テスト用の報酬です1",
						Point:       50,
						CreatedAt:   time.Now(),
					},
					{
						ID:          "test-id-2",
						Title:       "テスト報酬2",
						Description: "テスト用の報酬です2",
						Point:       100,
						CreatedAt:   time.Now(),
					},
				}
				rewardRepo.On("List").Return(rewards, nil)
			},
			expectedRewards: []*models.Reward{
				{
					ID:          "test-id-1",
					Title:       "テスト報酬1",
					Description: "テスト用の報酬です1",
					Point:       50,
				},
				{
					ID:          "test-id-2",
					Title:       "テスト報酬2",
					Description: "テスト用の報酬です2",
					Point:       100,
				},
			},
			expectedError: nil,
		},
		{
			name: "空の報酬一覧",
			setupMocks: func(rewardRepo *MockRewardRepository, pointRepo *MockPointRepository) {
				rewardRepo.On("List").Return([]*models.Reward{}, nil)
			},
			expectedRewards: []*models.Reward{},
			expectedError:   nil,
		},
		{
			name: "データベースエラー",
			setupMocks: func(rewardRepo *MockRewardRepository, pointRepo *MockPointRepository) {
				rewardRepo.On("List").Return(nil, &errors.DatabaseError{})
			},
			expectedRewards: nil,
			expectedError:   &errors.DatabaseError{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rewardRepo := new(MockRewardRepository)
			pointRepo := new(MockPointRepository)

			tt.setupMocks(rewardRepo, pointRepo)

			service := NewRewardService(rewardRepo, pointRepo)
			rewards, err := service.List()

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Nil(t, rewards)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, len(tt.expectedRewards), len(rewards))

				for i, expected := range tt.expectedRewards {
					assert.Equal(t, expected.ID, rewards[i].ID)
					assert.Equal(t, expected.Title, rewards[i].Title)
					assert.Equal(t, expected.Description, rewards[i].Description)
					assert.Equal(t, expected.Point, rewards[i].Point)
				}
			}

			rewardRepo.AssertExpectations(t)
			pointRepo.AssertExpectations(t)
		})
	}
}

func TestRewardService_Delete(t *testing.T) {
	tests := []struct {
		name               string
		id                 string
		setupMocks         func(*MockRewardRepository, *MockPointRepository)
		expectedError      error
		expectedErrorType  interface{}
		expectedErrorField string
	}{
		{
			name: "正常な報酬削除",
			id:   "test-id",
			setupMocks: func(rewardRepo *MockRewardRepository, pointRepo *MockPointRepository) {
				rewardRepo.On("Delete", "test-id").Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "IDが空",
			id:   "",
			setupMocks: func(rewardRepo *MockRewardRepository, pointRepo *MockPointRepository) {
				// モックの設定は不要
			},
			expectedError:      &errors.ValidationError{},
			expectedErrorType:  &errors.ValidationError{},
			expectedErrorField: "id",
		},
		{
			name: "存在しない報酬の削除",
			id:   "non-existent-id",
			setupMocks: func(rewardRepo *MockRewardRepository, pointRepo *MockPointRepository) {
				rewardRepo.On("Delete", "non-existent-id").Return(errors.ErrNotFound)
			},
			expectedError: errors.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rewardRepo := new(MockRewardRepository)
			pointRepo := new(MockPointRepository)

			tt.setupMocks(rewardRepo, pointRepo)

			service := NewRewardService(rewardRepo, pointRepo)
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

			rewardRepo.AssertExpectations(t)
			pointRepo.AssertExpectations(t)
		})
	}
}

func TestRewardService_Redeem(t *testing.T) {
	tests := []struct {
		name               string
		rewardID           string
		setupMocks         func(*MockRewardRepository, *MockPointRepository)
		expectedError      error
		expectedErrorType  interface{}
		expectedErrorField string
	}{
		{
			name:     "正常な報酬獲得",
			rewardID: "test-reward-id",
			setupMocks: func(rewardRepo *MockRewardRepository, pointRepo *MockPointRepository) {
				reward := &models.Reward{
					ID:          "test-reward-id",
					Title:       "テスト報酬",
					Description: "テスト用の報酬です",
					Point:       50,
					CreatedAt:   time.Now(),
				}
				currentPoints := &models.CurrentPoints{
					ID:        "current",
					Point:     100,
					UpdatedAt: time.Now(),
				}
				rewardRepo.On("GetByID", "test-reward-id").Return(reward, nil)
				pointRepo.On("GetCurrentPoints").Return(currentPoints, nil)
				pointRepo.On("TransactPointsAndHistory", 
					mock.MatchedBy(func(p *models.CurrentPoints) bool {
						return p.Point == 50 // 100 - 50 = 50
					}),
					mock.MatchedBy(func(h *models.RewardHistory) bool {
						return h.RewardID == "test-reward-id" && h.RewardTitle == "テスト報酬" && h.PointCost == 50
					}),
				).Return(nil)
			},
			expectedError: nil,
		},
		{
			name:     "rewardIDが空",
			rewardID: "",
			setupMocks: func(rewardRepo *MockRewardRepository, pointRepo *MockPointRepository) {
				// モックの設定は不要
			},
			expectedError:      &errors.ValidationError{},
			expectedErrorType:  &errors.ValidationError{},
			expectedErrorField: "rewardID",
		},
		{
			name:     "存在しない報酬",
			rewardID: "non-existent-id",
			setupMocks: func(rewardRepo *MockRewardRepository, pointRepo *MockPointRepository) {
				rewardRepo.On("GetByID", "non-existent-id").Return(nil, errors.ErrNotFound)
			},
			expectedError: errors.ErrNotFound,
		},
		{
			name:     "現在のポイント取得エラー",
			rewardID: "test-reward-id",
			setupMocks: func(rewardRepo *MockRewardRepository, pointRepo *MockPointRepository) {
				reward := &models.Reward{
					ID:          "test-reward-id",
					Title:       "テスト報酬",
					Description: "テスト用の報酬です",
					Point:       50,
					CreatedAt:   time.Now(),
				}
				rewardRepo.On("GetByID", "test-reward-id").Return(reward, nil)
				pointRepo.On("GetCurrentPoints").Return(nil, &errors.DatabaseError{})
			},
			expectedError:     &errors.DatabaseError{},
			expectedErrorType: &errors.DatabaseError{},
		},
		{
			name:     "ポイント不足",
			rewardID: "test-reward-id",
			setupMocks: func(rewardRepo *MockRewardRepository, pointRepo *MockPointRepository) {
				reward := &models.Reward{
					ID:          "test-reward-id",
					Title:       "テスト報酬",
					Description: "テスト用の報酬です",
					Point:       100,
					CreatedAt:   time.Now(),
				}
				currentPoints := &models.CurrentPoints{
					ID:        "current",
					Point:     50, // 報酬のポイント（100）より少ない
					UpdatedAt: time.Now(),
				}
				rewardRepo.On("GetByID", "test-reward-id").Return(reward, nil)
				pointRepo.On("GetCurrentPoints").Return(currentPoints, nil)
			},
			expectedError:     &errors.BusinessLogicError{},
			expectedErrorType: &errors.BusinessLogicError{},
		},
		{
			name:     "トランザクション実行エラー",
			rewardID: "test-reward-id",
			setupMocks: func(rewardRepo *MockRewardRepository, pointRepo *MockPointRepository) {
				reward := &models.Reward{
					ID:          "test-reward-id",
					Title:       "テスト報酬",
					Description: "テスト用の報酬です",
					Point:       50,
					CreatedAt:   time.Now(),
				}
				currentPoints := &models.CurrentPoints{
					ID:        "current",
					Point:     100,
					UpdatedAt: time.Now(),
				}
				rewardRepo.On("GetByID", "test-reward-id").Return(reward, nil)
				pointRepo.On("GetCurrentPoints").Return(currentPoints, nil)
				pointRepo.On("TransactPointsAndHistory", 
					mock.MatchedBy(func(p *models.CurrentPoints) bool {
						return p.Point == 50
					}),
					mock.MatchedBy(func(h *models.RewardHistory) bool {
						return h.RewardID == "test-reward-id"
					}),
				).Return(&errors.DatabaseError{})
			},
			expectedError:     &errors.DatabaseError{},
			expectedErrorType: &errors.DatabaseError{},
		},
		{
			name:     "ちょうどのポイントで報酬獲得",
			rewardID: "test-reward-id",
			setupMocks: func(rewardRepo *MockRewardRepository, pointRepo *MockPointRepository) {
				reward := &models.Reward{
					ID:          "test-reward-id",
					Title:       "テスト報酬",
					Description: "テスト用の報酬です",
					Point:       100,
					CreatedAt:   time.Now(),
				}
				currentPoints := &models.CurrentPoints{
					ID:        "current",
					Point:     100, // ちょうど同じポイント
					UpdatedAt: time.Now(),
				}
				rewardRepo.On("GetByID", "test-reward-id").Return(reward, nil)
				pointRepo.On("GetCurrentPoints").Return(currentPoints, nil)
				pointRepo.On("TransactPointsAndHistory", 
					mock.MatchedBy(func(p *models.CurrentPoints) bool {
						return p.Point == 0 // 100 - 100 = 0
					}),
					mock.MatchedBy(func(h *models.RewardHistory) bool {
						return h.RewardID == "test-reward-id" && h.RewardTitle == "テスト報酬" && h.PointCost == 100
					}),
				).Return(nil)
			},
			expectedError: nil,
		},
		{
			name:     "1ポイント不足",
			rewardID: "test-reward-id",
			setupMocks: func(rewardRepo *MockRewardRepository, pointRepo *MockPointRepository) {
				reward := &models.Reward{
					ID:          "test-reward-id",
					Title:       "テスト報酬",
					Description: "テスト用の報酬です",
					Point:       100,
					CreatedAt:   time.Now(),
				}
				currentPoints := &models.CurrentPoints{
					ID:        "current",
					Point:     99, // 1ポイント不足
					UpdatedAt: time.Now(),
				}
				rewardRepo.On("GetByID", "test-reward-id").Return(reward, nil)
				pointRepo.On("GetCurrentPoints").Return(currentPoints, nil)
			},
			expectedError:     &errors.BusinessLogicError{},
			expectedErrorType: &errors.BusinessLogicError{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rewardRepo := new(MockRewardRepository)
			pointRepo := new(MockPointRepository)

			tt.setupMocks(rewardRepo, pointRepo)

			service := NewRewardService(rewardRepo, pointRepo)
			err := service.Redeem(tt.rewardID)

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

			rewardRepo.AssertExpectations(t)
			pointRepo.AssertExpectations(t)
		})
	}
}