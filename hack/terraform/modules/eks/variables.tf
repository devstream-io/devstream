# mandatory variables

# cluster
variable "cluster_name" {
  type = string
}

variable "nodegroup_name" {
  type = string
}

# networking
variable "vpc_id" {
  description = "The existing VPC"
}

variable "worker_subnet_ids" {
  type        = list(any)
  description = "the subnets where to deploy EKS"
}

# optional variables
variable "k8s_version" {
  type    = string
  default = "1.21"
}

variable "worker_instance_type" {
  description = "worker type, like t2.medium"
  type        = string
  default     = "t2.medium"
}

variable "min_worker_node_number" {
  type    = number
  default = 2
}

variable "desired_worker_node_number" {
  type    = number
  default = 2
}

variable "max_worker_node_number" {
  type    = number
  default = 2
}

variable "team" {
  type        = string
  description = "used for tagging, to which team the resource belongs"
}
