NAME := tempt
VERSION := $(shell git describe --long --tags)
BINARIES := $(patsubst %,$(NAME)-%,\
	linux-amd64 \
	darwin-amd64)
SHAS := $(BINARIES:%=%.sha256)
RELEASE := $(BINARIES) $(SHAS)

.PHONY: release

release: $(RELEASE)
	hub release create $(RELEASE:%=-a %) -c master -m "$(VERSION)" "$(VERSION)"

clean:
	rm -f $(RELEASE)

%.sha256: %
	sha256sum $< > $@

$(SHAS): $(BINARIES)

$(BINARIES):
	./release.sh $(NAME) $(VERSION) $@
