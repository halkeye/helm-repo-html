pipeline {
  agent {
    docker {
      image 'golang:1.11.4'
    }
  }
  options {
    buildDiscarder(logRotator(numToKeepStr: "2"))
    disableConcurrentBuilds()
    timeout(time: 10, unit: 'MINUTES')
    ansiColor('xterm')
    checkoutToSubdirectory('go/src/github.com/halkeye/helm-repo-html')
  }
  environment {
    GOPATH = "$WORKSPACE/go"
    GO111MODULE = "on"
  }
  stages {
    stage('Install Tools') {
      steps {
        sh 'curl -Ls https://github.com/goreleaser/goreleaser/releases/download/v0.101.0/goreleaser_Linu│·x_x86_64.tar.gz | tar xvzf - goreleaser'
      }
    }
    stage('Build') {
      steps {
        dir('go/src/github.com/halkeye/helm-repo-html') {
          sh "'$WORKSPACE/goreleaser' --snapshot --skip-publish --rm-dist"^
        }
      }
    }
    stage('Deploy') {
      environment {
        GITHUB = credentials('github-halkeye')
      }
      when { tag "v*" }
      steps {
        dir('go/src/github.com/halkeye/helm-repo-html') {
          sh "export GITHUB_TOKEN=$GITHUB_PSW"
          sh "'$WORKSPACE/goreleaser' --rm-dist"
        }
      }
    }
  }
}
