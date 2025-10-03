# VPC
resource "aws_vpc" "main" {
  cidr_block           = var.vpc_cidr
  enable_dns_hostnames = true
  enable_dns_support   = true

  tags = merge(var.common_tags, {
    Name = "${var.app_name}-${var.environment}-vpc"
    ResourceType = "vpc"
    NetworkTier = "core"
  })
}

# Internet Gateway
resource "aws_internet_gateway" "main" {
  vpc_id = aws_vpc.main.id

  tags = merge(var.common_tags, {
    Name = "${var.app_name}-${var.environment}-igw"
    ResourceType = "internet-gateway"
    NetworkTier = "public"
  })
}

# Get availability zones
data "aws_availability_zones" "available" {
  state = "available"
}

# Public Subnets
resource "aws_subnet" "public" {
  count = var.public_subnet_count

  vpc_id                  = aws_vpc.main.id
  cidr_block              = var.public_subnet_cidrs[count.index]
  availability_zone       = data.aws_availability_zones.available.names[count.index]
  map_public_ip_on_launch = true

  tags = merge(var.common_tags, {
    Name = "${var.app_name}-${var.environment}-public-subnet-${count.index + 1}"
    ResourceType = "subnet"
    SubnetType = "public"
    NetworkTier = "public"
    AvailabilityZone = data.aws_availability_zones.available.names[count.index]
  })
}

# Private Subnets
resource "aws_subnet" "private" {
  count = var.private_subnet_count

  vpc_id            = aws_vpc.main.id
  cidr_block        = var.private_subnet_cidrs[count.index]
  availability_zone = data.aws_availability_zones.available.names[count.index]

  tags = merge(var.common_tags, {
    Name = "${var.app_name}-${var.environment}-private-subnet-${count.index + 1}"
    ResourceType = "subnet"
    SubnetType = "private"
    NetworkTier = "private"
    AvailabilityZone = data.aws_availability_zones.available.names[count.index]
  })
}

# Elastic IPs for NAT Gateways
resource "aws_eip" "nat" {
  count = var.nat_gateway_count

  domain = "vpc"
  depends_on = [aws_internet_gateway.main]

  tags = merge(var.common_tags, {
    Name = "${var.app_name}-${var.environment}-nat-eip-${count.index + 1}"
    ResourceType = "elastic-ip"
    Purpose = "nat-gateway"
    AvailabilityZone = data.aws_availability_zones.available.names[count.index]
  })
}

# NAT Gateways
resource "aws_nat_gateway" "main" {
  count = var.nat_gateway_count

  allocation_id = aws_eip.nat[count.index].id
  subnet_id     = aws_subnet.public[count.index].id
  depends_on    = [aws_internet_gateway.main]

  tags = merge(var.common_tags, {
    Name = "${var.app_name}-${var.environment}-nat-gateway-${count.index + 1}"
    ResourceType = "nat-gateway"
    NetworkTier = "public"
    AvailabilityZone = data.aws_availability_zones.available.names[count.index]
  })
}

# Route Table for Public Subnets
resource "aws_route_table" "public" {
  vpc_id = aws_vpc.main.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.main.id
  }

  tags = merge(var.common_tags, {
    Name = "${var.app_name}-${var.environment}-public-rt"
    ResourceType = "route-table"
    RouteTableType = "public"
    NetworkTier = "public"
  })
}

# Route Tables for Private Subnets
resource "aws_route_table" "private" {
  count = var.private_subnet_count

  vpc_id = aws_vpc.main.id

  route {
    cidr_block     = "0.0.0.0/0"
    nat_gateway_id = aws_nat_gateway.main[count.index % var.nat_gateway_count].id
  }

  tags = merge(var.common_tags, {
    Name = "${var.app_name}-${var.environment}-private-rt-${count.index + 1}"
    ResourceType = "route-table"
    RouteTableType = "private"
    NetworkTier = "private"
    AvailabilityZone = data.aws_availability_zones.available.names[count.index]
  })
}

# Route Table Associations for Public Subnets
resource "aws_route_table_association" "public" {
  count = var.public_subnet_count

  subnet_id      = aws_subnet.public[count.index].id
  route_table_id = aws_route_table.public.id
}

# Route Table Associations for Private Subnets
resource "aws_route_table_association" "private" {
  count = var.private_subnet_count

  subnet_id      = aws_subnet.private[count.index].id
  route_table_id = aws_route_table.private[count.index].id
}