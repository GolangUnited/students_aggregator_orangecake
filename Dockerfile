ARG APP_NAME="UrlScrapper"

FROM golang:1.19 as builder

WORKDIR /build
COPY . .
RUN go get -v ./... && GOOS=linux go build -o /build/${APP_NAME}/ ./...

FROM alpine:latest
COPY --from=builder /build/${APP_NAME}/ /bin/${APP_NAME}/
ENTRYPOINT [ "/bin/${APP_NAME}" ]
