package utils

import (
	"os"
	"path/filepath"
	"strings"
)

// GetRouteRelativePath get the relative path folder based in number of relative paths indicated
func GetRouteRelativePath(numRelativePath int, relativePath string) (*string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	parent := filepath.Dir(wd)
	s := strings.Split(parent, "/")
	s = s[:len(s)-numRelativePath]
	finalPath := strings.Join(s, "/")
	finalPath = finalPath + relativePath
	return &finalPath, nil
}

// CreateAndWriteContentInTempFile writes the string content in a temporary file
// based in a pattern and creating the file in the default tmpDir of the machine
func CreateAndWriteContentInTempFile(tempFilePattern string, content string) (*os.File, error) {
	// create and open a temporary file
	f, err := os.CreateTemp(os.TempDir(), tempFilePattern)
	if err != nil {
		return nil, err
	}

	// write data to the temporary file
	data := []byte(content)
	if _, err := f.Write(data); err != nil {
		return nil, err
	}

	return f, nil
}
