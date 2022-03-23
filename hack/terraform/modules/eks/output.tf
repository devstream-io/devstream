output "cluster_name" {
  value = aws_eks_cluster.cluster.id
}

output "cluster_arn" {
  value = aws_eks_cluster.cluster.arn
}

output "endpoint" {
  value = aws_eks_cluster.cluster.endpoint
}

output "kubeconfig_certificate_authority_data" {
  value = aws_eks_cluster.cluster.certificate_authority[0].data
}

output "security_group_id" {
  value = aws_eks_cluster.cluster.vpc_config[0].cluster_security_group_id
}

output "worker_node_role_name" {
  value = aws_iam_role.worker_role.name
}
