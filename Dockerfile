FROM golang:latest as builder
WORKDIR /app
COPY . .

RUN GOOS=linux CGO_ENABLED=0 go build -ldflags="-w -s" -o build cmd/main.go


FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/build .
COPY --from=builder /app/.env .

RUN apk --no-cache add ca-certificates

EXPOSE 8080

CMD ["./build"]