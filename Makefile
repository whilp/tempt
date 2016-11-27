NAME := tempt
VERSION := $(git describe --long --tags)
BINARIES := $(patsubst %,release/$(NAME)-%,\
	linux-amd64 \
	darwin-amd64)
SHAS := $(BINARIES:%=%.sha256)
RELEASE := $(BINARIES) $(SHAS)

export CGO_ENABLED ?= 0

.PHONY: release

release: release.sh $(RELEASE)

%.sha256: %
	sha256sum $< > $@

$(SHAS): $(BINARIES)

$(BINARIES):
	./release.sh $(NAME) $@
