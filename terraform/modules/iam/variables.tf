# Variables for IAM Module

variable "app_name" {
  description = "Name of the application"
  type        = string
}

variable "environment" {
  description = "Environment name (dev, staging, prod)"
  type        = string
}

variable "dynamodb_table_names" {
  description = "List of DynamoDB table names that the application needs access to"
  type        = list(string)
  default     = ["achievements", "rewards", "current_points", "reward_history"]
}

variable "tags" {
  description = "Tags to apply to all resources"
  type        = map(string)
  default     = {}
}