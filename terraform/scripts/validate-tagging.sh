#!/bin/bash

# Terraform Tagging and Naming Validation Script
# This script validates that all resources follow the established tagging and naming conventions

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Counters
TOTAL_CHECKS=0
PASSED_CHECKS=0
FAILED_CHECKS=0

# Function to print colored output
print_status() {
    local status=$1
    local message=$2
    
    TOTAL_CHECKS=$((TOTAL_CHECKS + 1))
    
    if [ "$status" = "PASS" ]; then
        echo -e "${GREEN}✓ PASS${NC}: $message"
        PASSED_CHECKS=$((PASSED_CHECKS + 1))
    elif [ "$status" = "FAIL" ]; then
        echo -e "${RED}✗ FAIL${NC}: $message"
        FAILED_CHECKS=$((FAILED_CHECKS + 1))
    elif [ "$status" = "WARN" ]; then
        echo -e "${YELLOW}⚠ WARN${NC}: $message"
    else
        echo -e "${NC}ℹ INFO${NC}: $message"
    fi
}

# Function to check if Terraform is installed
check_terraform() {
    if command -v terraform &> /dev/null; then
        print_status "PASS" "Terraform is installed"
        terraform version
    else
        print_status "FAIL" "Terraform is not installed"
        exit 1
    fi
}

# Function to validate Terraform syntax
validate_terraform_syntax() {
    echo -e "\n${YELLOW}Validating Terraform syntax...${NC}"
    
    if terraform validate; then
        print_status "PASS" "Terraform syntax is valid"
    else
        print_status "FAIL" "Terraform syntax validation failed"
    fi
}

# Function to check for required tagging variables
check_tagging_variables() {
    echo -e "\n${YELLOW}Checking required tagging variables...${NC}"
    
    local required_vars=(
        "environment"
        "app_name"
        "cost_center"
        "owner"
        "team"
        "business_unit"
        "service_tier"
        "backup_policy"
        "data_classification"
        "compliance_scope"
    )
    
    for var in "${required_vars[@]}"; do
        if grep -q "variable \"$var\"" variables.tf; then
            print_status "PASS" "Required variable '$var' is defined"
        else
            print_status "FAIL" "Required variable '$var' is missing"
        fi
    done
}

# Function to check for common_tags implementation
check_common_tags() {
    echo -e "\n${YELLOW}Checking common_tags implementation...${NC}"
    
    if grep -q "common_tags = {" main.tf; then
        print_status "PASS" "common_tags is defined in main.tf"
    else
        print_status "FAIL" "common_tags is not defined in main.tf"
    fi
    
    # Check for required tags in common_tags
    local required_tags=(
        "Environment"
        "Application"
        "Project"
        "ManagedBy"
        "Owner"
        "Team"
        "CostCenter"
    )
    
    for tag in "${required_tags[@]}"; do
        if grep -A 20 "common_tags = {" main.tf | grep -q "$tag"; then
            print_status "PASS" "Required tag '$tag' is in common_tags"
        else
            print_status "FAIL" "Required tag '$tag' is missing from common_tags"
        fi
    done
}

# Function to check naming conventions in modules
check_naming_conventions() {
    echo -e "\n${YELLOW}Checking naming conventions in modules...${NC}"
    
    # Check VPC module naming
    if grep -q 'Name = "${var.app_name}-${var.environment}-vpc"' modules/vpc/main.tf; then
        print_status "PASS" "VPC naming convention is correct"
    else
        print_status "FAIL" "VPC naming convention is incorrect"
    fi
    
    # Check ECS cluster naming
    if grep -q 'name = "${local.name_prefix}-cluster"' modules/ecs/main.tf; then
        print_status "PASS" "ECS cluster naming convention is correct"
    else
        print_status "FAIL" "ECS cluster naming convention is incorrect"
    fi
    
    # Check DynamoDB table naming
    if grep -q 'name = "${local.name_prefix}-${each.key}"' modules/dynamodb/main.tf; then
        print_status "PASS" "DynamoDB table naming convention is correct"
    else
        print_status "FAIL" "DynamoDB table naming convention is incorrect"
    fi
}

# Function to check resource-specific tags
check_resource_specific_tags() {
    echo -e "\n${YELLOW}Checking resource-specific tags...${NC}"
    
    # Check if ResourceType tag is used
    if grep -r "ResourceType" modules/*/main.tf | wc -l | grep -q "[1-9]"; then
        print_status "PASS" "ResourceType tags are implemented"
    else
        print_status "FAIL" "ResourceType tags are missing"
    fi
    
    # Check if merge function is used for tags
    if grep -r "merge(.*tags" modules/*/main.tf | wc -l | grep -q "[1-9]"; then
        print_status "PASS" "Tag merge function is used in modules"
    else
        print_status "FAIL" "Tag merge function is not used consistently"
    fi
}

