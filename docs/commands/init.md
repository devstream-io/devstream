# dtm init

`dtm` releases plugins to an AWS S3 bucket. 

When running `dtm init`, it will download plugins from the AWS S3 bucket(through Cloudflare CDN.)

There are two ways to download plugins:

1. from config file: download plugins according to `tools` and `apps` in the config
2. from command line: download plugins according to the command line arguments

## Downloading Plugins from Config File

In this way, `dtm init` will download the required plugins according to `tools` and `apps` in the config.

**command:** `dtm init -f <config file/config dir>`.

You can put all the configuration in one file, or you can spread it out into multiple files in the same directory with `yaml` or `yaml` as a suffix.

For config file, tools and apps, see the [config](../core-concepts/config.md) section of this documentation.

## Downloading Plugins from Command Line

This can be used to pre-download the plugin to use `dtm` in an **offline environment**.

command: 

- download specify plugins. e.g. `dtm init --download-only --plugins="repo-scaffolding, github-actions" -d=.devstream/plugins`
- download all plugins. e.g. `dtm init --download-only --all -d=.devstream/plugins`

## Init Logic

- Based on the config file and the tool file or command line flags, decide which plugins are required.
- If the plugin exists locally and the version is correct, do nothing.
- If the plugin is missing, download the plugin.
- After downloading, `dtm` also downloads the md5 value of that plugin and compares them. If matched, succeed.

## Flags

| Short | Long            | Default                  | Description                                                                       |
|-------|-----------------|--------------------------|-----------------------------------------------------------------------------------|
| -f    | --config-file   | `"config.yaml"`          | The config file to use                                                            |
| -d    | --plugin-dir    | `"~/.devstream/plugins"` | The directory to store plugins                                                    |
|       | --download-only | `false`                  | Download plugins only, generally used for using dtm offline                       |
| -p    | --plugins       | `""`                     | Plugins to be downloaded, seperated with ",", should be used with --download-only |
| -a    | --all           | `false`                  | Download all plugins, should be used with --download-only                         |
|       | --os            | OS of this machine       | The OS to download plugins for.                                                   |
|       | --arch          | Arch of this machine     | The architecture to download plugins for.                                         |

