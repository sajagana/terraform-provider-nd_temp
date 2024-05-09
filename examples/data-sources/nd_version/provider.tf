terraform {
  required_providers {
    nd = {
      source = "hashicorp.com/edu/nd"
    }
  }
}

provider "nd" {
  username = "admin"
  password = "ins3965!"
  url      = "https://173.36.219.32/" # Sabari System
  insecure = true
  platform = "nd"
}
