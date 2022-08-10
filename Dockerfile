FROM balenalib/raspberry-pi-golang AS build-env
ADD . /src
RUN cd /src && apt-get update && apt-get install libusb-dev && go env -w GO111MODULE=off && go get github.com/eclipse/paho.mqtt.golang && go build -ldflags "-linkmode external -extldflags -static" -x -o mqtt-hub-ctrl .

FROM hypriot/rpi-alpine-scratch
WORKDIR /app
COPY --from=build-env /src/mqtt-hub-ctrl /app/
ENTRYPOINT ["./mqtt-hub-ctrl"]
