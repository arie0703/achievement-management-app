package handlers

import (
	"achievement-management/internal/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAchievementService モックの達成目録サービス
type MockAchievementService struct {
	mock.Mock
}

func (m *MockAchievementService) Create(achievement *models.Achievement) error {
	args := m.Called(achievement)
	return args.Error(0)
}

func (m *MockAchievementService) Update(id string, achievement *models.Achievement) error {
	args := m.Called(id, achievement)
	return args.Error(0)
}

func (m *MockAchievementService) GetByID(id string) (*models.Achievement, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Achievement), args.Error(1)
}

func (m *MockAchievementService) List() ([]*models.Achievement, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Achievement), args.Error(1)
}

func (m *MockAchievementService) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

// MockRewardService モックの報酬サービス
type MockRewardService struct {
	mock.Mock
}

func (m *MockRewardService) Create(reward *models.Reward) error {
	args := m.Called(reward)
	return args.Error(0)
}

func (m *MockRewardService) Update(id string, reward *models.Reward) error {
	args := m.Called(id, reward)
	return args.Error(0)
}

func (m *MockRewardService) GetByID(id string) (*models.Reward, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Reward), args.Error(1)
}

func (m *MockRewardService) List() ([]*models.Reward, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Reward), args.Error(1)
}

func (m *MockRewardService) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockRewardService) Redeem(rewardID string) error {
	args := m.Called(rewardID)
	return args.Error(0)
}

// MockPointService モックのポイントサービス
type MockPointService struct {
	mock.Mock
}

func (m *MockPointService) GetCurrentPoints() (*models.CurrentPoints, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.CurrentPoints), args.Error(1)
}

func (m *MockPointService) AddPoints(points int) error {
	args := m.Called(points)
	return args.Error(0)
}

func (m *MockPointService) SubtractPoints(points int) error {
	args := m.Called(points)
	return args.Error(0)
}

func (m *MockPointService) AggregatePoints() (*models.PointSummary, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.PointSummary), args.Error(1)
}

func (m *MockPointService) GetRewardHistory() ([]*models.RewardHistory, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.RewardHistory), args.Error(1)
}

func TestNewServer(t *testing.T) {
	// モックサービスを作成
	mockAchievementService := &MockAchievementService{}
	mockRewardService := &MockRewardService{}
	mockPointService := &MockPointService{}

	// サーバーを作成
	server := NewServer(mockAchievementService, mockRewardService, mockPointService)

	// サーバーが正しく初期化されていることを確認
	assert.NotNil(t, server)
	assert.NotNil(t, server.router)
	assert.Equal(t, mockAchievementService, server.achievementService)
	assert.Equal(t, mockRewardService, server.rewardService)
	assert.Equal(t, mockPointService, server.pointService)
}

func TestHealthCheck(t *testing.T) {
	// モックサービスを作成
	mockAchievementService := &MockAchievementService{}
	mockRewardService := &MockRewardService{}
	mockPointService := &MockPointService{}

	// サーバーを作成
	server := NewServer(mockAchievementService, mockRewardService, mockPointService)

	// テストリクエストを作成
	req, err := http.NewRequest("GET", "/health", nil)
	assert.NoError(t, err)

	// レスポンスレコーダーを作成
	rr := httptest.NewRecorder()

	// リクエストを実行
	server.router.ServeHTTP(rr, req)

	// レスポンスを検証
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Body.String(), "Achievement Management API is running")
}

func TestRouteSetup(t *testing.T) {
	// モックサービスを作成
	mockAchievementService := &MockAchievementService{}
	mockRewardService := &MockRewardService{}
	mockPointService := &MockPointService{}

	// サーバーを作成
	server := NewServer(mockAchievementService, mockRewardService, mockPointService)

	// 各エンドポイントが正しく設定されていることを確認
	testCases := []struct {
		method string
		path   string
		status int
	}{
		{"GET", "/health", http.StatusOK},
		// Achievement endpoints are now implemented - they will return 400/500 due to missing mock setup
		{"POST", "/api/achievements", http.StatusBadRequest},    // バリデーションエラー
		{"GET", "/api/achievements", http.StatusInternalServerError}, // モックが設定されていないためパニック
		{"GET", "/api/achievements/test-id", http.StatusInternalServerError}, // モックが設定されていないためパニック
		{"PUT", "/api/achievements/test-id", http.StatusBadRequest},    // バリデーションエラー
		{"DELETE", "/api/achievements/test-id", http.StatusInternalServerError}, // モックが設定されていないためパニック
		// Reward endpoints are now implemented - they will return 400/500 due to missing mock setup
		{"POST", "/api/rewards", http.StatusBadRequest},    // バリデーションエラー
		{"GET", "/api/rewards", http.StatusInternalServerError}, // モックが設定されていないためパニック
		{"GET", "/api/rewards/test-id", http.StatusInternalServerError}, // モックが設定されていないためパニック
		{"PUT", "/api/rewards/test-id", http.StatusBadRequest},    // バリデーションエラー
		{"DELETE", "/api/rewards/test-id", http.StatusInternalServerError}, // モックが設定されていないためパニック
		{"POST", "/api/rewards/test-id/redeem", http.StatusInternalServerError}, // モックが設定されていないためパニック
		// Point endpoints are now implemented - they will return 500 due to missing mock setup
		{"GET", "/api/points/current", http.StatusInternalServerError}, // モックが設定されていないためパニック
		{"GET", "/api/points/aggregate", http.StatusInternalServerError}, // モックが設定されていないためパニック
		{"GET", "/api/points/history", http.StatusInternalServerError}, // モックが設定されていないためパニック
	}

	for _, tc := range testCases {
		t.Run(tc.method+" "+tc.path, func(t *testing.T) {
			req, err := http.NewRequest(tc.method, tc.path, nil)
			assert.NoError(t, err)

			rr := httptest.NewRecorder()
			server.router.ServeHTTP(rr, req)

			assert.Equal(t, tc.status, rr.Code)
		})
	}
}