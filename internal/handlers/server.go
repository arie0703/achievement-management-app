package handlers

import (
	"achievement-management/internal/config"
	"achievement-management/internal/errors"
	"achievement-management/internal/logging"
	"achievement-management/internal/models"
	"achievement-management/internal/services"
	"crypto/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/oklog/ulid/v2"
)

// Server HTTPサーバー
type Server struct {
	achievementService services.AchievementService
	rewardService      services.RewardService
	pointService       services.PointService
	router             *gin.Engine
	logger             logging.Logger
	accessLogger       *logging.AccessLogger
	errorLogger        *logging.ErrorLogger
}

// NewServer 新しいサーバーインスタンスを作成
func NewServer(
	achievementService services.AchievementService,
	rewardService services.RewardService,
	pointService services.PointService,
	config *config.Config,
) *Server {
	// ログ設定に基づいてGinのモードを設定
	if config.Logging.Level == "debug" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// ロガーを初期化
	logger, err := logging.NewLogger(config)
	if err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}

	accessLogger, err := logging.NewAccessLogger(config)
	if err != nil {
		panic("Failed to initialize access logger: " + err.Error())
	}

	errorLogger, err := logging.NewErrorLogger(config)
	if err != nil {
		panic("Failed to initialize error logger: " + err.Error())
	}

	router := gin.New()

	server := &Server{
		achievementService: achievementService,
		rewardService:      rewardService,
		pointService:       pointService,
		router:             router,
		logger:             logger,
		accessLogger:       accessLogger,
		errorLogger:        errorLogger,
	}

	// ミドルウェアの設定
	router.Use(logging.LoggingMiddleware(accessLogger))
	router.Use(logging.ErrorLoggingMiddleware(errorLogger))
	router.Use(logging.RecoveryMiddleware(errorLogger))
	router.Use(server.CORSMiddleware())

	// ルートの設定
	server.setupRoutes()

	return server
}

// setupRoutes ルートの設定
func (s *Server) setupRoutes() {
	// ヘルスチェックエンドポイント
	s.router.GET("/health", s.healthCheck)

	// APIルートグループ
	api := s.router.Group("/api")
	{
		// 達成目録エンドポイント（後で実装）
		achievements := api.Group("/achievements")
		{
			achievements.POST("", s.createAchievement)
			achievements.GET("", s.listAchievements)
			achievements.GET("/:id", s.getAchievement)
			achievements.PUT("/:id", s.updateAchievement)
			achievements.DELETE("/:id", s.deleteAchievement)
		}

		// 報酬エンドポイント（後で実装）
		rewards := api.Group("/rewards")
		{
			rewards.POST("", s.createReward)
			rewards.GET("", s.listRewards)
			rewards.GET("/:id", s.getReward)
			rewards.PUT("/:id", s.updateReward)
			rewards.DELETE("/:id", s.deleteReward)
			rewards.POST("/:id/redeem", s.redeemReward)
		}

		// ポイント管理エンドポイント（後で実装）
		points := api.Group("/points")
		{
			points.GET("/current", s.getCurrentPoints)
			points.GET("/aggregate", s.aggregatePoints)
			points.GET("/history", s.getPointsHistory)
		}
	}
}

// healthCheck ヘルスチェックハンドラー
func (s *Server) healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "Achievement Management API is running",
	})
}

// Run サーバーを起動
func (s *Server) Run(addr string) error {
	return s.router.Run(addr)
}

// GetRouter ルーターを取得（テスト用）
func (s *Server) GetRouter() *gin.Engine {
	return s.router
}

// Achievement API endpoints implementation

// createAchievement POST /api/achievements - 達成目録作成
func (s *Server) createAchievement(c *gin.Context) {
	s.logger.WithField("endpoint", "create_achievement").Debug("Processing create achievement request")

	var req CreateAchievementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		s.errorLogger.LogAPIError("/api/achievements", "POST", 400, err)
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: "Invalid request body: " + err.Error(),
			Code:    400,
		})
		return
	}

	achievement := req.ToModel()
	if err := s.achievementService.Create(achievement); err != nil {
		s.errorLogger.LogServiceError("achievement", "create", err)
		handleServiceError(c, err)
		return
	}

	s.logger.WithFields(map[string]interface{}{
		"achievement_id": achievement.ID,
		"title":          achievement.Title,
		"point":          achievement.Point,
	}).Info("Achievement created successfully")

	c.JSON(http.StatusCreated, AchievementResponse{
		ID:          achievement.ID,
		Title:       achievement.Title,
		Description: achievement.Description,
		Point:       achievement.Point,
		CreatedAt:   achievement.CreatedAt,
	})
}

