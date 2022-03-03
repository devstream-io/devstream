resource "aws_eks_node_group" "default" {
  depends_on = [
    aws_iam_role_policy_attachment.eks_worker_AmazonEKSWorkerNodePolicy,
    aws_iam_role_policy_attachment.eks_worker_AmazonEKS_CNI_Policy,
    aws_iam_role_policy_attachment.eks_worker_AmazonEC2ContainerRegistryReadOnly,
    aws_iam_role_policy_attachment.eks_worker_autoscaling
  ]

  cluster_name    = aws_eks_cluster.cluster.name
  node_group_name = "${aws_eks_cluster.cluster.name}-managed-${var.nodegroup_name}"
  node_role_arn   = aws_iam_role.worker_role.arn
  subnet_ids      = var.worker_subnet_ids
  instance_types  = [var.worker_instance_type]

  scaling_config {
    min_size     = var.min_worker_node_number
    desired_size = var.desired_worker_node_number
    max_size     = var.max_worker_node_number
  }

  tags = {
    Team = var.team
  }
}
