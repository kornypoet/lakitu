variable "aws_profile" {
  default = "personal"
}

variable "aws_region" {
  default = "us-east-1"
}

provider "aws" {
  profile    = var.aws_profile
  region     = var.aws_region
}
