resource "aws_security_group" "cluster" {
  name        = "${var.cluster_name}_eks_cluster_sg"
  description = "EKS cluster security group."
  vpc_id      = var.vpc_id

  tags = {
    Name = "${var.cluster_name}_eks_cluster_sg"
    Team = var.team
  }
}

resource "aws_security_group_rule" "cluster_egress_internet" {
  description       = "Allow cluster egress access to the Internet."
  protocol          = "-1"
  security_group_id = aws_security_group.cluster.id
  cidr_blocks       = ["0.0.0.0/0"]
  from_port         = 0
  to_port           = 0
  type              = "egress"
}

resource "aws_security_group_rule" "cluster_https_worker_ingress" {
  description              = "Allow pods to communicate with the EKS cluster API."
  protocol                 = "tcp"
  security_group_id        = aws_security_group.cluster.id
  source_security_group_id = aws_security_group.worker.id
  from_port                = 443
  to_port                  = 443
  type                     = "ingress"
}

resource "aws_eks_cluster" "cluster" {
  depends_on = [
    aws_iam_role_policy_attachment.AmazonEKSClusterPolicy,
    aws_iam_role_policy_attachment.AmazonEKSServicePolicy,
  ]
  version  = var.k8s_version
  name     = var.cluster_name
  role_arn = aws_iam_role.eks_role.arn
  vpc_config {
    subnet_ids         = var.worker_subnet_ids
    security_group_ids = [aws_security_group.cluster.id]
  }
  enabled_cluster_log_types = ["api", "audit", "authenticator", "controllerManager", "scheduler"]

  tags = {
    Team = var.team
  }
}
