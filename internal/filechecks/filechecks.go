package filechecks

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// getFileInfos returns list of Fileinformations for every file in a specified folder
func getFileInfos(folder string) ([]os.FileInfo, error) {
	fileInfos, err := ioutil.ReadDir(folder)
	if err != nil {
		return nil, err
	}
	if len(fileInfos) == 0 {
		noFilesErr := "No files in folder: " + folder
		return nil, errors.New(noFilesErr)
	}
	return fileInfos, nil
}

// isMP3 checks if the suffix of a file is MP3 or mp3
func isMP3(file string) (bool, error) {
	if file == "" {
		return false, errors.New("Empty filename found!")
	}

	matchMP3, err := regexp.MatchString("mp3|MP3", filepath.Ext(file))

	if err != nil {
		return false, err
	}
	return matchMP3, nil
}

// filterFileList adds files with mp3 or MP3 suffix to fileList
func filterFileList(fileList []os.FileInfo) ([]string, error) {
	filteredList := []string{}

	for _, file := range fileList {
		desiredType, err := isMP3(file.Name())
		if err != nil {
			return nil, err
		}
		if desiredType {
			filteredList = append(filteredList, file.Name())
		}
	}
	return filteredList, nil
}

// getFilenameList extracts filenames received by getFileInfos function
// only mp3 files are allowed
func getFilenameList(folder string) ([]string, error) {
	fileList, err := getFileInfos(folder)
	if err != nil {
		return nil, err
	}
	filenameList, err := filterFileList(fileList)
	if err != nil {
		return nil, err
	}
	if len(filenameList) == 0 {
		return nil, errors.New("No mp3 files found in folder.")
	} else {
		return filenameList, nil
	}
}

// CheckFiles validates that only audio/mpeg files are used
// input desired filetype p.e. "audio/mpeg"
func CheckFiletype(folder, desiredFileType string) error {
	filenames, err := getFilenameList(folder)
	if err != nil {
		return err
	}

	// take only first 512 bytes into consideration
	buff := make([]byte, 512)

	// iterate over files
	for _, filename := range filenames {
		path := folder + filename
		// open input file
		fileInput, err := os.Open(path)
		if err != nil {
			return err
		}
		// close fileInput on exit
		// use closure to collect error
		defer func() {
			if err := fileInput.Close(); err != nil {
				panic(err)
			}
		}()

		if _, err := fileInput.Read(buff); err != nil {
			return err
		}
		// validation
		fileType := http.DetectContentType(buff)

		if fileType != desiredFileType {
			return fmt.Errorf(
				"Wrong filetype for: %q.\nExpecting audio/mpeg, got %q \n",
				filename, fileType)
		}
	}
	return nil
}

// FileMapper assigns a filename to the corresponding button
// button number has to be a two digit string
// example: "01" = Button01
func FileMapper(folder, buttonNr string) (string, error) {
	fileMap, err := getFilenameList(folder)
	if err != nil {
		return "", err
	}
	mappedName := ""

	for _, fileName := range fileMap {
		if strings.Compare(buttonNr, fileName[0:2]) == 0 {
			mappedName = fileName
			return mappedName, nil
		}
	}
	return "", errors.New(fmt.Sprintln("No file found for button nr.:", buttonNr))
}
