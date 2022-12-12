# dtm-repo-scaffolding-java-springboot

This repo contains templates used by DevStream plugin "repo-scaffolding" (thereafter: the plugin).

This repo isn't intended to be used directly without DevStream. It should only be consumed by the plugin automatically.

The plugin (together with this repo of templates) can create a repo in GitHub and set up the project layout and initialize the reop with necessary files that are typical for a Spring web app. The followings can be created automatically:

- a Java Spring app example
- directory structure
- Dockerfile, with multistage build

## Where does this repo come from?

`dtm-repo-scaffolding-java-springboot` is synced from https://github.com/devstream-io/devstream/blob/main/staging/dtm-repo-scaffolding-java-springboot. 
Code changes are made in that location, merged into `devstream-io/devstream` and later synced here.
