GO ?= go
output ?= dbg-go
TEST_FLAGS ?= -v -race -tags=test_all
DRIVERKIT_VERSION=v0.15.0
LDFLAGS := -X github.com/falcosecurity/driverkit/pkg/driverbuilder/builder.defaultImageTag=${DRIVERKIT_VERSION}

.PHONY: build
build: ${output}

.PHONY: ${output}
${output}:
	CGO_ENABLED=0 GOEXPERIMENT=loopvar $(GO) build -ldflags '${LDFLAGS}' -o $@

.PHONY: clean
clean:
	$(RM) -R ${output}

.PHONY: test
test:
	GOEXPERIMENT=loopvar $(GO) test ${TEST_FLAGS} ./...

.PHONY: bump-driverkit
bump-driverkit:
	go get github.com/falcosecurity/driverkit@$(DRIVERKIT_VER)
	sed -E "s/DRIVERKIT_VERSION=${DRIVERKIT_VERSION}/DRIVERKIT_VERSION=${DRIVERKIT_VER}/" Makefile
