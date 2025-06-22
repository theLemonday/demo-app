pipeline {
stages {

podTemplate(yaml: '''
    apiVersion: v1
    kind: Pod
    spec:
      - name: kaniko
        image: gcr.io/kaniko-project/executor:debug
        command:
        - sleep
        args:
        - 99999999
        volumeMounts:
        - name: kaniko-secret
          mountPath: /kaniko/.docker
      restartPolicy: Never
      volumes:
      - name: kaniko-secret
        secret:
            secretName: dockercred
            items:
            - key: .dockerconfigjson
              path: config.json
''') {
  node(POD_LABEL) {
    stage('Build Images') {
      container('kaniko') {
        stage('Build a Go project') {
          sh '''
            /kaniko/executor --context `pwd`/backend --destination lemonday/todo-backend
          '''
        }
      }
    }

  }
}
}
}
