package config

import (
	"reflect"
	"testing"
)

func TestInitializeButtonMap(t *testing.T) {
	expected := map[uint]string{5: "01", 6: "02", 13: "03", 16: "06", 19: "04", 20: "07", 21: "08", 26: "05"}

	buttonMap := InitializeButtonMap()

	eq := reflect.DeepEqual(expected, buttonMap)
	if eq {
		t.Log("Maps are equal.")
	} else {
		t.Error("Maps are unequal.")
	}
}
