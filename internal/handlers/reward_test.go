package handlers

import (
	"achievement-management/internal/errors"
	"achievement-management/internal/models"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Using MockRewardService from server_test.go

func TestCreateReward(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		requestBody    interface{}
		setupMock      func(*MockRewardService)
		expectedStatus int
		expectedError  string
	}{
		{
			name: "正常な報酬作成",
			requestBody: CreateRewardRequest{
				Title:       "Test Reward",
				Description: "Test Description",
				Point:       100,
			},
			setupMock: func(m *MockRewardService) {
				m.On("Create", mock.AnythingOfType("*models.Reward")).Return(nil)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "タイトル未入力エラー",
			requestBody: CreateRewardRequest{
				Description: "Test Description",
				Point:       100,
			},
			setupMock:      func(m *MockRewardService) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "validation_error",
		},
		{
			name: "ポイント未入力エラー",
			requestBody: CreateRewardRequest{
				Title:       "Test Reward",
				Description: "Test Description",
			},
			setupMock:      func(m *MockRewardService) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "validation_error",
		},
		{
			name: "ポイント負数エラー",
			requestBody: CreateRewardRequest{
				Title:       "Test Reward",
				Description: "Test Description",
				Point:       -1,
			},
			setupMock:      func(m *MockRewardService) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "validation_error",
		},
		{
			name: "サービスエラー",
			requestBody: CreateRewardRequest{
				Title:       "Test Reward",
				Description: "Test Description",
				Point:       100,
			},
			setupMock: func(m *MockRewardService) {
				m.On("Create", mock.AnythingOfType("*models.Reward")).Return(&errors.DatabaseError{
					Operation: "Create",
					Cause:     fmt.Errorf("database error"),
				})
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  "internal_error",
		},
		{
			name:           "不正なJSONエラー",
			requestBody:    "invalid json",
			setupMock:      func(m *MockRewardService) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "validation_error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// モックサービスの設定
			mockRewardService := new(MockRewardService)
			mockAchievementService := new(MockAchievementService)
			mockPointService := new(MockPointService)
			tt.setupMock(mockRewardService)

			// サーバーの作成
			server := NewServer(mockAchievementService, mockRewardService, mockPointService)

			// リクエストボディの作成
			var body []byte
			var err error
			if str, ok := tt.requestBody.(string); ok {
				body = []byte(str)
			} else {
				body, err = json.Marshal(tt.requestBody)
				assert.NoError(t, err)
			}

			// HTTPリクエストの作成
			req, err := http.NewRequest("POST", "/api/rewards", bytes.NewBuffer(body))
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			// レスポンスレコーダーの作成
			w := httptest.NewRecorder()

			// リクエストの実行
			server.GetRouter().ServeHTTP(w, req)

			// ステータスコードの確認
			assert.Equal(t, tt.expectedStatus, w.Code)

			// エラーレスポンスの確認
			if tt.expectedError != "" {
				var errorResponse ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedError, errorResponse.Error)
			}

			// 正常レスポンスの確認
			if tt.expectedStatus == http.StatusCreated {
				var response RewardResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.NotEmpty(t, response.ID)
				assert.Equal(t, "Test Reward", response.Title)
				assert.Equal(t, "Test Description", response.Description)
				assert.Equal(t, 100, response.Point)
				assert.NotZero(t, response.CreatedAt)
			}

			// モックの検証
			mockRewardService.AssertExpectations(t)
		})
	}
}

func TestListRewards(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		setupMock      func(*MockRewardService)
		expectedStatus int
		expectedCount  int
		expectedError  string
	}{
		{
			name: "正常な報酬一覧取得",
			setupMock: func(m *MockRewardService) {
				rewards := []*models.Reward{
					{
						ID:          "reward1",
						Title:       "Reward 1",
						Description: "Description 1",
						Point:       100,
						CreatedAt:   time.Now(),
					},
					{
						ID:          "reward2",
						Title:       "Reward 2",
						Description: "Description 2",
						Point:       200,
						CreatedAt:   time.Now(),
					},
				}
				m.On("List").Return(rewards, nil)
			},
			expectedStatus: http.StatusOK,
			expectedCount:  2,
		},
		{
			name: "空の報酬一覧",
			setupMock: func(m *MockRewardService) {
				m.On("List").Return([]*models.Reward{}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedCount:  0,
		},
		{
			name: "サービスエラー",
			setupMock: func(m *MockRewardService) {
				m.On("List").Return(nil, &errors.DatabaseError{
					Operation: "List",
					Cause:     fmt.Errorf("database error"),
				})
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  "internal_error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// モックサービスの設定
			mockRewardService := new(MockRewardService)
			mockAchievementService := new(MockAchievementService)
			mockPointService := new(MockPointService)
			tt.setupMock(mockRewardService)

			// サーバーの作成
			server := NewServer(mockAchievementService, mockRewardService, mockPointService)

			// HTTPリクエストの作成
			req, err := http.NewRequest("GET", "/api/rewards", nil)
			assert.NoError(t, err)

			// レスポンスレコーダーの作成
			w := httptest.NewRecorder()

			// リクエストの実行
			server.GetRouter().ServeHTTP(w, req)

			// ステータスコードの確認
			assert.Equal(t, tt.expectedStatus, w.Code)

			// エラーレスポンスの確認
			if tt.expectedError != "" {
				var errorResponse ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedError, errorResponse.Error)
			} else {
				// 正常レスポンスの確認
				var response ListRewardsResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedCount, response.Count)
				assert.Equal(t, tt.expectedCount, len(response.Rewards))
			}

			// モックの検証
			mockRewardService.AssertExpectations(t)
		})
	}
}

