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
	create() *Watcher
}

// create instantiates a new watcher
func (w gpioWatcher) create() *gpio.Watcher {
	return gpio.NewWatcher()
}

// PushedButtons creates a "watcher" listening on the specified buttons
// if a button is pushed the assigned sound-file will be played by the player
func PushedButtons(folder string) error {
	numberOfButtons := len(config.ButtonMapping)
	fmt.Println("Nr. of buttons", numberOfButtons)
	// The Watcher is a type which listens on the GPIO pins you specify and then notifies you when the values of those pins change. It uses a select() call so that it does not need to actively poll, which saves CPU time and gives you better latencies from your inputs.
	// TODO write function to mock gpio
	newWatcher := gpioWatcher{}
	watcher := newWatcher.create()
	defer watcher.Close()

	// range over list of Buttons and add pin to watcher
	for i, mapping := range config.ButtonMapping {
		log.Println("Cycle nr:", i, "Add PIN:", mapping.Pin)
		// TODO mock watcher
		watcher.AddPin(mapping.Pin)
		fileName, err := filechecks.FileMapper(folder, mapping.Name)
		log.Println("Filename", fileName)
		if err != nil {
			return err
		}
	}

	// GPIO logic
	go func() {
		// build a map for buttons an name
		// Build a config map:
		confMap := map[uint]string{}
		for _, v := range config.ButtonMapping {
			confMap[v.Pin] = v.Name
		}
		for {
			var pin uint = 0
			// pullup resistor defines 1 as default value for pins
			var value uint = 1
			// TODO mock
			pin, value = watcher.Watch()
			log.Printf("read %d from gpio %d\n", value, pin)
			time.Sleep(time.Second / 2)
			if value == 0 {
				log.Println("Pressed button: ", pin)
				// And then to find values by key:
				if v, ok := confMap[pin]; ok {
					// Found
					log.Println("Corresponding name :", v)

					fN, err := filechecks.FileMapper(folder, v)

					log.Println("Filename :", fN, err)

					soundboard.PlaySound(config.SoundfilesFolder + fN)
				}

			}
		}
	}()

	// TODO replace with infinite loop
	time.Sleep(time.Second * 120)
	return nil
}
