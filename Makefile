.PHONY: build test clean docker

GO = CGO_ENABLED=0 GO111MODULE=on go

MICROSERVICES=cmd/device-gps

.PHONY: $(MICROSERVICES)

DOCKERS=docker_device_gps
.PHONY: $(DOCKERS)

VERSION=$(shell cat ./VERSION)
GIT_SHA=$(shell git rev-parse HEAD)
GOFLAGS=-ldflags "-X github.com/edgexfoundry/device-gps.Version=$(VERSION)"

build: $(MICROSERVICES)
	$(GO) build ./...

cmd/device-gps:
	$(GO) build $(GOFLAGS) -o $@ ./cmd

test:
	$(GO) test ./... -cover

clean:
	rm -f $(MICROSERVICES)

docker: $(DOCKERS)

docker_device_gps:
	docker build \
		--label "git_sha=$(GIT_SHA)" \
		-t edgexfoundry/docker-device-gps:$(GIT_SHA) \
		-t edgexfoundry/docker-device-gps:$(VERSION)-dev \
		.