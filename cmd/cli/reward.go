package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/spf13/cobra"

	"achievement-management/internal/models"
)

// rewardCmd represents the reward command
var rewardCmd = &cobra.Command{
	Use:   "reward",
	Short: "Manage rewards",
	Long: `Manage rewards in the system.

You can create, list, update, redeem, and delete rewards using this command.
Each reward has a title, description, and point cost.`,
}

// rewardCreateCmd represents the reward create command
var rewardCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new reward",
	Long: `Create a new reward with the specified title, description, and point cost.

Example:
  achievement-app reward create --title "Coffee Voucher" --description "Free coffee at the office" --point 50`,
	RunE: func(cmd *cobra.Command, args []string) error {
		title, _ := cmd.Flags().GetString("title")
		description, _ := cmd.Flags().GetString("description")
		point, _ := cmd.Flags().GetInt("point")

		if title == "" {
			return fmt.Errorf("title is required")
		}
		if point <= 0 {
			return fmt.Errorf("point must be a positive integer")
		}

		_, rewardService, _, err := initServices()
		if err != nil {
			return fmt.Errorf("failed to initialize services: %w", err)
		}

		reward := &models.Reward{
			Title:       title,
			Description: description,
			Point:       point,
			CreatedAt:   time.Now(),
		}

		if err := rewardService.Create(reward); err != nil {
			return fmt.Errorf("failed to create reward: %w", err)
		}

		fmt.Printf("✅ Reward created successfully!\n")
		fmt.Printf("ID: %s\n", reward.ID)
		fmt.Printf("Title: %s\n", reward.Title)
		fmt.Printf("Description: %s\n", reward.Description)
		fmt.Printf("Point Cost: %d\n", reward.Point)
		fmt.Printf("Created: %s\n", reward.CreatedAt.Format("2006-01-02 15:04:05"))

		return nil
	},
}

// rewardListCmd represents the reward list command
var rewardListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all rewards",
	Long: `List all rewards in the system.

Example:
  achievement-app reward list`,
	RunE: func(cmd *cobra.Command, args []string) error {
		_, rewardService, _, err := initServices()
		if err != nil {
			return fmt.Errorf("failed to initialize services: %w", err)
		}

		rewards, err := rewardService.List()
		if err != nil {
			return fmt.Errorf("failed to list rewards: %w", err)
		}

		if len(rewards) == 0 {
			fmt.Println("No rewards found.")
			return nil
		}

		fmt.Printf("Found %d reward(s):\n\n", len(rewards))
		for i, reward := range rewards {
			fmt.Printf("%d. %s (ID: %s)\n", i+1, reward.Title, reward.ID)
			fmt.Printf("   Description: %s\n", reward.Description)
			fmt.Printf("   Point Cost: %d\n", reward.Point)
			fmt.Printf("   Created: %s\n", reward.CreatedAt.Format("2006-01-02 15:04:05"))
			fmt.Println()
		}

		return nil
	},
}

// rewardUpdateCmd represents the reward update command
var rewardUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update an existing reward",
	Long: `Update an existing reward by ID.

Example:
  achievement-app reward update --id "01234567890" --title "Updated Title" --point 75`,
	RunE: func(cmd *cobra.Command, args []string) error {
		id, _ := cmd.Flags().GetString("id")
		title, _ := cmd.Flags().GetString("title")
		description, _ := cmd.Flags().GetString("description")
		pointStr, _ := cmd.Flags().GetString("point")

		if id == "" {
			return fmt.Errorf("id is required")
		}

		_, rewardService, _, err := initServices()
		if err != nil {
			return fmt.Errorf("failed to initialize services: %w", err)
		}

		// Get existing reward
		existing, err := rewardService.GetByID(id)
		if err != nil {
			return fmt.Errorf("failed to get reward: %w", err)
		}

		// Update fields if provided
		updated := &models.Reward{
			ID:          existing.ID,
			Title:       existing.Title,
			Description: existing.Description,
			Point:       existing.Point,
			CreatedAt:   existing.CreatedAt,
		}

		if title != "" {
			updated.Title = title
		}
		if description != "" {
			updated.Description = description
		}
		if pointStr != "" {
			point, err := strconv.Atoi(pointStr)
			if err != nil {
				return fmt.Errorf("invalid point value: %w", err)
			}
			if point <= 0 {
				return fmt.Errorf("point must be a positive integer")
			}
			updated.Point = point
		}

		if err := rewardService.Update(id, updated); err != nil {
			return fmt.Errorf("failed to update reward: %w", err)
		}

		fmt.Printf("✅ Reward updated successfully!\n")
		fmt.Printf("ID: %s\n", updated.ID)
		fmt.Printf("Title: %s\n", updated.Title)
		fmt.Printf("Description: %s\n", updated.Description)
		fmt.Printf("Point Cost: %d\n", updated.Point)
		fmt.Printf("Created: %s\n", updated.CreatedAt.Format("2006-01-02 15:04:05"))

		return nil
	},
}

