# Terraform Module Integration Guide

This document describes how the Terraform modules are integrated in the root configuration and the data flow between modules.

## Module Architecture Overview

The Achievement Management infrastructure is composed of six main modules that work together to create a complete, scalable application platform:

```
┌─────────────────────────────────────────────────────────────────┐
│                        Root Module (main.tf)                    │
├─────────────────────────────────────────────────────────────────┤
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────┐ │
│  │     VPC     │  │  DynamoDB   │  │     IAM     │  │Security │ │
│  │   Module    │  │   Module    │  │   Module    │  │ Groups  │ │
│  │             │  │             │  │             │  │ Module  │ │
│  └─────────────┘  └─────────────┘  └─────────────┘  └─────────┘ │
│         │                 │                 │            │      │
│         └─────────────────┼─────────────────┼────────────┘      │
│                           │                 │                   │
│  ┌─────────────────────────┼─────────────────┼─────────────────┐ │
│  │                ECS Module                │                 │ │
│  │        (Cluster, Services, ALB)          │                 │ │
│  └─────────────────────────┼─────────────────┼─────────────────┘ │
│                           │                 │                   │
│  ┌─────────────────────────┼─────────────────┼─────────────────┐ │
│  │              Monitoring Module           │                 │ │
│  │         (CloudWatch, Logs, Alarms)       │                 │ │
│  └─────────────────────────┼─────────────────┼─────────────────┘ │
└─────────────────────────────────────────────────────────────────┘
```

## Module Dependencies and Data Flow

### 1. Foundation Layer (Level 1)

#### VPC Module
- **Purpose**: Creates the foundational networking infrastructure
- **Dependencies**: None (foundation layer)
- **Outputs Used By**: Security, ECS modules
- **Key Resources**: VPC, Subnets, Internet Gateway, NAT Gateways

#### DynamoDB Module
- **Purpose**: Creates application data storage tables
- **Dependencies**: None (foundation layer)
- **Outputs Used By**: IAM, Monitoring modules
- **Key Resources**: DynamoDB tables with encryption and backup

### 2. Security Layer (Level 2)

#### Security Groups Module
- **Purpose**: Creates network security rules
- **Dependencies**: VPC module
- **Inputs From VPC**: `vpc_id`, `vpc_cidr_block`
- **Outputs Used By**: ECS module
- **Key Resources**: ALB and ECS security groups

#### IAM Module
- **Purpose**: Creates roles and policies for ECS tasks
- **Dependencies**: DynamoDB module
- **Inputs From DynamoDB**: `table_names`, `table_arns`
- **Outputs Used By**: ECS module
- **Key Resources**: Task execution roles, task roles, CLI roles

### 3. Compute Layer (Level 3)

#### ECS Module
- **Purpose**: Creates container orchestration infrastructure
- **Dependencies**: VPC, Security, IAM, DynamoDB modules
- **Inputs From Multiple Modules**:
  - VPC: `vpc_id`, `public_subnet_ids`, `private_subnet_ids`
  - Security: `alb_security_group_id`, `ecs_security_group_id`
  - IAM: `ecs_task_execution_role_arn`, `ecs_task_role_arn`, `cli_task_role_arns`
  - DynamoDB: `table_names` (for environment variables)
- **Outputs Used By**: Monitoring module
- **Key Resources**: ECS cluster, services, task definitions, ALB

### 4. Observability Layer (Level 4)

#### Monitoring Module
- **Purpose**: Creates logging, monitoring, and alerting infrastructure
- **Dependencies**: ECS, DynamoDB modules
- **Inputs From Multiple Modules**:
  - ECS: `cluster_name`, `service_name`, `alb_arn_suffix`, `target_group_arn_suffix`
  - DynamoDB: `table_names`, `table_arns`
- **Key Resources**: CloudWatch log groups, dashboards, alarms

## Data Flow Patterns

### 1. Resource Identification Flow
```
Variables → Local Values → Module Inputs → Resource Creation → Outputs
```

### 2. Cross-Module Reference Flow
```
Module A Outputs → Local Values → Module B Inputs → Resource Dependencies
```

### 3. Configuration Inheritance Flow
```
Root Variables → Common Tags/Naming → All Modules → Consistent Resource Attributes
```

## Module Integration Implementation

