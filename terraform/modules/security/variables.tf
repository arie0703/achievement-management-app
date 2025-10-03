# Variables for Security Groups Module
# These variables define the input parameters for security group configuration

variable "environment" {
  description = "Environment name (dev, staging, prod)"
  type        = string
}

variable "app_name" {
  description = "Application name used for resource naming"
  type        = string
}

variable "vpc_id" {
  description = "VPC ID where security groups will be created"
  type        = string
}

variable "private_subnet_cidrs" {
  description = "CIDR blocks for private subnets (used for VPC endpoint access)"
  type        = list(string)
  default     = ["10.0.10.0/24", "10.0.20.0/24"]
}

variable "allowed_cidr_blocks" {
  description = "CIDR blocks allowed to access the ALB (default: internet)"
  type        = list(string)
  default     = ["0.0.0.0/0"]
}

variable "container_port" {
  description = "Port on which the container application runs"
  type        = number
  default     = 8080
}

variable "tags" {
  description = "Tags to apply to all security group resources"
  type        = map(string)
  default     = {}
}