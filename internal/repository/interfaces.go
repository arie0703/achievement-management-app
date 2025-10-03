package repository

import "achievement-management/internal/models"

// TransactWriteItem DynamoDB トランザクション書き込みアイテム
type TransactWriteItem struct {
	TableName string
	Item      interface{}
	Operation string // "PUT", "UPDATE", "DELETE"
}

// Repository DynamoDB操作の抽象化
type Repository interface {
	PutItem(tableName string, item interface{}) error
	GetItem(tableName string, key map[string]interface{}, result interface{}) error
	UpdateItem(tableName string, key map[string]interface{}, updateExpression string, expressionAttributeValues map[string]interface{}) error
	Scan(tableName string, result interface{}) error
	DeleteItem(tableName string, key map[string]interface{}) error
	TransactWrite(items []TransactWriteItem) error
}

// AchievementRepository 達成目録リポジトリ
type AchievementRepository interface {
	Create(achievement *models.Achievement) error
	Update(achievement *models.Achievement) error
	GetByID(id string) (*models.Achievement, error)
	List() ([]*models.Achievement, error)
	Delete(id string) error
}

// RewardRepository 報酬リポジトリ
type RewardRepository interface {
	Create(reward *models.Reward) error
	Update(reward *models.Reward) error
	GetByID(id string) (*models.Reward, error)
	List() ([]*models.Reward, error)
	Delete(id string) error
}

// PointRepository ポイントリポジトリ
type PointRepository interface {
	GetCurrentPoints() (*models.CurrentPoints, error)
	UpdateCurrentPoints(points *models.CurrentPoints) error
	CreateRewardHistory(history *models.RewardHistory) error
	GetRewardHistory() ([]*models.RewardHistory, error)
	TransactPointsAndHistory(pointsUpdate *models.CurrentPoints, history *models.RewardHistory) error
	AddPoints(points int) error
	SubtractPoints(points int) error
}