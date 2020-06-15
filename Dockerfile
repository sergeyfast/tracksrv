FROM golang:alpine AS builder

RUN apk add --no-cache git

COPY . .
RUN go build -o /tracksrv .

# Start fresh from a smaller image
FROM alpine:3.9
RUN apk add ca-certificates

COPY --from=builder /tracksrv /tracksrv

EXPOSE 8090

ENTRYPOINT ["/tracksrv"]