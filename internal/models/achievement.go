package models

import "time"

// Achievement 達成目録
type Achievement struct {
	ID          string    `json:"id" dynamodbav:"id"`
	Title       string    `json:"title" dynamodbav:"title"`
	Description string    `json:"description" dynamodbav:"description"`
	Point       int       `json:"point" dynamodbav:"point"`
	CreatedAt   time.Time `json:"created_at" dynamodbav:"created_at"`
}
