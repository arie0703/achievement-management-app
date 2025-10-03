package services

import (
	"testing"
	"time"

	"achievement-management/internal/errors"
	"achievement-management/internal/models"

	"github.com/stretchr/testify/assert"
)

// Mock implementations are already defined in achievement_test.go

func TestNewPointService(t *testing.T) {
	mockPointRepo := &MockPointRepository{}
	mockAchievementRepo := &MockAchievementRepository{}

	service := NewPointService(mockPointRepo, mockAchievementRepo)

	assert.NotNil(t, service)
	assert.IsType(t, &PointServiceImpl{}, service)
}

func TestPointService_GetCurrentPoints(t *testing.T) {
	tests := []struct {
		name           string
		mockSetup      func(*MockPointRepository)
		expectedResult *models.CurrentPoints
		expectedError  error
	}{
		{
			name: "正常系: 現在のポイントを取得",
			mockSetup: func(m *MockPointRepository) {
				m.On("GetCurrentPoints").Return(&models.CurrentPoints{
					ID:        "current",
					Point:     100,
					UpdatedAt: time.Now(),
				}, nil)
			},
			expectedResult: &models.CurrentPoints{
				ID:        "current",
				Point:     100,
				UpdatedAt: time.Now(),
			},
			expectedError: nil,
		},
		{
			name: "異常系: リポジトリエラー",
			mockSetup: func(m *MockPointRepository) {
				m.On("GetCurrentPoints").Return(nil, &errors.DatabaseError{
					Operation: "GetCurrentPoints",
					Table:     "current_points",
					Cause:     assert.AnError,
				})
			},
			expectedResult: nil,
			expectedError: &errors.DatabaseError{
				Operation: "GetCurrentPoints",
				Table:     "current_points",
				Cause:     assert.AnError,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockPointRepo := &MockPointRepository{}
			mockAchievementRepo := &MockAchievementRepository{}
			tt.mockSetup(mockPointRepo)

			service := NewPointService(mockPointRepo, mockAchievementRepo)
			result, err := service.GetCurrentPoints()

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.expectedResult.ID, result.ID)
				assert.Equal(t, tt.expectedResult.Point, result.Point)
			}

			mockPointRepo.AssertExpectations(t)
		})
	}
}

func TestPointService_AddPoints(t *testing.T) {
	tests := []struct {
		name          string
		points        int
		mockSetup     func(*MockPointRepository)
		expectedError error
	}{
		{
			name:   "正常系: ポイント加算",
			points: 50,
			mockSetup: func(m *MockPointRepository) {
				m.On("AddPoints", 50).Return(nil)
			},
			expectedError: nil,
		},
		{
			name:          "異常系: 負のポイント",
			points:        -10,
			mockSetup:     func(m *MockPointRepository) {},
			expectedError: &errors.ValidationError{Field: "points", Message: "points must be positive"},
		},
		{
			name:          "異常系: ゼロポイント",
			points:        0,
			mockSetup:     func(m *MockPointRepository) {},
			expectedError: &errors.ValidationError{Field: "points", Message: "points must be positive"},
		},
		{
			name:   "異常系: リポジトリエラー",
			points: 50,
			mockSetup: func(m *MockPointRepository) {
				m.On("AddPoints", 50).Return(&errors.DatabaseError{
					Operation: "AddPoints",
					Table:     "current_points",
					Cause:     assert.AnError,
				})
			},
			expectedError: &errors.DatabaseError{
				Operation: "AddPoints",
				Table:     "current_points",
				Cause:     assert.AnError,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockPointRepo := &MockPointRepository{}
			mockAchievementRepo := &MockAchievementRepository{}
			tt.mockSetup(mockPointRepo)

			service := NewPointService(mockPointRepo, mockAchievementRepo)
			err := service.AddPoints(tt.points)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			mockPointRepo.AssertExpectations(t)
		})
	}
}

