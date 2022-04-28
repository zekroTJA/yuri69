.PHONY = run

APPNAME = yuri
CONFIGDIR = $(CURDIR)/config/private.config.toml

run:
	go run ./cmd/$(APPNAME)/main.go \
		-c $(CONFIGDIR) \
		-l 5 \
		-debug

setup-lavalink:
	docker-compose -f docker-compose.dev.yml \
		up -d lavalink