### Local Values for Integration

The root module uses local values to:
- Standardize resource naming across modules
- Apply consistent tagging strategy
- Compute derived configuration values
- Validate configuration consistency

```hcl
locals {
  name_prefix = "${var.app_name}-${var.environment}"
  
  common_tags = {
    Environment = var.environment
    Application = var.app_name
    # ... additional tags
  }
  
  # Derived values for cross-module integration
  dynamodb_table_names = keys(var.dynamodb_tables)
  availability_zones = length(var.availability_zones) > 0 ? var.availability_zones : slice(data.aws_availability_zones.available.names, 0, 2)
}
```

### Explicit Dependencies

Dependencies are managed through:
1. **Implicit Dependencies**: Terraform automatically detects dependencies through resource references
2. **Explicit Dependencies**: Using `depends_on` for complex dependency chains

```hcl
module "ecs" {
  # ... configuration
  
  depends_on = [
    module.vpc,
    module.security,
    module.iam,
    module.dynamodb
  ]
}
```

### Variable Passing Patterns

#### Direct Reference Pattern
```hcl
module "security" {
  vpc_id = module.vpc.vpc_id
}
```

#### Computed Value Pattern
```hcl
module "iam" {
  dynamodb_table_names = local.dynamodb_table_names
  dynamodb_table_arns  = local.dynamodb_table_arns
}
```

#### Configuration Inheritance Pattern
```hcl
module "vpc" {
  name_prefix = local.name_prefix
  common_tags = local.common_tags
}
```

## Integration Validation

### Configuration Validation
The root module includes validation checks to ensure:
- Subnet counts match availability zone counts
- Required variables are provided
- Environment-specific constraints are met

### Output Validation
Special outputs provide integration status:
- `module_integration_status`: Verifies cross-module references
- `deployment_readiness`: Checks if infrastructure is ready for application deployment
- `cross_module_references`: Validates data flow between modules

### Runtime Validation
```hcl
resource "null_resource" "validate_configuration" {
  count = local.validate_subnet_count ? 0 : 1
  
  provisioner "local-exec" {
    command = "echo 'ERROR: Configuration validation failed' && exit 1"
  }
}
```

## Best Practices for Module Integration

### 1. Consistent Naming
- Use `name_prefix` for all resource names
- Apply consistent naming patterns across modules
- Include environment and application identifiers

### 2. Comprehensive Tagging
- Define common tags in root module
- Pass tags to all modules
- Include operational and compliance tags

### 3. Explicit Dependencies
- Use `depends_on` for complex dependency chains
- Document dependency relationships
- Ensure proper resource creation order

### 4. Configuration Validation
- Validate inputs at the root level
- Use variable validation rules
- Implement runtime checks for complex scenarios

### 5. Output Organization
- Provide outputs for external integration
- Include validation outputs for troubleshooting
- Document output usage patterns

## Troubleshooting Integration Issues

### Common Issues

#### Circular Dependencies
- **Symptom**: Terraform reports circular dependency errors
- **Solution**: Review module dependencies and use data sources where appropriate

#### Missing Resource References
- **Symptom**: Resources not found or null values
- **Solution**: Check module outputs and ensure proper dependency order

#### Inconsistent Configuration
- **Symptom**: Resources created with different naming or tagging
- **Solution**: Verify common_tags and name_prefix usage across modules

### Debugging Steps

1. **Check Module Outputs**: Verify that modules are producing expected outputs
2. **Validate Dependencies**: Ensure `depends_on` is used where needed
3. **Review Local Values**: Check that computed values are correct
4. **Use Validation Outputs**: Check `module_integration_status` output

### Integration Testing

Use the provided validation outputs to verify integration:

```bash
# Check integration status
terraform output module_integration_status

# Verify deployment readiness
terraform output deployment_readiness

# Review cross-module references
terraform output cross_module_references
```

## Future Enhancements

### Planned Improvements
1. **Module Versioning**: Pin module versions for stability
2. **Enhanced Validation**: Add more comprehensive validation checks
3. **Integration Tests**: Automated testing of module integration
4. **Documentation Generation**: Auto-generate integration documentation

### Extension Points
- Additional modules can be integrated following the same patterns
- New validation checks can be added to the root module
- Additional outputs can be provided for external integrations