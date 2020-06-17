// pushbuttons processes the GPIO logic and triggers playback

package pushbuttons

import (
	"errors"
	"log"
	"soundboard/internal/config"
	"soundboard/internal/filechecks"
	"soundboard/internal/soundboard"
	"time"

	"github.com/micha37-martins/gpio"
)

// Initialize ButtonMap
var buttonMap = config.InitializeButtonMap()

type watchConfig struct {
	errChan chan error
	folder  string
	player  Player
	watcher Watcher
}

// Watcher defines interface for gpio watcher
type Watcher interface {
	Watch() (uint, uint)
	AddPin(uint)
	Close()
}

// Player defines interface for music player
type Player interface {
	play(string)
}

type playerConfig struct{}

func newWatchConfig(folder string, player Player, watcher Watcher) *watchConfig {
	return &watchConfig{
		folder:  folder,
		player:  player,
		errChan: make(chan error),
		watcher: watcher,
	}
}

func (pC *playerConfig) play(path string) {
	log.Println("folder + fileName =", path)
	soundboard.PlaySound(path)
}

// find values by key and play file
// TODO function so schreiben das sie fertig ist
func (wConf *watchConfig) mapAndPlay(watcherPin uint, watcherValue uint) {
	if name, ok := buttonMap[watcherPin]; ok {
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

// checkPins contains the main logic of this package
// it checks if a button has been pushed
func (wConf *watchConfig) checkPins() {
	// declare and initialize variables
	// pullup resistor has "1" set as default value for pins
	var watcherPin uint = 0
	var watcherValue uint = 1

	watcherPin, watcherValue = wConf.watcher.Watch()

	log.Printf("Reading watcherValue %d and watcherPin: %d\n", watcherValue, watcherPin)

	wConf.mapAndPlay(watcherValue, watcherValue)
}

// TODO als nächstes error channel implementieren
// Im Prinzip error channel erstellen und in der (go) Funktion if err!=nil errorchannel <- err
// z.B errc := make(chan error, 1) vor der gofunc deklarieren und handlen

// TODO nur einmal testen mit echten pfaden ob files abgespielt werden und zwar im "player" package
func (wConf *watchConfig) WatchPins() {
	log.Println("ButtonMap: ", buttonMap)

	for {
		// slow down loop to go easy on resources
		time.Sleep(time.Second / 2)

		wConf.checkPins()
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

	wConf := newWatchConfig(folder, player, watcher)
	defer close(wConf.errChan)
	go wConf.WatchPins()

	// TODO evaluieren: Idee statt loop durch channel der fehler aufnimmt blocken
	err := <-wConf.errChan
	return err

	//time.Sleep(time.Second * 120)
}
