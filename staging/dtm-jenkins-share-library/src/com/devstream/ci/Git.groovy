package com.devstream.ci

def getChangeString() {
    def changeString = ""
    def MAX_MSG_LEN = 10
    def changeLogSets = currentBuild.changeSets
    for (int i = 0; i < changeLogSets.size(); i++) {
        def entries = changeLogSets[i].items
        for (int j = 0; j < entries.length; j++) {
            def entry = entries[j]
            truncatedMsg = entry.msg.take(MAX_MSG_LEN)
            commitTime = new Date(entry.timestamp).format("yyyy-MM-dd HH:mm:ss")
            changeString += " - ${truncatedMsg} [${entry.author} ${commitTime}]\n"
        }
    }
    if (!changeString) {
        changeString = " - No new changes"
    }
    return (changeString)
}

def getCommitIDHead() {
    String gitCommit
    if (env.GIT_COMMIT) {
        gitCommit = env.GIT_COMMIT.substring(0, 8)
    } else {
        sh "git config --global --add safe.directory '*'"
        String gitCommitLang = sh (script: "git log -n 1 --pretty=format:'%H'", returnStdout: true)
        gitCommit = gitCommitLang.substring(0, 8)
    }
    return gitCommit
}
