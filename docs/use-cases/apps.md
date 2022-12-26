# GitOps with "Apps"

## 0 Goal

In this tutorial, we will try to use DevStream's new feature "Apps" to achieve similar result to the [GitOps](./gitops.md), but using much less configuration, to show the power of "Apps" in the config. If you haven't read the original GitOps best practice, click the above link first.

Two applications will be created (one Python, one Golang), CI/CD pipelines will be set up for both, and both apps will be deployed via Argo CD, just like the GitOps best practice.

---

## 1 Overview

DevStream will use the "Apps" concept to create repositories with scaffolding code, GitHub Actions CI pipelines, install Argo CD, and deploy the apps via Argo CD.

---

## 2 Create the Config File

Create a temporary working directory for this tutorial:

```bash
mkdir test
cd test/
```

Download dtm (see the [GitOps](./gitops.md) best practice if you haven't).

Then generate the config file by running:

```bash
./dtm show config --template=apps > config.yaml
```

Then set the following environment variables by running (replace values within the double quotes):

```shell
export GITHUB_USER="<YOUR_GITHUB_PERSONAL_ACCESS_TOKEN_HERE>"
export DOCKERHUB_USERNAME="<YOUR_DOCKER_HUB_USER_NAME_HERE>"
```

Then we run the following commands to update our config file with those env vars:

===  "**macOS** or **FreeBSD** based systems"

    ```shell
    sed -i.bak "s@YOUR_GITHUB_USER@${GITHUB_USER}@g" config.yaml
    sed -i.bak "s@YOUR_DOCKERHUB_USER@${DOCKERHUB_USERNAME}@g" config.yaml
    ```

=== "**GNU** Linux users"
 
    ```shell
    sed -i "s@YOUR_GITHUB_USERNAME_CASE_SENSITIVE@${GITHUB_USER}@g" config.yaml
    sed -i "s@YOUR_DOCKER_USERNAME@${DOCKERHUB_USERNAME}@g" config.yaml
    ```

---

## 3 Init and Apply

Run:

```bash
./dtm init -f config.yaml
```

This downloads the required plugins, according to the config file, automatically.

Then we apply it by running:

```bash
./dtm apply -f config.yaml -y
```

You will see similar outputs as the following:

```bash
tiexin@mbp ~/work/devstream-io/test $ ./dtm apply -f config.yaml -y
2022-12-16 14:30:27 ℹ [INFO]  Apply started.
2022-12-16 14:30:28 ℹ [INFO]  Using local backend. State file: devstream.state.
2022-12-16 14:30:28 ℹ [INFO]  Tool (helm-installer/argocd) found in config but doesn't exist in the state, will be created.
2022-12-16 14:30:28 ℹ [INFO]  Tool (repo-scaffolding/myapp1) found in config but doesn't exist in the state, will be created.
2022-12-16 14:30:28 ℹ [INFO]  Tool (repo-scaffolding/myapp2) found in config but doesn't exist in the state, will be created.
2022-12-16 14:30:28 ℹ [INFO]  Tool (github-actions/myapp1) found in config but doesn't exist in the state, will be created.
2022-12-16 14:30:28 ℹ [INFO]  Tool (argocdapp/myapp1) found in config but doesn't exist in the state, will be created.
2022-12-16 14:30:28 ℹ [INFO]  Tool (github-actions/myapp2) found in config but doesn't exist in the state, will be created.
2022-12-16 14:30:28 ℹ [INFO]  Tool (argocdapp/myapp2) found in config but doesn't exist in the state, will be created.
2022-12-16 14:30:28 ℹ [INFO]  Start executing the plan.
2022-12-16 14:30:28 ℹ [INFO]  Changes count: 7.
2022-12-16 14:30:28 ℹ [INFO]  -------------------- [  Processing progress: 1/7.  ] --------------------
2022-12-16 14:30:28 ℹ [INFO]  Processing: (helm-installer/argocd) -> Create ...
2022-12-16 14:30:29 ℹ [INFO]  Filling default config with instance: argocd.
2022-12-16 14:30:29 ℹ [INFO]  Creating or updating helm chart ...
... (omitted)
... (omitted)
2022-12-16 14:32:09 ℹ [INFO]  -------------------- [  Processing progress: 7/7.  ] --------------------
2022-12-16 14:32:09 ℹ [INFO]  Processing: (argocdapp/myapp2) -> Create ...
2022-12-16 14:32:19 ℹ [INFO]  application.argoproj.io/myapp2 created
2022-12-16 14:32:19 ✔ [SUCCESS]  Tool (argocdapp/myapp2) Create done.
2022-12-16 14:32:19 ℹ [INFO]  -------------------- [  Processing done.  ] --------------------
2022-12-16 14:32:19 ✔ [SUCCESS]  All plugins applied successfully.
2022-12-16 14:32:19 ✔ [SUCCESS]  Apply finished.
```

---

## 4 Check the Results

Let's continue to look at the results of the `apply` command.

Similar to what we did in the [GitOps](./gitops.md) best practice, we can check that repositories for two applications are created, CI pipelins are created for both apps too, Argo CD is installed, and both apps are deployed by Argo CD into our Kubernetes cluster.

---

## 5 Clean Up

Run:

```bash
./dtm delete -f config.yaml -y
```

And you will get similar outputs to the following:

```bash
2022-12-16 14:34:30 ℹ [INFO]  Delete started.
2022-12-16 14:34:31 ℹ [INFO]  Using local backend. State file: devstream.state.
2022-12-16 14:34:31 ℹ [INFO]  Tool (github-actions/myapp1) will be deleted.
2022-12-16 14:34:31 ℹ [INFO]  Tool (argocdapp/myapp1) will be deleted.
2022-12-16 14:34:31 ℹ [INFO]  Tool (github-actions/myapp2) will be deleted.
2022-12-16 14:34:31 ℹ [INFO]  Tool (argocdapp/myapp2) will be deleted.
2022-12-16 14:34:31 ℹ [INFO]  Tool (repo-scaffolding/myapp1) will be deleted.
2022-12-16 14:34:31 ℹ [INFO]  Tool (repo-scaffolding/myapp2) will be deleted.
2022-12-16 14:34:31 ℹ [INFO]  Tool (helm-installer/argocd) will be deleted.
2022-12-16 14:34:31 ℹ [INFO]  Start executing the plan.
2022-12-16 14:34:31 ℹ [INFO]  Changes count: 7.
2022-12-16 14:34:31 ℹ [INFO]  -------------------- [  Processing progress: 1/7.  ] --------------------
2022-12-16 14:34:31 ℹ [INFO]  Processing: (github-actions/myapp1) -> Delete ...
2022-12-16 14:34:33 ℹ [INFO]  Prepare to delete 'github-actions_myapp1' from States.
2022-12-16 14:34:33 ✔ [SUCCESS]  Tool (github-actions/myapp1) delete done.
... (omitted)
... (omitted)
2022-12-16 14:34:40 ✔ [SUCCESS]  Tool (helm-installer/argocd) delete done.
2022-12-16 14:34:40 ℹ [INFO]  -------------------- [  Processing done.  ] --------------------
2022-12-16 14:34:40 ✔ [SUCCESS]  All plugins deleted successfully.
2022-12-16 14:34:40 ✔ [SUCCESS]  Delete finished.
```

Then you can delete what we created:

```bash
cd ../
rm -rf test/
rm -rf ~/.devstream/
```
