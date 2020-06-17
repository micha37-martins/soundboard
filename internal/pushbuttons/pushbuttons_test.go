package pushbuttons

import (
	"log"
	"soundboard/internal/config"
	"testing"
)

var buttons = config.InitializeButtonMap()

type testNotification struct {
	pin   uint
	value uint
}

type testWatcher struct {
	pin          uint
	closed       bool
	notification chan testNotification
}

type testPlayer struct{}

func newTestPlayer() *testPlayer {
	return &testPlayer{}
}

// The following testcode is inspired by martinohmanns code:
// https://github.com/martinohmann/rfoutlet/
func newTestWatcher() *testWatcher {
	return &testWatcher{notification: make(chan testNotification, 1)}
}

// The following three methodes implement the Watcher interface
// Watch blocks until one change occurs on one of the watched pins
// It returns the pin which changed and its new value
// Users can either use Watch() or receive from Watcher.Notification directly
func (w *testWatcher) Watch() (uint, uint) {
	notification := <-w.notification
	return notification.pin, notification.value
}

func (w *testWatcher) AddPin(pin uint) {
	w.pin = pin
}

func (w *testWatcher) Close() {
	w.closed = true
}

//todo watch pins func beenden mit error z.B.
func (tp *testPlayer) play(fileName string) {
	log.Printf("Testing not playing %s", fileName)
}

// getSomeKey gets a random key from a map
// it can be used to get a valid GPIO Pin number from ButtonMap
func getSomeKey(m map[uint]string) uint {
	for k := range m {
		return k
	}
	return 0
}

//end testcode mohmann
// es muss ein Watcher und ein folder an die Funktion übergeben werden
// testsound can be found here:
// https://freesound.org/people/ramsamba/sounds/318687/
func TestCheckPins(t *testing.T) {
	watcher := newTestWatcher()
	somePin := getSomeKey(buttons)
	testPlayer := newTestPlayer()

	tests := []struct {
		name             string
		testNotification testNotification
		folder           string
		expectedError    string
	}{
		{
			name:             "checkPinsTest01",
			testNotification: testNotification{pin: somePin, value: 0},
			folder:           "./testdata/onlymp3/",
			expectedError:    "",
		},
		{
			name:             "checkPinsTest02",
			testNotification: testNotification{pin: somePin, value: 0},
			folder:           "./testdata/onlymp3/",
			expectedError:    "",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			watcher.notification <- test.testNotification
			wConf := newWatchConfig(test.folder, testPlayer, watcher)
			wConf.checkPins()
			log.Println("Test", test.name, "went through")
			//if err.Error() != test.expectedError {
			//	t.Errorf("Got %s, expected %s", err, test.expectedError)
			//}
		})
	}
}
