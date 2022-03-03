variable "vpc_cidr_block" {
  type    = string
  default = "10.0.0.0/16"
}

variable "vpc_name" {
  type    = string
  default = "main"
}

variable "public_subnets" {
  type = map(any)
  default = {
    "eu-central-1a" = "10.0.1.0/24"
    "eu-central-1b" = "10.0.2.0/24"
    "eu-central-1c" = "10.0.3.0/24"
  }
}

variable "private_subnets" {
  type = map(any)
  default = {
    "eu-central-1a" = "10.0.11.0/24"
    "eu-central-1b" = "10.0.12.0/24"
    "eu-central-1c" = "10.0.13.0/24"
  }
}

variable "team" {
  type        = string
  description = "used for tagging, to which team the resource belongs"
}
