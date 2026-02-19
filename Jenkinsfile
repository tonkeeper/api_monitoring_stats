#!groovy

pipeline {
  agent any
  options {
    timestamps()
    disableConcurrentBuilds()
  }
  environment {
    APP_DOCKER_IMAGE_NAME = "${env.JOB_BASE_NAME}"
    APP_GIT_REPO = "git@github.com:tonkeeper/${env.JOB_BASE_NAME}"
    APP_GIT_DEFAULT_BRANCH = 'main'
    JENKINS_GIT_CRED_NAME = 'tonkeeper-build-bot_ssh'
    TG_TOKEN = credentials('tg_bot_zaya_monitoring')
    TG_CHAT_ID = credentials('tg_chan_tonapi_build')
  }
  triggers {
    GenericTrigger(
     genericVariables: [[key: 'REF', value: '$.ref']],
     causeString: 'Triggered on tag $REF',
     regexpFilterExpression: '^v[0-9].*',
     regexpFilterText: '$REF',
     printContributedVariables: true,
     printPostContent: false,
     token: "${env.JOB_BASE_NAME}-mildly-sensitive-token"
    )
  }
  stages {
    stage('checkout') {
      steps {
        script {

          env.APP_GIT_BRANCH = env.APP_GIT_DEFAULT_BRANCH
          if (env.REF) {
            env.APP_GIT_BRANCH = env.REF
            env.APP_DOCKER_REPO_TAG = env.REF
          }

          def scmVars = checkout([
                          $class: 'GitSCM',
                          userRemoteConfigs: [[
                            url: "${APP_GIT_REPO}",
                            credentialsId: "${JENKINS_GIT_CRED_NAME}"
                          ]],
                          branches: [[name: "${APP_GIT_BRANCH}"]],
                          extensions: [[
                            $class: 'SubmoduleOption',
                            disableSubmodules: false,
                            parentCredentials: true,
                            recursiveSubmodules: true,
                            reference: '',
                            shallow: true,
                            trackingSubmodules: true
                          ]]
                        ])
          env.APP_GIT_COMMIT = scmVars.GIT_COMMIT[0..6]
          env.APP_DOCKER_REPO_TAG = env.APP_DOCKER_REPO_TAG ?: "${APP_GIT_COMMIT}"
          env.APP_BUILD_DATE = sh(script: "date --rfc-3339=seconds --utc", returnStdout: true).trim()
        }
        sh """
           docker build . \
           --build-arg DOCKER_IMAGE_VERSION="${APP_DOCKER_REPO_TAG}" \
           --label org.opencontainers.image.source="${APP_GIT_REPO}" \
           --label org.opencontainers.image.created="${APP_BUILD_DATE}" \
           --label org.opencontainers.image.revision="${APP_GIT_COMMIT}" \
           --label docker.repo.tag="${APP_DOCKER_REPO_TAG}" \
           --tag "${DOCKER_REPOSITORY}/${APP_DOCKER_IMAGE_NAME}:${APP_DOCKER_REPO_TAG}" \
           --tag "${DOCKER_REPOSITORY}/${APP_DOCKER_IMAGE_NAME}:${APP_GIT_COMMIT}" \
           --tag "${DOCKER_REPOSITORY}/${APP_DOCKER_IMAGE_NAME}:latest"
           docker push "${DOCKER_REPOSITORY}/${APP_DOCKER_IMAGE_NAME}:${APP_DOCKER_REPO_TAG}"
           docker push "${DOCKER_REPOSITORY}/${APP_DOCKER_IMAGE_NAME}:${APP_GIT_COMMIT}"
           docker push "${DOCKER_REPOSITORY}/${APP_DOCKER_IMAGE_NAME}:latest"
           """
      }
    }
  }
  post {
    success {
      sh '''
      curl -s https://api.telegram.org/bot${TG_TOKEN}/sendMessage \
      --request POST --header 'Content-Type: multipart/form-data' \
      --form chat_id=${TG_CHAT_ID} \
      --form text="${APP_DOCKER_IMAGE_NAME}:${APP_DOCKER_REPO_TAG} OK"
      '''
    }
    failure {
      sh '''
      curl -s https://api.telegram.org/bot${TG_TOKEN}/sendMessage \
      --request POST --header 'Content-Type: multipart/form-data' \
      --form chat_id=${TG_CHAT_ID} \
      --form text="${APP_DOCKER_IMAGE_NAME}:${APP_DOCKER_REPO_TAG} FAIL"
      '''
    }
  }
}
