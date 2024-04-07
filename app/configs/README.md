# versiond config

## Example

```json
{
  "default_version": "v1",

  "cache": {
    "dir": "/etc/versiond"
  },

  "version": {
    "source": {
      "address": "http://localhost:8888/api/verson",
      "timeout": "3s"
    },

    "monitor": {
      "period": "300s"
    }
  },

  "on_change_cmds": {
    "before": "/opt/foo/before.sh",
    "main": "/opt/bar/main.sh",
    "after": "/opt/baz/after.sh"
  }
}
```

## Sections

### default_version
Default version that will be used at demon start.

### cache
Cache settings, for now - only directory (`dir`) where current version will be stored

### version
Those settings are about version monitoring:

* source:
  * address - http-address;
  * timeout - http-request timeout.
* monitor:
  * period - compare versions period.

Expected response:
```json
{
  "newest_version": "v2"
}
```
key `newest_version` with string value that will be compared to current version.

### on_change_cmds
Commands executed after version change will happen:

* before - path of script that will be executed before main;
* main - path of script that will be executed for new version to take place;
* after - path of script that will be executed after new version is handled.

For each command will be provided environment variables CURRENT_VERSION and NEW_VERSION with corresponding values.

### config
TODO add impl of config updates