// listAchievements GET /api/achievements - 達成目録一覧取得
func (s *Server) listAchievements(c *gin.Context) {
	achievements, err := s.achievementService.List()
	if err != nil {
		handleServiceError(c, err)
		return
	}

	response := make([]AchievementResponse, len(achievements))
	for i, achievement := range achievements {
		response[i] = AchievementResponse{
			ID:          achievement.ID,
			Title:       achievement.Title,
			Description: achievement.Description,
			Point:       achievement.Point,
			CreatedAt:   achievement.CreatedAt,
		}
	}

	c.JSON(http.StatusOK, ListAchievementsResponse{
		Achievements: response,
		Count:        len(response),
	})
}

// getAchievement GET /api/achievements/{id} - 達成目録詳細取得
func (s *Server) getAchievement(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: "Achievement ID is required",
			Code:    400,
		})
		return
	}

	achievement, err := s.achievementService.GetByID(id)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, AchievementResponse{
		ID:          achievement.ID,
		Title:       achievement.Title,
		Description: achievement.Description,
		Point:       achievement.Point,
		CreatedAt:   achievement.CreatedAt,
	})
}

// updateAchievement PUT /api/achievements/{id} - 達成目録更新
func (s *Server) updateAchievement(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: "Achievement ID is required",
			Code:    400,
		})
		return
	}

	var req UpdateAchievementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: "Invalid request body: " + err.Error(),
			Code:    400,
		})
		return
	}

	achievement := req.ToModel()
	if err := s.achievementService.Update(id, achievement); err != nil {
		handleServiceError(c, err)
		return
	}

	// 更新後のデータを取得して返す
	updatedAchievement, err := s.achievementService.GetByID(id)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, AchievementResponse{
		ID:          updatedAchievement.ID,
		Title:       updatedAchievement.Title,
		Description: updatedAchievement.Description,
		Point:       updatedAchievement.Point,
		CreatedAt:   updatedAchievement.CreatedAt,
	})
}

// deleteAchievement DELETE /api/achievements/{id} - 達成目録削除
func (s *Server) deleteAchievement(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: "Achievement ID is required",
			Code:    400,
		})
		return
	}

	if err := s.achievementService.Delete(id); err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Achievement deleted successfully",
	})
}

// Reward API endpoints implementation

// createReward POST /api/rewards - 報酬作成
func (s *Server) createReward(c *gin.Context) {
	var req CreateRewardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: "Invalid request body: " + err.Error(),
			Code:    400,
		})
		return
	}

	reward := req.ToModel()
	if err := s.rewardService.Create(reward); err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusCreated, RewardResponse{
		ID:          reward.ID,
		Title:       reward.Title,
		Description: reward.Description,
		Point:       reward.Point,
		CreatedAt:   reward.CreatedAt,
	})
}

// listRewards GET /api/rewards - 報酬一覧取得
func (s *Server) listRewards(c *gin.Context) {
	rewards, err := s.rewardService.List()
	if err != nil {
		handleServiceError(c, err)
		return
	}

	response := make([]RewardResponse, len(rewards))
	for i, reward := range rewards {
		response[i] = RewardResponse{
			ID:          reward.ID,
			Title:       reward.Title,
			Description: reward.Description,
			Point:       reward.Point,
			CreatedAt:   reward.CreatedAt,
		}
	}

	c.JSON(http.StatusOK, ListRewardsResponse{
		Rewards: response,
		Count:   len(response),
	})
}

// getReward GET /api/rewards/{id} - 報酬詳細取得
func (s *Server) getReward(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: "Reward ID is required",
			Code:    400,
		})
		return
	}

	reward, err := s.rewardService.GetByID(id)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, RewardResponse{
		ID:          reward.ID,
		Title:       reward.Title,
		Description: reward.Description,
		Point:       reward.Point,
		CreatedAt:   reward.CreatedAt,
	})
}

