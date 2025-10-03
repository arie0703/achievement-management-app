# Variables for ECS Module

# Basic Configuration
variable "environment" {
  description = "Environment name (dev, staging, prod)"
  type        = string
}

variable "app_name" {
  description = "Application name"
  type        = string
}

variable "tags" {
  description = "Tags to apply to all resources"
  type        = map(string)
  default     = {}
}

# Network Configuration
variable "vpc_id" {
  description = "VPC ID where ECS resources will be created"
  type        = string
}

variable "public_subnet_ids" {
  description = "List of public subnet IDs for ALB"
  type        = list(string)
}

variable "private_subnet_ids" {
  description = "List of private subnet IDs for ECS tasks"
  type        = list(string)
}

variable "alb_security_group_id" {
  description = "Security group ID for ALB"
  type        = string
}

variable "ecs_security_group_id" {
  description = "Security group ID for ECS tasks"
  type        = string
}

# IAM Configuration
variable "ecs_task_execution_role_arn" {
  description = "ARN of the ECS task execution role"
  type        = string
}

variable "ecs_task_role_arn" {
  description = "ARN of the ECS task role for API service"
  type        = string
}

variable "cli_task_role_arns" {
  description = "Map of CLI task role ARNs by operation type"
  type        = map(string)
  default = {
    achievement = ""
    points      = ""
    reward      = ""
  }
}

# ECS Cluster Configuration
variable "enable_container_insights" {
  description = "Enable CloudWatch Container Insights for the cluster"
  type        = bool
  default     = true
}

variable "fargate_base_capacity" {
  description = "Base capacity for Fargate capacity provider"
  type        = number
  default     = 1
}

variable "fargate_weight" {
  description = "Weight for Fargate capacity provider"
  type        = number
  default     = 1
}

variable "enable_fargate_spot" {
  description = "Enable Fargate Spot capacity provider"
  type        = bool
  default     = false
}

variable "fargate_spot_base_capacity" {
  description = "Base capacity for Fargate Spot capacity provider"
  type        = number
  default     = 0
}

variable "fargate_spot_weight" {
  description = "Weight for Fargate Spot capacity provider"
  type        = number
  default     = 1
}

# Container Configuration
variable "container_port" {
  description = "Port on which the container application runs"
  type        = number
  default     = 8080
}

variable "api_container_image" {
  description = "Container image for API service"
  type        = string
}

variable "cli_container_image" {
  description = "Container image for CLI tasks"
  type        = string
}

# API Task Configuration
variable "api_task_cpu" {
  description = "CPU units for API task (1024 = 1 vCPU)"
  type        = number
  default     = 512
}

variable "api_task_memory" {
  description = "Memory for API task in MB"
  type        = number
  default     = 1024
}

variable "api_desired_count" {
  description = "Desired number of API service tasks"
  type        = number
  default     = 2
}

variable "api_min_capacity" {
  description = "Minimum number of API service tasks"
  type        = number
  default     = 1
}

variable "api_max_capacity" {
  description = "Maximum number of API service tasks"
  type        = number
  default     = 10
}

# CLI Task Configuration
variable "cli_task_cpu" {
  description = "CPU units for CLI tasks (1024 = 1 vCPU)"
  type        = number
  default     = 256
}

variable "cli_task_memory" {
  description = "Memory for CLI tasks in MB"
  type        = number
  default     = 512
}

variable "cli_task_timeout" {
  description = "Timeout for CLI tasks in seconds"
  type        = number
  default     = 300
}

variable "enable_cli_specific_tasks" {
  description = "Enable creation of specific CLI task definitions for detailed operations"
  type        = bool
  default     = true
}

variable "cli_task_types" {
  description = "List of CLI task types to create"
  type        = list(string)
  default     = ["achievement", "points", "reward"]
}

variable "cli_commands" {
  description = "Commands to run for each CLI task type"
  type        = map(list(string))
  default = {
    achievement = ["./achievement-app", "achievement"]
    points      = ["./achievement-app", "points"]
    reward      = ["./achievement-app", "reward"]
  }
}

# Environment Variables
variable "api_environment_variables" {
  description = "Environment variables for API service"
  type = list(object({
    name  = string
    value = string
  }))
  default = []
}

variable "cli_environment_variables" {
  description = "Environment variables for CLI tasks"
  type = list(object({
    name  = string
    value = string
  }))
  default = []
}

# Load Balancer Configuration
variable "enable_deletion_protection" {
  description = "Enable deletion protection for ALB"
  type        = bool
  default     = false
}

variable "enable_https" {
  description = "Enable HTTPS listener for ALB"
  type        = bool
  default     = false
}

variable "enable_https_redirect" {
  description = "Enable HTTP to HTTPS redirect"
  type        = bool
  default     = false
}

variable "certificate_arn" {
  description = "ARN of the SSL certificate for HTTPS listener"
  type        = string
  default     = ""
}

variable "ssl_policy" {
  description = "SSL policy for HTTPS listener"
  type        = string
  default     = "ELBSecurityPolicy-TLS-1-2-2017-01"
}

# Health Check Configuration
variable "health_check_path" {
  description = "Health check path for ALB target group"
  type        = string
  default     = "/health"
}

variable "health_check_healthy_threshold" {
  description = "Number of consecutive health checks successes required"
  type        = number
  default     = 2
}

variable "health_check_unhealthy_threshold" {
  description = "Number of consecutive health check failures required"
  type        = number
  default     = 3
}

variable "health_check_timeout" {
  description = "Health check timeout in seconds"
  type        = number
  default     = 5
}

variable "health_check_interval" {
  description = "Health check interval in seconds"
  type        = number
  default     = 30
}

variable "health_check_matcher" {
  description = "HTTP response codes to indicate a healthy check"
  type        = string
  default     = "200"
}

# Auto Scaling Configuration
variable "cpu_target_value" {
  description = "Target CPU utilization percentage for auto scaling"
  type        = number
  default     = 70
}

variable "memory_target_value" {
  description = "Target memory utilization percentage for auto scaling"
  type        = number
  default     = 80
}

# Deployment Configuration
variable "deployment_maximum_percent" {
  description = "Maximum percentage of tasks that can be running during deployment"
  type        = number
  default     = 200
}

variable "deployment_minimum_healthy_percent" {
  description = "Minimum percentage of tasks that must remain healthy during deployment"
  type        = number
  default     = 50
}

# Logging Configuration
variable "log_retention_days" {
  description = "CloudWatch log retention in days"
  type        = number
  default     = 7
}