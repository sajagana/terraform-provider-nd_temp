terraform {
  required_providers {
    nd = {
      source = "hashicorp.com/edu/nd"
    }
  }
}

provider "nd" {
  username = "admin"
  password = "password"
  url      = "https://my-cisco-nd.com"
  insecure = true
  platform = "nd"
}

data "nd_version" "example" {
}