// updateReward PUT /api/rewards/{id} - 報酬更新
func (s *Server) updateReward(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: "Reward ID is required",
			Code:    400,
		})
		return
	}

	var req UpdateRewardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: "Invalid request body: " + err.Error(),
			Code:    400,
		})
		return
	}

	reward := req.ToModel()
	if err := s.rewardService.Update(id, reward); err != nil {
		handleServiceError(c, err)
		return
	}

	// 更新後のデータを取得して返す
	updatedReward, err := s.rewardService.GetByID(id)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, RewardResponse{
		ID:          updatedReward.ID,
		Title:       updatedReward.Title,
		Description: updatedReward.Description,
		Point:       updatedReward.Point,
		CreatedAt:   updatedReward.CreatedAt,
	})
}

// deleteReward DELETE /api/rewards/{id} - 報酬削除
func (s *Server) deleteReward(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: "Reward ID is required",
			Code:    400,
		})
		return
	}

	if err := s.rewardService.Delete(id); err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Reward deleted successfully",
	})
}

// redeemReward POST /api/rewards/{id}/redeem - 報酬獲得
func (s *Server) redeemReward(c *gin.Context) {
	s.logger.WithField("endpoint", "redeem_reward").Debug("Processing reward redemption request")

	id := c.Param("id")
	if id == "" {
		s.errorLogger.LogAPIError("/api/rewards/{id}/redeem", "POST", 400,
			&ValidationError{Message: "Reward ID is required"})
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: "Reward ID is required",
			Code:    400,
		})
		return
	}

	if err := s.rewardService.Redeem(id); err != nil {
		s.errorLogger.LogServiceError("reward", "redeem", err)
		handleServiceError(c, err)
		return
	}

	s.logger.WithField("reward_id", id).Info("Reward redeemed successfully")

	c.JSON(http.StatusOK, gin.H{
		"message": "Reward redeemed successfully",
	})
}

// getCurrentPoints GET /api/points/current - 現在のポイント取得
func (s *Server) getCurrentPoints(c *gin.Context) {
	currentPoints, err := s.pointService.GetCurrentPoints()
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, CurrentPointsResponse{
		ID:        currentPoints.ID,
		Point:     currentPoints.Point,
		UpdatedAt: currentPoints.UpdatedAt,
	})
}

// aggregatePoints GET /api/points/aggregate - ポイント集計
func (s *Server) aggregatePoints(c *gin.Context) {
	summary, err := s.pointService.AggregatePoints()
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, PointSummaryResponse{
		TotalAchievements: summary.TotalAchievements,
		TotalPoints:       summary.TotalPoints,
		CurrentBalance:    summary.CurrentBalance,
		Difference:        summary.Difference,
	})
}

// getPointsHistory GET /api/points/history - 報酬獲得履歴取得
func (s *Server) getPointsHistory(c *gin.Context) {
	history, err := s.pointService.GetRewardHistory()
	if err != nil {
		handleServiceError(c, err)
		return
	}

	response := make([]RewardHistoryResponse, len(history))
	for i, record := range history {
		response[i] = RewardHistoryResponse{
			ID:          record.ID,
			RewardID:    record.RewardID,
			RewardTitle: record.RewardTitle,
			PointCost:   record.PointCost,
			RedeemedAt:  record.RedeemedAt,
		}
	}

	c.JSON(http.StatusOK, ListRewardHistoryResponse{
		History: response,
		Count:   len(response),
	})
}

// Achievement API request/response types

// CreateAchievementRequest 達成目録作成リクエスト
type CreateAchievementRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Point       int    `json:"point" binding:"required,min=1"`
}

// ToModel リクエストをモデルに変換
func (r *CreateAchievementRequest) ToModel() *models.Achievement {
	return &models.Achievement{
		ID:          ulid.MustNew(ulid.Timestamp(time.Now()), rand.Reader).String(),
		Title:       r.Title,
		Description: r.Description,
		Point:       r.Point,
		CreatedAt:   time.Now(),
	}
}

