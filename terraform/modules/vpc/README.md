# VPC Module

This module creates a VPC with public and private subnets across multiple availability zones, along with the necessary networking components for a secure and scalable infrastructure.

## Features

- VPC with DNS support enabled
- Public and private subnets across multiple AZs
- Internet Gateway for public subnet internet access
- NAT Gateways for private subnet outbound internet access
- Route tables and associations for proper traffic routing
- Configurable subnet counts and CIDR blocks
- Comprehensive resource tagging

## Architecture

```
Internet
    |
Internet Gateway
    |
Public Subnets (Multi-AZ)
    |
NAT Gateways
    |
Private Subnets (Multi-AZ)
```

## Usage

```hcl
module "vpc" {
  source = "./modules/vpc"

  environment = "dev"
  app_name    = "achievement-management"
  
  vpc_cidr               = "10.0.0.0/16"
  public_subnet_cidrs    = ["10.0.1.0/24", "10.0.2.0/24"]
  private_subnet_cidrs   = ["10.0.10.0/24", "10.0.20.0/24"]
  public_subnet_count    = 2
  private_subnet_count   = 2
  nat_gateway_count      = 2

  common_tags = {
    Environment = "dev"
    Project     = "achievement-management"
  }
}
```

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| environment | Environment name (dev, staging, prod) | `string` | n/a | yes |
| app_name | Application name | `string` | n/a | yes |
| vpc_cidr | CIDR block for VPC | `string` | `"10.0.0.0/16"` | no |
| public_subnet_count | Number of public subnets | `number` | `2` | no |
| private_subnet_count | Number of private subnets | `number` | `2` | no |
| public_subnet_cidrs | CIDR blocks for public subnets | `list(string)` | `["10.0.1.0/24", "10.0.2.0/24"]` | no |
| private_subnet_cidrs | CIDR blocks for private subnets | `list(string)` | `["10.0.10.0/24", "10.0.20.0/24"]` | no |
| nat_gateway_count | Number of NAT gateways | `number` | `2` | no |
| common_tags | Common tags to be applied to all resources | `map(string)` | `{}` | no |

## Outputs

| Name | Description |
|------|-------------|
| vpc_id | ID of the VPC |
| vpc_cidr_block | CIDR block of the VPC |
| internet_gateway_id | ID of the Internet Gateway |
| public_subnet_ids | IDs of the public subnets |
| private_subnet_ids | IDs of the private subnets |
| public_subnet_cidrs | CIDR blocks of the public subnets |
| private_subnet_cidrs | CIDR blocks of the private subnets |
| nat_gateway_ids | IDs of the NAT Gateways |
| nat_gateway_public_ips | Public IPs of the NAT Gateways |
| public_route_table_id | ID of the public route table |
| private_route_table_ids | IDs of the private route tables |
| availability_zones | List of availability zones used |

## Requirements

| Name | Version |
|------|---------|
| terraform | >= 1.0 |
| aws | >= 5.0 |

## Resources Created

- 1 VPC
- 1 Internet Gateway
- 2 Public Subnets (default)
- 2 Private Subnets (default)
- 2 NAT Gateways (default)
- 2 Elastic IPs for NAT Gateways
- 1 Public Route Table
- 2 Private Route Tables (default)
- Route Table Associations

## Cost Considerations

- NAT Gateways incur hourly charges and data processing fees
- Elastic IPs are free when associated with running instances
- Consider using a single NAT Gateway for cost optimization in non-production environments