func TestPointService_SubtractPoints(t *testing.T) {
	tests := []struct {
		name          string
		points        int
		mockSetup     func(*MockPointRepository)
		expectedError error
	}{
		{
			name:   "正常系: ポイント減算",
			points: 30,
			mockSetup: func(m *MockPointRepository) {
				m.On("SubtractPoints", 30).Return(nil)
			},
			expectedError: nil,
		},
		{
			name:          "異常系: 負のポイント",
			points:        -5,
			mockSetup:     func(m *MockPointRepository) {},
			expectedError: &errors.ValidationError{Field: "points", Message: "points must be positive"},
		},
		{
			name:          "異常系: ゼロポイント",
			points:        0,
			mockSetup:     func(m *MockPointRepository) {},
			expectedError: &errors.ValidationError{Field: "points", Message: "points must be positive"},
		},
		{
			name:   "異常系: ポイント不足",
			points: 100,
			mockSetup: func(m *MockPointRepository) {
				m.On("SubtractPoints", 100).Return(errors.ErrInsufficientPoints)
			},
			expectedError: errors.ErrInsufficientPoints,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockPointRepo := &MockPointRepository{}
			mockAchievementRepo := &MockAchievementRepository{}
			tt.mockSetup(mockPointRepo)

			service := NewPointService(mockPointRepo, mockAchievementRepo)
			err := service.SubtractPoints(tt.points)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			mockPointRepo.AssertExpectations(t)
		})
	}
}

