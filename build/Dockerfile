# Build the manager binary
FROM golang:1.19 as builder

WORKDIR /workspace

# Copy the Go Modules manifests
COPY go.mod go.sum ./
# cache deps before building and copying source so that we don't need to re-build as much
# and so that source changes don't invalidate our built layer
#COPY vendor/ vendor/

# Copy the project source
COPY Makefile ./
COPY mock/api/ mock/api/
COPY cmd/metering-reporter cmd/metering-reporter
COPY hack/ hack/
COPY config/ config/
# Run tests and linting
RUN make go-test

# Build
RUN make go-build

# Use distroless as minimal base image to package the application binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /workspace/bin/metering-reporter .
USER 65532:65532

ENTRYPOINT ["/metering-reporter"]
