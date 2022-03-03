resource "aws_security_group" "control_plane" {
  name        = "eks_cluster_${var.cluster_name}_control_plane_sg"
  description = "EKS cluster ${var.cluster_name} control plane security group."

  vpc_id = var.vpc_id

  tags = {
    Name = "eks_cluster_${var.cluster_name}_control_plane_sg"
    Team = var.team

  }
}
