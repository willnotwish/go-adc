# go-adc
A really simple sampling routine to grab voltages from a MCP3008 ADC connected to a Raspberry Pi. Written in Go.

Uses [go-rpio](https://github.com/stianeikeland/go-rpio) to talk to the [MCP3008 A/D converter chip](http://ww1.microchip.com/downloads/en/DeviceDoc/21295d.pdf) using the the Raspberry Pi's SPI interface.

## Usage

First, set up docker-machine to connect to your pi. I used a pi zero W, hence the comment in the Dockerfile.

Get, build and install the source on the pi, with something like this

```
docker build --tag go-adc .
```

Run the app like this:

```
docker run -it --privileged go-adc adc-cli --count=1000 --interval=10ms --output=/results/fast.txt
```

Don't forget the --privileged flag. If you see a dev/mem error then the chances are you forgot to include it. Your Docker container needs access to /dev/mem in order to use memory-mapped io.

## To do

TODO: try cross compiling.
