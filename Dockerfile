# build
FROM golang:1.16-alpine AS builder

WORKDIR /src/upload-example

RUN apk add build-base git

COPY . .

RUN make build

# deployment
FROM alpine:latest
WORKDIR /home/example
COPY --from=builder /src/upload-example/upload-example ./

EXPOSE 3000

ENTRYPOINT ["./upload-example"]
