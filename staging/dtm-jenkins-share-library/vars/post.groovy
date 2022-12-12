#!/usr/bin/env groovy

def success() {
    // config notify
    if(Config.notifySettings) {
        String successHeadMsg = "✅✅✅✅✅✅✅✅✅"
        String successStatusMsg = "构建成功✅"
        util.notify(successHeadMsg, successStatusMsg)
    }
}

def failure() {
    // config notify
    if(Config.notifySettings) {
        String failureHeadMsg = "❌❌❌❌❌❌❌❌❌"
        String failureStatusMsg = "构建失败❌"
        util.notify(failureHeadMsg, failureStatusMsg)
    }
}

def aborted() {
    // config notify
    if(Config.notifySettings) {
        String failureHeadMsg = "🟠🟠🟠🟠🟠🟠🟠"
        String failureStatusMsg = "构建中断🟠"
        util.notify(failureHeadMsg, failureStatusMsg)
    }
}
