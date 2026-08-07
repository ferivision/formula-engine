# Runs everything inside a golang:1.25-bookworm container, mounting
# the current directory — your laptop's local Go install (if any) is
# never used. This keeps every contributor and CI on the exact same
# pinned Go version regardless of what's installed on their machine.

GO_IMAGE := golang:1.25-bookworm
DOCKER_RUN := docker run --rm -v "$(CURDIR):/app" -w /app $(GO_IMAGE)

.PHONY: test test-race build tidy fmt vet shell

test:
	$(DOCKER_RUN) go test ./...

test-race:
	$(DOCKER_RUN) go test -race ./...

build:
	$(DOCKER_RUN) go build ./...

tidy:
	$(DOCKER_RUN) go mod tidy

fmt:
	$(DOCKER_RUN) gofmt -l -w .

vet:
	$(DOCKER_RUN) go vet ./...

# Drop into an interactive shell inside the pinned Go environment,
# useful for anything not covered by the targets above.
shell:
	docker run --rm -it -v "$(CURDIR):/app" -w /app $(GO_IMAGE) bash
