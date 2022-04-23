# *vault* Plugin

This plugin installs [vault](https://www.vaultproject.io/) in an existing Kubernetes cluster using the Helm chart for you test or develop vault.

### Usage

```yaml
tools:
- name: vault
  # name of the plugin
  plugin: vault
  options:
    # need to create the namespace or not, default: false
    create_namespace: true
    repo:
      # name of the Helm repo
      name: hashicorp
      # url of the Helm repo
      url: https://helm.releases.hashicorp.com
    # Helm chart information
    chart:
      # name of the chart
      chart_name: hashicorp/vault
      # release name of the chart
      release_name: vault
      # k8s namespace where Vault will be installed
      namespace: hashicorp-vault
      # whether to wait for the release to be deployed or not
      wait: true
      # the time to wait for any individual Kubernetes operation (like Jobs for hooks). This defaults to 5m0s
      timeout: 5m
      values_yaml: |
        global:
          enabled: true
        server:
          affinity: ""
          ha:
            enabled: true
            replicas: 3
            raft:
              enabled: true
              setNodeId: true
          namespaceSelector:
            matchLabels:
              injection: enabled
```

Currently, except for `values_yaml`, all the parameters in the example above are mandatory.

After vault pods deployed  by the plugin vault, you can follow the instructions below to initialize the vault:
Before you follow the instructions bellow, you must install `jq` tool: a command tool to analyze json

If you work on MacOS:
```
brew install jq
```

If you work on linux OS like CentOS:
```
sudo yum install jq
```

If you work on Linux OS like Ubuntu:
```
sudo apt-get install jq
```


**in the  command below, you should replace the variable $NAMESPACE with "hashicorp-vault" if you not modify the namespace variable.**
**Otherwise, use the namespace name you replaced**

1. 
```
   # Initialize Vault with one key share and one key threshold.
   kubectl exec vault-0 -n $NAMESPACE -- vault operator init -key-shares=1 -key-threshold=1 -format=json > cluster-keys.json
```
2. 
  ```
   # Display the unseal key found in cluster-keys.json
   cat cluster-keys.json | jq -r ".unseal_keys_b64[]"
   ```
3. 
```
   # Create a variable named VAULT_UNSEAL_KEY to capture the Vault unseal key.
   VAULT_UNSEAL_KEY=$(cat cluster-keys.json | jq -r ".unseal_keys_b64[]")
```
4. 
```
   # Unseal Vault running on the vault-0 pod.
   kubectl exec vault-0  -n $NAMESPACE -- vault operator unseal $VAULT_UNSEAL_KEY
```
you will see the above command's output like this, make sure the item "Initialized" value is "true" and the item "Sealed" value is "false"
```shell
Key                     Value
---                     -----
Seal Type               shamir
Initialized             true
Sealed                  false
Total Shares            1
Threshold               1
Version                 1.9.2
Storage Type            raft
Cluster Name            vault-cluster-14052440
Cluster ID              7630cd33-2ee1-39c1-db3f-e48a6d79970a
HA Enabled              true
HA Cluster              n/a
HA Mode                 standby
Active Node Address     <none>
Raft Committed Index    25
Raft Applied Index      25
```

5. Initialize vault-1 and  vault-2 like vault-0:

```shell
# Initialize vault-1
kubectl exec vault-1 -n hashicorp-vault -- vault operator init -key-shares=1 -key-threshold=1 -format=json > cluster-keys.json
VAULT_UNSEAL_KEY=$(cat cluster-keys.json | jq -r ".unseal_keys_b64[]")
kubectl exec vault-1  -n hashicorp-vault -- vault operator unseal $VAULT_UNSEAL_KEY
# Initialize vault-1
kubectl exec vault-2 -n hashicorp-vault -- vault operator init -key-shares=1 -key-threshold=1 -format=json > cluster-keys.json
VAULT_UNSEAL_KEY=$(cat cluster-keys.json | jq -r ".unseal_keys_b64[]")
kubectl exec vault-2  -n hashicorp-vault -- vault operator unseal $VAULT_UNSEAL_KEY
```

6
```
   # Verify  all the Vault pods are running and ready.
   kubectl get pods -n $NAMESPACE
```

you will see the above command's output like this bellow, make sure all the pods are running and ready 
```
NAME                                 READY   STATUS    RESTARTS   AGE
vault-0                              1/1     Running   0          2m29s
vault-1                              1/1     Running   0          2m29s
vault-2                              1/1     Running   0          2m29s
vault-agent-injector-68dc986-bnsj2   1/1     Running   0          2m28s
```
7. After these above operations, you will want to use vault to write/read/delete secrets, how to do this?
   You should follow the hashicorp vault documentation:
 - https://learn.hashicorp.com/tutorials/vault/kubernetes-minikube?in=vault/kubernetes#set-a-secret-in-vault
 - https://learn.hashicorp.com/tutorials/vault/getting-started-first-secret