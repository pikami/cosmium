FROM scratch

WORKDIR /app
COPY cosmium /app/cosmium

ENTRYPOINT ["/app/cosmium"]
