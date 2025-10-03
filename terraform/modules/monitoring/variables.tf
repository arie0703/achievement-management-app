# Variables for the CloudWatch Monitoring Module

variable "environment" {
  description = "Environment name (dev, staging, prod)"
  type        = string
}

variable "app_name" {
  description = "Application name used for resource naming"
  type        = string
}

variable "log_retention_days" {
  description = "CloudWatch log retention period in days"
  type        = number
  default     = 14
}

variable "cli_task_types" {
  description = "List of CLI task types for log group creation"
  type        = list(string)
  default     = ["achievement", "points", "reward"]
}

variable "enable_alb_logs" {
  description = "Enable ALB access logs to CloudWatch"
  type        = bool
  default     = false
}

variable "enable_alb_monitoring" {
  description = "Enable ALB CloudWatch alarms"
  type        = bool
  default     = true
}

variable "enable_dashboard" {
  description = "Enable CloudWatch dashboard creation"
  type        = bool
  default     = true
}

# ECS Service Configuration for Monitoring
variable "ecs_service_name" {
  description = "Name of the ECS service to monitor"
  type        = string
}

variable "ecs_cluster_name" {
  description = "Name of the ECS cluster to monitor"
  type        = string
}

# ALB Configuration for Monitoring
variable "alb_arn_suffix" {
  description = "ARN suffix of the Application Load Balancer"
  type        = string
  default     = ""
}

variable "target_group_arn_suffix" {
  description = "ARN suffix of the ALB target group"
  type        = string
  default     = ""
}

# DynamoDB Configuration for Monitoring
variable "dynamodb_table_names" {
  description = "Map of DynamoDB table logical names to actual table names"
  type        = map(string)
  default     = {}
}

# Alarm Thresholds
variable "cpu_alarm_threshold" {
  description = "CPU utilization threshold for alarms (percentage)"
  type        = number
  default     = 80
}

variable "memory_alarm_threshold" {
  description = "Memory utilization threshold for alarms (percentage)"
  type        = number
  default     = 80
}

variable "min_task_count_threshold" {
  description = "Minimum task count threshold for alarms"
  type        = number
  default     = 1
}

variable "response_time_threshold" {
  description = "ALB response time threshold for alarms (seconds)"
  type        = number
  default     = 2.0
}

variable "http_5xx_threshold" {
  description = "HTTP 5XX error count threshold for alarms"
  type        = number
  default     = 10
}

variable "unhealthy_host_threshold" {
  description = "Unhealthy host count threshold for alarms"
  type        = number
  default     = 0
}

variable "dynamodb_throttle_threshold" {
  description = "DynamoDB throttle events threshold for alarms"
  type        = number
  default     = 0
}

# Alarm Actions
variable "alarm_actions" {
  description = "List of ARNs to notify when alarm triggers"
  type        = list(string)
  default     = []
}

# Common Tags
variable "tags" {
  description = "Common tags to apply to all resources"
  type        = map(string)
  default     = {}
}