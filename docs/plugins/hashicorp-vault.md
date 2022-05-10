# hashicorp-vault Plugin

This plugin installs [hashicorp-vault](https://www.vaultproject.io/) in an existing Kubernetes cluster using the Helm chart for your tests or develops hashicorp-vault.

This plugin installs hashicorp-vault with replicas:3 by default value.

## Usage

```yaml
tools:
# name of the tool
- name: hashicorp-vault
  # the id of the tool instance
  instanceID: default
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
      # The k8s namespace is where you deploy the Vault to k8s
      namespace: hashicorp
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

## Initialize all the Vault pods

After installing the Vault on k8s, you can initialize all pods of the Vault on k8s. To know more about the Vault, you can refer to:

- [Vault init](https://www.vaultproject.io/docs/commands/operator/init)
- [Vault:Seal/Unseal](https://www.vaultproject.io/docs/concepts/seal)

At first, you must install [jq](https://stedolan.github.io/jq/) tool: jq is a lightweight and flexible command-line JSON processor.
[Download jq](https://stedolan.github.io/jq/download/)

In the command below, the variable `$NAMESPACE` you should replace with "hashicorp" if you do not modify the namespace variable.
Otherwise, use the namespace name you replaced.

1. Initialize vault-0
```shell
# Initialize vault-0 with one key share and one key threshold.
kubectl exec vault-0 -n $NAMESPACE -- vault operator init -key-shares=1 -key-threshold=1 -format=json > cluster-keys.json
```
2. Display the unseal key
```shell
# Display the unseal key found in cluster-keys.json
cat cluster-keys.json | jq -r ".unseal_keys_b64[]"
```
3. Create a variable to capture the Vault unseal key
```shell
# Create a variable named VAULT_UNSEAL_KEY to capture the Vault unseal key.
VAULT_UNSEAL_KEY=$(cat cluster-keys.json | jq -r ".unseal_keys_b64[]")
```
4. Unseal vault-0
```shell
# Unseal vault-0 running on the vault-0 pod.
kubectl exec vault-0  -n $NAMESPACE -- vault operator unseal $VAULT_UNSEAL_KEY
```
You will see the above command's output like this. Make sure the value of `Initialized` is 'true' and the value of `Sealed` is 'false'.
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
HA Cluster              https://vault-0.vault-internal:8201
HA Mode                 active
Active Since            2022-04-23T16:45:47.6060163Z
Raft Committed Index    30
Raft Applied Index      30
```

5. Initialize vault-1 and vault-2 like vault-0

```shell
# Initialize vault-1
kubectl exec vault-1 -n $NAMESPACE -- vault operator init -key-shares=1 -key-threshold=1 -format=json > cluster-keys.json
VAULT_UNSEAL_KEY=$(cat cluster-keys.json | jq -r ".unseal_keys_b64[]")
kubectl exec vault-1  -n $NAMESPACE -- vault operator unseal $VAULT_UNSEAL_KEY
# Initialize vault-2
kubectl exec vault-1 -n $NAMESPACE -- vault operator init -key-shares=1 -key-threshold=1 -format=json > cluster-keys.json
VAULT_UNSEAL_KEY=$(cat cluster-keys.json | jq -r ".unseal_keys_b64[]")
kubectl exec vault-1  -n $NAMESPACE -- vault operator unseal $VAULT_UNSEAL_KEY
```

6. Verify all the pods status
```shell
# Verify all the Vault pods are running and ready.
kubectl get pods -n $NAMESPACE
```

You will see the above command's outputs like this below. Make sure all the pods are running and ready.
```shell
NAME                                 READY   STATUS    RESTARTS   AGE
vault-0                              1/1     Running   0          2m29s
vault-1                              1/1     Running   0          2m29s
vault-2                              1/1     Running   0          2m29s
vault-agent-injector-68dc986-bnsj2   1/1     Running   0          2m28s
```

7. After the above operations, you want to use the Vault to write/read secrets. You need to follow the documentation of the hashicorp Vault:
- [Set a secret in Vault](https://learn.hashicorp.com/tutorials/vault/kubernetes-minikube?in=vault/kubernetes#set-a-secret-in-vault)
- [Your First Secret](https://learn.hashicorp.com/tutorials/vault/getting-started-first-secret)
