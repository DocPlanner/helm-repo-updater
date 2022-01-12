package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// UpdateApplication get the relative path folder based in number of relative paths indicated
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
	fmt.Printf("finalPath: %s\n", finalPath)
	_, err = os.Stat(finalPath)
	if err != nil {
		fmt.Printf("Directory doesn't %s exist", finalPath)
	} else {
		fmt.Printf("Directory %s exist", finalPath)
	}
	return &finalPath, nil
}
