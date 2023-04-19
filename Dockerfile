FROM golang:1.20-alpine AS builder

RUN apk update && apk add --no-cache git

WORKDIR /go/src/app
COPY . /go/src/app

# build
RUN go get -v ./...
RUN go build -ldflags "-X github.com/noqqe/relaystation/src/relaystation.Version=`git describe --tags`" -v . 

# copy
FROM scratch
WORKDIR /go/src/app
COPY --from=builder /go/src/app/relaystation /go/src/app/relaystation

# run
CMD [ "/go/src/app/relaystation" ]
