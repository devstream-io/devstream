# dtm apply

When _applying_ a config file using `dtm`, here's what happens:

## 1 For Each _Tool_ Defined in the _Config_

We compare the _Tool_, its _State_, and the _Resoruce_ it has created before (if the state exists).

We generate a plan of changes according to the comparison result:
- If the _Tool_ isn't in the _State_, the `Create` interface will be called.
- If the _Tool_ is in the _State_, but the _Config_ is different than the _State_ (meaning users probably updated the config after the last `apply`,) the `Update` interface will be called.
- If the _Tool_ is in the _State_, and the _Config_ is the same as the _State_, we try to read the _Resource_.
    - If the _Resource_ doesn't exist, the `Create` interface will be called. It probably suggests that the _Resource_ got deleted manually after the last successful `apply`.
    - If the _Resource_ does exist but drifted from the _State_ (meaning somebody modified it), the `Update` interface will be called.
    - Last but not least, nothing would happen if the _Resource_ is exactly the same as the _State_.

## 2 For Each _State_ That Doesn't Have a _Tool_ in the _Config_

We generate a "Delete" change to delete the _Resource_. Since there isn't a _Tool_ in the config but there is a _State_, it means maybe the _Resource_ had been created previously then the user removed the _Tool_ from the _Config_, which means the user doesn't want the _Resource_ any more.
