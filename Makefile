GOBIN ?= go
GOPATH = $(shell pwd)/lib

GCHATD_SRC = ./src
GCHATD_BIN = ./bin
GCHATD_SOURCES = $(GCHATD_SRC)
GCHATD_VERSION ?= 0.0.1
GCHATD_APPS = $(GCHATD_BIN)/gchatd

all: $(GCHATD_APPS)

#$(GOPATH)/src/%:
#    GOPATH=$(GOPATH) $(GOBIN) get $*

$(GCHATD_BIN)/gchatd: deps
	GOPATH=$(GOPATH) $(GOBIN) build -o gchatd $(GCHATD_SOURCES)
	mv gchatd $(GCHATD_BIN)

deps:
	[ -d lib ] || mkdir lib
	[ -d bin ] || mkdir bin
	[ -d logs ] || mkdir logs

clean:
	rm -rf lib bin logs
