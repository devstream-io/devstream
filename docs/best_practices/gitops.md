# GitOps Toolchain

## Plugins needed

1. [github-repo-scaffolding-golang](../plugins/github-repo-scaffolding-golang.md)
2. [jira-github](../plugins/jira-github-integ.md)
3. [githubactions-golang](../plugins/githubactions-golang.md)
4. [argocd](../plugins/argocd.md)
5. [argocdapp](../plugins/argocdapp.md)

The dependencies of these plugins are(`a -> b` means for `a depends on b`):

- `jira-github` -> `github-repo-scaffolding-golang`
- `githubactions-golang` -> `github-repo-scaffolding-golang`
- `argocdapp` -> `argocd`

Note: These dependencies are not consistent, such as when the repo operated by `jira-github` and `github-repo-scaffolding-golang` are not the same, the dependencies disappear.

We should use the `dependency` according to the actual usage situation.

## Download DevStream (`dtm`)

Download the appropriate `dtm` version for your platform from [DevStream Releases](https://github.com/devstream-io/devstream/releases).

> Remember to rename the binary file to `dtm` so that it's easier to use. For example: `mv dtm-darwin-arm64 dtm`.

> Once downloaded, you can run the binary from anywhere. Ideally, you want to put it in a place that is in your PATH (e.g., `/usr/local/bin`).

## Prepare the Config File

Copy the [examples/gitops.yaml](../../examples/gitops.yaml) to your working directory:

```bash
cp examples/gitops.yaml config-gitops.yaml
```

Then modify the file accordingly. The variables you need to care about are the following: 

- YOUR_GITHUB_USERNAME
- YOUR_REPO_NAME
- YOUR_DOCKERHUB_USERNAME
- YOUR_DOCKERHUB_IMAGE_REPO_NAME
- JIRA_ID
- JIRA_USER_EMAIL
- JIRA_PROJECT_KEY
- YOUR_CHART_REPO_URL

For me I can set these variables like:

| Variable                       | Example           | Note                                                         |
| ------------------------------ | ----------------- | ------------------------------------------------------------ |
| YOUR_GITHUB_USERNAME           | daniel-hutao      | It should be case-sensitive here; strictly use your GitHub username |
| YOUR_REPO_NAME                 | go-webapp         | As long as it doesn't exist in your GitHub account and the name is legal |
| YOUR_DOCKERHUB_USERNAME        | exploitht         | It should be case-sensitive here; strictly use your DockerHub username |
| YOUR_DOCKERHUB_IMAGE_REPO_NAME | go-webapp         | It is recommended that you use the same name as the project, please make sure that there is no project with the same name in your DockerHub account |
| JIRA_ID                        | merico            | This is a domain name prefix like merico in https://merico.atlassian.net |
| JIRA_PROJECT_KEY               | DT                | A descriptive prefix for your project’s issue keys to recognize work from this project |
| JIRA_USER_EMAIL                | tao.hu@merico.dev | The email you use to log in to Jira |
| YOUR_CHART_REPO_URL            | https://github.com/daniel-hutao/dtm-gitops-test.git | Helm chart repo URL |

These plugins require some environment variables to work, so let's set them:

```bash
export GITHUB_TOKEN="YOUR_GITHUB_TOKEN_HERE"
export JIRA_API_TOKEN="YOUR_JIRA_API_TOKEN_HERE"
```

If you don't know how to create these two tokens, check out:

- [Manage API tokens for your Atlassian account](https://support.atlassian.com/atlassian-account/docs/manage-api-tokens-for-your-atlassian-account/)
- [Creating a personal access token](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token)

## 3. Initialize

Run:

```bash
dtm init -f config.yaml
```

## 4. Apply

Run:

```bash
dtm apply -f config.yaml
```

and confirm to continue, then you should see similar output to:

```
...
2022-03-11 13:36:11 ✔ [SUCCESS]  All plugins applied successfully.
2022-03-11 13:36:11 ✔ [SUCCESS]  Apply finished.
```

## 5. Check the Results

// TODO(daniel-hutao): I'll add the docs here later.

## 6. Clean Up

Run:

```bash
dtm destroy
```

and you should see similar output:

```
2022-03-11 13:39:11 ✔ [SUCCESS]  All plugins destroyed successfully.
2022-03-11 13:39:11 ✔ [SUCCESS]  Destroy finished.
```

```{toctree}
---
maxdepth: 1
---
```