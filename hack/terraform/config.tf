terraform {
  required_version = ">= 1.1.7"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "4.3.0"
    }
  }

  backend "s3" {
    bucket = "devstream-terraform-state"
    key    = "test.tfstate"
    region = "ap-southeast-1"
  }
}
provider "aws" {
  region = "ap-southeast-1"
}

variable "s3_bucket_names" {
  type    = set(string)
  default = ["download.devstream.io"]
}

resource "aws_s3_bucket" "devstream_buckets" {
  for_each = var.s3_bucket_names
  bucket = each.key
}
