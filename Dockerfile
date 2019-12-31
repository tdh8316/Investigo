FROM golang:1.13-alpine AS builder

# Copy source code
COPY .  /go/src/github.com/tdh8316/Investigo
WORKDIR /go/src/github.com/tdh8316/Investigo

# Build executable
RUN cd /go/src/github.com/tdh8316/Investigo \
 	&& go build -v \
 	&& ls -l

FROM alpine:3.11 AS runtime

ARG TINI_VERSION=${TINI_VERSION:-"0.18.0"}
ARG BUILD_DATE
ARG VCS_REF

LABEL org.label-schema.build-date=$BUILD_DATE \
      org.label-schema.description="Investigo is a minimal Go implementation of Sherlock allowing to find usernames across social media." \
      org.label-schema.name="investigo" \
      org.label-schema.schema-version="1.0.0-rc1" \
      org.label-schema.usage="https://github.com/tdh8316/Investigo/blob/master/README.md" \
      org.label-schema.vcs-url="https://github.com/tdh8316/Investigo" \
      org.label-schema.vcs-ref=$VCS_REF \
      org.label-schema.vendor="x0rzkov" \
      org.label-schema.version="latest"

# Install tini to /usr/local/sbin
ADD https://github.com/krallin/tini/releases/download/v${TINI_VERSION}/tini-muslc-amd64 /usr/local/sbin/tini

# Install runtime dependencies & create runtime user
RUN apk --no-cache --no-progress add ca-certificates \
 && chmod +x /usr/local/sbin/tini && mkdir -p /opt \
 && adduser -D investigo -h /opt/investigo -s /bin/sh \
 && su investigo -c 'cd /opt/investigo; mkdir -p bin config data'

# Switch to user context
USER investigo
WORKDIR /opt/investigo

# Copy investigo binary to /opt/investigo/bin
COPY --from=builder /go/src/github.com/tdh8316/Investigo/Investigo /opt/investigo/bin/investigo
COPY config.example.yaml /opt/investigo/config/config.yaml
COPY data.json /opt/investigo/bin/data.json
ENV PATH $PATH:/opt/investigo/bin

# Container configuration
# EXPOSE 4242
VOLUME ["/opt/investigo/data"]
ENTRYPOINT ["tini", "-g", "--"]
CMD ["/opt/investigo/bin/investigo"]
