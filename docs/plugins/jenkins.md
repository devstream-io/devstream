# jenkins Plugin

This plugin installs [jenkins](https://jenkins.io) in an existing Kubernetes cluster using the Helm chart.

## Usage

NOTICE:

Create the data directory manually:

```bash
mkdir -p ~/data/jenkins-volumes/
chown -R 1000:1000 ~/data/jenkins-volumes/
```

```yaml
--8<-- "jenkins.yaml"
```

Currently, all the parameters in the example above are mandatory.
