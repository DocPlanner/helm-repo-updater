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
