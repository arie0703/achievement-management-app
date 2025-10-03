# CLI Task Usage Examples

This document provides comprehensive examples of how to use the ECS CLI tasks for the Achievement Management application.

## Prerequisites

1. Terraform infrastructure deployed
2. AWS CLI configured with appropriate permissions
3. Container image built and pushed to ECR (or other registry)

## Getting Started

### 1. Get Infrastructure Information

First, get the necessary information from Terraform outputs:

```bash
# Get all outputs
terraform output

# Get specific task execution info
terraform output task_execution_info

# Get CLI task definition ARNs
terraform output cli_task_definition_arns
```

### 2. Basic Task Execution

Use the AWS CLI to run ECS tasks:

```bash
# Basic template
aws ecs run-task \
  --cluster <cluster-name> \
  --task-definition <task-definition-arn> \
  --launch-type FARGATE \
  --network-configuration "awsvpcConfiguration={subnets=[<subnet-ids>],securityGroups=[<security-group-id>],assignPublicIp=DISABLED}" \
  --overrides '<command-overrides>'
```

## Achievement Management Examples

### Create Achievement

```bash
# Using general CLI task with command override
aws ecs run-task \
  --cluster achievement-management-dev-cluster \
  --task-definition achievement-management-dev-cli-achievement \
  --launch-type FARGATE \
  --network-configuration "awsvpcConfiguration={subnets=[subnet-xxx,subnet-yyy],securityGroups=[sg-xxx],assignPublicIp=DISABLED}" \
  --overrides '{
    "containerOverrides": [
      {
        "name": "cli-achievement",
        "command": ["./achievement-app", "achievement", "create", "--title", "First Login", "--description", "Log in for the first time", "--point", "10"]
      }
    ]
  }'

# Using specific CLI task (if enabled)
aws ecs run-task \
  --cluster achievement-management-dev-cluster \
  --task-definition achievement-management-dev-cli-achievement-create \
  --launch-type FARGATE \
  --network-configuration "awsvpcConfiguration={subnets=[subnet-xxx,subnet-yyy],securityGroups=[sg-xxx],assignPublicIp=DISABLED}" \
  --overrides '{
    "containerOverrides": [
      {
        "name": "cli-achievement-create",
        "command": ["./achievement-app", "achievement", "create", "--title", "First Login", "--description", "Log in for the first time", "--point", "10"]
      }
    ]
  }'
```

### List Achievements

```bash
aws ecs run-task \
  --cluster achievement-management-dev-cluster \
  --task-definition achievement-management-dev-cli-achievement \
  --launch-type FARGATE \
  --network-configuration "awsvpcConfiguration={subnets=[subnet-xxx,subnet-yyy],securityGroups=[sg-xxx],assignPublicIp=DISABLED}" \
  --overrides '{
    "containerOverrides": [
      {
        "name": "cli-achievement",
        "command": ["./achievement-app", "achievement", "list"]
      }
    ]
  }'
```

### Update Achievement

```bash
aws ecs run-task \
  --cluster achievement-management-dev-cluster \
  --task-definition achievement-management-dev-cli-achievement \
  --launch-type FARGATE \
  --network-configuration "awsvpcConfiguration={subnets=[subnet-xxx,subnet-yyy],securityGroups=[sg-xxx],assignPublicIp=DISABLED}" \
  --overrides '{
    "containerOverrides": [
      {
        "name": "cli-achievement",
        "command": ["./achievement-app", "achievement", "update", "--id", "achievement-id-here", "--title", "Updated Title", "--point", "20"]
      }
    ]
  }'
```

### Delete Achievement

```bash
aws ecs run-task \
  --cluster achievement-management-dev-cluster \
  --task-definition achievement-management-dev-cli-achievement \
  --launch-type FARGATE \
  --network-configuration "awsvpcConfiguration={subnets=[subnet-xxx,subnet-yyy],securityGroups=[sg-xxx],assignPublicIp=DISABLED}" \
  --overrides '{
    "containerOverrides": [
      {
        "name": "cli-achievement",
        "command": ["./achievement-app", "achievement", "delete", "--id", "achievement-id-here"]
      }
    ]
  }'
```

