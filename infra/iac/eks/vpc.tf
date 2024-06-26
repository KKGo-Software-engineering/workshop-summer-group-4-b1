resource "aws_vpc" "workshop" {
  cidr_block = "10.0.0.0/16"

  tags = {
    Cluster = var.cluster_name
    Name    = "${var.cluster_name}-vpc"
  }

  enable_dns_support   = true
  enable_dns_hostnames = true
}
