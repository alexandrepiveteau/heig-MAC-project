FROM golang:1.15

WORKDIR /go/src/climb
COPY . .

RUN go get -d -v ./...
RUN go install -v climb/pkg

CMD ["pkg"]
