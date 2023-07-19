FROM golang:1.20 as builder

ARG ARG_VERSION=1.0.0

WORKDIR /workspace

COPY Makefile Makefile

COPY go.mod go.mod
COPY go.sum go.sum

COPY main.go main.go
COPY api/ api/
COPY x-metrics/ x-metrics/
COPY hack/ hack/

COPY controllers/ controllers/

RUN make build CGO_ENABLED=0 -e VERSION="$ARG_VERSION"

RUN ls /workspace/bin

FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /workspace/bin/x-metrics .
USER nonroot:nonroot
ENTRYPOINT ["/x-metrics"]
