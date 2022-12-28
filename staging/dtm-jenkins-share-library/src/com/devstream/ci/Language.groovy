package com.devstream.ci

def selector(String language){
    // config default options for different language
    switch(language.toLowerCase()){
        case "java":
            return javaDefault()
        default:
            if (!Config.generalSettings.ci_test_command) {
                throw new Exception("Language %s language should set ci_test_command and ci_test_options in generalSettings")
            }
    }
}


def javaDefault() {
    return [
        ci_test_command: 'mvn',
        ci_test_options: '-B test',
        ci_test_container_repo: 'maven:3.8.1-jdk-8',
        container_requests_cpu: "512m",
        container_requests_memory: "2Gi",
        container_limit_cpu: "512m",
        container_limit_memory: "2Gi",
    ]
}
