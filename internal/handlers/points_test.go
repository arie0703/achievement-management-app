package handlers

import (
	"achievement-management/internal/errors"
	"achievement-management/internal/models"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetCurrentPoints_Success(t *testing.T) {
	// モックサービスを作成
	mockAchievementService := &MockAchievementService{}
	mockRewardService := &MockRewardService{}
	mockPointService := &MockPointService{}

	// モックの期待値を設定
	expectedPoints := &models.CurrentPoints{
		ID:        "current",
		Point:     150,
		UpdatedAt: time.Now(),
	}
	mockPointService.On("GetCurrentPoints").Return(expectedPoints, nil)

	// サーバーを作成
	server := NewServer(mockAchievementService, mockRewardService, mockPointService)

	// テストリクエストを作成
	req, err := http.NewRequest("GET", "/api/points/current", nil)
	assert.NoError(t, err)

	// レスポンスレコーダーを作成
	rr := httptest.NewRecorder()

	// リクエストを実行
	server.router.ServeHTTP(rr, req)

	// レスポンスを検証
	assert.Equal(t, http.StatusOK, rr.Code)

	var response CurrentPointsResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedPoints.ID, response.ID)
	assert.Equal(t, expectedPoints.Point, response.Point)
	assert.Equal(t, expectedPoints.UpdatedAt.Unix(), response.UpdatedAt.Unix())

	// モックが呼ばれたことを確認
	mockPointService.AssertExpectations(t)
}

func TestGetCurrentPoints_ServiceError(t *testing.T) {
	// モックサービスを作成
	mockAchievementService := &MockAchievementService{}
	mockRewardService := &MockRewardService{}
	mockPointService := &MockPointService{}

	// モックの期待値を設定（エラーを返す）
	mockPointService.On("GetCurrentPoints").Return(nil, &errors.DatabaseError{
		Operation: "get_current_points",
		Cause:     fmt.Errorf("database connection failed"),
	})

	// サーバーを作成
	server := NewServer(mockAchievementService, mockRewardService, mockPointService)

	// テストリクエストを作成
	req, err := http.NewRequest("GET", "/api/points/current", nil)
	assert.NoError(t, err)

	// レスポンスレコーダーを作成
	rr := httptest.NewRecorder()

	// リクエストを実行
	server.router.ServeHTTP(rr, req)

	// レスポンスを検証
	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	var response ErrorResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "internal_error", response.Error)
	assert.Equal(t, "Internal server error", response.Message)
	assert.Equal(t, 500, response.Code)

	// モックが呼ばれたことを確認
	mockPointService.AssertExpectations(t)
}

func TestAggregatePoints_Success(t *testing.T) {
	// モックサービスを作成
	mockAchievementService := &MockAchievementService{}
	mockRewardService := &MockRewardService{}
	mockPointService := &MockPointService{}

	// モックの期待値を設定
	expectedSummary := &models.PointSummary{
		TotalAchievements: 5,
		TotalPoints:       500,
		CurrentBalance:    150,
		Difference:        -350,
	}
	mockPointService.On("AggregatePoints").Return(expectedSummary, nil)

	// サーバーを作成
	server := NewServer(mockAchievementService, mockRewardService, mockPointService)

	// テストリクエストを作成
	req, err := http.NewRequest("GET", "/api/points/aggregate", nil)
	assert.NoError(t, err)

	// レスポンスレコーダーを作成
	rr := httptest.NewRecorder()

	// リクエストを実行
	server.router.ServeHTTP(rr, req)

	// レスポンスを検証
	assert.Equal(t, http.StatusOK, rr.Code)

	var response PointSummaryResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedSummary.TotalAchievements, response.TotalAchievements)
	assert.Equal(t, expectedSummary.TotalPoints, response.TotalPoints)
	assert.Equal(t, expectedSummary.CurrentBalance, response.CurrentBalance)
	assert.Equal(t, expectedSummary.Difference, response.Difference)

	// モックが呼ばれたことを確認
	mockPointService.AssertExpectations(t)
}

func TestAggregatePoints_ServiceError(t *testing.T) {
	// モックサービスを作成
	mockAchievementService := &MockAchievementService{}
	mockRewardService := &MockRewardService{}
	mockPointService := &MockPointService{}

	// モックの期待値を設定（エラーを返す）
	mockPointService.On("AggregatePoints").Return(nil, &errors.BusinessLogicError{
		Operation: "aggregate_points",
		Reason:    "aggregation failed",
	})

	// サーバーを作成
	server := NewServer(mockAchievementService, mockRewardService, mockPointService)

	// テストリクエストを作成
	req, err := http.NewRequest("GET", "/api/points/aggregate", nil)
	assert.NoError(t, err)

	// レスポンスレコーダーを作成
	rr := httptest.NewRecorder()

	// リクエストを実行
	server.router.ServeHTTP(rr, req)

	// レスポンスを検証
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	var response ErrorResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "business_logic_error", response.Error)
	assert.Equal(t, "business logic error in operation 'aggregate_points': aggregation failed", response.Message)
	assert.Equal(t, 400, response.Code)

	// モックが呼ばれたことを確認
	mockPointService.AssertExpectations(t)
}

