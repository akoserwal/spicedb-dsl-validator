package util

import (
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var projectRootDirectory = GetProjectRootDir()

func GetProjectRootDir() string {
	workingDir, err := os.Getwd()
	if err != nil {
		log.Err(err)
	}
	dirs := strings.Split(workingDir, "/")
	var goModPath string
	var rootPath string
	for _, d := range dirs {
		rootPath = rootPath + "/" + d
		goModPath = rootPath + "/go.mod"
		_, err := ioutil.ReadFile(goModPath)
		if err != nil { // if the file doesn't exist, continue searching
			continue
		}
	}
	return rootPath
}

// Read the contents of file into string value
func ReadFileValueString(file string, val *string) error {
	fileContents, err := readFile(file)
	if err != nil {
		return err
	}

	*val = strings.TrimSuffix(fileContents, "\n")
	return err
}

func readFile(file string) (string, error) {
	absFilePath := BuildFullFilePath(file)

	// If no file is provided then we don't try to read it
	if absFilePath == "" {
		return "", nil
	}

	// Read the file
	buf, err := ioutil.ReadFile(absFilePath)
	if err != nil {
		return "", err
	}
	return string(buf), nil
}

func BuildFullFilePath(filename string) string {
	// If the value is in quotes, unquote it
	unquotedFile, err := strconv.Unquote(filename)
	if err != nil {
		// values without quotes will raise an error, ignore it.
		unquotedFile = filename
	}

	// If no file is provided, leave val unchanged.
	if unquotedFile == "" {
		return ""
	}

	// Ensure the absolute file path is used
	absFilePath := unquotedFile
	if !filepath.IsAbs(unquotedFile) {
		absFilePath = filepath.Join(projectRootDirectory, unquotedFile)
	}
	return absFilePath
}

func SafeString(ptr *string) string {
	if ptr == nil {
		return ""
	}
	return *ptr
}
