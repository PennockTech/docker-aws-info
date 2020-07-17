ARG BUILDER_IMAGE="golang:1.14.6"
ARG RUNTIME_BASE_IMAGE="scratch"
ARG PORT=8080

# =============================8< Builder >8==============================

FROM ${BUILDER_IMAGE} AS builder

ADD . docker-aws-info
WORKDIR docker-aws-info
RUN CGO_ENABLED=0 GOOS=linux go build -tags "docker" -ldflags -s -o /tmp/dai .

# ===========================8< Final Image >8============================

FROM ${RUNTIME_BASE_IMAGE}
ARG PORT
ENV PORT=${PORT}

COPY --from=builder /tmp/dai /docker-aws-info
# Putting an ${ARG} into CMD forces shell, there's no way to have a const number
# baked in, that I can tell.  So we switched to ENV.
CMD ["/docker-aws-info"]
EXPOSE ${PORT}

# ARG repetition because Docker doesn't let you mark one as persisting across contexts
ARG BUILDER_IMAGE
ARG RUNTIME_BASE_IMAGE
LABEL maintainer="noc+di@pennock-tech.com"
LABEL tech.pennock.builder.image="${BUILDER_IMAGE}"
LABEL tech.pennock.baseimage="${RUNTIME_BASE_IMAGE}"
LABEL tech.pennock.portlist="${PORT}"
