#!/usr/bin/env groovy

/* utils.groovy
This package is used for utils func
*/
import com.devstream.notification.Dingtalk
import com.devstream.ci.Git

/*
send notifucation
input params =>
statusMessage: jenkins status message
headMessage: jenkins head message
*/
def notify(String headMessage, String statusMessage) {
    def gitUtils = new Git()
    String changeString = gitUtils.getChangeString()
    switch (Config.notifySettings.notify_type) {
        case "dingding":
            dingtalk = new Dingtalk()
            dingtalk.send(changeString, headMessage, statusMessage)
            break
        default:
            throw new Exception("jenkins notify type ${Config.notifySettings.notifyType} doesn't support")
    }
}
