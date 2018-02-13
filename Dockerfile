ARG BUILDER_IMAGE="golang:1.9.4"
FROM ${BUILDER_IMAGE} AS builder

ADD . /go/src/go.pennock.tech/docker-aws-info
WORKDIR /go/src/go.pennock.tech/docker-aws-info
RUN CGO_ENABLED=0 GOOS=linux go build -tags "docker" -ldflags -s .

ARG RUNTIME_BASE_IMAGE="scratch"
FROM ${RUNTIME_BASE_IMAGE}
ARG PORT=8080
ENV PORT=${PORT}

COPY --from=builder /go/src/go.pennock.tech/docker-aws-info/docker-aws-info /
# Putting an ${ARG} into CMD forces shell, there's no way to have a const number
# baked in, that I can tell.  So we switched to ENV.
CMD ["/docker-aws-info"]
EXPOSE ${PORT}

# ARG repetition because Docker doesn't let you mark one as persisting across contexts
ARG BUILDER_IMAGE="golang:1.9.4"
ARG RUNTIME_BASE_IMAGE="scratch"
LABEL maintainer="noc+di@pennock-tech.com"
LABEL tech.pennock.builder.image="${BUILDER_IMAGE}"
LABEL tech.pennock.baseimage="${RUNTIME_BASE_IMAGE}"
LABEL tech.pennock.portlist="${PORT}"
