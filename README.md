# miniproxy

miniproxy is a small proxy server written in Go, primarily for use in local development
where you wish to route requests on a single port to one of a number of other ports.

Routing is done by request path.

miniproxy will take care of starting your services and restarting them if they crash.

In future, I would like to include a file watcher as part of miniproxy, but there don't
seem to be any good recursive, cross-platform file-watching libraries available yet,
and I don't have the inclination to create one. 

## Installation

To install miniproxy, run:

```
go install github.com/isaacharrisholt/miniproxy@latest
```

Note: miniproxy requires Go version 1.21 or later.

## Usage

To start using miniproxy, simply create a `miniproxy.json` file in the directory you'd
like to use miniproxy from.

Here's a full example of a `miniproxy.json` file:

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
  "default": "go",
  "routes": {
    "/api/foo/*": "go",
    "/api/bar/*": "go",
    "/api/baz": "rust"
  }
}
```

Once you have a `miniproxy.json`, simply run `miniproxy` in the same directory.
In the future, miniproxy will look for a `miniproxy.json` in parent directories,
but this hasn't been implemented yet.

This will start miniproxy on port 3000 (which is also the default) and start two other
services - `go` and `rust` - in their respective directories.

You may use glob pattern matching in `routes`, and it's also worth mentioning that
`targets.<target>.service` is optional, as is `targets.<target>.service.workDir`.
The latter will default to the current directory.

If the `default` key is specified, any request that does not match a route in `routes`
will be routed to the default service. This is equivalent to having `"*"` as your
last route.

You may also specify the log level (`debug`, `info`, `error`) with the `LOG_LEVEL`
environment variable.

## Contributing

Contributions are welcome! However, if you'd like a new feature, please submit an
issue first - I only have a limited capacity for maintenance.

## License

miniproxy is MIT licensed, so do what you like!
