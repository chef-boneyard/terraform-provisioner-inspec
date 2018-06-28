variable "instance_key" {
  description = "SSH Key for the instance"
  default     = "chartmann"
}

variable "aws_region" {
  description = "AWS region to launch servers."
  default     = "eu-central-1"
}

# Ubuntu Precise 12.04 LTS (x64)
variable "aws_amis" {
  default = {
    eu-central-1 = "ami-fa2fb595"
    eu-west-1 = "ami-674cbc1e"
    us-east-1 = "ami-1d4e7a66"
    us-west-1 = "ami-969ab1f6"
    us-west-2 = "ami-8803e0f0"
  }
}
