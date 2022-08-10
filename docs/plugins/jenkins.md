# jenkins Plugin

This plugin installs [Jenkins](https://jenkins.io) in an existing Kubernetes cluster using the Helm chart.

It also installs [GitHub Pull Request Builder(ghprb)](https://plugins.jenkins.io/ghprb/) and [OWASP Markup Formatter](https://plugins.jenkins.io/antisamy-markup-formatter/) plugins. Then enable HTML parsing using OWASP Markup Formatter Plugin , useful with ghprb plugin.

## Usage

### Production Environment

Please be sure to change the `storageClass` in the options of the config to an existing StorageClass.

### Test/Local Dev Environment

If you want to **test** the plugin in a **local environment**:

1. Please change the `test_env` to `true` in the config file.
2. Create the data directory manually in the host where Kubernetes is running. Here's how:

If you run Kubernetes and `dtm` on the same host:

```bash
mkdir -p ~/data/jenkins-volume/
chown -R 1000:1000 ~/data/jenkins-volume/
```

Or, if you run Kubernetes and `dtm` on different "hosts," such as running Kubernetes in an VM or a Docker container:

```bash
# 1 get the home directory of the user who runs dtm
cd ~ && pwd
# 2 copy the result of the above command
# 3 enter the host on which k8s is running
  # 3.1 for minikube
  minikube ssh
  # 3.2 for kind
  docker exec -it <kind-container-name-or-id> bash
# 4 Create the data directory manually in the host where k8s is running:
mkdir -p <your-dtm-home-dir>/data/jenkins-volume/
chown -R 1000:1000 <your-dtm-home-dir>/data/jenkins-volume/
```

### Config

```yaml
--8<-- "jenkins.yaml"
```

#### Default Configs

| key                | default value             | description                                        |
| ----               | ----                      | ----                                               |
| chart.chart_name   | jenkins/jenkins           | chart name                                         |
| chart.timeout      | 5m                        | this config will wait 5 minutes to deploy          |
| chart.upgradeCRDs  | true                      | default update CRD config                          |
| chart.release_name | jenkins                   | helm release name                                  |
| chart.wait         | true                      | whether to wait until installation is complete     |
| repo.url           | https://charts.jenkins.io | helm official repo address                         |
| repo.name          | jenkins                   | helm repo name                                     |
| create_namespace   | false                     | whether to create namespace if namespace not eixst |


Currently, expect default configs all the parameters in the example above are mandatory.

## Outputs

This plugin has two outputs:

- `jenkinsURL` (format: `hostname:port`, example: "localhost:8080")
- `jenkinsPasswordOfAdmin` 
