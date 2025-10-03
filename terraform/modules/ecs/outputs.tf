# Outputs for ECS Module

# ECS Cluster
output "cluster_id" {
  description = "ID of the ECS cluster"
  value       = aws_ecs_cluster.main.id
}

output "cluster_name" {
  description = "Name of the ECS cluster"
  value       = aws_ecs_cluster.main.name
}

output "cluster_arn" {
  description = "ARN of the ECS cluster"
  value       = aws_ecs_cluster.main.arn
}

# Load Balancer
output "load_balancer_arn" {
  description = "ARN of the Application Load Balancer"
  value       = aws_lb.main.arn
}

output "load_balancer_dns_name" {
  description = "DNS name of the Application Load Balancer"
  value       = aws_lb.main.dns_name
}

output "load_balancer_zone_id" {
  description = "Zone ID of the Application Load Balancer"
  value       = aws_lb.main.zone_id
}

output "alb_arn_suffix" {
  description = "ARN suffix of the Application Load Balancer for CloudWatch metrics"
  value       = aws_lb.main.arn_suffix
}

output "load_balancer_url" {
  description = "URL of the Application Load Balancer"
  value       = var.enable_https ? "https://${aws_lb.main.dns_name}" : "http://${aws_lb.main.dns_name}"
}

output "load_balancer_http_url" {
  description = "HTTP URL of the Application Load Balancer"
  value       = "http://${aws_lb.main.dns_name}"
}

output "load_balancer_https_url" {
  description = "HTTPS URL of the Application Load Balancer (if enabled)"
  value       = var.enable_https ? "https://${aws_lb.main.dns_name}" : null
}

# Target Group
output "target_group_arn" {
  description = "ARN of the API target group"
  value       = aws_lb_target_group.api.arn
}

output "target_group_name" {
  description = "Name of the API target group"
  value       = aws_lb_target_group.api.name
}

output "target_group_arn_suffix" {
  description = "ARN suffix of the API target group for CloudWatch metrics"
  value       = aws_lb_target_group.api.arn_suffix
}

# Listeners
output "http_listener_arn" {
  description = "ARN of the HTTP listener"
  value       = aws_lb_listener.api_http.arn
}

output "https_listener_arn" {
  description = "ARN of the HTTPS listener (if enabled)"
  value       = var.enable_https ? aws_lb_listener.api_https[0].arn : null
}

# API Service
output "api_service_id" {
  description = "ID of the API ECS service"
  value       = aws_ecs_service.api.id
}

output "api_service_name" {
  description = "Name of the API ECS service"
  value       = aws_ecs_service.api.name
}

output "service_name" {
  description = "Name of the API ECS service (alias for monitoring)"
  value       = aws_ecs_service.api.name
}

output "api_service_cluster" {
  description = "Cluster name of the API ECS service"
  value       = aws_ecs_service.api.cluster
}

# Task Definitions
output "api_task_definition_arn" {
  description = "ARN of the API task definition"
  value       = aws_ecs_task_definition.api.arn
}

output "api_task_definition_family" {
  description = "Family of the API task definition"
  value       = aws_ecs_task_definition.api.family
}

output "api_task_definition_revision" {
  description = "Revision of the API task definition"
  value       = aws_ecs_task_definition.api.revision
}

output "cli_task_definition_arns" {
  description = "Map of CLI task definition ARNs by operation type"
  value = {
    for task_type in var.cli_task_types :
    task_type => aws_ecs_task_definition.cli[task_type].arn
  }
}

output "cli_task_definition_families" {
  description = "Map of CLI task definition families by operation type"
  value = {
    for task_type in var.cli_task_types :
    task_type => aws_ecs_task_definition.cli[task_type].family
  }
}

# Specific CLI Task Definition ARNs for detailed operations
output "cli_specific_task_definition_arns" {
  description = "Map of specific CLI task definition ARNs for detailed operations"
  value = var.enable_cli_specific_tasks ? {
    # Achievement operations
    achievement_create = aws_ecs_task_definition.cli_achievement_create[0].arn
    achievement_list   = aws_ecs_task_definition.cli_achievement_list[0].arn
    
    # Points operations
    points_current   = aws_ecs_task_definition.cli_points_current[0].arn
    points_aggregate = aws_ecs_task_definition.cli_points_aggregate[0].arn
    
    # Reward operations
    reward_create = aws_ecs_task_definition.cli_reward_create[0].arn
    reward_list   = aws_ecs_task_definition.cli_reward_list[0].arn
    reward_redeem = aws_ecs_task_definition.cli_reward_redeem[0].arn
  } : {}
}

output "cli_specific_task_definition_families" {
  description = "Map of specific CLI task definition families for detailed operations"
  value = var.enable_cli_specific_tasks ? {
    # Achievement operations
    achievement_create = aws_ecs_task_definition.cli_achievement_create[0].family
    achievement_list   = aws_ecs_task_definition.cli_achievement_list[0].family
    
    # Points operations
    points_current   = aws_ecs_task_definition.cli_points_current[0].family
    points_aggregate = aws_ecs_task_definition.cli_points_aggregate[0].family
    
    # Reward operations
    reward_create = aws_ecs_task_definition.cli_reward_create[0].family
    reward_list   = aws_ecs_task_definition.cli_reward_list[0].family
    reward_redeem = aws_ecs_task_definition.cli_reward_redeem[0].family
  } : {}
}

# Auto Scaling
output "autoscaling_target_resource_id" {
  description = "Resource ID of the auto scaling target"
  value       = aws_appautoscaling_target.api.resource_id
}

output "autoscaling_policies" {
  description = "Map of auto scaling policy ARNs"
  value = {
    cpu    = aws_appautoscaling_policy.api_cpu.arn
    memory = aws_appautoscaling_policy.api_memory.arn
  }
}

# CloudWatch Log Groups are managed by the monitoring module

# Service Discovery (for future use)
output "service_discovery_namespace" {
  description = "Service discovery namespace (placeholder for future implementation)"
  value       = null
}

# Capacity Provider Information
output "capacity_providers" {
  description = "List of capacity providers configured for the cluster"
  value       = aws_ecs_cluster_capacity_providers.main.capacity_providers
}

# Task Execution Information
output "task_execution_info" {
  description = "Information needed to run CLI tasks"
  value = {
    cluster_name           = aws_ecs_cluster.main.name
    subnet_ids            = var.private_subnet_ids
    security_group_id     = var.ecs_security_group_id
    task_definition_arns  = {
      for task_type in var.cli_task_types :
      task_type => aws_ecs_task_definition.cli[task_type].arn
    }
    specific_task_definition_arns = var.enable_cli_specific_tasks ? {
      # Achievement operations
      achievement_create = aws_ecs_task_definition.cli_achievement_create[0].arn
      achievement_list   = aws_ecs_task_definition.cli_achievement_list[0].arn
      
      # Points operations
      points_current   = aws_ecs_task_definition.cli_points_current[0].arn
      points_aggregate = aws_ecs_task_definition.cli_points_aggregate[0].arn
      
      # Reward operations
      reward_create = aws_ecs_task_definition.cli_reward_create[0].arn
      reward_list   = aws_ecs_task_definition.cli_reward_list[0].arn
      reward_redeem = aws_ecs_task_definition.cli_reward_redeem[0].arn
    } : {}
  }
}