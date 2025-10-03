#!/bin/bash

# Script to run CLI tasks in ECS
# Usage: ./run-cli-task.sh <environment> <task-type> <command> [args...]

set -e

# Check if required arguments are provided
if [ $# -lt 3 ]; then
    echo "Usage: $0 <environment> <task-type> <command> [args...]"
    echo ""
    echo "Examples:"
    echo "  $0 dev achievement create --title \"Test Achievement\" --description \"Test\" --point 10"
    echo "  $0 dev points current"
    echo "  $0 dev reward list"
    echo ""
    echo "Task types: achievement, points, reward"
    exit 1
fi

ENVIRONMENT=$1
TASK_TYPE=$2
COMMAND=$3
shift 3
ARGS="$@"

# Configuration
APP_NAME="achievement-management"
CLUSTER_NAME="${APP_NAME}-${ENVIRONMENT}-cluster"
TASK_DEFINITION="${APP_NAME}-${ENVIRONMENT}-cli-${TASK_TYPE}"

# Get network configuration from Terraform outputs
echo "Getting network configuration..."
TERRAFORM_OUTPUT=$(terraform output -json)

if [ $? -ne 0 ]; then
    echo "Error: Failed to get Terraform outputs. Make sure you're in the terraform directory and have applied the configuration."
    exit 1
fi

# Extract network configuration
SUBNET_IDS=$(echo "$TERRAFORM_OUTPUT" | jq -r '.ecs_task_execution_info.value.subnet_ids | join(",")')
SECURITY_GROUP_ID=$(echo "$TERRAFORM_OUTPUT" | jq -r '.ecs_task_execution_info.value.security_group_id')

if [ "$SUBNET_IDS" = "null" ] || [ "$SECURITY_GROUP_ID" = "null" ]; then
    echo "Error: Could not extract network configuration from Terraform outputs."
    exit 1
fi

echo "Cluster: $CLUSTER_NAME"
echo "Task Definition: $TASK_DEFINITION"
echo "Command: ./achievement-app $TASK_TYPE $COMMAND $ARGS"
echo "Subnets: $SUBNET_IDS"
echo "Security Group: $SECURITY_GROUP_ID"
echo ""

# Build the command override
if [ -n "$ARGS" ]; then
    FULL_COMMAND="[\"./achievement-app\", \"$TASK_TYPE\", \"$COMMAND\", $ARGS]"
else
    FULL_COMMAND="[\"./achievement-app\", \"$TASK_TYPE\", \"$COMMAND\"]"
fi

# Convert space-separated args to JSON array format
if [ -n "$ARGS" ]; then
    # Split args and format as JSON array elements
    JSON_ARGS=""
    for arg in $ARGS; do
        if [ -n "$JSON_ARGS" ]; then
            JSON_ARGS="$JSON_ARGS, \"$arg\""
        else
            JSON_ARGS="\"$arg\""
        fi
    done
    FULL_COMMAND="[\"./achievement-app\", \"$TASK_TYPE\", \"$COMMAND\", $JSON_ARGS]"
else
    FULL_COMMAND="[\"./achievement-app\", \"$TASK_TYPE\", \"$COMMAND\"]"
fi

# Create the overrides JSON
OVERRIDES=$(cat <<EOF
{
  "containerOverrides": [
    {
      "name": "cli-$TASK_TYPE",
      "command": $FULL_COMMAND
    }
  ]
}
EOF
)

echo "Running ECS task..."
echo "Command override: $FULL_COMMAND"
echo ""

# Run the ECS task
TASK_ARN=$(aws ecs run-task \
    --cluster "$CLUSTER_NAME" \
    --task-definition "$TASK_DEFINITION" \
    --launch-type FARGATE \
    --network-configuration "awsvpcConfiguration={subnets=[$SUBNET_IDS],securityGroups=[$SECURITY_GROUP_ID],assignPublicIp=DISABLED}" \
    --overrides "$OVERRIDES" \
    --query 'tasks[0].taskArn' \
    --output text)

if [ $? -ne 0 ]; then
    echo "Error: Failed to run ECS task."
    exit 1
fi

echo "Task started successfully!"
echo "Task ARN: $TASK_ARN"
echo ""
echo "To view logs:"
echo "  aws logs tail /ecs/${APP_NAME}-${ENVIRONMENT}-cli-${TASK_TYPE} --follow"
echo ""
echo "To check task status:"
echo "  aws ecs describe-tasks --cluster $CLUSTER_NAME --tasks $TASK_ARN"