def buildUserName(){
    /*
       Requeire 'build user vars' plugin, See https://plugins.jenkins.io/build-user-vars-plugin for more information
    */
    wrap([$class: 'BuildUser']) {
        return BUILD_USER
    }
}

def checkGitRepo() {
    def repoType = Config.generalSettings.get("repo_type")
    return repoType && repoType == "gitlab"
}

def checkPushImage() {
    String imageName = Config.imageRepoSettings.get("image_name")
    return (imageName && imageName != "")
}
