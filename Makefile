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

#######################
# For local development

.PHONY: all
all: browserpass test

browserpass: *.go **/*.go
	go build -ldflags $(GO_LDFLAGS) -gcflags $(GO_GCFLAGS) -asmflags $(GO_ASMFLAGS) -o $@

browserpass-linux64: *.go **/*.go
	env GOOS=linux GOARCH=amd64 go build -ldflags $(GO_LDFLAGS) -gcflags $(GO_GCFLAGS) -asmflags $(GO_ASMFLAGS) -o $@

browserpass-darwinx64: *.go **/*.go
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
dist: clean browserpass-linux64 browserpass-darwinx64 browserpass-openbsd64 browserpass-freebsd64 browserpass-windows64
	mkdir -p dist

	git archive -o dist/$(VERSION).tar.gz --format tar.gz --prefix=browserpass-native-$(VERSION)/ $(VERSION)

	zip -FSr dist/browserpass-linux64   browserpass-linux64       browser-files/* Makefile README.md LICENSE
	zip -FSr dist/browserpass-darwinx64 browserpass-darwinx64     browser-files/* Makefile README.md LICENSE
	zip -FSr dist/browserpass-openbsd64 browserpass-openbsd64     browser-files/* Makefile README.md LICENSE
	zip -FSr dist/browserpass-freebsd64 browserpass-freebsd64     browser-files/* Makefile README.md LICENSE
	zip -FSr dist/browserpass-windows64 browserpass-windows64.exe browser-files/* Makefile README.md LICENSE

	for file in dist/*; do \
	    gpg --detach-sign "$$file"; \
	done

	rm -f dist/$(VERSION).tar.gz

#######################
# For user installation

.PHONY: install
install:
	install -Dm755 -t "$(BIN_DIR)/" $(BIN)
	install -Dm644 -t "$(LIB_DIR)/browserpass/" Makefile
	install -Dm644 -t "$(SHARE_DIR)/licenses/browserpass/" LICENSE
	install -Dm644 -t "$(SHARE_DIR)/doc/browserpass/" README.md

	install -Dm644 browser-files/chromium-host.json   "$(LIB_DIR)/browserpass/hosts/chromium/$(APP_ID).json"
	install -Dm644 browser-files/chromium-policy.json "$(LIB_DIR)/browserpass/policies/chromium/$(APP_ID).json"
	install -Dm644 browser-files/firefox-host.json    "$(LIB_DIR)/browserpass/hosts/firefox/$(APP_ID).json"

	sed -i "s|%%replace%%|/usr/bin/$(BIN)|" "$(LIB_DIR)/browserpass/hosts/chromium/$(APP_ID).json"
	sed -i "s|%%replace%%|/usr/bin/$(BIN)|" "$(LIB_DIR)/browserpass/hosts/firefox/$(APP_ID).json"

# Browser-specific hosts targets

.PHONY: hosts-chromium
hosts-chromium:
	@case $(OS) in \
	Linux)      mkdir -p "/etc/chromium/native-messaging-hosts/"; \
	            ln -sf "$(LIB_DIR)/browserpass/hosts/chromium/$(APP_ID).json" "/etc/chromium/native-messaging-hosts/" \
	            ;; \
	Darwin)     mkdir -p "/Library/Application Support/Chromium/NativeMessagingHosts/"; \
	            ln -sf "$(LIB_DIR)/browserpass/hosts/chromium/$(APP_ID).json" "/Library/Application Support/Chromium/NativeMessagingHosts/" \
	            ;; \
	*)          echo "The operating system $(OS) is not supported"; exit 1 ;; \
	esac

.PHONY: hosts-chromium-user
hosts-chromium-user:
	@case $(OS) in \
	Linux|*BSD) mkdir -p "${HOME}/.config/chromium/NativeMessagingHosts/"; \
	            ln -sf "$(LIB_DIR)/browserpass/hosts/chromium/$(APP_ID).json" "${HOME}/.config/chromium/NativeMessagingHosts/" \
	            ;; \
	Darwin)     mkdir -p "${HOME}/Library/Application Support/Chromium/NativeMessagingHosts/"; \
	            ln -sf "$(LIB_DIR)/browserpass/hosts/chromium/$(APP_ID).json" "${HOME}/Library/Application Support/Chromium/NativeMessagingHosts/" \
	            ;; \
	*)          echo "The operating system $(OS) is not supported"; exit 1 ;; \
	esac

.PHONY: hosts-chrome
hosts-chrome:
	@case $(OS) in \
	Linux)      mkdir -p "/etc/opt/chrome/native-messaging-hosts/"; \
	            ln -sf "$(LIB_DIR)/browserpass/hosts/chromium/$(APP_ID).json" "/etc/opt/chrome/native-messaging-hosts/" \
	            ;; \
	Darwin)     mkdir -p "/Library/Google/Chrome/NativeMessagingHosts/"; \
	            ln -sf "$(LIB_DIR)/browserpass/hosts/chromium/$(APP_ID).json" "/Library/Google/Chrome/NativeMessagingHosts/" \
	            ;; \
	*)          echo "The operating system $(OS) is not supported"; exit 1 ;; \
	esac

.PHONY: hosts-chrome-user
hosts-chrome-user:
	@case $(OS) in \
	Linux|*BSD) mkdir -p "${HOME}/.config/google-chrome/NativeMessagingHosts/"; \
	            ln -sf "$(LIB_DIR)/browserpass/hosts/chromium/$(APP_ID).json" "${HOME}/.config/google-chrome/NativeMessagingHosts/" \
	            ;; \
	Darwin)     mkdir -p "${HOME}/Library/Application Support/Google/Chrome/NativeMessagingHosts/"; \
	            ln -sf "$(LIB_DIR)/browserpass/hosts/chromium/$(APP_ID).json" "${HOME}/Library/Application Support/Google/Chrome/NativeMessagingHosts/" \
	            ;; \
	*)          echo "The operating system $(OS) is not supported"; exit 1 ;; \
	esac

.PHONY: hosts-vivaldi
hosts-vivaldi:
	@case $(OS) in \
	Linux)      mkdir -p "/etc/opt/vivaldi/native-messaging-hosts/"; \
	            ln -sf "$(LIB_DIR)/browserpass/hosts/chromium/$(APP_ID).json" "/etc/opt/vivaldi/native-messaging-hosts/" \
	            ;; \
	Darwin)     mkdir -p "/Library/Application Support/Vivaldi/NativeMessagingHosts/"; \
	            ln -sf "$(LIB_DIR)/browserpass/hosts/chromium/$(APP_ID).json" "/Library/Application Support/Vivaldi/NativeMessagingHosts/" \
	            ;; \
	*)          echo "The operating system $(OS) is not supported"; exit 1 ;; \
	esac

.PHONY: hosts-vivaldi-user
hosts-vivaldi-user:
	@case $(OS) in \
	Linux|*BSD) mkdir -p "${HOME}/.config/vivaldi/NativeMessagingHosts/"; \
	            ln -sf "$(LIB_DIR)/browserpass/hosts/chromium/$(APP_ID).json" "${HOME}/.config/vivaldi/NativeMessagingHosts/" \
	            ;; \
	Darwin)     mkdir -p "${HOME}/Library/Application Support/Vivaldi/NativeMessagingHosts/"; \
	            ln -sf "$(LIB_DIR)/browserpass/hosts/chromium/$(APP_ID).json" "${HOME}/Library/Application Support/Vivaldi/NativeMessagingHosts/" \
	            ;; \
	*)          echo "The operating system $(OS) is not supported"; exit 1 ;; \
	esac

.PHONY: hosts-brave
hosts-brave:
	@case $(OS) in \
	Linux)      mkdir -p "/etc/opt/brave/native-messaging-hosts/"; \
	            ln -sf "$(LIB_DIR)/browserpass/hosts/chromium/$(APP_ID).json" "/etc/opt/brave/native-messaging-hosts/" \
	            ;; \
	Darwin)     mkdir -p "/Library/Application Support/BraveSoftware/Brave-Browser/NativeMessagingHosts/"; \
	            ln -sf "$(LIB_DIR)/browserpass/hosts/chromium/$(APP_ID).json" "/Library/Application Support/BraveSoftware/Brave-Browser/NativeMessagingHosts/" \
	            ;; \
	*)          echo "The operating system $(OS) is not supported"; exit 1 ;; \
	esac

.PHONY: hosts-brave-user
hosts-brave-user:
	@case $(OS) in \
	Linux|*BSD) mkdir -p "${HOME}/.config/BraveSoftware/Brave-Browser/NativeMessagingHosts/"; \
	            ln -sf "$(LIB_DIR)/browserpass/hosts/chromium/$(APP_ID).json" "${HOME}/.config/BraveSoftware/Brave-Browser/NativeMessagingHosts/" \
	            ;; \
	Darwin)     mkdir -p "${HOME}/Library/Application Support/BraveSoftware/Brave-Browser/NativeMessagingHosts/"; \
	            ln -sf "$(LIB_DIR)/browserpass/hosts/chromium/$(APP_ID).json" "${HOME}/Library/Application Support/BraveSoftware/Brave-Browser/NativeMessagingHosts/" \
	            ;; \
	*)          echo "The operating system $(OS) is not supported"; exit 1 ;; \
	esac

.PHONY: hosts-firefox
hosts-firefox:
	@case $(OS) in \
	Linux)      mkdir -p "$(LIB_DIR)/mozilla/native-messaging-hosts/"; \
	            ln -sf "$(LIB_DIR)/browserpass/hosts/firefox/$(APP_ID).json" "/usr/lib/mozilla/native-messaging-hosts/" \
	            ;; \
	Darwin)     mkdir -p "/Library/Application Support/Mozilla/NativeMessagingHosts/"; \
	            ln -sf "$(LIB_DIR)/browserpass/hosts/firefox/$(APP_ID).json" "/Library/Application Support/Mozilla/NativeMessagingHosts/" \
	            ;; \
	*)          echo "The operating system $(OS) is not supported"; exit 1 ;; \
	esac

.PHONY: hosts-firefox-user
hosts-firefox-user:
	@case $(OS) in \
	Linux|*BSD) mkdir -p "${HOME}/.mozilla/native-messaging-hosts/"; \
	            ln -sf "$(LIB_DIR)/browserpass/hosts/firefox/$(APP_ID).json" "${HOME}/.mozilla/native-messaging-hosts/" \
	            ;; \
	Darwin)     mkdir -p "${HOME}/Library/Application Support/Mozilla/NativeMessagingHosts/"; \
	            ln -sf "$(LIB_DIR)/browserpass/hosts/firefox/$(APP_ID).json" "${HOME}/Library/Application Support/Mozilla/NativeMessagingHosts/" \
	            ;; \
	*)          echo "The operating system $(OS) is not supported"; exit 1 ;; \
	esac

# Browser-specific policies targets

.PHONY: policies-chromium
policies-chromium:
	@case $(OS) in \
	Linux)      mkdir -p "/etc/chromium/policies/managed/"; \
	            ln -sf "$(LIB_DIR)/browserpass/hosts/chromium/$(APP_ID).json" "/etc/chromium/policies/managed/" \
	            ;; \
	Darwin)     mkdir -p "/Library/Application Support/Chromium/policies/managed/"; \
	            ln -sf "$(LIB_DIR)/browserpass/hosts/chromium/$(APP_ID).json" "/Library/Application Support/Chromium/policies/managed/" \
	            ;; \
	*)          echo "The operating system $(OS) is not supported"; exit 1 ;; \
	esac

.PHONY: policies-chromium-user
policies-chromium-user:
	@case $(OS) in \
	Linux|*BSD) mkdir -p "${HOME}/.config/chromium/policies/managed/"; \
	            ln -sf "$(LIB_DIR)/browserpass/hosts/chromium/$(APP_ID).json" "${HOME}/.config/chromium/policies/managed/" \
	            ;; \
	Darwin)     mkdir -p "${HOME}/Library/Application Support/Chromium/policies/managed/"; \
	            ln -sf "$(LIB_DIR)/browserpass/hosts/chromium/$(APP_ID).json" "${HOME}/Library/Application Support/Chromium/policies/managed/" \
	            ;; \
	*)          echo "The operating system $(OS) is not supported"; exit 1 ;; \
	esac

.PHONY: policies-chrome
policies-chrome:
	@case $(OS) in \
	Linux)      mkdir -p "/etc/opt/chrome/policies/managed/"; \
	            ln -sf "$(LIB_DIR)/browserpass/hosts/chromium/$(APP_ID).json" "/etc/opt/chrome/policies/managed/" \
	            ;; \
	Darwin)     mkdir -p "/Library/Google/Chrome/policies/managed/"; \
	            ln -sf "$(LIB_DIR)/browserpass/hosts/chromium/$(APP_ID).json" "/Library/Google/Chrome/policies/managed/" \
	            ;; \
	*)          echo "The operating system $(OS) is not supported"; exit 1 ;; \
	esac

.PHONY: policies-chrome-user
policies-chrome-user:
	@case $(OS) in \
	Linux|*BSD) mkdir -p "${HOME}/.config/google-chrome/policies/managed/"; \
	            ln -sf "$(LIB_DIR)/browserpass/hosts/chromium/$(APP_ID).json" "${HOME}/.config/google-chrome/policies/managed/" \
	            ;; \
	Darwin)     mkdir -p "${HOME}/Library/Application Support/Google/Chrome/policies/managed/"; \
	            ln -sf "$(LIB_DIR)/browserpass/hosts/chromium/$(APP_ID).json" "${HOME}/Library/Application Support/Google/Chrome/policies/managed/" \
	            ;; \
	*)          echo "The operating system $(OS) is not supported"; exit 1 ;; \
	esac

.PHONY: policies-vivaldi
policies-vivaldi:
	@case $(OS) in \
	Linux)      mkdir -p "/etc/opt/vivaldi/policies/managed/"; \
	            ln -sf "$(LIB_DIR)/browserpass/hosts/chromium/$(APP_ID).json" "/etc/opt/vivaldi/policies/managed/" \
	            ;; \
	Darwin)     mkdir -p "/Library/Application Support/Vivaldi/policies/managed/"; \
	            ln -sf "$(LIB_DIR)/browserpass/hosts/chromium/$(APP_ID).json" "/Library/Application Support/Vivaldi/policies/managed/" \
	            ;; \
	*)          echo "The operating system $(OS) is not supported"; exit 1 ;; \
	esac

.PHONY: policies-vivaldi-user
policies-vivaldi-user:
	@case $(OS) in \
	Linux|*BSD) mkdir -p "${HOME}/.config/vivaldi/policies/managed/"; \
	            ln -sf "$(LIB_DIR)/browserpass/hosts/chromium/$(APP_ID).json" "${HOME}/.config/vivaldi/policies/managed/" \
	            ;; \
	Darwin)     mkdir -p "${HOME}/Library/Application Support/Vivaldi/policies/managed/"; \
	            ln -sf "$(LIB_DIR)/browserpass/hosts/chromium/$(APP_ID).json" "${HOME}/Library/Application Support/Vivaldi/policies/managed/" \
	            ;; \
	*)          echo "The operating system $(OS) is not supported"; exit 1 ;; \
	esac

.PHONY: policies-brave
policies-brave:
	@case $(OS) in \
	Linux)      mkdir -p "/etc/opt/brave/policies/managed/"; \
	            ln -sf "$(LIB_DIR)/browserpass/hosts/chromium/$(APP_ID).json" "/etc/opt/brave/policies/managed/" \
	            ;; \
	Darwin)     mkdir -p "/Library/Application Support/BraveSoftware/Brave-Browser/policies/managed/"; \
	            ln -sf "$(LIB_DIR)/browserpass/hosts/chromium/$(APP_ID).json" "/Library/Application Support/BraveSoftware/Brave-Browser/policies/managed/" \
	            ;; \
	*)          echo "The operating system $(OS) is not supported"; exit 1 ;; \
	esac

.PHONY: policies-brave-user
policies-brave-user:
	@case $(OS) in \
	Linux|*BSD) mkdir -p "${HOME}/.config/BraveSoftware/Brave-Browser/policies/managed/"; \
	            ln -sf "$(LIB_DIR)/browserpass/hosts/chromium/$(APP_ID).json" "${HOME}/.config/BraveSoftware/Brave-Browser/policies/managed/" \
	            ;; \
	Darwin)     mkdir -p "${HOME}/Library/Application Support/BraveSoftware/Brave-Browser/policies/managed/"; \
	            ln -sf "$(LIB_DIR)/browserpass/hosts/chromium/$(APP_ID).json" "${HOME}/Library/Application Support/BraveSoftware/Brave-Browser/policies/managed/" \
	            ;; \
	*)          echo "The operating system $(OS) is not supported"; exit 1 ;; \
	esac
