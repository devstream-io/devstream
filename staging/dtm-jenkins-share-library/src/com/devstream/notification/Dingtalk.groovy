package com.devstream.notification


def send(changeString, headMessage, statusMessage, Integer _timeout=60) {
    // String buildUser = variable.buildUserName()
    String notifyUser = Config.notifySettings.get("at_user")
    String robotID = Config.notifySettings.get("robot_id")
    List<String> atUsers = [] as String[]
    if (notifyUser != null && notifyUser != "") {
        atUsers = notifyUser.split(",") as String[]
    }
    timeout(time: _timeout, unit: 'SECONDS') {
             dingtalk (
              robot: "${robotID}",
              type: 'MARKDOWN',
              title: "${env.JOB_NAME}[${env.BRANCH_NAME}]构建通知",
              text: [
                  "# $headMessage",
                  "# 构建详情",
                  "- 构建变更: ${changeString}",
                  "- 构建结果: ${statusMessage}",
                  // "- 构建人: **${buildUser}**",
                  "- 持续时间: ${currentBuild.durationString}",
                  "# 构建日志",
                  "[日志](${env.BUILD_URL}console)"
              ],
              at: atUsers
            )
    }

}
