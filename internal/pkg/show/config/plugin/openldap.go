package plugin

var OpenldapDefaultConfig = `tools:
# name of the tool
- name: openldap
  # id of the tool instance
  instanceID: default
  # format: name.instanceID; If specified, dtm will make sure the dependency is applied first before handling this tool.
  dependsOn: []
  # options for the plugin  
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
        replicaCount: 1
        service: 
          type: NodePort  
        env:
          LDAP_ORGANISATION: "DevStream Inc."
          LDAP_DOMAIN: "devstream.io"
        persistence:
          enabled: false
        adminPassword: Not@SecurePassw0rd
        configPassword: Not@SecurePassw0rd

        ltb-passwd:
          enabled : false

        phpldapadmin:
          enabled: true
          ingress:
            enabled: false`
