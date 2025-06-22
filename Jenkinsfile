pipeline {
  agent {
    kubernetes {
      yaml """
apiVersion: v1
kind: Pod
spec:
  containers:
  - name: kaniko
    image: gcr.io/kaniko-project/executor:debug
    imagePullPolicy: Always
    command: [sleep]
    args: [9999999]
    volumeMounts:
    - name: workspace-volume
      mountPath: /home/jenkins/agent
    - name: docker-config
      mountPath: /kaniko/.docker/
  volumes:
  # Volume để chia sẻ workspace giữa tất cả các container
  - name: workspace-volume
    emptyDir: {}
  # Volume từ Secret để Kaniko xác thực với Docker Hub
  - name: docker-config
    secret:
      secretName: dockercred
      items:
        - key: .dockerconfigjson
          path: config.json
"""
    }
  }

  environment {
    DOCKER_HUB_USER = "lemonday"
    GIT_APP_REPO_URL = "https://github.com/theLemonday/demo-app"
    GIT_CONFIG_REPO_URL = "https://github.com/theLemonday/demo-app-values"
  }

  stages {
    stage('Checkout code') {
      steps {
        script {
          echo "Checking out source code"
          git url: env.GIT_APP_REPO_URL,
            branch: 'main',
          echo "Checkout completed"
        }
      }
    }

    stage('Build & Push Docker Image (with Kaniko)') {
      steps {
        // THAY ĐỔI: Chạy `script` ở ngoài `container` để lấy git commit trước
        script {
          // Bước 1: Lấy git commit hash trong container mặc định 'jnlp' (nơi có git)
          def gitCommit = sh(script: 'git rev-parse HEAD', returnStdout: true).trim().substring(0, 8)
          def images = ["frontend", "backend"]

          for (image in images) {
            def dockerImageTag = "todo-${image}:${gitCommit}"
            echo "Building ${dockerImageTag}"

            container('kaniko') {
              sh """
              /kaniko/executor \\
              --context `pwd`/${image} \\
              --dockerfile `pwd`/${image}/Dockerfile \\
              --destination ${dockerImageTag}
              """
            }

            echo "Finished building ${dockerImageTag}"
          }
        }
      }
    }
  }
}