// rewardRedeemCmd represents the reward redeem command
var rewardRedeemCmd = &cobra.Command{
	Use:   "redeem",
	Short: "Redeem a reward",
	Long: `Redeem a reward by ID. This will deduct the required points from your current balance.

Example:
  achievement-app reward redeem --id "01234567890"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		id, _ := cmd.Flags().GetString("id")

		if id == "" {
			return fmt.Errorf("id is required")
		}

		_, rewardService, pointService, err := initServices()
		if err != nil {
			return fmt.Errorf("failed to initialize services: %w", err)
		}

		// Get reward details before redemption
		reward, err := rewardService.GetByID(id)
		if err != nil {
			return fmt.Errorf("failed to get reward: %w", err)
		}

		// Get current points to show before/after
		currentPoints, err := pointService.GetCurrentPoints()
		if err != nil {
			return fmt.Errorf("failed to get current points: %w", err)
		}

		fmt.Printf("Redeeming reward: %s\n", reward.Title)
		fmt.Printf("Point cost: %d\n", reward.Point)
		fmt.Printf("Current balance: %d\n", currentPoints.Point)

		if currentPoints.Point < reward.Point {
			return fmt.Errorf("insufficient points. Required: %d, Available: %d", reward.Point, currentPoints.Point)
		}

		if err := rewardService.Redeem(id); err != nil {
			return fmt.Errorf("failed to redeem reward: %w", err)
		}

		// Get updated points
		updatedPoints, err := pointService.GetCurrentPoints()
		if err != nil {
			fmt.Printf("⚠️  Reward redeemed but failed to get updated balance: %v\n", err)
		} else {
			fmt.Printf("✅ Reward redeemed successfully!\n")
			fmt.Printf("Reward: %s\n", reward.Title)
			fmt.Printf("Points deducted: %d\n", reward.Point)
			fmt.Printf("New balance: %d\n", updatedPoints.Point)
		}

		return nil
	},
}

// rewardDeleteCmd represents the reward delete command
var rewardDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a reward",
	Long: `Delete a reward by ID.

Example:
  achievement-app reward delete --id "01234567890"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		id, _ := cmd.Flags().GetString("id")

		if id == "" {
			return fmt.Errorf("id is required")
		}

		_, rewardService, _, err := initServices()
		if err != nil {
			return fmt.Errorf("failed to initialize services: %w", err)
		}

		// Get reward details before deletion for confirmation
		reward, err := rewardService.GetByID(id)
		if err != nil {
			return fmt.Errorf("failed to get reward: %w", err)
		}

		if err := rewardService.Delete(id); err != nil {
			return fmt.Errorf("failed to delete reward: %w", err)
		}

		fmt.Printf("✅ Reward deleted successfully!\n")
		fmt.Printf("Deleted: %s (ID: %s)\n", reward.Title, reward.ID)

		return nil
	},
}

func init() {
	// Add subcommands to reward command
	rewardCmd.AddCommand(rewardCreateCmd)
	rewardCmd.AddCommand(rewardListCmd)
	rewardCmd.AddCommand(rewardUpdateCmd)
	rewardCmd.AddCommand(rewardRedeemCmd)
	rewardCmd.AddCommand(rewardDeleteCmd)

	// Flags for create command
	rewardCreateCmd.Flags().String("title", "", "Reward title (required)")
	rewardCreateCmd.Flags().String("description", "", "Reward description")
	rewardCreateCmd.Flags().Int("point", 0, "Reward point cost (required)")
	rewardCreateCmd.MarkFlagRequired("title")
	rewardCreateCmd.MarkFlagRequired("point")

	// Flags for update command
	rewardUpdateCmd.Flags().String("id", "", "Reward ID (required)")
	rewardUpdateCmd.Flags().String("title", "", "New reward title")
	rewardUpdateCmd.Flags().String("description", "", "New reward description")
	rewardUpdateCmd.Flags().String("point", "", "New reward point cost")
	rewardUpdateCmd.MarkFlagRequired("id")

	// Flags for redeem command
	rewardRedeemCmd.Flags().String("id", "", "Reward ID (required)")
	rewardRedeemCmd.MarkFlagRequired("id")

	// Flags for delete command
	rewardDeleteCmd.Flags().String("id", "", "Reward ID (required)")
	rewardDeleteCmd.MarkFlagRequired("id")
}