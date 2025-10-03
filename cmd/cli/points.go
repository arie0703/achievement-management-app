package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

// pointsCmd represents the points command
var pointsCmd = &cobra.Command{
	Use:   "points",
	Short: "Manage points",
	Long: `Manage points in the system.

You can view current points, aggregate points from achievements, and view reward redemption history.`,
}

// pointsCurrentCmd represents the points current command
var pointsCurrentCmd = &cobra.Command{
	Use:   "current",
	Short: "Show current point balance",
	Long: `Show the current point balance.

Example:
  achievement-app points current`,
	RunE: func(cmd *cobra.Command, args []string) error {
		_, _, pointService, err := initServices()
		if err != nil {
			return fmt.Errorf("failed to initialize services: %w", err)
		}

		currentPoints, err := pointService.GetCurrentPoints()
		if err != nil {
			return fmt.Errorf("failed to get current points: %w", err)
		}

		fmt.Printf("💰 Current Point Balance\n")
		fmt.Printf("Points: %d\n", currentPoints.Point)
		fmt.Printf("Last Updated: %s\n", currentPoints.UpdatedAt.Format("2006-01-02 15:04:05"))

		return nil
	},
}

// pointsAggregateCmd represents the points aggregate command
var pointsAggregateCmd = &cobra.Command{
	Use:   "aggregate",
	Short: "Show point aggregation summary",
	Long: `Show a summary of points aggregated from all achievements and compare with current balance.

Example:
  achievement-app points aggregate`,
	RunE: func(cmd *cobra.Command, args []string) error {
		_, _, pointService, err := initServices()
		if err != nil {
			return fmt.Errorf("failed to initialize services: %w", err)
		}

		summary, err := pointService.AggregatePoints()
		if err != nil {
			return fmt.Errorf("failed to aggregate points: %w", err)
		}

		fmt.Printf("📊 Point Aggregation Summary\n")
		fmt.Printf("═══════════════════════════════\n")
		fmt.Printf("Total Achievements: %d\n", summary.TotalAchievements)
		fmt.Printf("Total Points from Achievements: %d\n", summary.TotalPoints)
		fmt.Printf("Current Balance: %d\n", summary.CurrentBalance)
		fmt.Printf("Difference: %d\n", summary.Difference)

		if summary.Difference == 0 {
			fmt.Printf("✅ Points are in sync!\n")
		} else if summary.Difference > 0 {
			fmt.Printf("⚠️  Current balance is %d points higher than expected.\n", summary.Difference)
			fmt.Printf("   This might indicate a data inconsistency.\n")
		} else {
			fmt.Printf("⚠️  Current balance is %d points lower than expected.\n", -summary.Difference)
			fmt.Printf("   This is normal if rewards have been redeemed.\n")
		}

		return nil
	},
}

// pointsHistoryCmd represents the points history command
var pointsHistoryCmd = &cobra.Command{
	Use:   "history",
	Short: "Show reward redemption history",
	Long: `Show the history of reward redemptions.

Example:
  achievement-app points history`,
	RunE: func(cmd *cobra.Command, args []string) error {
		_, _, pointService, err := initServices()
		if err != nil {
			return fmt.Errorf("failed to initialize services: %w", err)
		}

		history, err := pointService.GetRewardHistory()
		if err != nil {
			return fmt.Errorf("failed to get reward history: %w", err)
		}

		fmt.Printf("📜 Reward Redemption History\n")
		fmt.Printf("═══════════════════════════════\n")

		if len(history) == 0 {
			fmt.Println("No reward redemptions found.")
			return nil
		}

		fmt.Printf("Found %d redemption(s):\n\n", len(history))
		for i, record := range history {
			fmt.Printf("%d. %s (ID: %s)\n", i+1, record.RewardTitle, record.RewardID)
			fmt.Printf("   Points Used: %d\n", record.PointCost)
			fmt.Printf("   Redeemed: %s\n", record.RedeemedAt.Format("2006-01-02 15:04:05"))
			fmt.Println()
		}

		return nil
	},
}

func init() {
	// Add subcommands to points command
	pointsCmd.AddCommand(pointsCurrentCmd)
	pointsCmd.AddCommand(pointsAggregateCmd)
	pointsCmd.AddCommand(pointsHistoryCmd)
}