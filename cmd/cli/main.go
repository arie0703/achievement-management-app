package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"

	"achievement-management/internal/config"
	"achievement-management/internal/repository"
	"achievement-management/internal/services"
)

// Version information (set by build flags)
var (
	Version    = "dev"
	BuildTime  = "unknown"
	CommitHash = "unknown"
)

var (
	cfgFile   string
	logLevel  string
	verbose   bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "achievement-app",
	Short: "Achievement Management CLI Tool",
	Long: `Achievement Management CLI Tool

A command-line interface for managing achievements, rewards, and points.
This tool allows you to create, update, list, and delete achievements and rewards,
as well as manage points and view aggregation reports.`,
	Version: fmt.Sprintf("%s (built: %s, commit: %s)", Version, BuildTime, CommitHash),
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.achievement-app.yaml)")
	rootCmd.PersistentFlags().StringVar(&logLevel, "log-level", "info", "log level (debug, info, warn, error)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")

	// Add subcommands
	rootCmd.AddCommand(achievementCmd)
	rootCmd.AddCommand(rewardCmd)
	rootCmd.AddCommand(pointsCmd)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if verbose {
		log.SetOutput(os.Stdout)
	}
}

// initServices initializes the services with DynamoDB repository
func initServices() (services.AchievementService, services.RewardService, services.PointService, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to load configuration: %w", err)
	}
	
	// Initialize DynamoDB repository
	repo, err := repository.NewDynamoDBRepository(context.Background(), cfg)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to initialize repository: %w", err)
	}

	// Initialize services
	achievementRepo := repository.NewAchievementRepository(repo, cfg)
	rewardRepo := repository.NewRewardRepository(repo, cfg)
	pointRepo := repository.NewPointRepository(repo, cfg)

	achievementService := services.NewAchievementService(achievementRepo, pointRepo)
	rewardService := services.NewRewardService(rewardRepo, pointRepo)
	pointService := services.NewPointService(pointRepo, achievementRepo)

	return achievementService, rewardService, pointService, nil
}

func main() {
	Execute()
}