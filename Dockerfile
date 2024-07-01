# Build Stage
FROM leozqin/venmo-calculator:1.13 AS build-stage

LABEL app="build-venmo-calculator"
LABEL REPO="https://github.com/leozqin/venmo-calculator"

ENV PROJPATH=/go/src/github.com/leozqin/venmo-calculator

# Because of https://github.com/docker/docker/issues/14914
ENV PATH=$PATH:$GOROOT/bin:$GOPATH/bin

ADD . /go/src/github.com/leozqin/venmo-calculator
WORKDIR /go/src/github.com/leozqin/venmo-calculator

RUN make build-alpine

# Final Stage
FROM leozqin/venmo-calculator

ARG GIT_COMMIT
ARG VERSION
LABEL REPO="https://github.com/leozqin/venmo-calculator"
LABEL GIT_COMMIT=$GIT_COMMIT
LABEL VERSION=$VERSION

# Because of https://github.com/docker/docker/issues/14914
ENV PATH=$PATH:/opt/venmo-calculator/bin

WORKDIR /opt/venmo-calculator/bin

COPY --from=build-stage /go/src/github.com/leozqin/venmo-calculator/bin/venmo-calculator /opt/venmo-calculator/bin/
RUN chmod +x /opt/venmo-calculator/bin/venmo-calculator

# Create appuser
RUN adduser -D -g '' venmo-calculator
USER venmo-calculator

ENTRYPOINT ["/usr/bin/dumb-init", "--"]

CMD ["/opt/venmo-calculator/bin/venmo-calculator"]
