## 1 `argocd` Plugin

This plugin installs [OpenLDAP](https://argoproj.github.io/cd/) in an existing Kubernetes cluster using the Helm chart.

## 2 Usage:

```yaml
tools:
- name: openldap
  plugin:
    # name of the plugin
    kind: openldap
    # version of the plugin
    version: 0.2.0
  options:
    # need to create the namespace or not, default: false
    create_namespace: true
    repo:
      # name of the Helm repo
      name: helm-openldap
      # url of the Helm repo
      url: https://jp-gouin.github.io/helm-openldap/
    # Helm chart information
    chart:
      # name of the chart
      chart_name: helm-openldap/openldap-stack-ha
      # release name of the chart
      release_name: openldap
      # k8s namespace where OpenLDAP will be installed
      namespace: openldap
      # whether to wait for the release to be deployed or not
      wait: true
      # the time to wait for any individual Kubernetes operation (like Jobs for hooks). This defaults to 5m0s
      timeout: 5m
      # custom configuration (Optional). You can refer to https://github.com/jp-gouin/helm-openldap/blob/master/values.yaml
      values_yaml: |
        service: 
          type: NodePort
        env:
          LDAP_LOG_LEVEL: "256"
          LDAP_ORGANISATION: "Example Inc."
          LDAP_DOMAIN: "example.org"
          LDAP_READONLY_USER: "false"
          LDAP_READONLY_USER_USERNAME: "readonly"
          LDAP_READONLY_USER_PASSWORD: "readonly"
          LDAP_RFC2307BIS_SCHEMA: "false"
          LDAP_BACKEND: "mdb"
          LDAP_TLS: "true"
          LDAP_TLS_CRT_FILENAME: "tls.crt"
          LDAP_TLS_KEY_FILENAME: "tls.key"
          LDAP_TLS_DH_PARAM_FILENAME: "dhparam.pem"
          LDAP_TLS_CA_CRT_FILENAME: "ca.crt"
          LDAP_TLS_ENFORCE: "false"
          CONTAINER_LOG_LEVEL: "4"
          LDAP_TLS_REQCERT: "never"
          KEEP_EXISTING_CONFIG: "false"
          LDAP_REMOVE_CONFIG_AFTER_SETUP: "true"
          LDAP_SSL_HELPER_PREFIX: "ldap"
          LDAP_TLS_VERIFY_CLIENT: "never"
          LDAP_TLS_PROTOCOL_MIN: "3.0"
          LDAP_TLS_CIPHER_SUITE: "NORMAL"
        persistence:
          enabled: true
          storageClass: "alicloud-nas-subpath"
          accessModes:
            - ReadWriteOnce
          size: 8Gi
        adminPassword: Not@SecurePassw0rd
        configPassword: Not@SecurePassw0rd

        # you can modify user's password by this config
        ltb-passwd:
          enabled : true
          ingress:
            enabled: true
            annotations: {}
            path: /
            pathType: Prefix
            ## Ingress Host
            hosts:
            - "ssl-ldap2.example"
          ldap:
            server: ldap://openldap-openldap-stack-ha
            searchBase: dc=example,dc=org
            # existingSecret: openldaptest
            bindDN: cn=admin,dc=example,dc=org
            bindPWKey: LDAP_ADMIN_PASSWORD

        # web 
        phpldapadmin:
          enabled: true
          ingress:
            enabled: true
            annotations: {}
            path: /
            pathType: Prefix
            ## Ingress Host
            hosts:
            - phpldapadmin.example
          env:
            PHPLDAPADMIN_LDAP_HOSTS: openldap-openldap-stack-ha
```
