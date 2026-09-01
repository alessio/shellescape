
#!/usr/bin/make -f

VERSION := $(shell git describe)

all: build

build-stamp:
	go build -a -v
	touch $@
build: build-stamp

install-stamp: build
	go install -v \
            -ldflags="X 'main.version=$(VERSION)'" \
            ./cmd/escargs
	touch $@
install: install-stamp

escargs: build
	go build -v \
            -ldflags="-X 'main.version=$(VERSION)'" \
            ./cmd/escargs

clean:
	rm -f escargs

distclean: clean
	rm -f build-stamp install-stamp

uninstall:
	rm -fv $(shell go env GOPATH)/bin/escargs

FUZZTIME ?= 30s
FUZZTARGETS := FuzzQuote FuzzQuoteCommand FuzzStripUnsafe FuzzStripSpaces FuzzScanTokens

# go test runs at most one fuzz target per invocation, so iterate over them.
# Failing inputs are written to testdata/fuzz/<target>/ and must be committed.
fuzz:
	@for t in $(FUZZTARGETS); do \
	    echo "=== $$t ==="; \
	    go test -run '^'"$$t"'$$' -fuzz '^'"$$t"'$$' -fuzztime=$(FUZZTIME) . || exit 1; \
	done

.PHONY: clean distclean fuzz install uninstall 
