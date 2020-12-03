FROM golang:1.15
WORKDIR /go/src
COPY . .
RUN rm /go/src/Dockerfile
RUN GO111MODULE=on CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -a -installsuffix cgo -v -o api

FROM alpine:latest  
WORKDIR /root/
COPY --from=0 /go/src/api .
CMD ["./api"]
