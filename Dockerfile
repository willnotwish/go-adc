# This is for ARM 32 v6 (eg, rpi zero w)
FROM arm32v6/golang:alpine

# FROM arm32v7/golang # (eg, rpi model 3B+)

WORKDIR /go/src/github.com/willnotwish/go-adc
COPY . .

RUN apk add --no-cache git \
    && go get -d -v ./... \
    && go install -v ./... \
    && mkdir /results \
    && apk del git

CMD ["adc-cli", "--help"]

# CMD ["ash"]
