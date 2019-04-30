FROM golang:1.10.8-alpine3.8 as builder

RUN apk update && apk add --no-cache --virtual .build-deps wget
RUN wget https://github.com/squizzi/s3-dtr-tester/blob/master/s3Test.go && \
    mv s3Test.go /go/

RUN go get github.com/aws/aws-sdk-go/aws
RUN go get github.com/aws/aws-sdk-go/service/s3

RUN cd /go && \
    go build s3Test.go

FROM alpine:3.8
COPY --from=builder /go/s3Test.go /usr/bin/s3Test
ENTRYPOINT ["s3Test"]
