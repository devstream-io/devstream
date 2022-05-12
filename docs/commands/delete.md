#  dtm delete

## 1 Normal (Non-Force) Delete

When _deleting_ using `dtm`, here's what happens:

- Read the _Config_
- For each _Tool_ defined in the _Config_, if there is a corresponding _State_, the `Delete` interface will be called.

_Note: the `Delete` change will be executed only if the _State_ exists._

## 2 Force Delete

When _deleting_ using `dtm delete --force`, here's what happens:

- Read the _Config_
- The `Delete` interface will be called for each _Tool_ defined in the _Config_.

_Note: the difference between "force delete" and "normal delete" is that in force mode, no matter the _State_ exists or not, `dtm` will try to trigger the `Delete` interface. The purpose is for corner cases where the state gets corrupted or even lost during testing (I hope this only happens in the development environment), there is still a way to clean up the _Tools_._
