terraform {
  required_version = ">= 1.5.0"

  required_providers {
    openstack = {
      source  = "terraform-provider-openstack/openstack"
      version = ">= 2.1, < 4.0"
    }
    local = {
      source  = "hashicorp/local"
      version = "~> 2.5"
    }
  }
}

# The credentials are read from clouds.yaml (cloud name in var.cloud) or from the OS_* environment
# variables of an OpenStack RC file when var.cloud is empty
provider "openstack" {
  cloud = var.cloud != "" ? var.cloud : null
}
