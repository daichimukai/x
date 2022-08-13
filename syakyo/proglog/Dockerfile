FROM golang:1.18 AS build
WORKDIR /go/src/proglog
COPY . .
RUN CGO_ENABLED=0 go build -o /go/bin/proglog ./cmd/proglog
RUN CGO_ENABLED=0 go build -o /go/bin/client ./cmd/getservers
RUN GRPC_HEALTH_PROBE_VERSION=v0.4.11 && \
        curl -sSLf -o /go/bin/grpc_health_probe \
        https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 && \
        chmod +x /go/bin/grpc_health_probe

FROM alpine
COPY --from=build /go/bin/proglog /bin/proglog
COPY --from=build /go/bin/client /bin/client
COPY --from=build /go/bin/grpc_health_probe /bin/grpc_health_probe
ENTRYPOINT ["/bin/proglog"]
