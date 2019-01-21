# go-adc
A really simple Raspberry Pi to MCP3008 A/D sampler written in Go.

Uses go-rpio to talk to the MCP3008 A/D converter chip using the the Raspberry Pi's SPI interface.

## Usage

First, set up docker-machine to connect to your pi. I used a pi zero W, hence the comment in the Dockerfile.

Get, build and install the source on the pi, with something like this:

```docker build --tag mytag .

TODO: try cross compiling.

Run the app like this:

```docker run -it --privileged mytag adc-cli --count=1000 --interval=10ms --output /results/fast.txt

