# Terraform configuration block
# Defines required providers, Terraform version, and backend configuration

terraform {
  # Minimum Terraform version required
  required_version = ">= 1.0"

  # Required providers and their versions
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }

    random = {
      source  = "hashicorp/random"
      version = "~> 3.1"
    }
  }
}

# AWS Provider configuration
provider "aws" {
  region = var.aws_region

  # Default tags applied to all resources
  default_tags {
    tags = {
      Environment = var.environment
      Application = var.app_name
      ManagedBy   = "terraform"
      Project     = "achievement-management"
    }
  }
}

# Random provider for generating unique resource names when needed
provider "random" {
  # No specific configuration required
}
