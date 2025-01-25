FROM ghcr.io/cosmiumdev/cosmos-explorer-base:latest AS explorer-base
FROM alpine:latest

COPY --from=explorer-base /cosmos-explorer /cosmos-explorer

WORKDIR /app
COPY cosmium /app/cosmium

ENTRYPOINT ["/app/cosmium", "-ExplorerDir", "/cosmos-explorer"]
