resource "aws_iam_group" "devstream" {
  name = "DevStream"
  path = "/"
}

resource "aws_iam_group_policy" "devstream-enforce-mfa" {
  name  = "devstream_enforce_mfa"
  group = aws_iam_group.devstream.name

  policy = jsonencode({
    "Version" : "2012-10-17",
    "Statement" : [
      {
        "Sid" : "AllowViewAccountInfo",
        "Effect" : "Allow",
        "Action" : [
          "iam:GetAccountPasswordPolicy",
          "iam:GetAccountSummary",
          "iam:ListVirtualMFADevices"
        ],
        "Resource" : "*"
      },
      {
        "Sid" : "AllowManageOwnPasswords",
        "Effect" : "Allow",
        "Action" : [
          "iam:ChangePassword",
          "iam:GetUser"
        ],
        "Resource" : "arn:aws:iam::*:user/$${aws:username}"
      },
      {
        "Sid" : "AllowManageOwnAccessKeys",
        "Effect" : "Allow",
        "Action" : [
          "iam:CreateAccessKey",
          "iam:DeleteAccessKey",
          "iam:ListAccessKeys",
          "iam:UpdateAccessKey"
        ],
        "Resource" : "arn:aws:iam::*:user/$${aws:username}"
      },
      {
        "Sid" : "AllowManageOwnSigningCertificates",
        "Effect" : "Allow",
        "Action" : [
          "iam:DeleteSigningCertificate",
          "iam:ListSigningCertificates",
          "iam:UpdateSigningCertificate",
          "iam:UploadSigningCertificate"
        ],
        "Resource" : "arn:aws:iam::*:user/$${aws:username}"
      },
      {
        "Sid" : "AllowManageOwnSSHPublicKeys",
        "Effect" : "Allow",
        "Action" : [
          "iam:DeleteSSHPublicKey",
          "iam:GetSSHPublicKey",
          "iam:ListSSHPublicKeys",
          "iam:UpdateSSHPublicKey",
          "iam:UploadSSHPublicKey"
        ],
        "Resource" : "arn:aws:iam::*:user/$${aws:username}"
      },
      {
        "Sid" : "AllowManageOwnGitCredentials",
        "Effect" : "Allow",
        "Action" : [
          "iam:CreateServiceSpecificCredential",
          "iam:DeleteServiceSpecificCredential",
          "iam:ListServiceSpecificCredentials",
          "iam:ResetServiceSpecificCredential",
          "iam:UpdateServiceSpecificCredential"
        ],
        "Resource" : "arn:aws:iam::*:user/$${aws:username}"
      },
      {
        "Sid" : "AllowManageOwnVirtualMFADevice",
        "Effect" : "Allow",
        "Action" : [
          "iam:CreateVirtualMFADevice",
          "iam:DeleteVirtualMFADevice"
        ],
        "Resource" : "arn:aws:iam::*:mfa/$${aws:username}"
      },
      {
        "Sid" : "AllowManageOwnUserMFA",
        "Effect" : "Allow",
        "Action" : [
          "iam:DeactivateMFADevice",
          "iam:EnableMFADevice",
          "iam:ListMFADevices",
          "iam:ResyncMFADevice"
        ],
        "Resource" : "arn:aws:iam::*:user/$${aws:username}"
      }
    ]
  })
}

resource "aws_iam_group_policy_attachment" "devstream-iam" {
  group      = aws_iam_group.devstream.name
  policy_arn = "arn:aws:iam::aws:policy/IAMReadOnlyAccess"
}

resource "aws_iam_group_policy" "devstream-eks" {
  name  = "devstream-eks"
  group = aws_iam_group.devstream.name

  policy = jsonencode({
    "Version" : "2012-10-17",
    "Statement" : [
      {
        "Effect" : "Allow",
        "Action" : [
          "eks:DescribeNodegroup",
          "eks:ListNodegroups",
          "eks:DescribeCluster",
          "eks:ListClusters",
          "eks:AccessKubernetesApi",
          "ssm:GetParameter",
          "eks:ListUpdates",
          "eks:ListFargateProfiles"
        ],
        "Resource" : "*"
      }
    ]
  })
}

resource "aws_iam_group_policy" "DevStream-Download-Bucket-RW-Policy" {
  name  = "DevStream-Download-Bucket-RW-Policy"
  group = aws_iam_group.devstream.name

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

locals {
  users = ["fangbao", "hutao"]
}

resource "aws_iam_user" "devstream" {
  for_each = toset(local.users)

  name          = each.key
  force_destroy = true
}

resource "aws_iam_user_login_profile" "devstream" {
  for_each = toset(local.users)

  user = aws_iam_user.devstream[each.key].name
  # password_reset_required = true
}

output "users" {
  value = {
    for k, v in aws_iam_user_login_profile.devstream : k => v.password
  }
}

resource "aws_iam_user_group_membership" "devstream" {
  for_each = toset(local.users)

  user   = aws_iam_user.devstream[each.key].name
  groups = [aws_iam_group.devstream.name]
}
