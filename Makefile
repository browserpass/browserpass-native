BIN ?= browserpass

APP_ID = com.github.browserpass.native
OS = $(shell uname -s)

#######################
# For local development

.PHONY: all
all: browserpass test

browserpass: *.go **/*.go
	go build -o $@

browserpass-linux64: *.go **/*.go
	env GOOS=linux GOARCH=amd64 go build -o $@

browserpass-darwinx64: *.go **/*.go
	env GOOS=darwin GOARCH=amd64 go build -o $@

browserpass-openbsd64: *.go **/*.go
	env GOOS=openbsd GOARCH=amd64 go build -o $@

browserpass-freebsd64: *.go **/*.go
	env GOOS=freebsd GOARCH=amd64 go build -o $@

browserpass-windows64: *.go **/*.go
	env GOOS=windows GOARCH=amd64 go build -o $@.exe

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
	zip -FS dist/browserpass-linux64   browserpass-linux64       browser-files/* Makefile README.md LICENSE
	zip -FS dist/browserpass-darwinx64 browserpass-darwinx64     browser-files/* Makefile README.md LICENSE
	zip -FS dist/browserpass-openbsd64 browserpass-openbsd64     browser-files/* Makefile README.md LICENSE
	zip -FS dist/browserpass-freebsd64 browserpass-freebsd64     browser-files/* Makefile README.md LICENSE
	zip -FS dist/browserpass-windows64 browserpass-windows64.exe browser-files/* Makefile README.md LICENSE

	for file in dist/*; do \
        gpg --detach-sign "$$file"; \
    done

#######################
# For user installation

.PHONY: install
install:
	install -Dm755 -t "$(DESTDIR)/usr/bin/" $(BIN)
	install -Dm644 -t "$(DESTDIR)/usr/lib/browserpass/" Makefile
	install -Dm644 -t "$(DESTDIR)/usr/share/licenses/browserpass/" LICENSE
	install -Dm644 -t "$(DESTDIR)/usr/share/doc/browserpass/" README.md

	install -Dm644 browser-files/chromium-host.json   "$(DESTDIR)/usr/lib/browserpass/hosts/chromium/$(APP_ID).json"
	install -Dm644 browser-files/chromium-policy.json "$(DESTDIR)/usr/lib/browserpass/policies/chromium/$(APP_ID).json"
	install -Dm644 browser-files/firefox-host.json    "$(DESTDIR)/usr/lib/browserpass/hosts/firefox/$(APP_ID).json"

	sed -i "s|%%replace%%|$(DESTDIR)/usr/bin/$(BIN)|" "$(DESTDIR)/usr/lib/browserpass/hosts/chromium/$(APP_ID).json"
	sed -i "s|%%replace%%|$(DESTDIR)/usr/bin/$(BIN)|" "$(DESTDIR)/usr/lib/browserpass/hosts/firefox/$(APP_ID).json"

# Browser-specific hosts targets

.PHONY: hosts-chromium
hosts-chromium:
	@case $(OS) in \
	Linux)  			   ln -sf "$(DESTDIR)/usr/lib/browserpass/hosts/chromium/$(APP_ID).json" "/etc/chromium/native-messaging-hosts/" ;; \
	Darwin) 			   ln -sf "$(DESTDIR)/usr/lib/browserpass/hosts/chromium/$(APP_ID).json" "/Library/Application Support/Chromium/NativeMessagingHosts/" ;; \
	*)      			   echo "The operating system $(OS) is not supported"; exit 1 ;; \
	esac

.PHONY: hosts-chromium-user
hosts-chromium-user:
	@case $(OS) in \
	Linux|OpenBSD|FreeBSD) ln -sf "$(DESTDIR)/usr/lib/browserpass/hosts/chromium/$(APP_ID).json" "${HOME}/.config/chromium/NativeMessagingHosts/" ;; \
	Darwin)                ln -sf "$(DESTDIR)/usr/lib/browserpass/hosts/chromium/$(APP_ID).json" "${HOME}/Library/Application Support/Chromium/NativeMessagingHosts/" ;; \
	*)                     echo "The operating system $(OS) is not supported"; exit 1 ;; \
	esac

.PHONY: hosts-chrome
hosts-chrome:
	@case $(OS) in \
	Linux)  			   ln -sf "$(DESTDIR)/usr/lib/browserpass/hosts/chromium/$(APP_ID).json" "/etc/opt/chrome/native-messaging-hosts/" ;; \
	Darwin) 			   ln -sf "$(DESTDIR)/usr/lib/browserpass/hosts/chromium/$(APP_ID).json" "/Library/Google/Chrome/NativeMessagingHosts/" ;; \
	*)      			   echo "The operating system $(OS) is not supported"; exit 1 ;; \
	esac

.PHONY: hosts-chrome-user
hosts-chrome-user:
	@case $(OS) in \
	Linux|OpenBSD|FreeBSD) ln -sf "$(DESTDIR)/usr/lib/browserpass/hosts/chromium/$(APP_ID).json" "${HOME}/.config/google-chrome/NativeMessagingHosts/" ;; \
	Darwin)                ln -sf "$(DESTDIR)/usr/lib/browserpass/hosts/chromium/$(APP_ID).json" "${HOME}/Library/Application Support/Google/Chrome/NativeMessagingHosts/" ;; \
	*)                     echo "The operating system $(OS) is not supported"; exit 1 ;; \
	esac

.PHONY: hosts-vivaldi
hosts-vivaldi:
	@case $(OS) in \
	Linux)  			   ln -sf "$(DESTDIR)/usr/lib/browserpass/hosts/chromium/$(APP_ID).json" "/etc/opt/vivaldi/native-messaging-hosts/" ;; \
	Darwin) 			   ln -sf "$(DESTDIR)/usr/lib/browserpass/hosts/chromium/$(APP_ID).json" "/Library/Application Support/Vivaldi/NativeMessagingHosts/" ;; \
	*)      			   echo "The operating system $(OS) is not supported"; exit 1 ;; \
	esac

.PHONY: hosts-vivaldi-user
hosts-vivaldi-user:
	@case $(OS) in \
	Linux|OpenBSD|FreeBSD) ln -sf "$(DESTDIR)/usr/lib/browserpass/hosts/chromium/$(APP_ID).json" "${HOME}/.config/vivaldi/NativeMessagingHosts/" ;; \
	Darwin)                ln -sf "$(DESTDIR)/usr/lib/browserpass/hosts/chromium/$(APP_ID).json" "${HOME}/Library/Application Support/Vivaldi/NativeMessagingHosts/" ;; \
	*)                     echo "The operating system $(OS) is not supported"; exit 1 ;; \
	esac

.PHONY: hosts-brave
hosts-brave:
	@case $(OS) in \
	Linux)  			   ln -sf "$(DESTDIR)/usr/lib/browserpass/hosts/chromium/$(APP_ID).json" "/etc/opt/brave/native-messaging-hosts/" ;; \
	Darwin) 			   ln -sf "$(DESTDIR)/usr/lib/browserpass/hosts/chromium/$(APP_ID).json" "/Library/Application Support/BraveSoftware/Brave-Browser/NativeMessagingHosts/" ;; \
	*)      			   echo "The operating system $(OS) is not supported"; exit 1 ;; \
	esac

.PHONY: hosts-brave-user
hosts-brave-user:
	@case $(OS) in \
	Linux|OpenBSD|FreeBSD) ln -sf "$(DESTDIR)/usr/lib/browserpass/hosts/chromium/$(APP_ID).json" "${HOME}/.config/BraveSoftware/Brave-Browser/NativeMessagingHosts/" ;; \
	Darwin)                ln -sf "$(DESTDIR)/usr/lib/browserpass/hosts/chromium/$(APP_ID).json" "${HOME}/Library/Application Support/BraveSoftware/Brave-Browser/NativeMessagingHosts/" ;; \
	*)                     echo "The operating system $(OS) is not supported"; exit 1 ;; \
	esac

.PHONY: hosts-firefox
hosts-firefox:
	@case $(OS) in \
	Linux)  			   ln -sf "$(DESTDIR)/usr/lib/browserpass/hosts/firefox/$(APP_ID).json" "/usr/lib/mozilla/native-messaging-hosts/" ;; \
	Darwin) 			   ln -sf "$(DESTDIR)/usr/lib/browserpass/hosts/firefox/$(APP_ID).json" "/Library/Application Support/Mozilla/NativeMessagingHosts/" ;; \
	*)      			   echo "The operating system $(OS) is not supported"; exit 1 ;; \
	esac

.PHONY: hosts-firefox-user
hosts-firefox-user:
	@case $(OS) in \
	Linux|OpenBSD|FreeBSD) ln -sf "$(DESTDIR)/usr/lib/browserpass/hosts/firefox/$(APP_ID).json" "${HOME}/.mozilla/native-messaging-hosts/" ;; \
	Darwin)                ln -sf "$(DESTDIR)/usr/lib/browserpass/hosts/firefox/$(APP_ID).json" "${HOME}/Library/Application Support/Mozilla/NativeMessagingHosts/" ;; \
	*)                     echo "The operating system $(OS) is not supported"; exit 1 ;; \
	esac

# Browser-specific policies targets

.PHONY: policies-chromium
policies-chromium:
	@case $(OS) in \
	Linux)  			   ln -sf "$(DESTDIR)/usr/lib/browserpass/hosts/chromium/$(APP_ID).json" "/etc/chromium/policies/managed/" ;; \
	Darwin) 			   ln -sf "$(DESTDIR)/usr/lib/browserpass/hosts/chromium/$(APP_ID).json" "/Library/Application Support/Chromium/policies/managed/" ;; \
	*)      			   echo "The operating system $(OS) is not supported"; exit 1 ;; \
	esac

.PHONY: policies-chromium-user
policies-chromium-user:
	@case $(OS) in \
	Linux|OpenBSD|FreeBSD) ln -sf "$(DESTDIR)/usr/lib/browserpass/hosts/chromium/$(APP_ID).json" "${HOME}/.config/chromium/policies/managed/" ;; \
	Darwin)                ln -sf "$(DESTDIR)/usr/lib/browserpass/hosts/chromium/$(APP_ID).json" "${HOME}/Library/Application Support/Chromium/policies/managed/" ;; \
	*)                     echo "The operating system $(OS) is not supported"; exit 1 ;; \
	esac

.PHONY: policies-chrome
policies-chrome:
	@case $(OS) in \
	Linux)  			   ln -sf "$(DESTDIR)/usr/lib/browserpass/hosts/chromium/$(APP_ID).json" "/etc/opt/chrome/policies/managed/" ;; \
	Darwin) 			   ln -sf "$(DESTDIR)/usr/lib/browserpass/hosts/chromium/$(APP_ID).json" "/Library/Google/Chrome/policies/managed/" ;; \
	*)      			   echo "The operating system $(OS) is not supported"; exit 1 ;; \
	esac

.PHONY: policies-chrome-user
policies-chrome-user:
	@case $(OS) in \
	Linux|OpenBSD|FreeBSD) ln -sf "$(DESTDIR)/usr/lib/browserpass/hosts/chromium/$(APP_ID).json" "${HOME}/.config/google-chrome/policies/managed/" ;; \
	Darwin)                ln -sf "$(DESTDIR)/usr/lib/browserpass/hosts/chromium/$(APP_ID).json" "${HOME}/Library/Application Support/Google/Chrome/policies/managed/" ;; \
	*)                     echo "The operating system $(OS) is not supported"; exit 1 ;; \
	esac

.PHONY: policies-vivaldi
policies-vivaldi:
	@case $(OS) in \
	Linux)  			   ln -sf "$(DESTDIR)/usr/lib/browserpass/hosts/chromium/$(APP_ID).json" "/etc/opt/vivaldi/policies/managed/" ;; \
	Darwin) 			   ln -sf "$(DESTDIR)/usr/lib/browserpass/hosts/chromium/$(APP_ID).json" "/Library/Application Support/Vivaldi/policies/managed/" ;; \
	*)      			   echo "The operating system $(OS) is not supported"; exit 1 ;; \
	esac

.PHONY: policies-vivaldi-user
policies-vivaldi-user:
	@case $(OS) in \
	Linux|OpenBSD|FreeBSD) ln -sf "$(DESTDIR)/usr/lib/browserpass/hosts/chromium/$(APP_ID).json" "${HOME}/.config/vivaldi/policies/managed/" ;; \
	Darwin)                ln -sf "$(DESTDIR)/usr/lib/browserpass/hosts/chromium/$(APP_ID).json" "${HOME}/Library/Application Support/Vivaldi/policies/managed/" ;; \
	*)                     echo "The operating system $(OS) is not supported"; exit 1 ;; \
	esac

.PHONY: policies-brave
policies-brave:
	@case $(OS) in \
	Linux)  			   ln -sf "$(DESTDIR)/usr/lib/browserpass/hosts/chromium/$(APP_ID).json" "/etc/opt/brave/policies/managed/" ;; \
	Darwin) 			   ln -sf "$(DESTDIR)/usr/lib/browserpass/hosts/chromium/$(APP_ID).json" "/Library/Application Support/BraveSoftware/Brave-Browser/policies/managed/" ;; \
	*)      			   echo "The operating system $(OS) is not supported"; exit 1 ;; \
	esac

.PHONY: policies-brave-user
policies-brave-user:
	@case $(OS) in \
	Linux|OpenBSD|FreeBSD) ln -sf "$(DESTDIR)/usr/lib/browserpass/hosts/chromium/$(APP_ID).json" "${HOME}/.config/BraveSoftware/Brave-Browser/policies/managed/" ;; \
	Darwin)                ln -sf "$(DESTDIR)/usr/lib/browserpass/hosts/chromium/$(APP_ID).json" "${HOME}/Library/Application Support/BraveSoftware/Brave-Browser/policies/managed/" ;; \
	*)                     echo "The operating system $(OS) is not supported"; exit 1 ;; \
	esac
