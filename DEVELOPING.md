# Developing Guide


## Build Collector
```shell
builder --config=ocb.yaml
```
if you are not in a dev container and want to debug, you may need to run this
command to rebuild with full paths
```shell
go build -C ./dist -o otelcol-custom "-gcflags=all=-N -l" -ldflags=
```

## Run Collector
The Collector uses a `HONEYCOMB_API_KEY` environment variable to send data to
Honeycomb. Make sure you have one exported in your shell before running these
commands.
```shell
./dist/otelcol-custom --config=config.yaml
```

or debug mode

```shell
dlv --listen=:2345 --headless=true --api-version=2 --accept-multiclient --log exec ./dist/otelcol-custom -- --config=config.yaml

```

If dlv is running, can attach to it with vscode debugger and hit breakpoints if ocb.yaml has debug mode enabled.

Note: ctrl + c does not stop dlv.
To stop dlv run the following command in a separate terminal
```shell
killall dlv
```

## Send a single span using otel-cli
Use this to send a single span with a host.name attribute using otel-cli
(requires otel-cli to be installed)
```shell
otel-cli span --endpoint localhost:4317 --service "test" --name "Test Span" --attrs "host.name=$(hostname)" --tp-print
```
