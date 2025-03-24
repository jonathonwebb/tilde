# ~ (tilde)

## Setup

```bash
# install dependencies
pnpm install

# compile project source
./bin/build

# run app
./bin/tilde
```

## Usage

### `./bin/tilde`

```text
usage: tilde [options] [command]

options:
  -e, --env <env>      application environment (default: "production", env: MDW_ENV)
  -l, --level <level>  log level (default: "info", env: MDW_LEVEL)
  -V, --version        show version info and exit
  -h, --help           display help for command

commands:
  serve [options]
  help [command]       display help for command
```

### `./bin/tilde serve`

```text
usage: tilde serve [options]

options:
  -a, --host <host>    listener host (default: "localhost", env: MDW_HOST)
  -p, --port <number>  listener port (default: ephemeral, env: MDW_PORT)
  -h, --help           display help for command
```

## Development

```bash
./bin/build    # compile source
./bin/check    # lint and type check source
./bin/coverage # run tests with coverage
./bin/fix      # lint with autofix
./bin/test     # run tests
./bin/watch    # run app and watch for changes
```
