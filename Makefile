JOBMOND_DEPS = $(wildcard src/jobmon/job/*.go src/jobmon/logger/*.go src/jobmon/server/*.go)
JOBMON_DEPS = $(wildcard src/jobmon/job/*.go src/jobmon/logger/*.go src/jobmon/client/*.go)

all: bin/jobmond bin/jobmon

bin/jobmond: $(JOBMOND_DEPS)
	GOPATH=$(shell pwd) go build -o $@ jobmon/server

bin/jobmon: $(JOBMON_DEPS)
	GOPATH=$(shell pwd) go build -o $@ jobmon/client

.PHONY: clean
clean:
	rm -f bin/jobmon
	rm -f bin/jobmond
