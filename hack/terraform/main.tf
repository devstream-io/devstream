module "network" {
  source = "./modules/networking"

  vpc_name       = "DevStream"
  vpc_cidr_block = "10.0.0.0/16"

  public_subnets = {
    "ap-southeast-1a" = "10.0.1.0/24"
    "ap-southeast-1b" = "10.0.2.0/24"
  }

  private_subnets = {
    "ap-southeast-1a" = "10.0.11.0/24"
    "ap-southeast-1b" = "10.0.12.0/24"
  }

  team = "DevStream"
}

module "cluster" {
  source = "./modules/eks"

  cluster_name               = "dtm-test"
  nodegroup_name             = "dtm-test-1"
  vpc_id                     = module.network.vpc_id
  worker_subnet_ids          = module.network.private_subnet_ids
  worker_instance_type       = "t2.medium"
  min_worker_node_number     = 1
  desired_worker_node_number = 1
  max_worker_node_number     = 1

  team = "DevStream"
}
