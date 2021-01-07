TEST?=$$(go list ./... | grep -v 'vendor')
HOSTNAME=registry.terraform.io
NAMESPACE=hkak03key
MKFILE_PATH=$(abspath $(lastword $(MAKEFILE_LIST)))
BINARY=$(notdir $(patsubst %/,%,$(dir $(MKFILE_PATH))))
NAME=$(shell echo $(BINARY) | sed "s/terraform-provider-//")
VERSION=0.0.1
OS_ARCH=linux_amd64

default: testacc

.PHONY: build
build:
	go build -o ${BINARY}

.PHONY: install
install: build
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	mv ${BINARY} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}

.PHONY: test
test: 
	go test -i $(TEST) || exit 1                                                   
	echo $(TEST) | xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4                    

# Run acceptance tests
.PHONY: testacc
testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m
