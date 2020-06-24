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

type playerConfig struct{}

type watcherConfig struct {
	errChan chan error
	folder  string
	player  Player
	watcher Watcher
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

// newWatcherConfig constructor for instanciating watcherConfig
func newWatcherConfig(folder string, player Player, watcher Watcher) *watcherConfig {
	w := new(watcherConfig)
	w.folder = folder
	w.player = player
	w.errChan = make(chan error, 1)
	w.watcher = watcher
	return w
}

// play calls PlaySound function start playing file
func (p *playerConfig) play(path string) {
	log.Println("folder + fileName =", path)

	soundboard.PlaySound(path)
}

// mapAndPlay finds values by key and hands path to player
func (wConf *watcherConfig) mapAndPlay(pin uint) {
	if name, ok := buttonMap[pin]; ok {
		fileName, err := filechecks.FileMapper(wConf.folder, name)

		if err != nil {
			wConf.errChan <- err
		}

		log.Printf("Folder: %s - Filename: %s", wConf.folder, fileName)

		wConf.player.play(wConf.folder + fileName)
	} else {
		wConf.errChan <- errors.New("No file mapped to pushed button")
	}
}

// checkPins checks if a button has been pushed
// pullup resistor has "1" set as default value for pins
// notice: zero value of uint is 0
func (wConf *watcherConfig) checkPins() (uint, uint) {
	var pin uint = 0
	var value uint = 1

	pin, value = wConf.watcher.Watch()

	if value != 0 && value != 1 {
		e := fmt.Sprintf("Pin value ot of scope. Only 0 or 1 allowed, got %d.", value)
		wConf.errChan <- errors.New(e)
		return 0, 0
	}

	log.Printf("Got value %d and pin: %d from watcher.Watch()\n",
		value, pin)

	return pin, value
}

// TODO als nächstes error channel testen
// z.B errc := make(chan error, 1) vor der gofunc deklarieren und handlen
// TODO nur einmal testen mit echten pfaden ob files abgespielt werden und zwar im "player" package
//WatchPins continously calls the checkPins to test if a button has been pushed
func (wConf *watcherConfig) WatchPins() {
	log.Println("ButtonMap: ", buttonMap)

	for {
		// slow down loop to go easy on resources
		time.Sleep(time.Second / 2)

		pin, value := wConf.checkPins()

		if value == 0 {
			wConf.mapAndPlay(pin)
		}
	}
}

// PushedButtons creates a "watcher" listening on the specified buttons
// if a button is pressed the assigned sound-file will be played
func PushedButtons(folder string) error {
	// Watcher is a type which listens on the GPIO pins you specify and
	// then notifies you when the values of those pins change.
	watcher := gpio.NewWatcher()
	defer watcher.Close()

	player := &playerConfig{}

	// range over list of Buttons and add pin to watcher
	for pin, name := range buttonMap {
		log.Println("Name:", name, "Adding Pin:", pin)

		watcher.AddPin(pin)
	}

	wConf := newWatcherConfig(folder, player, watcher)
	defer close(wConf.errChan)
	go wConf.WatchPins()

	// TODO evaluieren: Idee statt loop durch channel der fehler aufnimmt blocken
	// eventuell ergibt es sinn den errChan auf länge 0 zu setzen um ihn zum blcoken zu nutzen
	err := <-wConf.errChan
	return err
	// TODO replace
	//time.Sleep(time.Second * 120)
}
