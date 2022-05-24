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
    resources = [module.cluster.cluster_arn]
  }
}

resource "aws_iam_user_policy" "githubactions" {
  name   = "devstream-githubactions-eks-policy"
  user   = aws_iam_user.githubactions.name
  policy = data.aws_iam_policy_document.githubactions.json
}

resource "aws_iam_user_policy" "GitHubActions-Download-Bucket-RW-Policy" {
  name  = "GitHubActions-Download-Bucket-RW-Policy"
  user = aws_iam_user.githubactions.name

  policy = jsonencode({
    "Version": "2012-10-17",
    "Statement": [
      {
        "Sid": "PermissionForDownloadUploadForTheBucket",
        "Effect": "Allow",
        "Action": [
          "s3:GetObject",
          "s3:GetObjectAcl",
          "s3:PutObject",
          "s3:PutObjectAcl",
          "s3:DeleteObject",
          "s3:RestoreObject",
          "s3:ListBucket",
          "s3:GetBucketPolicy",
          "s3:ReplicateObject",
          "s3:GetBucketWebsite",
          "s3:PutBucketWebsite",
          "s3:GetBucketCORS",
        ],
        "Resource": [
          "arn:aws:s3:::download.devstream.io",
          "arn:aws:s3:::download.devstream.io/*"
        ]
      }
    ]
  })
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
