<p align="center"><img src="https://github.com/browserpass/browserpass-extension/raw/master/images/logotype-horizontal.png"></p>

# Browserpass - native messaging host

This is a host application for [browserpass](https://github.com/browserpass/browserpass-extension) browser extension providing it access to your password store. The communication is handled through [Native Messaging API](https://developer.chrome.com/extensions/nativeMessaging).

## Table of Contents

-   [Installation](#installation)
    -   [Install via package manager](#install-via-package-manager)
    -   [Install manually](#install-manually)
        -   [Install on Nix / NixOS](#install-on-nix--nixos)
        -   [Install on Windows](#install-on-windows)
        -   [Install on Windows through WSL](#install-on-windows-through-wsl)
    -   [Configure browsers](#configure-browsers)
-   [Building the app](#building-the-app)
    -   [Build locally](#build-locally)
    -   [Build using Docker](#build-using-docker)
-   [Updates](#updates)
-   [FAQ](#faq)
    -   [Error: Unable to fetch and parse login fields](#error-unable-to-fetch-and-parse-login-fields)
-   [Uninstallation](#uninstallation)
-   [Contributing](#contributing)

## Installation

### Install via package manager

The following operating systems provide a browserpass package that can be installed using a package manager:

-   Arch Linux: [browserpass](https://www.archlinux.org/packages/community/x86_64/browserpass/)
-   Gentoo Linux: [browserpass](https://packages.gentoo.org/packages/www-plugins/browserpass)
-   Debian sid: [webext-browserpass](https://packages.debian.org/sid/webext-browserpass) (note: users report ([#126](https://github.com/browserpass/browserpass-native/issues/126)) that it's not suitable for any browser besides Chromium and Firefox, use manual installation if you plan to use other browsers).
-   openSUSE Tumbleweed: [browserpass-native](https://software.opensuse.org/package/browserpass-native)
-   NixOS: [browserpass](https://github.com/NixOS/nixpkgs/blob/master/pkgs/tools/security/browserpass/default.nix) - also read [Install on Nix / NixOS](#install-on-nix--nixos)
-   macOS: [browserpass](https://github.com/Amar1729/homebrew-formulae/blob/master/browserpass.rb) in a user-contributed tap [amar1729/formulae](https://github.com/amar1729/homebrew-formulae)

Once the package is installed, **refer to the section [Configure browsers](#configure-browsers)**.

If your OS is not listed above, proceed with the manual installation steps below.

### Install manually

Download [the latest Github release](https://github.com/browserpass/browserpass-native/releases), choose either the source code archive (if you want to compile the app yourself) or an archive for your operating system (it contains a pre-built binary).

All release files are signed with a PGP key that is available on [maximbaz.com](https://maximbaz.com/), [keybase.io](https://keybase.io/maximbaz) and various OpenPGP key servers. First, import the public key using any of these commands:

```
$ curl https://maximbaz.com/pgp_keys.asc | gpg --import
$ curl https://keybase.io/maximbaz/pgp_keys.asc | gpg --import
$ gpg --recv-keys 56C3E775E72B0C8B1C0C1BD0B5DB77409B11B601
```

To verify the signature of a given file, use `$ gpg --verify <file>.asc`.

It should report:

```
gpg: Signature made ...
gpg:                using EDDSA key 04D7A219B0ABE4C2B62A5E654A2B758631E1FD91
gpg: Good signature from "Maxim Baz <...>"
gpg:                 aka ...
Primary key fingerprint: 56C3 E775 E72B 0C8B 1C0C  1BD0 B5DB 7740 9B11 B601
     Subkey fingerprint: 04D7 A219 B0AB E4C2 B62A  5E65 4A2B 7586 31E1 FD91
```

Unpack the archive. If you decided to compile the application yourself, refer to the [Building the app](#building-the-app) section on how to do so. Once complete, continue with the steps below.

If you are on macOS, first install the necessary tools: `brew install coreutils gnu-sed`.

If you are on FreeBSD, first install the GNU tools: `pkg install coreutils gmake gsed'.` Use `gmake` in place of `make` below.

If you downloaded a release archive with pre-compiled binary, follow these steps to install the app:

```
# IMPORTANT: replace XXXX with OS name depending on the archive you downloaded, e.g. "linux64"
make BIN=browserpass-XXXX configure      # Configure the hosts json files
sudo make BIN=browserpass-XXXX install   # Install the app
```

In addition, both `configure` and `install` targets respect `PREFIX`, `DESTDIR` parameters if you want to customize the install location (e.g. to install to a `$HOME` dir to avoid using `sudo`).

For example, if you are on macOS or FreeBSD, you probably want to install Browserpass in `/usr/local/bin`, therefore you need to run:

```
make BIN=browserpass-darwin64 PREFIX=/usr/local configure      # Configure the hosts json files
sudo make BIN=browserpass-darwin64 PREFIX=/usr/local install   # Install the app
```

If you compiled the app yourself, you can omit `BIN` parameter:

```
make configure      # Configure the hosts json files
sudo make install   # Install the app
```

Finally proceed to the [Configure browsers](#configure-browsers) section.

#### Install on Nix / NixOS

For a declarative NixOS installation, update your channel with `sudo nix-channel --update`, use the following to your `/etc/nixos/configuration.nix` and rebuild your system:

```nix
{ pkgs, ... }: {
  programs.browserpass.enable = true;
  environment.systemPackages = with pkgs; [
    # All of these browsers will work
    chromium firefox google-chrome vivaldi
    # firefox*-bin versions do *not* work with this. If you require such Firefox versions, use the stateful setup described below.
  ];
}
```

For a stateful Nix setup, update your channel, install Browserpass and link the necessary files with the Makefile (see [Configure browsers](#configure-browsers) section), but pass `DESTDIR=~/.nix-profile`:

```bash
$ nix-channel --update
$ nix-env -iA nixpkgs.browserpass # Or nix-env -iA nixos.browserpass on NixOS
$ DESTDIR=~/.nix-profile make -f ~/.nix-profile/lib/browserpass/Makefile <desired make goal>
```

#### Install on Windows

Download [the latest Github release](https://github.com/browserpass/browserpass-native/releases/latest) for `windows64`.

Run the installer, it will install all the necessary files in `C:\Program Files\Browserpass` and it will also [configure browsers](#configure-browsers).

#### Install on Windows through WSL

If you want to use WSL instead

1. Follow the [installation](#installation) steps for the WSL distribution you are using. There is no need to configure the browser as your browser does not run in WSL.
2. Follow the then [installation](#install-on-windows) steps for installing on Windows.
3. Create `C:\Program Files\Browserpass\browserpass-wsl.bat` with the following contents:

```
@echo off
bash -c "/usr/bin/browserpass-linux64 2>/dev/null"
```

4. Edit the hosts json files (in our example `C:\Program Files\Browserpass\browser-files\*-host.json`) and replace `browserpass-windows64.exe` with `browserpass-wsl.bat` you've just created.

Remember to check [Hints for configuring gpg](#error-unable-to-fetch-and-parse-login-fields) on how to configure pinentry to unlock your PGP key.

### Configure browsers

The following operating systems provide packages for certain browsers that can be installed using a package manager:

-   Arch Linux: [browserpass-chromium](https://www.archlinux.org/packages/community/any/browserpass-chromium/) and [browserpass-firefox](https://www.archlinux.org/packages/community/any/browserpass-firefox/)
    -   AUR: [browserpass-chrome](https://aur.archlinux.org/packages/browserpass-chrome/)

If you installed a distro package above, you are done!

If something went wrong, if there's no package for your OS and/or a browser of your choice, or for whatever reason you just don't want to use them, proceed with the steps below.

First, enter the directory with installed Browserpass, by default it is `/usr/lib/browserpass/`, but if you used `PREFIX` or `DESTDIR` when running `make install`, it might be different for you. For example, on macOS the directory is likely to be `/usr/local/lib/browserpass/`.

See below the list of available `make` goals to configure various browsers. Use `gmake` on FreeBSD in place of `make`.

**It is recommended to use `*-user` make goals**, as more people had luck with them. But if they don't work as expected, try other available goals.

If you provided `PREFIX` and/or `DESTDIR` while running `make install`, remember that you must provide the same parameters, for example `make PREFIX=/usr/local hosts-chromium-user`:

| Command                     | Description                                                                          |
| --------------------------- | ------------------------------------------------------------------------------------ |
| `make hosts-chromium-user`  | Configure browserpass for Chromium browser, for the current user only                |
| `make hosts-firefox-user`   | Configure browserpass for Firefox browser, for the current user only                 |
| `make hosts-librewolf-user` | Configure browserpass for Librewolf browser, for the current user only               |
| `make hosts-chrome-user`    | Configure browserpass for Google Chrome or Opera browsers, for the current user only |
| `make hosts-edge-user`      | Configure browserpass for Microsoft Edge browser, for the current user only          |
| `make hosts-brave-user`     | Configure browserpass for Brave browser, for the current user only                   |
| `make hosts-iridium-user`   | Configure browserpass for Iridium browser, for the current user only                 |
| `make hosts-vivaldi-user`   | Configure browserpass for Vivaldi browser, for the current user only                 |
| `make hosts-yandex-user`    | Configure browserpass for Yandex browser, for the current user only                  |
| `make hosts-slimjet-user`   | Configure browserpass for Slimjet browser, for the current user only                 |
| `sudo make hosts-chromium`  | Configure browserpass for Chromium browser, system-wide                              |
| `sudo make hosts-firefox`   | Configure browserpass for Firefox browser, system-wide                               |
| `sudo make hosts-librewolf` | Configure browserpass for Librewolf browser, system-wide                             |
| `sudo make hosts-chrome`    | Configure browserpass for Google Chrome or Opera browsers, system-wide               |
| `sudo make hosts-edge`      | Configure browserpass for Microsoft Edge browser, system-wide                        |
| `sudo make hosts-brave`     | Configure browserpass for Brave browser, system-wide                                 |
| `sudo make hosts-iridium`   | Configure browserpass for Iridium browser, system-wide                               |
| `sudo make hosts-vivaldi`   | Configure browserpass for Vivaldi browser, system-wide                               |
| `sudo make hosts-yandex`    | Configure browserpass for Yandex browser, system-wide                                |
| `sudo make hosts-slimjet`   | Configure browserpass for Slimjet browser, system-wide                               |

In addition, Chromium-based browsers support the following `make` goals:

| Command                       | Description                                                                                                  |
| ----------------------------- | ------------------------------------------------------------------------------------------------------------ |
| `make policies-chromium-user` | Automatically install browser extension from Web Store for Chromium browser, for the current user only       |
| `make policies-chrome-user`   | Automatically install browser extension from Web Store for Google Chrome browser, for the current user only  |
| `make policies-edge-user`     | Automatically install browser extension from Web Store for Microsoft Edge browser, for the current user only |
| `make policies-brave-user`    | Automatically install browser extension from Web Store for Brave browser, for the current user only          |
| `make policies-iridium-user`  | Automatically install browser extension from Web Store for Iridium browser, for the current user only        |
| `make policies-slimjet-user`  | Automatically install browser extension from Web Store for Slimjet browser, for the current user only        |
| `make policies-vivaldi-user`  | Automatically install browser extension from Web Store for Vivaldi browser, for the current user only        |
| `make policies-yandex-user`   | Automatically install browser extension from Web Store for Yandex browser, for the current user only         |
| `sudo make policies-chromium` | Automatically install browser extension from Web Store for Chromium browser, system-wide                     |
| `sudo make policies-chrome`   | Automatically install browser extension from Web Store for Google Chrome browser, system-wide                |
| `sudo make policies-edge`     | Automatically install browser extension from Web Store for Microsoft Edge browser, system-wide               |
| `sudo make policies-brave`    | Automatically install browser extension from Web Store for Brave browser, system-wide                        |
| `sudo make policies-iridium`  | Automatically install browser extension from Web Store for Iridium browser, system-wide                      |
| `sudo make policies-slimjet`  | Automatically install browser extension from Web Store for Slimjet browser, system-wide                      |
| `sudo make policies-vivaldi`  | Automatically install browser extension from Web Store for Vivaldi browser, system-wide                      |
| `sudo make policies-yandex`   | Automatically install browser extension from Web Store for Yandex browser, system-wide                       |

## Building the app

### Build locally

Make sure you have the latest stable Go installed.

The following `make` goals are available (check Makefile for more details):

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

## FAQ

### Error: Unable to fetch and parse login fields

If you can see passwords, but unable to fill forms or copy credentials, you likely have issues with your `gpg` setup.

First things first, make sure that `gpg` and some GUI `pinentry` are installed.

-   on macOS many people succeeded with `pinentry-mac`
-   on Linux [users report](https://github.com/browserpass/browserpass-extension/issues/155) that `pinentry-gnome3` does not work well with GNOME 3 and Firefox, use e.g. `pinentry-gtk-2`
-   on Windows WSL people succeded with [pinentry-wsl-ps1](https://github.com/diablodale/pinentry-wsl-ps1)

`pinentry` is the application that asks you your password to unlock PGP key when you for example use `pass`.

The selected `pinentry` **must have GUI**, console-based (like `pinentry-tty` or `pinentry-curses`) **are not supported** (unless you know what you are doing).

Ensure that `gpg-agent` process is actually running, if not you need to investigate how to enable it.

Finally configure a GUI pinentry program in `~/.gnupg/gpg-agent.conf`:

```
pinentry-program /full/path/to/pinentry
```

You will need to restart `gpg-agent` using: `$ gpgconf --kill gpg-agent`

If Browserpass is unable to locate the proper `gpg` binary, try configuring a full path to your `gpg` in the browser extension settings or in `.browserpass.json` file in the root of your password store:

```json
{
    "gpgPath": "/full/path/to/gpg"
}
```

## Uninstallation

If you installed Browserpass via a package manager, uninstalling the package will hopefully do all the necessary cleanup.

For manual installations, uninstallation procedure basically consists of following the steps you ran in Makefile in reverse.

## Contributing

1. Fork [the repo](https://github.com/browserpass/browserpass-extension)
2. Create your feature branch
    - `git checkout -b my-new-feature`
3. Commit your changes
    - `git commit -am 'Add some feature'`
4. Push the branch
    - `git push origin my-new-feature`
5. Create a new pull request
