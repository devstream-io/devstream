output "vpc_id" {
  value = aws_vpc.main.id
}

output "private_subnet_route_table_ids" {
  value = [
    for r in aws_route_table.private :
    r.id
  ]
}

output "private_subnet_ids" {
  value = [
    for r in aws_subnet.private :
    r.id
  ]
}
