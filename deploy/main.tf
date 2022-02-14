// Only one OIDC provider needs to exist per account
// This code would typically live somewhere more centralized
resource "aws_iam_openid_connect_provider" "github" {
  url = "https://token.actions.githubusercontent.com"

  client_id_list = [
    "sts.amazonaws.com"
  ]

  thumbprint_list = [
    "6938fd4d98bab03faadb97b34396831e3780aea1"
  ]
}

// The top-level role for Github Actions
resource "aws_iam_role" "github-actions" {
  name               = "GithubActions"
  assume_role_policy = data.aws_iam_policy_document.github-actions.json
}

// Only allow assume role if the repo is an exact match
// and the origin is the OIDC provided above
data "aws_iam_policy_document" "github-actions" {
  statement {
    effect  = "Allow"
    actions = ["sts:AssumeRoleWithWebIdentity"]
    principals {
      type        = "Federated"
      identifiers = [aws_iam_openid_connect_provider.github.arn]
    }

    condition {
      test     = "StringLike"
      variable = "token.actions.githubusercontent.com:sub"
      values   = ["repo:kornypoet/lakitu:*"]
    }
  }
}

// Policy to allow pushes to ECR
resource "aws_iam_policy" "allow-push-to-ecr" {
  name   = "AllowPushToEcr"
  policy = data.aws_iam_policy_document.allow-push-to-ecr.json
}

data "aws_iam_policy_document" "allow-push-to-ecr" {
  statement {
    effect    = "Allow"
    actions   = ["ecr:GetAuthorizationToken"]
    resources = ["*"]
  }

  statement {
    effect  = "Allow"
    actions = [
      "ecr:BatchCheckLayerAvailability",
      "ecr:CompleteLayerUpload",
      "ecr:InitiateLayerUpload",
      "ecr:PutImage",
      "ecr:UploadLayerPart",
    ]
    resources = [aws_ecr_repository.lakitu.arn]
  }
}

resource "aws_iam_role_policy_attachment" "github-actions-allow-push-to-ecr" {
  role       = aws_iam_role.github-actions.name
  policy_arn = aws_iam_policy.allow-push-to-ecr.arn
}

// Private repo for our containers
resource "aws_ecr_repository" "lakitu" {
  name                 = "lakitu"
  image_tag_mutability = "IMMUTABLE"
}
