FROM resin/raspberry-pi-golang AS build-env
ADD . /src
RUN cd /src && apt-get update && apt-get install libusb-dev && go get -u github.com/eclipse/paho.mqtt.golang && go build -x -o mqtt-hub-ctrl .

FROM hypriot/rpi-alpine-scratch
WORKDIR /app
COPY --from=build-env /src/mqtt-hub-ctrl /app/
ENTRYPOINT ["./mqtt-hub-ctrl"]