## Points Management Examples

### Check Current Points

```bash
aws ecs run-task \
  --cluster achievement-management-dev-cluster \
  --task-definition achievement-management-dev-cli-points \
  --launch-type FARGATE \
  --network-configuration "awsvpcConfiguration={subnets=[subnet-xxx,subnet-yyy],securityGroups=[sg-xxx],assignPublicIp=DISABLED}" \
  --overrides '{
    "containerOverrides": [
      {
        "name": "cli-points",
        "command": ["./achievement-app", "points", "current"]
      }
    ]
  }'
```

### Aggregate Points

```bash
aws ecs run-task \
  --cluster achievement-management-dev-cluster \
  --task-definition achievement-management-dev-cli-points \
  --launch-type FARGATE \
  --network-configuration "awsvpcConfiguration={subnets=[subnet-xxx,subnet-yyy],securityGroups=[sg-xxx],assignPublicIp=DISABLED}" \
  --overrides '{
    "containerOverrides": [
      {
        "name": "cli-points",
        "command": ["./achievement-app", "points", "aggregate"]
      }
    ]
  }'
```

### View Reward History

```bash
aws ecs run-task \
  --cluster achievement-management-dev-cluster \
  --task-definition achievement-management-dev-cli-points \
  --launch-type FARGATE \
  --network-configuration "awsvpcConfiguration={subnets=[subnet-xxx,subnet-yyy],securityGroups=[sg-xxx],assignPublicIp=DISABLED}" \
  --overrides '{
    "containerOverrides": [
      {
        "name": "cli-points",
        "command": ["./achievement-app", "points", "history"]
      }
    ]
  }'
```

## Reward Management Examples

### Create Reward

```bash
aws ecs run-task \
  --cluster achievement-management-dev-cluster \
  --task-definition achievement-management-dev-cli-reward \
  --launch-type FARGATE \
  --network-configuration "awsvpcConfiguration={subnets=[subnet-xxx,subnet-yyy],securityGroups=[sg-xxx],assignPublicIp=DISABLED}" \
  --overrides '{
    "containerOverrides": [
      {
        "name": "cli-reward",
        "command": ["./achievement-app", "reward", "create", "--title", "Coffee Voucher", "--description", "Free coffee at the office", "--point", "50"]
      }
    ]
  }'
```

### List Rewards

```bash
aws ecs run-task \
  --cluster achievement-management-dev-cluster \
  --task-definition achievement-management-dev-cli-reward \
  --launch-type FARGATE \
  --network-configuration "awsvpcConfiguration={subnets=[subnet-xxx,subnet-yyy],securityGroups=[sg-xxx],assignPublicIp=DISABLED}" \
  --overrides '{
    "containerOverrides": [
      {
        "name": "cli-reward",
        "command": ["./achievement-app", "reward", "list"]
      }
    ]
  }'
```

### Redeem Reward

```bash
aws ecs run-task \
  --cluster achievement-management-dev-cluster \
  --task-definition achievement-management-dev-cli-reward \
  --launch-type FARGATE \
  --network-configuration "awsvpcConfiguration={subnets=[subnet-xxx,subnet-yyy],securityGroups=[sg-xxx],assignPublicIp=DISABLED}" \
  --overrides '{
    "containerOverrides": [
      {
        "name": "cli-reward",
        "command": ["./achievement-app", "reward", "redeem", "--id", "reward-id-here"]
      }
    ]
  }'
```

## Using the Helper Script

For easier task execution, use the provided helper script:

```bash
# Make it executable
chmod +x terraform/scripts/run-cli-task.sh

# Examples
./terraform/scripts/run-cli-task.sh dev achievement create --title "Test Achievement" --description "Test" --point 10
./terraform/scripts/run-cli-task.sh dev points current
./terraform/scripts/run-cli-task.sh dev reward list
```

## Monitoring Task Execution

### View Task Status

```bash
# Get task ARN from the run-task command output, then:
aws ecs describe-tasks \
  --cluster achievement-management-dev-cluster \
  --tasks <task-arn>
```

