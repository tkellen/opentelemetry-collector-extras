FROM golang:1.20-bullseye as builder
RUN go install go.opentelemetry.io/collector/cmd/builder@latest

WORKDIR /

COPY metricsasattributesprocessor ./metricsasattributesprocessor
COPY ocb.yaml ./ocb.yaml

RUN builder --config=ocb.yaml


FROM alpine:3.16 as certs
RUN apk --update add ca-certificates


FROM scratch

ARG USER_UID=10001
USER ${USER_UID}

COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /dist/otelcol-custom /otelcol-custom
EXPOSE 4317 4318
ENTRYPOINT ["/otelcol-custom"]
CMD ["--config", "/etc/otel/config.yaml"]
