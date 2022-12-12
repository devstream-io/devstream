# dtm-jenkins-share-library

This repo contains share library used by DevStream plugin "jenkins-pipeline" (thereafter: the plugin). It currently has the following functions：
- send notification to dingtalk when pipeline result is success or failed.
- run test for java project.
- run sonar scanner for java project.
- build docker image by current project, and push this image to image repo.

## Example
```groovy dingtalk
/*
config dingtalk related info, if not set, pipeline result will not be sent
required plugins: dingding-notifications
*/
setting.configNotifyDingtalk([
    'robot_id': "dingdingtest", // robotID in jenkins config
    'at_user': ""
])

/*
config docker image repo info, if not set, image will not be sent to repo
*/
setting.configImageRepo([
    'image_repo': "test.com/library",
    'image_auth_secret_name': "docker-config",
])

/*
pipeline generic config
*/
runPipeline([
    'enable_test': true, // whether run test for code
    'name': "spring-test-github", // this name will be used for image repo name and sonar project name
    'enable_sonarqube': true, // whether use sonar to scan code
])

```

## General config variables

|  field   | description  | default_value  |
|  ----  | ----  | ----  |
|  repo_type  | whether code repo is github or gitlab, if repo_type is gitlab, we can use gitlab_connection to show jenkins pipeline status in gitlab  |   |
|  name  | this name will be used for image repo name and sonar project name  |  |
|  language  | the project language, current only support Java  | java |
|  container_requests_cpu  | jenkins worker container requests cpu  | 0.3 |
|  container_requests_memory  | jenkins worker container requests memory  | 512Mi |
|  container_limit_cpu  | jenkins worker container limit cpu  | 1 |
|  container_limit_memory  | jenkins worker container limit memory  | 2Gi |
|  enable_test  | whether run test for code  | true |
|  enable_sonarqube  | whether use sonar to scan code  | false |

## Where does this repo come from?

`dtm-jenkins-share-library` is synced from https://github.com/devstream-io/devstream/blob/main/staging/dtm-jenkins-share-library.
Code changes are made in that location, merged into `devstream-io/devstream` and later synced here.
