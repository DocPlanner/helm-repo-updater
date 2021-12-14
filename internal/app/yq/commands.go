package yq

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/argoproj-labs/argocd-image-updater/pkg/log"
)

// InplaceApply applies the yq expression to the given file
func InplaceApply(key, value string, targetFile string) error {
	if !strings.HasPrefix(key, ".") {
		return fmt.Errorf("key %s doesn't start with '.'", key)
	}

	cmd := fmt.Sprintf("yq eval -i '%s=\"%s\"' %s", key, value, targetFile)
	command := exec.Command("/bin/sh", "-c", cmd)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	log.Debugf(command.String())

	err := command.Run()
	if err != nil {
		return fmt.Errorf("cmd.Run() failed with %s", command)
	}

	return nil
}

// ReadKey reads the value of the given key from the given file
func ReadKey(key string, targetFile string) string {
	if !strings.HasPrefix(key, ".") {
		return ""
	}

	cmd := fmt.Sprintf("yq eval '%s' %s", key, targetFile)
	command := exec.Command("/bin/sh", "-c", cmd)

	log.Debugf(command.String())

	out, err := command.CombinedOutput()
	if err != nil {
		return ""
	}

	return strings.TrimSpace(string(out))
}
