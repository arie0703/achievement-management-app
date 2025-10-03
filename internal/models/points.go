package models

import "time"

// CurrentPoints 現在のポイント
type CurrentPoints struct {
	ID        string    `json:"id" dynamodbav:"id"` // 固定値 "current"
	Point     int       `json:"point" dynamodbav:"point"`
	UpdatedAt time.Time `json:"updated_at" dynamodbav:"updated_at"`
}

// RewardHistory 報酬獲得履歴
type RewardHistory struct {
	ID          string    `json:"id" dynamodbav:"id"`
	RewardID    string    `json:"reward_id" dynamodbav:"reward_id"`
	RewardTitle string    `json:"reward_title" dynamodbav:"reward_title"`
	PointCost   int       `json:"point_cost" dynamodbav:"point_cost"`
	RedeemedAt  time.Time `json:"redeemed_at" dynamodbav:"redeemed_at"`
}

// PointSummary ポイント集計結果
type PointSummary struct {
	TotalAchievements int `json:"total_achievements"`
	TotalPoints       int `json:"total_points"`
	CurrentBalance    int `json:"current_balance"`
	Difference        int `json:"difference"`
}
