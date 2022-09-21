# dtm show

`dtm show` shows plugins' configuration templates or status.

## 1 show config template of plugin/toolchain

`dtm show config` shows configuration templates of plugins or sample toolchains.

Flags:

| Short | Long       | Default | Description                                                                |
|-------|------------|---------|----------------------------------------------------------------------------|
| -p    | --plugin   | `""`    | plugin name                                                                |
| -t    | --template | `""`    | name of template tool chain, currently supports "quickstart" and "gitopts" |

## 2 show plugin status

`dtm show status` shows status of plugin instances.

Flags:

| Short | Long          | Default                  | Description                              |
|-------|---------------|--------------------------|------------------------------------------|
| -p    | --plugin      | `""`                     | plugin name                              |
| -i    | --id          | `""`                     | plugin instance id                       |
| -a    | --all         | `false`                  | show all instances of all plugins status |
| -d    | --plugin-dir  | `"~/.devstream/plugins"` | plugins directory                        |
| -f    | --config-file | `"config.yaml"`          | config file                              |


_Note: If `-p` and `-i` are both empty, `dtm` will show all instances of all plugins' status._


