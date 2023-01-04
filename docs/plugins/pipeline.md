# Pipeline Config Option

## imageRepo

Devstream will add a build stage in `CI` if you have configured this option. This stage will build a docker image and push this image to the configured image repository(like Dockerhub or self-host Harbor).

### Options

| Option  | Description                                       |
| ------ | ------------------------------------------ |
| user   | the user name of image repository          |
| url    | the address of image repository, the default value is Dockerhub address     |
| password| the password of `user`, this config option is used for auth when pushing image       |

## dingTalk

Devstream will add a notify stage in `CI` if you have configured this option. This stage will notify DingDing whether this `CI` is successful or failed.

### Options

| Option        | Description                                                         |
| ------------- | ------------------------------------------------------------ |
| name          | the name of dingding robot                                             |
| webhook       | the webhook of dingding robot                      |
| securityType  | the auth method of dingding robot, support SECRET/KEYWORD for now               |
| securityValue | if you set `securityType` to SECRET, you can set this value to secret value |

## sonarqube

Devstream will add a scanner stage in `CI` if you have configured this option. This stage will use [sonarqube](https://www.sonarqube.org/) to scan code.

### Options

| Option | Description                                                                                                                    |
| ------ | ----------------------------------------------------------------------------------------------------------------------- |
| name   | the name of sonarqube                                                                                                     |
| token  | the token of sonarqube, you can refer to this [doc](https://sonarqube.inria.fr/sonarqube/documentation/user-guide/user-token/) to get this token |
| url    | the url address of soanrqube                                                                                                   |

## language

If you config this option, Devstream will add the language's default option in the pipeline. For example, if you set `name` to `golang`, `CI` test stage will use `go test ./...` to run the test command.

### Options

| Option    | Description             |
| --------- | ---------------- |
| name      | name of programing language   |
| framework | the framework     |
| version   | the version of programing language  |
