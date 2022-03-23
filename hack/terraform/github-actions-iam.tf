resource "aws_iam_user" "githubactions" {
  name          = "devstream-github-actions"
  force_destroy = true
}

data "aws_iam_policy_document" "githubactions" {
  statement {
    actions = [
      "eks:DescribeCluster",
      "eks:AccessKubernetesApi",
    ]
    resources = ["${module.cluster.cluster_arn}"]
  }
}

resource "aws_iam_user_policy" "githubactions" {
  name   = "devstream-githubactions-eks-policy"
  user   = aws_iam_user.githubactions.name
  policy = data.aws_iam_policy_document.githubactions.json
}

resource "aws_iam_access_key" "githubactions" {
  user = aws_iam_user.githubactions.name
}

output "githubactions_iam_id" {
  sensitive = true
  value     = aws_iam_access_key.githubactions.id
}

output "githubactions_iam_secret" {
  sensitive = true
  value     = aws_iam_access_key.githubactions.secret
}
