# Outputs for IAM Module

# ECS Task Execution Role
output "ecs_task_execution_role_arn" {
  description = "ARN of the ECS task execution role"
  value       = aws_iam_role.ecs_task_execution_role.arn
}

output "ecs_task_execution_role_name" {
  description = "Name of the ECS task execution role"
  value       = aws_iam_role.ecs_task_execution_role.name
}

# ECS Task Role (for API service)
output "ecs_task_role_arn" {
  description = "ARN of the ECS task role for API service"
  value       = aws_iam_role.ecs_task_role.arn
}

output "ecs_task_role_name" {
  description = "Name of the ECS task role for API service"
  value       = aws_iam_role.ecs_task_role.name
}

# CLI Task Roles
output "cli_achievement_task_role_arn" {
  description = "ARN of the CLI achievement task role"
  value       = aws_iam_role.cli_achievement_task_role.arn
}

output "cli_achievement_task_role_name" {
  description = "Name of the CLI achievement task role"
  value       = aws_iam_role.cli_achievement_task_role.name
}

output "cli_points_task_role_arn" {
  description = "ARN of the CLI points task role"
  value       = aws_iam_role.cli_points_task_role.arn
}

output "cli_points_task_role_name" {
  description = "Name of the CLI points task role"
  value       = aws_iam_role.cli_points_task_role.name
}

output "cli_reward_task_role_arn" {
  description = "ARN of the CLI reward task role"
  value       = aws_iam_role.cli_reward_task_role.arn
}

output "cli_reward_task_role_name" {
  description = "Name of the CLI reward task role"
  value       = aws_iam_role.cli_reward_task_role.name
}

# All CLI task role ARNs for convenience
output "cli_task_role_arns" {
  description = "Map of CLI task role ARNs by operation type"
  value = {
    achievement = aws_iam_role.cli_achievement_task_role.arn
    points      = aws_iam_role.cli_points_task_role.arn
    reward      = aws_iam_role.cli_reward_task_role.arn
  }
}