# Plugins List

| Type                           | Plugin                      | Note                           | Usage/Doc                             |
|--------------------------------|-----------------------------|--------------------------------|---------------------------------------|
| Issue Tracking                 | trello-github-integ         | Trello/GitHub integration      | [doc](trello-github-integ.md)         |
| Issue Tracking                 | trello                      | Trello                         | [doc](trello.md)                      |
| Issue Tracking                 | jira-github-integ           | Jira/GitHub integration        | [doc](jira-github-integ.md)           |
| Issue Tracking                 | zentao                      | Zentao installation            | [doc](zentao.md)                      |
| Source Code Management         | repo-scaffolding            | App scaffolding                | [doc](repo-scaffolding.md)            |
| Source Code Management         | gitlab-ce-docker            | GitLab CE version installation | [doc](gitlab-ce-docker.md)            |
| CI                             | jenkins                     | Jenkins installation           | [doc](jenkins.md)                     |
| CI                             | jenkins-pipeline            | Jenkins pipeline creation      | [doc](jenkins-pipeline.md)            |
| CI                             | githubactions-golang        | GitHub Actions CI for Golang   | [doc](githubactions-golang.md)        |
| CI                             | githubactions-python        | GitHub Actions CI for Python   | [doc](githubactions-python.md)        |
| CI                             | githubactions-nodejs        | GitHub Actions CI for Nodejs   | [doc](githubactions-nodejs.md)        |
| CI                             | gitlabci-golang             | GitLab CI for Golang           | [doc](gitlabci-golang.md)             |
| CI                             | gitlabci-java               | GitLab CI for Java             | [doc](gitlabci-java.md)               |
| CI                             | gitlabci-generic            | Generic GitLab CI              | [doc](gitlabci-generic.md)            |
| CI                             | ci-generic                  | Generic CI plugin              | [doc](ci-generic.md)                  |
| CI                             | tekton                      | Tekton CI installation         | [doc](tekton.md)                      |
| Code Quality/Security          | sonarqube                   | SonarQube                      | [doc](sonarqube.md)
| CD/GitOps                      | argocd                      | ArgoCD installation            | [doc](argocd.md)                      |
| CD/GitOps                      | argocdapp                   | ArgoCD Application creation    | [doc](argocdapp.md)                   |
| Image Repository               | artifactory                 | Artifactory installation       | [doc](artifactory.md)                 |
| Image Repository               | harbor                      | Harbor helm installation       | [doc](harbor.md)                      |
| Image Repository               | harbor-docker               | Harbor Docker compose install  | [doc](harbor-docker.md)               |
| Deployment                     | helm-generic                | Helm chart install             | [doc](helm-generic.md)               |
| Monitoring                     | kube-prometheus             | Prometheus/Grafana K8s install | [doc](kube-prometheus.md)             |
| Observability                  | devlake                     | DevLake installation           | [doc](devlake.md)                     |
| LDAP                           | openldap                    | OpenLDAP installation          | [doc](openldap.md)                    |
| Secrets/Credentials Management | hashicorp-vault             | Hashicorp Vault installation   | [doc](hashicorp-vault.md)             |

Or, to get a list of plugins, run:

```shell
$ dtm list plugins
argocd
argocdapp
artifactory
ci-generic
devlake
githubactions-golang
githubactions-nodejs
githubactions-python
gitlab-ce-docker
gitlabci-generic
gitlabci-golang
gitlabci-java
harbor
harbor-docker
hashicorp-vault
helm-generic
jenkins
jenkins-pipeline
jira-github-integ
kube-prometheus
openldap
repo-scaffolding
sonarqube
tekton
trello
trello-github-integ
zentao
```
