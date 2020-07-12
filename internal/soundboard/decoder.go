package soundboard

import (
	"bytes"
	"encoding/binary"
	"log"

	"github.com/bobertlo/go-mpg123/mpg123"
	"github.com/gordonklaus/portaudio"
)

// PlaySound plays selected sound-file
// it is based on the example from gordonklaus/portaudio:
// https://github.com/gordonklaus/portaudio/blob/master/examples/mp3.go
func PlaySound(filePath string) {
	log.Println("Playing: ", filePath)

	// create mpg123 decoder instance
	decoder, err := mpg123.NewDecoder("")
	chk(err)

	chk(decoder.Open(filePath))
	defer decoder.Close()

	// get audio format information
	rate, channels, _ := decoder.GetFormat()

	// make sure output format does not change
	decoder.FormatNone()
	decoder.Format(rate, channels, mpg123.ENC_SIGNED_16)

	portaudio.Initialize()
	defer portaudio.Terminate()
	// default [out := make([]int16, 8192)]
	// for raspberry pi set frames per buffer to 128
	// http://www.portaudio.com/docs/v19-doxydocs/open_default_stream.html
	out := make([]int16, 128)

	stream, err := portaudio.OpenDefaultStream(0, channels, float64(rate), len(out), &out)
	chk(err)
	defer stream.Close()

	chk(stream.Start())
	defer stream.Stop()

	for {
		audio := make([]byte, 2*len(out))
		_, err = decoder.Read(audio)
		if err == mpg123.EOF {
			break
		}
		chk(err)

		chk(binary.Read(bytes.NewBuffer(audio), binary.LittleEndian, out))
		chk(stream.Write())
	}
}

func chk(err error) {
	if err != nil {
		panic(err)
	}
}
