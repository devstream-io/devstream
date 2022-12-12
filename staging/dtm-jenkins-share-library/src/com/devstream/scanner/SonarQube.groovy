package com.devstream.scanner

import com.devstream.ci.Git

def scanner(
    String name,
    String lang,
    String options='') {
    try {
        println('Info: Preparing SonarQube Scanner')
        gitUtil = new Git()
        version = gitUtil.getCommitIDHead()
        withSonarQubeEnv(){
            def private opts

            opts  = ' -Dsonar.projectKey='       + name
            opts += ' -Dsonar.projectName='      + name
            opts += ' -Dsonar.projectVersion='   + version
            opts += ' -Dsonar.language='         + lang
            opts += ' -Dsonar.projectBaseDir=.'
            opts += ' -Dsonar.sources=.'
            opts += ' -Dsonar.java.binaries=.'
            sonar_exec  = 'sonar-scanner' + opts + ' ' + options

            sh(sonar_exec)
        }
    }
    catch (e) {
        println('Error: Failed with SonarQube Scanner')
        throw e
    }
}

def qualityGateStatus(){
    try {
        timeout(time: Config.generalSettings.sonarqube_timeout_minutes, unit: 'MINUTES') {
            def qg_stats = waitForQualityGate()
            if (qg_stats.status != 'SUCCESS') {
                println('Error: Pipeline aborted due to quality gate failure: ' + qg.stats)
                error "Pipeline aborted due to quality gate failure: ${qg.status}"
            }
        }
    }
    catch (e) {
        throw e
    }
}
