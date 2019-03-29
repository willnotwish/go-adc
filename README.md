# go-adc
Grabs voltages from a MCP3008 ADC connected to a Raspberry Pi. The back end is in Go, the front end in Vue.js.

The back end uses the excellent [go-rpio](https://github.com/stianeikeland/go-rpio) to talk to an [MCP3008 A/D converter chip](http://ww1.microchip.com/downloads/en/DeviceDoc/21295d.pdf) connected via the Raspberry Pi's SPI. There are many tutorials which show you the hardware arrangement: it's pretty standard. Most of them gloss over the (mostly Python) software though. I wanted to do this in Go, not C.

The whole app is orchestrated using docker compose.

The back end reminds me of my early days of C coroutines and cooperative multitasking. The front end uses Vuex, websockets and Bulma css.

## Why?

I needed to monitor the battery voltage on my old Mercedes.

I wanted to get a feel for what the voltage was doing when pulling up at traffic lights: I could hear a change in the whine of the fuel pump and sometimes see the headlights going dim. At the same time I wanted to monitor the long-term state of the battery, which needed charging every few months during the winter.

I started by attaching my multimeter to the 12V supply in the car. But I couldn't watch the meter while driving and it didn't record any activity for me to "replay" later. I needed some form of logger.

## Hardware

The basic arrangement is shown in Figure 1. The battery voltage in the car is dropped to a safe level with a simple voltage divider. It is then buffered with a unity gain op amp, before being fed into the ADC3008 A/D converter. The sampled voltages are read periodically (*i.e.*, **logged**) by the Raspberry Pi via its SPI (serial peripheral interface).

I originally started with a Raspberry Pi 3 Model B. I fried it when I accidentally shorted the 3V & 5V pins together. A £35 mistake. After that I took the cheaper route of a Zero W. I added some overvoltage protection in the form of another resistor and diode D1 to dump overvoltages into the 5V supply rail and limit the input into the ADC.

## Usage

First, set up docker-machine to connect to your pi. I used a pi zero W. This uses ARMv6, not v7 as used by the pi 3.

Get, build and install the source on the pi, with something like this

```
docker build --tag go-adc .
```

Run the app (from the command line at first - to check it's working) like this:

```
docker run -it --privileged go-adc adc-cli --count=1000 --interval=10ms --output=/results/fast.txt
```

Don't forget the --privileged flag. If you see a /dev/mem error then the chances are you forgot to include it. The Docker container needs access to /dev/mem in order to use memory-mapped io. See [go-rpio](https://github.com/stianeikeland/go-rpio) for details.

### As an IoT device

Run up the included web server on the pi

```
docker run -it -p 3000:3000 go-adc httpd
```

and then browse to your pi's IP address on port 3000. Use docker-machine like this

```
docker-machine ls
```

if you've forgotten it. You should see an index page saying it's working.

### Cross compiling
In the end I wrote an old-school Makefile to do the dirty work. I had to re-learn some old stuff. But it works OK.

### Using the rpi's serial port to access external peripherals
I purchased an ELM327L directly from the manufacturer in Canada. I avoided the temptation to buy a much cheaper Chinese clone. It runs directly from a 3.3V supply, so it should connect directly to the pi.

When I started to look at the pi's serial port offerings, I got confused. Here's my take on it, from https://www.raspberrypi.org/documentation/configuration/uart.md

There are two UARTs: a PL011 (aka "full") and a cut down "mini". The full is (by default) used for Bluetooth on the zero W. To use the mini, you need to do two things:

1. Reduce the pi's core_freq on of the GPU to 250MHz from its default of 400MHz. Do this by adding ```enable_uart=1``` to ```config.txt```. I'm not sure of the side effects, if any. Something will be slower, I presume. We'll have to wait and see.
2. Tell Raspian not to use the mini for Linux console access. You can do this by editing ```/boot/cmdline.txt``` directly, removing the text ```console=serial0,115200``` (but leaving the rest of the line intact).

I was able to do both these things with ```raspi-config```. Select options 5 & 6. After a reboot, I could see ```/dev/ttyS0``` and its symlink ```/dev/serial0```.

The mini's transmit and receive pins are GPIO pins 14 & 15, which translates to pins 8 & 10 on the GPIO header.

*Note.* When I accessed the pi via a Docker container, I was unable to see the ```/dev/serial0``` symlink, but I could access ```/dev/ttyS0``` OK using minicom. I shorted pins 14 & 15 together, and I could see the output transmitted back to the input.

Hopefully, by using the mini UART for ELM327 comms, I can still use Bluetooth to control the pi (should I develop this idea further) in real time.

### Extending the basic A/D sampler as a more functional data (voltage) logger
Suppose I want to be able to access the sampler via a custom designed IoT web interface.

Here's my idea.

1. Serve a single page application, written in Vue.js via Nginx.
2. The spa makes API requests to a back end Go server, which interacts with the MCP3008.

The user interacts with the spa to view the latest voltage samples, to control the MCP3008 sampling rate and select the input channel(s) to use.

### Workflow
#### Back end
1. Cross compile the back end Go code for the pi0w in a dedicated Docker container, producing the ```go-adc``` executable in the ```backend/dist``` directory.
2. Build a Docker image based on ```arm32v6/alpine``` tagged ```dl-backend-arm32v6``` from the Go executable, using a custom Dockerfile ```backend/Dockerfile``` and including the back end executable(s).
3. Push the go-adc-arm32v6 image to Docker Hub.
4. On the pi0w (using Docker Machine), pull the go-adc-arm32v6 image.
5. Run the back end with
docker run --rm -it --privileged -p 8001:8001 -d willnotwish/go-adc-arm32v6:latest
6. Monitor log with docker logs -f <id>

#### Front end
1. Use ```yarn run build``` to build the Vue app in the ```frontend/dist``` directory.
2. Write an Nginx config file to serve the front end spa, fwarding requests beginning /api to the Go back end.
3. Build a Docker image based on ```arm32v6/nginx``` tagged ```dl-frontend-arm32v6``` including files from the frontend dist directory.

### Deployment
Run docker-compose to orchestrate the two containers.





Use vue-cli to generate the SPA. Develop this independently of the back end API.

Build the SPA
