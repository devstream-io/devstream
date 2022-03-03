# AWS VPC Module

This module creates:

- VPC
- public subnets
- private subnets
- NAT gateways in each public subnet (including the EIP of course)
- routing tables and associations (each private subnet routes to the NAT gateway in the same AZ)
- Internet gateway
