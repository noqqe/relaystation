FROM golang:1.19

WORKDIR /go/src/app
COPY . /go/src/app

RUN go get -v ./...
RUN go build -ldflags "-X main.Version=`git describe --tags`"  -v .

# Run radsportsalat
CMD [ "./relaystation" ]
