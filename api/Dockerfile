# Build Go
FROM golang:1.22.4 as go-build

LABEL app="build-venmo-calculator"
LABEL REPO="https://github.com/leozqin/venmo-calculator"

ENV PROJPATH=/go/src/github.com/leozqin/venmo-calculator

# Because of https://github.com/docker/docker/issues/14914
ENV PATH=$PATH:$GOROOT/bin:$GOPATH/bin

ADD . /go/src/github.com/leozqin/venmo-calculator
WORKDIR /go/src/github.com/leozqin/venmo-calculator

RUN go build

# Serve
FROM debian:stable-slim
COPY --from=go-build /go/src/github.com/leozqin/venmo-calculator /opt/bin/
RUN chmod +x /opt/bin/venmo-calculator
CMD ["/opt/bin/venmo-calculator"]
EXPOSE 3333