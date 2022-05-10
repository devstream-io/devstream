##  `dtm verify`

The command `dtm verify` checks the following:

### 1 Config File

`dtm verify` first verifies if the config file can be loaded successfully.

If not, the following information might be printed out:

- if the config file doesn't exist, it reminds you if you forgot to specify the config file by using the "-f" parameter;
- if the config format isn't correct, it would print some error.

### 2 Plugins

`dtm verify` then checks if all required plugins (according to the config file) exist.

If not, it tries to give you a hint that maybe you forgot to run `dtm init` first.

### 3 State

`dtm verify` also tries to create a state manager that operates a backend. If something is wrong with the state, it generates an error telling you what exactly the error is.

### 4 Config / State / Resource

For definitions of _Config_, _State_, and _Resource_, see [Core Concepts](../core-concepts/core-concepts).

`dtm verify` tries to see if the _Config_ matches the _State_ and the _Resource_ or not. If not, it tells you what exactly is not the same, and what would happen if you run `dtm apply`.

If all the above checks are successful, `dtm verify` finishes with a success log "Verify succeeded."
