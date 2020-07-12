package pushbuttons

import (
	"errors"
	"log"
	"testing"
)

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

type testPlayerConfig struct {
	folder string
	player Player
}

type testWatcherConfig struct {
	watcher Watcher
}

func newTestPlayer() *testPlayer {
	return &testPlayer{}
}

func newTestWatcherConfig(watcher Watcher) *watcherConfig {
	return &watcherConfig{watcher: watcher}
}

func NewTestPlayerConfig(folder string, player Player) *testPlayerConfig {
	pc := new(testPlayerConfig)
	pc.folder = folder
	pc.player = player
	return pc
}

// The following testcode is inspired by code from:
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

func (tp *testPlayerConfig) watchPins(wConf *watcherConfig, errChan chan error) {
	log.Println("Mocked watchPins() function. Just returning error.")

	errChan <- errors.New("Some error occurred!")
}

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

// testsound can be found at:
// https://freesound.org/people/ramsamba/sounds/318687/
func TestCheckPins(t *testing.T) {
	watcher := newTestWatcher()
	somePin := getSomeKey(buttonMap)
	testPlayer := newTestPlayer()

	tests := []struct {
		name             string
		testNotification testNotification
		folder           string
		expectedErr      string
	}{
		{
			name:             "checkPins value=0",
			testNotification: testNotification{pin: somePin, value: 0},
			folder:           "./testdata/onlymp3/",
			expectedErr:      "",
		},
		{
			name:             "checkPins value=1",
			testNotification: testNotification{pin: somePin, value: 1},
			folder:           "./testdata/onlymp3/",
			expectedErr:      "",
		},
		{
			name:             "checkPins value=5",
			testNotification: testNotification{pin: somePin, value: 5},
			folder:           "./testdata/onlymp3/",
			expectedErr:      "Pin value ot of scope. Only 0 or 1 allowed, got 5.",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			watcher.notification <- test.testNotification
			pConf := NewPlayerConfig(test.folder, testPlayer)
			errChan := make(chan error, 1)

			pin, value := pConf.checkPins(watcher, errChan)

			t.Log("Pin: ", pin, "Value: ", value)

			resultedErr := ""

			select {
			case err, ok := <-errChan:
				if ok {
					t.Log("Error detected: ", err)
					resultedErr = err.Error() // error type to string
				} else {
					t.Log("Channel closed")
				}
			default:
				t.Log("No errors")
			}
			if resultedErr != test.expectedErr {
				t.Errorf("Got %q, expected %q", resultedErr, test.expectedErr)
			}
		})
	}
}

func TestMapAndPlay(t *testing.T) {
	somePin := getSomeKey(buttonMap)
	testPlayer := newTestPlayer()

	tests := []struct {
		name        string
		pin         uint
		folder      string
		expectedErr string
	}{
		{
			name:        "Existing Pin",
			pin:         somePin,
			folder:      "./testdata/onlymp3/",
			expectedErr: "",
		},
		{
			name:        "Pin should not exist",
			pin:         999,
			folder:      "./testdata/onlymp3/",
			expectedErr: "No file mapped to pushed button.",
		},
		{
			name:        "Folder should not exist",
			pin:         somePin,
			folder:      "./notexisting/",
			expectedErr: "open ./notexisting/: no such file or directory",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			pConf := NewPlayerConfig(test.folder, testPlayer)
			errChan := make(chan error, 1)

			pConf.mapAndPlay(errChan, test.pin)
			resultedErr := ""

			select {
			case err, ok := <-errChan:
				if ok {
					t.Log("Error detected: ", err)
					resultedErr = err.Error() // error type to string
				} else {
					t.Log("Channel closed")
				}
			default:
				t.Log("No errors")
			}
			if resultedErr != test.expectedErr {
				t.Errorf("Got %q, expected %q", resultedErr, test.expectedErr)
			}
		})
	}
}

func TestPushedButtons(t *testing.T) {
	somePin := getSomeKey(buttonMap)
	testPlayer := newTestPlayer()
	watcher := newTestWatcher()
	testWatcherConfig := NewWatcherConfig(watcher)

	tests := []struct {
		name        string
		pin         uint
		folder      string
		expectedErr string
	}{
		{
			name:        "Loop stops because of error",
			pin:         somePin,
			folder:      "./testdata/onlymp3/",
			expectedErr: "Some error occurred!",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			errChan := make(chan error, 1)
			testPlayerConfig := NewTestPlayerConfig(test.folder, testPlayer)
			var resultedErr error
			resultedErr = PushedButtons(testPlayerConfig, testWatcherConfig, errChan)

			select {
			case err, ok := <-errChan:
				if ok {
					t.Log("Error detected: ", err)
					resultedErr = err // error type to string
				} else {
					t.Log("Channel closed")
				}
			default:
				t.Log("No errors")
			}
			if resultedErr.Error() != test.expectedErr {
				t.Errorf("Got %q, expected %q", resultedErr, test.expectedErr)
			}
		})
	}
}
