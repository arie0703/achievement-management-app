# Input variables for the root Terraform module
# These variables define the configuration parameters for the infrastructure

variable "environment" {
  description = "Environment"
  default     = "sandbox"
  type        = string
}

variable "app_name" {
  description = "Application name used for resource naming"
  type        = string
  default     = "achievement-management"
}

variable "aws_region" {
  description = "AWS region for resource deployment"
  type        = string
  default     = "ap-northeast-1"
}

variable "cost_center" {
  description = "Cost center for resource tagging and billing"
  type        = string
  default     = "engineering"
}

# Enhanced Tagging Variables for Resource Organization and Cost Tracking
variable "owner" {
  description = "Owner of the resources (team or individual responsible)"
  type        = string
  default     = "platform-team"
}

variable "team" {
  description = "Team responsible for the resources"
  type        = string
  default     = "engineering"
}

variable "business_unit" {
  description = "Business unit that owns the resources"
  type        = string
  default     = "technology"
}

variable "service_tier" {
  description = "Service tier classification (critical, important, standard)"
  type        = string
  default     = "standard"

  validation {
    condition     = contains(["critical", "important", "standard"], var.service_tier)
    error_message = "Service tier must be one of: critical, important, standard."
  }
}

variable "backup_policy" {
  description = "Backup policy for resources (daily, weekly, none)"
  type        = string
  default     = "daily"

  validation {
    condition     = contains(["daily", "weekly", "monthly", "none"], var.backup_policy)
    error_message = "Backup policy must be one of: daily, weekly, monthly, none."
  }
}

variable "data_classification" {
  description = "Data classification level (public, internal, confidential, restricted)"
  type        = string
  default     = "internal"

  validation {
    condition     = contains(["public", "internal", "confidential", "restricted"], var.data_classification)
    error_message = "Data classification must be one of: public, internal, confidential, restricted."
  }
}

variable "compliance_scope" {
  description = "Compliance requirements that apply to these resources"
  type        = string
  default     = "none"
}

# VPC Configuration Variables
variable "vpc_cidr" {
  description = "CIDR block for the VPC"
  type        = string
  default     = "10.0.0.0/16"
}

variable "availability_zones" {
  description = "List of availability zones to use"
  type        = list(string)
  default     = ["us-east-1a", "us-east-1b"]
}

variable "public_subnet_cidrs" {
  description = "CIDR blocks for public subnets"
  type        = list(string)
  default     = ["10.0.1.0/24", "10.0.2.0/24"]
}

variable "private_subnet_cidrs" {
  description = "CIDR blocks for private subnets"
  type        = list(string)
  default     = ["10.0.10.0/24", "10.0.20.0/24"]
}

# ECS Configuration Variables
variable "ecs_config" {
  description = "ECS service configuration parameters"
  type = object({
    cpu               = number
    memory            = number
    desired_count     = number
    max_capacity      = number
    min_capacity      = number
    container_port    = number
    health_check_path = string
    container_image   = string
  })

  default = {
    cpu               = 256
    memory            = 512
    desired_count     = 1
    max_capacity      = 10
    min_capacity      = 1
    container_port    = 8080
    health_check_path = "/health"
    container_image   = "achievement-management:latest"
  }
}

# Load Balancer Configuration Variables
variable "alb_config" {
  description = "Application Load Balancer configuration parameters"
  type = object({
    enable_https               = bool
    enable_https_redirect      = bool
    certificate_arn            = string
    ssl_policy                 = string
    enable_deletion_protection = bool
  })

  default = {
    enable_https               = false
    enable_https_redirect      = false
    certificate_arn            = ""
    ssl_policy                 = "ELBSecurityPolicy-TLS-1-2-2017-01"
    enable_deletion_protection = false
  }
}

# DynamoDB Configuration Variables
variable "dynamodb_tables" {
  description = "DynamoDB table configurations"
  type = map(object({
    hash_key               = string
    range_key              = optional(string)
    billing_mode           = string
    read_capacity          = optional(number)
    write_capacity         = optional(number)
    point_in_time_recovery = bool
    server_side_encryption = bool
  }))

  default = {
    achievements = {
      hash_key               = "id"
      billing_mode           = "PAY_PER_REQUEST"
      point_in_time_recovery = true
      server_side_encryption = true
    }
    rewards = {
      hash_key               = "id"
      billing_mode           = "PAY_PER_REQUEST"
      point_in_time_recovery = true
      server_side_encryption = true
    }
    current_points = {
      hash_key               = "id"
      billing_mode           = "PAY_PER_REQUEST"
      point_in_time_recovery = true
      server_side_encryption = true
    }
    reward_history = {
      hash_key               = "id"
      billing_mode           = "PAY_PER_REQUEST"
      point_in_time_recovery = true
      server_side_encryption = true
    }
  }
}

# Monitoring Configuration Variables
variable "log_retention_days" {
  description = "CloudWatch log retention period in days"
  type        = number
  default     = 14
}

variable "cli_task_types" {
  description = "List of CLI task types for monitoring"
  type        = list(string)
  default     = ["achievement", "points", "reward"]
}

variable "enable_monitoring_dashboard" {
  description = "Enable CloudWatch dashboard creation"
  type        = bool
  default     = true
}

variable "monitoring_thresholds" {
  description = "CloudWatch alarm thresholds"
  type = object({
    cpu_alarm_threshold         = number
    memory_alarm_threshold      = number
    response_time_threshold     = number
    http_5xx_threshold          = number
    dynamodb_throttle_threshold = number
  })

  default = {
    cpu_alarm_threshold         = 80
    memory_alarm_threshold      = 80
    response_time_threshold     = 2.0
    http_5xx_threshold          = 10
    dynamodb_throttle_threshold = 0
  }
}

variable "alarm_actions" {
  description = "List of ARNs to notify when CloudWatch alarms trigger"
  type        = list(string)
  default     = []
}

variable "enable_container_insights" {
  description = "Enable CloudWatch Container Insights"
  type        = bool
  default     = false
}
