BIN ?= browserpass
VERSION ?= $(shell cat .version)

PREFIX ?= /usr
BIN_DIR = $(DESTDIR)$(PREFIX)/bin
LIB_DIR = $(DESTDIR)$(PREFIX)/lib
SHARE_DIR = $(DESTDIR)$(PREFIX)/share

GO_GCFLAGS := "all=-trimpath=${PWD}"
GO_ASMFLAGS := "all=-trimpath=${PWD}"
GO_LDFLAGS := "-extldflags ${LDFLAGS}"

APP_ID := com.github.browserpass.native
OS := $(shell uname -s)

# GNU tools
SED := $(shell which gsed 2>/dev/null || which sed 2>/dev/null)
INSTALL := $(shell which ginstall 2>/dev/null || which install 2>/dev/null)

#######################
# For local development

.PHONY: all
all: browserpass test

browserpass: *.go **/*.go
	go build -ldflags $(GO_LDFLAGS) -gcflags $(GO_GCFLAGS) -asmflags $(GO_ASMFLAGS) -o $@

browserpass-linux64: *.go **/*.go
	env GOOS=linux GOARCH=amd64 go build -ldflags $(GO_LDFLAGS) -gcflags $(GO_GCFLAGS) -asmflags $(GO_ASMFLAGS) -o $@

browserpass-darwin64: *.go **/*.go
	env GOOS=darwin GOARCH=amd64 go build -ldflags $(GO_LDFLAGS) -gcflags $(GO_GCFLAGS) -asmflags $(GO_ASMFLAGS) -o $@

browserpass-openbsd64: *.go **/*.go
	env GOOS=openbsd GOARCH=amd64 go build -ldflags $(GO_LDFLAGS) -gcflags $(GO_GCFLAGS) -asmflags $(GO_ASMFLAGS) -o $@

browserpass-freebsd64: *.go **/*.go
	env GOOS=freebsd GOARCH=amd64 go build -ldflags $(GO_LDFLAGS) -gcflags $(GO_GCFLAGS) -asmflags $(GO_ASMFLAGS) -o $@

browserpass-windows64: *.go **/*.go
	env GOOS=windows GOARCH=amd64 go build -ldflags $(GO_LDFLAGS) -gcflags $(GO_GCFLAGS) -asmflags $(GO_ASMFLAGS) -o $@.exe

.PHONY: test
test:
	go test ./...

#######################
# For official releases

.PHONY: clean
clean:
	rm -f browserpass browserpass-*
	rm -rf dist

