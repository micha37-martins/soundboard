Soundboard
===================

[![GitHub go.mod Go version of a Go module](https://img.shields.io/github/go-mod/go-version/gomods/athens.svg)](https://github.com/micha37-martins/soundboard)
[![GoDoc reference example](https://img.shields.io/badge/godoc-reference-blue.svg)](https://pkg.go.dev/micha37-martins/soundboard)
[![GoReportCard example](https://goreportcard.com/badge/github.com/micha37-martins/soundboard)](https://goreportcard.com/report/github.com/micha37-martins/soundboard)
[![codecov.io](https://codecov.io/gh/micha37-martins/soundboard/branch/master/graph/badge.svg?token=:graph_token )](https://codecov.io/gh/micha37-martins/soundboard)

A Soundboard for Raspberry Pi.

![pushbuttons](assets/pushbuttons.JPG)

Contents
----
- [What is the purpose of this tool?](#what-is-the-purpose-of-this-tool?)
- [How does it work?](#how-does-it-work?)
- [Build](#build)
- [Prerequisites](#prerequisites)
- [Build and run soundboard container](#build-and-run-soundboard-container)
- [Recommendations](#recommendations)
- [Configuration](#configuration)
- [Testing](#testing)
- [Troubleshoot](#troubleshoot)
- [Todo](#todo)

What is the purpose of this tool?
----
This tool enables your Raspberry Pi to play a specific sound-file by pressing a  
button connected via GPIO.

How does it work?
----
Depending on the number of buttons to use, the soundboard creates listeners  
for the GPIO pins. If a button is pushed the associated sound-file will be played.  
During playback no inputs from push buttons will be accepted.

Build
----
Linux (Raspbian)

Install needed Debian packages:
```sh
sudo apt-get install build-essential libasound2-dev alsa-utils
```

Get dependencies:
```sh
go get -u github.com/hajimehoshi/oto
go get -u github.com/faiface/beep
go get -u github.com/micha37-martins/gpio
```

Prerequisites
----

### Hardware

- Raspberry PI 3 (other versions should also work but I did not test them)
- SD Card
- Power supply
- Wiring
- Pushbuttons
- A suitable case for a soundboard
- An output device like a speaker

### Recommendations

It is expected to have a [`Pull-Up-Resistor`](https://en.wikipedia.org/wiki/Pull-up_resistor)  
for every push button connected. The resistor should have between 10 and 100 k&#8486;.

### Software

- I tested with Raspbian but other operating systems should also work
- golang 1.14+
- `make`

Older software versions may also work, but I did not test that.

### Soundfiles

- .mp3 files have to be stored in ./soundfiles folder
- the .gitignore file in this folder have to be removed

Build and run soundboard container
----

### Install Docker

To install Docker on a raspberry pi you basically need:
```sh
curl -sSL https://get.docker.com | sh
```
And after the command finished you most likely want to run:
```sh
sudo usermod -aG docker pi
```
Source: [Docker install guide](https://dev.to/rohansawant/installing-docker-and-docker-compose-on-the-raspberry-pi-in-5-simple-steps-3mgl "Install guide")

### Use Docker

Find your sound device by using
```sh
aplay -l
```

Use environment variable to set your desired sound device

examples:
```sh
ALSA_CARD=PCH
ALSA_CARD=HDMI
ALSA_CARD=0
```

Build container
```sh
docker build -t soundboard .
```

Run container
```sh
docker run -it --rm --device /dev/snd -e "ALSA_CARD=SET_YOUR_SOUND_DEVICE" -v /sys:/sys soundboard:0.0.1 /bin/sh
```
Play test sound:
```sh
speaker-test
```

Example:
```sh
docker run -it --rm --device /dev/snd -e "ALSA_CARD=0" -v /sys:/sys soundboard:0.0.1 /bin/sh
```

### Use Docker-Compose
You can also use docker-compose to run the soundboard container. The sound_device 
you discovered using `aplay -l` can be configured in the `.env` file. The default 
is "PCH".

```sh
docker-compose up --build
```

The `--build` flag is needed at first start only.

Configuration
----
Configuration is done in `/internal/config/config.go`

Here the mapping between file and pin is set. The `FileMapper`  
function assigns a filename to the corresponding button. The button  
number has to be a two digit string:
`For example: "01" for your first button`

The pin number for Raspberry Pi 2/3/4 can be found at  
[`raspberrypi.org`](https://www.raspberrypi.org/documentation/usage/gpio/)  

Example:
Map pin 5 to file 01

```go
ButtonMap[5] = "01"
```

The number of entries in `ButtonMap` have to equal the number of soundfiles.

Testing
----
A usual go test is configured in the `Makefile` but it will not do an integration test. If you want to check if the soundboard plays a file use:
```sh
go test -tags=integration
```

Troubleshoot
----
- "No files in folder: ./soundfiles/"
  - Remove the hidden `.gitignore` file from the `./soundfiles` folder
- "failed to open gpio 13 direction file for writing"
  - use sudo to run soundboard
  - e.g. `sudo /usr/local/go/bin/go run cmd/main.go`
- Cannot play sound when using a Docker container
  - mount the `/sys` directory to have acess to gpio file
  - select the correct sound_card

Todo
----
- write a soundboard.service example to document how to start
  soundboard using Systemd
