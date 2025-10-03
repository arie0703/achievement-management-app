package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/spf13/cobra"

	"achievement-management/internal/models"
)

// achievementCmd represents the achievement command
var achievementCmd = &cobra.Command{
	Use:   "achievement",
	Short: "Manage achievements",
	Long: `Manage achievements in the system.

You can create, list, update, and delete achievements using this command.
Each achievement has a title, description, and point value.`,
}

// achievementCreateCmd represents the achievement create command
var achievementCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new achievement",
	Long: `Create a new achievement with the specified title, description, and point value.

Example:
  achievement-app achievement create --title "First Login" --description "Log in for the first time" --point 10`,
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

		achievementService, _, _, err := initServices()
		if err != nil {
			return fmt.Errorf("failed to initialize services: %w", err)
		}

		achievement := &models.Achievement{
			Title:       title,
			Description: description,
			Point:       point,
			CreatedAt:   time.Now(),
		}

		if err := achievementService.Create(achievement); err != nil {
			return fmt.Errorf("failed to create achievement: %w", err)
		}

		fmt.Printf("✅ Achievement created successfully!\n")
		fmt.Printf("ID: %s\n", achievement.ID)
		fmt.Printf("Title: %s\n", achievement.Title)
		fmt.Printf("Description: %s\n", achievement.Description)
		fmt.Printf("Points: %d\n", achievement.Point)
		fmt.Printf("Created: %s\n", achievement.CreatedAt.Format("2006-01-02 15:04:05"))

		return nil
	},
}

// achievementListCmd represents the achievement list command
var achievementListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all achievements",
	Long: `List all achievements in the system.

Example:
  achievement-app achievement list`,
	RunE: func(cmd *cobra.Command, args []string) error {
		achievementService, _, _, err := initServices()
		if err != nil {
			return fmt.Errorf("failed to initialize services: %w", err)
		}

		achievements, err := achievementService.List()
		if err != nil {
			return fmt.Errorf("failed to list achievements: %w", err)
		}

		if len(achievements) == 0 {
			fmt.Println("No achievements found.")
			return nil
		}

		fmt.Printf("Found %d achievement(s):\n\n", len(achievements))
		for i, achievement := range achievements {
			fmt.Printf("%d. %s (ID: %s)\n", i+1, achievement.Title, achievement.ID)
			fmt.Printf("   Description: %s\n", achievement.Description)
			fmt.Printf("   Points: %d\n", achievement.Point)
			fmt.Printf("   Created: %s\n", achievement.CreatedAt.Format("2006-01-02 15:04:05"))
			fmt.Println()
		}

		return nil
	},
}

// achievementUpdateCmd represents the achievement update command
var achievementUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update an existing achievement",
	Long: `Update an existing achievement by ID.

Example:
  achievement-app achievement update --id "01234567890" --title "Updated Title" --point 20`,
	RunE: func(cmd *cobra.Command, args []string) error {
		id, _ := cmd.Flags().GetString("id")
		title, _ := cmd.Flags().GetString("title")
		description, _ := cmd.Flags().GetString("description")
		pointStr, _ := cmd.Flags().GetString("point")

		if id == "" {
			return fmt.Errorf("id is required")
		}

		achievementService, _, _, err := initServices()
		if err != nil {
			return fmt.Errorf("failed to initialize services: %w", err)
		}

		// Get existing achievement
		existing, err := achievementService.GetByID(id)
		if err != nil {
			return fmt.Errorf("failed to get achievement: %w", err)
		}

		// Update fields if provided
		updated := &models.Achievement{
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

		if err := achievementService.Update(id, updated); err != nil {
			return fmt.Errorf("failed to update achievement: %w", err)
		}

		fmt.Printf("✅ Achievement updated successfully!\n")
		fmt.Printf("ID: %s\n", updated.ID)
		fmt.Printf("Title: %s\n", updated.Title)
		fmt.Printf("Description: %s\n", updated.Description)
		fmt.Printf("Points: %d\n", updated.Point)
		fmt.Printf("Created: %s\n", updated.CreatedAt.Format("2006-01-02 15:04:05"))

		return nil
	},
}

// achievementDeleteCmd represents the achievement delete command
var achievementDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete an achievement",
	Long: `Delete an achievement by ID.

Example:
  achievement-app achievement delete --id "01234567890"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		id, _ := cmd.Flags().GetString("id")

		if id == "" {
			return fmt.Errorf("id is required")
		}

		achievementService, _, _, err := initServices()
		if err != nil {
			return fmt.Errorf("failed to initialize services: %w", err)
		}

		// Get achievement details before deletion for confirmation
		achievement, err := achievementService.GetByID(id)
		if err != nil {
			return fmt.Errorf("failed to get achievement: %w", err)
		}

		if err := achievementService.Delete(id); err != nil {
			return fmt.Errorf("failed to delete achievement: %w", err)
		}

		fmt.Printf("✅ Achievement deleted successfully!\n")
		fmt.Printf("Deleted: %s (ID: %s)\n", achievement.Title, achievement.ID)

		return nil
	},
}

func init() {
	// Add subcommands to achievement command
	achievementCmd.AddCommand(achievementCreateCmd)
	achievementCmd.AddCommand(achievementListCmd)
	achievementCmd.AddCommand(achievementUpdateCmd)
	achievementCmd.AddCommand(achievementDeleteCmd)

	// Flags for create command
	achievementCreateCmd.Flags().String("title", "", "Achievement title (required)")
	achievementCreateCmd.Flags().String("description", "", "Achievement description")
	achievementCreateCmd.Flags().Int("point", 0, "Achievement point value (required)")
	achievementCreateCmd.MarkFlagRequired("title")
	achievementCreateCmd.MarkFlagRequired("point")

	// Flags for update command
	achievementUpdateCmd.Flags().String("id", "", "Achievement ID (required)")
	achievementUpdateCmd.Flags().String("title", "", "New achievement title")
	achievementUpdateCmd.Flags().String("description", "", "New achievement description")
	achievementUpdateCmd.Flags().String("point", "", "New achievement point value")
	achievementUpdateCmd.MarkFlagRequired("id")

	// Flags for delete command
	achievementDeleteCmd.Flags().String("id", "", "Achievement ID (required)")
	achievementDeleteCmd.MarkFlagRequired("id")
}