.PHONY: build test clean docker

GO=CGO_ENABLED=0 GO111MODULE=on go

MICROSERVICES=cmd/device-gps-go/device-gps-go
.PHONY: $(MICROSERVICES)

VERSION=$(shell cat ./VERSION)

GOFLAGS=-ldflags "-X github.com/edgexfoundry/device-gps-go.Version=$(VERSION)"

GIT_SHA=$(shell git rev-parse HEAD)

build: $(MICROSERVICES)
	$(GO) install -tags=safe

cmd/device-gps-go/device-gps-go:
	$(GO) build $(GOFLAGS) -o $@ ./cmd/device-gps-go

docker:
	docker build \
		-f cmd/device-gps-go/Dockerfile \
		--label "git_sha=$(GIT_SHA)" \
		-t edgexfoundry/docker-device-sdk-simple:$(GIT_SHA) \
		-t edgexfoundry/docker-device-sdk-simple:$(VERSION)-dev \
		.

test:
	$(GO) test ./... -cover

clean:
	rm -f $(MICROSERVICES)
