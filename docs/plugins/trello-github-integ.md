# trello-github-integ Plugin

This plugin creates a new GitHub Actions workflow(trello-github-integration) and uploads it to your GitHub repo.

## Usage

This plugin depends on and can be used together with the `trello` plugin (see document [here](./trello.md)).

`trello-github-integ` plugin can also use `trello` plugin's outputs as input. See the example below:

```yaml
--8<-- "trello-github-integ.yaml"
```

Replace the following from the config above:

- `YOUR_GITHUB_USERNAME`
- `YOUR_REPO_NAME`
- `KANBAN_BOARD_NAME`

In the example above:

- We create a Trello board using `trello` plugin, and the board is marked to be used for repo YOUR_GITHUB_USERNAME/YOUR_REPO_NAME.
- `trello-github-integ` plugin depends on `trello` plugin, because we use `trello` plugin's outputs as the input for the `trello-github-integ` plugin.

Pay attention to the `${{ xxx }}` part in the example. `${{ TOOL_NAME.TOOL_INSTANCE_ID.outputs.var}}` is the syntax for using an output.
