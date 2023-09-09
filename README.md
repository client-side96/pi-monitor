# PI Monitor

pi-monitor` is a software that provides real-time system information about your
Raspberry PI. The application monitors system data like for example CPU
temperature, CPU load and memory usage.

## Local Development

The application is written in `golang` for a compile target with a `arm-v5`
architecture. The local development is possible with a docker container that
simulates the underlying Raspberry PI.

### Building and starting docker image

```sh
# Build golang binary and docker container
make image

# Start the container
make start
```

### Unit testing

Component tests can be run with:

```sh
make test
```
