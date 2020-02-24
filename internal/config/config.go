package config

import (
	"os"
)

const SoundfilesFolder string = "./soundfiles" + string(os.PathSeparator)

// ButtonMap implements the mapping of button, GPIO-pin and sound-file
type ButtonMap struct {
	Name string
	Pin  uint
}

// ButtonMapping must be set to fit your device GPIO pins
var ButtonMapping = []*ButtonMap{
	&ButtonMap{
		Name: "01",
		Pin:  5,
	},
	&ButtonMap{
		Name: "02",
		Pin:  6,
	},
	&ButtonMap{
		Name: "03",
		Pin:  13,
	},
	&ButtonMap{
		Name: "04",
		Pin:  19,
	},
	&ButtonMap{
		Name: "05",
		Pin:  26,
	},
	&ButtonMap{
		Name: "06",
		Pin:  16,
	},
	&ButtonMap{
		Name: "07",
		Pin:  20,
	},
	&ButtonMap{
		Name: "08",
		Pin:  21,
	},
}
