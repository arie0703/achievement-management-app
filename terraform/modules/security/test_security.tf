# Test configuration for security module
# This file can be used to test the security module in isolation

# Uncomment the following to test the module independently
# terraform {
#   required_providers {
#     aws = {
#       source  = "hashicorp/aws"
#       version = ">= 5.0"
#     }
#   }
# }

# provider "aws" {
#   region = "us-east-1"
# }

# # Test VPC for security group testing
# resource "aws_vpc" "test" {
#   cidr_block           = "10.0.0.0/16"
#   enable_dns_hostnames = true
#   enable_dns_support   = true
#   
#   tags = {
#     Name = "test-vpc"
#   }
# }

# # Test the security module
# module "test_security" {
#   source = "./"
#   
#   environment = "test"
#   app_name    = "test-app"
#   vpc_id      = aws_vpc.test.id
#   
#   private_subnet_cidrs = ["10.0.10.0/24", "10.0.20.0/24"]
#   container_port       = 8080
#   
#   tags = {
#     Environment = "test"
#     Purpose     = "testing"
#   }
# }