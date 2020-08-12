package config

import (
	"os"
)

// SoundfilesFolder creates compatible path for soundfiles
const SoundfilesFolder string = "." + string(os.PathSeparator) + "soundfiles" + string(os.PathSeparator)

// InitializeButtonMap implements the mapping of button, GPIO-pin and sound-file
// it must be set to fit your device GPIO pins
// key defines the GPIO Pin
// value defines the Name
func InitializeButtonMap() map[uint]string {
	ButtonMap := make(map[uint]string)

	ButtonMap[5] = "01"
	ButtonMap[6] = "02"
	ButtonMap[13] = "03"
	ButtonMap[19] = "04"
	ButtonMap[26] = "05"
	ButtonMap[16] = "06"
	ButtonMap[20] = "07"
	ButtonMap[21] = "08"

	return ButtonMap
}
