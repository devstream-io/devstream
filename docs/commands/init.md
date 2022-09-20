# dtm init

`dtm init` will download the required plugins according to the tool file in the config.

For config file, tool file, see the [config](../core-concepts/config.md) section of this documentation.

## Downloading Plugins

There are two ways to download plugins:
1. from config file: download plugins according to the tool file in the config
2. from command line: download plugins according to the command line arguments

### Downloading Plugins from Config File

command: `dtm init` or `dtm init -f <config file>`. The default config file path is `config.yaml`.

Before v0.5.0 (included,) `dtm` releases plugins to the GitHub releases page.

When you run `dtm init`, `dtm` will decide which plugins exist and which do not (based on the config file/tool file,) then download the plugins that are needed but don't exist locally from the GitHub release page.

After v0.5.0 (feature work in progress now), `dtm` will release plugins to an AWS S3 bucket. When running `dtm init`, it will download plugins from the AWS S3 bucket instead (through Cloudflare CDN.)

### Downloading Plugins from Command Line

This can be used to pre-download the plugin to use dtm in an offline environment.

command: 

1. download specify plugins. e.g. `dtm init --download-only --plugins=repo-scaffolding,githubactions-golang -d=.devstream/plugins`
2. download all plugins. e.g. `dtm init --download-only --all -d=.devstream/plugins`


## Flags

| Short | Long            | Default              | Description                                                                       |
|-------|-----------------|----------------------|-----------------------------------------------------------------------------------|
| -f    | --config-file   | `"config.yaml"`      | The config file to use                                                            |
| -d    | --plugin-dir    | `"~/.dtm/plugins"`   | The directory to store plugins                                                    |
|       | --download-only | `false`              | Download plugins only                                                             |
| -p    | --plugins       | `""`                 | Plugins to be downloaded, seperated with ",", should be used with --download-only |
| -a    | --all           | `false`              | Download all plugins, should be used with --download-only                         |
|       | --os            | OS of this machine   | The OS to download plugins for.                                                   |
|       | --arch          | Arch of this machine | The architecture to download plugins for.                                         |

_Note: You could use `dtm init --all --os=DestinyOS --arch=DestinyArch -d=PLUGIN_DIR` to download all plugins and use dtm offline._

## Init Logic

- Based on the config file and the tool file or command line flags, decide which plugins are required.
- If the plugin exists locally and the version is correct, do nothing.
- If the plugin is missing, download the plugin.
- After downloading, `dtm` also downloads the md5 value of that plugin and compares them. If matched, succeed.
