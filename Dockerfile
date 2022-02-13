FROM golang:1.17 AS build

WORKDIR /build
COPY . ./
RUN make clean
RUN make build

FROM golang:1.17

RUN mkdir -p /opt/lakitu
COPY --from=build /build/bin/linux_amd64/lakitu /opt/lakitu/lakitu

ENTRYPOINT ["/opt/lakitu/lakitu", "--logging", "--bind", "0.0.0.0", "--port", "9090", "--assets", "/data"]
