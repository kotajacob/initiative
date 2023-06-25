.POSIX:

PREFIX ?= /usr/local
GOFLAGS ?= -buildvcs=false

all: clean build-server

build-client:
	go build $(GOFLAGS) -o initiative-client ./cmd/client

build-server:
	go build $(GOFLAGS) -o initiative-server ./cmd/server

clean:
	rm -f initiative-client
	rm -f initiative-server

install-client: build-client
	mkdir -p $(DESTDIR)$(PREFIX)/bin
	cp -f initiative-client $(DESTDIR)$(PREFIX)/bin
	chmod 755 $(DESTDIR)$(PREFIX)/bin/initiative-client

install-server: build-server
	mkdir -p $(DESTDIR)$(PREFIX)/bin
	cp -f initiative-server $(DESTDIR)$(PREFIX)/bin
	chmod 755 $(DESTDIR)$(PREFIX)/bin/initiative-server

uninstall:
	rm -f $(DESTDIR)$(PREFIX)/bin/initiative

.PHONY: all build-server clean install-server uninstall