# Function to check environment-specific configurations
check_environment_configs() {
    echo -e "\n${YELLOW}Checking environment-specific configurations...${NC}"
    
    local environments=("dev" "staging" "prod")
    
    for env in "${environments[@]}"; do
        if [ -f "environments/${env}.tfvars" ]; then
            print_status "PASS" "Environment file for $env exists"
            
            # Check if environment-specific tagging variables are set
            if grep -q "owner.*=" "environments/${env}.tfvars"; then
                print_status "PASS" "Owner is set in $env environment"
            else
                print_status "FAIL" "Owner is not set in $env environment"
            fi
            
            if grep -q "service_tier.*=" "environments/${env}.tfvars"; then
                print_status "PASS" "Service tier is set in $env environment"
            else
                print_status "FAIL" "Service tier is not set in $env environment"
            fi
        else
            print_status "FAIL" "Environment file for $env is missing"
        fi
    done
}

# Function to check for tag validation rules
check_tag_validation() {
    echo -e "\n${YELLOW}Checking tag validation rules...${NC}"
    
    # Check for service_tier validation
    if grep -A 5 'variable "service_tier"' variables.tf | grep -q "validation"; then
        print_status "PASS" "Service tier validation is implemented"
    else
        print_status "FAIL" "Service tier validation is missing"
    fi
    
    # Check for backup_policy validation
    if grep -A 5 'variable "backup_policy"' variables.tf | grep -q "validation"; then
        print_status "PASS" "Backup policy validation is implemented"
    else
        print_status "FAIL" "Backup policy validation is missing"
    fi
    
    # Check for data_classification validation
    if grep -A 5 'variable "data_classification"' variables.tf | grep -q "validation"; then
        print_status "PASS" "Data classification validation is implemented"
    else
        print_status "FAIL" "Data classification validation is missing"
    fi
}

# Function to check documentation
check_documentation() {
    echo -e "\n${YELLOW}Checking documentation...${NC}"
    
    if [ -f "TAGGING_STRATEGY.md" ]; then
        print_status "PASS" "Tagging strategy documentation exists"
    else
        print_status "FAIL" "Tagging strategy documentation is missing"
    fi
    
    if [ -f "RESOURCE_NAMING.md" ]; then
        print_status "PASS" "Resource naming documentation exists"
    else
        print_status "FAIL" "Resource naming documentation is missing"
    fi
}

# Function to run terraform plan and check for issues
check_terraform_plan() {
    echo -e "\n${YELLOW}Running terraform plan to check for issues...${NC}"
    
    # Initialize terraform if needed
    if [ ! -d ".terraform" ]; then
        print_status "INFO" "Initializing Terraform..."
        terraform init -backend=false
    fi
    
    # Run terraform plan
    if terraform plan -var-file="environments/dev.tfvars" -out=plan.out > /dev/null 2>&1; then
        print_status "PASS" "Terraform plan completed successfully"
    else
        print_status "WARN" "Terraform plan had issues (may be due to missing AWS credentials)"
    fi
}

# Main execution
main() {
    echo -e "${GREEN}=== Terraform Tagging and Naming Validation ===${NC}\n"
    
    # Change to terraform directory if not already there
    if [ ! -f "main.tf" ]; then
        if [ -f "terraform/main.tf" ]; then
            cd terraform
        else
            echo -e "${RED}Error: Cannot find main.tf file${NC}"
            exit 1
        fi
    fi
    
    # Run all checks
    check_terraform
    validate_terraform_syntax
    check_tagging_variables
    check_common_tags
    check_naming_conventions
    check_resource_specific_tags
    check_environment_configs
    check_tag_validation
    check_documentation
    check_terraform_plan
    
    # Print summary
    echo -e "\n${YELLOW}=== Validation Summary ===${NC}"
    echo -e "Total checks: $TOTAL_CHECKS"
    echo -e "${GREEN}Passed: $PASSED_CHECKS${NC}"
    echo -e "${RED}Failed: $FAILED_CHECKS${NC}"
    
    if [ $FAILED_CHECKS -eq 0 ]; then
        echo -e "\n${GREEN}🎉 All validation checks passed!${NC}"
        exit 0
    else
        echo -e "\n${RED}❌ Some validation checks failed. Please review and fix the issues.${NC}"
        exit 1
    fi
}

# Run main function
main "$@"