# CloudWatch Monitoring Implementation

This document describes the CloudWatch monitoring implementation for the Achievement Management application infrastructure.

## Overview

The monitoring module provides comprehensive observability for the Achievement Management application deployed on AWS ECS with DynamoDB backend. It includes log aggregation, metrics collection, alerting, and dashboards for proactive monitoring and troubleshooting.

## Components Implemented

### 1. CloudWatch Log Groups

**API Service Logs**
- Log Group: `/ecs/{app_name}-{environment}-api`
- Purpose: Centralized logging for the API service container
- Retention: Environment-specific (7 days dev, 14 days staging, 30 days prod)

**CLI Task Logs**
- Log Groups: `/ecs/{app_name}-{environment}-cli-{task_type}`
- Task Types: achievement, points, reward
- Purpose: Separate log streams for different CLI operations
- Retention: Same as API service logs

**Log Configuration**
- Driver: `awslogs`
- Stream Prefix: `ecs` (with specific prefixes for CLI operations)
- Region: Automatically configured based on deployment region

### 2. CloudWatch Alarms

**ECS Service Monitoring**
- **CPU Utilization High**: Triggers when CPU usage exceeds threshold
  - Dev: 90%, Staging: 80%, Prod: 70%
  - Evaluation: 2 periods of 5 minutes
- **Memory Utilization High**: Triggers when memory usage exceeds threshold
  - Dev: 90%, Staging: 80%, Prod: 75%
  - Evaluation: 2 periods of 5 minutes
- **Task Count Low**: Triggers when running tasks fall below minimum
  - Threshold: 1 task (configurable per environment)
  - Evaluation: 2 periods of 5 minutes

**Application Load Balancer Monitoring**
- **Response Time High**: Monitors target response time
  - Dev: 5.0s, Staging: 3.0s, Prod: 2.0s
  - Evaluation: 2 periods of 5 minutes
- **HTTP 5XX Errors High**: Monitors server error rate
  - Dev: 20 errors, Staging: 10 errors, Prod: 5 errors
  - Evaluation: 2 periods of 5 minutes
- **Unhealthy Hosts**: Monitors target health
  - Threshold: 0 unhealthy hosts
  - Evaluation: 2 periods of 5 minutes

**DynamoDB Monitoring**
- **Read Throttle Events**: Monitors read capacity throttling
  - Dev: 5 events, Staging: 1 event, Prod: 0 events
  - Evaluation: 2 periods of 5 minutes
- **Write Throttle Events**: Monitors write capacity throttling
  - Dev: 5 events, Staging: 1 event, Prod: 0 events
  - Evaluation: 2 periods of 5 minutes

### 3. CloudWatch Dashboard

**Dashboard Name**: `{app_name}-{environment}-monitoring`

**Widgets Included**:
- ECS Service Metrics (CPU, Memory, Task Count)
- DynamoDB Capacity Metrics (Read/Write Capacity Units)
- ALB Performance Metrics (when enabled)

**Features**:
- Real-time metric visualization
- 5-minute data points
- Environment-specific configuration
- Direct links to detailed CloudWatch metrics

### 4. Log Insights Queries

**Pre-defined Queries**:
- **API Errors**: Filters and displays ERROR level logs from API service
- **API Performance**: Analyzes response time patterns (when logged)
- **CLI Task Status**: Tracks completion status of CLI tasks

**Query Features**:
- Optimized for common troubleshooting scenarios
- Time-based sorting
- Configurable result limits
- Cross-log-group analysis for CLI tasks

## Environment-Specific Configuration

### Development Environment
```hcl
log_retention_days = 7
monitoring_thresholds = {
  cpu_alarm_threshold         = 90
  memory_alarm_threshold      = 90
  response_time_threshold     = 5.0
  http_5xx_threshold          = 20
  dynamodb_throttle_threshold = 5
}
alarm_actions = []  # No notifications in dev
```

### Staging Environment
```hcl
log_retention_days = 14
monitoring_thresholds = {
  cpu_alarm_threshold         = 80
  memory_alarm_threshold      = 80
  response_time_threshold     = 3.0
  http_5xx_threshold          = 10
  dynamodb_throttle_threshold = 1
}
alarm_actions = []  # Can be configured with SNS topics
```

### Production Environment
```hcl
log_retention_days = 30
monitoring_thresholds = {
  cpu_alarm_threshold         = 70
  memory_alarm_threshold      = 75
  response_time_threshold     = 2.0
  http_5xx_threshold          = 5
  dynamodb_throttle_threshold = 0
}
alarm_actions = [
  # "arn:aws:sns:us-east-1:123456789012:production-alerts"
]
```

## Integration with ECS Module

The monitoring module integrates seamlessly with the ECS module:

1. **Log Group Creation**: Monitoring module creates log groups that ECS tasks reference
2. **Service Monitoring**: Alarms monitor ECS services created by the ECS module
3. **Dependency Management**: Proper Terraform dependencies ensure correct creation order

## Usage Examples

### Deploying with Monitoring

```bash
# Initialize Terraform
terraform init -backend-config=backend-dev.hcl

# Plan deployment with monitoring
terraform plan -var-file=environments/dev.tfvars

# Apply configuration
terraform apply -var-file=environments/dev.tfvars
```

### Accessing Monitoring Resources

**CloudWatch Dashboard**:
```bash
# Get dashboard URL from Terraform output
terraform output cloudwatch_dashboard_url
```

**Log Groups**:
```bash
# View API logs
aws logs describe-log-groups --log-group-name-prefix "/ecs/achievement-management-dev-api"

# Stream CLI logs
aws logs tail "/ecs/achievement-management-dev-cli-achievement" --follow
```

**Alarms**:
```bash
# List all alarms
aws cloudwatch describe-alarms --alarm-name-prefix "achievement-management-dev"

# Get alarm history
aws cloudwatch describe-alarm-history --alarm-name "achievement-management-dev-ecs-cpu-high"
```

### Running Log Insights Queries

```bash
# Start query for API errors
aws logs start-query \
  --log-group-name "/ecs/achievement-management-dev-api" \
  --start-time $(date -d '1 hour ago' +%s) \
  --end-time $(date +%s) \
  --query-string 'fields @timestamp, @message | filter @message like /ERROR/ | sort @timestamp desc | limit 100'
```

## Alarm Actions Configuration

### Setting up SNS Notifications

1. **Create SNS Topic**:
```bash
aws sns create-topic --name achievement-management-alerts
```

2. **Subscribe to Topic**:
```bash
aws sns subscribe \
  --topic-arn arn:aws:sns:us-east-1:123456789012:achievement-management-alerts \
  --protocol email \
  --notification-endpoint admin@example.com
```

3. **Update Environment Variables**:
```hcl
alarm_actions = [
  "arn:aws:sns:us-east-1:123456789012:achievement-management-alerts"
]
```

### Integration with PagerDuty or Slack

The alarm actions can be configured to integrate with various notification systems:
- SNS → Lambda → PagerDuty
- SNS → Lambda → Slack
- SNS → SQS → Custom processing

## Cost Optimization

### Log Retention Strategy
- **Development**: 7 days (minimal cost)
- **Staging**: 14 days (moderate retention)
- **Production**: 30 days (compliance balance)

### Alarm Optimization
- Environment-specific thresholds reduce false positives
- Consolidated alarms reduce CloudWatch costs
- Proper evaluation periods prevent alarm flapping

### Dashboard Efficiency
- Single dashboard per environment
- Focused on essential metrics
- Automatic refresh intervals

## Troubleshooting Guide

### Common Issues

**1. Log Groups Not Created**
- Check Terraform dependencies
- Verify IAM permissions for log group creation
- Ensure monitoring module is applied before ECS tasks start

**2. Alarms Not Triggering**
- Verify metric names and namespaces
- Check alarm thresholds against actual metrics
- Ensure sufficient data points for evaluation

**3. Dashboard Not Loading**
- Verify CloudWatch permissions
- Check metric availability in the region
- Ensure dashboard JSON is valid

### Debugging Commands

```bash
# Check log group existence
aws logs describe-log-groups --log-group-name-prefix "/ecs/achievement-management"

# Verify alarm configuration
aws cloudwatch describe-alarms --alarm-names "achievement-management-dev-ecs-cpu-high"

# Test metric availability
aws cloudwatch get-metric-statistics \
  --namespace AWS/ECS \
  --metric-name CPUUtilization \
  --dimensions Name=ServiceName,Value=achievement-management-dev-api-service \
  --start-time $(date -d '1 hour ago' --iso-8601) \
  --end-time $(date --iso-8601) \
  --period 300 \
  --statistics Average
```

## Best Practices

1. **Environment Separation**: Use different thresholds and retention policies per environment
2. **Alarm Actions**: Configure appropriate notification channels for each environment
3. **Log Analysis**: Regularly review Log Insights queries and optimize for common use cases
4. **Cost Management**: Monitor CloudWatch costs and adjust retention policies as needed
5. **Documentation**: Keep alarm descriptions clear and actionable
6. **Testing**: Regularly test alarm conditions to ensure they trigger correctly

## Future Enhancements

1. **Custom Metrics**: Add application-specific metrics
2. **Anomaly Detection**: Implement CloudWatch Anomaly Detection
3. **Cross-Region Monitoring**: Extend monitoring to multiple regions
4. **Advanced Dashboards**: Create role-specific dashboards
5. **Automated Remediation**: Implement Lambda-based auto-remediation
6. **Cost Alerting**: Add billing and cost monitoring alarms