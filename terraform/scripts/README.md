# Terraform Scripts

This directory contains utility scripts for managing the Achievement Management infrastructure.

## Application Management Scripts

### run-cli-task.sh

A helper script to run CLI tasks in ECS Fargate.

### Usage

```bash
# Make the script executable
chmod +x run-cli-task.sh

# Run CLI tasks
./run-cli-task.sh <environment> <task-type> <command> [args...]
```

### Examples

```bash
# Create an achievement
./run-cli-task.sh dev achievement create --title "First Login" --description "Log in for the first time" --point 10

# List achievements
./run-cli-task.sh dev achievement list

# Show current points
./run-cli-task.sh dev points current

# Create a reward
./run-cli-task.sh dev reward create --title "Coffee Voucher" --description "Free coffee" --point 50

# Redeem a reward
./run-cli-task.sh dev reward redeem --id "reward-id-here"
```

### Prerequisites

1. AWS CLI configured with appropriate permissions
2. Terraform configuration applied
3. jq installed for JSON parsing

### Task Types

- `achievement`: Manage achievements (create, list, update, delete)
- `points`: Manage points (current, aggregate, history)
- `reward`: Manage rewards (create, list, update, redeem, delete)

### Monitoring

After running a task, you can monitor its execution:

```bash
# View logs (replace with actual log group name)
aws logs tail /ecs/achievement-management-dev-cli-achievement --follow

# Check task status
aws ecs describe-tasks --cluster achievement-management-dev-cluster --tasks <task-arn>
```

## Notes

- All CLI tasks run in private subnets without public IP addresses
- Tasks use IAM roles for authentication (no hardcoded credentials)
- Logs are automatically sent to CloudWatch Logs
- Tasks have a default timeout of 5 minutes (configurable via Terraform variables)