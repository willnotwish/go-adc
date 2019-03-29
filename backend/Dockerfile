# This is for ARM 32 v6 (eg, rpi zero w)
FROM arm32v6/alpine

# Need CA certs in order to make requests using https
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

WORKDIR /app
COPY dist/arm32v6/. .

RUN mkdir /app/results

CMD ["/app/go-adc"]
