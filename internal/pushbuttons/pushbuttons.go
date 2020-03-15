// pushbuttons processes the GPIO logic and triggers playback

package pushbuttons

import (
	"fmt"
	"log"
	"soundboard/internal/config"
	"soundboard/internal/filechecks"
	"soundboard/internal/soundboard"
	"time"

	"github.com/micha37-martins/gpio"
)

type gpioWatcher struct{}

type gpioInterface interface {
	create() *gpio.Watcher
}

// create instantiates a new watcher
// The Watcher is a type which listens on the GPIO pins you specify and
// then notifies you when the values of those pins change.
func (w gpioWatcher) create() *gpio.Watcher {
	return gpio.NewWatcher()
}

func watchPins(watcher *gpio.Watcher, folder string) {

	// build a map for buttons and name
	confMap := map[uint]string{}
	for _, v := range config.ButtonMapping {
		confMap[v.Pin] = v.Name
	}

	for {
		var pin uint = 0
		// pullup resistor defines 1 as default value for pins
		var value uint = 1
		pin, value = watcher.Watch()
		log.Printf("read %d from gpio %d\n", value, pin)
		time.Sleep(time.Second / 2)
		if value == 0 {
			log.Println("Pressed button: ", pin)
			// And then to find values by key:
			if v, ok := confMap[pin]; ok {
				// Found
				log.Println("Corresponding name :", v)

				fileName, err := filechecks.FileMapper(folder, v)

				log.Println("Filename :", fileName, err)

				soundboard.PlaySound(config.SoundfilesFolder + fileName)
			}

		}
	}
}

// PushedButtons creates a "watcher" listening on the specified buttons
// if a button is pushed the assigned sound-file will be played
func PushedButtons(folder string) error {
	numberOfButtons := len(config.ButtonMapping)
	fmt.Println("Nr. of buttons", numberOfButtons)
	gpioWatcher := gpioWatcher{}
	watcher := gpioWatcher.create()
	defer watcher.Close()

	// range over list of Buttons and add pin to watcher
	for i, mapping := range config.ButtonMapping {
		log.Println("Cycle nr:", i, "Add PIN:", mapping.Pin)

		watcher.AddPin(mapping.Pin)
	}
	go watchPins(watcher, folder)

	// TODO replace with infinite loop
	time.Sleep(time.Second * 120)
	return nil
}
