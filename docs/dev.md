# Development

## Local setup:
Only for Unix OS
```bash
$ make setup
```

## Run Tests
```bash
$ make test
```

## Lint
```bash
$ make lint
```

## Run demon

### Run test server:
```bash
$ make test-server
```

### Run demon
```bash
$ make run
```

### Stop demon
```bash
$ kill -2 $(cat /tmp/versiond/pid)
```

### Check test scripts logs:
```bash
$ cat /tmp/versiond/before_log
$ cat /tmp/versiond/main_log
$ cat /tmp/versiond/after_log
```
