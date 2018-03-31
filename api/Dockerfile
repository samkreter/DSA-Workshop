FROM golang:1.9.2 as builder
WORKDIR /go/src/github.com/samkreter/DSA-Workshop/api
COPY . /go/src/github.com/samkreter/DSA-Workshop/api
RUN go get -u github.com/golang/dep/cmd/dep
RUN dep ensure
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o run .

FROM alpine:latest
RUN apk --update add ca-certificates
WORKDIR /root/
COPY --from=builder /go/src/github.com/samkreter/DSA-Workshop/api/run .
CMD ["./run"]