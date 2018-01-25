FROM golang:1.9.2 as builder
WORKDIR /go/src/github.com/samkreter/DSA-Workshop/cmd/postgres/ 
COPY . /go/src/github.com/samkreter/DSA-Workshop/cmd/postgres/
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o run .

FROM scratch
WORKDIR /root/
COPY --from=builder /go/src/github.com/samkreter/DSA-Workshop/cmd/postgres/run .
CMD ["./run"]