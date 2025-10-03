# Outputs for the CloudWatch Monitoring Module

# CloudWatch Log Groups
output "api_log_group_name" {
  description = "Name of the API CloudWatch log group"
  value       = aws_cloudwatch_log_group.ecs_api.name
}

output "api_log_group_arn" {
  description = "ARN of the API CloudWatch log group"
  value       = aws_cloudwatch_log_group.ecs_api.arn
}

output "cli_log_group_names" {
  description = "Names of the CLI CloudWatch log groups"
  value       = { for k, v in aws_cloudwatch_log_group.ecs_cli : k => v.name }
}

output "cli_log_group_arns" {
  description = "ARNs of the CLI CloudWatch log groups"
  value       = { for k, v in aws_cloudwatch_log_group.ecs_cli : k => v.arn }
}

output "alb_log_group_name" {
  description = "Name of the ALB CloudWatch log group"
  value       = var.enable_alb_logs ? aws_cloudwatch_log_group.alb[0].name : null
}

output "alb_log_group_arn" {
  description = "ARN of the ALB CloudWatch log group"
  value       = var.enable_alb_logs ? aws_cloudwatch_log_group.alb[0].arn : null
}

# CloudWatch Alarms
output "ecs_cpu_alarm_arn" {
  description = "ARN of the ECS CPU utilization alarm"
  value       = aws_cloudwatch_metric_alarm.ecs_cpu_high.arn
}

output "ecs_memory_alarm_arn" {
  description = "ARN of the ECS memory utilization alarm"
  value       = aws_cloudwatch_metric_alarm.ecs_memory_high.arn
}

output "ecs_task_count_alarm_arn" {
  description = "ARN of the ECS task count alarm"
  value       = aws_cloudwatch_metric_alarm.ecs_task_count_low.arn
}

output "alb_response_time_alarm_arn" {
  description = "ARN of the ALB response time alarm"
  value       = var.enable_alb_monitoring ? aws_cloudwatch_metric_alarm.alb_response_time_high[0].arn : null
}

output "alb_http_5xx_alarm_arn" {
  description = "ARN of the ALB HTTP 5XX alarm"
  value       = var.enable_alb_monitoring ? aws_cloudwatch_metric_alarm.alb_http_5xx_high[0].arn : null
}

output "alb_unhealthy_hosts_alarm_arn" {
  description = "ARN of the ALB unhealthy hosts alarm"
  value       = var.enable_alb_monitoring ? aws_cloudwatch_metric_alarm.alb_unhealthy_hosts[0].arn : null
}

output "dynamodb_read_throttle_alarm_arns" {
  description = "ARNs of the DynamoDB read throttle alarms"
  value       = { for k, v in aws_cloudwatch_metric_alarm.dynamodb_read_throttle : k => v.arn }
}

output "dynamodb_write_throttle_alarm_arns" {
  description = "ARNs of the DynamoDB write throttle alarms"
  value       = { for k, v in aws_cloudwatch_metric_alarm.dynamodb_write_throttle : k => v.arn }
}

# CloudWatch Dashboard
output "dashboard_url" {
  description = "URL of the CloudWatch dashboard"
  value       = var.enable_dashboard ? "https://console.aws.amazon.com/cloudwatch/home?region=${data.aws_region.current.name}#dashboards:name=${aws_cloudwatch_dashboard.main[0].dashboard_name}" : null
}

output "dashboard_name" {
  description = "Name of the CloudWatch dashboard"
  value       = var.enable_dashboard ? aws_cloudwatch_dashboard.main[0].dashboard_name : null
}

# CloudWatch Log Insights Queries
output "log_insights_queries" {
  description = "CloudWatch Log Insights query definitions"
  value = {
    api_errors      = aws_cloudwatch_query_definition.api_errors.name
    api_performance = aws_cloudwatch_query_definition.api_performance.name
    cli_task_status = aws_cloudwatch_query_definition.cli_task_status.name
  }
}