func TestGetReward(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		rewardID       string
		setupMock      func(*MockRewardService)
		expectedStatus int
		expectedError  string
	}{
		{
			name:     "正常な報酬取得",
			rewardID: "reward1",
			setupMock: func(m *MockRewardService) {
				reward := &models.Reward{
					ID:          "reward1",
					Title:       "Test Reward",
					Description: "Test Description",
					Point:       100,
					CreatedAt:   time.Now(),
				}
				m.On("GetByID", "reward1").Return(reward, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:     "存在しない報酬",
			rewardID: "nonexistent",
			setupMock: func(m *MockRewardService) {
				m.On("GetByID", "nonexistent").Return(nil, fmt.Errorf("resource not found"))
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  "not_found",
		},
		{
			name:     "サービスエラー",
			rewardID: "reward1",
			setupMock: func(m *MockRewardService) {
				m.On("GetByID", "reward1").Return(nil, &errors.DatabaseError{
					Operation: "GetByID",
					Cause:     fmt.Errorf("database error"),
				})
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  "internal_error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// モックサービスの設定
			mockRewardService := new(MockRewardService)
			mockAchievementService := new(MockAchievementService)
			mockPointService := new(MockPointService)
			tt.setupMock(mockRewardService)

			// サーバーの作成
			server := NewServer(mockAchievementService, mockRewardService, mockPointService)

			// HTTPリクエストの作成
			url := "/api/rewards/" + tt.rewardID
			req, err := http.NewRequest("GET", url, nil)
			assert.NoError(t, err)

			// レスポンスレコーダーの作成
			w := httptest.NewRecorder()

			// リクエストの実行
			server.GetRouter().ServeHTTP(w, req)

			// ステータスコードの確認
			assert.Equal(t, tt.expectedStatus, w.Code)

			// エラーレスポンスの確認
			if tt.expectedError != "" {
				var errorResponse ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedError, errorResponse.Error)
			}

			// 正常レスポンスの確認
			if tt.expectedStatus == http.StatusOK {
				var response RewardResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "reward1", response.ID)
				assert.Equal(t, "Test Reward", response.Title)
				assert.Equal(t, "Test Description", response.Description)
				assert.Equal(t, 100, response.Point)
				assert.NotZero(t, response.CreatedAt)
			}

			// モックの検証
			mockRewardService.AssertExpectations(t)
		})
	}
}

func TestUpdateReward(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		rewardID       string
		requestBody    interface{}
		setupMock      func(*MockRewardService)
		expectedStatus int
		expectedError  string
	}{
		{
			name:     "正常な報酬更新",
			rewardID: "reward1",
			requestBody: UpdateRewardRequest{
				Title:       "Updated Reward",
				Description: "Updated Description",
				Point:       150,
			},
			setupMock: func(m *MockRewardService) {
				m.On("Update", "reward1", mock.AnythingOfType("*models.Reward")).Return(nil)
				updatedReward := &models.Reward{
					ID:          "reward1",
					Title:       "Updated Reward",
					Description: "Updated Description",
					Point:       150,
					CreatedAt:   time.Now(),
				}
				m.On("GetByID", "reward1").Return(updatedReward, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:     "タイトル未入力エラー",
			rewardID: "reward1",
			requestBody: UpdateRewardRequest{
				Description: "Updated Description",
				Point:       150,
			},
			setupMock:      func(m *MockRewardService) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "validation_error",
		},
		{
			name:     "存在しない報酬更新",
			rewardID: "nonexistent",
			requestBody: UpdateRewardRequest{
				Title:       "Updated Reward",
				Description: "Updated Description",
				Point:       150,
			},
			setupMock: func(m *MockRewardService) {
				m.On("Update", "nonexistent", mock.AnythingOfType("*models.Reward")).Return(fmt.Errorf("resource not found"))
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  "not_found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// モックサービスの設定
			mockRewardService := new(MockRewardService)
			mockAchievementService := new(MockAchievementService)
			mockPointService := new(MockPointService)
			tt.setupMock(mockRewardService)

			// サーバーの作成
			server := NewServer(mockAchievementService, mockRewardService, mockPointService)

			// リクエストボディの作成
			body, err := json.Marshal(tt.requestBody)
			assert.NoError(t, err)

			// HTTPリクエストの作成
			url := "/api/rewards/" + tt.rewardID
			req, err := http.NewRequest("PUT", url, bytes.NewBuffer(body))
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			// レスポンスレコーダーの作成
			w := httptest.NewRecorder()

			// リクエストの実行
			server.GetRouter().ServeHTTP(w, req)

			// ステータスコードの確認
			assert.Equal(t, tt.expectedStatus, w.Code)

			// エラーレスポンスの確認
			if tt.expectedError != "" {
				var errorResponse ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedError, errorResponse.Error)
			}

			// 正常レスポンスの確認
			if tt.expectedStatus == http.StatusOK {
				var response RewardResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "reward1", response.ID)
				assert.Equal(t, "Updated Reward", response.Title)
				assert.Equal(t, "Updated Description", response.Description)
				assert.Equal(t, 150, response.Point)
				assert.NotZero(t, response.CreatedAt)
			}

			// モックの検証
			mockRewardService.AssertExpectations(t)
		})
	}
}

