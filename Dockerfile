FROM alpine:latest

WORKDIR /app
COPY cosmium /app/cosmium

ENTRYPOINT ["/app/cosmium"]
