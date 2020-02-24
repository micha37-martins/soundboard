// this package is for debugging and plays every file once

package pushbuttons

import (
	"fmt"
	"soundboard/internal/config"
	"soundboard/internal/filechecks"
	"soundboard/internal/soundboard"
)

// PushedButtons creates a "watcher" listening on the specified buttons
// if a button is pushed the assigned sound-file will be played by the player
func PushedButtons() error {
	numberOfButtons := len(config.ButtonMapping)
	fmt.Println("Nr. of buttons", numberOfButtons)

	// range over list of Buttons and add pin to watcher
	for i, mapping := range config.ButtonMapping {
		fmt.Println("Cycle nr:", i, "PIN:", mapping.Pin)
		fileName, err := filechecks.FileMapper(mapping.Name)
		fmt.Println("Filename", fileName)

		soundboard.PlaySound(config.SoundfilesFolder + fileName)

		if err != nil {
			return err
		}
	}
	return nil
}
