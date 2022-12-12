def call(Map config=[:]) {
    setting.configGeneral(config)
    templates {
        entry()
    }
}

def entry() {
    node(POD_LABEL) {
        try {
            controller.cloneCode()
            parallel(
                'Test': {controller.testCode()},
                'Sonar Scan': {controller.sonarScan()},
            )
            controller.pushCodeImage()
            post.success()
        } catch (org.jenkinsci.plugins.workflow.steps.FlowInterruptedException err) {
            post.aborted()
            throw err
        } catch (Exception err) {
            post.failure()
            throw err
        }
    }
}



def templates(Closure body) {
    def s = Config.generalSettings
    if (s.repo_type == "gitlab") {
        gitlabCommitStatus {
            testTemplate(body)
        }
    } else {
        testTemplate(body)
    }
}

def testTemplate(Closure body) {
    def s = Config.generalSettings
    if (!s.enable_test) {
        if (s.enable_sonarqube) {
            pod.scannerTemplate {
                pod.buildTemplate {
                    body.call()
                }
            }
        } else {
            pod.buildTemplate {
                body.call()
            }
        }
    } else {
        if (s.enable_sonarqube) {
            pod.testTemplate {
                pod.scannerTemplate {
                    pod.buildTemplate {
                        body.call()
                    }
                }
            }
        } else {
            pod.testTemplate {
                pod.buildTemplate {
                    body.call()
                }
            }
        }
    }
}