func TestPointService_AggregatePoints(t *testing.T) {
	tests := []struct {
		name           string
		mockSetup      func(*MockPointRepository, *MockAchievementRepository)
		expectedResult *models.PointSummary
		expectedError  error
	}{
		{
			name: "正常系: ポイント集計（差異なし）",
			mockSetup: func(mp *MockPointRepository, ma *MockAchievementRepository) {
				achievements := []*models.Achievement{
					{ID: "1", Title: "Achievement 1", Point: 50},
					{ID: "2", Title: "Achievement 2", Point: 30},
					{ID: "3", Title: "Achievement 3", Point: 20},
				}
				ma.On("List").Return(achievements, nil)
				
				currentPoints := &models.CurrentPoints{
					ID:        "current",
					Point:     100,
					UpdatedAt: time.Now(),
				}
				mp.On("GetCurrentPoints").Return(currentPoints, nil)
			},
			expectedResult: &models.PointSummary{
				TotalAchievements: 3,
				TotalPoints:       100,
				CurrentBalance:    100,
				Difference:        0,
			},
			expectedError: nil,
		},
		{
			name: "正常系: ポイント集計（現在のポイントが少ない）",
			mockSetup: func(mp *MockPointRepository, ma *MockAchievementRepository) {
				achievements := []*models.Achievement{
					{ID: "1", Title: "Achievement 1", Point: 60},
					{ID: "2", Title: "Achievement 2", Point: 40},
				}
				ma.On("List").Return(achievements, nil)
				
				currentPoints := &models.CurrentPoints{
					ID:        "current",
					Point:     80,
					UpdatedAt: time.Now(),
				}
				mp.On("GetCurrentPoints").Return(currentPoints, nil)
			},
			expectedResult: &models.PointSummary{
				TotalAchievements: 2,
				TotalPoints:       100,
				CurrentBalance:    80,
				Difference:        20,
			},
			expectedError: nil,
		},
		{
			name: "正常系: ポイント集計（現在のポイントが多い）",
			mockSetup: func(mp *MockPointRepository, ma *MockAchievementRepository) {
				achievements := []*models.Achievement{
					{ID: "1", Title: "Achievement 1", Point: 30},
				}
				ma.On("List").Return(achievements, nil)
				
				currentPoints := &models.CurrentPoints{
					ID:        "current",
					Point:     50,
					UpdatedAt: time.Now(),
				}
				mp.On("GetCurrentPoints").Return(currentPoints, nil)
			},
			expectedResult: &models.PointSummary{
				TotalAchievements: 1,
				TotalPoints:       30,
				CurrentBalance:    50,
				Difference:        -20,
			},
			expectedError: nil,
		},
		{
			name: "正常系: 達成目録が空",
			mockSetup: func(mp *MockPointRepository, ma *MockAchievementRepository) {
				achievements := []*models.Achievement{}
				ma.On("List").Return(achievements, nil)
				
				currentPoints := &models.CurrentPoints{
					ID:        "current",
					Point:     0,
					UpdatedAt: time.Now(),
				}
				mp.On("GetCurrentPoints").Return(currentPoints, nil)
			},
			expectedResult: &models.PointSummary{
				TotalAchievements: 0,
				TotalPoints:       0,
				CurrentBalance:    0,
				Difference:        0,
			},
			expectedError: nil,
		},
		{
			name: "正常系: nilの達成目録を含む",
			mockSetup: func(mp *MockPointRepository, ma *MockAchievementRepository) {
				achievements := []*models.Achievement{
					{ID: "1", Title: "Achievement 1", Point: 25},
					nil, // nilの達成目録
					{ID: "2", Title: "Achievement 2", Point: 35},
				}
				ma.On("List").Return(achievements, nil)
				
				currentPoints := &models.CurrentPoints{
					ID:        "current",
					Point:     60,
					UpdatedAt: time.Now(),
				}
				mp.On("GetCurrentPoints").Return(currentPoints, nil)
			},
			expectedResult: &models.PointSummary{
				TotalAchievements: 3,
				TotalPoints:       60,
				CurrentBalance:    60,
				Difference:        0,
			},
			expectedError: nil,
		},
		{
			name: "異常系: 達成目録取得エラー",
			mockSetup: func(mp *MockPointRepository, ma *MockAchievementRepository) {
				ma.On("List").Return(nil, &errors.DatabaseError{
					Operation: "List",
					Table:     "achievements",
					Cause:     assert.AnError,
				})
			},
			expectedResult: nil,
			expectedError: &errors.ServiceError{
				Operation: "AggregatePoints",
				Message:   "failed to get achievements list",
				Cause: &errors.DatabaseError{
					Operation: "List",
					Table:     "achievements",
					Cause:     assert.AnError,
				},
			},
		},
		{
			name: "異常系: 現在のポイント取得エラー",
			mockSetup: func(mp *MockPointRepository, ma *MockAchievementRepository) {
				achievements := []*models.Achievement{
					{ID: "1", Title: "Achievement 1", Point: 50},
				}
				ma.On("List").Return(achievements, nil)
				
				mp.On("GetCurrentPoints").Return(nil, &errors.DatabaseError{
					Operation: "GetCurrentPoints",
					Table:     "current_points",
					Cause:     assert.AnError,
				})
			},
			expectedResult: nil,
			expectedError: &errors.ServiceError{
				Operation: "AggregatePoints",
				Message:   "failed to get current points",
				Cause: &errors.DatabaseError{
					Operation: "GetCurrentPoints",
					Table:     "current_points",
					Cause:     assert.AnError,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockPointRepo := &MockPointRepository{}
			mockAchievementRepo := &MockAchievementRepository{}
			tt.mockSetup(mockPointRepo, mockAchievementRepo)

			service := NewPointService(mockPointRepo, mockAchievementRepo)
			result, err := service.AggregatePoints()

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.expectedResult.TotalAchievements, result.TotalAchievements)
				assert.Equal(t, tt.expectedResult.TotalPoints, result.TotalPoints)
				assert.Equal(t, tt.expectedResult.CurrentBalance, result.CurrentBalance)
				assert.Equal(t, tt.expectedResult.Difference, result.Difference)
			}

			mockPointRepo.AssertExpectations(t)
			mockAchievementRepo.AssertExpectations(t)
		})
	}
}