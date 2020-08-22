[![GitHub go.mod Go version of a Go module](https://img.shields.io/github/go-mod/go-version/gomods/athens.svg)](https://github.com/micha37-martins/soundboard)
[![GoDoc reference example](https://img.shields.io/badge/godoc-reference-blue.svg)](https://pkg.go.dev/micha37-martins/soundboard)
[![GoReportCard example](https://goreportcard.com/badge/github.com/micha37-martins/soundboard)](https://goreportcard.com/report/github.com/micha37-martins/soundboard)
[![codecov](https://codecov.io/gh/micha37-martins/soundboard/branch/master/graph/badge.svg)](https://codecov.io/gh/micha37-martins)

![pushbuttons](assets/pushbuttons.JPG)

Soundboard
===================
A Soundboard for Raspberry Pi.

- [Why?](#why)
- [How does it work?](#how-does-it-work)
- [Build](#how-to-build)
- [Docker](#build-and-run-soundboard-container)
- [Recommendations](#recommendations)
- [Configuration](#config)
- [Todo](#optional features)

Why?
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
sudo apt-get install libmpg123-dev portaudio19-dev
```

Get dependencies:
```sh
go get github.com/bobertlo/go-mpg123/mpg123
go get github.com/gordonklaus/portaudio
go get github.com/micha37-martins/gpio
```

Docker
----
Find your sound device by using
```sh
aplay -l
```

Use environment variable to set your desired sound device

examples:
```sh
ALSA_CARD=PCH
ALSA_CARD=HDMI

```

Build container
```sh
docker build -t soundboard .
```

Run container
```sh
docker run -it --rm --device /dev/snd -e "ALSA_CARD=SET_SOUND_DEVICE" test_soundboard:0.0.4 /bin/sh
```
Play test sound:
```sh
speaker-test
```

Recommendations
----
It is expected to have a [`Pull-Up-Resistor`](https://en.wikipedia.org/wiki/Pull-up_resistor)  
for every push button connected. The resistor should have between 10 and 100 k&#8486;.

Configuration
----
Configuration is don in `/internal/config/config.go`

Here the mapping between file and pin is set. The `FileMapper`  
function assigns a filename to the corresponding button. The button  
number has to be a two digit string:
`example: "01" = Button01`

The pin number for Raspberry Pi 2/3/4 can be found at  
[`raspberrypi.org`](https://www.raspberrypi.org/documentation/usage/gpio/)  

Example:
Map pin 5 to file 01

```go
&ButtonMap{
	Name: "01",
	Pin:  5,
},
```

The number of entries in `ButtonMap` have to equal the number of soundfiles.

Todo
----
- add picture of buttons
- document how to start integrationtest
  - go test -tags=integration
- write a soundboard.service example to document how to start
  soundboard using Systemd
- add troubleshoot if filecheck fails cause of empty dir
