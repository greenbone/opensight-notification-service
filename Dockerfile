FROM golang:1.22.2-alpine  AS builder
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

# copy rest of the source files and build
COPY cmd cmd
COPY pkg pkg
RUN make test
RUN make build

FROM busybox
# service files
COPY --from=builder /src/api /api
COPY --from=builder /src/bin/notification-service /bin/

ENTRYPOINT ["./bin/notification-service"]