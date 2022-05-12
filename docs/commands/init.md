# dtm init

`dtm init` will download the required plugins according to the tool file in the config.

For config file, tool file, see the [config](../core-concepts/config.md) section of this documentation.

## Downloading Plugins

Before v0.5.0 (included,) `dtm` releases plugins to the GitHub releases page.

When you run `dtm init`, `dtm` will decide which plugins exist and which do not (based on the config file/tool file,) then download the plugins that are needed but don't exist locally from the GitHub release page.

After v0.5.0 (feature work in progress now), `dtm` will release plugins to an AWS S3 bucket. When running `dtm init`, it will download plugins from the AWS S3 bucket instead (through Cloudflare CDN.)

## Init Logic

- Based on the config file and the tool file, decide which plugins are required.
- If the plugin exists locally and the version is correct, do nothing.
- If the plugin is missing, download the plugin.
- After downloading, `dtm` also downloads the md5 value of that plugin and compares them. If matched, succeed.
