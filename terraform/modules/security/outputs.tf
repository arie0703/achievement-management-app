# Outputs for Security Groups Module
# These outputs provide security group IDs for use by other modules

output "alb_security_group_id" {
  description = "ID of the Application Load Balancer security group"
  value       = aws_security_group.alb.id
}

output "ecs_security_group_id" {
  description = "ID of the ECS services security group"
  value       = aws_security_group.ecs.id
}

output "vpc_endpoints_security_group_id" {
  description = "ID of the VPC endpoints security group"
  value       = aws_security_group.vpc_endpoints.id
}

output "alb_security_group_arn" {
  description = "ARN of the Application Load Balancer security group"
  value       = aws_security_group.alb.arn
}

output "ecs_security_group_arn" {
  description = "ARN of the ECS services security group"
  value       = aws_security_group.ecs.arn
}

output "vpc_endpoints_security_group_arn" {
  description = "ARN of the VPC endpoints security group"
  value       = aws_security_group.vpc_endpoints.arn
}

# Security group names for reference
output "security_group_names" {
  description = "Map of security group names"
  value = {
    alb           = aws_security_group.alb.name
    ecs           = aws_security_group.ecs.name
    vpc_endpoints = aws_security_group.vpc_endpoints.name
  }
}

# Security group rule IDs for reference
output "security_group_rule_ids" {
  description = "Map of security group rule IDs"
  value = {
    alb_to_ecs         = aws_security_group_rule.alb_to_ecs.id
    alb_egress_to_ecs  = aws_security_group_rule.alb_egress_to_ecs.id
  }
}