FROM golang:1.9.3 AS builder

ADD . /go/src/go.pennock.tech/docker-aws-info
WORKDIR /go/src/go.pennock.tech/docker-aws-info
RUN go build .

FROM scratch
MAINTAINER Phil Pennock "phil@pennock-tech.com"
EXPOSE 8080

COPY --from=builder /go/src/go.pennock.tech/docker-aws-info/docker-aws-info /
ENV PORT=8080
CMD ["/docker-aws-info", "-port", "8080]
