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

## Installation

There are several ways to install the `pi-monitor` onto your Raspberry PI.
I recommend using `systemd` service and use the raw binary as the bash scripts
need to access the bare kernel. Installation via docker is supported, but 
OS functionality needs to be mapped into the container.

### Service (recommended)

1. Create a new service unit file

```
  sudo vim /lib/systemd/system/pi-monitor.service
```

2. Paste the following in the newly created file
```
After=network.target

[Service]
Type=simple
ExecStart=<location of pi-monitor binary> -script-dir=<location of scripts>

[Install]
WantedBy=multi-user.target
```

3. Enable and start the service

```
  sudo systemctl enable pi-monitor.service
  sudo systemctl start pi-monitor.service
```

> The service runs (unless configured otherwise) on port 8000 and can be accessed
> by adding a reverse proxy on your web server of choice (Nginx, Caddy, Apache).