.PHONY: dist
dist: clean browserpass-linux64 browserpass-darwin64 browserpass-openbsd64 browserpass-freebsd64 browserpass-windows64
	mkdir -p dist

	git archive -o dist/browserpass-native-$(VERSION).tar.gz --format tar.gz --prefix=browserpass-native-$(VERSION)/ $(VERSION)

	$(eval TMP := $(shell mktemp -d))

	for os in linux64 darwin64 openbsd64 freebsd64 windows64; do \
	    mkdir $(TMP)/browserpass-"$$os"-$(VERSION); \
	    cp -a browserpass-"$$os"* browser-files Makefile README.md LICENSE $(TMP)/browserpass-"$$os"-$(VERSION); \
	    if [ "$$os" = "windows64" ]; then \
	        (cd $(TMP) && zip -r ${CURDIR}/dist/browserpass-"$$os"-$(VERSION).zip browserpass-"$$os"-$(VERSION)); \
	    else \
	        (cd $(TMP) && tar -cvzf ${CURDIR}/dist/browserpass-"$$os"-$(VERSION).tar.gz browserpass-"$$os"-$(VERSION)); \
	    fi \
	done

	rm -rf $(TMP)

	for file in dist/*; do \
	    gpg --detach-sign --armor "$$file"; \
	done

	rm -f dist/browserpass-native-$(VERSION).tar.gz

#######################
# For user installation

.PHONY: configure
configure:
	$(SED) -i 's|"path": ".*"|"path": "'"$(BIN_DIR)/$(BIN)"'"|' browser-files/chromium-host.json
	$(SED) -i 's|"path": ".*"|"path": "'"$(BIN_DIR)/$(BIN)"'"|' browser-files/firefox-host.json

.PHONY: install
install:
	$(INSTALL) -Dm755 -t "$(BIN_DIR)/" $(BIN)
	$(INSTALL) -Dm644 -t "$(LIB_DIR)/browserpass/" Makefile
	$(INSTALL) -Dm644 -t "$(SHARE_DIR)/licenses/browserpass/" LICENSE
	$(INSTALL) -Dm644 -t "$(SHARE_DIR)/doc/browserpass/" README.md

	$(INSTALL) -Dm644 browser-files/chromium-host.json   "$(LIB_DIR)/browserpass/hosts/chromium/$(APP_ID).json"
	$(INSTALL) -Dm644 browser-files/chromium-policy.json "$(LIB_DIR)/browserpass/policies/chromium/$(APP_ID).json"
	$(INSTALL) -Dm644 browser-files/firefox-host.json    "$(LIB_DIR)/browserpass/hosts/firefox/$(APP_ID).json"

# Browser-specific hosts targets

.PHONY: hosts-chromium
hosts-chromium:
	@case $(OS) in \
	Linux)      mkdir -p "/etc/chromium/native-messaging-hosts/"; \
	            ln -sfv "$(LIB_DIR)/browserpass/hosts/chromium/$(APP_ID).json" "/etc/chromium/native-messaging-hosts/$(APP_ID).json"; \
	            [ -e "/etc/chromium/native-messaging-hosts/$(APP_ID).json" ] || echo "Error: the symlink points to a non-existent location" >&2; \
	            ;; \
	Darwin)     mkdir -p "/Library/Application Support/Chromium/NativeMessagingHosts/"; \
	            ln -sfv "$(LIB_DIR)/browserpass/hosts/chromium/$(APP_ID).json" "/Library/Application Support/Chromium/NativeMessagingHosts/$(APP_ID).json"; \
	            [ -e "/Library/Application Support/Chromium/NativeMessagingHosts/$(APP_ID).json" ] || echo "Error: the symlink points to a non-existent location" >&2; \
	            ;; \
	*)          echo "The operating system $(OS) is not supported"; exit 1 ;; \
	esac

.PHONY: hosts-chromium-user
hosts-chromium-user:
	@case $(OS) in \
	Linux|*BSD) mkdir -p "${HOME}/.config/chromium/NativeMessagingHosts/"; \
	            ln -sfv "$(LIB_DIR)/browserpass/hosts/chromium/$(APP_ID).json" "${HOME}/.config/chromium/NativeMessagingHosts/$(APP_ID).json"; \
	            [ -e "${HOME}/.config/chromium/NativeMessagingHosts/$(APP_ID).json" ] || echo "Error: the symlink points to a non-existent location" >&2; \
	            ;; \
	Darwin)     mkdir -p "${HOME}/Library/Application Support/Chromium/NativeMessagingHosts/"; \
	            ln -sfv "$(LIB_DIR)/browserpass/hosts/chromium/$(APP_ID).json" "${HOME}/Library/Application Support/Chromium/NativeMessagingHosts/$(APP_ID).json"; \
	            [ -e "${HOME}/Library/Application Support/Chromium/NativeMessagingHosts/$(APP_ID).json" ] || echo "Error: the symlink points to a non-existent location" >&2; \
	            ;; \
	*)          echo "The operating system $(OS) is not supported"; exit 1 ;; \
	esac

.PHONY: hosts-chrome
hosts-chrome:
	@case $(OS) in \
	Linux)      mkdir -p "/etc/opt/chrome/native-messaging-hosts/"; \
	            ln -sfv "$(LIB_DIR)/browserpass/hosts/chromium/$(APP_ID).json" "/etc/opt/chrome/native-messaging-hosts/$(APP_ID).json"; \
	            [ -e "/etc/opt/chrome/native-messaging-hosts/$(APP_ID).json" ] || echo "Error: the symlink points to a non-existent location" >&2; \
	            ;; \
	Darwin)     mkdir -p "/Library/Google/Chrome/NativeMessagingHosts/"; \
	            ln -sfv "$(LIB_DIR)/browserpass/hosts/chromium/$(APP_ID).json" "/Library/Google/Chrome/NativeMessagingHosts/$(APP_ID).json"; \
	            [ -e "/Library/Google/Chrome/NativeMessagingHosts/$(APP_ID).json" ] || echo "Error: the symlink points to a non-existent location" >&2; \
	            ;; \
	*)          echo "The operating system $(OS) is not supported"; exit 1 ;; \
	esac

.PHONY: hosts-chrome-user
hosts-chrome-user:
	@case $(OS) in \
	Linux|*BSD) mkdir -p "${HOME}/.config/google-chrome/NativeMessagingHosts/"; \
	            ln -sfv "$(LIB_DIR)/browserpass/hosts/chromium/$(APP_ID).json" "${HOME}/.config/google-chrome/NativeMessagingHosts/$(APP_ID).json"; \
	            [ -e "${HOME}/.config/google-chrome/NativeMessagingHosts/$(APP_ID).json" ] || echo "Error: the symlink points to a non-existent location" >&2; \
	            ;; \
	Darwin)     mkdir -p "${HOME}/Library/Application Support/Google/Chrome/NativeMessagingHosts/"; \
	            ln -sfv "$(LIB_DIR)/browserpass/hosts/chromium/$(APP_ID).json" "${HOME}/Library/Application Support/Google/Chrome/NativeMessagingHosts/$(APP_ID).json"; \
	            [ -e "${HOME}/Library/Application Support/Google/Chrome/NativeMessagingHosts/$(APP_ID).json" ] || echo "Error: the symlink points to a non-existent location" >&2; \
	            ;; \
	*)          echo "The operating system $(OS) is not supported"; exit 1 ;; \
	esac

.PHONY: hosts-vivaldi
hosts-vivaldi:
	@case $(OS) in \
	Linux)      mkdir -p "/etc/opt/vivaldi/native-messaging-hosts/"; \
	            ln -sfv "$(LIB_DIR)/browserpass/hosts/chromium/$(APP_ID).json" "/etc/opt/vivaldi/native-messaging-hosts/$(APP_ID).json"; \
	            [ -e "/etc/opt/vivaldi/native-messaging-hosts/$(APP_ID).json" ] || echo "Error: the symlink points to a non-existent location" >&2; \
	            ;; \
	Darwin)     mkdir -p "/Library/Application Support/Vivaldi/NativeMessagingHosts/"; \
	            ln -sfv "$(LIB_DIR)/browserpass/hosts/chromium/$(APP_ID).json" "/Library/Application Support/Vivaldi/NativeMessagingHosts/$(APP_ID).json"; \
	            [ -e "/Library/Application Support/Vivaldi/NativeMessagingHosts/$(APP_ID).json" ] || echo "Error: the symlink points to a non-existent location" >&2; \
	            ;; \
	*)          echo "The operating system $(OS) is not supported"; exit 1 ;; \
	esac

.PHONY: hosts-vivaldi-user
hosts-vivaldi-user:
	@case $(OS) in \
	Linux|*BSD) mkdir -p "${HOME}/.config/vivaldi/NativeMessagingHosts/"; \
	            ln -sfv "$(LIB_DIR)/browserpass/hosts/chromium/$(APP_ID).json" "${HOME}/.config/vivaldi/NativeMessagingHosts/$(APP_ID).json"; \
	            [ -e "${HOME}/.config/vivaldi/NativeMessagingHosts/$(APP_ID).json" ] || echo "Error: the symlink points to a non-existent location" >&2; \
	            ;; \
	Darwin)     mkdir -p "${HOME}/Library/Application Support/Vivaldi/NativeMessagingHosts/"; \
	            ln -sfv "$(LIB_DIR)/browserpass/hosts/chromium/$(APP_ID).json" "${HOME}/Library/Application Support/Vivaldi/NativeMessagingHosts/$(APP_ID).json"; \
	            [ -e "${HOME}/Library/Application Support/Vivaldi/NativeMessagingHosts/$(APP_ID).json" ] || echo "Error: the symlink points to a non-existent location" >&2; \
	            ;; \
	*)          echo "The operating system $(OS) is not supported"; exit 1 ;; \
	esac

.PHONY: hosts-brave
hosts-brave:
	@case $(OS) in \
	Linux)      mkdir -p "/etc/opt/chrome/native-messaging-hosts/"; \
	            ln -sfv "$(LIB_DIR)/browserpass/hosts/chromium/$(APP_ID).json" "/etc/opt/chrome/native-messaging-hosts/$(APP_ID).json"; \
	            [ -e "/etc/opt/chrome/native-messaging-hosts/$(APP_ID).json" ] || echo "Error: the symlink points to a non-existent location" >&2; \
	            ;; \
	Darwin)     mkdir -p "/Library/Application Support/Chromium/NativeMessagingHosts/"; \
	            ln -sfv "$(LIB_DIR)/browserpass/hosts/chromium/$(APP_ID).json" "/Library/Application Support/Chromium/NativeMessagingHosts/$(APP_ID).json"; \
	            [ -e "/Library/Application Support/Chromium/NativeMessagingHosts/$(APP_ID).json" ] || echo "Error: the symlink points to a non-existent location" >&2; \
	            ;; \
	*)          echo "The operating system $(OS) is not supported"; exit 1 ;; \
	esac

.PHONY: hosts-brave-user
hosts-brave-user:
	@case $(OS) in \
	Linux|*BSD) mkdir -p "${HOME}/.config/BraveSoftware/Brave-Browser/NativeMessagingHosts/"; \
	            ln -sfv "$(LIB_DIR)/browserpass/hosts/chromium/$(APP_ID).json" "${HOME}/.config/BraveSoftware/Brave-Browser/NativeMessagingHosts/$(APP_ID).json"; \
	            [ -e "${HOME}/.config/BraveSoftware/Brave-Browser/NativeMessagingHosts/$(APP_ID).json" ] || echo "Error: the symlink points to a non-existent location" >&2; \
	            ;; \
	Darwin)     mkdir -p "${HOME}/Library/Application Support/Google/Chrome/NativeMessagingHosts/"; \
	            ln -sfv "$(LIB_DIR)/browserpass/hosts/chromium/$(APP_ID).json" "${HOME}/Library/Application Support/Google/Chrome/NativeMessagingHosts/$(APP_ID).json"; \
	            [ -e "${HOME}/Library/Application Support/Google/Chrome/NativeMessagingHosts/$(APP_ID).json" ] || echo "Error: the symlink points to a non-existent location" >&2; \
	            ;; \
	*)          echo "The operating system $(OS) is not supported"; exit 1 ;; \
	esac

.PHONY: hosts-firefox
hosts-firefox:
	@case $(OS) in \
	Linux)      mkdir -p "$(LIB_DIR)/mozilla/native-messaging-hosts/"; \
	            ln -sfv "$(LIB_DIR)/browserpass/hosts/firefox/$(APP_ID).json" "/usr/lib/mozilla/native-messaging-hosts/$(APP_ID).json"; \
	            [ -e "/usr/lib/mozilla/native-messaging-hosts/$(APP_ID).json" ] || echo "Error: the symlink points to a non-existent location" >&2; \
	            ;; \
	Darwin)     mkdir -p "/Library/Application Support/Mozilla/NativeMessagingHosts/"; \
	            ln -sfv "$(LIB_DIR)/browserpass/hosts/firefox/$(APP_ID).json" "/Library/Application Support/Mozilla/NativeMessagingHosts/$(APP_ID).json"; \
	            [ -e "/Library/Application Support/Mozilla/NativeMessagingHosts/$(APP_ID).json" ] || echo "Error: the symlink points to a non-existent location" >&2; \
	            ;; \
	*)          echo "The operating system $(OS) is not supported"; exit 1 ;; \
	esac

.PHONY: hosts-firefox-user
hosts-firefox-user:
	@case $(OS) in \
	Linux|*BSD) mkdir -p "${HOME}/.mozilla/native-messaging-hosts/"; \
	            ln -sfv "$(LIB_DIR)/browserpass/hosts/firefox/$(APP_ID).json" "${HOME}/.mozilla/native-messaging-hosts/$(APP_ID).json"; \
	            [ -e "${HOME}/.mozilla/native-messaging-hosts/$(APP_ID).json" ] || echo "Error: the symlink points to a non-existent location" >&2; \
	            ;; \
	Darwin)     mkdir -p "${HOME}/Library/Application Support/Mozilla/NativeMessagingHosts/"; \
	            ln -sfv "$(LIB_DIR)/browserpass/hosts/firefox/$(APP_ID).json" "${HOME}/Library/Application Support/Mozilla/NativeMessagingHosts/$(APP_ID).json"; \
	            [ -e "${HOME}/Library/Application Support/Mozilla/NativeMessagingHosts/$(APP_ID).json" ] || echo "Error: the symlink points to a non-existent location" >&2; \
	            ;; \
	*)          echo "The operating system $(OS) is not supported"; exit 1 ;; \
	esac

# Browser-specific policies targets

.PHONY: policies-chromium
policies-chromium:
	@case $(OS) in \
	Linux)      mkdir -p "/etc/chromium/policies/managed/"; \
	            ln -sfv "$(LIB_DIR)/browserpass/policies/chromium/$(APP_ID).json" "/etc/chromium/policies/managed/$(APP_ID).json"; \
	            [ -e "/etc/chromium/policies/managed/$(APP_ID).json" ] || echo "Error: the symlink points to a non-existent location" >&2; \
	            ;; \
	Darwin)     mkdir -p "/Library/Application Support/Chromium/policies/managed/"; \
	            ln -sfv "$(LIB_DIR)/browserpass/policies/chromium/$(APP_ID).json" "/Library/Application Support/Chromium/policies/managed/$(APP_ID).json"; \
	            [ -e "/Library/Application Support/Chromium/policies/managed/$(APP_ID).json" ] || echo "Error: the symlink points to a non-existent location" >&2; \
	            ;; \
	*)          echo "The operating system $(OS) is not supported"; exit 1 ;; \
	esac

.PHONY: policies-chromium-user
policies-chromium-user:
	@case $(OS) in \
	Linux|*BSD) mkdir -p "${HOME}/.config/chromium/policies/managed/"; \
	            ln -sfv "$(LIB_DIR)/browserpass/policies/chromium/$(APP_ID).json" "${HOME}/.config/chromium/policies/managed/$(APP_ID).json"; \
	            [ -e "${HOME}/.config/chromium/policies/managed/$(APP_ID).json" ] || echo "Error: the symlink points to a non-existent location" >&2; \
	            ;; \
	Darwin)     mkdir -p "${HOME}/Library/Application Support/Chromium/policies/managed/"; \
	            ln -sfv "$(LIB_DIR)/browserpass/policies/chromium/$(APP_ID).json" "${HOME}/Library/Application Support/Chromium/policies/managed/$(APP_ID).json"; \
	            [ -e "${HOME}/Library/Application Support/Chromium/policies/managed/$(APP_ID).json" ] || echo "Error: the symlink points to a non-existent location" >&2; \
	            ;; \
	*)          echo "The operating system $(OS) is not supported"; exit 1 ;; \
	esac

.PHONY: policies-chrome
policies-chrome:
	@case $(OS) in \
	Linux)      mkdir -p "/etc/opt/chrome/policies/managed/"; \
	            ln -sfv "$(LIB_DIR)/browserpass/policies/chromium/$(APP_ID).json" "/etc/opt/chrome/policies/managed/$(APP_ID).json"; \
	            [ -e "/etc/opt/chrome/policies/managed/$(APP_ID).json" ] || echo "Error: the symlink points to a non-existent location" >&2; \
	            ;; \
	Darwin)     mkdir -p "/Library/Google/Chrome/policies/managed/"; \
	            ln -sfv "$(LIB_DIR)/browserpass/policies/chromium/$(APP_ID).json" "/Library/Google/Chrome/policies/managed/$(APP_ID).json"; \
	            [ -e "/Library/Google/Chrome/policies/managed/$(APP_ID).json" ] || echo "Error: the symlink points to a non-existent location" >&2; \
	            ;; \
	*)          echo "The operating system $(OS) is not supported"; exit 1 ;; \
	esac

.PHONY: policies-chrome-user
policies-chrome-user:
	@case $(OS) in \
	Linux|*BSD) mkdir -p "${HOME}/.config/google-chrome/policies/managed/"; \
	            ln -sfv "$(LIB_DIR)/browserpass/policies/chromium/$(APP_ID).json" "${HOME}/.config/google-chrome/policies/managed/$(APP_ID).json"; \
	            [ -e "${HOME}/.config/google-chrome/policies/managed/$(APP_ID).json" ] || echo "Error: the symlink points to a non-existent location" >&2; \
	            ;; \
	Darwin)     mkdir -p "${HOME}/Library/Application Support/Google/Chrome/policies/managed/"; \
	            ln -sfv "$(LIB_DIR)/browserpass/policies/chromium/$(APP_ID).json" "${HOME}/Library/Application Support/Google/Chrome/policies/managed/$(APP_ID).json"; \
	            [ -e "${HOME}/Library/Application Support/Google/Chrome/policies/managed/$(APP_ID).json" ] || echo "Error: the symlink points to a non-existent location" >&2; \
	            ;; \
	*)          echo "The operating system $(OS) is not supported"; exit 1 ;; \
	esac

.PHONY: policies-vivaldi
policies-vivaldi:
	@case $(OS) in \
	Linux)      mkdir -p "/etc/opt/vivaldi/policies/managed/"; \
	            ln -sfv "$(LIB_DIR)/browserpass/policies/chromium/$(APP_ID).json" "/etc/opt/vivaldi/policies/managed/$(APP_ID).json"; \
	            [ -e "/etc/opt/vivaldi/policies/managed/$(APP_ID).json" ] || echo "Error: the symlink points to a non-existent location" >&2; \
	            ;; \
	Darwin)     mkdir -p "/Library/Application Support/Vivaldi/policies/managed/"; \
	            ln -sfv "$(LIB_DIR)/browserpass/policies/chromium/$(APP_ID).json" "/Library/Application Support/Vivaldi/policies/managed/$(APP_ID).json"; \
	            [ -e "/Library/Application Support/Vivaldi/policies/managed/$(APP_ID).json" ] || echo "Error: the symlink points to a non-existent location" >&2; \
	            ;; \
	*)          echo "The operating system $(OS) is not supported"; exit 1 ;; \
	esac

.PHONY: policies-vivaldi-user
policies-vivaldi-user:
	@case $(OS) in \
	Linux|*BSD) mkdir -p "${HOME}/.config/vivaldi/policies/managed/"; \
	            ln -sfv "$(LIB_DIR)/browserpass/policies/chromium/$(APP_ID).json" "${HOME}/.config/vivaldi/policies/managed/$(APP_ID).json"; \
	            [ -e "${HOME}/.config/vivaldi/policies/managed/$(APP_ID).json" ] || echo "Error: the symlink points to a non-existent location" >&2; \
	            ;; \
	Darwin)     mkdir -p "${HOME}/Library/Application Support/Vivaldi/policies/managed/"; \
	            ln -sfv "$(LIB_DIR)/browserpass/policies/chromium/$(APP_ID).json" "${HOME}/Library/Application Support/Vivaldi/policies/managed/$(APP_ID).json"; \
	            [ -e "${HOME}/Library/Application Support/Vivaldi/policies/managed/$(APP_ID).json" ] || echo "Error: the symlink points to a non-existent location" >&2; \
	            ;; \
	*)          echo "The operating system $(OS) is not supported"; exit 1 ;; \
	esac

.PHONY: policies-brave
policies-brave:
	@case $(OS) in \
	Linux)      mkdir -p "/etc/opt/chrome/policies/managed/"; \
	            ln -sfv "$(LIB_DIR)/browserpass/policies/chromium/$(APP_ID).json" "/etc/opt/chrome/policies/managed/$(APP_ID).json"; \
	            [ -e "/etc/opt/chrome/policies/managed/$(APP_ID).json" ] || echo "Error: the symlink points to a non-existent location" >&2; \
	            ;; \
	Darwin)     mkdir -p "/Library/Application Support/Chromium/policies/managed/"; \
	            ln -sfv "$(LIB_DIR)/browserpass/policies/chromium/$(APP_ID).json" "/Library/Application Support/Chromium/policies/managed/$(APP_ID).json"; \
	            [ -e "/Library/Application Support/Chromium/policies/managed/$(APP_ID).json" ] || echo "Error: the symlink points to a non-existent location" >&2; \
	            ;; \
	*)          echo "The operating system $(OS) is not supported"; exit 1 ;; \
	esac

.PHONY: policies-brave-user
policies-brave-user:
	@case $(OS) in \
	Linux|*BSD) mkdir -p "${HOME}/.config/BraveSoftware/Brave-Browser/policies/managed/"; \
	            ln -sfv "$(LIB_DIR)/browserpass/policies/chromium/$(APP_ID).json" "${HOME}/.config/BraveSoftware/Brave-Browser/policies/managed/$(APP_ID).json"; \
	            [ -e "${HOME}/.config/BraveSoftware/Brave-Browser/policies/managed/$(APP_ID).json" ] || echo "Error: the symlink points to a non-existent location" >&2; \
	            ;; \
	Darwin)     mkdir -p "${HOME}/Library/Application Support/Google/Chrome/policies/managed/"; \
	            ln -sfv "$(LIB_DIR)/browserpass/policies/chromium/$(APP_ID).json" "${HOME}/Library/Application Support/Google/Chrome/policies/managed/$(APP_ID).json"; \
	            [ -e "${HOME}/Library/Application Support/Google/Chrome/policies/managed/$(APP_ID).json" ] || echo "Error: the symlink points to a non-existent location" >&2; \
	            ;; \
	*)          echo "The operating system $(OS) is not supported"; exit 1 ;; \
	esac
