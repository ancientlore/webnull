FROM golang as builder
WORKDIR /go/src/github.com/ancientlore/webnull
ADD . .
WORKDIR /go/src/github.com/ancientlore/webnull
RUN CGO_ENABLED=0 GOOS=linux GO111MODULE=on go get .
RUN CGO_ENABLED=0 GOOS=linux GO111MODULE=on go install

FROM gcr.io/distroless/base
WORKDIR /go
COPY --from=builder /go/bin/webnull /go/bin/webnull

ENTRYPOINT ["/go/bin/webnull", "-addr", ":8080"]

EXPOSE 8080
