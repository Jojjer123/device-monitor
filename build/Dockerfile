FROM onosproject/golang-build:v1 as build

ENV ONOS_MONITOR_ROOT=/go/src/github.com/onosproject/device-monitor
ENV CGO_ENABLED=0

RUN mkdir -p $ONOS_MONITOR_ROOT/

COPY . $ONOS_MONITOR_ROOT

RUN cd $ONOS_MONITOR_ROOT && GO111MODULE=on go build -o /go/bin/main ./cmd/


FROM alpine:3.11
RUN apk add bash
ENV HOME=/home/device-monitor
RUN mkdir $HOME
WORKDIR $HOME

COPY --from=build /go/bin/main /usr/local/bin/

CMD ["main"]