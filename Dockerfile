FROM golang:1.9.3 AS builder

ADD . /go/src/go.pennock.tech/docker-aws-info
WORKDIR /go/src/go.pennock.tech/docker-aws-info
RUN CGO_ENABLED=0 GOOS=linux go build -tags "docker" -ldflags -s .

FROM scratch
MAINTAINER Phil Pennock "phil@pennock-tech.com"

COPY --from=builder /go/src/go.pennock.tech/docker-aws-info/docker-aws-info /
CMD ["/docker-aws-info", "-port", "8080"]
EXPOSE 8080
