package main

import (
	"log"
	"os"
	"soundboard/internal/config"
	fake "soundboard/internal/fake/pushbuttons"
	"soundboard/internal/filechecks"
)

func main() {
	log.SetOutput(os.Stdout)

	// check if all files use thre correct format (audio/mpeg)
	fileErr := filechecks.CheckFiletype(config.SoundfilesFolder, "audio/mpeg")
	if fileErr != nil {
		log.Fatal(fileErr)
	}
	/*

		err := pushbuttons.PushedButtons()
		if err != nil {
			log.Fatal(err)
		}
	*/
	err2 := fake.PushedButtons(config.SoundfilesFolder)
	if err2 != nil {
		log.Fatal(err2)
	}
}
