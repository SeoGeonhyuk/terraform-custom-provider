terraform {
  required_providers {
    hashicups = {
      // .terraformmrc를 통해 dev_override 블럭을 추가하고 그 안에 GOBIN을 넣으면
      // terraform이 환경 변수 GOBIN 위치에 빌드되어 있는 프로바이더 바이너리 파일을 가지고 와서 한번에 실행한다.
      // 이미 빌드되어 있는 파일을 인식함으로써 테라폼 명령어를 통해서 테라폼은 프로바이더를 인식할 수 있다.
      source = "hashicorp.com/edu/hashicups"
    }
  }
}

provider "hashicups" {
  host     = "http://localhost:19090"
  username = "education"
  password = "test123"
}

data "hashicups_coffees" "edu" {}

output "edu_coffees" {
    value = data.hashicups_coffees.edu
}
