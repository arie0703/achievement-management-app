package models

import "time"

// Reward 報酬
type Reward struct {
	ID          string    `json:"id" dynamodbav:"id"`
	Title       string    `json:"title" dynamodbav:"title"`
	Description string    `json:"description" dynamodbav:"description"`
	Point       int       `json:"point" dynamodbav:"point"`
	CreatedAt   time.Time `json:"created_at" dynamodbav:"created_at"`
}