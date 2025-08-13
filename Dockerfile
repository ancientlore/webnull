ARG GO_VERSION=1.25
ARG IMG_VERSION=1.25

FROM --platform=${BUILDPLATFORM} golang:${GO_VERSION} AS builder
WORKDIR /go/src/github.com/ancientlore/webnull
COPY . .
WORKDIR /go/src/github.com/ancientlore/webnull
RUN go version
ARG TARGETOS TARGETARCH
RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} CGO_ENABLED=0 go build -o /go/bin/webnull

FROM ancientlore/goimg:${IMG_VERSION}
COPY --from=builder /go/bin/webnull /usr/local/bin/webnull
EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/webnull", "-addr", ":8080"]
