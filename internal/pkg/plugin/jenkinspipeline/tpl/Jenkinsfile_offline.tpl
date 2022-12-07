podTemplate(containers: [
    containerTemplate(name: 'maven', image: 'maven:3.8.6-openjdk-18', command: 'sleep', args: '99d'),
    containerTemplate(name: 'buildkit', image: 'moby/buildkit:master', ttyEnabled: true, privileged: true),
  ], volumes: [
    secretVolume(secretName: '[[ .ImageRepoDockerSecret ]]', mountPath: '/root/.docker')
  ]) {
    node(POD_LABEL) {
        stage("Get Project") {
            checkout scm
        }
        stage('Run Maven test') {
            gitlabCommitStatus("test") {
                container('maven') {
                    stage('run mvn test') {
                        sh 'mvn -B test'
                    }
                }
            }
        }
        stage("Build Docker image") {
            gitlabCommitStatus("build image") {
                container('buildkit') {
                    stage('build a Maven project') {
                        String opts = ""
                        String imageRepo = "[[ .imageRepo.user ]]/[[ .AppName ]]"
                        String imageURL = "[[ .imageRepo.url ]]"
                        if (imageURL) {
                            imageRepo = "${imageURL}/${imageRepo}"
                        }
                        if (imageRepo.contains("http://")) {
                            opts = ",registry.insecure=true"
                            imageRepo = imageRepo.replace("http://", "")
                        }
                        String version
                        if (env.GIT_COMMIT) {
                            version = env.GIT_COMMIT.substring(0, 8)
                        } else {
                            sh "git config --global --add safe.directory '*'"
                            String gitCommitLang = sh (script: "git log -n 1 --pretty=format:'%H'", returnStdout: true)
                            version = gitCommitLang.substring(0, 8)
                        }
                        sh """
                          buildctl build --frontend dockerfile.v0 --local context=. --local dockerfile=. --output type=image,name=${imageRepo}:latest,push=true${opts}
                          buildctl build --frontend dockerfile.v0 --local context=. --local dockerfile=. --output type=image,name=${imageRepo}:${version},push=true${opts}
                        """
                    }
                }
            }
        }
    }
}