### View Logs

```bash
# View logs in real-time
aws logs tail /ecs/achievement-management-dev-cli-achievement --follow

# View logs for specific time range
aws logs filter-log-events \
  --log-group-name /ecs/achievement-management-dev-cli-achievement \
  --start-time 1640995200000 \
  --end-time 1640998800000
```

### List Recent Tasks

```bash
# List recent tasks for a service
aws ecs list-tasks \
  --cluster achievement-management-dev-cluster \
  --family achievement-management-dev-cli-achievement
```

## Automation Examples

### Using in CI/CD Pipeline

```yaml
# GitHub Actions example
- name: Run Achievement CLI Task
  run: |
    aws ecs run-task \
      --cluster ${{ env.CLUSTER_NAME }} \
      --task-definition ${{ env.TASK_DEFINITION }} \
      --launch-type FARGATE \
      --network-configuration "awsvpcConfiguration={subnets=[${{ env.SUBNET_IDS }}],securityGroups=[${{ env.SECURITY_GROUP_ID }}],assignPublicIp=DISABLED}" \
      --overrides '{
        "containerOverrides": [
          {
            "name": "cli-achievement",
            "command": ["./achievement-app", "achievement", "create", "--title", "CI Achievement", "--description", "Created by CI", "--point", "5"]
          }
        ]
      }'
```

### Using with Terraform

```hcl
# Run a CLI task as part of Terraform deployment
resource "null_resource" "create_initial_achievement" {
  provisioner "local-exec" {
    command = <<-EOT
      aws ecs run-task \
        --cluster ${module.ecs.cluster_name} \
        --task-definition ${module.ecs.cli_task_definition_arns.achievement} \
        --launch-type FARGATE \
        --network-configuration "awsvpcConfiguration={subnets=[${join(",", module.ecs.task_execution_info.subnet_ids)}],securityGroups=[${module.ecs.task_execution_info.security_group_id}],assignPublicIp=DISABLED}" \
        --overrides '{
          "containerOverrides": [
            {
              "name": "cli-achievement",
              "command": ["./achievement-app", "achievement", "create", "--title", "Welcome", "--description", "Welcome to the system", "--point", "10"]
            }
          ]
        }'
    EOT
  }

  depends_on = [module.ecs]
}
```

## Troubleshooting

### Common Issues

1. **Task fails to start**
   - Check IAM permissions
   - Verify network configuration
   - Ensure container image is accessible

2. **Command not found**
   - Verify the CLI binary exists in the container image
   - Check the working directory and PATH

3. **Database connection errors**
   - Verify security group rules
   - Check IAM permissions for DynamoDB
   - Ensure DynamoDB tables exist

4. **Task timeout**
   - Increase task timeout in Terraform variables
   - Optimize CLI operations
   - Check for infinite loops or hanging operations

### Debugging Commands

```bash
# Get detailed task information
aws ecs describe-tasks --cluster <cluster> --tasks <task-arn> --include TAGS

# Get task definition details
aws ecs describe-task-definition --task-definition <task-definition-arn>

# Check CloudWatch logs
aws logs describe-log-groups --log-group-name-prefix /ecs/achievement-management

# Test network connectivity (if needed)
aws ecs run-task \
  --cluster <cluster> \
  --task-definition <task-definition> \
  --launch-type FARGATE \
  --network-configuration "..." \
  --overrides '{
    "containerOverrides": [
      {
        "name": "cli-achievement",
        "command": ["ping", "-c", "3", "dynamodb.us-east-1.amazonaws.com"]
      }
    ]
  }'
```

## Best Practices

1. **Use specific task definitions** when available for better performance
2. **Monitor task execution** through CloudWatch Logs and metrics
3. **Set appropriate timeouts** for long-running operations
4. **Use IAM roles** instead of hardcoded credentials
5. **Tag resources** for better cost tracking and management
6. **Test in development** before running in production
7. **Implement retry logic** for critical operations
8. **Use environment variables** for configuration instead of hardcoded values