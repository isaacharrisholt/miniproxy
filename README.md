# tinyproxy

tinyproxy is a small proxy server written in Go, primarily for use in local development
where you wish to route requests on a single port to one of a number of other ports.

Routing is done by request path.

tinyproxy will take care of starting your services and restarting them if they crash.

In future, I would like to include a file watcher as part of tinyproxy, but there don't
seem to be any good recursive, cross-platform file-watching libraries available yet,
and I don't have the inclination to create one. 

## Installation

To install tinyproxy, run:

```
go install github.com/isaacharrisholt/tinyproxy@latest
```

Note: tinyproxy requires Go version 1.21 or later.

## Usage

To start using tinyproxy, simply create a `tinyproxy.json` file in the directory you'd
like to use tinyproxy from.

Here's a full example of a `tinyproxy.json` file:

```json
{
  "port": 3000,
  "targets": {
    "go": {
      "port": 3001,
      "service": {
        "command": ["go", "run", "."],
        "workDir": "./go"
      }
    },
    "rust": {
      "port": 3002,
      "service": {
        "command": ["cargo", "watch", "-s", "cargo run --bin local"],
        "workDir": "./rust"
      }
    }
  },
  "routes": {
    "/api/foo/*": "go",
    "/api/bar/*": "go",
    "/api/baz": "rust"
  }
}
```

Once you have a `tinyproxy.json`, simply run `tinyproxy` in the same directory.
In the future, tinyproxy will look for a `tinyproxy.json` in parent directories,
but this hasn't been implemented yet.

This will start tinyproxy on port 3000 (which is also the default) and start two other
services - `go` and `rust` - in their respective directories.

You may use glob pattern matching in `routes`, and it's also worth mentioning that
`targets.<target>.service` is optional, as is `targets.<target>.service.workDir`.
The latter will default to the current directory.

You may also specify the log level (`debug`, `info`, `error`) with the `LOG_LEVEL`
environment variable.

## Contributing

Contributions are welcome! However, if you'd like a new feature, please submit an
issue first - I only have a limited capacity for maintenance.

## License

tinyproxy is MIT licensed, so do what you like!