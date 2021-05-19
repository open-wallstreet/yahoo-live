.PHONY: compile assets

PROTOC_GEN_GO := $(GOPATH)/bin/protoc-gen-go
GO_BINDATA := $(GOPATH)/bin/go-bindata

PROTOC := $(shell which protoc)
# If protoc isn't on the path, set it to a target that's never up to date, so
# the install command always runs.
ifeq ($(PROTOC),)
    PROTOC = must-rebuild
endif

# Figure out which machine we're running on.
UNAME := $(shell uname)

$(PROTOC):
# Run the right installation command for the operating system.
ifeq ($(UNAME), Darwin)
	brew install protobuf
endif
ifeq ($(UNAME), Linux)
	sudo apt-get install protobuf-compiler
endif
# You can add instructions for other operating systems here, or use different
# branching logic as appropriate.

# If $GOPATH/bin/protoc-gen-go does not exist, we'll run this command to install
# it.
$(PROTOC_GEN_GO):
	go get -u github.com/golang/protobuf/protoc-gen-go

yaticker.pb.go: proto/yaticker.proto | $(PROTOC_GEN_GO) $(PROTOC)
	protoc --go_out=. proto/yaticker.proto

# This is a "phony" target - an alias for the above command, so "make compile"
# still works.
compile: yaticker.pb.go

$(GOPATH)/bin/proto-make-example: $(shell find . -name '*.go')
	go install .

$(GO_BINDATA):
	go get github.com/kevinburke/go-bindata