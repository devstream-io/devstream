## 1 `openldap` Plugin

This plugin installs [OpenLDAP](https://www.openldap.org/) in an existing Kubernetes cluster using the Helm chart. Please at least make sure your Kubernetes's version is greater than 1.18.

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
          LDAP_ORGANISATION: "DevStream Inc."
          LDAP_DOMAIN: "devstream.org"
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
            searchBase: dc=devstream,dc=org
            # existingSecret: openldaptest
            bindDN: cn=admin,dc=devstream,dc=org
            bindPWKey: LDAP_ADMIN_PASSWORD

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

## Description of key fields in values_yaml
- `service.type`: The default value is `ClusterIP`, if you have services outside the Kubernetes cluster that require ldap integration, the value preferably be set to `NodePort`, so that services outside the Kubernetes cluster can access the ldap service via `ldap://ip:389` instead of `ldap://openldap.openldap-openldap-stack-ha`
- `env`: List of key value pairs as env variables to be sent to the docker image. See https://github.com/osixia/docker-openldap for available ones. Please change the value of `LDAP_DOMAIN`. Of course even if you use **devstream.org** it will be fine, except that the search base will be `dc=devstream,dc=org`. 
- `persistence.storageClass`: Please use your own `storage class`, or use the `storage class` provided by the Kubernetes cluster hosted in the public cloud directly. The above example uses the NFS-based `storage class` provided by AliCloud ACK
- `adminPassword`: Use your own custom password
- `configPassword`: Use your own custom password
- `ltb-passwd.ingress`: Ingress of the Ltb-Passwd service by which you can modify your password. Please change **ssl-ldap2.example** to your own domain name
- `ltb-passwd.ldap`: Ldap configuration for the Ltb-Passwd service. If you change the `env.LDAP_DOMAIN`, don't forget to change the values of `ltb-passwd.ldap.searchBase` and `ltb-passwd.ldap.bindDN`
- `phpldapadmin.ingress`: Ingress of Phpldapadmin service by which you can manage your ldap service. Please change **phpldapadmin.example** to your own domain name

## Post-installation operations

Once the installation is complete, you can manage ldap service through **phpldapadmin**. Access the service by visiting the domain name (e.g. **phpldapadmin.example**) in the `phpldapadmin.ingress` configuration section of the above example. If you have not changed the default values in the above example, its account will be **cn=admin,dc=devstream,dc=org** and password will be **Not@SecurePassw0rd**.

If you're familiar with OpenLDAP, then you don't need to continue reading the tutorial below, you can just go ahead and integrate ldap for your service.

### Importing your data

The following is a sample file, if you have changed the above configuration, remember to replace `dc=devstream,dc=org` with your own.

```
dn: cn=admin,dc=devstream,dc=org
cn: admin
objectclass: organizationalRole

dn: ou=Group,dc=devstream,dc=org
cn: Group
objectclass: organizationalRole
ou: Group

# confluence organizationalUnit
dn: ou=confluence,ou=Group,dc=devstream,dc=org
objectclass: organizationalUnit
objectclass: top
ou: confluence

# confluence administrators group
dn: cn=confluence-administrators,ou=confluence,ou=Group,dc=devstream,dc=org
cn: confluence-administrators
description:: d2lraeeuoeeQhue7hA==
objectclass: groupOfUniqueNames
uniquemember: uid=example,ou=People,dc=devstream,dc=org

# confluence users group
dn: cn=confluence-users,ou=confluence,ou=Group,dc=devstream,dc=org
cn: confluence-users
description:: d2lraeaZrumAmueUqOaItw==
objectclass: groupOfUniqueNames
uniquemember: uid=example,ou=People,dc=devstream,dc=org

# jira organizationalUnit
dn: ou=jira,ou=Group,dc=devstream,dc=org
objectclass: organizationalUnit
objectclass: top
ou: jira

# jira administrators Group
dn: cn=jira-administrators,ou=jira,ou=Group,dc=devstream,dc=org
cn: jira-administrators
description:: amlyYeeuoeeQhue7hA==
objectclass: groupOfUniqueNames
uniquemember: uid=example,ou=People,dc=devstream,dc=org

# jira users group
dn: cn=jira-software-users,ou=jira,ou=Group,dc=devstream,dc=org
cn: jira-software-users
description:: amlyYeeuoeeQhue7hA==
objectclass: groupOfUniqueNames
uniquemember: uid=example,ou=People,dc=devstream,dc=org

dn: ou=People,dc=devstream,dc=org
objectclass: organizationalUnit
ou: People

# People for example
dn: uid=example,ou=People,dc=devstream,dc=org
cn: example
gidnumber: 500
givenname: example
homedirectory: /home/example
loginshell: /bin/sh
mail: example@devstream.org
objectclass: inetOrgPerson
objectclass: posixAccount
objectclass: top
sn: example
uid: example
uidnumber: 1007
userpassword: example@123456
```

Login your `phpldapadmin` service and import the sample configuration above.After importing the data successfully, the result is as follows.

![](../images/openldap-example.png)

### Verify the ldap service

Log in to the container where the ldap service is located, and then use the `ldapsearch` command to query the user(`uid=example,ou=people,dc=devstream,dc=org`) created above

```bash
root@openldap-openldap-stack-ha-0:/# ldapsearch -x -H ldap://127.0.0.1:389 -b uid=example,ou=people,dc=devstream,dc=org -D "cn=admin,dc=devstream,dc=org" -w Not@SecurePassw0rd

# extended LDIF
#
# LDAPv3
# base <uid=example,ou=people,dc=devstream,dc=org> with scope subtree
# filter: (objectclass=*)
# requesting: ALL
#

# example, People, devstream.org
dn: uid=example,ou=People,dc=devstream,dc=org
cn: example
gidNumber: 500
givenName: example
homeDirectory: /home/example
loginShell: /bin/sh
mail: example@devstream.org
objectClass: inetOrgPerson
objectClass: posixAccount
objectClass: top
sn: example
uid: example
uidNumber: 1007
userPassword:: ZXhhbXBsZUAxMjM0NTY=

# search result
search: 2
result: 0 Success

# numResponses: 2
# numEntries: 1
```


Then you can create users in the **People** group, assign them to different groups, and integrate with ldap-enabled services, and you can implement unified authentication based on OpenLDAP.