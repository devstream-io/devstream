# GitOps Toolchain

If you are interested in watching a video demo, see the youtube video below:

<iframe width="100%" height="500" src="https://www.youtube.com/embed/q7TK3vFr1kg" title="YouTube video player" frameborder="0" allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture" allowfullscreen></iframe>

[https://www.youtube.com/watch?v=q7TK3vFr1kg](https://www.youtube.com/watch?v=q7TK3vFr1kg)

For Chinese readers, watch this one instead:

<iframe src="//player.bilibili.com/player.html?aid=426762434&bvid=BV1W3411P7oW&cid=728576152&high_quality=1&danmaku=0" allowfullscreen="allowfullscreen" width="100%" height="500" scrolling="no" frameborder="0" sandbox="allow-top-navigation allow-same-origin allow-forms allow-scripts"></iframe>



[https://www.bilibili.com/video/BV1W3411P7oW/](https://www.bilibili.com/video/BV1W3411P7oW/)

## Plugins needed

1. [repo-scaffolding](../plugins/repo-scaffolding.md)
2. [jira-github](../plugins/jira-github-integ.md)
3. [githubactions-golang](../plugins/githubactions-golang.md)
4. [argocd](../plugins/argocd.md)
5. [argocdapp](../plugins/argocdapp.md)

The dependencies of these plugins are(`a -> b` means for `a depends on b`):

- `jira-github` -> `repo-scaffolding`
- `githubactions-golang` -> `repo-scaffolding`
- `argocdapp` -> `argocd`, `githubactions-golang` and `repo-scaffolding`

Note: These dependencies are optional; you can use dependency to make sure a certain tool is installed before another. We should use dependency according to the actual usage situation.

## 1 Download DevStream (`dtm`)

In your working directory, run:

```shell
sh -c "$(curl -fsSL https://raw.githubusercontent.com/devstream-io/devstream/main/hack/install/download.sh)"
```

This will download the corresponding `dtm` binary to your working directory according to your OS and chip architecture, and grant the binary execution permission.

> Optional: you can then move `dtm` to a place which is in your PATH. For example: `mv dtm /usr/local/bin/`.

_For more details on how to install, see [install dtm](../install.md)._

## 2 Prepare the Config File

Run the following command to generate a template configuration file for gitops `gitops.yaml`.

```shell
./dtm show config -t gitops > gitops.yaml
```

Then modify the `gitops.yaml` file accordingly.

For me I can set these variables like:

| Variable                       | Example           | Note                                                         |
| ------------------------------ | ----------------- | ------------------------------------------------------------ |
| defaultBranch                  | main              | The branch name you want to use |
| githubUsername                 | daniel-hutao      | It should be case-sensitive here; strictly use your GitHub username |
| repoName                       | go-webapp         | As long as it doesn't exist in your GitHub account and the name is legal |
| dockerhubUsername              | exploitht         | It should be case-sensitive here; strictly use your DockerHub username |
| jiraID                         | merico            | This is a domain name prefix like merico in https://merico.atlassian.net |
| jiraProjectKey                 | DT                | A descriptive prefix for your project’s issue keys to recognize work from this project |
| jiraUserEmail                  | tao.hu@merico.dev | The email you use to log in to Jira |
| argocdNameSpace                | argocd            | The namespace used by ArgoCD |
| argocdDeployTimeout            | 10m               | How long does ArgoCD deployment timeout |


These plugins require some environment variables to work, so let's set them:

```bash
export GITHUB_TOKEN="YOUR_GITHUB_TOKEN_HERE"
export JIRA_API_TOKEN="YOUR_JIRA_API_TOKEN_HERE"
export DOCKERHUB_TOKEN="YOUR_DOCKERHUB_TOKEN_HERE"
```

If you don't know how to create these three tokens, check out:

- GITHUB_TOKEN: [Manage API tokens for your Atlassian account](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token)
- JIRA_API_TOKEN: [Creating a personal access token](https://support.atlassian.com/atlassian-account/docs/manage-api-tokens-for-your-atlassian-account/)
- DOCKERHUB_TOKEN: [Manage access tokens](https://docs.docker.com/docker-hub/access-tokens/)

## 3 Initialize

Run:

```bash
dtm init -f gitops.yaml
```

## 4 Apply

Run:

```bash
dtm apply -f gitops.yaml
```

and confirm to continue, then you should see similar output to:

```
...
2022-03-11 13:36:11 ✔ [SUCCESS]  All plugins applied successfully.
2022-03-11 13:36:11 ✔ [SUCCESS]  Apply finished.
```

## 5 Check the Results

Let's continue to look at the results of the `apply` command.

### 5.1 Repository Scaffolding

- The repository scaffolding we got looks like this:

![](gitops/a.png)

### 5.2 Jira-Github Integration

- How do Jira and Github integrate? Let's create a new issue:

![](gitops/b.png)

- The issue will be renamed automatically like this:

![](gitops/c.png)

- We can find this auto-synced `Story` in Jira:

![](gitops/d.png)

- If we continue to leave a comment on this issue:

![](gitops/d1.png)

- The comment will also be automatically synced to Jira:

![](gitops/e.png)

### 5.3 GitHub Actions CI for Golang

- What does CI do here?

![](gitops/f.png)

- The CI processes also build an image, and this image is automatically pushed to DockerHub:

![](gitops/g.png)

### 5.4 ArgoCD Deployment

- Of course, the ArgoCD must have been installed as expected.

![](gitops/h.png)

### 5.5 ArgoCD Application Deployment

- Our code has just been built into an image, at this time the image is automatically deployed to our k8s as a Pod:

![](gitops/i.png)

## 6 Clean Up

Run:

```bash
dtm destroy -f gitops.yaml
```

and you should see similar output:

```
2022-03-11 13:39:11 ✔ [SUCCESS]  All plugins destroyed successfully.
2022-03-11 13:39:11 ✔ [SUCCESS]  Destroy finished.
```
