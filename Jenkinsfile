#!/usr/bin/env groovy


def projectName = "smackapi"
def imageName = "alex202/brigade-smackapi"

def gitCommit = null
def gitBranch = null
def imageTag = null
def buildDate = null

podTemplate(label: 'mypod', containers: [
    containerTemplate(name: 'golang', image: 'golang:1.9', ttyEnabled: true, command: 'cat'),
    containerTemplate(name: 'docker', image: 'docker', ttyEnabled: true, command: 'cat'),
    containerTemplate(name: 'kubectl', image: 'lachlanevenson/k8s-kubectl:v1.8.0', command: 'cat', ttyEnabled: true),
    containerTemplate(name: 'helm', image: 'lachlanevenson/k8s-helm:latest', command: 'cat', ttyEnabled: true)
  ],
  envVars: [

  ],
  volumes: [
    hostPathVolume(mountPath: '/var/run/docker.sock', hostPath: '/var/run/docker.sock'),
  ]) {

    node('mypod') {

        checkout scm

        // print environment variables
        //echo sh(script: 'env|sort', returnStdout: true)

        sh "git rev-parse --short HEAD > .git/commit-id"
        gitCommit = readFile('.git/commit-id').trim()

        // git branch name is taken from an env var for multi-branch pipeline project, or from git for other projects
        if (env['BRANCH_NAME']) {
            gitBranch = BRANCH_NAME
        } else {

            //sh "git rev-parse --symbolic-full-name --abbrev-ref HEAD > .git/branch-name"
            //gitBranch = readFile('.git/branch-name').trim()
            gitBranch = sh returnStdout: true, script: 'git rev-parse --abbrev-ref HEAD'
            gitBranch = gitBranch.trim()

        }

        imageTag = "${gitBranch}-${gitCommit}"

        sh "date +'%Y-%m-%d %H-%M-%S' > .git/build-date"
        buildDate = readFile('.git/build-date').trim()


        def buildInfo = """# Build info
BUILD_NUMBER=${env.BUILD_NUMBER}
BUILD_DATE=${buildDate}
BUILD_GIT_COMMIT=${gitCommit}
BUILD_GIT_BRANCH=${gitBranch}
DOCKER_IMAGE_TAG=${imageTag}
"""

        echo buildInfo

        stage('Build go binaries') {
            container('golang') {

            sh "git rev-parse --symbolic-full-name --abbrev-ref HEAD"

                def pwd = pwd()

                sh """
                    go env
                    ls -la
                    mkdir -p /go/src/github.com/alex-egorov
                    ln -s $pwd /go/src/github.com/alex-egorov/
                    cd /go/src/github.com/alex-egorov/${projectName}
                    go get -v
                    GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -v -o target/${projectName}

                    ls -la
                """

                archiveArtifacts artifacts: 'target/*', excludes: 'target/*.tmp'
            }
        }

        stage('Build and push docker image') {
            container('docker') {

                withCredentials([[$class: 'UsernamePasswordMultiBinding',
                        credentialsId: 'dockerhub',
                        usernameVariable: 'DOCKER_HUB_USER',
                        passwordVariable: 'DOCKER_HUB_PASSWORD']]) {

                    sh """
                      docker build --force-rm \
                            --build-arg BUILD_DATE="${buildDate}" \
                            --build-arg IMAGE_TAG_REF=${imageTag} \
                            --build-arg VCS_REF=${gitCommit} \
                            -t ${imageName}:${imageTag} .
                      """
                    sh "docker login -u ${DOCKER_HUB_USER} -p ${DOCKER_HUB_PASSWORD} "
                    sh "docker push ${imageName}:${imageTag} "
                }
            }
        }

        stage('do some kubectl work') {
            container('kubectl') {

                sh "kubectl get nodes --all-namespaces"
            }
        }
        stage('do some helm work') {
            container('helm') {

                dir("charts") {

                    sh "helm ls"

                    sh "helm lint smackapi-release"
                    sh "helm upgrade -i smackapi-v${env.BUILD_NUMBER} --set image.tag=${imageTag} smackapi-release"
                }
            }
        }
    }
}