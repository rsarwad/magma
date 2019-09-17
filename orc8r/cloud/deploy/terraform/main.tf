terraform {
  required_version = ">= 0.12.0"
}

provider "aws" {
  version = ">= 2.6.0"
  region  = var.region
}

module "vpc" {
  source = "terraform-aws-modules/vpc/aws"
  name   = var.vpc_name

  cidr = "10.10.0.0/16"

  azs = data.aws_availability_zones.available.names

  public_subnets   = ["10.10.1.0/24", "10.10.2.0/24", "10.10.3.0/24"]
  private_subnets  = []
  database_subnets = ["10.10.11.0/24", "10.10.12.0/24", "10.10.13.0/24"]

  create_database_subnet_group = true

  tags = {
    "kubernetes.io/cluster/${var.cluster_name}" = "shared"
  }

  public_subnet_tags = {
    "kubernetes.io/cluster/${var.cluster_name}" = "shared"
    "kubernetes.io/role/elb"                    = 1
  }
}

data "template_file" "metrics_userdata" {
  template = file("${path.module}/scripts/prepare_metrics_instance.sh.tpl")
}

data "aws_iam_policy_document" "worker_node_policy_doc" {
  statement {
    effect = "Allow"

    actions = [
      "ec2:DescribeVolumes",
      "ec2:AttachVolume",
      "ec2:DetachVolume",
    ]

    resources = [
      "arn:aws:ec2:*:*:volume/*",
      "arn:aws:ec2:*:*:instance/*",
    ]
  }
}

resource "aws_iam_policy" "worker_node_policy" {
  name   = "magma_eks_worker_node_policy"
  policy = data.aws_iam_policy_document.worker_node_policy_doc.json
}

module "eks" {
  source       = "terraform-aws-modules/eks/aws"
  cluster_name = var.cluster_name
  vpc_id       = module.vpc.vpc_id
  subnets      = module.vpc.public_subnets

  worker_additional_security_group_ids = [aws_security_group.default.id]
  workers_additional_policies          = [aws_iam_policy.worker_node_policy.arn]

  # asg max capacity is 3
  # 1 worker group for orc8r (3 boxes total)
  # 1 worker group for metrics (1 box)
  worker_groups = [
    {
      name                 = "wg-1"
      instance_type        = "t3.small"
      asg_desired_capacity = 2
      key_name             = var.key_name
      kubelet_extra_args   = "--node-labels=worker-type=controller"

      tags = [
        {
          key                 = "orc8r-node-type"
          value               = "orc8r-worker-node"
          propagate_at_launch = true
        },
      ]
    },
    {
      name                 = "wg-metrics"
      instance_type        = "t3.medium"
      asg_desired_capacity = 1
      key_name             = var.key_name
      kubelet_extra_args   = "--node-labels=worker-type=metrics"

      # we put the metrics nodes into 1 specific subnet because EBS volumes
      # can only be mounted into the same AZ
      subnets = [module.vpc.public_subnets[0]]

      additional_userdata = data.template_file.metrics_userdata.rendered

      tags = [
        {
          key                 = "orc8r-node-type"
          value               = "orc8r-prometheus-node"
          propagate_at_launch = true
        },
      ]
    },
  ]

  map_users = var.map_users

  write_kubeconfig      = true
  write_aws_auth_config = true
}

resource "aws_security_group" "default" {
  name        = "orc8r-default-sg"
  description = "Default orc8r SG"
  vpc_id      = module.vpc.vpc_id

  ingress {
    from_port = 0
    to_port   = 0
    protocol  = "-1"
    self      = "true"
  }

  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

# EBS volume for prometheus metrics.
resource "aws_ebs_volume" "prometheus-ebs-eks" {
  availability_zone = data.aws_availability_zones.available.names[0]
  size              = 64

  tags = {
    Name = "orc8r-prometheus-data"
  }
}

# EBS volume for prometheus configs.
resource "aws_ebs_volume" "prometheus-configs-ebs-eks" {
  availability_zone = data.aws_availability_zones.available.names[0]
  size              = 1

  tags = {
    Name = "orc8r-prometheus-configs"
  }
}
