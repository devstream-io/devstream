#  dtm destroy

`dtm destroy` acts like `dtm apply -f an_empty_config.yaml`.

The purpose of `destroy` is that in case you accidentally deleted your config file during testing, there would still be a way to destroy everything that is defined in the _State_ so that you can have a clean slate.
