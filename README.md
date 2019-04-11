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

The basic arrangement is shown in Figure 1. 

<TDOO> Figure 1 here

The battery voltage in the car is dropped to a safe level with a simple voltage divider. It is then buffered with a unity gain op amp, before being fed into the ADC3008 A/D converter. The sampled voltages are read periodically (*i.e.*, **logged**) by the Raspberry Pi via its SPI (serial peripheral interface).

I originally started with a Raspberry Pi 3 Model B. I fried it when I accidentally shorted the 3V & 5V pins together. A £35 mistake. After that I took the cheaper route of a Zero W. I added some overvoltage protection in the form of another resistor and diode D1 to dump overvoltages into the 5V supply rail and limit the input into the ADC.

### Workflow

#### Back end
1. Cross compile the back end Go code for the pi0w in a dedicated Docker container, producing the ```go-adc``` executable in the ```backend/dist/arm32v6``` directory.
2. Build a Docker image based on ```arm32v6/alpine``` tagged ```dl-backend``` from the Go executable, using a the custom Dockerfile ```backend/Dockerfile``` and including the back end executable(s).
3. Push the dl-backend image to Docker Hub.
4. On the pi0w (using Docker Machine), pull the dl-backend image.

#### Front end
1. Use ```yarn run build``` to build the Vue app in the ```frontend/dist``` directory.
2. Write an Nginx config file to serve the front end spa, fwarding requests beginning /api to the Go back end.
3. Build a Docker image based on ```arm32v6/nginx``` tagged ```dl-frontend``` including files from the frontend dist directory.

### In production
Run docker-compose to orchestrate the two containers.

Don't forget the --privileged (or equivalent) flag. If you see a /dev/mem error then the chances are you forgot to include it. The Docker container needs access to /dev/mem in order to use memory-mapped io. See [go-rpio](https://github.com/stianeikeland/go-rpio) for details.

### Cross compiling
In the end I wrote an old-school Makefile to do the dirty work. I had to re-learn some old stuff. But it works OK.
