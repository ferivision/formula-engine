# Runs everything inside a golang:1.25-bookworm container, mounting
# the current directory — your laptop's local Go install (if any) is
# never used. This keeps every contributor and CI on the exact same
# pinned Go version regardless of what's installed on their machine.

GO_IMAGE := golang:1.25-bookworm
DOCKER_RUN := docker run --rm -v "$(CURDIR):/app" -w /app $(GO_IMAGE)

.PHONY: test test-race build tidy fmt fmt-check vet shell

test:
	$(DOCKER_RUN) go test -buildvcs=false ./...

test-race:
	$(DOCKER_RUN) go test -buildvcs=false -race ./...

build:
	$(DOCKER_RUN) go build -buildvcs=false ./...

tidy:
	$(DOCKER_RUN) go mod tidy

fmt:
	$(DOCKER_RUN) gofmt -l -w .

# Like fmt, but only reports -- never rewrites files. Used by CI, where
# silently reformatting and passing anyway would hide the problem.
fmt-check:
	$(DOCKER_RUN) sh -c 'unformatted=$$(gofmt -l .); if [ -n "$$unformatted" ]; then echo "$$unformatted"; echo "run '"'"'make fmt'"'"' to fix"; exit 1; fi'

vet:
	$(DOCKER_RUN) go vet ./...

# Drop into an interactive shell inside the pinned Go environment,
# useful for anything not covered by the targets above.
shell:
	docker run --rm -it -v "$(CURDIR):/app" -w /app $(GO_IMAGE) bash
