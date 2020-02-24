package filechecks

import (
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
			_, err := PushedButtons(test.folder)
			if test.expectedError == "" {
				if err != nil {
					t.Errorf("got %q, want no error", err)
				} else {
					t.Logf("No error for test %q", test.name)
				}
			} else if err.Error() != test.expectedError {
				t.Errorf("got %q, want %q", err, test.expectedError)
			}
		})
	}
}
