# Plugins List

| Type                           | Plugin                         | Note                           | Usage/Doc                                |
|--------------------------------|--------------------------------|--------------------------------|------------------------------------------|
| Issue Tracking                 | trello-github-integ            | Trello/GitHub integration      | [doc](trello-github-integ.md)            |
| Issue Tracking                 | trello                         | Trello                         | [doc](trello.md)                         |
| Issue Tracking                 | jira-github-integ              | Jira/GitHub integration        | [doc](jira-github-integ.md)              |
| Source Code Management         | github-repo-scaffolding-golang | Go WebApp scaffolding          | [doc](github-repo-scaffolding-golang.md) |
| Source Code Management         | gitlab-repo-scaffolding-golang | Go WebApp scaffolding          | [doc](gitlab-repo-scaffolding-golang.md) |
| CI                             | jenkins                        | Jenkins installation           | [doc](jenkins.md)                        |
| CI                             | githubactions-golang           | GitHub Actions CI for Golang   | [doc](githubactions-golang.md)           |
| CI                             | githubactions-python           | GitHub Actions CI for Python   | [doc](githubactions-python.md)           |
| CI                             | githubactions-nodejs           | GitHub Actions CI for Nodejs   | [doc](githubactions-nodejs.md)           |
| CI                             | gitlabci-golang                | GitLab CI for Golang           | [doc](gitlabci-golang.md)                |
| CI                             | gitlabci-generic               | Generic GitLab CI              | [doc](gitlabci-generic.md)               |
| CD/GitOps                      | argocd                         | ArgoCD installation            | [doc](argocd.md)                         |
| CD/GitOps                      | argocdapp                      | ArgoCD Application creation    | [doc](argocdapp.md)                      |
| Monitoring                     | kube-prometheus                | Prometheus/Grafana K8s install | [doc](kube-prometheus.md)                |
| Observability                  | devlake                        | DevLake installation           | [doc](devlake.md)                        |
| LDAP                           | openldap                       | OpenLDAP installation          | [doc](openldap.md)                       |
| Secrets/Credentials Management | hashicorp-vault                | Hashicorp Vault installation   | [doc](hashicorp-vault.md)                |

Or, to get a list of plugins, run:

```shell
$ dtm list plugins
argocd
argocdapp
devlake
github-repo-scaffolding-golang
githubactions-golang
githubactions-nodejs
githubactions-python
gitlabci-generic
gitlabci-golang
hashicorp-vault
jenkins
jira-github-integ
kube-prometheus
openldap
trello
trello-github-integ
```
