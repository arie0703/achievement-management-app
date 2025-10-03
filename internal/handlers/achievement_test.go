package handlers

import (
	"achievement-management/internal/errors"
	"achievement-management/internal/models"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock services are already defined in server_test.go

func setupTestServer() (*Server, *MockAchievementService, *MockRewardService, *MockPointService) {
	gin.SetMode(gin.TestMode)
	
	mockAchievementService := &MockAchievementService{}
	mockRewardService := &MockRewardService{}
	mockPointService := &MockPointService{}
	
	server := NewServer(mockAchievementService, mockRewardService, mockPointService)
	
	return server, mockAchievementService, mockRewardService, mockPointService
}

func TestCreateAchievement(t *testing.T) {
	server, mockAchievementService, _, _ := setupTestServer()

	tests := []struct {
		name           string
		requestBody    interface{}
		setupMock      func()
		expectedStatus int
		expectedError  string
	}{
		{
			name: "正常な達成目録作成",
			requestBody: CreateAchievementRequest{
				Title:       "テスト達成目録",
				Description: "テスト用の達成目録です",
				Point:       100,
			},
			setupMock: func() {
				mockAchievementService.On("Create", mock.AnythingOfType("*models.Achievement")).Return(nil)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "タイトルが空の場合",
			requestBody: CreateAchievementRequest{
				Title:       "",
				Description: "テスト用の達成目録です",
				Point:       100,
			},
			setupMock:      func() {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "validation_error",
		},
		{
			name: "ポイントが0の場合",
			requestBody: CreateAchievementRequest{
				Title:       "テスト達成目録",
				Description: "テスト用の達成目録です",
				Point:       0,
			},
			setupMock:      func() {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "validation_error",
		},
		{
			name: "サービスエラーの場合",
			requestBody: CreateAchievementRequest{
				Title:       "テスト達成目録",
				Description: "テスト用の達成目録です",
				Point:       100,
			},
			setupMock: func() {
				mockAchievementService.On("Create", mock.AnythingOfType("*models.Achievement")).Return(&errors.DatabaseError{
					Operation: "Create",
					Table:     "achievements",
					Cause:     errors.ErrDatabaseOperation,
				})
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  "internal_error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// モックのセットアップ
			tt.setupMock()

			// リクエストボディの作成
			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/api/achievements", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			// レスポンスレコーダーの作成
			w := httptest.NewRecorder()

			// リクエストの実行
			server.GetRouter().ServeHTTP(w, req)

			// ステータスコードの検証
			assert.Equal(t, tt.expectedStatus, w.Code)

			// エラーレスポンスの検証
			if tt.expectedError != "" {
				var errorResponse ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedError, errorResponse.Error)
			}

			// モックの検証
			mockAchievementService.AssertExpectations(t)

			// モックのリセット
			mockAchievementService.ExpectedCalls = nil
		})
	}
}

func TestListAchievements(t *testing.T) {
	server, mockAchievementService, _, _ := setupTestServer()

	tests := []struct {
		name           string
		setupMock      func()
		expectedStatus int
		expectedCount  int
	}{
		{
			name: "正常な一覧取得",
			setupMock: func() {
				achievements := []*models.Achievement{
					{
						ID:          "test-id-1",
						Title:       "達成目録1",
						Description: "説明1",
						Point:       100,
						CreatedAt:   time.Now(),
					},
					{
						ID:          "test-id-2",
						Title:       "達成目録2",
						Description: "説明2",
						Point:       200,
						CreatedAt:   time.Now(),
					},
				}
				mockAchievementService.On("List").Return(achievements, nil)
			},
			expectedStatus: http.StatusOK,
			expectedCount:  2,
		},
		{
			name: "空の一覧",
			setupMock: func() {
				achievements := []*models.Achievement{}
				mockAchievementService.On("List").Return(achievements, nil)
			},
			expectedStatus: http.StatusOK,
			expectedCount:  0,
		},
		{
			name: "サービスエラー",
			setupMock: func() {
				mockAchievementService.On("List").Return(nil, &errors.DatabaseError{
					Operation: "List",
					Table:     "achievements",
					Cause:     errors.ErrDatabaseOperation,
				})
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// モックのセットアップ
			tt.setupMock()

			// リクエストの作成
			req := httptest.NewRequest(http.MethodGet, "/api/achievements", nil)
			w := httptest.NewRecorder()

			// リクエストの実行
			server.GetRouter().ServeHTTP(w, req)

			// ステータスコードの検証
			assert.Equal(t, tt.expectedStatus, w.Code)

			// 正常な場合のレスポンス検証
			if tt.expectedStatus == http.StatusOK {
				var response ListAchievementsResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedCount, response.Count)
				assert.Len(t, response.Achievements, tt.expectedCount)
			}

			// モックの検証
			mockAchievementService.AssertExpectations(t)

			// モックのリセット
			mockAchievementService.ExpectedCalls = nil
		})
	}
}

func TestGetAchievement(t *testing.T) {
	server, mockAchievementService, _, _ := setupTestServer()

	tests := []struct {
		name           string
		achievementID  string
		setupMock      func()
		expectedStatus int
	}{
		{
			name:          "正常な詳細取得",
			achievementID: "test-id",
			setupMock: func() {
				achievement := &models.Achievement{
					ID:          "test-id",
					Title:       "テスト達成目録",
					Description: "テスト用の達成目録です",
					Point:       100,
					CreatedAt:   time.Now(),
				}
				mockAchievementService.On("GetByID", "test-id").Return(achievement, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:          "存在しない達成目録",
			achievementID: "non-existent",
			setupMock: func() {
				mockAchievementService.On("GetByID", "non-existent").Return(nil, errors.ErrNotFound)
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:          "空のID",
			achievementID: "",
			setupMock:     func() {},
			expectedStatus: http.StatusMovedPermanently, // Ginのルーティングで301になる（/api/achievements/ -> /api/achievements）
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// モックのセットアップ
			tt.setupMock()

			// リクエストの作成
			url := "/api/achievements/" + tt.achievementID
			req := httptest.NewRequest(http.MethodGet, url, nil)
			w := httptest.NewRecorder()

			// リクエストの実行
			server.GetRouter().ServeHTTP(w, req)

			// ステータスコードの検証
			assert.Equal(t, tt.expectedStatus, w.Code)

			// 正常な場合のレスポンス検証
			if tt.expectedStatus == http.StatusOK {
				var response AchievementResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.achievementID, response.ID)
			}

			// モックの検証
			mockAchievementService.AssertExpectations(t)

			// モックのリセット
			mockAchievementService.ExpectedCalls = nil
		})
	}
}

func TestUpdateAchievement(t *testing.T) {
	server, mockAchievementService, _, _ := setupTestServer()

	tests := []struct {
		name           string
		achievementID  string
		requestBody    interface{}
		setupMock      func()
		expectedStatus int
	}{
		{
			name:          "正常な更新",
			achievementID: "test-id",
			requestBody: UpdateAchievementRequest{
				Title:       "更新されたタイトル",
				Description: "更新された説明",
				Point:       150,
			},
			setupMock: func() {
				mockAchievementService.On("Update", "test-id", mock.AnythingOfType("*models.Achievement")).Return(nil)
				updatedAchievement := &models.Achievement{
					ID:          "test-id",
					Title:       "更新されたタイトル",
					Description: "更新された説明",
					Point:       150,
					CreatedAt:   time.Now(),
				}
				mockAchievementService.On("GetByID", "test-id").Return(updatedAchievement, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:          "バリデーションエラー",
			achievementID: "test-id",
			requestBody: UpdateAchievementRequest{
				Title:       "",
				Description: "説明",
				Point:       100,
			},
			setupMock:      func() {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:          "存在しない達成目録の更新",
			achievementID: "non-existent",
			requestBody: UpdateAchievementRequest{
				Title:       "タイトル",
				Description: "説明",
				Point:       100,
			},
			setupMock: func() {
				mockAchievementService.On("Update", "non-existent", mock.AnythingOfType("*models.Achievement")).Return(errors.ErrNotFound)
			},
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// モックのセットアップ
			tt.setupMock()

			// リクエストボディの作成
			body, _ := json.Marshal(tt.requestBody)
			url := "/api/achievements/" + tt.achievementID
			req := httptest.NewRequest(http.MethodPut, url, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			// レスポンスレコーダーの作成
			w := httptest.NewRecorder()

			// リクエストの実行
			server.GetRouter().ServeHTTP(w, req)

			// ステータスコードの検証
			assert.Equal(t, tt.expectedStatus, w.Code)

			// モックの検証
			mockAchievementService.AssertExpectations(t)

			// モックのリセット
			mockAchievementService.ExpectedCalls = nil
		})
	}
}

func TestDeleteAchievement(t *testing.T) {
	server, mockAchievementService, _, _ := setupTestServer()

	tests := []struct {
		name           string
		achievementID  string
		setupMock      func()
		expectedStatus int
	}{
		{
			name:          "正常な削除",
			achievementID: "test-id",
			setupMock: func() {
				mockAchievementService.On("Delete", "test-id").Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:          "存在しない達成目録の削除",
			achievementID: "non-existent",
			setupMock: func() {
				mockAchievementService.On("Delete", "non-existent").Return(errors.ErrNotFound)
			},
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// モックのセットアップ
			tt.setupMock()

			// リクエストの作成
			url := "/api/achievements/" + tt.achievementID
			req := httptest.NewRequest(http.MethodDelete, url, nil)
			w := httptest.NewRecorder()

			// リクエストの実行
			server.GetRouter().ServeHTTP(w, req)

			// ステータスコードの検証
			assert.Equal(t, tt.expectedStatus, w.Code)

			// 正常な場合のレスポンス検証
			if tt.expectedStatus == http.StatusOK {
				var response map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "Achievement deleted successfully", response["message"])
			}

			// モックの検証
			mockAchievementService.AssertExpectations(t)

			// モックのリセット
			mockAchievementService.ExpectedCalls = nil
		})
	}
}