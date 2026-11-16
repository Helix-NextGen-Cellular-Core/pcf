# Builds the Helix image ghcr.io/helix-nextgen-cellular-core/pcf:v1.0.0
# the same way it was built and verified:
#   docker build -t pcf:local .
FROM golang:1.16.0-stretch AS builder
WORKDIR /go/src/pcf
COPY . .
RUN go mod vendor \
 && make all && mkdir -p /out/free5gc/pcf && cp bin/* /out/free5gc/pcf/

FROM alpine:3.16
LABEL org.opencontainers.image.source="https://github.com/Helix-NextGen-Cellular-Core/pcf" \
      org.opencontainers.image.licenses="Apache-2.0"
RUN apk update && apk add -U vim strace net-tools curl netcat-openbsd bind-tools
COPY --from=builder /out/ /
WORKDIR /free5gc/pcf
