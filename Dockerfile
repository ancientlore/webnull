FROM golang:1.19 as builder
WORKDIR /go/src/github.com/ancientlore/webnull
COPY . .
WORKDIR /go/src/github.com/ancientlore/webnull
RUN go version
RUN CGO_ENABLED=0 GOOS=linux GO111MODULE=on go get .
RUN CGO_ENABLED=0 GOOS=linux GO111MODULE=on go install

FROM ancientlore/goimg:1.19
COPY --from=builder /go/bin/webnull /usr/bin/webnull
EXPOSE 8080
ENTRYPOINT ["/usr/bin/webnull", "-addr", ":8080"]
