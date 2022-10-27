# Build

```shell
cd path/to/devstream
make clean
make build -j8 # multi-threaded build
```

This builds everything: `dtm` and all the plugins.

We also support the following build modes:

- Build `dtm` only: `make build-core`.
- Build a specific plugin: `make build-plugin.PLUGIN_NAME`. Example: `make build-plugin.argocd`.
- Build all plugins: `make build-plugins -j8` (multi-threaded build.)

See `make help` for more information.
