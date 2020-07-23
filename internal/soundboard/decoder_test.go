package soundboard

import (
	"bytes"
	"log"
	"os"
	"testing"
)

// TestPlaySound actually plays a Testfile
// if you are getting an error first check on the path
func TestPlaySound(t *testing.T) {
	//read log output
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer func() {
		log.SetOutput(os.Stderr)
	}()

	PlaySound("./testdata/radio.mp3")
	PlaySound("./testdata/bell.mp3")

	t.Log(buf.String())
}
