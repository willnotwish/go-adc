# This is for ARM 32 v6 (eg, rpi zero w)
FROM arm32v6/node:alpine as builder

# Install build toolchain, install node deps and compile native add-ons
RUN apk update && \
  apk add --no-cache make gcc g++ python linux-headers udev

RUN npm install serialport --build-from-source --unsafe-perm

# Copy built node modules and binaries without including the toolchain
FROM arm32v6/node:alpine as app
COPY --from=builder node_modules .

COPY . .

CMD ["node", "--version"]