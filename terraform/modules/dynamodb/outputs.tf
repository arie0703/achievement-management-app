# Output values for the DynamoDB module
# These outputs provide information about created DynamoDB tables

output "table_names" {
  description = "Map of table logical names to actual table names"
  value = merge(
    {
      for key, table in aws_dynamodb_table.tables : key => table.name
    },
    length(aws_dynamodb_table.achievements) > 0 ? {
      achievements = aws_dynamodb_table.achievements[0].name
    } : {}
  )
}

output "table_arns" {
  description = "Map of table logical names to table ARNs"
  value = merge(
    {
      for key, table in aws_dynamodb_table.tables : key => table.arn
    },
    length(aws_dynamodb_table.achievements) > 0 ? {
      achievements = aws_dynamodb_table.achievements[0].arn
    } : {}
  )
}

output "achievements_table_name" {
  description = "Name of the achievements table"
  value       = length(aws_dynamodb_table.achievements) > 0 ? aws_dynamodb_table.achievements[0].name : null
}

output "achievements_table_arn" {
  description = "ARN of the achievements table"
  value       = length(aws_dynamodb_table.achievements) > 0 ? aws_dynamodb_table.achievements[0].arn : null
}

output "rewards_table_name" {
  description = "Name of the rewards table"
  value       = try(aws_dynamodb_table.tables["rewards"].name, null)
}

output "rewards_table_arn" {
  description = "ARN of the rewards table"
  value       = try(aws_dynamodb_table.tables["rewards"].arn, null)
}

output "current_points_table_name" {
  description = "Name of the current points table"
  value       = try(aws_dynamodb_table.tables["current_points"].name, null)
}

output "current_points_table_arn" {
  description = "ARN of the current points table"
  value       = try(aws_dynamodb_table.tables["current_points"].arn, null)
}

output "reward_history_table_name" {
  description = "Name of the reward history table"
  value       = try(aws_dynamodb_table.tables["reward_history"].name, null)
}

output "reward_history_table_arn" {
  description = "ARN of the reward history table"
  value       = try(aws_dynamodb_table.tables["reward_history"].arn, null)
}

output "achievements_gsi_name" {
  description = "Name of the achievements GSI"
  value       = length(aws_dynamodb_table.achievements) > 0 ? "achievement-id-index" : null
}