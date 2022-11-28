FROM golang:1.19 as builder

WORKDIR /build
COPY . .
# CGO_ENABLED=0 for alpine (will not work without)
RUN go get -v ./... && GOOS=linux CGO_ENABLED=0 go build -o /myApp/ ./...


FROM alpine:latest

COPY --from=builder /myApp/ /myApp/
ENTRYPOINT [ "myApp/aggregator"]

