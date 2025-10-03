package main

import (
	"achievement-management/internal/config"
	"achievement-management/internal/handlers"
	"achievement-management/internal/repository"
	"achievement-management/internal/services"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

// Version information (set by build flags)
var (
	Version    = "dev"
	BuildTime  = "unknown"
	CommitHash = "unknown"
)

func main() {
	log.Printf("Starting Achievement Management API Server v%s (built: %s, commit: %s)", Version, BuildTime, CommitHash)

	// 設定を読み込み
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// コンテキストを作成
	ctx := context.Background()

	// DynamoDBリポジトリを初期化
	dynamoRepo, err := repository.NewDynamoDBRepository(ctx, cfg)
	if err != nil {
		log.Fatalf("Failed to initialize DynamoDB repository: %v", err)
	}

	// 各リポジトリを初期化
	achievementRepo := repository.NewAchievementRepository(dynamoRepo, cfg)
	rewardRepo := repository.NewRewardRepository(dynamoRepo, cfg)
	pointRepo := repository.NewPointRepository(dynamoRepo, cfg)

	// サービス層を初期化
	achievementService := services.NewAchievementService(achievementRepo, pointRepo)
	rewardService := services.NewRewardService(rewardRepo, pointRepo)
	pointService := services.NewPointService(pointRepo, achievementRepo)

	// HTTPサーバーを初期化
	server := handlers.NewServer(achievementService, rewardService, pointService, cfg)

	// サーバーを起動
	serverAddr := fmt.Sprintf(":%s", cfg.Server.Port)
	log.Printf("Server starting on port %s", cfg.Server.Port)

	// グレースフルシャットダウンの設定
	go func() {
		if err := server.Run(serverAddr); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// シグナルを待機してグレースフルシャットダウン
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Server shutting down...")
}