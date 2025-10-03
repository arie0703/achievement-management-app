# Input variables for the DynamoDB module
# These variables define the configuration parameters for DynamoDB tables

variable "environment" {
  description = "Environment name (dev, staging, prod)"
  type        = string
}

variable "app_name" {
  description = "Application name used for resource naming"
  type        = string
}

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
}

variable "tags" {
  description = "Common tags to apply to all resources"
  type        = map(string)
  default     = {}
}