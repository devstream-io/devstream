# dtm-pipeline-templates

This Repo is used to store devstream's best pipeline practice templates, If you like the template and want to use it, a fork is welcome.

## Pipeline Templates

It support pipeline templates below for now：

| ciType                    | configPath                        | description                                                                      |
| ----                      | ----                              | ----                                                                             |
| Github actions            | github-actions/workflows          | github actions template, support test, build, pushImage, notify stage            |
| Gitlab ci                 | gitlab-ci/.gitlab-ci.yml          | gitlab ci template, support test, build, pushImage, notify stage                 |
| Genkins pipeline          | jenkins-pipeline/Jenkinsfile      | jenkins-pipeline template, support test, build, pushImage, notify stage          |
| Argocd helm               | argocdapp/helm                    | argocdapp helm config                                                            |
| Argocd kustomize          | argocdapp/kustomize               | argocd kustomize config                                                          |

## Where does this repo come from?

`dtm-pipeline-templates` is synced from https://github.com/devstream-io/devstream/blob/main/staging/dtm-pipeline-templates.
Code changes are made in that location, merged into `devstream-io/devstream` and later synced here.
