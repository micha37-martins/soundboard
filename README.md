Soundboard
===================
A Soundboard for Raspberry Pi.

- [Why?](#why)
- [How does it work?](#how-does-it-work)
- [Build](#how-to-build)
- [Recommendations](#recommendations)
- [Configuration](#config)

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
