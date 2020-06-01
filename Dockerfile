FROM golang:alpine as builder

LABEL maintainer="Thiago Zilli Sarmento <thiago.zilli@gmail.com>"
LABEL builder="pismo-challenge"

RUN apk update && \
    apk upgrade && \
    apk add --update alpine-sdk && \
    apk add --no-cache bash git openssh make cmake

RUN mkdir /build 

ADD . /build/

WORKDIR /build 

RUN make alpine

FROM alpine:latest

RUN apk add --update sqlite curl

RUN mkdir -p app

COPY --from=builder /build/out/server.lin /app/

RUN chmod a+rw /app/ /app/* -R

WORKDIR /app

EXPOSE 8080

CMD ["./server.lin"]

HEALTHCHECK --interval=5s --timeout=2s --start-period=2s --retries=5 CMD [ "curl", "--silent", "--fail", "http://localhost:8080/ping" ]
