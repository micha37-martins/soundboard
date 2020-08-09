// For playback the following beep project is used:
// (https://github.com/faiface/beep)
// It uses oto (https://github.com/hajimehoshi/oto) as dependency.

package soundboard

import (
	"log"
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

// PlaySound plays selected sound-file
// it is based on the example from gordonklaus/portaudio:
// https://github.com/gordonklaus/portaudio/blob/master/examples/mp3.go
func PlaySound(filePath string) {
	log.Println("Playing: ", filePath)

	if err := run(filePath); err != nil {
		log.Fatal(err)
	}

}

func run(filePath string) error {
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}

	streamer, format, err := mp3.Decode(f)
	log.Println("Format: ", format)
	if err != nil {
		return err
	}
	defer streamer.Close()

	err = speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	if err != nil {
		return err
	}

	done := make(chan bool)
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		done <- true
	})))

	<-done
	return nil
}
