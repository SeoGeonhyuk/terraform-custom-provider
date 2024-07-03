terraform {
  required_providers {
    hashicups = {
      source = "hashicorp.com/edu/hashicups"
    }
  }
  required_version = ">= 1.1.0"
}

provider "hashicups" {
  username = "education"
  password = "test123"
  host     = "http://localhost:19090"
}

resource "hashicups_game" "edu" {
  name       = "good game"
  star_point = 3.4
  player_num = 3
}

output "edu_game" {
  value = hashicups_game.edu
}
