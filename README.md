# Browserpass - native messaging host

This is a host application for [browserpass](https://github.com/browserpass/browserpass-extension) browser extension providing it access to your password store. The communication is handled through [Native Messaging API](https://developer.chrome.com/extensions/nativeMessaging).

## Installation

### Install via package manager

The following operating systems provide a browserpass package that can be installed using a package manager:

-   TODO
-   TODO

Once the package is installed, refer to the section [Configure browsers](#configure-browsers).

If your OS is not listed above, proceed with the manual installation steps below.

### Install manually

Download [the latest Github release](https://github.com/browserpass/browserpass-native/releases), choose either the source code archive (if you want to compile the app yourself) or an archive for your operating system (it contains a pre-built binary).

All release files are signed with [this PGP key](https://keybase.io/maximbaz). To verify the signature of a given file, use `$ gpg --verify <file>.sig`.

It should report:

```
gpg: Signature made ...
gpg:                using RSA key 8053EB88879A68CB4873D32B011FDC52DA839335
gpg: Good signature from "Maxim Baz <...>"
gpg:                 aka ...
Primary key fingerprint: EB4F 9E5A 60D3 2232 BB52  150C 12C8 7A28 FEAC 6B20
     Subkey fingerprint: 8053 EB88 879A 68CB 4873  D32B 011F DC52 DA83 9335
```

Unpack the archive. If you decided to compile the application yourself, refer to the [Building the app](#building-the-app) section on how to do so. Once complete, continue with the steps below.

Finally install the app using `make install` (if you compiled it using `make browserpass`) or `make BIN=browserpass-XXXX install` (if you downloaded a release with pre-built binary).

### Configure browsers

The Makefile (which is also available in `/usr/lib/browserpass/`, if you installed via package manager) contains the following `make` goals to configure the browsers you use:

| Command                    | Description                                                                |
| -------------------------- | -------------------------------------------------------------------------- |
| `make hosts-chromium`      | Configure browserpass for Chromium browser, system-wide                    |
| `make hosts-chromium-user` | Configure browserpass for Chromium browser, for the current user only      |
| `make hosts-chrome`        | Configure browserpass for Google Chrome browser, system-wide               |
| `make hosts-chrome-user`   | Configure browserpass for Google Chrome browser, for the current user only |
| `make hosts-vivaldi`       | Configure browserpass for Vivaldi browser, system-wide                     |
| `make hosts-vivaldi-user`  | Configure browserpass for Vivaldi browser, for the current user only       |
| `make hosts-firefox`       | Configure browserpass for Firefox browser, system-wide                     |
| `make hosts-firefox-user`  | Configure browserpass for Firefox browser, for the current user only       |

In addition, Chromium-based browsers support the following `make` goals:

| Command                       | Description                                                                                  |
| ----------------------------- | -------------------------------------------------------------------------------------------- |
| `make policies-chromium`      | Automatically install browser extension for Chromium browser, system-wide                    |
| `make policies-chromium-user` | Automatically install browser extension for Chromium browser, for the current user only      |
| `make policies-chrome`        | Automatically install browser extension for Google Chrome browser, system-wide               |
| `make policies-chrome-user`   | Automatically install browser extension for Google Chrome browser, for the current user only |
| `make policies-vivaldi`       | Automatically install browser extension for Vivaldi browser, system-wide                     |
| `make policies-vivaldi-user`  | Automatically install browser extension for Vivaldi browser, for the current user only       |

## Building the app

### Build locally

Make sure you have the latest stable Go installed.

The following `make` goals are available:

| Command                      | Description                         |
| ---------------------------- | ----------------------------------- |
| `make` or `make all`         | Compile the app and run tests       |
| `make browserpass`           | Compile the app for your OS         |
| `make browserpass-linux64`   | Compile the app for Linux 64-bit    |
| `make browserpass-windows64` | Compile the app for Windows 64-bit  |
| `make browserpass-darwin64`  | Compile the app for Mac OS X 64-bit |
| `make browserpass-openbsd64` | Compile the app for OpenBSD 64-bit  |
| `make browserpass-freebsd64` | Compile the app for FreeBSD 64-bit  |
| `make test`                  | Run tests                           |

### Build using Docker

First build the docker image using the following command in the project root:

```shell
docker build -t browserpass-native .
```

The entry point in the docker image is the `make` command. To run it:

```shell
docker run --rm -v "$(pwd)":/src browserpass-native
```

Specify `make` goal(s) as the last parameter, for example:

```shell
docker run --rm -v "$(pwd)":/src browserpass-native test
```

Refer to the list of available `make` goals above.

## Updates

If you installed the app using a package manager for your OS, you will likely update it in the same way.

If you installed manually, repeat the steps in the [Install manually](#install-manually) section.
