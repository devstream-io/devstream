# jenkins Plugin

This plugin installs [jenkins](https://jenkins.io) in an existing Kubernetes cluster using the Helm chart.

## Usage

### NOTICE

#### Production Environments
Please be sure to change the `storageClass` in the configuration item to an existing StorageClass.

#### Test Environments
If you want to **test** the plugin in a **local environment**:

1. Please change the `test_env` to `true` in the config file.
2. Create the data directory manually in the host where kubernetes is running:

```bash
# If you run k8s and dtm in the same host
mkdir -p ~/data/jenkins-volume/
chown -R 1000:1000 ~/data/jenkins-volume/

-------------------------

# If you run k8s and dtm in different hosts, such as run k8s in docker
# 1 Get the home directory of the user who runs dtm
cd ~ && pwd
# 2 Copy the result of the above command
# 3 Enter the host on which k8s is running
  # 3.1 for minikube
  minikube ssh
  # 3.2 for kind
  docker exec -it <kind-container-name-or-id> bash
# 4 Create the data directory manually in the host where kubernetes is running:
mkdir -p <your-dtm-home-dir>/data/jenkins-volume/
chown -R 1000:1000 <your-dtm-home-dir>/data/jenkins-volume/
```

### Config

```yaml
--8<-- "jenkins.yaml"
```

Currently, all the parameters in the example above are mandatory.
