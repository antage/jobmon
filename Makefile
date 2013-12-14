JOBMOND_DEPS := $(wildcard src/jobmon/job/*.go src/jobmon/logger/*.go src/jobmon/server/*.go)
JOBMON_DEPS := $(wildcard src/jobmon/job/*.go src/jobmon/logger/*.go src/jobmon/client/*.go)
DEB_VERSION := $(shell dpkg-parsechangelog | egrep '^Version: ' | sed -e 's/^Version: //')

all: bin/jobmond bin/jobmon

bin/jobmond: $(JOBMOND_DEPS)
	GOPATH=$(shell pwd) go build -o $@ jobmon/server

bin/jobmon: $(JOBMON_DEPS)
	GOPATH=$(shell pwd) go build -o $@ jobmon/client

.PHONY: clean deb
clean:
	rm -f bin/jobmon
	rm -f bin/jobmond

deb: clean
	rm -rf pkg
	-quilt pop -aq
	git archive HEAD | gzip > ../jobmon_$(DEB_VERSION).orig.tar.gz
	dpkg-buildpackage -sa
