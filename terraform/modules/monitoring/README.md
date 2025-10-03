# CloudWatch Monitoring Module

This module creates comprehensive CloudWatch monitoring resources for the Achievement Management application, including log groups, metrics, alarms, and dashboards.

## Features

- **CloudWatch Log Groups**: Centralized logging for ECS API service and CLI tasks with configurable retention policies
- **CloudWatch Alarms**: Proactive monitoring for ECS services, ALB, and DynamoDB with customizable thresholds
- **CloudWatch Dashboard**: Visual monitoring dashboard for key application metrics
- **Log Insights Queries**: Pre-defined queries for common troubleshooting scenarios

## Resources Created

### Log Groups
- `/ecs/{app_name}-{environment}-api` - API service logs
- `/ecs/{app_name}-{environment}-cli-{task_type}` - CLI task logs for each task type
- `/aws/applicationloadbalancer/{app_name}-{environment}-alb` - ALB access logs (optional)

### CloudWatch Alarms
- **ECS Service Monitoring**:
  - CPU utilization high
  - Memory utilization high
  - Running task count low
- **ALB Monitoring**:
  - Response time high
  - HTTP 5XX errors high
  - Unhealthy host count
- **DynamoDB Monitoring**:
  - Read throttle events
  - Write throttle events

### Dashboard
- ECS service metrics (CPU, Memory, Task Count)
- DynamoDB capacity metrics
- ALB performance metrics (if enabled)

### Log Insights Queries
- API error analysis
- API performance analysis
- CLI task status tracking

## Usage

```hcl
module "monitoring" {
  source = "./modules/monitoring"

  environment = var.environment
  app_name    = var.app_name

  # ECS Configuration
  ecs_service_name = module.ecs.service_name
  ecs_cluster_name = module.ecs.cluster_name

  # ALB Configuration
  alb_arn_suffix           = module.ecs.alb_arn_suffix
  target_group_arn_suffix  = module.ecs.target_group_arn_suffix

  # DynamoDB Configuration
  dynamodb_table_names = module.dynamodb.table_names

  # Monitoring Configuration
  log_retention_days = var.log_retention_days
  enable_dashboard   = true

  # Alarm Thresholds
  cpu_alarm_threshold    = 80
  memory_alarm_threshold = 80

  # Alarm Actions (SNS topics, etc.)
  alarm_actions = var.alarm_actions

  tags = local.common_tags
}
```

## Variables

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| environment | Environment name (dev, staging, prod) | `string` | n/a | yes |
| app_name | Application name used for resource naming | `string` | n/a | yes |
| log_retention_days | CloudWatch log retention period in days | `number` | `14` | no |
| ecs_service_name | Name of the ECS service to monitor | `string` | n/a | yes |
| ecs_cluster_name | Name of the ECS cluster to monitor | `string` | n/a | yes |
| dynamodb_table_names | Map of DynamoDB table logical names to actual table names | `map(string)` | `{}` | no |
| cpu_alarm_threshold | CPU utilization threshold for alarms (percentage) | `number` | `80` | no |
| memory_alarm_threshold | Memory utilization threshold for alarms (percentage) | `number` | `80` | no |
| enable_dashboard | Enable CloudWatch dashboard creation | `bool` | `true` | no |
| alarm_actions | List of ARNs to notify when alarm triggers | `list(string)` | `[]` | no |

## Outputs

| Name | Description |
|------|-------------|
| api_log_group_name | Name of the API CloudWatch log group |
| cli_log_group_names | Names of the CLI CloudWatch log groups |
| dashboard_url | URL of the CloudWatch dashboard |
| ecs_cpu_alarm_arn | ARN of the ECS CPU utilization alarm |
| ecs_memory_alarm_arn | ARN of the ECS memory utilization alarm |

## Alarm Thresholds

The module provides sensible defaults for alarm thresholds, but they can be customized based on your application's requirements:

- **CPU Utilization**: 80% (triggers when ECS service CPU usage exceeds this threshold)
- **Memory Utilization**: 80% (triggers when ECS service memory usage exceeds this threshold)
- **Response Time**: 2 seconds (triggers when ALB response time exceeds this threshold)
- **HTTP 5XX Errors**: 10 errors (triggers when 5XX error count exceeds this threshold)
- **DynamoDB Throttling**: 0 events (triggers on any throttling events)

## Log Retention

Log retention is configurable per environment:
- **Development**: 7 days (cost optimization)
- **Staging**: 14 days (moderate retention)
- **Production**: 30+ days (compliance and troubleshooting)

## Integration with Other Modules

This module is designed to work with:
- **ECS Module**: Monitors ECS services and tasks
- **DynamoDB Module**: Monitors database performance
- **ALB**: Monitors load balancer metrics

## Best Practices

1. **Environment-specific thresholds**: Adjust alarm thresholds based on environment requirements
2. **Alarm actions**: Configure SNS topics or other notification mechanisms for critical alarms
3. **Log retention**: Balance cost and compliance requirements when setting retention periods
4. **Dashboard customization**: Extend the dashboard with application-specific metrics
5. **Regular review**: Periodically review and adjust alarm thresholds based on application behavior