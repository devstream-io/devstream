# dtm-gitlab-share-library

This repo contains share library used by DevStream plugin "gitlab-ci" (thereafter: the plugin). It currently has the following functions：
- send notification to dingtalk when pipeline result is success or failed.
- build docker image by current project, and push this image to image repo.

## General config variables

|  field   | description  | default_value  |
|  ----  | ----  | ----  |
|  CI_REGISTRY_USER  | image repo owner  |   |
|  CI_REGISTRY_PASSWORD  | image repo owner's password  |  |

## Where does this repo come from?

`dtm-gitlab-share-library` is synced from https://github.com/devstream-io/devstream/blob/main/staging/dtm-gitlab-share-library.
Code changes are made in that location, merged into `devstream-io/devstream` and later synced here.
