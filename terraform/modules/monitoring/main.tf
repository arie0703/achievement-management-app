# CloudWatch Monitoring Module for Achievement Management Application
# This module creates CloudWatch log groups, metrics, and alarms for monitoring

# Data sources
data "aws_caller_identity" "current" {}
data "aws_region" "current" {}

# Local values for common configurations
locals {
  name_prefix = "${var.app_name}-${var.environment}"
}

# CloudWatch Log Groups for ECS API Service
resource "aws_cloudwatch_log_group" "ecs_api" {
  name              = "/ecs/${local.name_prefix}-api"
  retention_in_days = var.log_retention_days

  tags = merge(var.tags, {
    Name = "${local.name_prefix}-api-log-group"
    ResourceType = "cloudwatch-log-group"
    LogGroupType = "ecs-api"
    ServiceType = "api"
    RetentionDays = tostring(var.log_retention_days)
  })
}

# CloudWatch Log Groups for CLI tasks
resource "aws_cloudwatch_log_group" "ecs_cli" {
  for_each = toset(var.cli_task_types)
  
  name              = "/ecs/${local.name_prefix}-cli-${each.key}"
  retention_in_days = var.log_retention_days

  tags = merge(var.tags, {
    Name = "${local.name_prefix}-cli-${each.key}-log-group"
    ResourceType = "cloudwatch-log-group"
    LogGroupType = "ecs-cli"
    ServiceType = "cli"
    CLITaskType = each.key
    RetentionDays = tostring(var.log_retention_days)
  })
}

# CloudWatch Log Group for Application Load Balancer
resource "aws_cloudwatch_log_group" "alb" {
  count = var.enable_alb_logs ? 1 : 0
  
  name              = "/aws/applicationloadbalancer/${local.name_prefix}-alb"
  retention_in_days = var.log_retention_days

  tags = var.tags
}

# CloudWatch Metric Alarms for ECS Service
resource "aws_cloudwatch_metric_alarm" "ecs_cpu_high" {
  alarm_name          = "${local.name_prefix}-ecs-cpu-high"
  comparison_operator = "GreaterThanThreshold"
  evaluation_periods  = "2"
  metric_name         = "CPUUtilization"
  namespace           = "AWS/ECS"
  period              = "300"
  statistic           = "Average"
  threshold           = var.cpu_alarm_threshold
  alarm_description   = "This metric monitors ECS service CPU utilization"
  alarm_actions       = var.alarm_actions

  dimensions = {
    ServiceName = var.ecs_service_name
    ClusterName = var.ecs_cluster_name
  }

  tags = merge(var.tags, {
    Name = "${local.name_prefix}-ecs-cpu-high-alarm"
    ResourceType = "cloudwatch-alarm"
    AlarmType = "ecs-cpu"
    MetricName = "CPUUtilization"
    Threshold = tostring(var.cpu_alarm_threshold)
    Severity = "warning"
  })
}

resource "aws_cloudwatch_metric_alarm" "ecs_memory_high" {
  alarm_name          = "${local.name_prefix}-ecs-memory-high"
  comparison_operator = "GreaterThanThreshold"
  evaluation_periods  = "2"
  metric_name         = "MemoryUtilization"
  namespace           = "AWS/ECS"
  period              = "300"
  statistic           = "Average"
  threshold           = var.memory_alarm_threshold
  alarm_description   = "This metric monitors ECS service memory utilization"
  alarm_actions       = var.alarm_actions

  dimensions = {
    ServiceName = var.ecs_service_name
    ClusterName = var.ecs_cluster_name
  }

  tags = var.tags
}

# CloudWatch Metric Alarm for ECS Service Task Count
resource "aws_cloudwatch_metric_alarm" "ecs_task_count_low" {
  alarm_name          = "${local.name_prefix}-ecs-task-count-low"
  comparison_operator = "LessThanThreshold"
  evaluation_periods  = "2"
  metric_name         = "RunningTaskCount"
  namespace           = "AWS/ECS"
  period              = "300"
  statistic           = "Average"
  threshold           = var.min_task_count_threshold
  alarm_description   = "This metric monitors ECS service running task count"
  alarm_actions       = var.alarm_actions

  dimensions = {
    ServiceName = var.ecs_service_name
    ClusterName = var.ecs_cluster_name
  }

  tags = var.tags
}

# CloudWatch Metric Alarms for Application Load Balancer
resource "aws_cloudwatch_metric_alarm" "alb_response_time_high" {
  count = var.enable_alb_monitoring ? 1 : 0
  
  alarm_name          = "${local.name_prefix}-alb-response-time-high"
  comparison_operator = "GreaterThanThreshold"
  evaluation_periods  = "2"
  metric_name         = "TargetResponseTime"
  namespace           = "AWS/ApplicationELB"
  period              = "300"
  statistic           = "Average"
  threshold           = var.response_time_threshold
  alarm_description   = "This metric monitors ALB target response time"
  alarm_actions       = var.alarm_actions

  dimensions = {
    LoadBalancer = var.alb_arn_suffix
  }

  tags = var.tags
}

resource "aws_cloudwatch_metric_alarm" "alb_http_5xx_high" {
  count = var.enable_alb_monitoring ? 1 : 0
  
  alarm_name          = "${local.name_prefix}-alb-http-5xx-high"
  comparison_operator = "GreaterThanThreshold"
  evaluation_periods  = "2"
  metric_name         = "HTTPCode_Target_5XX_Count"
  namespace           = "AWS/ApplicationELB"
  period              = "300"
  statistic           = "Sum"
  threshold           = var.http_5xx_threshold
  alarm_description   = "This metric monitors ALB 5XX error count"
  alarm_actions       = var.alarm_actions

  dimensions = {
    LoadBalancer = var.alb_arn_suffix
  }

  tags = var.tags
}

