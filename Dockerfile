FROM golang:1.24.6-alpine  AS builder
RUN apk add --no-cache make

# swagger docs generation will fail if cgo is used
ENV CGO_ENABLED=0

WORKDIR /src/

# preinstall dependencies for better build caching
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY Makefile .

# preinstall code generation tools for better build caching
RUN make install-code-generation-tools

# copy api related source files and generate api docs
COPY pkg/web pkg/web
COPY pkg/models pkg/models
RUN make api-docs

# copy rest of the source files
COPY cmd cmd
COPY pkg pkg
COPY version.go .

# (re)generate mocks
COPY .mockery.yaml .mockery.yaml
RUN make generate-code

# test and build
RUN make test
RUN make build

FROM busybox
# service files
COPY --from=builder /src/api /api
COPY --from=builder /src/bin/notification-service /bin/

ENTRYPOINT ["./bin/notification-service"]
