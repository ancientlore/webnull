FROM golang:1.15 as builder
WORKDIR /go/src/github.com/ancientlore/webnull
COPY . .
WORKDIR /go/src/github.com/ancientlore/webnull
RUN go version
RUN CGO_ENABLED=0 GOOS=linux GO111MODULE=on go get .
RUN CGO_ENABLED=0 GOOS=linux GO111MODULE=on go install

FROM gcr.io/distroless/static:nonroot
COPY --from=builder /go/bin/webnull /go/bin/webnull
EXPOSE 8080
ENTRYPOINT ["/go/bin/webnull", "-addr", ":8080"]