resource "aws_cloudwatch_metric_alarm" "alb_unhealthy_hosts" {
  count = var.enable_alb_monitoring ? 1 : 0
  
  alarm_name          = "${local.name_prefix}-alb-unhealthy-hosts"
  comparison_operator = "GreaterThanThreshold"
  evaluation_periods  = "2"
  metric_name         = "UnHealthyHostCount"
  namespace           = "AWS/ApplicationELB"
  period              = "300"
  statistic           = "Average"
  threshold           = var.unhealthy_host_threshold
  alarm_description   = "This metric monitors ALB unhealthy host count"
  alarm_actions       = var.alarm_actions

  dimensions = {
    TargetGroup  = var.target_group_arn_suffix
    LoadBalancer = var.alb_arn_suffix
  }

  tags = var.tags
}

# CloudWatch Metric Alarms for DynamoDB Tables
resource "aws_cloudwatch_metric_alarm" "dynamodb_read_throttle" {
  for_each = var.dynamodb_table_names
  
  alarm_name          = "${local.name_prefix}-dynamodb-${each.key}-read-throttle"
  comparison_operator = "GreaterThanThreshold"
  evaluation_periods  = "2"
  metric_name         = "ReadThrottledEvents"
  namespace           = "AWS/DynamoDB"
  period              = "300"
  statistic           = "Sum"
  threshold           = var.dynamodb_throttle_threshold
  alarm_description   = "This metric monitors DynamoDB read throttling for ${each.key} table"
  alarm_actions       = var.alarm_actions

  dimensions = {
    TableName = each.value
  }

  tags = var.tags
}

resource "aws_cloudwatch_metric_alarm" "dynamodb_write_throttle" {
  for_each = var.dynamodb_table_names
  
  alarm_name          = "${local.name_prefix}-dynamodb-${each.key}-write-throttle"
  comparison_operator = "GreaterThanThreshold"
  evaluation_periods  = "2"
  metric_name         = "WriteThrottledEvents"
  namespace           = "AWS/DynamoDB"
  period              = "300"
  statistic           = "Sum"
  threshold           = var.dynamodb_throttle_threshold
  alarm_description   = "This metric monitors DynamoDB write throttling for ${each.key} table"
  alarm_actions       = var.alarm_actions

  dimensions = {
    TableName = each.value
  }

  tags = var.tags
}

# CloudWatch Dashboard for Application Monitoring
resource "aws_cloudwatch_dashboard" "main" {
  count = var.enable_dashboard ? 1 : 0
  
  dashboard_name = "${local.name_prefix}-monitoring"

  dashboard_body = jsonencode({
    widgets = [
      {
        type   = "metric"
        x      = 0
        y      = 0
        width  = 12
        height = 6

        properties = {
          metrics = [
            ["AWS/ECS", "CPUUtilization", "ServiceName", var.ecs_service_name, "ClusterName", var.ecs_cluster_name],
            [".", "MemoryUtilization", ".", ".", ".", "."],
            [".", "RunningTaskCount", ".", ".", ".", "."]
          ]
          view    = "timeSeries"
          stacked = false
          region  = data.aws_region.current.name
          title   = "ECS Service Metrics"
          period  = 300
        }
      },
      {
        type   = "metric"
        x      = 0
        y      = 6
        width  = 12
        height = 6

        properties = {
          metrics = concat([
            for table_key, table_name in var.dynamodb_table_names : [
              "AWS/DynamoDB", "ConsumedReadCapacityUnits", "TableName", table_name
            ]
          ], [
            for table_key, table_name in var.dynamodb_table_names : [
              "AWS/DynamoDB", "ConsumedWriteCapacityUnits", "TableName", table_name
            ]
          ])
          view    = "timeSeries"
          stacked = false
          region  = data.aws_region.current.name
          title   = "DynamoDB Capacity Metrics"
          period  = 300
        }
      }
    ]
  })
}

# CloudWatch Log Insights Queries for common troubleshooting
resource "aws_cloudwatch_query_definition" "api_errors" {
  name = "${local.name_prefix}-api-errors"

  log_group_names = [
    aws_cloudwatch_log_group.ecs_api.name
  ]

  query_string = <<EOF
fields @timestamp, @message
| filter @message like /ERROR/
| sort @timestamp desc
| limit 100
EOF
}

resource "aws_cloudwatch_query_definition" "api_performance" {
  name = "${local.name_prefix}-api-performance"

  log_group_names = [
    aws_cloudwatch_log_group.ecs_api.name
  ]

  query_string = <<EOF
fields @timestamp, @message
| filter @message like /response_time/
| stats avg(response_time) by bin(5m)
| sort @timestamp desc
EOF
}

resource "aws_cloudwatch_query_definition" "cli_task_status" {
  name = "${local.name_prefix}-cli-task-status"

  log_group_names = [
    for log_group in aws_cloudwatch_log_group.ecs_cli : log_group.name
  ]

  query_string = <<EOF
fields @timestamp, @message
| filter @message like /COMPLETED/ or @message like /FAILED/ or @message like /ERROR/
| sort @timestamp desc
| limit 50
EOF
}