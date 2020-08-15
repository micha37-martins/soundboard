package filechecks

import (
	"log"
	"os"
	"strings"
	"testing"
)

func TestGetFileInfos(t *testing.T) {
	tests := []struct {
		name          string
		folder        string
		expectedError string
	}{
		{
			name:          "folder is ok",
			folder:        "./testdata/mixed/",
			expectedError: "",
		},
		{
			name:          "folder is empty",
			folder:        "./testdata/nofiles/",
			expectedError: "No files in folder: ./testdata/nofiles/",
		},
		{
			name:          "path is invalid",
			folder:        "./invalid/path/",
			expectedError: "open ./invalid/path/: no such file or directory",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := getFileInfos(test.folder)
			log.Println("Returned error: ", err)
			if test.expectedError == "" {
				if err != nil {
					t.Errorf("got %q, want no error", err)
				} else {
					t.Logf("No error for test %q", test.name)
				}
			} else if (err == nil) || (err.Error() != test.expectedError) {
				t.Errorf("got %q, want %q", err, test.expectedError)
			}
		})
	}
}

func TestIsMP3(t *testing.T) {
	tests := []struct {
		name          string
		file          string
		expected      bool
		expectedError string
	}{
		{
			name:     "mp3File",
			file:     "test.mp3",
			expected: true,
		},
		{
			name:          "empty",
			file:          "",
			expected:      false,
			expectedError: "Empty filename found!",
		},
		{
			name:     "multi.mp3",
			file:     "mp3.mp3.mp3",
			expected: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mp3OrNot, err := isMP3(test.file)
			if test.expectedError == "" {
				if err != nil {
					t.Errorf("got error %v", err)
				} else if mp3OrNot == test.expected {
					t.Logf("Test %q is successful", test.name)
				}
			} else if test.expectedError == err.Error() {
				t.Logf("Test %q is successful", test.name)
			}
		})
	}
}

// helper functions for TestFilterFileList
// using GetFileInfos function to create testdata of type []os.FileInfo
func createTestFileInfos(folder string) []os.FileInfo {
	testFileInfos, err := getFileInfos(folder)
	if err != nil {
		log.Fatal("Could not create TestFileInfos", err)
	}
	return testFileInfos
}

func compareLists(expected, result []string) bool {
	for i, expectedElement := range expected {
		if strings.Compare(expectedElement, result[i]) == 0 {
			return true
		}
	}

	return false
}

func TestFilterFileList(t *testing.T) {
	tests := []struct {
		name     string
		input    []os.FileInfo
		expected []string
	}{
		{
			name:     "working file list",
			input:    createTestFileInfos("./testdata/mixed/"),
			expected: []string{"silence_1.mp3", "silence_2.mp3"},
		},
		{
			name:     "no mp3 files",
			input:    createTestFileInfos("./testdata/nomp3/"),
			expected: []string{},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			filteredList, err := filterFileList(test.input)
			if err != nil {
				t.Errorf("got error %v", err)
			} else if len(filteredList) == 0 {
				t.Logf("expected %v, got %v, for test %q", test.expected, filteredList, test.name)
			} else if compareLists(test.expected, filteredList) == false {
				t.Errorf("got %v, want %v for test %q", test.expected, filteredList, test.name)
			}
		})
	}
}

// main.go CheckFiletype("audio/mpeg")
//"Wrong filetype for: \"silence_1.mp3\".\nExpecting audio/mpeg, got \"application/octet-stream\" \n"
func TestCheckFiletype(t *testing.T) {
	tests := []struct {
		name          string
		folder        string
		fileType      string
		expectedError string
	}{
		{
			name:          "onlyMP3",
			folder:        "./testdata/onlymp3/",
			fileType:      "audio/mpeg",
			expectedError: "",
		},
		{
			name:          "noMP3",
			folder:        "./testdata/nomp3/",
			fileType:      "audio/mpeg",
			expectedError: "No mp3 files found in folder.",
		},
		{
			name:          "mixed",
			folder:        "./testdata/mixed/",
			fileType:      "audio/mpeg",
			expectedError: "",
		},
		{
			name:          "brokenmp3",
			folder:        "./testdata/brokenmp3/",
			fileType:      "audio/mpeg",
			expectedError: "Wrong filetype for: \"silence_1.mp3\".\nExpecting audio/mpeg, got \"application/octet-stream\" \n",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := CheckFiletype(test.folder, test.fileType)
			if test.expectedError == "" {
				log.Println(err)
				if err != nil {
					t.Errorf("got %q, expected %q", err, test.expectedError)
				} else {
					t.Logf("Test is ok")
				}
			} else if test.expectedError == err.Error() {
				t.Log("Test is ok")
			}
		})
	}
}

func TestFileMapper(t *testing.T) {
	tests := []struct {
		name           string
		folder         string
		buttonNr       string
		expectedOutput string
		expectedError  string
	}{
		{
			name:           "working_mapping",
			folder:         "./testdata/onlymp3/",
			buttonNr:       "01",
			expectedOutput: "",
			expectedError:  "",
		},
		{
			name:           "no matching button",
			folder:         "./testdata/onlymp3/",
			buttonNr:       "99",
			expectedOutput: "",
			expectedError:  "No file found for button nr.: 99\n",
		},
		{
			name:           "invalid folder",
			folder:         "./testdata/invalid/",
			buttonNr:       "01",
			expectedOutput: "",
			expectedError:  "open ./testdata/invalid/: no such file or directory",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			filename, err := FileMapper(test.folder, test.buttonNr)
			if test.expectedError == "" {
				log.Println(err)
				if err != nil {
					t.Errorf("got %q, expected %q", err, test.expectedError)
				} else { // TODO add equal check
					t.Logf("Button name %q", filename)
				}
			} else if test.expectedError == err.Error() {
				t.Log("Test is ok")
			}
		})
	}
}