// UpdateAchievementRequest 達成目録更新リクエスト
type UpdateAchievementRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Point       int    `json:"point" binding:"required,min=1"`
}

// ToModel リクエストをモデルに変換
func (r *UpdateAchievementRequest) ToModel() *models.Achievement {
	return &models.Achievement{
		Title:       r.Title,
		Description: r.Description,
		Point:       r.Point,
	}
}

// AchievementResponse 達成目録レスポンス
type AchievementResponse struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Point       int       `json:"point"`
	CreatedAt   time.Time `json:"created_at"`
}

// ListAchievementsResponse 達成目録一覧レスポンス
type ListAchievementsResponse struct {
	Achievements []AchievementResponse `json:"achievements"`
	Count        int                   `json:"count"`
}

// Reward API request/response types

// CreateRewardRequest 報酬作成リクエスト
type CreateRewardRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Point       int    `json:"point" binding:"required,min=1"`
}

// ToModel リクエストをモデルに変換
func (r *CreateRewardRequest) ToModel() *models.Reward {
	return &models.Reward{
		ID:          ulid.MustNew(ulid.Timestamp(time.Now()), rand.Reader).String(),
		Title:       r.Title,
		Description: r.Description,
		Point:       r.Point,
		CreatedAt:   time.Now(),
	}
}

// UpdateRewardRequest 報酬更新リクエスト
type UpdateRewardRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Point       int    `json:"point" binding:"required,min=1"`
}

// ToModel リクエストをモデルに変換
func (r *UpdateRewardRequest) ToModel() *models.Reward {
	return &models.Reward{
		Title:       r.Title,
		Description: r.Description,
		Point:       r.Point,
	}
}

// RewardResponse 報酬レスポンス
type RewardResponse struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Point       int       `json:"point"`
	CreatedAt   time.Time `json:"created_at"`
}

// ListRewardsResponse 報酬一覧レスポンス
type ListRewardsResponse struct {
	Rewards []RewardResponse `json:"rewards"`
	Count   int              `json:"count"`
}

// Points API response types

// CurrentPointsResponse 現在のポイントレスポンス
type CurrentPointsResponse struct {
	ID        string    `json:"id"`
	Point     int       `json:"point"`
	UpdatedAt time.Time `json:"updated_at"`
}

// PointSummaryResponse ポイント集計レスポンス
type PointSummaryResponse struct {
	TotalAchievements int `json:"total_achievements"`
	TotalPoints       int `json:"total_points"`
	CurrentBalance    int `json:"current_balance"`
	Difference        int `json:"difference"`
}

// RewardHistoryResponse 報酬獲得履歴レスポンス
type RewardHistoryResponse struct {
	ID          string    `json:"id"`
	RewardID    string    `json:"reward_id"`
	RewardTitle string    `json:"reward_title"`
	PointCost   int       `json:"point_cost"`
	RedeemedAt  time.Time `json:"redeemed_at"`
}

// ListRewardHistoryResponse 報酬獲得履歴一覧レスポンス
type ListRewardHistoryResponse struct {
	History []RewardHistoryResponse `json:"history"`
	Count   int                     `json:"count"`
}

// handleServiceError サービス層のエラーをHTTPレスポンスに変換
func handleServiceError(c *gin.Context, err error) {
	switch e := err.(type) {
	case *errors.ValidationError:
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: e.Error(),
			Code:    400,
		})
	case *errors.BusinessLogicError:
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "business_logic_error",
			Message: e.Error(),
			Code:    400,
		})
	case *errors.DatabaseError:
		// データベースエラーの詳細は隠して一般的なメッセージを返す
		if e.Cause != nil && e.Cause.Error() == "resource not found" {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   "not_found",
				Message: "Resource not found",
				Code:    404,
			})
		} else {
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error:   "internal_error",
				Message: "Internal server error",
				Code:    500,
			})
		}
	default:
		// その他のエラーは内部サーバーエラーとして扱う
		if err.Error() == "resource not found" {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   "not_found",
				Message: "Resource not found",
				Code:    404,
			})
		} else {
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error:   "internal_error",
				Message: "Internal server error",
				Code:    500,
			})
		}
	}
}