func TestGetPointsHistory_Success(t *testing.T) {
	// モックサービスを作成
	mockAchievementService := &MockAchievementService{}
	mockRewardService := &MockRewardService{}
	mockPointService := &MockPointService{}

	// モックの期待値を設定
	now := time.Now()
	expectedHistory := []*models.RewardHistory{
		{
			ID:          "history-1",
			RewardID:    "reward-1",
			RewardTitle: "Test Reward 1",
			PointCost:   50,
			RedeemedAt:  now,
		},
		{
			ID:          "history-2",
			RewardID:    "reward-2",
			RewardTitle: "Test Reward 2",
			PointCost:   100,
			RedeemedAt:  now.Add(-time.Hour),
		},
	}
	mockPointService.On("GetRewardHistory").Return(expectedHistory, nil)

	// サーバーを作成
	server := NewServer(mockAchievementService, mockRewardService, mockPointService)

	// テストリクエストを作成
	req, err := http.NewRequest("GET", "/api/points/history", nil)
	assert.NoError(t, err)

	// レスポンスレコーダーを作成
	rr := httptest.NewRecorder()

	// リクエストを実行
	server.router.ServeHTTP(rr, req)

	// レスポンスを検証
	assert.Equal(t, http.StatusOK, rr.Code)

	var response ListRewardHistoryResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 2, response.Count)
	assert.Len(t, response.History, 2)

	// 最初の履歴項目を検証
	assert.Equal(t, expectedHistory[0].ID, response.History[0].ID)
	assert.Equal(t, expectedHistory[0].RewardID, response.History[0].RewardID)
	assert.Equal(t, expectedHistory[0].RewardTitle, response.History[0].RewardTitle)
	assert.Equal(t, expectedHistory[0].PointCost, response.History[0].PointCost)
	assert.Equal(t, expectedHistory[0].RedeemedAt.Unix(), response.History[0].RedeemedAt.Unix())

	// 2番目の履歴項目を検証
	assert.Equal(t, expectedHistory[1].ID, response.History[1].ID)
	assert.Equal(t, expectedHistory[1].RewardID, response.History[1].RewardID)
	assert.Equal(t, expectedHistory[1].RewardTitle, response.History[1].RewardTitle)
	assert.Equal(t, expectedHistory[1].PointCost, response.History[1].PointCost)
	assert.Equal(t, expectedHistory[1].RedeemedAt.Unix(), response.History[1].RedeemedAt.Unix())

	// モックが呼ばれたことを確認
	mockPointService.AssertExpectations(t)
}

func TestGetPointsHistory_EmptyHistory(t *testing.T) {
	// モックサービスを作成
	mockAchievementService := &MockAchievementService{}
	mockRewardService := &MockRewardService{}
	mockPointService := &MockPointService{}

	// モックの期待値を設定（空の履歴）
	expectedHistory := []*models.RewardHistory{}
	mockPointService.On("GetRewardHistory").Return(expectedHistory, nil)

	// サーバーを作成
	server := NewServer(mockAchievementService, mockRewardService, mockPointService)

	// テストリクエストを作成
	req, err := http.NewRequest("GET", "/api/points/history", nil)
	assert.NoError(t, err)

	// レスポンスレコーダーを作成
	rr := httptest.NewRecorder()

	// リクエストを実行
	server.router.ServeHTTP(rr, req)

	// レスポンスを検証
	assert.Equal(t, http.StatusOK, rr.Code)

	var response ListRewardHistoryResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 0, response.Count)
	assert.Len(t, response.History, 0)

	// モックが呼ばれたことを確認
	mockPointService.AssertExpectations(t)
}

func TestGetPointsHistory_ServiceError(t *testing.T) {
	// モックサービスを作成
	mockAchievementService := &MockAchievementService{}
	mockRewardService := &MockRewardService{}
	mockPointService := &MockPointService{}

	// モックの期待値を設定（エラーを返す）
	mockPointService.On("GetRewardHistory").Return(nil, &errors.DatabaseError{
		Operation: "get_reward_history",
		Cause:     fmt.Errorf("table not found"),
	})

	// サーバーを作成
	server := NewServer(mockAchievementService, mockRewardService, mockPointService)

	// テストリクエストを作成
	req, err := http.NewRequest("GET", "/api/points/history", nil)
	assert.NoError(t, err)

	// レスポンスレコーダーを作成
	rr := httptest.NewRecorder()

	// リクエストを実行
	server.router.ServeHTTP(rr, req)

	// レスポンスを検証
	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	var response ErrorResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "internal_error", response.Error)
	assert.Equal(t, "Internal server error", response.Message)
	assert.Equal(t, 500, response.Code)

	// モックが呼ばれたことを確認
	mockPointService.AssertExpectations(t)
}