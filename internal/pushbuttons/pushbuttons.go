// pushbuttons processes the GPIO logic and triggers playback

package pushbuttons

import (
	"errors"
	"fmt"
	"log"
	"soundboard/internal/config"
	"soundboard/internal/filechecks"
	"soundboard/internal/soundboard"
	"time"

	"github.com/micha37-martins/gpio"
)

// Initialize ButtonMap
var buttonMap = config.InitializeButtonMap()

// realPlay is used to implement Player interface in a non test situation
type RealPlay struct{}

type watcherConfig struct {
	watcher Watcher
}

type playerConfig struct {
	folder string
	player Player
}

// Player defines interface for music player
type Player interface {
	play(string)
}

// Watcher defines interface for gpio watcher
type Watcher interface {
	Watch() (uint, uint)
	AddPin(uint)
	Close()
}

type PinWatcher interface {
	watchPins(wConf *watcherConfig, errChan chan error)
}

// NewGpioWatcher creates a gpio-Watcher for checking gpio pins
func NewGpioWatcher() Watcher {
	return gpio.NewWatcher()
}

// NewWatcherConfig constructor for instanciating watcherConfig
func NewWatcherConfig(watcher Watcher) *watcherConfig {
	return &watcherConfig{watcher: watcher}
}

// NewPlayerConfig constructor for instanciating watchConfig
func NewPlayerConfig(folder string, player Player) *playerConfig {
	pc := new(playerConfig)
	pc.folder = folder
	pc.player = player
	return pc
}

// play calls PlaySound function to start playing file
func (p *RealPlay) play(path string) {
	log.Println("folder + fileName =", path)

	soundboard.PlaySound(path)
}

// mapAndPlay finds values by key and hands path to player
func (pConf *playerConfig) mapAndPlay(errChan chan error, pin uint) {
	if name, ok := buttonMap[pin]; ok {
		fileName, err := filechecks.FileMapper(pConf.folder, name)

		if err != nil {
			errChan <- err
		}

		log.Printf("Folder: %s - Filename: %s", pConf.folder, fileName)

		pConf.player.play(pConf.folder + fileName)
	} else {
		errChan <- errors.New("No file mapped to pushed button.")
	}
}

// checkPins checks if a button has been pushed
// pullup resistor has "1" set as default value for pins
// notice: zero-value of uint is 0
func (pConf *playerConfig) checkPins(watcher Watcher, errChan chan error) (uint, uint) {
	var pin uint = 0
	var value uint = 1

	// To avoiud linting error ineffectual assignment logging var values
	log.Println("Pin/Value: ", pin, value)

	pin, value = watcher.Watch()

	if value != 0 && value != 1 {
		e := fmt.Sprintf("Pin value ot of scope. Only 0 or 1 allowed, got %d.", value)
		errChan <- errors.New(e)
		return 0, 0
	}

	log.Printf("Got value %d and pin: %d from watcher.Watch()\n",
		value, pin)

	return pin, value
}

// watchPins continously calls checkPins to test if a button has been pushed
func (pConf *playerConfig) watchPins(wConf *watcherConfig, errChan chan error) {
	log.Println("ButtonMap: ", buttonMap)

	for {
		// slow down loop to go easy on resources
		time.Sleep(time.Second / 2)

		pin, value := pConf.checkPins(wConf.watcher, errChan)

		if value == 0 {
			pConf.mapAndPlay(errChan, pin)
		}
	}
}

// PushedButtons assigns a pin his corresponding name and adds it
// to the watcher. If a button is pressed the assigned sound-file
// will be played
func PushedButtons(pConf PinWatcher, wConf *watcherConfig, errChan chan error) error {
	// Watcher is a type which listens on the GPIO pins you specify and
	// then notifies you when the values of those pins change.
	defer wConf.watcher.Close()

	// range over list of Buttons and add pin to watcher
	for pin, name := range buttonMap {
		log.Println("Name:", name, "Adding Pin:", pin)

		wConf.watcher.AddPin(pin)
	}

	defer close(errChan)

	go pConf.watchPins(wConf, errChan)

	return <-errChan
}
