pipeline {
  agent {
    docker {
      image 'golang:1.11.4'
    }
  }
  options {
    buildDiscarder(logRotator(numToKeepStr: "2"))
    disableConcurrentBuilds()
  }
  stages {
    environment {
      GOPATH = "$WORKSPACE"
    }
    stage('Install Tools') {
      steps {
        sh 'go get github.com/goreleaser/goreleaser'
      }
    }
    stage('Build') {
      environment {
        GO111MODULE = "on"
      }
      steps {
        sh "'$GOPATH/bin/goreleaser' --snapshot --skip-publish --rm-dist"
      }
    }
    stage('Deploy') {
      environment {
        GO111MODULE = "on"
        GITHUB = credentials('github-halkeye')
      }
      when { tag "v*" }
      steps {
        sh "export GITHUB_TOKEN=$GITHUB_PSW"
        sh "'$GOPATH/bin/goreleaser' --rm-dist"
      }
    }
  }
}