func TestDeleteReward(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		rewardID       string
		setupMock      func(*MockRewardService)
		expectedStatus int
		expectedError  string
	}{
		{
			name:     "正常な報酬削除",
			rewardID: "reward1",
			setupMock: func(m *MockRewardService) {
				m.On("Delete", "reward1").Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:     "存在しない報酬削除",
			rewardID: "nonexistent",
			setupMock: func(m *MockRewardService) {
				m.On("Delete", "nonexistent").Return(fmt.Errorf("resource not found"))
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  "not_found",
		},
		{
			name:     "サービスエラー",
			rewardID: "reward1",
			setupMock: func(m *MockRewardService) {
				m.On("Delete", "reward1").Return(&errors.DatabaseError{
					Operation: "Delete",
					Cause:     fmt.Errorf("database error"),
				})
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  "internal_error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// モックサービスの設定
			mockRewardService := new(MockRewardService)
			mockAchievementService := new(MockAchievementService)
			mockPointService := new(MockPointService)
			tt.setupMock(mockRewardService)

			// サーバーの作成
			server := NewServer(mockAchievementService, mockRewardService, mockPointService)

			// HTTPリクエストの作成
			url := "/api/rewards/" + tt.rewardID
			req, err := http.NewRequest("DELETE", url, nil)
			assert.NoError(t, err)

			// レスポンスレコーダーの作成
			w := httptest.NewRecorder()

			// リクエストの実行
			server.GetRouter().ServeHTTP(w, req)

			// ステータスコードの確認
			assert.Equal(t, tt.expectedStatus, w.Code)

			// エラーレスポンスの確認
			if tt.expectedError != "" {
				var errorResponse ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedError, errorResponse.Error)
			}

			// 正常レスポンスの確認
			if tt.expectedStatus == http.StatusOK {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "Reward deleted successfully", response["message"])
			}

			// モックの検証
			mockRewardService.AssertExpectations(t)
		})
	}
}

func TestRedeemReward(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		rewardID       string
		setupMock      func(*MockRewardService)
		expectedStatus int
		expectedError  string
	}{
		{
			name:     "正常な報酬獲得",
			rewardID: "reward1",
			setupMock: func(m *MockRewardService) {
				m.On("Redeem", "reward1").Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:     "存在しない報酬獲得",
			rewardID: "nonexistent",
			setupMock: func(m *MockRewardService) {
				m.On("Redeem", "nonexistent").Return(fmt.Errorf("resource not found"))
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  "not_found",
		},
		{
			name:     "ポイント不足エラー",
			rewardID: "reward1",
			setupMock: func(m *MockRewardService) {
				m.On("Redeem", "reward1").Return(&errors.BusinessLogicError{
					Operation: "Redeem",
					Reason:    "insufficient points",
				})
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "business_logic_error",
		},
		{
			name:     "サービスエラー",
			rewardID: "reward1",
			setupMock: func(m *MockRewardService) {
				m.On("Redeem", "reward1").Return(&errors.DatabaseError{
					Operation: "Redeem",
					Cause:     fmt.Errorf("database error"),
				})
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  "internal_error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// モックサービスの設定
			mockRewardService := new(MockRewardService)
			mockAchievementService := new(MockAchievementService)
			mockPointService := new(MockPointService)
			tt.setupMock(mockRewardService)

			// サーバーの作成
			server := NewServer(mockAchievementService, mockRewardService, mockPointService)

			// HTTPリクエストの作成
			url := "/api/rewards/" + tt.rewardID + "/redeem"
			req, err := http.NewRequest("POST", url, nil)
			assert.NoError(t, err)

			// レスポンスレコーダーの作成
			w := httptest.NewRecorder()

			// リクエストの実行
			server.GetRouter().ServeHTTP(w, req)

			// ステータスコードの確認
			assert.Equal(t, tt.expectedStatus, w.Code)

			// エラーレスポンスの確認
			if tt.expectedError != "" {
				var errorResponse ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedError, errorResponse.Error)
			}

			// 正常レスポンスの確認
			if tt.expectedStatus == http.StatusOK {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "Reward redeemed successfully", response["message"])
			}

			// モックの検証
			mockRewardService.AssertExpectations(t)
		})
	}
}