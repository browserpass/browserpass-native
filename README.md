# Browserpass - native host

## Build locally

Make sure you have Golang and [dep](https://github.com/golang/dep) installed.

The following `make` goals are available:

| Command                      | Description                                                  |
| ---------------------------- | ------------------------------------------------------------ |
| `make` or `make all`         | Fetch the latest dependencies, compile the app and run tests |
| `make deps`                  | Fetch the latest dependencies                                |
| `make browserpass`           | Compile the app for your OS                                  |
| `make browserpass-linux64`   | Compile the app for Linux 64-bit                             |
| `make browserpass-windows64` | Compile the app for Windows 64-bit                           |
| `make browserpass-darwin64`  | Compile the app for Mac OS X 64-bit                          |
| `make browserpass-openbsd64` | Compile the app for OpenBSD 64-bit                           |
| `make browserpass-freebsd64` | Compile the app for FreeBSD 64-bit                           |
| `make test`                  | Run tests                                                    |

## Build using Docker

First build the docker image using the following command in the project root:

```shell
docker build -t browserpass .
```

The entry point in the docker image is the `make` command. To run it:

```shell
docker run --rm -v "$(pwd)":/browserpass-native browserpass
```

Specify `make` goal(s) as the last parameter, for example:

```shell
docker run --rm -v "$(pwd)":/browserpass-native browserpass test
```

Refer to the list of available `make` goals above.
