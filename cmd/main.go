package main

import (
	"log"
	"os"
	"soundboard/internal/config"
	"soundboard/internal/filechecks"
	"soundboard/internal/pushbuttons"
)

func main() {
	log.SetOutput(os.Stdout)

	soundfiles := config.SoundfilesFolder

	// check if all files use thre correct format (audio/mpeg)
	fileErr := filechecks.CheckFiletype(soundfiles, "audio/mpeg")
	if fileErr != nil {
		log.Fatal(fileErr)
	}

	wConf := pushbuttons.NewWatcherConfig(pushbuttons.NewGpioWatcher())
	player := new(pushbuttons.RealPlay)
	pConf := pushbuttons.NewPlayerConfig(soundfiles, player)
	errChan := make(chan error, 1)

	err := pushbuttons.PushedButtons(pConf, wConf, errChan)
	if err != nil {
		log.Fatal(err)
	}
}
