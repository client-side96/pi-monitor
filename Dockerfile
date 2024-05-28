FROM balenalib/raspberry-pi-debian:buster

WORKDIR /service

RUN mkdir -p ./scripts

COPY ./bin/pi-monitor .
COPY ./scripts/* ./scripts/
COPY ./scripts/mocks/vcgencmd /usr/bin/vcgencmd

RUN chmod +x ./scripts/*.sh
RUN usermod -aG video root

EXPOSE 8000

CMD ["./pi-monitor", "-script-dir=/service/scripts/"]
