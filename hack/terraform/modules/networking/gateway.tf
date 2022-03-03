# internet gateway
resource "aws_internet_gateway" "main" {
  vpc_id = aws_vpc.main.id

}

# elastic IP for nat gateway
resource "aws_eip" "nat_gateway" {
  for_each = var.public_subnets

  vpc = true

  tags = {
    Team = var.team
  }
}

# nat gateway
resource "aws_nat_gateway" "main" {
  for_each = var.public_subnets

  allocation_id = aws_eip.nat_gateway[each.key].id
  subnet_id     = aws_subnet.public[each.key].id

  tags = {
    Team = var.team
  }
}
